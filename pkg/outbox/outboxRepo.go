package outbox

import (
	"context"
	"microService/pkg/database"
	"time"

	"gorm.io/gorm"
)

type OutboxRepository interface {
	Add(ctx context.Context, ob *Outbox) error
	PollPending(ctx context.Context, limit int) ([]Outbox, error)
	PollPendingLocked(ctx context.Context, limit int) ([]Outbox, error)
	MarkDispatched(ctx context.Context, id int64) error
}

type OutboxRepo struct {
	db database.DatabasesPostgres
}

func NewOutboxRepo(db database.DatabasesPostgres) OutboxRepository {
	return &OutboxRepo{db: db}
}

func (r *OutboxRepo) Add(ctx context.Context, ob *Outbox) error {
	db := database.GetDB(ctx, r.db.Connect()).WithContext(ctx)
	return db.Create(ob).Error
}

// แบบไม่ lock (single worker พอได้)
func (r *OutboxRepo) PollPending(ctx context.Context, limit int) ([]Outbox, error) {
	db := database.GetDB(ctx, r.db.Connect()).WithContext(ctx)
	var res []Outbox
	if err := db.Where("status = ?", OutboxStatusPending).
		Order("id ASC").
		Limit(limit).
		Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// แบบ lock ด้วย SKIP LOCKED (แนะนำเวลา scale หลาย worker)
// ควรเรียกใน Transaction ภายนอก
func (r *OutboxRepo) PollPendingLocked(ctx context.Context, limit int) ([]Outbox, error) {
	base := r.db.Connect()
	var res []Outbox

	err := base.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		rows, err := tx.Raw(`
			SELECT id, aggregate, event_type, key, payload, status, created_at, updated_at
			  FROM outbox
			 WHERE status = 'pending'
			 ORDER BY id ASC
			 LIMIT ?
			 FOR UPDATE SKIP LOCKED
		`, limit).Rows()
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var ob Outbox
			if err := rows.Scan(&ob.ID, &ob.Aggregate, &ob.EventType, &ob.Key, &ob.Payload, &ob.Status, &ob.CreatedAt, &ob.UpdatedAt); err != nil {
				return err
			}
			res = append(res, ob)
		}
		return rows.Err()
	})
	return res, err
}

func (r *OutboxRepo) MarkDispatched(ctx context.Context, id int64) error {
	db := database.GetDB(ctx, r.db.Connect()).WithContext(ctx)
	return db.Model(&Outbox{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":     OutboxStatusDispatched,
			"updated_at": time.Now().UTC(),
		}).Error
}
