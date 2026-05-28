package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the depcheck configuration
type Config struct {
	MinimumAgeDays int      `yaml:"minimum_age_days"`
	Exceptions     []string `yaml:"exceptions"`
	CheckMode      string   `yaml:"check_mode"`
}

// ModuleInfo holds proxy response data
type ModuleInfo struct {
	Version string    `json:"Version"`
	Time    time.Time `json:"Time"`
}

// Result holds the check result for a single module
type Result struct {
	Module    string
	Version   string
	Published time.Time
	AgeDays   int
	Passes    bool
	Skipped   bool
	Reason    string
}

const configPath = "configs/depcheck.yaml"
const proxyURL = "https://proxy.golang.org"

func main() {
	cfg := loadConfig()

	fmt.Printf("🔒 Dependency Age Gate\n")
	fmt.Printf("   Minimum age: %d day(s)\n", cfg.MinimumAgeDays)
	fmt.Printf("   Mode: %s\n\n", cfg.CheckMode)

	modules, err := parseGoMod("go.mod")
	if err != nil {
		fatal("Failed to parse go.mod: %v", err)
	}

	results := checkModules(modules, cfg)
	printResults(results, cfg)
}

func loadConfig() Config {
	cfg := Config{
		MinimumAgeDays: 1,
		CheckMode:      "block",
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// Config not found - use defaults
		return cfg
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fatal("Invalid config at %s: %v", configPath, err)
	}
	return cfg
}

func parseGoMod(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	modules := make(map[string]string)
	lines := strings.Split(string(data), "\n")
	inRequire := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "require (" {
			inRequire = true
			continue
		}
		if inRequire && line == ")" {
			inRequire = false
			continue
		}

		// Single-line require
		if strings.HasPrefix(line, "require ") {
			line = strings.TrimPrefix(line, "require ")
			inRequire = false
		}

		if inRequire || strings.HasPrefix(line, "require ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				mod := parts[0]
				ver := strings.TrimSuffix(parts[1], " // indirect")
				ver = strings.Fields(ver)[0]
				modules[mod] = ver
			}
		}
	}
	return modules, nil
}

func checkModules(modules map[string]string, cfg Config) []Result {
	results := make([]Result, 0, len(modules))
	exceptionSet := make(map[string]bool)
	for _, e := range cfg.Exceptions {
		exceptionSet[e] = true
	}

	for mod, ver := range modules {
		r := Result{Module: mod, Version: ver}

		if exceptionSet[mod] {
			r.Skipped = true
			r.Passes = true
			r.Reason = "exception"
			results = append(results, r)
			continue
		}

		info, err := fetchModuleInfo(mod, ver)
		if err != nil {
			r.Reason = fmt.Sprintf("fetch error: %v", err)
			r.Passes = cfg.CheckMode == "warn" // warn mode passes on error
			results = append(results, r)
			continue
		}

		r.Published = info.Time
		r.AgeDays = int(time.Since(info.Time).Hours() / 24)
		r.Passes = r.AgeDays >= cfg.MinimumAgeDays

		if !r.Passes {
			r.Reason = fmt.Sprintf("published %d day(s) ago (minimum: %d)", r.AgeDays, cfg.MinimumAgeDays)
		}

		results = append(results, r)
	}
	return results
}

func fetchModuleInfo(mod, ver string) (*ModuleInfo, error) {
	url := fmt.Sprintf("%s/%s/@v/%s.info", proxyURL, mod, ver)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("proxy returned %d for %s@%s", resp.StatusCode, mod, ver)
	}

	body, _ := io.ReadAll(resp.Body)
	var info ModuleInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("invalid response: %v", err)
	}
	return &info, nil
}

func printResults(results []Result, cfg Config) {
	failed := 0
	warnings := 0

	fmt.Printf("%-60s %-20s %-10s %s\n", "Module", "Version", "Age (days)", "Status")
	fmt.Println(strings.Repeat("─", 110))

	for _, r := range results {
		status := "✅ OK"
		if r.Skipped {
			status = "⏭  skipped"
		} else if !r.Passes {
			if cfg.CheckMode == "warn" {
				status = "⚠️  WARN"
				warnings++
			} else {
				status = "❌ BLOCKED"
				failed++
			}
		}

		age := "-"
		if !r.Published.IsZero() {
			age = fmt.Sprintf("%d", r.AgeDays)
		}

		fmt.Printf("%-60s %-20s %-10s %s\n", truncate(r.Module, 59), truncate(r.Version, 19), age, status)
		if r.Reason != "" && r.Reason != "exception" {
			fmt.Printf("  └─ %s\n", r.Reason)
		}
	}

	fmt.Println(strings.Repeat("─", 110))
	fmt.Printf("\n📊 Summary: %d modules checked", len(results))
	if failed > 0 {
		fmt.Printf(", %d BLOCKED", failed)
	}
	if warnings > 0 {
		fmt.Printf(", %d warnings", warnings)
	}
	fmt.Println()

	if failed > 0 {
		fmt.Printf("\n🚨 %d module(s) failed the minimum age policy. Aborting.\n", failed)
		fmt.Println("   To override, add the module to 'exceptions' in configs/depcheck.yaml")
		os.Exit(1)
	}

	if warnings > 0 {
		fmt.Printf("\n⚠️  %d module(s) are newer than the policy allows (warn mode — not blocking).\n", warnings)
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "depcheck: "+format+"\n", args...)
	os.Exit(2)
}
