package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/shresthaoshan/workman/internal/workman"
)

func main() {
	w := workman.Workman{}

	// Load the JSON file
	if err := w.LoadInstances(); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error loading instances: %v", err)
	}

	// Define the CLI app
	app := &cli.App{
		Name:  "workman",
		Usage: "A CLI tool to manage EC2 instances",
		Commands: []*cli.Command{
			{
				Name:  "start",
				Usage: "Start an EC2 instance",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "login",
						Aliases: []string{"l"},
						Usage:   "SSH into the instance after starting it",
					},
				},
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						return errors.New("label is required")
					}
					label := c.Args().Get(0)
					login := c.Bool("login")
					return w.StartInstance(label, login)
				},
			},
			{
				Name:  "stop",
				Usage: "Stop an EC2 instance",
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						return errors.New("label is required")
					}
					label := c.Args().Get(0)
					return w.StopInstance(label)
				},
			},
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "List all configured instances from the global config file",
				Action: func(c *cli.Context) error {
					return w.ListInstances()
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"rm"},
				Usage:   "Remove a configured instance from the global config file",
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						return errors.New("label is required")
					}
					label := c.Args().Get(0)
					return w.RemoveInstance(label)
				},
			},
			{
				Name:  "configure",
				Usage: "Interactively configure a new instance and save it to the global JSON config file",
				Action: func(c *cli.Context) error {
					var label string
					fmt.Print("Enter a label for the instance: ")
					fmt.Scanln(&label)

					if w.LabelExists(label) {
						return fmt.Errorf("label '%s' already exists. Please choose a different label", label)
					}

					var id string
					fmt.Print("Enter the instance ID: ")
					fmt.Scanln(&id)

					var instanceUser string
					fmt.Print("Enter the instance user (default: ubuntu): ")
					fmt.Scanln(&instanceUser)
					if instanceUser == "" {
						instanceUser = "ubuntu"
					}

					var pem string
					fmt.Print("Enter the path to the PEM file: ")
					fmt.Scanln(&pem)

					var awsProfile string = "default"
					fmt.Print("Enter the AWS profile (leave blank for default): ")
					fmt.Scanln(&awsProfile)
					if awsProfile == "" {
						awsProfile = "default"
					}

					var usePrivateIPInput string
					fmt.Print("Use private IP for SSH? (y/N): ")
					fmt.Scanln(&usePrivateIPInput)
					usePrivateIP := strings.ToLower(usePrivateIPInput) == "y"

					return w.ConfigureInstance(label, id, pem, awsProfile, usePrivateIP)
				},
			},
		},
	}

	// Run the CLI app
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
