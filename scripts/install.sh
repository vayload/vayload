#!/bin/bash

set -e

REPO="zsweiter/vayload"
RELEASE="latest"
OS="$(uname -s)"
SKIP_SETUP="false"
INSTALL_SERVICE="false"
SKIP_PATH="false"

# Normalize OS name
case "${OS}" in
   MINGW* | Win*) OS="Windows" ;;
esac

# Default installation directory
if [ -n "$VAYLOAD_DIR" ]; then
  INSTALL_DIR="$VAYLOAD_DIR"
elif [ -n "$XDG_DATA_HOME" ]; then
  INSTALL_DIR="$XDG_DATA_HOME/vayload"
elif [ "$OS" = "Darwin" ]; then
  INSTALL_DIR="$HOME/Library/Application Support/vayload"
else
  INSTALL_DIR="/usr/local/bin"
fi

# Project directory (for server, config, uploads, logs)
if [ "$OS" = "Windows" ]; then
  PROJECT_DIR="$HOME/vayload"
else
  PROJECT_DIR="/opt/vayload"
fi

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Parse command line arguments
parse_args() {
  while [[ $# -gt 0 ]]; do
    key="$1"

    case $key in
    -d | --install-dir)
      INSTALL_DIR="$2"
      shift
      shift
      ;;
    -p | --project-dir)
      PROJECT_DIR="$2"
      shift
      shift
      ;;
    -s | --skip-setup)
      SKIP_SETUP="true"
      shift
      ;;
    --skip-path)
      SKIP_PATH="true"
      shift
      ;;
    --install-service)
      INSTALL_SERVICE="true"
      shift
      ;;
    -r | --release)
      RELEASE="$2"
      shift
      shift
      ;;
    -h | --help)
      show_help
      exit 0
      ;;
    *)
      echo -e "${RED}Unrecognized argument: $key${NC}"
      echo "Use --help for usage information"
      exit 1
      ;;
    esac
  done
}

show_help() {
  cat << EOF
Vayload Installation Script

Usage:
  curl -fsSL https://raw.githubusercontent.com/$REPO/main/scripts/install.sh | bash

  Or with options:
  curl -fsSL https://raw.githubusercontent.com/$REPO/main/scripts/install.sh | bash -s -- [OPTIONS]

Options:
  -d, --install-dir DIR      Installation directory for binaries (default: /usr/local/bin)
  -p, --project-dir DIR      Project directory for data (default: /opt/vayload)
  -r, --release VERSION      Specific release version (default: latest)
  -s, --skip-setup           Skip the initial setup wizard
  --skip-path                Skip adding vayload to PATH
  --install-service          Install as system service (requires sudo)
  -h, --help                 Show this help message

Examples:
  # Basic installation
  curl -fsSL https://raw.githubusercontent.com/$REPO/main/scripts/install.sh | bash

  # Install specific version
  curl -fsSL https://raw.githubusercontent.com/$REPO/main/scripts/install.sh | bash -s -- -r v1.0.0

  # Custom installation directory
  curl -fsSL https://raw.githubusercontent.com/$REPO/main/scripts/install.sh | bash -s -- -d ~/bin

  # Install and setup as service
  curl -fsSL https://raw.githubusercontent.com/$REPO/main/scripts/install.sh | bash -s -- --install-service

EOF
}

set_filename() {
  if [ "$OS" = "Linux" ]; then
    case "$(uname -m)" in
      arm | armv7*)
        ARCH="armv7"
        ;;
      aarch64 | arm64)
        ARCH="arm64"
        ;;
      x86_64)
        ARCH="amd64"
        ;;
      *)
        echo -e "${RED}Unsupported architecture: $(uname -m)${NC}"
        exit 1
        ;;
    esac
    FILENAME="vayload-${OS,,}-$ARCH"
    ARCHIVE_EXT="tar.gz"
  elif [ "$OS" = "Darwin" ]; then
    case "$(uname -m)" in
      arm64)
        ARCH="arm64"
        ;;
      x86_64)
        ARCH="amd64"
        ;;
      *)
        echo -e "${RED}Unsupported architecture: $(uname -m)${NC}"
        exit 1
        ;;
    esac
    FILENAME="vayload-darwin-$ARCH"
    ARCHIVE_EXT="tar.gz"
  elif [ "$OS" = "Windows" ]; then
    case "$(uname -m)" in
      x86_64)
        ARCH="amd64"
        ;;
      arm64)
        ARCH="arm64"
        ;;
      *)
        echo -e "${RED}Unsupported architecture: $(uname -m)${NC}"
        exit 1
        ;;
    esac
    FILENAME="vayload-windows-$ARCH"
    ARCHIVE_EXT="zip"
  else
    echo -e "${RED}OS $OS is not supported.${NC}"
    echo "If you think that's a bug, please file an issue at https://github.com/$REPO/issues"
    exit 1
  fi
}

