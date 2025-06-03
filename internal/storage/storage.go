package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rinefica/voice_null_files/internal/domain/model"
	"github.com/rinefica/voice_null_files/internal/lib/sl"
	"github.com/rinefica/voice_null_files/internal/storage/mapper"
	"log/slog"
	"strconv"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrFileAlreadyExists = errors.New("file already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrFileNotFound      = errors.New("file not found")
)

type InfoDataSaver interface {
	SaveInfoData(
		ctx context.Context,
		data string,
		infoType string,
		additionalData string,
		uuid string,
		userID int64,
	) (err error)
}

type InfoData interface {
	InfoData(ctx context.Context, uuid string, userID int64) (dataModel *model.InfoDataModel, err error)
}

type FileSaver interface {
	SaveFile(
		ctx context.Context, filename string, uuid string, userID int64) (err error)
}

type File interface {
	File(ctx context.Context, uuid string, userID int64) (file *model.FileModel, err error)
}

type UserSaver interface {
	SaveUser(
		ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (user *model.User, err error)
}

type UserData interface {
	AllData(ctx context.Context, userID int64) (commonData []*model.CommonData, err error)
}

type Storage struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

// Инициализация хранилища.
func NewStorage(
	log *slog.Logger,
	storagePath string,
) (*Storage, error) {
	const tag = "storage.CreateStorage"
	logTag := log.With(slog.String("tag", tag))

	poolConfig, err := pgxpool.ParseConfig(storagePath)
	if err != nil {
		logTag.Info("Unable to parse DATABASE_URL:", sl.Err(err))
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		logTag.Info("Unable to create connection pool:", sl.Err(err))
		return nil, err
	}

	return &Storage{
		log:  log,
		pool: db,
	}, nil
}

// Закрытие хранилища.
func (s *Storage) Close() {
	s.pool.Close()
}

func (s *Storage) SaveUser(
	ctx context.Context,
	email string,
	passHash []byte,
) (uid int64, err error) {
	const tag = "storage.SaveUser"
	log := s.log.With(slog.String("tag", tag))

	row := s.pool.QueryRow(ctx, saveUserQuery, email, passHash)
	var userID int
	err = row.Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Info("user already exists")
			return 0, ErrUserAlreadyExists
		}
		log.Info("Can't save url in table " + err.Error())
		return 0, err
	}

	return int64(userID), nil
}

func (s *Storage) SaveInfoData(
	ctx context.Context,
	data string,
	infoType string,
	additionalData string,
	uuid string,
	userID int64,
) (err error) {
	const tag = "storage.SaveInfoData"
	log := s.log.With(slog.String("tag", tag))
	log.Info("Save data " + infoType + " user " + strconv.FormatInt(userID, 10))

	row := s.pool.QueryRow(ctx, saveInfoDataQuery, uuid, data, additionalData, infoType, userID)
	var file string
	err = row.Scan(&file)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Info("data already exists")
			return ErrFileAlreadyExists
		}
		log.Info("Can't save data in table " + err.Error())
		return err
	}
	log.Info("Save data in table " + file)
	return nil
}

func (s *Storage) SaveFile(
	ctx context.Context,
	filename string,
	uuid string,
	userID int64,
) (err error) {
	const tag = "storage.SaveFile"
	log := s.log.With(slog.String("tag", tag))
	log.Info("Save data " + filename + " user " + strconv.FormatInt(userID, 10))

	row := s.pool.QueryRow(ctx, saveFileQuery, uuid, filename, userID)
	var file string
	err = row.Scan(&file)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Info("file already exists")
			return ErrFileAlreadyExists
		}
		log.Info("Can't save file in table " + err.Error())
		return err
	}
	log.Info("Save file in table " + file)
	return nil
}

func (s *Storage) InfoData(ctx context.Context, uuid string, userID int64) (dataModel *model.InfoDataModel, err error) {
	const tag = "storage.InfoData"
	log := s.log.With(slog.String("tag", tag))
	log.Info("get InfoData " + uuid)

	row := s.pool.QueryRow(ctx, getInfoDataQuery, uuid)
	f := model.InfoDataModel{}
	err = row.Scan(&f.UUID, &f.Data, &f.AdditionalData, &f.Type, &f.UserID)
	if err != nil {
		log.Info("Can't find info data in table " + err.Error())
		return nil, err
	}
	if f.UserID != userID {
		log.Info("Can't find info data in table for user with id" + strconv.Itoa(int(userID)) + err.Error())
		return nil, ErrFileNotFound
	}

	return &f, nil
}

