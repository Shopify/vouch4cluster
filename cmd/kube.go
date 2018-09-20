package cmd

import (
	"fmt"
	"os"

	"github.com/Shopify/vouch4cluster/listers/kubernetes"
	"github.com/Shopify/vouch4cluster/process"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var kubeCfgFile string

// kubeCmd represents the kube command
var kubeCmd = &cobra.Command{
	Use:   "kube",
	Short: "Attest all images in the current cluster",
	Long:  `Connects to kubernetes and attests all images in the currently configured context.`,
	Run: func(cmd *cobra.Command, args []string) {
		if "" == kubeCfgFile {
			home, err := homedir.Dir()
			if err != nil {
				errorf("failed to get default kubeconfig: %s", err)
				os.Exit(1)
			}

			kubeCfgFile = home + "/.kube/config"
		}

		var err error
		k8sLister, err := kubernetes.NewImageLister(kubeCfgFile)
		if nil != err {
			errorf("getting kubernetes image lister failed: %s", err)
			os.Exit(1)
		}

		err = process.LookupAndAttest(k8sLister, os.Stdout)
		if nil != err {
			fmt.Printf("attesting images failed: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(kubeCmd)
	kubeCmd.Flags().StringVarP(&kubeCfgFile, "kubeconfig", "k", "", "kubernetes configuration file (default is $HOME/.kube/config")
}
