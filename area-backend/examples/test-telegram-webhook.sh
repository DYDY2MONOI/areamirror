#!/bin/bash

# Test Telegram Webhook Manually
# This script simulates a Telegram webhook call for testing purposes

BACKEND_URL=${1:-"http://localhost:8080"}

echo "=========================================="
echo "Testing Telegram Webhook"
echo "=========================================="
echo "Backend URL: ${BACKEND_URL}/webhooks/telegram"
echo ""

# Sample Telegram update payload
PAYLOAD='{
  "update_id": 123456789,
  "message": {
    "message_id": 1,
    "from": {
      "id": 987654321,
      "first_name": "Test",
      "username": "testuser"
    },
    "chat": {
      "id": 987654321,
      "type": "private"
    },
    "date": 1234567890,
    "text": "Hello from test!"
  }
}'

echo "Sending test payload:"
echo "$PAYLOAD" | jq '.' 2>/dev/null || echo "$PAYLOAD"
echo ""

RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d "$PAYLOAD" \
  "${BACKEND_URL}/webhooks/telegram")

echo "Response:"
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

if echo "$RESPONSE" | grep -q '"status":"ok"'; then
    echo "✓ Webhook test successful!"
else
    echo "✗ Webhook test failed"
    echo "Check backend logs for more details"
fi

echo "=========================================="




