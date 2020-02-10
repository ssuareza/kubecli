package main

import (
	"github.com/ansel1/merry"
	"github.com/spf13/cobra"
	"github.com/ssuareza/kubecli/cmd/kubecli/config"
	"github.com/ssuareza/kubecli/cmd/kubecli/logs"
	"github.com/ssuareza/kubecli/cmd/kubecli/exec"
	"github.com/ssuareza/kubecli/cmd/kubecli/ckill"
)

var versionString = "dev"

func main() {
	cmd := &cobra.Command{
		Use:     "kubecli",
		Short:   config.Style.Title(`kube CLI`),
		Version: versionString,
		PreRun: func(cmd *cobra.Command, args []string) {
			merry.SetStackCaptureEnabled(config.Config.Debug)
		},
	}

	cmd.PersistentFlags().BoolVar(&config.Config.Debug, "debug", false, "Enable debug output")
	cmd.AddCommand(logs.Command())
	cmd.AddCommand(exec.Command())
	cmd.AddCommand(ckill.Command())

	cmd.Execute()
}
