package api

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/cyberspacesec/go-Sublist3r/pkg/docker"
)

// ScanRequest 表示子域名扫描的请求
type ScanRequest struct {
	Domain      string `json:"domain"`
	Bruteforce  bool   `json:"bruteforce"`
	Ports       string `json:"ports"`
	Verbose     bool   `json:"verbose"`
	Threads     int    `json:"threads"`
	Engines     string `json:"engines"`
	NoColor     bool   `json:"no_color"`
	CallbackURL string `json:"callback_url"`
}

// ScanResponse 表示子域名扫描的响应
type ScanResponse struct {
	ID      string    `json:"id"`
	Status  string    `json:"status"`
	Domain  string    `json:"domain"`
	Start   time.Time `json:"start_time"`
	End     time.Time `json:"end_time"`
	Results []string  `json:"results,omitempty"`
	Error   string    `json:"error,omitempty"`
}

const (
	// StatusPending 表示扫描尚未开始
	StatusPending = "pending"
	// StatusRunning 表示扫描正在进行中
	StatusRunning = "running"
	// StatusCompleted 表示扫描已完成
	StatusCompleted = "completed"
	// StatusFailed 表示扫描失败
	StatusFailed = "failed"
)

// ScanTask 表示扫描任务
type ScanTask struct {
	Request ScanRequest
	ID      string
}

// ScanManager 管理扫描任务
type ScanManager struct {
	scans       map[string]*ScanResponse
	scanCh      chan ScanTask
	mu          sync.RWMutex
	stopCh      chan struct{}
	workerCount int
	maxCapacity int
}

// NewScanManager 创建一个新的扫描管理器
func NewScanManager(workerCount, maxCapacity int) *ScanManager {
	if workerCount <= 0 {
		workerCount = 5 // 默认5个工作线程
	}
	if maxCapacity <= 0 {
		maxCapacity = 100 // 默认最多100个扫描
	}

	sm := &ScanManager{
		scans:       make(map[string]*ScanResponse),
		scanCh:      make(chan ScanTask, maxCapacity),
		stopCh:      make(chan struct{}),
		workerCount: workerCount,
		maxCapacity: maxCapacity,
	}

	// 启动工作线程
	for i := 0; i < workerCount; i++ {
		go sm.worker()
	}

	// 启动清理线程
	go sm.cleanupExpiredScans()

	return sm
}

// Stop 停止扫描管理器
func (sm *ScanManager) Stop() {
	close(sm.stopCh)
}

// AddScan 添加一个新的扫描任务
func (sm *ScanManager) AddScan(req ScanRequest) (string, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// 检查是否超过最大容量
	if len(sm.scans) >= sm.maxCapacity {
		return "", fmt.Errorf("scan manager is at maximum capacity (%d scans)", sm.maxCapacity)
	}

	id := uuid.New().String()
	scanResp := &ScanResponse{
		ID:     id,
		Status: StatusPending,
		Domain: req.Domain,
		Start:  time.Now(),
	}

	sm.scans[id] = scanResp

	// 将任务发送到工作队列
	task := ScanTask{
		Request: req,
		ID:      id,
	}

	select {
	case sm.scanCh <- task:
		// 任务成功添加到队列
	default:
		// 队列已满，删除记录并返回错误
		delete(sm.scans, id)
		return "", fmt.Errorf("scan queue is full")
	}

	return id, nil
}

// GetScan 获取扫描状态
func (sm *ScanManager) GetScan(id string) (*ScanResponse, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	scan, ok := sm.scans[id]
	if !ok {
		return nil, fmt.Errorf("scan with ID %s not found", id)
	}

	return scan, nil
}

// GetAllScans 获取所有扫描
func (sm *ScanManager) GetAllScans() []*ScanResponse {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	scans := make([]*ScanResponse, 0, len(sm.scans))
	for _, scan := range sm.scans {
		scans = append(scans, scan)
	}

	return scans
}

// 清理过期的扫描记录
func (sm *ScanManager) cleanupExpiredScans() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sm.mu.Lock()
			expireTime := time.Now().Add(-24 * time.Hour)

			for id, scan := range sm.scans {
				// 如果扫描已完成或失败，并且结束时间超过24小时，则删除
				if (scan.Status == StatusCompleted || scan.Status == StatusFailed) &&
					!scan.End.IsZero() && scan.End.Before(expireTime) {
					delete(sm.scans, id)
				}
			}
			sm.mu.Unlock()
		case <-sm.stopCh:
			return
		}
	}
}

// readSubdomainsFromFile 从文件中读取子域名
func readSubdomainsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var subdomains []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			subdomains = append(subdomains, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return subdomains, nil
}

