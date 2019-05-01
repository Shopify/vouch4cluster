package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	k8slister "github.com/Shopify/vouch4cluster/listers/kubernetes"
	"github.com/Shopify/vouch4cluster/process"
	"github.com/spf13/cobra"

	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/tools/clientcmd"

	// add GCP authentication support to the kubernetes code
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// kubeCmd represents the kube command
var kubeCmd = &cobra.Command{
	Use:   "kube",
	Short: "Attest all images in the current cluster",
	Long:  `Connects to kubernetes and attests all images in the currently configured context.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := getVoucherCfg()

		var err error

		kubeconfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: ""}},
		).ClientConfig()
		if nil != err {
			errorf("loading Kubernetes config failed: %s", err)
			os.Exit(1)
		}

		client, err := kubernetes.NewForConfig(kubeconfig)
		if nil != err {
			errorf("initializing Kubernetes client failed: %s", err)
			os.Exit(1)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3600*time.Second)
		defer cancel()

		processor := process.NewProcessor(ctx, cfg, k8slister.NewImageLister(client), os.Stdout)

		err = processor.LookupAndAttest()
		if nil != err {
			fmt.Printf("attesting images failed: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(kubeCmd)
}
