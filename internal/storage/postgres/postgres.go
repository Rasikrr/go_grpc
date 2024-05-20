package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"sso/internal/config"
	"sso/internal/domain/models"
	"sso/internal/storage"
)

type Storage struct {
	db *sqlx.DB
}

func New(storageCfg *config.StorageConfig) (*Storage, error) {
	const op = "postgres.New"

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		storageCfg.Host, storageCfg.Port, storageCfg.User, storageCfg.Dbname,
		storageCfg.SslMode, storageCfg.Password,
	)

	db, err := sqlx.Connect("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil

}

func (s *Storage) SaveUser(ctx context.Context, email, passHash string) (int64, error) {
	const op = "storage.postgres.SaveUser"
	var lastInsertedId int64
	err := s.db.QueryRowx(
		`INSERT INTO users(email, pass_hash) VALUES($1, $2) RETURNING id`,
		email,
		passHash,
	).Scan(&lastInsertedId)

	if err != nil {
		var psqlError pq.Error

		if errors.Is(err, &psqlError) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return lastInsertedId, nil

}

func (s *Storage) User(ctx context.Context, email string) (*models.User, error) {
	const op = "storage.postgres.User"

	user := new(models.User)
	err := s.db.QueryRowx(`SELECT * FROM users WHERE email=$1`, email).StructScan(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	user := new(models.User)
	err := s.db.QueryRowx(`SELECT * FROM users WHERE id=$1`, userID).StructScan(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return user.IsAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (*models.App, error) {
	const op = "storage.postgres.App"

	app := new(models.App)
	err := s.db.QueryRowx(`SELECT * FROM apps WHERE id=$1`, appID).StructScan(app)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}
