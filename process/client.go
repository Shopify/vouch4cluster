package process

import (
	"context"
	"time"

	"github.com/Shopify/voucher/client"
)

// VoucherConfig is a structure which contains voucher authentication information.
type VoucherConfig struct {
	Hostname string   `json:"hostname"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Workers  int      `json:"workers"`
	Checks   []string `json:"checks"`
}

// newVoucherClient creates a new voucher.Client with the information passed
// in the passed VoucherConfig.
func newVoucherClient(ctx context.Context, cfg *VoucherConfig) (*client.VoucherClient, error) {
	var timeout time.Duration = 120 * time.Second

	deadline, hasDeadline := ctx.Deadline()
	if hasDeadline {
		timeout = time.Until(deadline)
	}

	client, err := client.NewClient(cfg.Hostname, timeout)
	if nil == err {
		client.SetBasicAuth(cfg.Username, cfg.Password)
	}

	return client, err
}
