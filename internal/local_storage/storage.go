package storage

import (
	"log"
	"os"
	"path"
)

// Временное хранилище, используется для сохранения токена между запросами.
type TempStorage interface {
	SaveToken(token string) error
	GetToken() (string, error)
}

type TempStorageImpl struct {
	dir          string
	tempFileName string
}

func NewTempStorage(dir string, filename string) TempStorage {
	return &TempStorageImpl{
		dir:          dir,
		tempFileName: filename,
	}
}

// Сохраняет токен во временном файле, создает файл при необходимости.
func (t *TempStorageImpl) SaveToken(token string) error {

	if _, err := os.Stat(t.dir); os.IsNotExist(err) {
		err := os.Mkdir(t.dir, 0777)
		if err != nil {
			return err
		}
	}

	filepath := path.Join(t.dir, t.tempFileName)
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatal(err)
	}
	println("create file " + f.Name())

	_, err = f.WriteString(token)
	if err != nil {
		return err
	}
	return nil
}

// Чтение токена из временного файла.
func (t *TempStorageImpl) GetToken() (string, error) {
	filepath := path.Join(t.dir, t.tempFileName)
	b, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	token := string(b)
	return token, nil
}
