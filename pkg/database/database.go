package database

import (
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/pkg/config"
)

var (
	DB   *sqlx.DB
	once sync.Once
)

func RunMigrations(migrationsDir string) {
	if DB == nil {
		log.Fatal("Database not initialized. Call Connection() first.")
	}

	log.Println("Running database migrations...")
	if err := goose.Up(DB.DB, migrationsDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations ran successfully!")
}

// Connection establishes a PostgreSQL connection using the ENV configuration.
// It sets the global DB variable and ensures it is initialized only once.
func Connection(env *config.ENV) {
	once.Do(func() {
		db, err := sqlx.Connect("postgres", env.DB_URL) // Supabase Postgres URL (with sslmode=require)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		// Optional: ping to ensure DB is reachable
		if err := db.Ping(); err != nil {
			log.Fatalf("Database ping failed: %v", err)
		}
		DB = db
		log.Println("✅ Database connection established successfully")
		RunMigrations("./pkg/migration")
	})
}
