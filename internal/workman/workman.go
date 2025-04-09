package workman

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/shresthaoshan/workman/internal/config"
	"github.com/shresthaoshan/workman/internal/models"
)

type Workman struct {
	instances map[string]models.InstanceDetails
}

// loadInstances loads the JSON file containing instance details
func (w *Workman) LoadInstances() error {
	file, err := os.ReadFile(config.LoadRuntimeConfig().CONFIG_PATH)
	if err != nil {
		if os.IsNotExist(err) {
			// Initialize an empty map if the file doesn't exist
			w.instances = make(map[string]models.InstanceDetails)
			return nil
		}
		return fmt.Errorf("failed to read file: %v", err)
	}

	if err := json.Unmarshal(file, &w.instances); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	return nil
}

// saveInstances saves the instance details to the JSON file
func (w *Workman) SaveInstances() error {
	data, err := json.MarshalIndent(w.instances, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	filePath := config.LoadRuntimeConfig().CONFIG_PATH

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("cannot write to the config file: %s, error: %s", filePath, err.Error())
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// updateLastAccessed updates the LastAccessed field for the given label
func (w *Workman) UpdateLastAccessed(label string) {
	instance := w.instances[label]
	instance.LastAccessed = time.Now().Format(time.RFC3339)
	w.instances[label] = instance

	// Save the updated instances to the JSON file
	if err := w.SaveInstances(); err != nil {
		log.Printf("Failed to update LastAccessed: %v", err)
	}
}

// listInstances lists all configured instances from the global config file
func (w *Workman) ListInstances() error {
	if len(w.instances) == 0 {
		fmt.Println("No instances are currently configured.")
		return nil
	}

	fmt.Println("Configured Instances:")
	for label, details := range w.instances {
		fmt.Printf("\nLabel: %s\n", label)
		fmt.Printf(" - AWS Profile: %s\n", details.AwsProfile)
		fmt.Printf(" - Instance ID: %s\n", details.ID)
		fmt.Printf(" - PEM File: %s\n", details.PEM)
		fmt.Printf(" - Use Private IP: %v\n", details.UsePrivateIP)

		if details.LastAccessed == "" {
			fmt.Println(" - Last Accessed: ---Configured only. Never accessed.---")
		} else {
			fmt.Printf(" - Last Accessed: %s\n", details.LastAccessed)
		}
	}
	return nil
}

// removeInstance removes a configured instance by label
func (w *Workman) RemoveInstance(label string) error {
	if _, exists := w.instances[label]; !exists {
		return fmt.Errorf("instance with label '%s' not found", label)
	}

	fmt.Printf("Are you sure you want to remove the instance '%s'? (y/N): ", label)
	var confirm string
	fmt.Scanln(&confirm)

	if strings.ToLower(confirm) != "y" {
		fmt.Println("Operation canceled.")
		return nil
	}

	delete(w.instances, label)

	// Save the updated instances to the JSON file
	if err := w.SaveInstances(); err != nil {
		return fmt.Errorf("failed to save instance details: %v", err)
	}

	fmt.Printf("Instance '%s' has been successfully removed.\n", label)
	return nil
}

// configureInstance interactively prompts the user for instance details and saves them to the JSON file
func (w *Workman) ConfigureInstance(label, id, pemPath, awsProfile string, usePrivateIP bool) error {

	// Save the instance details
	w.instances[label] = models.InstanceDetails{
		ID:           id,
		PEM:          pemPath,
		AwsProfile:   awsProfile,
		UsePrivateIP: usePrivateIP,
		// LastAccessed is intentionally left empty for newly configured instances
	}

	// Write the updated instances to the JSON file
	if err := w.SaveInstances(); err != nil {
		return fmt.Errorf("failed to save instance details: %v", err)
	}

	fmt.Printf("Instance '%s' has been successfully configured and saved.\n", label)
	return nil
}

func (w *Workman) LabelExists(label string) bool {
	_, ok := w.instances[label]
	return ok
}

// startInstance starts an EC2 instance and optionally SSHes into it
func (w *Workman) StartInstance(label string, login bool) error {
	instance, err := w.getInstanceDetails(label)
	if err != nil {
		return err
	}

	state, err := w.getInstanceState(instance.ID, instance.AwsProfile)
	if err != nil {
		return err
	}

	if state == "running" {
		fmt.Printf("Instance '%s' is already running.\n", instance.ID)
		w.UpdateLastAccessed(label) // Update LastAccessed on access
		if login {
			return w.sshIntoInstance(instance.ID, instance.PEM, instance.AwsProfile, instance.UsePrivateIP)
		}
		return nil
	}

	fmt.Printf("Starting instance '%s'...\n", instance.ID)

	sess, err := w.createAWSSession(instance.AwsProfile)

	if err != nil {
		return err
	}

	svc := ec2.New(sess)

	_, err = svc.StartInstances(&ec2.StartInstancesInput{
		InstanceIds: []*string{aws.String(instance.ID)},
	})
	if err != nil {
		return fmt.Errorf("failed to start instance: %v", err)
	}

	fmt.Printf("Waiting for instance '%s' to start...\n", instance.ID)

	err = svc.WaitUntilInstanceRunning(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instance.ID)},
	})
	if err != nil {
		return fmt.Errorf("failed while waiting for instance to start: %v", err)
	}

	fmt.Printf("Instance '%s' is now running.\n", instance.ID)
	w.UpdateLastAccessed(label) // Update LastAccessed on access

	if login {
		return w.sshIntoInstance(instance.ID, instance.PEM, instance.AwsProfile, instance.UsePrivateIP)
	}

	return nil
}