check_dependencies() {
  echo -e "${BLUE}Checking dependencies...${NC}"
  local SHOULD_EXIT="false"

  echo -n "  Checking curl... "
  if hash curl 2>/dev/null; then
    echo -e "${GREEN}✓${NC}"
  else
    echo -e "${RED}✗${NC}"
    SHOULD_EXIT="true"
  fi

  if [ "$ARCHIVE_EXT" = "tar.gz" ]; then
    echo -n "  Checking tar... "
    if hash tar 2>/dev/null; then
      echo -e "${GREEN}✓${NC}"
    else
      echo -e "${RED}✗${NC}"
      SHOULD_EXIT="true"
    fi
  else
    echo -n "  Checking unzip... "
    if hash unzip 2>/dev/null; then
      echo -e "${GREEN}✓${NC}"
    else
      echo -e "${RED}✗${NC}"
      SHOULD_EXIT="true"
    fi
  fi

  if [ "$SHOULD_EXIT" = "true" ]; then
    echo -e "${RED}Missing required dependencies. Please install them and try again.${NC}"
    exit 1
  fi
}

get_latest_release() {
  if [ "$RELEASE" = "latest" ]; then
    echo -e "${BLUE}Fetching latest release...${NC}"
    RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$RELEASE" ]; then
      echo -e "${RED}Failed to fetch latest release version${NC}"
      exit 1
    fi

    echo -e "  Latest version: ${GREEN}$RELEASE${NC}"
  fi
}

download_vayload() {
  local URL

  if [ "$RELEASE" = "latest" ]; then
    URL="https://github.com/$REPO/releases/latest/download/$FILENAME.$ARCHIVE_EXT"
  else
    URL="https://github.com/$REPO/releases/download/$RELEASE/$FILENAME.$ARCHIVE_EXT"
  fi

  local DOWNLOAD_DIR=$(mktemp -d)

  echo -e "${BLUE}Downloading Vayload...${NC}"
  echo "  URL: $URL"
  echo "  Platform: $OS/$ARCH"

  if ! curl --progress-bar --fail -L "$URL" -o "$DOWNLOAD_DIR/$FILENAME.$ARCHIVE_EXT"; then
    echo -e "${RED}Download failed. Check that the release/platform are correct.${NC}"
    echo "  Release: $RELEASE"
    echo "  Filename: $FILENAME.$ARCHIVE_EXT"
    exit 1
  fi

  echo -e "${BLUE}Extracting archive...${NC}"

  if [ "$ARCHIVE_EXT" = "tar.gz" ]; then
    tar -xzf "$DOWNLOAD_DIR/$FILENAME.$ARCHIVE_EXT" -C "$DOWNLOAD_DIR"
  else
    unzip -q "$DOWNLOAD_DIR/$FILENAME.$ARCHIVE_EXT" -d "$DOWNLOAD_DIR"
  fi

  echo -e "${BLUE}Installing binaries...${NC}"

  # Create installation directory if it doesn't exist
  if [ ! -d "$INSTALL_DIR" ]; then
    echo "  Creating directory: $INSTALL_DIR"
    if [ -w "$(dirname "$INSTALL_DIR")" ]; then
      mkdir -p "$INSTALL_DIR"
    else
      sudo mkdir -p "$INSTALL_DIR"
    fi
  fi

  # Determine the extracted directory structure
  if [ -d "$DOWNLOAD_DIR/$FILENAME" ]; then
    EXTRACT_DIR="$DOWNLOAD_DIR/$FILENAME"
  else
    EXTRACT_DIR="$DOWNLOAD_DIR"
  fi

  # Install CLI binary
  if [ -f "$EXTRACT_DIR/vayload" ] || [ -f "$EXTRACT_DIR/vayload.exe" ]; then
    local VAYLOAD_BIN="vayload"
    [ "$OS" = "Windows" ] && VAYLOAD_BIN="vayload.exe"

    echo "  Installing vayload CLI to $INSTALL_DIR/"
    if [ -w "$INSTALL_DIR" ]; then
      cp "$EXTRACT_DIR/$VAYLOAD_BIN" "$INSTALL_DIR/"
      chmod +x "$INSTALL_DIR/$VAYLOAD_BIN"
    else
      sudo cp "$EXTRACT_DIR/$VAYLOAD_BIN" "$INSTALL_DIR/"
      sudo chmod +x "$INSTALL_DIR/$VAYLOAD_BIN"
    fi
  fi

  # Install Server binary
  if [ -f "$EXTRACT_DIR/vayload-server" ] || [ -f "$EXTRACT_DIR/vayload-server.exe" ]; then
    local SERVER_BIN="vayload-server"
    [ "$OS" = "Windows" ] && SERVER_BIN="vayload-server.exe"

    echo "  Installing vayload-server to $INSTALL_DIR/"
    if [ -w "$INSTALL_DIR" ]; then
      cp "$EXTRACT_DIR/$SERVER_BIN" "$INSTALL_DIR/"
      chmod +x "$INSTALL_DIR/$SERVER_BIN"
    else
      sudo cp "$EXTRACT_DIR/$SERVER_BIN" "$INSTALL_DIR/"
      sudo chmod +x "$INSTALL_DIR/$SERVER_BIN"
    fi
  fi

  # Clean up
  rm -rf "$DOWNLOAD_DIR"
}

