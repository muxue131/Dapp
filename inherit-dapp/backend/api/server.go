package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/inherit-dapp/chain/backend/config"
	"github.com/inherit-dapp/chain/backend/db"
	"github.com/inherit-dapp/chain/crypto"
)

// Server represents the REST API server
type Server struct {
	config    *config.Config
	database  *db.Database
	router    *gin.Engine
	httpServer *http.Server
	ipfsClient *crypto.IPFSClient
}

// NewServer creates a new API server
func NewServer(cfg *config.Config, database *db.Database) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	server := &Server{
		config:     cfg,
		database:   database,
		router:     router,
		ipfsClient: crypto.NewIPFSClient(cfg.IPFSAPIEndpoint),
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Middleware
	s.router.Use(corsMiddleware())
	s.router.Use(requestLogger())

	// Health check
	s.router.GET("/health", s.healthCheck)
	s.router.GET("/api/v1/status", s.getStatus)

	// Plan routes
	planGroup := s.router.Group("/api/v1/plans")
	{
		planGroup.GET("", s.listPlans)
		planGroup.GET("/:id", s.getPlan)
		planGroup.GET("/:id/assets", s.getPlanAssets)
		planGroup.GET("/:id/heartbeat-logs", s.getHeartbeatLogs)
		planGroup.POST("/:id/heartbeat", s.sendHeartbeat)
		planGroup.POST("/:id/claim", s.claimInheritance)
	}

	// Asset routes
	assetGroup := s.router.Group("/api/v1/assets")
	{
		assetGroup.GET("/:id", s.getAsset)
	}

	// Creator routes
	creatorGroup := s.router.Group("/api/v1/creators")
	{
		creatorGroup.GET("/:address/plans", s.getPlansByCreator)
	}

	// Crypto routes
	cryptoGroup := s.router.Group("/api/v1/crypto")
	{
		cryptoGroup.POST("/encrypt", s.encryptData)
		cryptoGroup.POST("/decrypt", s.decryptData)
		cryptoGroup.POST("/shamir/split", s.splitSecret)
		cryptoGroup.POST("/shamir/combine", s.combineSecret)
	}

	// IPFS routes
	ipfsGroup := s.router.Group("/api/v1/ipfs")
	{
		ipfsGroup.POST("/upload", s.uploadToIPFS)
		ipfsGroup.GET("/:cid", s.downloadFromIPFS)
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.ServerHost, s.config.ServerPort)
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

// --- Handlers ---

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"service":   "legacy-dapp-api",
	})
}

func (s *Server) getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"chain_id":   s.config.ChainID,
		"chain_rpc":  s.config.ChainRPCEndpoint,
		"ipfs_api":   s.config.IPFSAPIEndpoint,
		"monitor":    s.config.MonitorEnabled,
		"version":    "0.1.0",
	})
}

func (s *Server) listPlans(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	creator := c.Query("creator")
	var plans []db.PlanRecord
	var err error

	if creator != "" {
		plans, err = s.database.GetPlansByCreator(ctx, creator)
	} else {
		plans = []db.PlanRecord{}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

func (s *Server) getPlan(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	planID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}

	plan, err := s.database.GetPlan(ctx, planID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if plan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plan": plan})
}

func (s *Server) getPlanAssets(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	planID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}

	assets, err := s.database.GetAssetsByPlan(ctx, planID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"assets": assets})
}

func (s *Server) getHeartbeatLogs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	planID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}

	logs, err := s.database.GetHeartbeatLogs(ctx, planID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

func (s *Server) sendHeartbeat(c *gin.Context) {
	planID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}

	var req struct {
		SenderAddress string `json:"sender_address" binding:"required"`
		TxHash        string `json:"tx_hash"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Update heartbeat in database
	plan, err := s.database.GetPlan(ctx, planID)
	if err != nil || plan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}

	newTriggerTime := time.Now().Add(time.Duration(plan.HeartbeatInterval) * time.Second)
	if err := s.database.UpdateHeartbeat(ctx, planID, time.Now(), newTriggerTime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log the heartbeat
	_ = s.database.LogHeartbeat(ctx, planID, req.SenderAddress, req.TxHash)

	c.JSON(http.StatusOK, gin.H{
		"message":      "heartbeat recorded",
		"plan_id":      planID,
		"trigger_time": newTriggerTime,
	})
}

func (s *Server) claimInheritance(c *gin.Context) {
	planID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}

	var req struct {
		BeneficiaryAddress string `json:"beneficiary_address" binding:"required"`
		KeyShares          []int  `json:"key_shares"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	plan, err := s.database.GetPlan(ctx, planID)
	if err != nil || plan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}

	if plan.Status != "triggered" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan is not in triggered status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "claim initiated",
		"plan_id":     planID,
		"beneficiary": req.BeneficiaryAddress,
	})
}

func (s *Server) getAsset(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "asset query endpoint"})
}

func (s *Server) getPlansByCreator(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	address := c.Param("address")
	plans, err := s.database.GetPlansByCreator(ctx, address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

func (s *Server) encryptData(c *gin.Context) {
	var req struct {
		Plaintext string `json:"plaintext" binding:"required"`
		Key       string `json:"key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key, err := crypto.GenerateKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	encrypted, err := crypto.Encrypt([]byte(req.Plaintext), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"encrypted": encrypted, "key": fmt.Sprintf("%x", key)})
}

func (s *Server) decryptData(c *gin.Context) {
	var req struct {
		Ciphertext string `json:"ciphertext" binding:"required"`
		Salt       string `json:"salt"`
		Nonce      string `json:"nonce" binding:"required"`
		Key        string `json:"key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var data crypto.EncryptedData
	data.Ciphertext = req.Ciphertext
	data.Salt = req.Salt
	data.Nonce = req.Nonce

	key, err := fmt.Sscanf(req.Key, "%x", new([]byte))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key format"})
		return
	}
	_ = key

	c.JSON(http.StatusOK, gin.H{"message": "decrypt endpoint - key hex parsing required"})
}

func (s *Server) splitSecret(c *gin.Context) {
	var req struct {
		Secret string `json:"secret" binding:"required"`
		N      int    `json:"n" binding:"required"`
		K      int    `json:"k" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sss := crypto.NewShamirSecretSharing()
	shares, err := sss.Split([]byte(req.Secret), req.N, req.K)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"shares": shares})
}

func (s *Server) combineSecret(c *gin.Context) {
	var req struct {
		Shares []crypto.Share `json:"shares" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sss := crypto.NewShamirSecretSharing()
	secret, err := sss.Combine(req.Shares)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"secret": string(secret)})
}

func (s *Server) uploadToIPFS(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file upload required"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	cid, err := s.ipfsClient.Upload(data, header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cid": cid, "filename": header.Filename})
}

func (s *Server) downloadFromIPFS(c *gin.Context) {
	cid := c.Param("cid")
	data, err := s.ipfsClient.Download(cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", data)
}

// --- Middleware ---

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		log.Printf("[%s] %s %s %d %v",
			time.Now().Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
		)
	}
}
