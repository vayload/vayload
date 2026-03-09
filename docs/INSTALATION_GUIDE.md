# Vayload - Installation Guide

## 🚀 Quick Installation

### Automatic Installation (Recommended)

**Basic installation:**

```bash
curl -fsSL https://raw.githubusercontent.com/zsweiter/vayload/main/scripts/install.sh | bash
```

Or with wget:

```bash
wget -qO- https://raw.githubusercontent.com/zsweiter/vayload/main/scripts/install.sh | bash
```

**Full installation with system service:**

```bash
curl -fsSL https://raw.githubusercontent.com/zsweiter/vayload/main/scripts/install.sh | bash -s -- --install-service
```

**Install specific version:**

```bash
curl -fsSL https://raw.githubusercontent.com/zsweiter/vayload/main/scripts/install.sh | bash -s -- -r v1.0.0
```

**Install to custom directory:**

```bash
curl -fsSL https://raw.githubusercontent.com/zsweiter/vayload/main/scripts/install.sh | bash -s -- -d ~/bin -p ~/vayload-project
```

**Available script options:**

- `-d, --install-dir DIR` - Installation directory for binaries (default: /usr/local/bin)
- `-p, --project-dir DIR` - Project directory for data (default: /opt/vayload)
- `-r, --release VERSION` - Specific version to install (default: latest)
- `-s, --skip-setup` - Skip initial setup wizard
- `--skip-path` - Don't add vayload to PATH
- `--install-service` - Install as system service (requires sudo)
- `-h, --help` - Show help message

### Manual Installation

```bash
# Download latest version
VERSION=$(curl -s https://api.github.com/repos/zsweiter/vayload/releases/latest | grep tag_name | cut -d '"' -f 4)
wget https://github.com/vayload/vayload/releases/download/${VERSION}/vayload-${VERSION}-linux-amd64.tar.gz

# Extract
tar -xzf vayload-${VERSION}-linux-amd64.tar.gz
cd vayload-${VERSION}-linux-amd64

# Install binaries
sudo cp vayload /usr/local/bin/
sudo cp vayload-server /usr/local/bin/
sudo chmod +x /usr/local/bin/vayload /usr/local/bin/vayload-server
```

---

## 📋 Prerequisites

- Unix-like operating system (Linux recommended)
- SSH access with sudo permissions
- Internet access
- Database (PostgreSQL or MySQL recommended)

---

## 🛠️ Step-by-Step Configuration

### 1. Create project directory

```bash
sudo mkdir -p /opt/vayload
cd /opt/vayload
```

### 2. Verify installation

```bash
vayload --version
vayload --help
```

### 3. Initialize the project

```bash
sudo vayload setup
```

This command will guide you through:

- Database configuration
- Admin user creation
- JWT key generation
- Creating necessary directories (uploads, logs)

### 4. Install as system service

```bash
sudo vayload install
```

This will create a systemd service that:

- Starts automatically on system boot
- Restarts automatically on failure
- Manages system logs

### 5. Verify the service

```bash
# Check status
sudo systemctl status vayload

# View logs in real-time
sudo journalctl -u vayload -f

# Start/stop/restart
sudo systemctl start vayload
sudo systemctl stop vayload
sudo systemctl restart vayload
```

### 6. Health Check

Verify the server is running:

```bash
curl http://localhost:8080/_rest/health
```

You should receive a `200 OK` response.

---

## 🌐 Reverse Proxy Configuration

### Nginx (Recommended)

Create the file `/etc/nginx/sites-available/vayload`:

```nginx
upstream vayload_backend {
    server 127.0.0.1:8080;
    keepalive 32;
}

server {
    listen 80;
    server_name vayload.yourdomain.com;

    # Redirect to HTTPS (uncomment when you have SSL)
    # return 301 https://$server_name$request_uri;

    client_max_body_size 100M;

    # IP whitelist (optional)
    # allow 192.168.1.0/24;
    # deny all;

    location / {
        proxy_pass http://vayload_backend;
        proxy_http_version 1.1;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        proxy_set_header Connection "";
        proxy_buffering off;
        proxy_request_buffering off;
    }

    # WebSocket support (if needed)
    location /ws {
        proxy_pass http://vayload_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # gRPC support
    location /grpc {
        grpc_pass grpc://127.0.0.1:9090;
        grpc_set_header X-Real-IP $remote_addr;
    }
}
```

Enable the site:

```bash
sudo ln -s /etc/nginx/sites-available/vayload /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### SSL with Let's Encrypt (Recommended)

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d vayload.yourdomain.com
```

