# Telegram Integration - Complete Guide

## 🎉 What's New

Telegram is now fully integrated as both a **Trigger** and an **Action** in the AREA platform!

### Before
- ✅ Telegram as Action (send messages)

### Now
- ✅ Telegram as Action (send messages)
- ✅ **Telegram as Trigger** (receive messages and trigger workflows)

## 📋 Quick Start

### 1. Setup Environment

Make sure `TELEGRAM_BOT_TOKEN` is set in your `.env` file or environment:

```bash
TELEGRAM_BOT_TOKEN=your_bot_token_here
```

### 2. Configure Webhook

Run the setup script:

```bash
./examples/setup-telegram-webhook.sh
```

Or manually:

```bash
curl -X POST "https://api.telegram.org/bot<YOUR_TOKEN>/setWebhook?url=https://your-domain.com/webhooks/telegram"
```

### 3. Get Your Chat ID

1. Send a message to your bot
2. Visit: `https://api.telegram.org/bot<YOUR_TOKEN>/getUpdates`
3. Find the `chat.id` in the response

### 4. Create Your First Telegram Trigger

Use one of the example files or create via API:

```bash
./examples/telegram-to-discord-curl.sh http://localhost:8080 YOUR_AUTH_TOKEN
```

## 🔧 Implementation Details

### Files Modified/Created

#### Backend Core
- `services/telegram.go` - Added `TelegramUpdate` struct for webhook payloads
- `services/scheduler.go` - Added `ExecuteAreaPublic()` method and Telegram template variables
- `controllers/telegram_webhook.go` - **NEW** - Handles incoming Telegram webhooks
- `main.go` - Added `/webhooks/telegram` route

#### Examples & Documentation
- `examples/TELEGRAM_SETUP.md` - Complete setup guide
- `examples/telegram-trigger-example.json` - Basic message trigger
- `examples/telegram-keyword-example.json` - Keyword matching
- `examples/telegram-command-example.json` - Command detection
- `examples/setup-telegram-webhook.sh` - Automated webhook setup
- `examples/test-telegram-webhook.sh` - Manual webhook testing
- `examples/telegram-to-discord-curl.sh` - Complete AREA creation example

## 🎯 Trigger Types

### 1. Message Received (Any Message)

Triggers on any message sent to the bot:

```json
{
  "trigger_service": "Telegram",
  "trigger_type": "message_received",
  "trigger_config": {
    "chatId": "123456789"
  }
}
```

### 2. Keyword Match

Triggers when message contains specific keyword:

```json
{
  "trigger_service": "Telegram",
  "trigger_type": "keyword_match",
  "trigger_config": {
    "chatId": "123456789",
    "keyword": "alert"
  }
}
```

### 3. Command Received

Triggers on specific commands:

```json
{
  "trigger_service": "Telegram",
  "trigger_type": "command_received",
  "trigger_config": {
    "chatId": "123456789",
    "command": "/start"
  }
}
```

## 📊 Template Variables

When Telegram triggers an AREA, these variables are available:

| Variable | Description | Example |
|----------|-------------|---------|
| `{{messageText}}` | Message content | "Hello world" |
| `{{chatId}}` | Chat ID | "123456789" |
| `{{username}}` | Sender username | "johndoe" |
| `{{firstName}}` | Sender first name | "John" |
| `{{messageId}}` | Message ID | "42" |
| `{{areaName}}` | AREA name | "My Telegram Bot" |
| `{{triggerService}}` | Always "Telegram" | "Telegram" |
| `{{actionService}}` | Action service name | "Discord" |
| `{{eventTime}}` | Trigger timestamp | "2025-10-20 15:30:45" |

## 💡 Use Case Examples

### 1. Telegram to Discord Bridge

Forward all Telegram messages to Discord:

```json
{
  "name": "Telegram to Discord",
  "trigger_service": "Telegram",
  "trigger_type": "message_received",
  "trigger_config": {
    "chatId": "123456789"
  },
  "action_service": "Discord",
  "action_config": {
    "webhookUrl": "YOUR_WEBHOOK",
    "message": "From {{firstName}}: {{messageText}}"
  }
}
```

### 2. Alert System

Send emails on urgent keywords:

