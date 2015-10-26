# consul-dynamic
Consul watcher for bootstrapping and dynamically configuring consul

# Bootstrap usage
```shell
Usage:
  consul-dynamic bootstrap [flags]

Flags:
  -c, --consul-config="/etc/consul/consul.json": Consul configuration file
  -i, --node-ip="": Consul node public IP
  -n, --node-name="": Consul node name
  -p, --peers="": Comma separated list of peers
  -s, --server[=false]: Consul server flag

Global Flags:
      --log-level="warn": Logging level
```