### Apache

Enable modules:

```bash
sudo a2enmod proxy proxy_http proxy_grpc headers ssl
```

Create `/etc/apache2/sites-available/vayload.conf`:

```apache
<VirtualHost *:80>
    ServerName vayload.yourdomain.com

    ProxyPreserveHost On
    ProxyRequests Off

    ProxyPass / http://127.0.0.1:8080/
    ProxyPassReverse / http://127.0.0.1:8080/

    # Security headers
    Header always set X-Content-Type-Options "nosniff"
    Header always set X-Frame-Options "SAMEORIGIN"

    # IP whitelist (optional)
    # <Location />
    #     Require ip 192.168.1.0/24
    # </Location>

    ErrorLog ${APACHE_LOG_DIR}/vayload-error.log
    CustomLog ${APACHE_LOG_DIR}/vayload-access.log combined
</VirtualHost>
```

Enable:

```bash
sudo a2ensite vayload
sudo systemctl reload apache2
```

---

## 🔐 Access Admin Panel

1. Open your browser and go to: `http://vayload.yourdomain.com`
2. Enter the credentials configured during `setup`
3. You're ready to use Vayload!

---

## 🔧 Advanced Configuration

### Environment Variables

Create or edit `/opt/vayload/.env`:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=vayload
DB_USER=vayload
DB_PASSWORD=your_secure_password

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# gRPC Server
GRPC_PORT=9090
GRPC_HOST=0.0.0.0

# MCP Server
MCP_PORT=9091
MCP_HOST=0.0.0.0

# JWT
JWT_SECRET=your_very_secure_secret_generated_by_setup

# Uploads
UPLOAD_MAX_SIZE=100M
UPLOAD_DIR=/opt/vayload/uploads

# Logs
LOG_LEVEL=info
LOG_DIR=/opt/vayload/logs
```

### Firewall Configuration

```bash
# Allow only Nginx/Apache to Vayload port
sudo ufw allow 'Nginx Full'
sudo ufw allow 22/tcp  # SSH
sudo ufw enable

# If using gRPC/MCP externally
sudo ufw allow 9090/tcp  # gRPC
sudo ufw allow 9091/tcp  # MCP
```

### Automatic Backups

Create a backup script at `/opt/vayload/backup.sh`:

```bash
#!/bin/bash
BACKUP_DIR="/opt/vayload/backups"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# Database backup
pg_dump -U vayload vayload > $BACKUP_DIR/db_$DATE.sql

# Uploads backup
tar -czf $BACKUP_DIR/uploads_$DATE.tar.gz /opt/vayload/uploads

# Config backup
cp /opt/vayload/.env $BACKUP_DIR/env_$DATE.bak

# Clean old backups (older than 7 days)
find $BACKUP_DIR -mtime +7 -delete

echo "Backup completed: $DATE"
```

Make it executable and add to crontab:

```bash
sudo chmod +x /opt/vayload/backup.sh
sudo crontab -e
# Add:
0 2 * * * /opt/vayload/backup.sh
```

---

## 📊 Monitoring

### System Logs

```bash
# Systemd logs
sudo journalctl -u vayload -n 100

# Application logs
sudo tail -f /opt/vayload/logs/vayload.log

# Follow logs in real-time
sudo journalctl -u vayload -f
```

### Resource Monitoring

```bash
# View resource usage
systemctl status vayload

# Detailed statistics
sudo systemd-cgtop

# Check port usage
sudo netstat -tlnp | grep vayload
# or
sudo ss -tlnp | grep vayload
```

### Performance Monitoring

```bash
# HTTP endpoint metrics
curl http://localhost:8080/_rest/metrics

# Health check
curl http://localhost:8080/_rest/health

# Server info
curl http://localhost:8080/_rest/info
```

---

## 🐛 Troubleshooting

### Service won't start

```bash
# Check permissions
sudo chown -R vayload:vayload /opt/vayload

# Validate configuration
vayload validate-config

# View detailed logs
sudo journalctl -u vayload -n 50 --no-pager

# Check if port is already in use
sudo lsof -i :8080
```

### Cannot access admin panel

1. Verify service is running: `systemctl status vayload`
2. Check firewall: `sudo ufw status`
3. Test direct access: `curl http://localhost:8080/_rest/health`
4. Check Nginx/Apache logs:
    ```bash
    sudo tail -f /var/log/nginx/error.log
    sudo tail -f /var/log/apache2/error.log
    ```

### Database connection issues

