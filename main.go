package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/schollz/progressbar/v3"
)

// Config represents application configuration
type Config struct {
	ApiID      string `json:"api_id"`
	ApiHash    string `json:"api_hash"`
	Phone      string `json:"phone"`
	BaseDelay  int    `json:"base_delay_seconds"`
	MaxRetries int    `json:"max_retries"`
	Workers    int    `json:"workers"`
	UseProxy   bool   `json:"use_proxy"`
	ProxyURL   string `json:"proxy_url,omitempty"`
}

// CheckResult represents the result of a username check
type CheckResult struct {
	Username  string    `json:"username"`
	Status    string    `json:"status"`
	CheckTime time.Time `json:"check_time"`
	Error     string    `json:"error,omitempty"`
}

// State represents the current state of the application
type State struct {
	Checked         map[string]CheckResult `json:"checked"`
	LastCheckTime   time.Time              `json:"last_check_time"`
	RemainingTasks  []string               `json:"remaining_tasks"`
	FloodWaitUntil  time.Time              `json:"flood_wait_until"`
	CurrentUsername string                 `json:"current_username"`
}

var (
	configFile     = flag.String("config", "config.json", "Path to config file")
	inputFile      = flag.String("input", "usernames.txt", "Path to input file")
	outputDir      = flag.String("output", "results", "Directory for output files")
	stateFile      = flag.String("state", "state.json", "Path to state file")
	verbose        = flag.Bool("verbose", false, "Enable verbose logging")
	generateCombos = flag.Bool("generate", false, "Generate username combinations")
	minLength      = flag.Int("min-length", 3, "Minimum username length")
	maxLength      = flag.Int("max-length", 30, "Maximum username length")
	interactive    = flag.Bool("interactive", false, "Interactive mode")
)

func displayHeader() {
	header := `
╔════════════════════════════════════════════╗
║        Telegram Username Checker           ║
║         Created by: xPOURY4                ║
║     GitHub: github.com/xPOURY4             ║
║       Extended Version (v2.0)              ║
╚════════════════════════════════════════════╝
`
	fmt.Println(header)
}

