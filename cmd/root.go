/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"hrms-penguin/internal/cli"
	"os"

	_ "embed"

	"github.com/charmbracelet/log"
	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
)

const (
	configName     string = ".hrms-penguin"
	keyringService string = "hrms-penguin"
)

var (
	enableDebugLog    bool
	logPath           string
	enableNoti        bool
	forcePromptConfig bool

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
		log.Debug("Root Command started.")

		logger, err := cli.SetupLogger(enableDebugLog, logPath)
		if err != nil {
			log.Fatal("Setup logger fails: %v", err)
		}

		hrmsClient, err := cli.SetupHrmsClient(logger, forcePromptConfig)
		if err != nil {
			logger.Fatal("Setup HRMS client fails: %v", err)
		}

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
	rootCmd.PersistentFlags().BoolVarP(&forcePromptConfig, "prompt", "p", false, "Ignore saved config and prompt for new one")
}
