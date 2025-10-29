#!/bin/bash

# Create Telegram to Discord AREA using curl
# This creates an AREA that forwards Telegram messages to Discord

BACKEND_URL=${1:-"http://localhost:8080"}
AUTH_TOKEN=${2:-"YOUR_AUTH_TOKEN"}

echo "=========================================="
echo "Creating Telegram to Discord AREA"
echo "=========================================="
echo ""

# Prompt for configuration
read -p "Enter your Telegram Chat ID: " CHAT_ID
read -p "Enter your Discord Webhook URL: " DISCORD_WEBHOOK

PAYLOAD=$(cat <<EOF
{
  "name": "Telegram to Discord Bridge",
  "description": "Automatically forward Telegram messages to Discord",
  "trigger_service": "Telegram",
  "trigger_type": "message_received",
  "trigger_config": {
    "chatId": "$CHAT_ID"
  },
  "trigger_icon_url": "https://cdn-icons-png.flaticon.com/512/2111/2111646.png",
  "action_service": "Discord",
  "action_type": "send_message",
  "action_config": {
    "webhookUrl": "$DISCORD_WEBHOOK",
    "message": "📱 **New Telegram Message**\n\nFrom: {{firstName}} (@{{username}})\nMessage: {{messageText}}\n\nChat ID: {{chatId}}\nTime: {{eventTime}}"
  },
  "action_icon_url": "https://cdn-icons-png.flaticon.com/512/5968/5968756.png",
  "is_active": true,
  "is_public": false
}
EOF
)

echo "Creating AREA with payload:"
echo "$PAYLOAD" | jq '.' 2>/dev/null || echo "$PAYLOAD"
echo ""

RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -d "$PAYLOAD" \
  "${BACKEND_URL}/areas")

echo "Response:"
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

if echo "$RESPONSE" | grep -q '"id"'; then
    AREA_ID=$(echo "$RESPONSE" | jq -r '.id' 2>/dev/null)
    echo "✓ AREA created successfully!"
    echo "  Area ID: $AREA_ID"
    echo ""
    echo "Now send a message to your Telegram bot to test it!"
else
    echo "✗ Failed to create AREA"
    echo "Check your authentication token and backend URL"
fi

echo "=========================================="




