/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"hrms-penguin/internal/cli"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// leaveCmd represents the leave command
var leaveCmd = &cobra.Command{
	Use:   "leave",
	Short: "Create Leave Application",
	Long: `Create Leave Application with startTime and endTime.
	Time Format: "2025-07-18 09:20"
	Example: hrms-penguin leave '2025-07-18 09:20' '2025-07-18 09:45'`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Optionally run one of the validators provided by cobra
		if err := cobra.ExactArgs(2)(cmd, args); err != nil {
			return err
		}

		_, err := time.Parse("2006-01-02 15:04", args[0])
		if err != nil {
			return fmt.Errorf("fails to parse startTime, format: 2006-01-02 15:04:05")
		}

		_, err = time.Parse("2006-01-02 15:04", args[1])
		if err != nil {
			return fmt.Errorf("fails to parse endTime, format: 2006-01-02 15:04:05")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Leave Command started.")
		logger, err := cli.SetupLogger(enableDebugLog, logPath)
		if err != nil {
			log.Fatal("Setup logger fails: %v", err)
		}

		startTime, err := time.Parse("2006-01-02 15:04", args[0])
		if err != nil {
			logger.Fatalf("parse startTime fails: %v", err)
		}
		endTime, err := time.Parse("2006-01-02 15:04", args[1])
		if err != nil {
			logger.Fatalf("parse entTime fails: %v", err)
		}

		hrmsClient, err := cli.SetupHrmsClient(logger, forcePromptConfig)
		if err != nil {
			logger.Fatal("Setup HRMS client fails: %v", err)
		}

		err = hrmsClient.Login()
		if err != nil {
			logger.Fatalf("Login Fails: %v", err)
		}

		err = hrmsClient.CreateLeaveApplication(startTime, endTime)
		if err != nil {
			logger.Fatalf("Create Leave Application fails: %v", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(leaveCmd)
}
