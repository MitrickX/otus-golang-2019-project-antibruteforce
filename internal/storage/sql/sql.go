package sql

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

const DefaultConnectRetries = 5

type Config struct {
	Host           string
	Port           string
	DbName         string
	User           string
	Password       string
	ConnectRetries int
}

func NewConfigByEnv() Config {
	retries, err := strconv.Atoi(os.Getenv("DB_CONNECT_RETRIES"))
	if err != nil {
		retries = DefaultConnectRetries
	}

	return Config{
		Host:           os.Getenv("DB_HOST"),
		Port:           os.Getenv("DB_PORT"),
		DbName:         os.Getenv("DB_DBNAME"),
		User:           os.Getenv("DB_USER"),
		Password:       os.Getenv("DB_PASSWORD"),
		ConnectRetries: retries,
	}
}

func Connect(cfg Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)

	db, err := sqlx.Open("pgx", dsn) // *sql.DB
	if err != nil {
		return nil, fmt.Errorf("failed to load driver %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var connectErr error
	for i := 0; i < cfg.ConnectRetries; i++ {
		connectErr = db.PingContext(ctx)
		if connectErr == nil {
			break
		}

		time.Sleep(time.Second)
	}

	if connectErr != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", connectErr)
	}

	return db, nil
}

func IsTableExists(ctx context.Context, db *sqlx.DB, dbName string, tableName string) (bool, error) {
	query := `SELECT EXISTS(
    	SELECT * 
    	FROM information_schema.tables 
    	WHERE 
      		table_schema = $1 AND 
      		table_name = $2
		)`

	var ok bool

	row := db.QueryRowxContext(ctx, query, dbName, tableName)

	err := row.Scan(&ok)
	if err != nil {
		return false, err
	}

	return ok, nil
}
