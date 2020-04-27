package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/csrwng/host-router/pkg/config"
	"github.com/csrwng/host-router/pkg/controller"
)

func main() {
	log.SetLogger(zap.New(zap.UseDevMode(false)))
	hostRouterCommand().Execute()
}

func hostRouterCommand() *cobra.Command {
	r := &config.HostRouter{
		SetupFuncs: []config.SetupFunc{controller.Setup},
	}
	cmd := &cobra.Command{
		Use:   "host-router",
		Short: "Manages a hosts file based on OpenShift routes",
		Run: func(cmd *cobra.Command, args []string) {
			err := r.Start()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v", err)
				os.Exit(1)
			}
		},
	}
	cmd.Flags().StringVar(&r.Kubeconfig, "kubeconfig", os.Getenv("KUBECONFIG"), "Path to kubeconfig file to access OpenShift cluster (KUBECONFIG)")
	cmd.Flags().StringVar(&r.HostsFile, "hostsfile", os.Getenv("HOSTSFILE"), "Path to hosts file to manage (HOSTSFILE)")
	cmd.Flags().StringVar(&r.RouterAddress, "router-ip", os.Getenv("ROUTERIP"), "IP address of OpenShift router (ROUTERIP)")

	return cmd
}
