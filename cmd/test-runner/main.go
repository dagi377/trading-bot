package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	log.Println("Starting End-to-End Test Runner for Hustler Trading Bot")
	
	// Create test results directory
	resultsDir := "./test_results"
	err := os.MkdirAll(resultsDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create results directory: %v", err)
	}
	
	// Run go mod tidy to ensure all dependencies are available
	log.Println("Running go mod tidy...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = filepath.Dir(filepath.Dir(resultsDir)) // Project root
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Warning: go mod tidy failed: %v\n%s", err, output)
	}
	
	// Build the e2e test binary
	log.Println("Building e2e test binary...")
	buildCmd := exec.Command("go", "build", "-o", "e2e-test", "./cmd/e2e-test")
	buildCmd.Dir = filepath.Dir(filepath.Dir(resultsDir)) // Project root
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to build e2e test binary: %v\n%s", err, buildOutput)
	}
	
	// Run the e2e test
	log.Println("Running e2e test...")
	testStartTime := time.Now()
	
	testCmd := exec.Command("./e2e-test")
	testCmd.Dir = filepath.Dir(filepath.Dir(resultsDir)) // Project root
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr
	
	err = testCmd.Start()
	if err != nil {
		log.Fatalf("Failed to start e2e test: %v", err)
	}
	
	// Wait for the test to complete
	err = testCmd.Wait()
	testDuration := time.Since(testStartTime)
	
	if err != nil {
		log.Printf("E2E test failed: %v", err)
		os.Exit(1)
	}
	
	// Check if test results exist
	signalsFile := filepath.Join(resultsDir, "signals.txt")
	messagesFile := filepath.Join(resultsDir, "messages.txt")
	
	if _, err := os.Stat(signalsFile); os.IsNotExist(err) {
		log.Printf("Warning: Signals file not found: %s", signalsFile)
	}
	
	if _, err := os.Stat(messagesFile); os.IsNotExist(err) {
		log.Printf("Warning: Messages file not found: %s", messagesFile)
	}
	
	// Generate test summary
	log.Println("Generating test summary...")
	
	summaryFile := filepath.Join(resultsDir, "summary.txt")
	summary, err := os.Create(summaryFile)
	if err != nil {
		log.Printf("Failed to create summary file: %v", err)
	} else {
		defer summary.Close()
		
		fmt.Fprintf(summary, "Hustler Trading Bot - E2E Test Summary\n")
		fmt.Fprintf(summary, "=====================================\n\n")
		fmt.Fprintf(summary, "Test completed at: %s\n", time.Now().Format(time.RFC1123))
		fmt.Fprintf(summary, "Test duration: %s\n\n", testDuration)
		
		// Count signals
		if signalData, err := os.ReadFile(signalsFile); err == nil {
			signalCount := countOccurrences(string(signalData), "Signal ")
			fmt.Fprintf(summary, "Total signals generated: %d\n", signalCount)
		}
		
		// Count messages
		if messageData, err := os.ReadFile(messagesFile); err == nil {
			messageCount := countOccurrences(string(messageData), "Message ")
			fmt.Fprintf(summary, "Total messages sent: %d\n", messageCount)
		}
		
		fmt.Fprintf(summary, "\nTest Result: SUCCESS\n")
	}
	
	log.Println("E2E test completed successfully!")
	log.Printf("Test duration: %s", testDuration)
	log.Printf("Test results available in: %s", resultsDir)
}

func countOccurrences(s, substr string) int {
	count := 0
	for i := 0; i < len(s); {
		j := indexOf(s[i:], substr)
		if j == -1 {
			break
		}
		count++
		i += j + len(substr)
	}
	return count
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
