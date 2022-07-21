package dbservice

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	// taskDomain "todo-go-grpc/app/task/domain"
	// userDomain "todo-go-grpc/app/user/domain"
)

type Database struct {
	Db *gorm.DB
}

func Init() *Database {
	url := "postgres://postgres:111200@localhost:5432/todo-go-grpc"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatalln(err)
	}

	// db.AutoMigrate(&taskDomain.Task{}, &taskDomain.Tag{}, &userDomain.User{})

	return &Database{Db: db}
}
