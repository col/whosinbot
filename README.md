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

## Dependencies

Install AWS CLI tools

```  
brew install awscli
aws configure --profile col.w.harris
```

Install serverless

```
npm install -g serverless

```


## Build

```
make build
```

## Deploy

```
source .env
serverless deploy
```
    
## Links

- https://github.com/go-telegram-bot-api/telegram-bot-api
- https://core.telegram.org/bots/api