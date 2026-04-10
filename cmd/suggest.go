package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/yagofontanez/repomind-cli/internal/api"
	"github.com/yagofontanez/repomind-cli/internal/reader"
)

var repo string
var context string

var suggestCmd = &cobra.Command{
	Use:   "suggest",
	Short: "Analyze your project and suggest features",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📦 Reading project...")

		dir, _ := os.Getwd()
		projectName := filepath.Base(dir)

		files, err := reader.ReadLocal(dir)
		if err != nil {
			fmt.Println("Error reading project:", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Found %d files\n", len(files))

		var payload []map[string]string
		for _, f := range files {
			payload = append(payload, map[string]string{
				"path":    f.Path,
				"content": f.Content,
			})
		}

		fmt.Println("🤖 Analyzing with AI...")

		result, err := api.Analyse(payload, context, projectName)
		if err != nil {
			if strings.HasPrefix(err.Error(), "rate_limit") {
				fmt.Println("❌ Daily limit reached. Upgrade your plan at https://repomind.dev")
			} else {
				fmt.Println("Error:", err)
			}
			os.Exit(1)
		}

		fmt.Print("\n💡 Suggestions:\n\n")
		for i, s := range result.Suggestions {
			fmt.Printf("%d. [%s] %s\n", i+1, s.Type, s.Title)
			fmt.Printf("   %s\n\n", s.Description)
		}

		fmt.Printf("🔗 View full analysis: %s\n", result.PanelURL)
		browser.OpenURL(result.PanelURL)
	},
}

func init() {
	suggestCmd.Flags().StringVar(&repo, "repo", "", "GitHub repository URL")
	suggestCmd.Flags().StringVar(&context, "context", "", "Product context description")
	rootCmd.AddCommand(suggestCmd)
}
