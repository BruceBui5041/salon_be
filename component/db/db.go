package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"salon_be/component"
	"salon_be/component/logger"
	models "salon_be/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type dbInstances struct {
	mySQL *gorm.DB
}

func ConnectToDB(ctx context.Context) component.DBInstances {
	// Get database connection details from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Construct the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormlogger.Config{
			SlowThreshold:             time.Second,     // Slow SQL threshold
			LogLevel:                  gormlogger.Info, // Log level
			IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,           // Don't include params in the SQL log
			Colorful:                  true,            // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		// uncomment on the first migration run to create tables without relation error
		// DisableForeignKeyConstraintWhenMigrating: true,
		IgnoreRelationshipsWhenMigrating: true,
	})

	if err != nil {
		logger.AppLogger.Fatal(ctx, err.Error())
	}

	return &dbInstances{mySQL: db}
}

func (d *dbInstances) GetMySQLDBConnection() *gorm.DB {
	return d.mySQL
}

func (d *dbInstances) AutoMigrateMySQL() error {
	migrator := d.mySQL.Migrator()

	tables := []interface{}{
		&models.User{},
		&models.Role{},
		&models.Category{},

		&models.UserAuth{},
		&models.UserProfile{},
		&models.UserRole{},
		&models.Permission{},
		&models.RolePermission{},
		&models.Video{},
		&models.VideoProcessInfo{},
		&models.ServiceVersion{},
		&models.Service{},

		&models.Payment{},
		&models.Enrollment{},

		&models.Comment{},
		&models.Image{},
	}

	for _, table := range tables {
		tableName := d.mySQL.NamingStrategy.TableName(fmt.Sprintf("%T", table))

		if !migrator.HasTable(table) {
			log.Printf("Creating table: %s", tableName)
			if err := migrator.CreateTable(table); err != nil {
				return fmt.Errorf("failed to create table %s: %w", tableName, err)
			}
			log.Printf("Table created successfully: %s", tableName)
		} else {
			log.Printf("Table already exists: %s", tableName)
			// Optionally, you can still run AutoMigrate to add new columns or indexes
			if err := migrator.AutoMigrate(table); err != nil {
				return fmt.Errorf("failed to auto migrate table %s: %w", tableName, err)
			}
			log.Printf("Table auto migrated successfully: %s", tableName)
		}
	}

	return nil
}
