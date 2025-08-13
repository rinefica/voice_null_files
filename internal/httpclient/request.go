package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	model "github.com/rinefica/voice_null_files/internal/domain/model/cli"
)

// CreateFileBodyRequest создать запрос для отправки файла filename, добавить token.
func CreateFileBodyRequest(url, token, filename string) (*http.Request, error) {
	imageBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	fileWriter, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return nil, err
	}
	_, err = fileWriter.Write(imageBytes)
	if err != nil {
		fmt.Println("Error writing file data:", err)
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", token)

	return req, nil
}

// CreateDataRequestBody создать body для json-значений отправки текстовой информации.
func CreateDataRequestBody(data string, additional string, dataType string) (*bytes.Buffer, error) {
	dataRequestModel := model.InfoDataRequestModel{
		Data:           data,
		AdditionalData: additional,
		Type:           dataType,
	}
	jsonData, err := json.Marshal(dataRequestModel)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(jsonData)
	return body, nil
}
