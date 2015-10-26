package commands

import (
	"strings"

	"github.com/asteris-llc/consul-dynamic/commands/bootstrap"

	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
)


func Run() {
	cmd := Init()
	if err := cmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

func Init() *cobra.Command {
	var logLevel string

	rval := &cobra.Command{
		Use: "consul-dynamic",
		Short: "Dynamic configuration for Consul",
		Long: "Dynamic configuration for Consul",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ConfigureLogging(logLevel)
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Fatal("No command specified (try consul-dynamic help)")
		},
	}

	rval.PersistentFlags().StringVarP(&logLevel, "log-level", "", "warn", "Logging level")

	bootstrap.Init(rval)

	return rval
}

func ConfigureLogging(level string) {
	l, err := log.ParseLevel(strings.ToLower(level))
	if err != nil {
		log.SetLevel(log.WarnLevel)
		log.Warnf("Invalid log level '%v'. Setting to WARN")
	} else {
		log.SetLevel(l)
	}
}
