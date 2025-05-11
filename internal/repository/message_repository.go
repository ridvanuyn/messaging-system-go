package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ridvanuyn/messaging-system-go/internal/domain"
	"github.com/jmoiron/sqlx"
)

type MessageRepository interface {
	GetUnsentMessages(ctx context.Context, limit int) ([]domain.Message, error)
	MarkAsSent(ctx context.Context, id int64, messageID string) error
	GetSentMessages(ctx context.Context) ([]domain.Message, error)
	CacheMessageID(ctx context.Context, messageID string, sentTime time.Time) error
}

type postgresMessageRepository struct {
	db    *sqlx.DB
	redis *redis.Client
}

func NewMessageRepository(db *sqlx.DB, redis *redis.Client) MessageRepository {
	return &postgresMessageRepository{
		db:    db,
		redis: redis,
	}
}

func (r *postgresMessageRepository) GetUnsentMessages(ctx context.Context, limit int) ([]domain.Message, error) {
	var messages []domain.Message
	query := `SELECT id, "to", content, sent, created_at FROM messages WHERE sent = false ORDER BY created_at LIMIT $1`
	
	err := r.db.SelectContext(ctx, &messages, query, limit)
	if err != nil {
		return nil, err
	}
	
	return messages, nil
}

func (r *postgresMessageRepository) MarkAsSent(ctx context.Context, id int64, messageID string) error {
	query := `UPDATE messages SET sent = true, sent_at = NOW(), message_id = $1 WHERE id = $2`
	
	_, err := r.db.ExecContext(ctx, query, messageID, id)
	return err
}

func (r *postgresMessageRepository) GetSentMessages(ctx context.Context) ([]domain.Message, error) {
	var messages []domain.Message
	query := `SELECT id, "to", content, sent, created_at, sent_at, message_id FROM messages WHERE sent = true ORDER BY sent_at DESC`
	
	err := r.db.SelectContext(ctx, &messages, query)
	if err != nil {
		return nil, err
	}
	
	return messages, nil
}

func (r *postgresMessageRepository) CacheMessageID(ctx context.Context, messageID string, sentTime time.Time) error {
	key := "message:" + messageID
	value := sentTime.Format(time.RFC3339)
	
	return r.redis.Set(ctx, key, value, 24*time.Hour).Err()
}
