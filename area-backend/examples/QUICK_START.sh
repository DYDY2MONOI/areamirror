#!/bin/bash

echo "╔════════════════════════════════════════════════════════════╗"
echo "║     Telegram Integration - Quick Start Guide             ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Step 1: Check bot token
echo -e "${BLUE}Step 1: Checking Telegram Bot Token...${NC}"
if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
    echo -e "${RED}✗ TELEGRAM_BOT_TOKEN not found in environment${NC}"
    echo ""
    read -p "Enter your Telegram Bot Token: " BOT_TOKEN
    export TELEGRAM_BOT_TOKEN=$BOT_TOKEN
else
    echo -e "${GREEN}✓ Token found!${NC}"
    BOT_TOKEN=$TELEGRAM_BOT_TOKEN
fi
echo ""

# Step 2: Get webhook info
echo -e "${BLUE}Step 2: Checking current webhook configuration...${NC}"
WEBHOOK_INFO=$(curl -s "https://api.telegram.org/bot${BOT_TOKEN}/getWebhookInfo")
echo "$WEBHOOK_INFO" | jq '.' 2>/dev/null || echo "$WEBHOOK_INFO"
echo ""

# Step 3: Setup webhook
echo -e "${BLUE}Step 3: Setting up webhook...${NC}"
echo -e "${YELLOW}What is your backend URL?${NC}"
echo "Examples:"
echo "  - http://localhost:8080 (local development)"
echo "  - https://your-domain.com (production)"
echo ""
read -p "Backend URL: " BACKEND_URL
BACKEND_URL=${BACKEND_URL%/}
WEBHOOK_URL="${BACKEND_URL}/webhooks/telegram"

echo ""
echo "Setting webhook to: $WEBHOOK_URL"
RESPONSE=$(curl -s -X POST "https://api.telegram.org/bot${BOT_TOKEN}/setWebhook?url=${WEBHOOK_URL}")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"

if echo "$RESPONSE" | grep -q '"ok":true'; then
    echo -e "${GREEN}✓ Webhook configured successfully!${NC}"
else
    echo -e "${RED}✗ Failed to configure webhook${NC}"
    exit 1
fi
echo ""

# Step 4: Get Chat ID
echo -e "${BLUE}Step 4: Getting your Chat ID...${NC}"
echo ""
echo -e "${YELLOW}Action required:${NC}"
echo "1. Open Telegram and find your bot"
echo "2. Send a message to your bot (any message)"
echo "3. Press Enter here when done"
echo ""
read -p "Press Enter after sending a message..."

UPDATES=$(curl -s "https://api.telegram.org/bot${BOT_TOKEN}/getUpdates")
CHAT_ID=$(echo "$UPDATES" | jq -r '.result[0].message.chat.id' 2>/dev/null)

if [ "$CHAT_ID" != "null" ] && [ ! -z "$CHAT_ID" ]; then
    echo -e "${GREEN}✓ Chat ID found: $CHAT_ID${NC}"
else
    echo -e "${YELLOW}⚠ Could not automatically detect Chat ID${NC}"
    echo "Visit this URL manually to get your Chat ID:"
    echo "https://api.telegram.org/bot${BOT_TOKEN}/getUpdates"
    echo ""
    read -p "Enter your Chat ID manually: " CHAT_ID
fi
echo ""

# Step 5: Test webhook
echo -e "${BLUE}Step 5: Testing webhook...${NC}"
echo "Sending test payload to $WEBHOOK_URL"

TEST_PAYLOAD=$(cat <<EOF
{
  "update_id": 999999,
  "message": {
    "message_id": 1,
    "from": {
      "id": ${CHAT_ID},
      "first_name": "Test",
      "username": "testuser"
    },
    "chat": {
      "id": ${CHAT_ID},
      "type": "private"
    },
    "date": $(date +%s),
    "text": "Test message from quick start script"
  }
}
EOF
)

WEBHOOK_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d "$TEST_PAYLOAD" \
  "$WEBHOOK_URL")

echo "$WEBHOOK_RESPONSE" | jq '.' 2>/dev/null || echo "$WEBHOOK_RESPONSE"

if echo "$WEBHOOK_RESPONSE" | grep -q '"status":"ok"'; then
    echo -e "${GREEN}✓ Webhook test successful!${NC}"
else
    echo -e "${YELLOW}⚠ Webhook response unclear. Check backend logs.${NC}"
fi
echo ""

# Step 6: Summary
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                    Setup Complete! 🎉                      ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo "Your configuration:"
echo "  • Bot Token: ${BOT_TOKEN:0:15}..."
echo "  • Webhook URL: $WEBHOOK_URL"
echo "  • Chat ID: $CHAT_ID"
echo ""
echo -e "${BLUE}Next Steps:${NC}"
echo ""
echo "1. Create an AREA with Telegram trigger:"
echo "   - Open your frontend"
echo "   - Click 'Create Area'"
echo "   - Select Telegram as trigger"
echo "   - Enter Chat ID: $CHAT_ID"
echo "   - Choose trigger type (message_received, keyword_match, or command_received)"
echo "   - Select an action service (Discord, Gmail, etc.)"
echo "   - Click 'Create Area'"
echo ""
echo "2. Test your AREA:"
echo "   - Send a message to your Telegram bot"
echo "   - Check that the action executes"
echo ""
echo "3. Examples to try:"
echo "   ./telegram-to-discord-curl.sh $BACKEND_URL YOUR_AUTH_TOKEN"
echo ""
echo -e "${GREEN}Happy automating! 🚀${NC}"
echo ""

# Save config for future reference
CONFIG_FILE=".telegram_config"
cat > "$CONFIG_FILE" <<EOF
TELEGRAM_BOT_TOKEN=$BOT_TOKEN
WEBHOOK_URL=$WEBHOOK_URL
CHAT_ID=$CHAT_ID
BACKEND_URL=$BACKEND_URL
EOF

echo -e "${BLUE}Configuration saved to $CONFIG_FILE${NC}"
echo "You can source this file in other scripts: source $CONFIG_FILE"
echo ""




