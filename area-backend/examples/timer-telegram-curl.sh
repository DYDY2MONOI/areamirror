#!/bin/bash

# Script de test pour créer une Area Timer → Telegram
# Usage: ./timer-telegram-curl.sh YOUR_JWT_TOKEN YOUR_CHAT_ID

JWT_TOKEN="${1:-YOUR_JWT_TOKEN}"
CHAT_ID="${2:-123456789}"
API_URL="${3:-http://localhost:8080}"

echo "🚀 Création d'une Area Timer → Telegram"
echo "📍 API: $API_URL"
echo "💬 Chat ID: $CHAT_ID"
echo ""

curl -X POST "$API_URL/api/areas" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d @- << EOF
{
  "name": "Test Timer 30s",
  "description": "Area de test avec intervalle de 30 secondes",
  "trigger_service": "Timer",
  "trigger_type": "interval",
  "trigger_config": {
    "interval": "30s"
  },
  "action_service": "Telegram",
  "action_type": "send_message",
  "action_config": {
    "chatId": "$CHAT_ID",
    "message": "✅ Test réussi !\n\n📋 Area : {{areaName}}\n🕐 Déclenché le : {{triggerTime}}\n⏱️ Intervalle : {{interval}}"
  },
  "is_active": true,
  "status": "enabled"
}
EOF

echo ""
echo ""
echo "✅ Requête envoyée !"
echo "⏳ Attendez environ 30-60 secondes pour recevoir le message Telegram"
echo "📋 Vérifiez les logs: docker logs -f area-backend"


