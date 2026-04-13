#!/bin/sh
set -e

REPO="diegovrocha/certool"
DEST="/usr/local/bin"

OS=$(uname -s | tr A-Z a-z)
ARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')

URL="https://github.com/${REPO}/releases/latest/download/certool_${OS}_${ARCH}.tar.gz"

echo "Installing certool..."
echo "  OS:   ${OS}"
echo "  Arch: ${ARCH}"
echo "  From: ${URL}"
echo ""

curl -sSLf "$URL" | sudo tar -xz -C "$DEST" certool

echo "✔ certool installed to ${DEST}/certool"
echo "  Run: certool"
