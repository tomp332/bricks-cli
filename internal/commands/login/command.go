package login

import (
	"bricks-cli/internal/utils"
	"github.com/enescakir/emoji"
	"github.com/spf13/cobra"
)

// LoginCommand CLI command that authenticates the user through the Github 0Auth application configured.
var LoginCommand = &cobra.Command{
	Use:   "login",
	Short: "Username using your browser",
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.CheckLoginStatus() {
			err := utils.PerformLogin()
			if err != nil {
				utils.ErrorPrint("Authentication failed %s %v", err, emoji.ThumbsDown)
			}
		} else {
			utils.SuccessPrint("You have already been authenticated %v", emoji.ThumbsUp)
		}
	},
}
