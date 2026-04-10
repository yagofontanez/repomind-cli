package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init <api-key>",
	Short: "Authenticate with your repomind API key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := args[0]

		configDir, _ := os.UserConfigDir()
		dir := configDir + "/repomind"
		os.MkdirAll(dir, 0755)

		err := os.WriteFile(dir+"/config", []byte(apiKey), 0600)
		if err != nil {
			fmt.Println("Error saving API key:", err)
			os.Exit(1)
		}

		apiURL := "https://repomind-api-htom.onrender.com"
		err = os.WriteFile(dir+"/url", []byte(apiURL), 0600)
		if err != nil {
			fmt.Println("Error saving API URL:", err)
			os.Exit(1)
		}

		fmt.Println("✅ API key saved. You're ready to use repomind!")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
