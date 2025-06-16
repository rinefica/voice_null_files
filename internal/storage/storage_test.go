package storage

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rinefica/voice_null_files/internal/lib/sl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDbInstance *pgxpool.Pool

func TestMain(m *testing.M) {
	testDB := SetupTestDatabase()
	testDbInstance = testDB.DbInstance
	defer testDB.TearDown()
	os.Exit(m.Run())
}

func TestCreateUser(t *testing.T) {
	ds := Storage{
		sl.SetupLogger(""),
		testDbInstance,
	}

	email := "test@mail.co"
	userID, err := ds.SaveUser(context.Background(), email, []byte("jkjkjk"))

	log.Println(userID)

	assert.NotNil(t, userID)
	assert.NoError(t, err)

	id := int64(1)
	assert.Equal(t, id, userID)

	user, err := ds.User(context.Background(), email)
	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, id, user.ID)

	user, err = ds.User(context.Background(), email+"123")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestAddFile(t *testing.T) {
	secret := "secret"
	ds := Storage{
		sl.SetupLogger(""),
		testDbInstance,
	}

	userID := doAuth(t, secret, ds)

	err := ds.SaveFile(context.Background(), "filename", "uuid", userID)
	assert.NoError(t, err)

	file, err := ds.File(context.Background(), "uuid", userID)

	assert.NoError(t, err)
	assert.Equal(t, "filename", file.Filename)
}

func TestSaveInfoData(t *testing.T) {
	secret := "secret"
	ds := Storage{
		sl.SetupLogger(""),
		testDbInstance,
	}

	userID := doAuth(t, secret, ds)
	data := "password"
	infoType := "txt"
	additional := "txt"
	uuid := "uuid"

	err := ds.SaveInfoData(context.Background(), data, infoType, additional, uuid, userID)
	assert.NoError(t, err)

	infoData, err := ds.InfoData(context.Background(), uuid, userID)

	assert.NoError(t, err)
	assert.Equal(t, data, infoData.Data)
}

func TestAllFiles(t *testing.T) {
	secret := "secret"
	ds := Storage{
		sl.SetupLogger(""),
		testDbInstance,
	}

	userID := doAuth(t, secret, ds)

	err := ds.SaveFile(context.Background(), "filename1", "uuid1", userID)
	assert.NoError(t, err)

	err = ds.SaveFile(context.Background(), "filename2", "uuid2", userID)
	assert.NoError(t, err)

	err = ds.SaveFile(context.Background(), "filename3", "uuid3", userID)
	assert.NoError(t, err)

	files, err := ds.AllFile(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, files, 3)
}

func doAuth(t *testing.T, secret string, ds Storage) int64 {
	email := "test@mail.co"
	userID, err := ds.SaveUser(context.Background(), email, []byte(secret))

	log.Println(userID)
	require.NoError(t, err)

	id := int64(1)
	assert.Equal(t, id, userID)

	user, err := ds.User(context.Background(), email)
	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, id, user.ID)

	return userID
}
