package orderRepo

import (
	"context"
	"microService/modules/order/domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func getDB(ctx context.Context, pool *pgxpool.Pool) DB {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return pool
}

type OrderRepository interface {
	Tx(ctx context.Context, fn func(ctx context.Context) error) error
	InsertOrder(ctx context.Context, o *domain.Order) error
}

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	ctx = context.WithValue(ctx, txKey{}, tx)

	if err := fn(ctx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *Repo) InsertOrder(ctx context.Context, o *domain.Order) error {
	db := getDB(ctx, r.db)

	_, err := db.Exec(ctx,
		`INSERT INTO orders(id, user_id, amount, currency, status, created_at, updated_at)
		 VALUES($1,$2,$3,$4,$5,$6,$7)`,
		o.ID, o.UserID, o.Amount, o.Currency, o.Status, o.CreatedAt, o.UpdatedAt)
	if err != nil {
		return err
	}

	for _, it := range o.Items {
		_, err = db.Exec(ctx,
			`INSERT INTO order_items(order_id, sku, qty, price) VALUES ($1,$2,$3,$4)`,
			o.ID, it.SKU, it.Qty, it.Price)
		if err != nil {
			return err
		}
	}
	return nil
}

// ---------- Outbox Repo ----------
type OutboxRepository interface {
	Add(ctx context.Context, ob *domain.Outbox) error
	PollPending(ctx context.Context, limit int) ([]domain.Outbox, error)
	MarkDispatched(ctx context.Context, id int64) error
}

type OutboxRepo struct {
	db *pgxpool.Pool
}

func NewOutboxRepo(db *pgxpool.Pool) OutboxRepository {
	return &OutboxRepo{db: db}
}

func (r *OutboxRepo) Add(ctx context.Context, ob *domain.Outbox) error {
	db := getDB(ctx, r.db)
	_, err := db.Exec(ctx,
		`INSERT INTO outbox(aggregate, event_type, key, payload, status, created_at, updated_at)
		 VALUES($1,$2,$3,$4,$5,$6,$7)`,
		ob.Aggregate, ob.EventType, ob.Key, ob.Payload, ob.Status, ob.CreatedAt, ob.UpdatedAt)
	return err
}

func (r *OutboxRepo) PollPending(ctx context.Context, limit int) ([]domain.Outbox, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, aggregate, event_type, key, payload
		 FROM outbox
		 WHERE status='pending'
		 ORDER BY id ASC
		 LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.Outbox
	for rows.Next() {
		var ob domain.Outbox
		var b []byte
		if err := rows.Scan(&ob.ID, &ob.Aggregate, &ob.EventType, &ob.Key, &b); err != nil {
			return nil, err
		}
		ob.Payload = b
		res = append(res, ob)
	}
	return res, rows.Err()
}

func (r *OutboxRepo) MarkDispatched(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE outbox
		 SET status='dispatched', updated_at=$2
		 WHERE id=$1`, id, time.Now().UTC())
	return err
}
