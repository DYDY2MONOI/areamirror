#!/bin/bash

# Telegram Webhook Setup Script
# This script helps you set up the Telegram webhook for your AREA backend

echo "=========================================="
echo "Telegram Webhook Setup for AREA"
echo "=========================================="
echo ""

# Check if TELEGRAM_BOT_TOKEN is set
if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
    echo "⚠️  TELEGRAM_BOT_TOKEN environment variable is not set"
    echo ""
    read -p "Enter your Telegram Bot Token: " BOT_TOKEN
else
    BOT_TOKEN=$TELEGRAM_BOT_TOKEN
    echo "✓ Found TELEGRAM_BOT_TOKEN in environment"
fi

echo ""
read -p "Enter your backend URL (e.g., https://your-domain.com): " BACKEND_URL

# Remove trailing slash if present
BACKEND_URL=${BACKEND_URL%/}

WEBHOOK_URL="${BACKEND_URL}/webhooks/telegram"

echo ""
echo "Setting webhook to: $WEBHOOK_URL"
echo ""

# Set the webhook
RESPONSE=$(curl -s -X POST "https://api.telegram.org/bot${BOT_TOKEN}/setWebhook?url=${WEBHOOK_URL}")

echo "Response from Telegram:"
echo "$RESPONSE"
echo ""

# Check if successful
if echo "$RESPONSE" | grep -q '"ok":true'; then
    echo "✓ Webhook set successfully!"
    echo ""
    echo "To verify, run:"
    echo "curl \"https://api.telegram.org/bot${BOT_TOKEN}/getWebhookInfo\""
else
    echo "✗ Failed to set webhook"
    echo "Please check your bot token and backend URL"
    exit 1
fi

echo ""
echo "=========================================="
echo "Next Steps:"
echo "=========================================="
echo "1. Send a message to your bot on Telegram"
echo "2. Get your Chat ID by visiting:"
echo "   https://api.telegram.org/bot${BOT_TOKEN}/getUpdates"
echo ""
echo "3. Create an AREA with Telegram as trigger:"
echo "   - trigger_service: Telegram"
echo "   - trigger_type: message_received"
echo "   - trigger_config: { \"chatId\": \"YOUR_CHAT_ID\" }"
echo ""
echo "4. Test by sending a message to your bot!"
echo "=========================================="




