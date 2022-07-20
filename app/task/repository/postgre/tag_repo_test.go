package postgre

import (
	"database/sql"
	"regexp"
	"testing"
	"todo-go-grpc/app/dbservice"
	"todo-go-grpc/app/task/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB   *dbservice.Database
	mock sqlmock.Sqlmock
	repo tagRepository
}

func (s *Suite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(s.T(), err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(s.T(), err)

	s.DB = &dbservice.Database{Db: gormDB}
	s.repo = tagRepository{Conn: *s.DB}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_tagRepository_GetByID() {
	var (
		id          int32 = 1
		value             = ""
		description       = ""
		query             = `SELECT * FROM "tags" WHERE (id = $1) ORDER BY "tags"."id" LIMIT 1`
		rows              = sqlmock.NewRows([]string{"id", "value", "description"}).AddRow(id, value, description)
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnRows(rows)

	res, err := s.repo.GetByID(nil, id)
	require.NoError(s.T(), err)

	s.NotEqual(res, &domain.Tag{ID: id, Value: value, Description: description})
	s.Equal(res, &domain.Tag{ID: id, Value: value, Description: description})
	s.True(false)
}
