package dbservice

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	domainTag "todo-go-grpc/app/tag/domain"
	domainTask "todo-go-grpc/app/task/domain"
	domainUser "todo-go-grpc/app/user/domain"
)

type Database struct {
	Db *gorm.DB
}

func Init() *Database {
	url := "postgres://postgres:111200@localhost:5432/todo-go-grpc"

	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	// 	logger.Config{
	// 		SlowThreshold:             time.Second, // Slow SQL threshold
	// 		LogLevel:                  logger.Info, // Log level
	// 		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
	// 		Colorful:                  false,       // Disable color
	// 	},
	// )

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		// Logger: newLogger,
	})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&domainTask.Task{}, &domainTag.Tag{}, &domainUser.User{})

	return &Database{Db: db}
}
