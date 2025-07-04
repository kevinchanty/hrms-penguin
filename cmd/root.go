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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hrms-penguin",
	Short: "HRMS is hard to use",
	Long:  `Trying to write a cli`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		loggerOptions := log.Options{
			ReportCaller: true,
		}
		if enableDebugLog {
			loggerOptions.Level = log.DebugLevel
		}

		logger := log.NewWithOptions(os.Stderr, loggerOptions)

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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hrms-penguin.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVarP(&enableDebugLog, "debug", "d", false, "Enable debug logging")
}
