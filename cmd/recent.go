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

		for i, data := range attendanceDataList {
			if i == 0 {
				fmt.Println("--------------------------------------------------")
			}
			fmt.Printf("%v %v %v %v", data.DateStr, data.Date.Weekday().String()[:3], data.OriginalInTimeStr, data.OriginalOutTimeStr)

			// 2025-08-29 Fri 09:57 18:35 LATE
			if data.IsLate {
				fmt.Print(" LATE")
			} else {
				fmt.Print("     ")
			}

			//2025-08-29 Fri 09:57 18:35 LATE LEAVE: 09:23-09:32 APPROVED
			if data.LeaveApplicationRecord != nil {
				fmt.Printf(" LEAVE: %s-%s", data.LeaveApplicationRecord.StartTime.Format("15:04"), data.LeaveApplicationRecord.EndTime.Format("15:04"))
			}

			fmt.Print("\n")

			if (i == len(attendanceDataList)-1) || data.Date.Weekday() == 5 {
				fmt.Println("--------------------------------------------------")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(recentCmd)
}