```json
{
  "name": "Urgent Alerts",
  "trigger_service": "Telegram",
  "trigger_type": "keyword_match",
  "trigger_config": {
    "chatId": "123456789",
    "keyword": "urgent"
  },
  "action_service": "Gmail",
  "action_config": {
    "toEmail": "admin@example.com",
    "subject": "Urgent Alert",
    "body": "{{messageText}} from {{firstName}}"
  }
}
```

### 3. Bot Commands

Auto-respond to commands:

```json
{
  "name": "Status Command",
  "trigger_service": "Telegram",
  "trigger_type": "command_received",
  "trigger_config": {
    "chatId": "123456789",
    "command": "/status"
  },
  "action_service": "Telegram",
  "action_config": {
    "chatId": "123456789",
    "message": "System is online! ✅"
  }
}
```

## 🧪 Testing

### Test Webhook Locally

```bash
./examples/test-telegram-webhook.sh http://localhost:8080
```

### Create Test AREA

```bash
./examples/telegram-to-discord-curl.sh http://localhost:8080 YOUR_AUTH_TOKEN
```

### Send Real Message

Just send a message to your Telegram bot and check the logs!

## 🔍 Troubleshooting

### Webhook Not Working

```bash
# Check webhook status
curl "https://api.telegram.org/bot<TOKEN>/getWebhookInfo"

# Remove webhook
curl -X POST "https://api.telegram.org/bot<TOKEN>/deleteWebhook"

# Re-setup
./examples/setup-telegram-webhook.sh
```

### Messages Not Triggering

1. Check AREA is active: `is_active = true`
2. Verify chat ID matches
3. Check backend logs for errors
4. Ensure trigger type is correct

### Bot Not Responding

1. Verify `TELEGRAM_BOT_TOKEN` is set
2. Check bot has message permissions
3. Ensure webhook URL is accessible (HTTPS required for production)

## 🚀 Advanced Features

### Multiple Chat Support

Create different AREAs for different chats:

```bash
# Personal notifications
AREA 1: chatId="123" -> Discord Channel A

# Work alerts  
AREA 2: chatId="456" -> Gmail + Discord Channel B

# Group monitoring
AREA 3: chatId="789" -> Slack + Email
```

### Keyword Filtering

Use multiple AREAs with different keywords:

```bash
AREA 1: keyword="bug" -> Notify dev team
AREA 2: keyword="urgent" -> Alert everyone
AREA 3: keyword="help" -> Create support ticket
```

### Chain Reactions

Telegram can trigger multiple actions via different AREAs:

```bash
Message with "deploy" ->
  AREA 1: Log to Discord
  AREA 2: Send email notification
  AREA 3: Update Google Sheets
```

## 📝 API Endpoints

### Webhook Endpoint

**POST** `/webhooks/telegram`

Receives Telegram updates automatically when configured.

### Related Endpoints

- `POST /areas` - Create new AREA with Telegram trigger
- `GET /areas` - List all AREAs
- `PATCH /areas/:id/toggle` - Enable/disable AREA
- `DELETE /areas/:id` - Delete AREA

## 🔒 Security Considerations

1. **HTTPS Required**: Telegram requires HTTPS for webhooks in production
2. **Token Security**: Never commit `TELEGRAM_BOT_TOKEN` to git
3. **Chat ID Validation**: Verify chat IDs to prevent unauthorized triggers
4. **Rate Limiting**: Consider implementing rate limits for production
5. **Webhook Validation**: Optionally validate Telegram webhook signatures

## 📚 Resources

- [Telegram Bot API Documentation](https://core.telegram.org/bots/api)
- [BotFather - Create Bots](https://t.me/botfather)
- [Webhook Guide](https://core.telegram.org/bots/webhooks)

## 🎓 Learning Path

1. ✅ Create a bot with BotFather
2. ✅ Set up webhook with our script
3. ✅ Create simple message_received trigger
4. ✅ Test with Telegram to Discord
5. ✅ Try keyword matching
6. ✅ Build custom command handlers
7. ✅ Combine multiple triggers and actions

## 🤝 Contributing

Found a bug or want to improve Telegram integration? 

1. Check existing issues
2. Create new issue with details
3. Submit PR with tests

---

**Happy Automating! 🤖✨**




