package conversation

var (
	DefaultSettings = RequestSettings{
		Model:            "text-davinci-003",
		MaxTokens:        500,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0.6,
		Temperature:      0.9,
	}
	QuickChatSettings = RequestSettings{
		Model:            "text-davinci-003",
		MaxTokens:        3500,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0.6,
		Temperature:      0.9,
	}
	SummarySettings = RequestSettings{
		Model:            "text-davinci-003",
		MaxTokens:        2048,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  1,
		Temperature:      0.7,
	}
	FriendSettings = RequestSettings{
		Model:            "text-davinci-003",
		MaxTokens:        512,
		TopP:             1,
		FrequencyPenalty: 0.5,
		PresencePenalty:  0,
		Temperature:      0.5,
	}
)
