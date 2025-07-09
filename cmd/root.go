/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"hrms-penguin/internal/hrmsclient"
	"io/fs"
	"os"
	"path"

	_ "embed"

	"github.com/charmbracelet/log"
	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

const (
	configName     string = ".hrms-penguin"
	keyringService string = "hrms-penguin"
)

var (
	enableDebugLog    bool
	logPath           string
	enableNoti        bool
	ErrConfigNotFound error = errors.New("secret not found in keyring")
)

//go:embed hrms-config.json
var configStr []byte

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hrms-penguin",
	Short: "HRMS is hard to use",
	Long:  `Trying to write a cli`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// Enable debug log
		loggerOptions := log.Options{
			ReportCaller: true,
		}
		if enableDebugLog {
			loggerOptions.Level = log.DebugLevel
		}

		// Enable log to file if logPath is set
		var logger *log.Logger
		if logPath != "" {
			logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
			if err != nil {
				log.Fatal("Error opening log file: %v", err)
			}
			defer logFile.Close()

			logger = log.NewWithOptions(logFile, loggerOptions)
		} else {
			logger = log.NewWithOptions(os.Stderr, loggerOptions)
		}

		logger.Debug("Root Command started.")

		config, err := getSavedConfig()
		if err != nil {
			if errors.Is(err, ErrConfigNotFound) {
				logger.Fatalf("config not found!", err)

			} else {
				logger.Fatalf("get saved config fails: %v\n", err)
			}
		}

		// var config hrmsclient.HrmsConfig
		// err := json.Unmarshal((configStr), &config)
		// if err != nil {
		// 	logger.Fatal("Error parsing config")
		// }

		// config, err := getSavedConfig()
		// if err != nil {
		// 	logger.Fatalf("Get config fails: %v\n", err)
		// }

		hrmsClient := hrmsclient.New(hrmsclient.ClientOption{
			HrmsConfig: config,
			Logger:     logger,
		})

		err = hrmsClient.Login()
		if err != nil {
			logger.Fatalf("Login Fails: %v", err)
		}

		todayAttendance, err := hrmsClient.GetTodayAttendance()
		if err != nil {
			logger.Fatalf("GetAttendance fails: %v", err)
		}

		if enableNoti {
			beeep.Notify("HRMS Penguin", fmt.Sprintf("Today's Attendance: %v %v %v\n", todayAttendance.Date, todayAttendance.OriginalInTime, todayAttendance.OriginalOutTime), "")
		} else {
			fmt.Printf("Today's Attendance: %v %v %v\n", todayAttendance.Date, todayAttendance.OriginalInTime, todayAttendance.OriginalOutTime)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&enableDebugLog, "debug", "d", false, "Enable debug logging")
	rootCmd.PersistentFlags().StringVarP(&logPath, "log", "l", "", "Log file path")
	rootCmd.PersistentFlags().BoolVarP(&enableNoti, "noti", "n", false, "Create push notification instead of print to stdOut")
}

func getSavedConfig() (hrmsclient.HrmsConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return hrmsclient.HrmsConfig{}, err
	}

	configPath := path.Join(homeDir, configName)

	_, err = os.Stat(configPath)
	if errors.Is(err, fs.ErrNotExist) {
		return hrmsclient.HrmsConfig{}, errors.New("saved config does not exist")
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return hrmsclient.HrmsConfig{}, err
	}

	var hrmsConfig hrmsclient.HrmsConfig
	err = json.Unmarshal(configBytes, &hrmsConfig)
	if err != nil {
		return hrmsclient.HrmsConfig{}, err
	}

	password, err := getKeyringPassword(hrmsConfig.UserName)
	if err != nil {
		return hrmsclient.HrmsConfig{}, err
	}

	hrmsConfig.Pwd = password

	return hrmsConfig, nil
}

func getKeyringPassword(user string) (string, error) {
	val, err := keyring.Get(keyringService, user)
	if err != nil {
		return "", err
	}
	return val, nil
}

func promptConfig() (hrmsclient.HrmsConfig, error) {
	log.Fatal("not done")
	return hrmsclient.HrmsConfig{}, nil
}
