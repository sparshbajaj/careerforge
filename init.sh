#!/bin/bash
echo "==================================="
echo "  CareerForge Studio Initializer"
echo "==================================="
echo ""

# Detect OS
OS_TYPE="unknown"
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OS_TYPE="linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    OS_TYPE="mac"
fi

install_node() {
    echo "[INFO] Node.js is missing. Attempting to auto-install..."
    if [ "$OS_TYPE" == "mac" ]; then
        if command -v brew &> /dev/null; then
            brew install node
        else
            echo "[ERROR] Homebrew not found. Please install Node manually."
            exit 1
        fi
    elif [ "$OS_TYPE" == "linux" ]; then
        if command -v apt-get &> /dev/null; then
            curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
            sudo apt-get install -y nodejs
        else
            echo "[ERROR] apt-get not found. Please install Node manually."
            exit 1
        fi
    fi
}

install_go() {
    echo "[INFO] Go is missing. Attempting to auto-install..."
    if [ "$OS_TYPE" == "mac" ]; then
        if command -v brew &> /dev/null; then
            brew install go
        else
            echo "[ERROR] Homebrew not found. Please install Go manually."
            exit 1
        fi
    elif [ "$OS_TYPE" == "linux" ]; then
        if command -v apt-get &> /dev/null; then
            sudo apt-get install -y golang-go
        else
            echo "[ERROR] apt-get not found. Please install Go manually."
            exit 1
        fi
    fi
}

if ! command -v node &> /dev/null; then
    install_node
    if ! command -v node &> /dev/null; then
        echo "[ERROR] Failed to install Node.js. Please restart your terminal or install manually."
        exit 1
    fi
fi

if ! command -v go &> /dev/null; then
    install_go
    if ! command -v go &> /dev/null; then
        echo "[ERROR] Failed to install Go. Please restart your terminal or install manually."
        exit 1
    fi
fi

if ! command -v agy &> /dev/null; then
    echo "[INFO] Antigravity CLI not found. Installing globally via npm..."
    npm install -g @google/antigravity-cli
fi

echo "[INFO] Installing project dependencies..."
npm install

echo "[INFO] Building Go Dashboard..."
cd dashboard
go build .
cd ..

chmod +x launch.sh

echo ""
echo "Setup complete! Launching CareerForge Studio..."
./launch.sh