ensure_dir_exists() {
  local DIR="$1"
  if [ ! -d "$DIR" ]; then
    echo "  Creating directory: $DIR"
    mkdir -p "$DIR"
  fi
}

setup_path() {
  if [ "$SKIP_PATH" = "true" ]; then
    return
  fi

  if [ "$INSTALL_DIR" = "/usr/local/bin" ] || [ "$INSTALL_DIR" = "/usr/bin" ]; then
    # Standard directories, already in PATH
    return
  fi

  local CURRENT_SHELL="$(basename "$SHELL")"
  local CONF_FILE=""

  echo -e "${BLUE}Setting up PATH...${NC}"

  if [ "$CURRENT_SHELL" = "zsh" ]; then
    CONF_FILE="${ZDOTDIR:-$HOME}/.zshrc"
  elif [ "$CURRENT_SHELL" = "bash" ]; then
    if [ "$OS" = "Darwin" ]; then
      CONF_FILE="$HOME/.bash_profile"
      [ ! -f "$CONF_FILE" ] && CONF_FILE="$HOME/.profile"
    else
      CONF_FILE="$HOME/.bashrc"
    fi
  elif [ "$CURRENT_SHELL" = "fish" ]; then
    CONF_FILE="$HOME/.config/fish/conf.d/vayload.fish"
    ensure_dir_exists "$(dirname "$CONF_FILE")"
  else
    echo -e "${YELLOW}Could not detect shell. Please add $INSTALL_DIR to your PATH manually.${NC}"
    return
  fi

  # Check if already in PATH config
  if grep -q "vayload" "$CONF_FILE" 2>/dev/null; then
    echo "  PATH already configured in $CONF_FILE"
    return
  fi

  ensure_dir_exists "$(dirname "$CONF_FILE")"

  echo "  Adding Vayload to PATH in $CONF_FILE"

  if [ "$CURRENT_SHELL" = "fish" ]; then
    cat >> "$CONF_FILE" << EOF

# Vayload
set -gx VAYLOAD_PATH "$INSTALL_DIR"
if test -d "\$VAYLOAD_PATH"
  set -gx PATH "\$VAYLOAD_PATH" \$PATH
end
EOF
  else
    cat >> "$CONF_FILE" << EOF

# Vayload
export VAYLOAD_PATH="$INSTALL_DIR"
if [ -d "\$VAYLOAD_PATH" ]; then
  export PATH="\$VAYLOAD_PATH:\$PATH"
fi
EOF
  fi

  echo -e "${GREEN}✓ PATH configured${NC}"
  echo -e "${YELLOW}  Run: source $CONF_FILE${NC}"
  echo -e "${YELLOW}  Or open a new terminal${NC}"
}

run_setup() {
  if [ "$SKIP_SETUP" = "true" ]; then
    return
  fi

  echo ""
  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo -e "${BLUE}  Running Vayload Setup${NC}"
  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo ""
  echo -e "${YELLOW}This will configure your Vayload installation.${NC}"
  echo ""

  # Create project directory
  if [ ! -d "$PROJECT_DIR" ]; then
    echo -e "${BLUE}Creating project directory: $PROJECT_DIR${NC}"
    if [ -w "$(dirname "$PROJECT_DIR")" ]; then
      mkdir -p "$PROJECT_DIR"
      cd "$PROJECT_DIR"
    else
      sudo mkdir -p "$PROJECT_DIR"
      sudo chown $USER:$USER "$PROJECT_DIR" 2>/dev/null || sudo chown $USER "$PROJECT_DIR"
      cd "$PROJECT_DIR"
    fi
  else
    cd "$PROJECT_DIR"
  fi

  # Run setup command
  if command -v vayload >/dev/null 2>&1; then
    vayload setup
  elif [ -x "$INSTALL_DIR/vayload" ]; then
    "$INSTALL_DIR/vayload" setup
  else
    echo -e "${RED}Could not find vayload binary${NC}"
    return 1
  fi
}

