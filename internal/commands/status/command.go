package status

import (
	"bricks-cli/internal/utils"
	"github.com/enescakir/emoji"
	"github.com/spf13/cobra"
)

// Authentication status command

// AuthStatusCommand CLI command that prints the status of the authentication state for the current user.
var AuthStatusCommand = &cobra.Command{
	Use:   "status",
	Short: "Check if current session is authenticated",
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.CheckLoginStatus() {
			utils.ErrorPrint("You are not authenticated %v", emoji.ThumbsDown)
			return
		}
		utils.SuccessPrint("You are authenticated %v", emoji.ThumbsUp)
	},
}
