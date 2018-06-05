# WhosInBot Plan


## Package whosinbot
    WhosInBot
        DataStore DataStore
	HandleCommand(command Command) (Response, Error)

## Package whosinbot/domain
	Command struct
		name string
		params string[]

	Response struct
		message string

	RollCall struct
		ChatID string
		Title string
	
	RollCallResponse struct
		ChatID int64
		UserID int64
		Response string
		Reasons string 
	
	DataStore interface
        SetRollCall(RollCall)
        DeleteRollCall(RollCall)
        SetResponse(RollCallResponse)	
        DeleteResponse(RollCallResponse)
    			
## Package whosinbot/dynamodb
	DynamoDataStore
    	SetRollCall(RollCall)
        DeleteRollCall(RollCall)
        SetResponse(RollCallResponse)	
        DeleteResponse(RollCallResponse)
        
	NewDynamoDataStore(DynamoConfig)

## Package whoinbot/helpers
	Helpers
	
## Package telegram
	NewTelegram(TelegramConfig)
	Telegram
		BotApi
		ParseUpdateRequest(string) (Command, error)
		SendMessage(Response) (error)

## Package cmd/telegram_lambda 	
    Main
        LoadDynamoConfig
        LoadTelegramConfig

## Package cmd/telegram_http 	
    Main
        LoadDynamoConfig
        LoadTelegramConfig
        LoadHttpConfig

##Package whosinbot/config
	TelegramConfig
	HttpConfig
	DynamoConfig
	
	
## Commands
- start_roll_call (title)
- in (reason)
- out (reason)
- set_title (title)
- set_in_for (name)
- end_roll_call
- ssh 


