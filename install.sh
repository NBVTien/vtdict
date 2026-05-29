#!/bin/sh
set -e

REPO="NBVTien/vtdict"
BIN="vtdict"
INSTALL_DIR="/usr/local/bin"

# detect OS and arch
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

ASSET="${BIN}_${OS}_${ARCH}"
URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"

echo "Installing vtdict..."
echo "Downloading ${ASSET}..."

if command -v curl >/dev/null 2>&1; then
  curl -fsSL "$URL" -o "/tmp/${BIN}"
elif command -v wget >/dev/null 2>&1; then
  wget -qO "/tmp/${BIN}" "$URL"
else
  echo "Error: curl or wget required" >&2
  exit 1
fi

chmod +x "/tmp/${BIN}"

# install to /usr/local/bin, fallback to ~/bin if no permission
if [ -w "$INSTALL_DIR" ]; then
  mv "/tmp/${BIN}" "${INSTALL_DIR}/${BIN}"
else
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
  mv "/tmp/${BIN}" "${INSTALL_DIR}/${BIN}"
  case ":$PATH:" in
    *":$INSTALL_DIR:"*) ;;
    *)
      echo ""
      echo "  Add to your shell config:"
      echo "    export PATH=\"\$HOME/.local/bin:\$PATH\""
      ;;
  esac
fi

echo "Done. Run: vtdict hello"