func (s *Storage) File(ctx context.Context, uuid string, userID int64) (file *model.FileModel, err error) {
	const tag = "storage.SaveFile"
	log := s.log.With(slog.String("tag", tag))
	log.Info("Get data " + uuid + " user " + strconv.FormatInt(userID, 10))

	row := s.pool.QueryRow(ctx, getFileQuery, uuid)
	f := model.FileModel{}
	err = row.Scan(&f.UUID, &f.Filename, &f.UserID)
	if err != nil {
		log.Info("Can't find file in table " + err.Error())
		return nil, err
	}
	if f.UserID != userID {
		log.Info("Can't find file in table for user with id" + strconv.Itoa(int(userID)) + err.Error())
		return nil, ErrFileNotFound
	}

	return &f, nil
}

func (s *Storage) User(
	ctx context.Context,
	email string,
) (user *model.User, err error) {
	const tag = "storage.User"
	log := s.log.With(slog.String("tag", tag))

	row := s.pool.QueryRow(ctx, getUserQuery, email)
	u := model.User{}
	err = row.Scan(&u.ID, &u.Email, &u.PasswordHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				log.Info("user already exists")
				return nil, ErrUserAlreadyExists
			default:
				log.Info("error " + pgErr.Error())
				return nil, ErrUserNotFound
			}
		}
		log.Info("Can't find user in table " + err.Error())
		return nil, fmt.Errorf("unable to scan row: %w", err)
	}

	return &u, nil
}

func (s *Storage) AllData(ctx context.Context, userID int64) (commonData []*model.CommonData, err error) {
	const tag = "storage.AllData"
	log := s.log.With(slog.String("tag", tag))
	log.Info("get AllData")

	files, err := s.AllFile(ctx, userID)
	if err != nil {
		log.Info("Can't find files in table " + err.Error())
		return nil, err
	}

	infoDatas, err := s.AllInfoData(ctx, userID)
	if err != nil {
		log.Info("Can't find info data in table " + err.Error())
		return nil, err
	}

	return append(files, infoDatas...), nil
}

func (s *Storage) AllFile(ctx context.Context, userID int64) (commonData []*model.CommonData, err error) {
	const tag = "storage.AllFile"
	log := s.log.With(slog.String("tag", tag))
	log.Info("get AllFile")

	rows, err := s.pool.Query(ctx, getFileByUserQuery, userID)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to get files: %w", err)
	}

	files := []model.FileModel{}
	for rows.Next() {
		f := model.FileModel{}
		err := rows.Scan(&f.UUID, &f.Filename, &f.UserID)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		files = append(files, f)
	}

	return mapper.MapFileToCommon(files), nil
}

func (s *Storage) AllInfoData(ctx context.Context, userID int64) (commonData []*model.CommonData, err error) {
	const tag = "storage.AllInfoData"
	log := s.log.With(slog.String("tag", tag))
	log.Info("get AllInfoData")

	rows, err := s.pool.Query(ctx, getInfoDataByUserQuery, userID)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to get info data: %w", err)
	}
	infoData := []model.InfoDataModel{}
	for rows.Next() {
		info := model.InfoDataModel{}
		err := rows.Scan(&info.UUID, &info.Data, &info.AdditionalData, &info.Type, &info.UserID)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		infoData = append(infoData, info)
	}

	return mapper.MapInfoDataToCommon(infoData), nil
}

const (
	saveInfoDataQuery      = `INSERT INTO info_data VALUES ($1, $2, $3, $4, $5) returning (uuid);`
	getInfoDataQuery       = `SELECT uuid, data, additional, type, user_id FROM info_data WHERE uuid = $1;`
	getInfoDataByUserQuery = `SELECT uuid, data, additional, type, user_id FROM info_data WHERE user_id = $1;`
	saveFileQuery          = `INSERT INTO files VALUES ($1, $2, $3) returning (uuid);`
	getFileQuery           = `SELECT uuid, filename, user_id FROM files WHERE uuid = $1;`
	getFileByUserQuery     = `SELECT uuid, filename, user_id FROM files WHERE user_id = $1;`
	saveUserQuery          = `INSERT INTO users VALUES (DEFAULT, $1, $2) returning (id);`
	getUserQuery           = `SELECT id, email, pass_hash FROM users WHERE email = $1`
)
