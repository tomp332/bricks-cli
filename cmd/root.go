package cmd

import (
	"bricks-cli/internal/commands/login"
	"bricks-cli/internal/commands/logout"
	"bricks-cli/internal/commands/status"
	"bricks-cli/internal/utils"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "bricks-cli",
	Short: "Bricks CLI",
	Long:  utils.MainArt,
	Run: func(cmd *cobra.Command, args []string) {
		color.Cyan(utils.MainArt)
	},
}

func init() {
	rootCmd.AddCommand(login.LoginCommand)
	rootCmd.AddCommand(status.AuthStatusCommand)
	rootCmd.AddCommand(logout.LogoutCommand)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
