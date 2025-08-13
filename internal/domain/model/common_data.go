package model

// CommonData общий интерфейс для получения данных пользователя - записях и файлах.
type CommonData struct {
	UUID           string `json:"uuid"`
	AdditionalData string `json:"additional_data"`
	Type           string `json:"type"`
}
