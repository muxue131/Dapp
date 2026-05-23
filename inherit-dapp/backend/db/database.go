package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Database wraps the SQL database connection
type Database struct {
	conn *sql.DB
}

// PlanRecord represents a plan stored in the database
type PlanRecord struct {
	ID                int64     `json:"id"`
	PlanID            uint64    `json:"plan_id"`
	CreatorAddress    string    `json:"creator_address"`
	HeartbeatInterval int64     `json:"heartbeat_interval"`
	LastHeartbeat     time.Time `json:"last_heartbeat"`
	TriggerTime       time.Time `json:"trigger_time"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// AssetRecord represents an asset stored in the database
type AssetRecord struct {
	ID            int64     `json:"id"`
	AssetID       uint64    `json:"asset_id"`
	PlanID        uint64    `json:"plan_id"`
	OwnerAddress  string    `json:"owner_address"`
	AssetType     string    `json:"asset_type"`
	Denom         string    `json:"denom"`
	Amount        string    `json:"amount"`
	IPFSCid       string    `json:"ipfs_cid"`
	Metadata      string    `json:"metadata"`
	CreatedAt     time.Time `json:"created_at"`
}

// WitnessRecord represents a witness/validator for heartbeat monitoring
type WitnessRecord struct {
	ID             int64     `json:"id"`
	PlanID         uint64    `json:"plan_id"`
	WitnessAddress string    `json:"witness_address"`
	AddedAt        time.Time `json:"added_at"`
}

// NewDatabase creates a new database connection
func NewDatabase(dsn string) (*Database, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(5 * time.Minute)

	db := &Database{conn: conn}
	if err := db.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// Migrate runs database migrations
func (db *Database) Migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS plans (
			id SERIAL PRIMARY KEY,
			plan_id BIGINT UNIQUE NOT NULL,
			creator_address VARCHAR(128) NOT NULL,
			heartbeat_interval BIGINT NOT NULL,
			last_heartbeat TIMESTAMP WITH TIME ZONE NOT NULL,
			trigger_time TIMESTAMP WITH TIME ZONE NOT NULL,
			status VARCHAR(32) NOT NULL DEFAULT 'active',
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_plans_creator ON plans(creator_address)`,
		`CREATE INDEX IF NOT EXISTS idx_plans_status ON plans(status)`,
		`CREATE INDEX IF NOT EXISTS idx_plans_trigger_time ON plans(trigger_time)`,

		`CREATE TABLE IF NOT EXISTS assets (
			id SERIAL PRIMARY KEY,
			asset_id BIGINT UNIQUE NOT NULL,
			plan_id BIGINT NOT NULL REFERENCES plans(plan_id),
			owner_address VARCHAR(128) NOT NULL,
			asset_type VARCHAR(32) NOT NULL,
			denom VARCHAR(256) NOT NULL,
			amount VARCHAR(128) NOT NULL,
			ipfs_cid VARCHAR(256),
			metadata TEXT,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_assets_plan ON assets(plan_id)`,

		`CREATE TABLE IF NOT EXISTS witnesses (
			id SERIAL PRIMARY KEY,
			plan_id BIGINT NOT NULL REFERENCES plans(plan_id),
			witness_address VARCHAR(128) NOT NULL,
			added_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			UNIQUE(plan_id, witness_address)
		)`,

		`CREATE TABLE IF NOT EXISTS heartbeat_logs (
			id SERIAL PRIMARY KEY,
			plan_id BIGINT NOT NULL,
			sender_address VARCHAR(128) NOT NULL,
			tx_hash VARCHAR(128),
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_heartbeat_logs_plan ON heartbeat_logs(plan_id)`,
	}

	for _, migration := range migrations {
		if _, err := db.conn.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w\nSQL: %s", err, migration)
		}
	}

	return nil
}

// --- Plan Operations ---

