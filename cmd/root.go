/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"hrms-penguin/internal/hrmsclient"
	"io/fs"
	"os"
	"path"
	"strings"

	_ "embed"

	"github.com/charmbracelet/log"
	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"golang.org/x/term"
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
		var err error
		if forcePromptConfig {
			config, err = promptConfig()
		} else {
			config, err = getSavedConfig()
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					fmt.Printf("No config found.")
				}
				config, err = promptConfig()
				if err != nil {
					logger.Fatalf("Failed to prompt for configuration: %v", err)
				}
			}
		}

		hrmsClient := hrmsclient.New(hrmsclient.NewClientOption{
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
	rootCmd.PersistentFlags().BoolVarP(&forcePromptConfig, "prompt", "p", false, "Ignore saved config and prompt for new one")
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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return hrmsclient.HrmsConfig{}, err
	}

	configPath := path.Join(homeDir, configName)

	// Prompt user for HRMS host
	fmt.Print("Enter HRMS host (e.g., https://hrms.example.com): ")
	reader := bufio.NewReader(os.Stdin)
	host, err := reader.ReadString('\n')
	if err != nil {
		return hrmsclient.HrmsConfig{}, fmt.Errorf("failed to read host: %w", err)
	}
	host = strings.TrimSpace(host)
	if host == "" {
		return hrmsclient.HrmsConfig{}, errors.New("host cannot be empty")
	}

	// Prompt user for HRMS username
	fmt.Print("Enter HRMS username: ")
	var username string
	username, err = reader.ReadString('\n')
	if err != nil {
		return hrmsclient.HrmsConfig{}, fmt.Errorf("failed to read username: %w", err)
	}
	username = strings.TrimSpace(username)
	if username == "" {
		return hrmsclient.HrmsConfig{}, errors.New("username cannot be empty")
	}

	// Prompt user for password (without echoing)
	fmt.Print("Enter HRMS password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return hrmsclient.HrmsConfig{}, fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println() // Add newline after password input
	password := strings.TrimSpace(string(passwordBytes))
	if password == "" {
		return hrmsclient.HrmsConfig{}, errors.New("password cannot be empty")
	}

	// Create config object
	config := hrmsclient.HrmsConfig{
		Host:     host,
		UserName: username,
		Pwd:      password, // This will be replaced with keyring password
	}

	// Save to keyring
	err = keyring.Set(keyringService, username, password)
	if err != nil {
		return hrmsclient.HrmsConfig{}, fmt.Errorf("failed to save password to keyring: %w", err)
	}

	// Create config for file (without password)
	configForFile := hrmsclient.HrmsConfig{
		Host:     host,
		UserName: username,
		Pwd:      "", // Don't save password in file
	}

	configData, err := json.MarshalIndent(configForFile, "", "  ")
	if err != nil {
		return hrmsclient.HrmsConfig{}, fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(configPath, configData, 0600) // Read/write for owner only
	if err != nil {
		return hrmsclient.HrmsConfig{}, fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("Configuration saved to %s\n", configPath)
	return config, nil
}
