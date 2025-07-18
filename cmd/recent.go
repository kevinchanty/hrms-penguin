/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"hrms-penguin/internal/cli"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// recentCmd represents the recent command
var recentCmd = &cobra.Command{
	Use:   "recent",
	Short: "Get Recent Attendance Records",
	Long:  `Get Recent Attendance Records, excluding Sat & Sun.`,
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

		attendanceDataList, err := hrmsClient.GetRecentAttendance()
		if err != nil {
			logger.Fatalf("GetRecentAttendance fails: %v", err)
		}

		for _, data := range attendanceDataList {
			fmt.Printf("%v %v %v %v", data.Date, data.OriginalInTime.Weekday().String()[:3], data.OriginalInTimeStr, data.OriginalOutTimeStr)

			if data.IsLate {
				fmt.Print(" LATE\n")
			} else {
				fmt.Print("\n")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(recentCmd)
}