func main() {
	flag.Parse()
	displayHeader()

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Load configuration
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Printf("Warning: Could not load config file: %v", err)
		log.Println("Using default configuration. Please create a config file for future use.")
		config = &Config{
			ApiID:      "YOUR_API_ID",
			ApiHash:    "YOUR_API_HASH",
			Phone:      "YOUR_PHONE_NUMBER",
			BaseDelay:  3,
			MaxRetries: 3,
			Workers:    1,
		}
	}

	// Load or initialize state
	state, err := loadState(*stateFile)
	if err != nil {
		log.Printf("Starting from fresh state: %v", err)
		state = &State{
			Checked:        make(map[string]CheckResult),
			LastCheckTime:  time.Now(),
			RemainingTasks: []string{},
		}
	}

	// Handle interactive mode
	if *interactive {
		runInteractiveMode(config, state)
		return
	}

	// Load usernames
	usernames, err := loadUsernames(*inputFile, state)
	if err != nil {
		log.Fatalf("Error loading usernames: %v", err)
	}

	// If we should generate combinations
	if *generateCombos {
		usernames = generateCombinations(usernames, *minLength, *maxLength)
		log.Printf("Generated %d username combinations", len(usernames))
	}

	// Initialize Telegram client
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	setupSignalHandler(cancel, state)

	// Create Telegram client
	client := createTelegramClient(config)

	// Run the client
	if err := client.Run(ctx, func(ctx context.Context) error {
		// Authenticate
		if err := authenticate(ctx, client, config.Phone); err != nil {
			return err
		}

		log.Println("Successfully logged in!")

		// Configure worker pool
		workerCount := config.Workers
		if workerCount < 1 {
			workerCount = 1
		}

		// Process usernames
		results := processUsernames(ctx, client, usernames, state, config, workerCount)

		// Save results
		for status, statusResults := range results {
			if len(statusResults) > 0 {
				filename := filepath.Join(*outputDir, fmt.Sprintf("%s_usernames.txt", status))
				if err := saveResults(filename, statusResults); err != nil {
					log.Printf("Error saving %s usernames: %v", status, err)
				} else {
					log.Printf("%s usernames have been saved to %s", strings.Title(status), filename)
				}
			}
		}

		// Save detailed JSON results
		detailedResults := make([]CheckResult, 0, len(state.Checked))
		for _, result := range state.Checked {
			detailedResults = append(detailedResults, result)
		}
		
		detailedFilename := filepath.Join(*outputDir, "detailed_results.json")
		if err := saveJSONResults(detailedFilename, detailedResults); err != nil {
			log.Printf("Error saving detailed results: %v", err)
		} else {
			log.Printf("Detailed results saved to %s", detailedFilename)
		}

		return nil
	}); err != nil {
		log.Printf("Error: %v", err)
		// Save state before exiting
		saveState(*stateFile, state)
		os.Exit(1)
	}

	// Final state save
	saveState(*stateFile, state)
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func loadState(filename string) (*State, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

func saveState(filename string, state *State) error {
	state.LastCheckTime = time.Now()
	
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func loadUsernames(filename string, state *State) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(string(content), "\n")
	var usernames []string
	
	// Filter out already checked usernames unless they had errors
	for _, line := range lines {
		username := strings.TrimSpace(line)
		if username == "" {
			continue
		}
		
		if result, checked := state.Checked[username]; checked {
			if result.Status == "error" {
				// Retry errors
				usernames = append(usernames, username)
			}
		} else {
			usernames = append(usernames, username)
		}
	}
	
	// Add any remaining tasks from previous run
	usernames = append(state.RemainingTasks, usernames...)
	
	return usernames, nil
}

func generateCombinations(seeds []string, minLen, maxLen int) []string {
	// This is a simplified version. You can expand this to generate more complex combinations.
	var results []string
	
	// Add the seeds
	for _, seed := range seeds {
		if len(seed) >= minLen && len(seed) <= maxLen {
			results = append(results, seed)
		}
	}
	
	// Generate simple variations
	suffixes := []string{"", "_", ".", "0", "1", "2", "3", "official", "real", "thereal"}
	
	for _, seed := range seeds {
		for _, suffix := range suffixes {
			combined := seed + suffix
			if len(combined) >= minLen && len(combined) <= maxLen && !contains(results, combined) {
				results = append(results, combined)
			}
		}
	}
	
	return results
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func createTelegramClient(config *Config) *telegram.Client {
	// Set up client options
	options := telegram.Options{}
	
	// Add more customization as needed
	
	return telegram.NewClient(config.ApiID, config.ApiHash, options)
}

func authenticate(ctx context.Context, client *telegram.Client, phone string) error {
	flow := auth.NewFlow(
		auth.Constant(phone, "", auth.CodeAuthenticatorFunc(func(ctx context.Context, _ *tg.AuthSentCode) (string, error) {
			fmt.Print("Enter the code you received: ")
			var code string
			_, err := fmt.Scan(&code)
			return code, err
		})),
		auth.SendCodeOptions{},
	)

	if err := flow.Run(ctx, client.Auth()); err != nil {
		return fmt.Errorf("auth flow: %w", err)
	}
	
	return nil
}

func processUsernames(ctx context.Context, client *telegram.Client, usernames []string, state *State, config *Config, workerCount int) map[string][]CheckResult {
	results := make(map[string][]CheckResult)
	resultsMu := sync.Mutex{}
	
	// Check if we need to wait due to previous flood wait
	if time.Now().Before(state.FloodWaitUntil) {
		waitDuration := time.Until(state.FloodWaitUntil)
		log.Printf("Waiting %.2f minutes due to previous rate limit...", waitDuration.Minutes())
		time.Sleep(waitDuration)
	}
	
	// Create task channel
	tasks := make(chan string, len(usernames))
	for _, username := range usernames {
		tasks <- username
	}
	close(tasks)
	
	// Set up progress bar
	bar := progressbar.Default(int64(len(usernames)))
	
	// Create wait group for workers
	var wg sync.WaitGroup
	
	// Rate limiter shared across workers
	rateLimiter := make(chan struct{}, 1)
	go func() {
		baseDelay := time.Duration(config.BaseDelay) * time.Second
		ticker := time.NewTicker(baseDelay)
		defer ticker.Stop()
		
		for range ticker.C {
			select {
			case <-ctx.Done():
				return
			case rateLimiter <- struct{}{}:
				// Token added to limiter
			}
		}
	}()
	
	// Launch workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			for username := range tasks {
				select {
				case <-ctx.Done():
					return
				case <-rateLimiter:
					// Got permission to proceed
				}
				
				state.CurrentUsername = username
				state.RemainingTasks = removeUsername(state.RemainingTasks, username)
				
				// Try with retries
				var status string
				var checkErr error
				var retryCount int
				
				for retryCount = 0; retryCount < config.MaxRetries; retryCount++ {
					status, checkErr = checkUsername(ctx, client.API(), username)
					
					if checkErr == nil {
						break
					}
					
					if strings.Contains(checkErr.Error(), "FLOOD_WAIT") {
						// Parse wait time and set global wait
						waitTimeSec := parseFloodWaitSeconds(checkErr.Error())
						waitTime := time.Duration(waitTimeSec) * time.Second
						
						// Set flood wait until time
						state.FloodWaitUntil = time.Now().Add(waitTime)
						
						log.Printf("Rate limit hit. Waiting for %.2f minutes", waitTime.Minutes())
						
						// Save state in case of program termination during wait
						saveState(*stateFile, state)
						
						time.Sleep(waitTime)
						
						// Increase the base delay for future requests
						rateLimiter = make(chan struct{}, 1)
						baseDelay := time.Duration(config.BaseDelay) * time.Second * time.Duration(math.Pow(2, float64(retryCount+1)))
						ticker := time.NewTicker(baseDelay)
						go func() {
							defer ticker.Stop()
							for range ticker.C {
								select {
								case <-ctx.Done():
									return
								case rateLimiter <- struct{}{}:
									// Token added
								}
							}
						}()
					} else {
						// For other errors, just wait a bit and retry
						time.Sleep(time.Second * time.Duration(retryCount+1))
					}
				}
				
				// Create result
				result := CheckResult{
					Username:  username,
					Status:    status,
					CheckTime: time.Now(),
				}
				
				if checkErr != nil {
					result.Status = "error"
					result.Error = checkErr.Error()
					if *verbose {
						log.Printf("Error checking %s (retry %d/%d): %v", username, retryCount, config.MaxRetries, checkErr)
					}
				} else {
					if *verbose {
						log.Printf("%s is %s", username, status)
					}
				}
				
				// Update results
				resultsMu.Lock()
				results[result.Status] = append(results[result.Status], result)
				state.Checked[username] = result
				resultsMu.Unlock()
				
				// Update progress bar
				bar.Add(1)
				
				// Periodically save state (every 10 usernames)
				if len(state.Checked)%10 == 0 {
					saveState(*stateFile, state)
				}
			}
		}()
	}
	
	// Wait for all workers to finish
	wg.Wait()
	
	return results
}

