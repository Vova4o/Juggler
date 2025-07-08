#!/bin/bash

# Test runner script for Juggler application

echo "🧪 Running Juggler Tests..."
echo "=========================="

# Run all tests
echo "📋 Running all tests..."
go test -v ./test/...

if [ $? -eq 0 ]; then
    echo "✅ All tests passed!"
else
    echo "❌ Some tests failed!"
    exit 1
fi

echo ""
echo "📊 Running benchmarks..."
go test -bench=. ./test/...

echo ""
echo "🔍 Running tests with coverage..."
go test -cover -coverpkg=./internal/... ./test/...

echo ""
echo "📈 Generating detailed coverage report..."
go test -coverprofile=coverage.out -coverpkg=./internal/... ./test/...
if [ -f coverage.out ]; then
    go tool cover -html=coverage.out -o coverage.html
    echo "Coverage report generated: coverage.html"
    echo "📊 Coverage summary:"
    go tool cover -func=coverage.out
else
    echo "⚠️  No coverage data generated"
fi

echo ""
echo "🏁 Test run complete!"
echo "Coverage report generated: coverage.html"
