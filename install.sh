#!/bin/sh
set -e

usage() {
  echo "Usage: $0 [-b bindir] [-d]"
  echo "  -b sets the directory for the binary installation, default is ./bin"
  echo "  -d enables debug mode"
  exit 1
}

# Default parameters
BINDIR=/usr/bin
DEBUG=0

# Parse command-line options
while getopts "b:d" opt; do
  case ${opt} in
    b )
      BINDIR=$OPTARG
      ;;
    d )
      set -x
      DEBUG=1
      ;;
    \? )
      usage
      ;;
  esac
done

echo "Installing to ${BINDIR}"
mkdir -p "${BINDIR}"

# Detect platform and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case $ARCH in
  x86_64)
    ARCH=amd64
    ;;
  aarch64)
    ARCH=arm64
    ;;
  *arm*)
    ARCH=arm
    ;;
  *)
    echo "Architecture $ARCH is not supported"
    exit 1
    ;;
esac

# Fetch the latest release tag from GitHub
REPO="anxuanzi/oneuptime-infrastructure-agent-go"
API_URL="https://api.github.com/repos/${REPO}/releases/latest"
TAG=$(curl -s ${API_URL} | grep '"tag_name":' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/')

if [ "$TAG" = "" ]; then
  echo "Failed to find the latest release. Please check your internet connection or GitHub API limits."
  exit 1
fi

echo "Fetching the latest release: $TAG"

# Construct the URL for the binary release
URL="https://github.com/${REPO}/releases/download/${TAG}/oneuptime-infrastructure-agent_${OS}_${ARCH}.tar.gz"

# Download and extract the binary
curl -sL "${URL}" | tar xz -C "${BINDIR}"

# Check if the binary is executable
if [ ! -x "${BINDIR}/oneuptime-infrastructure-agent" ]; then
  echo "Failed to install oneuptime-infrastructure-agent"
  exit 1
fi

echo "oneuptime-infrastructure-agent installed successfully to ${BINDIR}"