package migration

import (
	"log"
	"microService/config"
	"microService/modules/order/domain"
	"microService/pkg/database"
	"microService/pkg/outbox"
)

func Order(cfg *config.Config) {

	db := database.NewPostgresDatabase(cfg.Postgres).Connect()

	if err := db.AutoMigrate(&domain.Order{}, &domain.OrderItem{}, &outbox.Outbox{}); err != nil {
		panic(err)
	}
	//sigle นะจ๊ะ
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_order_user_status ON orders (user_id, status)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_order_user_created ON orders (user_id, created_at DESC)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_order_status_created ON orders (status, created_at DESC)`)

	// OrderItem table composite indexes
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_order_item_sku_order ON order_items (sku, order_id)`)

	db.Exec(`CREATE INDEX IF NOT EXISTS idx_outbox_status_created ON outbox (status, created_at ASC)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_outbox_aggregate_event ON outbox (aggregate, event_type)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_outbox_key_status ON outbox (key, status)`)

	log.Println("✅ Migration completed successfully with optimized indexes")
	log.Println("✅ Migration completed successfully")
}
