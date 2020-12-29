package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// RootCmd returns the root application cmd
func RootCmd() *cobra.Command {
	ctx := context.Background()

	cmd := &cobra.Command{
		Use:   "guardian",
		Short: "Gaurdian is an API gateway that guards your upstream servers.",
	}

	cmd.AddCommand(startServerCmd(ctx))
	return cmd
}
