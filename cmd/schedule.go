package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var scheduleTime string

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Schedule daily analysis for this project",
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS == "windows" {
			fmt.Println("⚠️  Scheduling is not supported on Windows yet.")
			fmt.Println("   Consider using WSL or running 'repomind suggest' manually.")
			return
		}

		dir, _ := os.Getwd()

		parts := strings.Split(scheduleTime, ":")
		if len(parts) != 2 {
			fmt.Println("❌ Invalid time format. Use --time 08:00")
			os.Exit(1)
		}
		hour := parts[0]
		minute := parts[1]

		binary, err := exec.LookPath("repomind")
		if err != nil {
			binary, _ = os.Executable()
		}

		cronLine := fmt.Sprintf("%s %s * * * cd %s && %s suggest\n", minute, hour, dir, binary)

		out, _ := exec.Command("crontab", "-l").Output()
		current := string(out)

		if strings.Contains(current, dir+" && "+binary+" suggest") {
			fmt.Println("✅ Already scheduled for this project.")
			return
		}

		newCrontab := current + cronLine

		cronCmd := exec.Command("crontab", "-")
		cronCmd.Stdin = strings.NewReader(newCrontab)
		if err := cronCmd.Run(); err != nil {
			fmt.Println("❌ Failed to update crontab:", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Scheduled! repomind will analyze this project every day at %s\n", scheduleTime)
		fmt.Printf("   Project: %s\n", dir)
	},
}

func init() {
	scheduleCmd.Flags().StringVar(&scheduleTime, "time", "08:00", "Time to run daily analysis (HH:MM)")
	rootCmd.AddCommand(scheduleCmd)
}
