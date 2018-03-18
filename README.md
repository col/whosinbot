# WhosInBot

## Register webhook

Set Webhook
```
curl -XPOST https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/setWebhook?url=$WEBHOOK_URL/$TELEGRAM_BOT_TOKEN
```

Get Webhook Info
```
curl -XPOST https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/getWebhookInfo
```
