package database

import (
	"mononoke-go/model"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	dialSqlite3  = "sqlite3"
	dialMysql    = "mysql"
	dialPostgres = "postgres"
	dialSQLSrv   = "sqlserver"
)

// New creates a new wrapper for the gorm database framework.
func New(dial, conn, defaultUser, defaultPass string, createDefault bool) (*GormDatabase, error) {
	createDirectoryIfSqlite(dial, conn)
	var db *gorm.DB
	var err error
	switch dial {
	case dialSqlite3:
		db, err = gorm.Open(sqlite.Open(conn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: "mogo_",
			},
		})
		if err != nil {
			return nil, err
		}
	case dialMysql:
		db, err = gorm.Open(mysql.Open(conn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: "mogo_",
			},
		})
		if err != nil {
			return nil, err
		}
	case dialPostgres:
		db, err = gorm.Open(postgres.Open(conn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: "mogo_",
			},
		})
		if err != nil {
			return nil, err
		}
	case dialSQLSrv:
		db, err = gorm.Open(sqlserver.Open(conn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: "mogo_",
			},
		})
		if err != nil {
			return nil, err
		}
	}

	// We normally don't need that much connections, so we limit them. F.ex. mysql complains about
	// "too many connections", while load testing Gotify.
	database, err := db.DB()
	if err != nil {
		return nil, err
	}
	database.SetMaxOpenConns(10) // Set max open conns to 10

	if dial == dialSqlite3 {
		// We use the database connection inside the handlers from the http
		// framework, therefore concurrent access occurs. Sqlite cannot handle
		// concurrent writes, so we limit sqlite to one connection.
		// see https://github.com/mattn/go-sqlite3/issues/274
		database.SetMaxOpenConns(1)
	}

	if dial == dialMysql {
		// Mysql has a setting called wait_timeout, which defines the duration
		// after which a connection may not be used anymore.
		// The default for this setting on mariadb is 10 minutes.
		// See https://github.com/docker-library/mariadb/issues/113
		database.SetConnMaxLifetime(9 * time.Minute)
	}

	if err = db.AutoMigrate(new(model.Accounts)); err != nil {
		return nil, err
	}

	userCount := int64(0)
	db.Find(new(model.Accounts)).Count(&userCount)
	if createDefault && userCount == 0 {
		db.Create(&model.Accounts{AccountName: defaultUser, Password: defaultPass})
	}

	return &GormDatabase{DB: db}, nil
}

func createDirectoryIfSqlite(dialect, connection string) {
	if dialect == dialSqlite3 {
		if _, err := os.Stat(filepath.Dir(connection)); os.IsNotExist(err) {
			if err = os.MkdirAll(filepath.Dir(connection), 0o777); err != nil {
				panic(err)
			}
		}
	}
}

// GormDatabase is a wrapper for the gorm framework.
type GormDatabase struct {
	DB *gorm.DB
}

func (d *GormDatabase) Close() {
	database, err := d.DB.DB()
	if err != nil {
		panic(err)
	}
	database.Close()
}
