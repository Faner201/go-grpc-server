package postrgres

import (
	"context"
	"errors"
	"fmt"

	dm_models "github.com/Faner201/go-grpc-server/sso/internal/domain/models"
	"github.com/Faner201/go-grpc-server/sso/internal/storage"
	"github.com/Faner201/go-grpc-server/sso/internal/storage/models"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func New(host, username, password, dbnmae, sslmode, timeZone, port string) (*Storage, error) {
	const op = "storage.postgres.New"

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s sslmode=%s TimeZone=%s", host, port, username, password, sslmode, timeZone,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var (
		user models.User
		app  models.App
	)

	db.AutoMigrate(&user, &app)
	db.Create(&models.App{
		Name:   "test",
		Secret: "secret",
	})

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	user := models.User{
		Email:    email,
		PassHash: passHash,
	}

	result := s.db.Create(&user)

	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, result.Error)
	}

	result.Last(&user)

	return int64(user.ID), nil
}

func (s *Storage) User(ctx context.Context, email string) (dm_models.User, error) {
	const op = "storage.postgres.User"

	var user_db models.User

	err := s.db.Where("email = ?", email).First(&user_db).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dm_models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return dm_models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return dm_models.User{
		ID:       int64(user_db.ID),
		Email:    user_db.Email,
		PassHash: user_db.PassHash,
	}, nil
}

func (s *Storage) App(ctx context.Context, id int) (dm_models.App, error) {
	const op = "storage.postgres.App"

	var app_db models.App

	err := s.db.Where("id = ?", id).First(&app_db).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dm_models.App{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return dm_models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return dm_models.App{
		ID:     int(app_db.ID),
		Name:   app_db.Name,
		Secret: app_db.Secret,
	}, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	var isAdmin bool

	err := s.db.Where("id = ?", userID).First(&isAdmin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
