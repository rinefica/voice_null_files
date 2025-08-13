package cli

// AllInfoDataResponseModel обобщенная модель для данных в общем запросе.
type AllInfoDataResponseModel struct {
	Data []struct {
		UUID           string `json:"uuid"`
		AdditionalData string `json:"additional_data"`
		Type           string `json:"type"`
	} `json:"data"`
}
