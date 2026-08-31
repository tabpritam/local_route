#!/usr/bin/env bash

set -euo pipefail

# Constants
REPO="tabpritam/local_route"
INSTALL_DIR="/usr/local/bin"
EXECUTABLE="routelocal"
DOWNLOAD_URL_BASE="https://github.com/${REPO}/releases/latest/download"

echo "Installing ${EXECUTABLE}..."

# OS Detection
OS="$(uname -s)"
case "${OS}" in
    Linux*)     OS="linux";;
    Darwin*)    OS="darwin";;
    *)          echo "Error: Unsupported operating system: ${OS}"; exit 1;;
esac

# Architecture Detection
ARCH="$(uname -m)"
case "${ARCH}" in
    x86_64)     ARCH="amd64";;
    arm64)      ARCH="arm64";;
    aarch64)    ARCH="arm64";;
    *)          echo "Error: Unsupported architecture: ${ARCH}"; exit 1;;
esac

BINARY_NAME="${EXECUTABLE}-${OS}-${ARCH}"
BINARY_URL="${DOWNLOAD_URL_BASE}/${BINARY_NAME}"
CHECKSUM_URL="${DOWNLOAD_URL_BASE}/checksums.txt"

# Temporary directory for secure download
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TMP_DIR}"' EXIT

cd "${TMP_DIR}"

echo "Downloading latest release for ${OS}-${ARCH}..."
curl -fsSLO "${BINARY_URL}"
curl -fsSLO "${CHECKSUM_URL}"

# Verify checksum
echo "Verifying checksum..."
EXPECTED_CHECKSUM=$(grep "${BINARY_NAME}" checksums.txt | awk '{print $1}')

if [ -z "${EXPECTED_CHECKSUM}" ]; then
    echo "Error: Checksum for ${BINARY_NAME} not found in checksums.txt"
    exit 1
fi

if command -v shasum >/dev/null 2>&1; then
    # macOS / Linux with perl
    ACTUAL_CHECKSUM=$(shasum -a 256 "${BINARY_NAME}" | awk '{print $1}')
elif command -v sha256sum >/dev/null 2>&1; then
    # Linux standard
    ACTUAL_CHECKSUM=$(sha256sum "${BINARY_NAME}" | awk '{print $1}')
else
    echo "Warning: Neither shasum nor sha256sum found. Skipping checksum verification."
    ACTUAL_CHECKSUM="${EXPECTED_CHECKSUM}"
fi

if [ "${EXPECTED_CHECKSUM}" != "${ACTUAL_CHECKSUM}" ]; then
    echo "Error: Checksum verification failed!"
    echo "Expected: ${EXPECTED_CHECKSUM}"
    echo "Actual:   ${ACTUAL_CHECKSUM}"
    exit 1
fi
echo "Checksum verified successfully."

# Prepare executable
chmod +x "${BINARY_NAME}"

# Install
echo "Installing to ${INSTALL_DIR}..."
if [ ! -d "${INSTALL_DIR}" ]; then
    echo "Directory ${INSTALL_DIR} does not exist. Creating it..."
    SUDO=""
    if [ "$(id -u)" != "0" ]; then
        echo "Root privileges required to create ${INSTALL_DIR}. Prompting for sudo..."
        SUDO="sudo"
    fi
    $SUDO mkdir -p "${INSTALL_DIR}"
fi

# Check if we have write permission to INSTALL_DIR
SUDO=""
if [ ! -w "${INSTALL_DIR}" ]; then
    if [ "$(id -u)" != "0" ]; then
        echo "Root privileges required to install to ${INSTALL_DIR}. Prompting for sudo..."
        SUDO="sudo"
    fi
fi

$SUDO mv "${BINARY_NAME}" "${INSTALL_DIR}/${EXECUTABLE}"

# Verify installation
echo "Installation complete!"
if command -v ${EXECUTABLE} >/dev/null 2>&1; then
    INSTALLED_VERSION=$(${EXECUTABLE} --version 2>/dev/null || echo "Unknown (version flag unsupported)")
    echo "Successfully installed RouteLocal: ${INSTALLED_VERSION}"
else
    echo "Successfully installed RouteLocal to ${INSTALL_DIR}/${EXECUTABLE}"
    echo "Warning: ${INSTALL_DIR} is not in your PATH. You may need to add it to use the command globally."
fi

echo ""
echo "Get started by running:"
echo "  routelocal --help"
