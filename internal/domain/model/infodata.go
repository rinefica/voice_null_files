package model

type InfoDataModel struct {
	UUID           string
	Data           string
	AdditionalData string
	Type           string
	UserID         int64
}

type InfoDataBody struct {
	Data           string
	AdditionalData string `json:"additional_data"`
	Type           string
}
