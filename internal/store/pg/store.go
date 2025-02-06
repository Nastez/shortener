package pg

import (
	"context"
	"database/sql"
)

// Store реализует интерфейс store.Store и позволяет взаимодействовать с СУБД PostgreSQL
type Store struct {
	// Поле conn содержит объект соединения с СУБД
	conn *sql.DB
}

//func (s Store) Save(originalURL string, generatedID string) {
//	//TODO implement me
//	panic("implement me")
//}

//func (s Store) Get(urlID string) string {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s Store) Save(originalURL string, generatedID string) {
//	s[generatedID] = originalURL
//}

//func (s Store) Get(urlID string) string {
//	var originalURL = s[urlID]
//
//	return originalURL
//}

// NewStore возвращает новый экземпляр PostgreSQL-хранилища
func NewStore(conn *sql.DB) *Store {
	return &Store{conn: conn}
}

type MemoryStorage map[string]string

type URLStorage interface {
	GetOriginalURL(ctx context.Context, id string) (string, error)
	Save(ctx context.Context, originalURL string, shortURL string, generatedID string) error
	Bootstrap(ctx context.Context) error
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

func (s Store) GetOriginalURL(ctx context.Context, id string) (string, error) {
	// запрашиваем originalURL по сгенерированному id
	row := s.conn.QueryRowContext(ctx, `
        SELECT
            u.original_url
        FROM urls u
        WHERE
            u.url_id = $1
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

func (s Store) Save(ctx context.Context, originalURL string, shortURL string, generatedID string) error {
	// добавляем новую запись с URLs в БД
	_, err := s.conn.ExecContext(ctx, `
        INSERT INTO urls
        (original_url, short_url, url_id)
        VALUES
        ($1, $2, $3)
    `, originalURL, shortURL, generatedID)

	return err
}

//func (m MemoryStorage) Save(originalURL string, generatedID string) {
//	m[generatedID] = originalURL
//}
//
//func (m MemoryStorage) Get(urlID string) string {
//	var originalURL = m[urlID]
//
//	return originalURL
//}
