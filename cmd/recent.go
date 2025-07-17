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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
			logger.Fatalf("GetRedentAttendance fails: %v", err)
		}

		for _, data := range attendanceDataList {
			fmt.Printf("%+v\n", data)
		}
	},
}

func init() {
	testCmd.AddCommand(recentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// recentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// recentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
