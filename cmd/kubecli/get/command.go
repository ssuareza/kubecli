package get

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// GEt is the command for get
func Command() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get <namespace> <label>",
		Short: "Get objects",
		Long:  `Get the objects for a specific namespace label.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(os.Args) < 4 {
				return fmt.Errorf("%s", "arguments wrong")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			namespace := os.Args[2]
			label := os.Args[3]
			getObjects(namespace, label)
		},
	}

	return cmd
}
