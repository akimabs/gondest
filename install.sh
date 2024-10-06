#!/bin/bash

# Set variables
REPO_URL="https://github.com/akimabs/gondest.git"
INSTALL_PATH=/usr/local/bin
TEMPLATE_PATH=/usr/local/share/gondest/templates

# Clone the repository
git clone $REPO_URL
cd gondest

# Run make install
make install

# Clean up
cd ..
rm -rf gondest

echo "Gondest installed successfully!"