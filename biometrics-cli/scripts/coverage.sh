#!/bin/bash

set -e

MIN_COVERAGE=70

echo "========================================"
echo "BIOMETRICS CLI - Coverage Report"
echo "========================================"
echo ""

cd "$(dirname "$0")/.."

echo "Running tests with coverage..."
echo ""

go test ./... -coverprofile=coverage.out -covermode=atomic

echo ""
echo "Generating HTML coverage report..."
echo ""

go tool cover -html=coverage.out -o coverage.html

echo ""
echo "========================================"
echo "Coverage Summary"
echo "========================================"
echo ""

TOTAL=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
echo "Total Coverage: ${TOTAL}%"

if [ "$MIN_COVERAGE" -gt 0 ]; then
    COVERAGE_INT=${TOTAL%.*}
    if [ "$COVERAGE_INT" -lt "$MIN_COVERAGE" ]; then
        echo ""
        echo "ERROR: Coverage ($COVERAGE_INT%) is below minimum threshold ($MIN_COVERAGE%)"
        echo ""
        exit 1
    fi
    echo ""
    echo "Coverage check PASSED (>= $MIN_COVERAGE%)"
fi

echo ""
echo "Reports generated:"
echo "  - coverage.out (raw data)"
echo "  - coverage.html (HTML report)"
echo ""
echo "View HTML report with: open coverage.html"
echo ""
