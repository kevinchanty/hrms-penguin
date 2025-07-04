/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"hrms-penguin/internal/hrmsclient"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var enableDebugLog bool
var logPath string

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
			logger = log.NewWithOptions(os.Stdout, loggerOptions)
		}

		logger.Debug("Root Command started.")

		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		hrmsHost := os.Getenv("HRMS_HOST")
		hrmsUserName := os.Getenv("HRMS_USER")
		hrmsPwd := os.Getenv("HRMS_PWD")

		hrmsClient := hrmsclient.New(hrmsclient.ClientOption{
			Host:     hrmsHost,
			UserName: hrmsUserName,
			Pwd:      hrmsPwd,
			Logger:   logger,
		})

		hrmsClient.Login()

		_, err = hrmsClient.GetAttendance("2025", "7")
		if err != nil {
			log.Fatalf("GetAttendance errored: %v", err)
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
}
