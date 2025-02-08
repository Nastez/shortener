package pg

import (
	"context"
	"database/sql"
	"github.com/Nastez/shortener/internal/store"
)

// Store реализует интерфейс store.Store и позволяет взаимодействовать с СУБД PostgreSQL
type Store struct {
	// Поле conn содержит объект соединения с СУБД
	conn *sql.DB
}

// NewStore возвращает новый экземпляр PostgreSQL-хранилища
func NewStore(conn *sql.DB) *Store {
	return &Store{conn: conn}
}

// Bootstrap подготавливает БД к работе, создавая необходимые таблицы и индексы
func (s Store) Bootstrap(ctx context.Context) error {
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
            original_url text,
            short_url text,
            url_id text
        )
    `)

	tx.ExecContext(ctx, `CREATE INDEX url_idx ON urls (url_id)`)

	// коммитим транзакцию
	return tx.Commit()
}

func (s Store) Get(ctx context.Context, id string) (string, error) {
	// запрашиваем originalURL по сгенерированному id
	row := s.conn.QueryRowContext(ctx, `
        SELECT
            original_url
        FROM urls 
        WHERE
            url_id = $1
    `,
		id,
	)

	// считываем значения из записи БД в соответствующие поля структуры
	var originalURL string
	err := row.Scan(&originalURL) // разбираем результат
	if err != nil {
		return "", err
	}

	return originalURL, nil
}

func (s Store) Save(ctx context.Context, urls store.URL) error {
	// добавляем новую запись с URLs в БД
	_, err := s.conn.ExecContext(ctx, `
        INSERT INTO urls
        (original_url, short_url, url_id)
        VALUES
        ($1, $2, $3)
    `, urls.OriginalURL, urls.ShortURL, urls.GeneratedID)

	return err
}