func removeUsername(slice []string, username string) []string {
	for i, u := range slice {
		if u == username {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func parseFloodWaitSeconds(errMsg string) int {
	// Parse error message like "FLOOD_WAIT_X" where X is seconds
	parts := strings.Split(errMsg, "_")
	if len(parts) >= 3 {
		var seconds int
		fmt.Sscanf(parts[2], "%d", &seconds)
		if seconds > 0 {
			return seconds
		}
	}
	
	// Default to 60 seconds if we couldn't parse
	return 60
}

func checkUsername(ctx context.Context, api *tg.Client, username string) (string, error) {
	_, err := api.AccountUpdateUsername(ctx, username)
	if err != nil {
		if strings.Contains(err.Error(), "USERNAME_OCCUPIED") {
			return "taken", nil
		}
		if strings.Contains(err.Error(), "USERNAME_INVALID") {
			return "invalid", nil
		}
		if strings.Contains(err.Error(), "USERNAME_PURCHASE_AVAILABLE") {
			return "purchasable", nil
		}
		return "error", err
	}
	
	_, err = api.AccountUpdateUsername(ctx, "") // Reset username
	if err != nil {
		return "error", err
	}
	
	return "available", nil
}

func saveResults(filename string, results []CheckResult) error {
	var lines []string
	for _, result := range results {
		lines = append(lines, result.Username)
	}
	
	content := strings.Join(lines, "\n")
	return os.WriteFile(filename, []byte(content), 0644)
}

func saveJSONResults(filename string, results []CheckResult) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(filename, data, 0644)
}

func setupSignalHandler(cancel context.CancelFunc, state *State) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-c
		log.Println("\nReceived interrupt signal. Saving state and exiting gracefully...")
		saveState(*stateFile, state)
		cancel()
		os.Exit(0)
	}()
}

func runInteractiveMode(config *Config, state *State) {
	ctx := context.Background()
	client := createTelegramClient(config)
	
	if err := client.Run(ctx, func(ctx context.Context) error {
		if err := authenticate(ctx, client, config.Phone); err != nil {
			return err
		}
		
		log.Println("Successfully logged in!")
		log.Println("Interactive mode: Type usernames to check (one per line). Type 'exit' to quit.")
		
		scanner := os.NewScanner(os.Stdin)
		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				break
			}
			
			username := strings.TrimSpace(scanner.Text())
			if username == "exit" {
				break
			}
			
			if username == "" {
				continue
			}
			
			status, err := checkUsername(ctx, client.API(), username)
			if err != nil {
				log.Printf("Error checking %s: %v", username, err)
				continue
			}
			
			log.Printf("%s is %s", username, status)
			
			// Save to state
			state.Checked[username] = CheckResult{
				Username:  username,
				Status:    status,
				CheckTime: time.Now(),
			}
		}
		
		return nil
	}); err != nil {
		log.Printf("Error: %v", err)
	}
}
