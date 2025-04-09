package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/shresthaoshan/workman/internal/models"
)

const CONFIG_KEY = "WORKMAN_CONFIG_PATH"

func LoadRuntimeConfig() *models.WorkmanConfig {

	return &models.WorkmanConfig{
		CONFIG_PATH: *getConfigPath(),
	}
}

func getConfigPath() *string {
	if val, ok := os.LookupEnv(CONFIG_KEY); ok {
		if stats, err := os.Stat(val); err == nil || os.IsExist(err) {
			if !stats.IsDir() && path.Ext(val) == "json" {
				return &val
			}
		}
	}

	app_data_path := getAppOsPath()

	err := os.MkdirAll(app_data_path, 0755)
	if err != nil {
		panic(fmt.Sprintf("app data path could not be resolved: %s", err.Error()))
	}

	default_path := path.Join(app_data_path, ".workman.config.json")

	return &default_path
}

func getAppOsPath() string {
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			panic("APPDATA environment variable not set")
		}
		return filepath.Join(appData, "workman")
	case "darwin", "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic("Failed to get user home directory")
		}
		return filepath.Join(homeDir, ".config", "workman")
	default:
		panic(fmt.Sprintf("Unsupported operating system: %s", runtime.GOOS))
	}
}
