/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"hrms-penguin/internal/hrmsclient"
	"os"

	_ "embed"

	"github.com/charmbracelet/log"
	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
)

var (
	enableDebugLog bool
	logPath        string
	enableNoti     bool
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

		var config hrmsclient.HrmsConfig
		err := json.Unmarshal((configStr), &config)
		if err != nil {
			logger.Fatal("Error parsing config")
		}

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
			logger.Fatalf("GetAttendance errored: %v", err)
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
