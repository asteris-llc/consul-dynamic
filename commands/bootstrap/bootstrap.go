package bootstrap

import (
	"fmt"
	"strings"

	"github.com/asteris-llc/consul-dynamic/config"
	"github.com/asteris-llc/consul-dynamic/consul"

	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
)

type bootstrapConfig struct {
	NodeName string
	NodeIp	string
	peersRaw string
	Peers	[]string
	IsServer bool
	ConsulConfig string
}

func Init(cmd *cobra.Command) {
	bc := new(bootstrapConfig)

	bootstrapCmd := &cobra.Command{
		Use: "bootstrap",
		Short: "Bootstrap a Consul node",
		Long: "Bootstrap a Consul node",
		PersistentPreRunE: func (c *cobra.Command, args[] string) error {
			if c.Parent().PersistentPreRun != nil {
				c.Parent().PersistentPreRun(c.Parent(), args)
			}
			switch {
			case bc.NodeIp == "":
				return fmt.Errorf("Node IP address not set (--node-ip)")
			case bc.peersRaw == "":
				return fmt.Errorf("Peer list not set (--peers)")
			} 

			bc.Peers = strings.Split(bc.peersRaw, ",")

			return nil
		},
		RunE: func (c *cobra.Command, args []string) error {
			return Bootstrap(bc, args)
		},
	}

	bootstrapCmd.Flags().StringVarP(&bc.NodeIp, "node-ip", "i", "", "Consul node public IP")
	bootstrapCmd.Flags().StringVarP(&bc.peersRaw, "peers", "p", "", "Comma separated list of peers")
	bootstrapCmd.Flags().BoolVarP(&bc.IsServer, "server", "s", false, "Consul server flag")
	bootstrapCmd.Flags().StringVarP(&bc.ConsulConfig, "consul-config", "c", "/etc/consul/consul.json", "Consul configuration file")
	bootstrapCmd.Flags().StringVarP(&bc.NodeName, "node-name", "n", "", "Consul node name")

	cmd.AddCommand(bootstrapCmd)
}

// Bootstrap
// Do the minimum amount of configuration to bring up a Consul datacenter
//
// Configure `ConsulConfig` with:
// {
//   advertise_addr = bc.NodeIp
//   client_addr = 0.0.0.0
//   node_name = bc.NodeName
//   server = bc.IsServer
//   retry_join = [ bc.Peers ]
//   bootstrap_expect = len(bc.Peers)
//   rejoin_after_leave = true
// }
func Bootstrap(bc *bootstrapConfig, args []string) error {
	cfg, err := config.ReadConfig(bc.ConsulConfig)
	if err != nil {
		log.Fatalf("Error reading consul config: %s", err.Error())
	}

	cfg.Data["server"] = bc.IsServer
	cfg.Data["rejoin_after_leave"] = bool(true)
	cfg.Data["advertise_addr"] = bc.NodeIp
	cfg.Data["client_addr"] = "0.0.0.0"

	cfg.Data["bootstrap_expect"] = len(bc.Peers)

	rj := make([]string, 0, len(bc.Peers))
	for _, p := range bc.Peers {
		rj = append(rj, p)
	}
	cfg.Data["retry_join"] = rj

	// If no node name provided, delete the value from consul.json. The hostname
	// will be used as the node name.
	//
	if bc.NodeName == "" {
		delete(cfg.Data, "node_name")
	} else {
		cfg.Data["node_name"] = bc.NodeName
	}

	if err := cfg.Write(bc.ConsulConfig); err != nil {
		log.Fatal(err)
	}

	if err := consul.RestartConsul(); err != nil {
		log.Fatal(err)
	}

	return nil
}
