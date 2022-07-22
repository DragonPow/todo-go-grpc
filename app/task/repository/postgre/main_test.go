package postgre

import (
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"
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
	DB       *dbservice.Database
	mock     sqlmock.Sqlmock
	tagRepo  tagRepository
	taskRepo taskRepository
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
	s.tagRepo = tagRepository{Conn: *s.DB}
	s.taskRepo = taskRepository{Conn: *s.DB}
}

var (
	mockDataTag = []domain.Tag{
		{
			ID:          1,
			Value:       "Value 1",
			Description: "Description 1",
			CreatedAt:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:          2,
			Value:       "Value 2",
			Description: "Description 2",
			CreatedAt:   time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:          3,
			Value:       "Value 3",
			Description: "Description 3",
			CreatedAt:   time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:          4,
			Value:       "Value 4",
			Description: "Description 4",
			CreatedAt:   time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC),
		},
	}
	mockDataTask = []domain.Task{
		{
			ID:          1,
			Name:        "Name 1",
			IsDone:      false,
			Description: "Description 1",
			CreatedAt:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			CreatorId:   1,
		},
		{
			ID:          2,
			Name:        "Name 2",
			IsDone:      false,
			Description: "Description 2",
			CreatedAt:   time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			CreatorId:   2,
		},
		{
			ID:          3,
			Name:        "Name 3",
			IsDone:      true,
			Description: "Description 3",
			CreatedAt:   time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
			DoneAt:      time.Date(2022, 2, 3, 10, 0, 0, 0, time.UTC),
			CreatorId:   1,
		},
		{
			ID:          4,
			Name:        "Name 4",
			IsDone:      true,
			Description: "Description 4",
			CreatedAt:   time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC),
			DoneAt:      time.Date(2022, 2, 1, 30, 0, 0, 0, time.UTC),
			CreatorId:   2,
		},
	}
)

type testcase_tag struct {
	Data domain.Tag `json:"-"`
	TestcaseTemplate
}

type testcase_task struct {
	Data domain.Task `json:"-"`
	Arg  any         `json:"-"`
	TestcaseTemplate
}

type StatusQuery int

const (
	Query StatusQuery = iota
	TransactionQuery
	TransactionExecute
)

type TestcaseTemplate struct {
	ExportRow    *sqlmock.Rows `json:"export_row"`
	ExportResult sql.Result    `json:"export_result"`
	ExportError  error         `json:"export_error"`
	Want         any           `json:"-`
	WantErr      error         `json:"-"`
}

func (s *Suite) BeginQuery() {
	s.mock.ExpectBegin()
}

func (s *Suite) EndQuery(testcase TestcaseTemplate) {
	if testcase.ExportError == nil {
		s.mock.ExpectCommit()
	} else {
		s.mock.ExpectRollback()
	}
}

func (s *Suite) GetQueryWithoutCommit(status StatusQuery, query string, testcase TestcaseTemplate, args ...driver.Value) {
	expecQuery := func() {
		mockQuery := s.mock.ExpectQuery(query).WithArgs(args...)
		if testcase.ExportRow != nil {
			mockQuery.WillReturnRows(testcase.ExportRow)
		} else {
			mockQuery.WillReturnError(testcase.ExportError)
		}
	}
	expectExec := func() {
		mockQuery := s.mock.ExpectExec(query).WithArgs(args...)
		if testcase.ExportResult != nil {
			mockQuery.WillReturnResult(testcase.ExportResult)
		} else {
			mockQuery.WillReturnError(testcase.ExportError)
		}
	}

	if status == Query {
		expecQuery()
	} else {
		if status == TransactionQuery {
			expecQuery()
		} else {
			expectExec()
		}
	}
}

func (s *Suite) GetQuery(status StatusQuery, query string, testcase TestcaseTemplate, args ...driver.Value) {
	expecQuery := func() {
		mockQuery := s.mock.ExpectQuery(query).WithArgs(args...)
		if testcase.ExportRow != nil {
			mockQuery.WillReturnRows(testcase.ExportRow)
		} else {
			mockQuery.WillReturnError(testcase.ExportError)
		}
	}
	expectExec := func() {
		mockQuery := s.mock.ExpectExec(query).WithArgs(args...)
		if testcase.ExportResult != nil {
			mockQuery.WillReturnResult(testcase.ExportResult)
		} else {
			mockQuery.WillReturnError(testcase.ExportError)
		}
	}

	if status == Query {
		expecQuery()
	} else {
		s.mock.ExpectBegin()

		if status == TransactionQuery {
			expecQuery()
		} else {
			expectExec()
		}

		if testcase.ExportError == nil {
			s.mock.ExpectCommit()
		} else {
			s.mock.ExpectRollback()
		}
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