```bash
# Verify database is running
sudo systemctl status postgresql

# Test manual connection
psql -U vayload -d vayload -h localhost

# Check connection in logs
sudo journalctl -u vayload | grep -i database
```

### gRPC/MCP not working

```bash
# Check if ports are open
sudo netstat -tlnp | grep -E '9090|9091'

# Test gRPC connection
grpcurl -plaintext localhost:9090 list

# Check firewall rules
sudo ufw status | grep -E '9090|9091'
```

### Performance issues

```bash
# Check system resources
top
htop

# Monitor disk I/O
iotop

# Check disk space
df -h

# View memory usage
free -h

# Optimize database
# PostgreSQL
sudo -u postgres vacuumdb -z vayload

# MySQL
mysqlcheck -u vayload -p --optimize vayload
```

---

## ☁️ Supported Providers

- ✅ DigitalOcean Droplets
- ✅ AWS EC2
- ✅ Google Cloud Compute Engine
- ✅ Azure Virtual Machines
- ✅ Generic VPS (Linode, Vultr, etc.)
- ✅ Dedicated servers
- ⚠️ Shared hosting (with limitations)
- ⚠️ cPanel (requires manual configuration)

---

## 🔄 Updating Vayload

### Using the install script

```bash
# Update to latest version
curl -fsSL https://raw.githubusercontent.com/zsweiter/vayload/main/scripts/install.sh | bash

# Update to specific version
curl -fsSL https://raw.githubusercontent.com/zsweiter/vayload/main/scripts/install.sh | bash -s -- -r v1.2.0
```

### Manual update

```bash
# Download new version
VERSION=v1.2.0
wget https://github.com/vayload/vayload/releases/download/${VERSION}/vayload-${VERSION}-linux-amd64.tar.gz

# Stop service
sudo systemctl stop vayload

# Backup current binaries
sudo cp /usr/local/bin/vayload /usr/local/bin/vayload.bak
sudo cp /usr/local/bin/vayload-server /usr/local/bin/vayload-server.bak

# Extract and install
tar -xzf vayload-${VERSION}-linux-amd64.tar.gz
cd vayload-${VERSION}-linux-amd64
sudo cp vayload vayload-server /usr/local/bin/

# Restart service
sudo systemctl start vayload

# Verify
vayload --version
```

### Rollback

```bash
# If update fails, rollback to backup
sudo systemctl stop vayload
sudo mv /usr/local/bin/vayload.bak /usr/local/bin/vayload
sudo mv /usr/local/bin/vayload-server.bak /usr/local/bin/vayload-server
sudo systemctl start vayload
```

---

## 🔒 Security Best Practices

### SSL/TLS Configuration

Always use HTTPS in production:

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d vayload.yourdomain.com

# Auto-renewal test
sudo certbot renew --dry-run
```

### Firewall Rules

```bash
# Minimal firewall setup
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 'Nginx Full'
sudo ufw enable
```

### IP Whitelisting

In Nginx:

```nginx
# Admin panel access
location /admin {
    allow 192.168.1.0/24;
    allow 10.0.0.0/8;
    deny all;

    proxy_pass http://vayload_backend;
}
```

### Regular Updates

```bash
# Create update script
cat > /opt/vayload/update.sh << 'EOF'
#!/bin/bash
echo "Updating system packages..."
sudo apt update && sudo apt upgrade -y

echo "Checking for Vayload updates..."
CURRENT=$(vayload --version | grep -oP 'v\d+\.\d+\.\d+')
LATEST=$(curl -s https://api.github.com/repos/zsweiter/vayload/releases/latest | grep tag_name | cut -d '"' -f 4)

if [ "$CURRENT" != "$LATEST" ]; then
    echo "New version available: $LATEST (current: $CURRENT)"
    echo "Run: curl -fsSL https://raw.githubusercontent.com/zsweiter/vayload/main/scripts/install.sh | bash"
else
    echo "Vayload is up to date: $CURRENT"
fi
EOF

chmod +x /opt/vayload/update.sh
```

---

## 📚 Additional Resources

- [Full Documentation](https://github.com/vayload/vayload/wiki)
- [API Reference](https://github.com/vayload/vayload/wiki/API)
- [gRPC Protocol](https://github.com/vayload/vayload/wiki/gRPC)
- [MCP Integration](https://github.com/vayload/vayload/wiki/MCP)
- [Report Issues](https://github.com/vayload/vayload/issues)
- [Discussions](https://github.com/vayload/vayload/discussions)

---

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md).

---

## 📄 License

Vayload is released under the [MIT License](LICENSE).

---

**Need help?** Open an issue on GitHub or check the complete documentation.
