#!/bin/sh
set -e

REPO="diegovrocha/certui"
DEST="/usr/local/bin"

OS=$(uname -s | tr A-Z a-z)
ARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')

URL="https://github.com/${REPO}/releases/latest/download/certui_${OS}_${ARCH}.tar.gz"

echo "Installing certui..."
echo "  OS:   ${OS}"
echo "  Arch: ${ARCH}"
echo "  From: ${URL}"
echo ""

# Use sudo only when needed (skip when already root, e.g. in Docker)
if [ "$(id -u)" -eq 0 ]; then
    SUDO=""
else
    SUDO="sudo"
fi

curl -sSLf "$URL" | $SUDO tar -xz -C "$DEST" certui

echo "✔ certui installed to ${DEST}/certui"
echo "  Run: certui"
