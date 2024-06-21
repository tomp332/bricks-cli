package login

import (
	"bricks-cli/internal/utils"
	"fmt"
	"github.com/enescakir/emoji"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	clientId string
)

func checkPrompts() {
	if clientId != "" {
		// Prompt for the client-secret only if client-id is provided
		prompt := promptui.Prompt{
			Label: "Client Secret",
			Mask:  '*',
		}
		clientSecret, err := prompt.Run()
		if err != nil {
			fmt.Println("Prompt failed:", err)
			return
		}
		if clientId != "" && clientSecret != "" {
			utils.Settings.GitHubAuthConfig.GitHubClientID = clientId
			utils.Settings.GitHubAuthConfig.GitHubClientSecret = clientSecret
		}
	}
}

// LoginCommand CLI command that authenticates the user through the Github 0Auth application configured.
var LoginCommand = &cobra.Command{
	Use:   "login",
	Short: "Username using your browser",
	Run: func(cmd *cobra.Command, args []string) {
		checkPrompts()
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

func init() {
	LoginCommand.Flags().StringVarP(&clientId, "client-id", "", "", "Github client ID for login, client secret will be prompted.")
}
