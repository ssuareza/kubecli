package get

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Command is the command for get
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
			log.SetFlags(0)

			// stdin
			namespace := os.Args[2]
			label := os.Args[3]

			// kubeconfig
			kubeconfig := os.Getenv("HOME") + "/.kube/config"
			cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
			if err != nil {
				log.Fatal(err)
			}

			// clientset
			clientset, err := kubernetes.NewForConfig(cfg)
			if err != nil {
				log.Fatal(err)
			}

			results, err := getObjects(clientset, namespace, label)
			fmt.Println(results)
		},
	}

	return cmd
}
