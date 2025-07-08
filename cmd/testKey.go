/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/user"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

// testKeyCmd represents the testKey command
var testKeyCmd = &cobra.Command{
	Use:   "testKey",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("testKey called")

		service := "hrms-penguin"
		// user := "anon"
		user, err := user.Current()
		if err != nil {
			log.Fatal("get user fails: %v", err)
		}
		password := "secret"

		// set password
		err = keyring.Set(service, user.Name, password)
		if err != nil {
			log.Fatal(err)
		}

		// val, err := keyring.Get(service, user.Name)

		fmt.Printf("ok!")
		// fmt.Printf("ok! %v", val)
	},
}

func init() {

	rootCmd.AddCommand(testKeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testKeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testKeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
