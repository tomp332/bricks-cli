package logout

import (
	"bricks-cli/internal/utils"
	"github.com/spf13/cobra"
)

// LogoutCommand CLI command that prints the status of the authentication state for the current user.
var LogoutCommand = &cobra.Command{
	Use:   "logout",
	Short: "Logout from the current user session",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logout()
	},
}
