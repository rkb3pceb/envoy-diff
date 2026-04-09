#!/bin/bash
# Demo script for envoy-diff

set -e

echo "=== envoy-diff Demo ==="
echo ""

# Create example directory
mkdir -p examples

# Create old.env
cat > examples/old.env << EOF
# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=old_secret_123

# API configuration
API_KEY=sk_test_old_key
API_ENDPOINT=https://api.example.com/v1
API_TIMEOUT=30

# Feature flags
FEATURE_NEW_UI=false
FEATURE_ANALYTICS=true

# Other settings
LOG_LEVEL=info
MAX_CONNECTIONS=100
EOF

# Create new.env
cat > examples/new.env << EOF
# Database configuration
DB_HOST=db.production.com
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=new_secret_456

# API configuration
API_KEY=sk_live_new_key
API_ENDPOINT=https://api.example.com/v2
API_TIMEOUT=60

# Feature flags
FEATURE_NEW_UI=true
FEATURE_ANALYTICS=true
FEATURE_BETA=true

# Other settings
LOG_LEVEL=debug
MAX_CONNECTIONS=200
CACHE_ENABLED=true
EOF

echo "Created example environment files:"
echo "  - examples/old.env"
echo "  - examples/new.env"
echo ""

echo "1. Basic diff (text format):"
echo "   $ envoy-diff examples/old.env examples/new.env"
echo ""

if command -v envoy-diff &> /dev/null; then
    envoy-diff examples/old.env examples/new.env
else
    echo "   (envoy-diff not installed, run 'make install' first)"
fi

echo ""
echo "2. JSON format output:"
echo "   $ envoy-diff --format json examples/old.env examples/new.env"
echo ""

if command -v envoy-diff &> /dev/null; then
    envoy-diff --format json examples/old.env examples/new.env
fi

echo ""
echo "3. With security audit:"
echo "   $ envoy-diff --audit examples/old.env examples/new.env"
echo ""

if command -v envoy-diff &> /dev/null; then
    envoy-diff --audit examples/old.env examples/new.env || true
fi

echo ""
echo "=== Demo Complete ==="