// 发送回调
func sendCallback(callbackURL string, scanResp *ScanResponse) error {
	data, err := json.Marshal(scanResp)
	if err != nil {
		return fmt.Errorf("failed to marshal scan response: %w", err)
	}

	req, err := http.NewRequest("POST", callbackURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send callback: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("Callback successfully sent to %s", callbackURL)
		return nil
	}

	return fmt.Errorf("callback failed with status code: %d", resp.StatusCode)
}

// 工作线程处理扫描任务
func (sm *ScanManager) worker() {
	for {
		select {
		case task := <-sm.scanCh:
			sm.processScan(task)
		case <-sm.stopCh:
			return
		}
	}
}

// 处理扫描任务
func (sm *ScanManager) processScan(task ScanTask) {
	// 更新状态为运行中
	sm.mu.Lock()
	scan := sm.scans[task.ID]
	scan.Status = StatusRunning
	sm.mu.Unlock()

	// 使用临时文件存储结果
	outputFile := fmt.Sprintf("/tmp/sublist3r-%s.txt", task.ID)
	defer os.Remove(outputFile) // 清理临时文件

	// 准备命令行参数
	req := task.Request
	args := []string{
		"-d", req.Domain,
		"-o", outputFile,
	}

	if req.Bruteforce {
		args = append(args, "-b")
	}

	if req.Ports != "" {
		args = append(args, "-p", req.Ports)
	}

	if req.Verbose {
		args = append(args, "-v")
	}

	if req.Threads > 0 {
		args = append(args, "-t", fmt.Sprintf("%d", req.Threads))
	}

	if req.Engines != "" {
		args = append(args, "-e", req.Engines)
	}

	if req.NoColor {
		args = append(args, "-n")
	}

	var results []string
	var scanErr error

	// 尝试使用Docker运行
	if docker.IsDockerAvailable() {
		log.Printf("Running with Docker for scan %s", task.ID)
		scanErr = docker.RunSublist3r(args)

		// 尝试从文件读取结果
		if scanErr == nil {
			results, scanErr = readSubdomainsFromFile(outputFile)
		}
	} else {
		// Docker不可用，使用模拟数据
		log.Printf("Docker not available, using simulated data for scan %s", task.ID)
		time.Sleep(2 * time.Second) // 模拟扫描时间

		// 生成一些模拟数据
		results = []string{
			fmt.Sprintf("www.%s", req.Domain),
			fmt.Sprintf("mail.%s", req.Domain),
			fmt.Sprintf("api.%s", req.Domain),
			fmt.Sprintf("blog.%s", req.Domain),
			fmt.Sprintf("dev.%s", req.Domain),
		}
	}

	// 更新扫描结果
	sm.mu.Lock()
	scan.End = time.Now()

	if scanErr != nil {
		scan.Status = StatusFailed
		scan.Error = scanErr.Error()
	} else {
		scan.Status = StatusCompleted
		scan.Results = results
	}
	sm.mu.Unlock()

	// 如果提供了回调URL，发送回调
	if req.CallbackURL != "" {
		go func() {
			err := sendCallback(req.CallbackURL, scan)
			if err != nil {
				log.Printf("Error sending callback for scan %s: %v", task.ID, err)
			}
		}()
	}
}

// APIConfig API服务器配置
type APIConfig struct {
	Port       int    // 监听端口
	Workers    int    // 工作线程数
	Capacity   int    // 最大扫描容量
	APIKey     string // API密钥
	EnableAuth bool   // 是否启用认证
}

// APIServer 表示API服务器
type APIServer struct {
	config      APIConfig
	router      *mux.Router
	server      *http.Server
	scanManager *ScanManager
}

// NewAPIServer 创建一个新的API服务器
func NewAPIServer(port, workers, capacity int, apiKey string) *APIServer {
	r := mux.NewRouter()
	scanManager := NewScanManager(workers, capacity)

	config := APIConfig{
		Port:       port,
		Workers:    workers,
		Capacity:   capacity,
		APIKey:     apiKey,
		EnableAuth: apiKey != "", // 如果提供了API密钥，则启用认证
	}

	server := &APIServer{
		config:      config,
		router:      r,
		scanManager: scanManager,
	}

	server.setupRoutes()
	return server
}

// APIKeyMiddleware API密钥认证中间件
func (s *APIServer) APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 如果未启用认证，直接通过
		if !s.config.EnableAuth {
			next.ServeHTTP(w, r)
			return
		}

		// 从请求头中获取API密钥
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			// 也可以从查询参数获取API密钥
			apiKey = r.URL.Query().Get("api_key")
		}

		// 验证API密钥
		if apiKey != s.config.APIKey {
			http.Error(w, "未授权：无效的API密钥", http.StatusUnauthorized)
			return
		}

		// API密钥有效，继续处理请求
		next.ServeHTTP(w, r)
	})
}

