package logs

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func Command() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "logs <namespace> <label>",
		Short: "Get service logs",
		Long:  `Print the logs for a specific namespace label.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(os.Args) < 4 {
				return fmt.Errorf("%s", "arguments wrong")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			namespace := os.Args[2]
			label := os.Args[3]
			tailLogs(namespace, label)
		},
	}

	return cmd
}
