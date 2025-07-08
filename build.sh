#!/bin/bash

# Build script for Juggler application

echo "Building Juggler application..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Build the application
go build -o bin/juggler cmd/app/main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo "Run the application with: ./bin/juggler [port]"
    echo "Example: ./bin/juggler (uses default port 8080)"
    echo "Example: ./bin/juggler 9000 (uses custom port 9000)"
    echo ""
    echo "Configure balls and time through the web interface"
    echo "Web interface will be available at: http://localhost:8080"
else
    echo "❌ Build failed!"
    exit 1
fi
