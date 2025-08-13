package main

import (
	"log"
	"os"

	"github.com/rinefica/voice_null_files/internal/cmd_handler"
	"github.com/rinefica/voice_null_files/internal/httpclient"
	storage "github.com/rinefica/voice_null_files/internal/local_storage"
	"github.com/urfave/cli/v2"
)

// Пример работы с приложением
// Регистрация
// go run main.go r testcli 12345
// Логин
// go run main.go l testcli 12345
// Посмотрели все загруженные данные
// go run main.go w
// Сохранили текст
// go run main.go sd_txt somedatatxt somecomment
// Просмотрели конкретную запись
// go run main.go w [UUID]
// Сохранили банковские данные
// go run main.go sd_bank 1234567 somecomment
// Сохранили логин и пароль
// go run main.go sd_login test_login test_pswrd somecomment
// Загрузили картинку
// go run main.go sf [path]
// Просмотрели файл
// go run main.go f [UUID] 123.png
func main() {
	client := httpclient.NewCliHTTPClient("http://localhost:8000/")
	// для сохранения токена между запусками записываем его во временный файл
	strg := storage.NewTempStorage("./storage", "access.log")
	cmds := cmd_handler.NewCommandsHandler(client, strg)
	app := &cli.App{
		Name:  "voice_null",
		Usage: "store passwords, logins & data",
		Commands: []*cli.Command{
			{
				Name:    "reg",
				Aliases: []string{"r"},
				Usage:   "введите логин и пароль через пробел для регистрации в системе",
				Action:  cmds.Register,
			},
			{
				Name:    "login",
				Aliases: []string{"l"},
				Usage:   "введите логин и пароль через пробел для авторизации в системе",
				Action:  cmds.Login,
			},
			{
				Name:    "watch",
				Aliases: []string{"w"},
				Usage:   "просмотреть все сохраненные в системе данные или запись по uuid",
				Action:  cmds.InfoData,
			},
			{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "введите ID и название файла. скачать файл по ID в текущую директорию и сохранить с текущим именем",
				Action:  cmds.File,
			},
			{
				Name:    "savefile",
				Aliases: []string{"sf"},
				Usage:   "введите путь до файла для загрузки на сервер",
				Action:  cmds.SaveFile,
			},
			{
				Name:    "savedata-text",
				Aliases: []string{"sd_txt"},
				Usage:   "сохранение текстовой информации: данные и комментарий через пробел",
				Action:  cmds.SaveDataText,
			},
			{
				Name:    "savedata-bank",
				Aliases: []string{"sd_bank"},
				Usage:   "сохранение банковской карты: номер и комментарий через пробел",
				Action:  cmds.SaveDataBank,
			},
			{
				Name:    "savedata-login",
				Aliases: []string{"sd_login"},
				Usage:   "сохранение данных авторизации: логин, пароль и комментарий через пробел",
				Action:  cmds.SaveDataLogin,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
