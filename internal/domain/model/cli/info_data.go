package cli

// InfoDataRequestModel модель для данных в запросе.
type InfoDataRequestModel struct {
	// Data текстовая сохраненная информация.
	Data string
	// AdditionalData дополнительная текстовая информация.
	AdditionalData string `json:"additional_data"`
	// Type тип сохраняемой информации: txt - текст, login - пара логин-пароль, bank - банковская карта.
	Type string
}

// InfoDataResponseModel модель для данных при ответе запроса.
type InfoDataResponseModel struct {
	// UUID уникальный идентификатор.
	UUID string
	// Data текстовая сохраненная информация.
	Data string
	// AdditionalData дополнительная текстовая информация.
	AdditionalData string
	// Type тип сохраняемой информации: txt - текст, login - пара логин-пароль, bank - банковская карта.
	Type string
}
