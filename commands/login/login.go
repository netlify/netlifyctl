package login

import (
	"os"

	"github.com/netlify/netlifyctl/auth"
	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/open-api/go/models"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

const defaultClientID = "5edad8f69d47ae8923d0cf0b4ab95ba1415e67492b5af26ad97f4709160bb31b"

func Setup(middlewares []middleware.Middleware) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log user in",
		Long:  "Log user in",
	}

	return middleware.SetupCommand(cmd, loginUser, middlewares)
}

func loginUser(ctx context.Context, cmd *cobra.Command, args []string) error {
	client := context.GetClient(ctx)

	clientID := os.Getenv("NETLIFY_CLIENT_ID")
	if clientID == "" {
		clientID = defaultClientID
	}
	ticket, err := client.CreateTicket(ctx, clientID)
	if err != nil {
		return err
	}

	if err := openAuthUI(ticket); err != nil {
		return err
	}

	if !ticket.Authorized {
		a, err := client.WaitUntilTicketAuthorized(ctx, ticket)
		if err != nil {
			return err
		}

		ticket = a
	}

	token, err := client.ExchangeTicket(ctx, ticket.ID)
	if err != nil {
		return err
	}

	return auth.SaveToken(token.AccessToken)
}

func openAuthUI(ticket *models.Ticket) error {
	return open.Run("https://app.netlify.com/authorize?response_type=ticket&ticket=" + ticket.ID)
}
