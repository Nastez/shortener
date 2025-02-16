package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Nastez/shortener/internal/app/models"
	"github.com/Nastez/shortener/internal/logger"
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

	// считываем значения из записи БД
	var originalURL string
	err := row.Scan(&originalURL) // разбираем результат
	if err != nil {
		return "", err
	}

	return originalURL, nil
}

func (s Store) Save(ctx context.Context, urls store.URL, userID string) (string, error) {
	// добавляем новую запись с URLs в БД
	res, err := s.conn.ExecContext(ctx, `
        INSERT INTO urls (original_url, short_url, url_id, user_id)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (original_url) DO NOTHING
    `, urls.OriginalURL, urls.ShortURL, urls.GeneratedID, userID)
	if err != nil {
		return "", fmt.Errorf("insert error: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("RowsAffected error: %w", err)
	}

	if rowsAffected == 0 {
		// проверяем, что ошибка сигнализирует о потенциальном нарушении целостности данных
		dataConflictErr := store.ErrConflict
		row := s.conn.QueryRowContext(ctx, `
			   SELECT
			       short_url
			   FROM urls
			   WHERE
			       original_url = $1
			`,
			urls.OriginalURL,
		)
		// считываем значения из записи БД в соответствующие поля структуры
		var oldShortURL string
		err = row.Scan(&oldShortURL) // разбираем результат
		if err != nil {
			logger.Log.Error("scan error")
		}

		return oldShortURL, dataConflictErr

	}

	return "", err
}

func (s Store) SaveBatch(ctx context.Context, requestBatch models.PayloadBatch, shortURLBatch models.ResponseBodyBatch, userID string) error {
	// запускаем транзакцию
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// в случае неуспешного коммита все изменения транзакции будут отменены
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO urls (short_url, url_id, user_id) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmtOriginalURL, err := tx.PrepareContext(ctx,
		"UPDATE urls SET original_url = $1 WHERE url_id = $2")
	if err != nil {
		return err
	}
	defer stmtOriginalURL.Close()

	for _, b := range shortURLBatch {
		_, err = stmt.ExecContext(ctx, b.ShortURL, b.CorrelationID, userID)
		if err != nil {
			return err
		}
		for _, req := range requestBatch {
			_, err = stmtOriginalURL.ExecContext(ctx, req.OriginalURL, req.CorrelationID)
			if err != nil {
				return err

			}
		}
	}

	// коммитим транзакцию
	return tx.Commit()
}

func (s Store) GetURLs(ctx context.Context, userID string) (models.URLSResponseArr, error) {
	// запрашиваем список с originalURL и shortURL по userID
	rows, err := s.conn.Query(`
       SELECT
           original_url, short_url
       FROM urls
       WHERE
           user_id = $1
   `,
		userID,
	)
	if err != nil {
		fmt.Println(err)
		logger.Log.Error("get urls from pg store error")
		return nil, err
	}
	defer rows.Close()

	// считываем значения из записи БД в соответствующие поля структуры
	var urls models.URLSResponse
	var urlsArr []models.URLSResponse
	for rows.Next() {
		err = rows.Scan(&urls.OriginalURL, &urls.ShortURL) // разбираем результат
		if err != nil {
			fmt.Println(err)
			logger.Log.Error("scan error")
		}
		urlsArr = append(urlsArr, urls)
	}
	if err = rows.Err(); err != nil {
		logger.Log.Error("get rows with urls error")
		return nil, err
	}

	return urlsArr, nil
}
