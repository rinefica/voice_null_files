package model

// FileModel данные загруженного файла.
type FileModel struct {
	// UUID уникальный идентификатор для файла, сохранен под этим именем в файловой системе.
	UUID string
	// Filename имя файла, заданное пользователем.
	Filename string
	// UserID идентификатор пользователя, которому принадлежит файл.
	UserID int64
}
