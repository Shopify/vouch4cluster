package cmd

import (
	"github.com/Shopify/vouch4cluster/process"
)

var defaultConfig = &process.VoucherConfig{}

func getVoucherCfg() *process.VoucherConfig {
	return defaultConfig
}
