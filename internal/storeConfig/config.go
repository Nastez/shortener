package storeConfig

import (
	"context"
	"database/sql"
)

// StoreConfig реализует интерфейс store.Store и позволяет взаимодействовать с СУБД PostgreSQL
type StoreConfig struct {
	// Поле conn содержит объект соединения с СУБД
	conn *sql.DB
}

func NewStoreConfig(conn *sql.DB) *StoreConfig {
	return &StoreConfig{conn: conn}
}

// Bootstrap подготавливает БД к работе, создавая необходимые таблицы и индексы
func (s StoreConfig) Bootstrap(ctx context.Context) error {
	// запускаем транзакцию
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// в случае неуспешного коммита все изменения транзакции будут отменены
	defer tx.Rollback()

	// создаём таблицу urls и необходимые индексы
	tx.ExecContext(ctx, `
       CREATE TABLE if NOT EXISTS urls (
           id SERIAL PRIMARY KEY,
           original_url text UNIQUE,
           short_url text,
           url_id text
       )
    `)

	tx.ExecContext(ctx, `CREATE INDEX url_idx ON urls (url_id)`)

	// коммитим транзакцию
	return tx.Commit()
}
