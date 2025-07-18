package cli

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

	"github.com/charmbracelet/log"
	"github.com/zalando/go-keyring"
	"golang.org/x/term"
)

const (
	configName     string = ".hrms-penguin"
	keyringService string = "hrms-penguin"
)

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

func SetupLogger(enableDebugLog bool, logPath string) (*log.Logger, error) {
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
			return nil, err
		}
		defer logFile.Close()

		logger = log.NewWithOptions(logFile, loggerOptions)
	} else {
		logger = log.NewWithOptions(os.Stderr, loggerOptions)
	}

	return logger, nil
}

func SetupHrmsClient(logger *log.Logger, forcePromptConfig bool) (*hrmsclient.HrmsClient, error) {
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

	return hrmsClient, nil
}