// stopInstance stops an EC2 instance
func (w *Workman) StopInstance(label string) error {
	instance, err := w.getInstanceDetails(label)
	if err != nil {
		return err
	}

	state, err := w.getInstanceState(instance.ID, instance.AwsProfile)
	if err != nil {
		return err
	}

	if state == "stopped" {
		fmt.Printf("Instance '%s' is already stopped.\n", instance.ID)
		w.UpdateLastAccessed(label) // Update LastAccessed on access
		return nil
	}

	if state == "stopping" {
		fmt.Printf("Instance '%s' is currently stopping. Waiting for it to stop...\n", instance.ID)

		sess, err := w.createAWSSession(instance.AwsProfile)

		if err != nil {
			return err
		}

		svc := ec2.New(sess)

		err = svc.WaitUntilInstanceStopped(&ec2.DescribeInstancesInput{
			InstanceIds: []*string{aws.String(instance.ID)},
		})

		if err != nil {
			return fmt.Errorf("failed while waiting for instance to stop: %v", err)
		}

		fmt.Printf("Instance '%s' has been stopped.\n", instance.ID)
		w.UpdateLastAccessed(label) // Update LastAccessed on access
		return nil
	}

	fmt.Printf("Stopping instance '%s'...\n", instance.ID)
	sess, err := w.createAWSSession(instance.AwsProfile)

	if err != nil {
		return err
	}

	svc := ec2.New(sess)

	_, err = svc.StopInstances(&ec2.StopInstancesInput{
		InstanceIds: []*string{aws.String(instance.ID)},
	})
	if err != nil {
		return fmt.Errorf("failed to stop instance: %v", err)
	}

	fmt.Printf("Waiting for instance '%s' to stop...\n", instance.ID)

	err = svc.WaitUntilInstanceStopped(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instance.ID)},
	})
	if err != nil {
		return fmt.Errorf("failed while waiting for instance to stop: %v", err)
	}

	fmt.Printf("Instance '%s' has been stopped.\n", instance.ID)
	w.UpdateLastAccessed(label) // Update LastAccessed on access

	return nil
}

// getInstanceDetails retrieves instance details by label
func (w *Workman) getInstanceDetails(label string) (models.InstanceDetails, error) {
	instance, exists := w.instances[label]
	if !exists {
		return models.InstanceDetails{}, fmt.Errorf("label '%s' not found", label)
	}
	return instance, nil
}

// createAWSSession creates an AWS session using the specified profile
func (w *Workman) createAWSSession(profile string) (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String("us-west-2")}, // Replace with your desired region
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}
	return sess, nil
}

// getInstanceState retrieves the current state of an EC2 instance
func (w *Workman) getInstanceState(instanceID, profile string) (string, error) {
	sess, err := w.createAWSSession(profile)

	if err != nil {
		return "", err
	}

	svc := ec2.New(sess)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		return "", fmt.Errorf("failed to describe instance: %v", err)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return "", fmt.Errorf("instance '%s' not found", instanceID)
	}

	return *result.Reservations[0].Instances[0].State.Name, nil
}

// sshIntoInstance connects to an EC2 instance via SSH
func (w *Workman) sshIntoInstance(instanceID, pemFile, profile string, usePrivateIP bool) error {
	sess, err := w.createAWSSession(profile)

	if err != nil {
		return err
	}

	svc := ec2.New(sess)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		return fmt.Errorf("failed to describe instance: %v", err)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return fmt.Errorf("instance '%s' not found", instanceID)
	}

	instance := result.Reservations[0].Instances[0]

	var ip string
	if usePrivateIP || instance.PublicIpAddress == nil {
		ip = *instance.PrivateIpAddress
		fmt.Println("Using private IP for SSH connection.")
	} else {
		ip = *instance.PublicIpAddress
	}

	if ip == "" {
		return fmt.Errorf("could not retrieve an IP address for instance '%s'", instanceID)
	}

	cmd := exec.Command("ssh", "-i", pemFile, fmt.Sprintf("ec2-user@%s", ip))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("Connecting to instance '%s' via SSH...\n", instanceID)
	return cmd.Run()
}
