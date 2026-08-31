#!/usr/bin/env bash

set -euo pipefail

INSTALL_DIR="/usr/local/bin"
EXECUTABLE="routelocal"
TARGET="${INSTALL_DIR}/${EXECUTABLE}"

echo "Uninstalling ${EXECUTABLE}..."

if [ ! -f "${TARGET}" ]; then
    echo "Error: ${EXECUTABLE} is not installed at ${TARGET}."
    
    # Try finding it in PATH
    if command -v ${EXECUTABLE} >/dev/null 2>&1; then
        FOUND_AT=$(command -v ${EXECUTABLE})
        echo "Found ${EXECUTABLE} at ${FOUND_AT} instead."
        TARGET="${FOUND_AT}"
    else
        echo "RouteLocal does not appear to be installed on this system."
        exit 0
    fi
fi

echo "Found RouteLocal at: ${TARGET}"

# Prompt for confirmation if running interactively
if [ -t 0 ]; then
    read -p "Are you sure you want to remove ${TARGET}? (y/N) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Uninstall cancelled."
        exit 0
    fi
fi

SUDO=""
if [ ! -w "${TARGET}" ]; then
    if [ "$(id -u)" != "0" ]; then
        echo "Root privileges required to remove ${TARGET}. Prompting for sudo..."
        SUDO="sudo"
    fi
fi

$SUDO rm -f "${TARGET}"

echo "✓ RouteLocal has been successfully uninstalled."
