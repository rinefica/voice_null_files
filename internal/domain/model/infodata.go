package model

// InfoDataModel модель для данных.
type InfoDataModel struct {
	// UUID уникальный идентификатор.
	UUID string
	// Data текстовая сохраненная информация.
	Data string
	// AdditionalData дополнительная текстовая информация.
	AdditionalData string
	// Type тип сохраняемой информации: txt - текст, login - пара логин-пароль, bank - банковская карта.
	Type string
	// UserID идентификатор пользователя, которому принадлежит запись.
	UserID int64
}

// InfoDataBody модель для данных в запросе.
type InfoDataBody struct {
	// Data текстовая сохраненная информация.
	Data string
	// AdditionalData дополнительная текстовая информация.
	AdditionalData string `json:"additional_data"`
	// Type тип сохраняемой информации: txt - текст, login - пара логин-пароль, bank - банковская карта.
	Type string
}
