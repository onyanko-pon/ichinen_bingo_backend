package sql_handler

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type SQLHandler struct {
	db *sql.DB
}

func NewHandler(dataSource string) (*SQLHandler, error) {
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, fmt.Errorf("init db connection: %w", err)
	}

	ctx := context.Background()
	// TODO タイムアウト時間を定数化
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("init db connection: %w", err)
	}

	return &SQLHandler{db: db}, nil
}

func (h *SQLHandler) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	log.Printf("[sql handler] QueryContext, query: %s, args: %v", strings.ReplaceAll(query, "\n", " "), args)
	// ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()
	return h.db.QueryContext(ctx, query, args...)
}

func (h *SQLHandler) CleanData(ctx context.Context) error {
	bytes, err := ioutil.ReadFile(os.Getenv("PROJECT_ROOT") + "/db/clean.sql")
	if err != nil {
		panic(err)
	}

	query := string(bytes)
	_, err = h.db.QueryContext(ctx, query)
	return err
}