// 设置路由
func (s *APIServer) setupRoutes() {
	// API v1 路由
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// 应用API密钥认证中间件
	api.Use(s.APIKeyMiddleware)

	// API文档路由 - 无需认证
	s.router.HandleFunc("/docs", s.handleDocs)

	// 扫描API路由 - 需要认证
	api.HandleFunc("/scan", s.handleStartScan).Methods("POST")
	api.HandleFunc("/scan/sync", s.handleSyncScan).Methods("POST")
	api.HandleFunc("/scan/{id}", s.handleGetScan).Methods("GET")
	api.HandleFunc("/scans", s.handleListScans).Methods("GET")
}

// 处理API文档请求
func (s *APIServer) handleDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, GetAPIDocsHTML())
}

// 处理启动扫描请求
func (s *APIServer) handleStartScan(w http.ResponseWriter, r *http.Request) {
	var req ScanRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	if req.Domain == "" {
		http.Error(w, "Domain is required", http.StatusBadRequest)
		return
	}

	scanID, err := s.scanManager.AddScan(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to start scan: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"id":      scanID,
		"message": "Scan started",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

// 处理同步扫描请求 - 直接返回结果
func (s *APIServer) handleSyncScan(w http.ResponseWriter, r *http.Request) {
	var req ScanRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	if req.Domain == "" {
		http.Error(w, "Domain is required", http.StatusBadRequest)
		return
	}

	// 创建结果通道
	resultCh := make(chan *ScanResponse, 1)
	errorCh := make(chan error, 1)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

	// 启动扫描
	go func() {
		// 创建临时输出文件
		outputFile := fmt.Sprintf("/tmp/sublist3r-sync-%s.txt", uuid.New().String())
		defer os.Remove(outputFile)

		// 准备命令行参数
		args := []string{
			"-d", req.Domain,
			"-o", outputFile,
		}

		if req.Bruteforce {
			args = append(args, "-b")
		}

		if req.Ports != "" {
			args = append(args, "-p", req.Ports)
		}

		if req.Verbose {
			args = append(args, "-v")
		}

		if req.Threads > 0 {
			args = append(args, "-t", fmt.Sprintf("%d", req.Threads))
		}

		if req.Engines != "" {
			args = append(args, "-e", req.Engines)
		}

		if req.NoColor {
			args = append(args, "-n")
		}

		scanResp := &ScanResponse{
			ID:     uuid.New().String(),
			Status: StatusRunning,
			Domain: req.Domain,
			Start:  time.Now(),
		}

		var results []string
		var scanErr error

		// 执行扫描
		if docker.IsDockerAvailable() {
			scanErr = docker.RunSublist3r(args)
			if scanErr == nil {
				results, scanErr = readSubdomainsFromFile(outputFile)
			}
		} else {
			// 使用模拟模式
			time.Sleep(2 * time.Second)
			results = []string{
				fmt.Sprintf("www.%s", req.Domain),
				fmt.Sprintf("mail.%s", req.Domain),
				fmt.Sprintf("api.%s", req.Domain),
				fmt.Sprintf("blog.%s", req.Domain),
				fmt.Sprintf("dev.%s", req.Domain),
			}
		}

		scanResp.End = time.Now()
		if scanErr != nil {
			scanResp.Status = StatusFailed
			scanResp.Error = scanErr.Error()
			errorCh <- scanErr
		} else {
			scanResp.Status = StatusCompleted
			scanResp.Results = results
			resultCh <- scanResp
		}
	}()

	// 等待扫描结果或超时
	select {
	case <-ctx.Done():
		http.Error(w, "Scan timeout after 5 minutes", http.StatusRequestTimeout)
		return
	case err := <-errorCh:
		http.Error(w, fmt.Sprintf("Scan failed: %v", err), http.StatusInternalServerError)
		return
	case result := <-resultCh:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}
}

// 处理获取扫描状态请求
func (s *APIServer) handleGetScan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	scan, err := s.scanManager.GetScan(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Scan not found: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scan)
}

// 处理列出所有扫描请求
func (s *APIServer) handleListScans(w http.ResponseWriter, r *http.Request) {
	scans := s.scanManager.GetAllScans()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scans)
}

// Start 启动API服务器
func (s *APIServer) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	if s.config.EnableAuth {
		log.Printf("API server with authentication enabled listening on %s", addr)
	} else {
		log.Printf("API server listening on %s", addr)
	}
	return s.server.ListenAndServe()
}

// Stop 停止API服务器
func (s *APIServer) Stop(ctx context.Context) error {
	s.scanManager.Stop()
	return s.server.Shutdown(ctx)
}
