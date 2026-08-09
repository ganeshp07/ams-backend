package main

import (
	"ams-backend/config"
	"ams-backend/repositories"
	"ams-backend/routes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	config.Load()

	if err := repositories.Connect(config.App.DatabaseURL); err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	if err := runMigrations(); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	// On first deploy (no users yet), run the seeder automatically.
	// This seeds departments, programs, semesters, courses, and the default admin.
	autoSeed()

	r := routes.SetupRouter(repositories.DB)

	addr := ":" + config.App.Port
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func runMigrations() error {
	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory not found")
	}

	matches, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(matches)

	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		for _, stmt := range strings.Split(string(data), ";") {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := repositories.DB.Exec(stmt); err != nil {
				return fmt.Errorf("migration %s: %w", path, err)
			}
		}
		log.Printf("Applied migration: %s", filepath.Base(path))
	}
	return nil
}

// autoSeed runs the seeder binary if no users exist in the database.
// This handles first-run on cloud platforms where we can't run commands manually.
func autoSeed() {
	var count int
	if err := repositories.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		log.Printf("autoSeed: could not check users table: %v", err)
		return
	}

	if count > 0 {
		log.Printf("Database already has %d user(s), skipping auto-seed", count)
		return
	}

	log.Println("No users found — running auto-seed...")

	seederPath := "./seeder"
	if _, err := os.Stat(seederPath); os.IsNotExist(err) {
		log.Println("Seeder binary not found, skipping auto-seed")
		return
	}

	cmd := exec.Command(seederPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ() // pass DATABASE_URL etc. to seeder

	if err := cmd.Run(); err != nil {
		log.Printf("Auto-seed warning: %v (server will still start)", err)
	}
}
