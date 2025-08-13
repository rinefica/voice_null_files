package cmd_handler

import (
	"github.com/rinefica/voice_null_files/internal/httpclient"
	storage "github.com/rinefica/voice_null_files/internal/local_storage"
	"github.com/urfave/cli/v2"
)

// CommandsHandler обрабатывает команды cli приложения.
type CommandsHandler interface {
	Register(*cli.Context) error
	Login(c *cli.Context) error
	InfoData(c *cli.Context) error
	File(c *cli.Context) error
	SaveFile(c *cli.Context) error
	SaveDataText(c *cli.Context) error
	SaveDataBank(c *cli.Context) error
	SaveDataLogin(c *cli.Context) error
}

type CommandsHandlerImpl struct {
	client  *httpclient.CliHTTPClient
	storage storage.TempStorage
}

func NewCommandsHandler(
	client *httpclient.CliHTTPClient,
	storage storage.TempStorage,
) CommandsHandler {
	return &CommandsHandlerImpl{
		client:  client,
		storage: storage,
	}
}

func (h *CommandsHandlerImpl) Register(c *cli.Context) error {
	email := c.Args().First()
	pass := c.Args().Get(1)
	err := h.client.Register(email, pass)
	if err != nil {
		return err
	}
	return nil
}

func (h *CommandsHandlerImpl) Login(c *cli.Context) error {
	email := c.Args().First()
	pass := c.Args().Get(1)
	token, err := h.client.Login(email, pass)
	if err != nil {
		return err
	}
	println("Token: " + token)
	return h.storage.SaveToken(token)
}

func (h *CommandsHandlerImpl) InfoData(c *cli.Context) error {
	token, err := h.storage.GetToken()
	if err != nil {
		println("CommandsHandlerImpl error")
		return err
	}
	if c.NArg() > 0 {
		return h.client.GetInfoData(token, c.Args().Get(0))
	} else {
		return h.client.GetAllInfoData(token)
	}
}

func (h *CommandsHandlerImpl) File(c *cli.Context) error {
	token, err := h.storage.GetToken()
	if err != nil {
		println("CommandsHandlerImpl error")
		return err
	}
	if c.NArg() < 2 {
		println("Введите данные для сохранения")
		return nil
	}
	return h.client.DownloadFile(token, c.Args().First(), c.Args().Get(1))
}

func (h *CommandsHandlerImpl) SaveFile(c *cli.Context) error {
	token, err := h.storage.GetToken()
	if err != nil {
		println("CommandsHandlerImpl error")
		return err
	}
	if c.NArg() < 1 {
		println("Введите данные для сохранения")
		return nil
	}
	return h.client.UploadFile(token, c.Args().First())
}

func (h *CommandsHandlerImpl) SaveDataText(c *cli.Context) error {
	token, err := h.storage.GetToken()
	if err != nil {
		println("CommandsHandlerImpl error")
		return err
	}
	dataType := "txt"
	if c.NArg() < 1 {
		println("Введите данные для сохранения")
		return nil
	}
	data := c.Args().First()
	additional := ""
	if c.NArg() > 1 {
		additional = c.Args().Get(1)
	}

	return h.client.SaveData(token, data, additional, dataType)
}

func (h *CommandsHandlerImpl) SaveDataBank(c *cli.Context) error {
	token, err := h.storage.GetToken()
	if err != nil {
		println("CommandsHandlerImpl error")
		return err
	}
	dataType := "bank"
	if c.NArg() < 1 {
		println("Введите данные для сохранения")
		return nil
	}
	data := c.Args().First()
	additional := ""
	if c.NArg() > 1 {
		additional = c.Args().Get(1)
	}

	return h.client.SaveData(token, data, additional, dataType)
}

func (h *CommandsHandlerImpl) SaveDataLogin(c *cli.Context) error {
	token, err := h.storage.GetToken()
	if err != nil {
		println("CommandsHandlerImpl error")
		return err
	}
	dataType := "login"
	if c.NArg() < 2 {
		println("Введите данные для сохранения")
		return nil
	}
	data := c.Args().First() + ":" + c.Args().Get(1)
	additional := ""
	if c.NArg() > 2 {
		additional = c.Args().Get(2)
	}

	return h.client.SaveData(token, data, additional, dataType)
}
