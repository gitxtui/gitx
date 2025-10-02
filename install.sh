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
  API_URL="https://api.github.com/repos/$REPO/releases"
  
  # Use python if available, otherwise fall back to perl.
  if command -v python >/dev/null 2>&1; then
    # This python script now checks if the response is a list and not empty.
    curl -sL "$API_URL" | python -c "import sys, json; data = json.load(sys.stdin); print(data[0]['tag_name'] if isinstance(data, list) and data else '')"
  elif command -v perl >/dev/null 2>&1; then
    curl -sL "$API_URL" | perl -ne 'if (/\"tag_name\":\s*\"([^\"]+)\"/) { print $1; exit }'
  else
    # Fallback to grep/sed for systems without python/perl
    curl -sL "$API_URL" |
    grep '"tag_name":' |
    head -n 1 |
    sed -E 's/.*"([^"]+)".*/\1/'
  fi
}

main() {
  OS=$(get_os)
  ARCH=$(get_arch)
  VERSION=$(get_latest_release)

  if [ -z "$VERSION" ]; then
    echo "Error: Could not find any release version for $REPO."
    echo "Please check that the repository has releases and that you are not being rate-limited by the GitHub API."
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

  # Ensure the installation directory exists.
  if [ ! -d "$INSTALL_DIR" ]; then
    echo "Creating ${INSTALL_DIR}..."
    if sudo mkdir -p "$INSTALL_DIR" 2>/dev/null; then
      echo "Directory created successfully."
    else
      echo "Warning: Could not create ${INSTALL_DIR}."
      # Fall back to user's local bin directory.
      INSTALL_DIR="$HOME/.local/bin"
      echo "Falling back to ${INSTALL_DIR}..."
      mkdir -p "$INSTALL_DIR"
    fi
  fi

  # Move the binary to the installation directory.
  # Use sudo if the directory is not writable by the current user.
  if [ -w "$INSTALL_DIR" ]; then
    mv "$TEMP_DIR/gitx" "${INSTALL_DIR}/gitx"
  else
    echo "Root permission is required to install gitx to ${INSTALL_DIR}"
    if ! sudo mv "$TEMP_DIR/gitx" "${INSTALL_DIR}/gitx" 2>/dev/null; then
      echo "Warning: Could not install to ${INSTALL_DIR}."
      # Fall back to user's local bin directory.
      INSTALL_DIR="$HOME/.local/bin"
      echo "Falling back to ${INSTALL_DIR}..."
      mkdir -p "$INSTALL_DIR"
      mv "$TEMP_DIR/gitx" "${INSTALL_DIR}/gitx"
    fi
  fi

  # Make the binary executable.
  chmod +x "${INSTALL_DIR}/gitx"

  # Clean up the temporary directory.
  rm -rf "$TEMP_DIR"

  echo ""
  echo "gitx has been installed successfully to ${INSTALL_DIR}!"
  
  # Check if the install directory is in PATH.
  case ":$PATH:" in
    *":${INSTALL_DIR}:"*) 
      echo "Run 'gitx' to get started."
      ;;
    *)
      echo "Note: ${INSTALL_DIR} is not in your PATH."
      echo "Add it to your PATH by running:"
      echo "  echo 'export PATH=\"${INSTALL_DIR}:\$PATH\"' >> ~/.zshrc"
      echo "  source ~/.zshrc"
      echo ""
      echo "Or run gitx directly: ${INSTALL_DIR}/gitx"
      ;;
  esac
}

# Run the main function.
main