package ckill

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// Command is the command for ckill
func Command() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "ckill <namespace> <label>",
		Short: "Kill CrashLoopBackOff pods",
		Long:  `Kill CrashLoopBackOff pods for a specific namespace label.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(os.Args) < 4 {
				return fmt.Errorf("%s", "arguments wrong")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.SetFlags(0)
			namespace := os.Args[2]
			label := os.Args[3]

			kill(namespace, label)

		},
	}

	return cmd
}
