#!/bin/sh
#
# This script downloads and installs the latest binary release of gitx.
# It detects the user's OS and architecture to download the correct binary.
#
# Usage:
# curl -sSL https://raw.githubusercontent.com/gitxtui/gitx/master/install.sh | bash

set -e

# The GitHub repository in the format "owner/repo".
REPO="gitxtui/gitx"
INSTALL_DIR="/usr/local/bin"

# Get the operating system.
get_os() {
  case "$(uname -s)" in
    Linux*)   OS='linux';;
    Darwin*)  OS='darwin';;
    *)
      echo "Unsupported operating system: $(uname -s)"
      exit 1
      ;;
  esac
  echo "$OS"
}

# Get the architecture.
get_arch() {
  case "$(uname -m)" in
    x86_64|amd64) ARCH='amd64';;
    aarch64|arm64) ARCH='arm64';;
    *)
      echo "Unsupported architecture: $(uname -m)"
      exit 1
      ;;
  esac
  echo "$ARCH"
}

# Get the latest release tag from the GitHub API.
get_latest_release() {
  curl --silent "https://api.github.com/repos/$REPO/releases/latest" |
  grep '"tag_name":' |
  sed -E 's/.*"([^"]+)".*/\1/'
}

main() {
  OS=$(get_os)
  ARCH=$(get_arch)
  VERSION=$(get_latest_release)

  if [ -z "$VERSION" ]; then
    echo "Error: Could not find the latest release version for $REPO."
    exit 1
  fi

  # Construct the archive filename and download URL.
  VERSION_NUM=$(echo "$VERSION" | sed 's/v//')
  FILENAME="gitx_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
  DOWNLOAD_URL="https://github.com/$REPO/releases/download/${VERSION}/${FILENAME}"

  # Download and extract the binary.
  echo "Downloading gitx ${VERSION} for ${OS}/${ARCH}..."
  TEMP_DIR=$(mktemp -d)
  # Download to a temporary directory.
  curl -sSL -o "$TEMP_DIR/$FILENAME" "$DOWNLOAD_URL"

  echo "Installing gitx to ${INSTALL_DIR}..."
  # Extract the archive.
  tar -xzf "$TEMP_DIR/$FILENAME" -C "$TEMP_DIR"

  # Move the binary to the installation directory.
  # Use sudo if the directory is not writable by the current user.
  if [ -w "$INSTALL_DIR" ]; then
    mv "$TEMP_DIR/gitx" "${INSTALL_DIR}/gitx"
  else
    echo "Root permission is required to install gitx to ${INSTALL_DIR}"
    sudo mv "$TEMP_DIR/gitx" "${INSTALL_DIR}/gitx"
  fi

  # Clean up the temporary directory.
  rm -rf "$TEMP_DIR"

  echo ""
  echo "gitx has been installed successfully!"
  echo "Run 'gitx' to get started."
}

# Run the main function.
main