func (db *Database) CreatePlan(ctx context.Context, record *PlanRecord) error {
	query := `INSERT INTO plans (plan_id, creator_address, heartbeat_interval, last_heartbeat, trigger_time, status)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	return db.conn.QueryRowContext(ctx, query,
		record.PlanID, record.CreatorAddress, record.HeartbeatInterval,
		record.LastHeartbeat, record.TriggerTime, record.Status,
	).Scan(&record.ID, &record.CreatedAt, &record.UpdatedAt)
}

func (db *Database) GetPlan(ctx context.Context, planID uint64) (*PlanRecord, error) {
	query := `SELECT id, plan_id, creator_address, heartbeat_interval, last_heartbeat, trigger_time, status, created_at, updated_at
		FROM plans WHERE plan_id = $1`
	record := &PlanRecord{}
	err := db.conn.QueryRowContext(ctx, query, planID).Scan(
		&record.ID, &record.PlanID, &record.CreatorAddress, &record.HeartbeatInterval,
		&record.LastHeartbeat, &record.TriggerTime, &record.Status, &record.CreatedAt, &record.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return record, err
}

func (db *Database) GetPlansByCreator(ctx context.Context, creatorAddress string) ([]PlanRecord, error) {
	query := `SELECT id, plan_id, creator_address, heartbeat_interval, last_heartbeat, trigger_time, status, created_at, updated_at
		FROM plans WHERE creator_address = $1 ORDER BY created_at DESC`
	rows, err := db.conn.QueryContext(ctx, query, creatorAddress)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []PlanRecord
	for rows.Next() {
		var p PlanRecord
		if err := rows.Scan(&p.ID, &p.PlanID, &p.CreatorAddress, &p.HeartbeatInterval,
			&p.LastHeartbeat, &p.TriggerTime, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, nil
}

func (db *Database) UpdatePlanStatus(ctx context.Context, planID uint64, status string) error {
	query := `UPDATE plans SET status = $1, updated_at = NOW() WHERE plan_id = $2`
	_, err := db.conn.ExecContext(ctx, query, status, planID)
	return err
}

func (db *Database) UpdateHeartbeat(ctx context.Context, planID uint64, lastHeartbeat, triggerTime time.Time) error {
	query := `UPDATE plans SET last_heartbeat = $1, trigger_time = $2, updated_at = NOW() WHERE plan_id = $3`
	_, err := db.conn.ExecContext(ctx, query, lastHeartbeat, triggerTime, planID)
	return err
}

func (db *Database) GetExpiredPlans(ctx context.Context) ([]PlanRecord, error) {
	query := `SELECT id, plan_id, creator_address, heartbeat_interval, last_heartbeat, trigger_time, status, created_at, updated_at
		FROM plans WHERE status = 'active' AND trigger_time < NOW()`
	rows, err := db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []PlanRecord
	for rows.Next() {
		var p PlanRecord
		if err := rows.Scan(&p.ID, &p.PlanID, &p.CreatorAddress, &p.HeartbeatInterval,
			&p.LastHeartbeat, &p.TriggerTime, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, nil
}

// --- Asset Operations ---

func (db *Database) CreateAsset(ctx context.Context, record *AssetRecord) error {
	query := `INSERT INTO assets (asset_id, plan_id, owner_address, asset_type, denom, amount, ipfs_cid, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at`
	return db.conn.QueryRowContext(ctx, query,
		record.AssetID, record.PlanID, record.OwnerAddress, record.AssetType,
		record.Denom, record.Amount, record.IPFSCid, record.Metadata,
	).Scan(&record.ID, &record.CreatedAt)
}

func (db *Database) GetAssetsByPlan(ctx context.Context, planID uint64) ([]AssetRecord, error) {
	query := `SELECT id, asset_id, plan_id, owner_address, asset_type, denom, amount, ipfs_cid, metadata, created_at
		FROM assets WHERE plan_id = $1 ORDER BY created_at`
	rows, err := db.conn.QueryContext(ctx, query, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []AssetRecord
	for rows.Next() {
		var a AssetRecord
		if err := rows.Scan(&a.ID, &a.AssetID, &a.PlanID, &a.OwnerAddress, &a.AssetType,
			&a.Denom, &a.Amount, &a.IPFSCid, &a.Metadata, &a.CreatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, nil
}

// --- Heartbeat Log Operations ---

func (db *Database) LogHeartbeat(ctx context.Context, planID uint64, senderAddress, txHash string) error {
	query := `INSERT INTO heartbeat_logs (plan_id, sender_address, tx_hash) VALUES ($1, $2, $3)`
	_, err := db.conn.ExecContext(ctx, query, planID, senderAddress, txHash)
	return err
}

func (db *Database) GetHeartbeatLogs(ctx context.Context, planID uint64) ([]map[string]interface{}, error) {
	query := `SELECT id, plan_id, sender_address, tx_hash, timestamp FROM heartbeat_logs WHERE plan_id = $1 ORDER BY timestamp DESC`
	rows, err := db.conn.QueryContext(ctx, query, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []map[string]interface{}
	for rows.Next() {
		var id int64
		var pID uint64
		var sender, txHash string
		var ts time.Time
		if err := rows.Scan(&id, &pID, &sender, &txHash, &ts); err != nil {
			return nil, err
		}
		logs = append(logs, map[string]interface{}{
			"id": id, "plan_id": pID, "sender": sender, "tx_hash": txHash, "timestamp": ts,
		})
	}
	return logs, nil
}

// Close closes the database connection
func (db *Database) Close() error {
	return db.conn.Close()
}
