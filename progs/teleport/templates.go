package teleport

import "text/template"

var confTmpl = template.Must(template.New("teleport.yaml").
	Parse(`# General configurations.
teleport:
  data_dir: /var/lib/teleport
  auth_token: %AUTH_TOKEN%
  advertise_ip: {{ .AdvertiseIP }}
  auth_servers: %AUTH_SERVERS%
  log:
    output: stderr
    severity: INFO
  storage:
    type: dir

# Node service specific configurations.
ssh_service:
  enabled: yes
  listen_addr: 0.0.0.0:3022
  pam:
    enabled: yes
    service_name: "teleport"
`))

var serviceTmpl = template.Must(template.New("teleport-node.service").
	Parse(`[Unit]
Description=Teleport node
Wants=network-online.target
After=network-online.target
StartLimitIntervalSec=600s
ConditionPathExists={{ .ConfigFile }}

[Service]
Slice=system.slice
Type=simple
# SSH processes should be maintained even if main process is killed
KillMode=process
Restart=on-failure
RestartSec=10s
ExecStart=/usr/local/bin/teleport start --roles=node -c {{ .ConfigFile }}

[Install]
WantedBy=multi-user.target
`))
