package httpclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	model "github.com/rinefica/voice_null_files/internal/domain/model/cli"
)

type CliHTTPClient struct {
	baseURL string
	client  *http.Client
}

func NewCliHTTPClient(baseURL string) *CliHTTPClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &CliHTTPClient{
		baseURL: baseURL,
		client:  &http.Client{Transport: tr},
	}
}

// Login авторизация в системе по email и password.
func (c *CliHTTPClient) Login(email, password string) (string, error) {
	url := c.baseURL + "login"
	loginRequest, err := json.Marshal(model.LoginRequestModel{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(loginRequest))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	var respModel model.LoginResponseModel
	if err := json.NewDecoder(resp.Body).Decode(&respModel); err != nil {
		return "", err
	}
	return respModel.Token, nil
}

// Register регистрация в системе по email и password.
func (c *CliHTTPClient) Register(email, password string) error {
	url := c.baseURL + "register"
	registerRequest, err := json.Marshal(model.RegisterRequestModel{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(registerRequest))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)
	return nil
}

// GetAllInfoData получить всю хранимую информацию по пользователю.
func (c *CliHTTPClient) GetAllInfoData(token string) error {
	url := c.baseURL + "api/info_data/all"
	resp, err := c.GetWithToken(url, token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		println("для начала работы, пожалуйста, авторизуйтесь в системе")
		return nil
	}

	var respModel model.AllInfoDataResponseModel
	if err := json.NewDecoder(resp.Body).Decode(&respModel); err != nil {
		return err
	}

	if len(respModel.Data) == 0 {
		fmt.Println("Сохраненных данных нет")
	}

	for _, data := range respModel.Data {
		fmt.Printf("%s %s ID %s\n", data.Type, data.AdditionalData, data.UUID)
	}

	fmt.Println("Response status:", resp.Status)
	return nil
}

// GetInfoData получить одну запись, принадлежащую пользователю.
func (c *CliHTTPClient) GetInfoData(token, uuid string) error {
	url := c.baseURL + "api/info_data/" + uuid
	resp, err := c.GetWithToken(url, token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		println("для начала работы, пожалуйста, авторизуйтесь в системе")
		return nil
	}
	var respModel model.InfoDataResponseModel
	if err := json.NewDecoder(resp.Body).Decode(&respModel); err != nil {
		return err
	}

	fmt.Printf("%s: %s - %s ID %s\n", respModel.Type, respModel.Data, respModel.AdditionalData, respModel.UUID)

	fmt.Println("Response status:", resp.Status)
	return nil
}

// DownloadFile скачать файл с сервера, сохранить под именем filename.
func (c *CliHTTPClient) DownloadFile(token, uuid, filename string) error {
	url := c.baseURL + "api/file/" + uuid
	resp, err := c.GetWithToken(url, token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		println("для начала работы, пожалуйста, авторизуйтесь в системе")
		return nil
	}

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)

	println("Response status:", resp.Status)
	return err
}

// UploadFile загрузить файл на сервер.
func (c *CliHTTPClient) UploadFile(token, filename string) error {
	url := c.baseURL + "api/file/"
	req, err := CreateFileBodyRequest(url, token, filename)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		println("для начала работы, пожалуйста, авторизуйтесь в системе")
		return nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	println("Answer:", string(bodyBytes))
	println("Response status:", resp.Status)
	return err
}

// SaveData сохранить текстовую информацию клиента.
func (c *CliHTTPClient) SaveData(token, data, additional, dataType string) error {
	url := c.baseURL + "api/info_data/"
	body, err := CreateDataRequestBody(data, additional, dataType)
	if err != nil {
		return err
	}
	resp, err := c.PostWithTokenJSON(url, token, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		println("для начала работы, пожалуйста, авторизуйтесь в системе")
		return nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	println("Answer:", string(bodyBytes))
	println("Response status:", resp.Status)
	return nil
}

// GetWithToken сделать GET-запрос, добавив в хедер token.
func (c *CliHTTPClient) GetWithToken(url, token string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", token)
	return c.client.Do(req)
}

// PostWithTokenJSON сделать POST-запрос, добавить token и тип для json-данных.
func (c *CliHTTPClient) PostWithTokenJSON(url, token string, body io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	return c.client.Do(req)
}
