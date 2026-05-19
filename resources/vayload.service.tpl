[Unit]
Description={{DESCRIPTION}}
Documentation={{DOCUMENTATION}}
After=network.target
Wants=network-online.target

[Service]
Type=simple
User={{USER}}
Group={{GROUP}}
WorkingDirectory={{WORKDIR}}

ExecStart={{EXEC_START}}
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
RestartSec=3

# Environment
Environment=ENV={{ENV}}
EnvironmentFile=-{{ENV_FILE}}

# Limits
LimitNOFILE=65535
LimitNPROC=4096

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=full
ProtectHome=true
ProtectKernelTunables=true
ProtectKernelModules=true
ProtectControlGroups=true
RestrictRealtime=true
RestrictSUIDSGID=true
LockPersonality=true
MemoryDenyWriteExecute=true

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier={{NAME}}

# Timeout
TimeoutStartSec=30
TimeoutStopSec=15

[Install]
WantedBy=multi-user.target
