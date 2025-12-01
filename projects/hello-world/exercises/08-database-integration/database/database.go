package database

import (
	"fmt"
	"log"

	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/08-database-integration/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// =============================================================================
// DATABASE CONNECTION
// =============================================================================
// In Java/Spring: @Configuration class with @Bean DataSource
// Go: Just a function that returns a configured connection

// DB is the global database instance
// In production, you'd use dependency injection instead
var DB *gorm.DB

// Connect initializes the database connection
// Similar to Spring Boot's auto-configuration but explicit
func Connect() error {
	var err error

	// Open SQLite database (file-based, no server needed)
	// Production: Use PostgreSQL, MySQL, etc.
	// gorm.io/driver/postgres, gorm.io/driver/mysql
	DB, err = gorm.Open(sqlite.Open("app.db"), &gorm.Config{
		// Logger configuration - like hibernate.show_sql=true
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("✅ Database connected successfully!")
	return nil
}

// AutoMigrate creates/updates database schema
// Similar to Hibernate's ddl-auto: update
// In production, use proper migrations (golang-migrate, goose, etc.)
func AutoMigrate() error {
	log.Println("🔄 Running database migrations...")

	// AutoMigrate creates tables, missing foreign keys, constraints, columns, indexes
	// It WON'T delete unused columns (safe for production)
	err := DB.AutoMigrate(
		// One-to-One
		&model.User{},
		&model.Profile{},

		// One-to-Many
		&model.Post{},
		&model.Author{},
		&model.Book{},

		// Many-to-Many
		&model.Tag{},
		&model.Student{},
		&model.Course{},
		&model.Enrollment{},

		// Self-referential
		&model.Category{},

		// Polymorphic
		&model.Comment{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("✅ Database migration completed!")
	return nil
}

// GetDB returns the database instance
// For dependency injection patterns
func GetDB() *gorm.DB {
	return DB
}

// =============================================================================
// TRANSACTION HELPER
// =============================================================================
// In Java/Spring: @Transactional annotation
// Go: Explicit transaction management

// Transaction executes a function within a database transaction
// If the function returns an error, the transaction is rolled back
// Usage:
//
//	err := database.Transaction(func(tx *gorm.DB) error {
//	    // All operations here are in one transaction
//	    if err := tx.Create(&user).Error; err != nil {
//	        return err  // Rollback
//	    }
//	    return nil  // Commit
//	})
func Transaction(fn func(tx *gorm.DB) error) error {
	return DB.Transaction(fn)
}

// =============================================================================
// CONNECTION POOL SETTINGS (Production)
// =============================================================================
// In Java/Spring: spring.datasource.hikari.* properties
// Go: Configure through sql.DB

// ConfigurePool sets up connection pooling
// Call this after Connect() in production
func ConfigurePool() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns: idle connections in pool
	// Java equivalent: HikariCP minimumIdle
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns: max open connections
	// Java equivalent: HikariCP maximumPoolSize
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime: max time a connection can be reused
	// Java equivalent: HikariCP maxLifetime
	// sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}

