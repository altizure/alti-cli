package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/service"
	"github.com/spf13/cobra"
)

var message string

// transferCoinsCmd transfer my coins to other user.
var transferCoinsCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer coins to others",
	Long:  "Transfer coins to others by email with custom message.",
	Run: func(cmd *cobra.Command, args []string) {
		// pre-checks general
		if err := service.Check(
			nil,
			service.CheckNonNegative(coins),
			service.CheckAPIServer(),
			service.CheckBalance(coins),
		); err != nil {
			log.Println(err)
			return
		}

		// current balance
		_, myself, err := gql.MySelf()
		if err != nil {
			log.Println(err)
			return
		}
		var ans string
		fmt.Printf("You have %.2f coins. Are you sure to transfer %.2f coins to %q? (Y/N): ", myself.Balance, coins, email)
		if assumeYes {
			fmt.Println("Yes")
		} else {
			fmt.Scanln(&ans)
			ans = strings.ToUpper(ans)
			if ans != "Y" && ans != service.Yes {
				log.Println("Cancelled.")
				return
			}
		}

		_, err = gql.TransferCoins(coins, email, message)
		if err != nil {
			log.Println(err)
			return
		}

		// balnce after transaction
		_, myself, err = gql.MySelf()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("Successfully transferred %.2f coins to %q\nCurrent balnce: %.2f coins", coins, email, myself.Balance)
	},
}

func init() {
	bankCmd.AddCommand(transferCoinsCmd)
	transferCoinsCmd.Flags().Float64VarP(&coins, "coins", "c", coins, "Number of coins to transfer")
	transferCoinsCmd.Flags().StringVarP(&email, "email", "e", email, "Recipient email")
	transferCoinsCmd.Flags().StringVarP(&message, "message", "m", message, "Message to recipient")
	transferCoinsCmd.Flags().BoolVarP(&assumeYes, "assumeyes", "y", assumeYes, "Assume yes; assume that the answer to any question which would be asked is yes")
	errors.Must(transferCoinsCmd.MarkFlagRequired("coins"))
	errors.Must(transferCoinsCmd.MarkFlagRequired("email"))
}
