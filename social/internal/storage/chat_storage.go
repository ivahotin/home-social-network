package storage

import (
	"context"
	"database/sql"
	"example.com/social/internal/domain"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	saveMessageStmt = "insert into messages (chat_id, message_id, published_at, message, author_id) values ($1, $2, current_timestamp, $3, $4)"
)

type CockroachChatStorage struct {
	db *pgxpool.Pool
}

func NewCockroachChatStorage(db *pgxpool.Pool) *CockroachChatStorage {
	return &CockroachChatStorage{
		db: db,
	}
}

func (storage *CockroachChatStorage) SaveMessage(message *domain.Message) error {
	row := storage.db.QueryRow(
		context.Background(),
		saveMessageStmt,
		message.ChatId,
		message.MessageId,
		message.Message,
		message.AuthorId)
	err := row.Scan()
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil
		}
		return err
	}
	return nil
}

const (
	insertMessage = "insert into messages (chat_id, message_id, published_at, message, author_id) values (?, ?, current_timestamp, ?, ?)"
)

type MySqlChatStorage struct {
	db *sql.DB
}

func NewMySqlChatStorage(db *sql.DB) *MySqlChatStorage {
	return &MySqlChatStorage{
		db: db,
	}
}

func (storage *MySqlChatStorage) SaveMessage(message *domain.Message) error {
	stmt, err := storage.db.Prepare(insertMessage)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(message.ChatId, message.MessageId, message.Message, message.AuthorId)
	if err != nil {
		return err
	}

	return nil
}