install_service() {
  if [ "$INSTALL_SERVICE" != "true" ]; then
    return
  fi

  if [ "$OS" = "Windows" ]; then
    echo -e "${YELLOW}Service installation on Windows is not yet supported.${NC}"
    echo -e "${YELLOW}Please run vayload-server manually or use Task Scheduler.${NC}"
    return
  fi

  echo ""
  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo -e "${BLUE}  Installing Vayload as System Service${NC}"
  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo ""

  cd "$PROJECT_DIR"

  if command -v vayload >/dev/null 2>&1; then
    sudo vayload install
  elif [ -x "$INSTALL_DIR/vayload" ]; then
    sudo "$INSTALL_DIR/vayload" install
  else
    echo -e "${RED}Could not find vayload binary${NC}"
    return 1
  fi

  echo -e "${GREEN}✓ Service installed${NC}"
  echo ""
  echo "Useful commands:"
  echo "  sudo systemctl status vayload   - Check service status"
  echo "  sudo systemctl start vayload    - Start service"
  echo "  sudo systemctl stop vayload     - Stop service"
  echo "  sudo systemctl restart vayload  - Restart service"
  echo "  sudo journalctl -u vayload -f   - View logs"
}

verify_installation() {
  echo ""
  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo -e "${BLUE}  Verifying Installation${NC}"
  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo ""

  local VAYLOAD_CMD=""

  if command -v vayload >/dev/null 2>&1; then
    VAYLOAD_CMD="vayload"
  elif [ -x "$INSTALL_DIR/vayload" ]; then
    VAYLOAD_CMD="$INSTALL_DIR/vayload"
  fi

  if [ -n "$VAYLOAD_CMD" ]; then
    echo -e "${GREEN}✓ Vayload CLI installed successfully${NC}"
    VERSION=$("$VAYLOAD_CMD" --version 2>/dev/null || echo "unknown")
    echo "  Version: $VERSION"
    echo "  Location: $(which vayload 2>/dev/null || echo "$INSTALL_DIR/vayload")"
  else
    echo -e "${RED}✗ Vayload CLI not found in PATH${NC}"
    echo -e "${YELLOW}  You may need to restart your terminal or run: source ~/.bashrc${NC}"
  fi

  if [ -f "$INSTALL_DIR/vayload-server" ] || [ -f "$INSTALL_DIR/vayload-server.exe" ]; then
    echo -e "${GREEN}✓ Vayload Server installed successfully${NC}"
    echo "  Location: $INSTALL_DIR/vayload-server"
  fi

  echo ""
}

show_next_steps() {
  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo -e "${BLUE}  Next Steps${NC}"
  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo ""

  if [ "$SKIP_SETUP" = "true" ]; then
    echo -e "${YELLOW}1.${NC} Initialize your project:"
    echo "   cd $PROJECT_DIR"
    echo "   vayload setup"
    echo ""
  fi

  if [ "$INSTALL_SERVICE" != "true" ]; then
    echo -e "${YELLOW}2.${NC} Install as a system service (optional):"
    echo "   sudo vayload install"
    echo ""
  fi

  echo -e "${YELLOW}3.${NC} Start the server:"
  if [ "$INSTALL_SERVICE" = "true" ]; then
    echo "   sudo systemctl start vayload"
  else
    echo "   cd $PROJECT_DIR"
    echo "   vayload-server"
  fi
  echo ""

  echo -e "${YELLOW}4.${NC} Access the admin panel:"
  echo "   http://localhost:8080"
  echo ""

  echo "📚 Documentation: https://github.com/$REPO/blob/main/GUIDE.md"
  echo "🐛 Issues: https://github.com/$REPO/issues"
  echo ""
  echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo -e "${GREEN}  Installation Complete! 🎉${NC}"
  echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

main() {
  echo -e "${GREEN}"
  cat << "EOF"
   __     __          _                 _
   \ \   / /_ _ _   _| | ___   __ _  __| |
    \ \ / / _` | | | | |/ _ \ / _` |/ _` |
     \ V / (_| | |_| | | (_) | (_| | (_| |
      \_/ \__,_|\__, |_|\___/ \__,_|\__,_|  0.1.0
                |___/
EOF
  echo -e "${NC}"
  echo -e "${BLUE}Vayload Installation Script${NC}"
  echo ""

  parse_args "$@"
  set_filename
  check_dependencies
  get_latest_release
  download_vayload
  setup_path
  verify_installation
  run_setup
  install_service
  show_next_steps
}

main "$@"
