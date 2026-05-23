package monitor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/inherit-dapp/chain/backend/config"
	"github.com/inherit-dapp/chain/backend/db"
)

// HeartbeatMonitor monitors heartbeat expirations and triggers inheritance
type HeartbeatMonitor struct {
	config   *config.Config
	database *db.Database
	stopCh   chan struct{}
}

// NewHeartbeatMonitor creates a new heartbeat monitor
func NewHeartbeatMonitor(cfg *config.Config, database *db.Database) *HeartbeatMonitor {
	return &HeartbeatMonitor{
		config:   cfg,
		database: database,
		stopCh:   make(chan struct{}),
	}
}

// Start begins the heartbeat monitoring loop
func (m *HeartbeatMonitor) Start() {
	log.Println("Starting heartbeat monitor...")
	ticker := time.NewTicker(m.config.HeartbeatCheckInterval)
	defer ticker.Stop()

	// Run immediately on start
	m.checkExpiredPlans()

	for {
		select {
		case <-ticker.C:
			m.checkExpiredPlans()
		case <-m.stopCh:
			log.Println("Heartbeat monitor stopped.")
			return
		}
	}
}

// Stop stops the heartbeat monitor
func (m *HeartbeatMonitor) Stop() {
	close(m.stopCh)
}

// checkExpiredPlans checks for plans with expired heartbeats
func (m *HeartbeatMonitor) checkExpiredPlans() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	expiredPlans, err := m.database.GetExpiredPlans(ctx)
	if err != nil {
		log.Printf("Error checking expired plans: %v", err)
		return
	}

	if len(expiredPlans) == 0 {
		return
	}

	log.Printf("Found %d expired plans", len(expiredPlans))

	for _, plan := range expiredPlans {
		log.Printf("Plan %d heartbeat expired. Creator: %s, Last heartbeat: %s, Trigger time: %s",
			plan.PlanID,
			plan.CreatorAddress,
			plan.LastHeartbeat.Format(time.RFC3339),
			plan.TriggerTime.Format(time.RFC3339),
		)

		// Update plan status to triggered
		if err := m.database.UpdatePlanStatus(ctx, plan.PlanID, "triggered"); err != nil {
			log.Printf("Error updating plan %d status: %v", plan.PlanID, err)
			continue
		}

		// In production: submit MsgCheckHeartbeat to the chain
		m.triggerInheritance(ctx, plan)
	}
}

// triggerInheritance submits the heartbeat check transaction to the chain
func (m *HeartbeatMonitor) triggerInheritance(ctx context.Context, plan db.PlanRecord) {
	log.Printf("Triggering inheritance for plan %d...", plan.PlanID)

	// In production, this would:
	// 1. Create and sign a MsgCheckHeartbeat transaction
	// 2. Submit it to the Cosmos SDK chain via RPC
	// 3. Wait for confirmation
	// 4. Notify beneficiaries

	// For now, log the event
	log.Printf("Inheritance triggered for plan %d. Beneficiaries will be notified.", plan.PlanID)

	// Notify via webhook/callback (placeholder)
	m.notifyBeneficiaries(plan)
}

// notifyBeneficiaries sends notifications to beneficiaries
func (m *HeartbeatMonitor) notifyBeneficiaries(plan db.PlanRecord) {
	// In production:
	// 1. Query the chain for plan beneficiaries
	// 2. Send email/push notifications
	// 3. Post to webhook endpoints
	// 4. Emit events for frontend to pick up

	log.Printf("Notifying beneficiaries for plan %d (creator: %s)",
		plan.PlanID, plan.CreatorAddress)
}

// GetMonitorStatus returns the current monitor status
func (m *HeartbeatMonitor) GetMonitorStatus() map[string]interface{} {
	return map[string]interface{}{
		"enabled":       m.config.MonitorEnabled,
		"check_interval": m.config.HeartbeatCheckInterval.String(),
		"chain_rpc":     m.config.ChainRPCEndpoint,
	}
}

// --- Keeper Bot (auto-heartbeat sender) ---

// KeeperBot automatically sends heartbeats for plans that are about to expire
type KeeperBot struct {
	config   *config.Config
	database *db.Database
	stopCh   chan struct{}
}

// NewKeeperBot creates a new keeper bot
func NewKeeperBot(cfg *config.Config, database *db.Database) *KeeperBot {
	return &KeeperBot{
		config:   cfg,
		database: database,
		stopCh:   make(chan struct{}),
	}
}

// Start begins the keeper bot
func (kb *KeeperBot) Start() {
	log.Println("Starting keeper bot...")
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			kb.checkAndRenew()
		case <-kb.stopCh:
			log.Println("Keeper bot stopped.")
			return
		}
	}
}

// Stop stops the keeper bot
func (kb *KeeperBot) Stop() {
	close(kb.stopCh)
}

// checkAndRenew checks plans nearing expiration and sends reminders
func (kb *KeeperBot) checkAndRenew() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// In production: query chain for plans with heartbeats due soon
	// Send reminder notifications to plan creators
	log.Println("Keeper bot: checking for plans needing heartbeat renewal...")

	// Placeholder: in production, this would query the chain
	_ = ctx
}
