package exec

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "exec <namespace> <label>",
		Short: "Open a terminal inside a container",
		Long:  `Open a terminal in a specific namespace label.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(os.Args) < 4 {
				return fmt.Errorf("%s", "arguments wrong")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			namespace := os.Args[2]
			label := os.Args[3]
			exec(namespace, label)
		},
	}

	return cmd
}
