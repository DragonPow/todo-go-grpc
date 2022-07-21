package postgre

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"
	"todo-go-grpc/app/dbservice"
	"todo-go-grpc/app/user/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB   *dbservice.Database
	mock sqlmock.Sqlmock
	repo userRepository
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
	s.repo = userRepository{Conn: *s.DB}
}

var (
	mockData = []domain.User{
		{
			ID:        1,
			Name:      "Vu Ngoc Thach",
			Username:  "vungocthach",
			Password:  "123456",
			CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:        2,
			Name:      "Huynh Thi Minh Nhuc",
			Username:  "huynhthiminhnhuc",
			Password:  "1234",
			CreatedAt: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:        3,
			Name:      "Vo Thi Thuy Tien",
			Username:  "thuytien",
			Password:  "12345678",
			CreatedAt: time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:        4,
			Name:      "Nguyen Pham Duy Bang",
			Username:  "duybang",
			Password:  "99999999",
			CreatedAt: time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC),
		},
	}
)

type testcase struct {
	data         domain.User
	exportRow    *sqlmock.Rows
	exportResult sql.Result
	exportError  error
	want         any
	wantErr      error
}

type StatusQuery int

const (
	Query StatusQuery = iota
	TransactionQuery
	TransactionExecute
)

func (s *Suite) GetQuery(status StatusQuery, query string, testcase testcase, args ...driver.Value) {
	expecQuery := func() {
		mockQuery := s.mock.ExpectQuery(query).WithArgs(args...)
		if testcase.exportRow != nil {
			mockQuery.WillReturnRows(testcase.exportRow)
		} else {
			mockQuery.WillReturnError(testcase.exportError)
		}
	}
	expectExec := func() {
		mockQuery := s.mock.ExpectExec(query).WithArgs(args...)
		if testcase.exportResult != nil {
			mockQuery.WillReturnResult(testcase.exportResult)
		} else {
			mockQuery.WillReturnError(testcase.exportError)
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

		if testcase.exportError == nil {
			s.mock.ExpectCommit()
		} else {
			s.mock.ExpectRollback()
		}
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_userRepository_GetByID() {
	var (
		query = `SELECT \* FROM "users"`
	)

	testcases := []testcase{
		{
			// success
			data: mockData[0],
			exportRow: sqlmock.NewRows([]string{"id", "name", "username", "password", "created_at"}).
				AddRow(mockData[0].ID, mockData[0].Name, mockData[0].Username, mockData[0].Password, mockData[0].CreatedAt),
			exportError: nil,
			want:        mockData[0],
			wantErr:     nil,
		},
		{
			// fail because id not exists
			data:        mockData[0],
			exportError: gorm.ErrRecordNotFound,
			want:        nil,
			wantErr:     domain.ErrUserNotExists,
		},
	}

	for _, testcase := range testcases {
		s.GetQuery(Query, query, testcase, testcase.data.ID)
		res, err := s.repo.GetByID(context.Background(), testcase.data.ID)
		if testcase.want != nil {
			s.Nil(err)
			s.Equal(testcase.want, *res)
		} else {
			s.ErrorIs(err, testcase.wantErr)

		}
	}
}

func (s *Suite) Test_userRepository_GetByUsernameAndPassword() {
	var (
		query = `SELECT \* FROM "users"`
	)

	testcases := []testcase{
		{
			// success
			data: mockData[0],
			exportRow: sqlmock.NewRows([]string{"id", "name", "username", "password", "created_at"}).
				AddRow(mockData[0].ID, mockData[0].Name, mockData[0].Username, mockData[0].Password, mockData[0].CreatedAt),
			exportError: nil,
			want:        mockData[0],
			wantErr:     nil,
		},
		{
			// fail because username or passwrod wrong
			data:        mockData[0],
			exportError: gorm.ErrRecordNotFound,
			want:        nil,
			wantErr:     domain.ErrUserNotExists,
		},
	}

	for _, testcase := range testcases {
		s.GetQuery(Query, query, testcase, testcase.data.Username, testcase.data.Password)
		res, err := s.repo.GetByUsernameAndPassword(context.Background(), testcase.data.Username, testcase.data.Password)
		if testcase.want != nil {
			s.Nil(err)
			s.Equal(testcase.want, *res)
		} else {
			s.ErrorIs(err, testcase.wantErr)

		}
	}
}

func (s *Suite) Test_userRepository_Create() {
	var (
		query = `INSERT INTO "users"`
	)

	testcases := []testcase{
		{
			// success
			data:        mockData[0],
			exportRow:   sqlmock.NewRows([]string{}),
			exportError: nil,
			want:        mockData[0],
			wantErr:     nil,
		},
		{
			// fail because duplicate username
			data:        mockData[1],
			exportError: &pgconn.PgError{Code: "23505"},
			want:        nil,
			wantErr:     domain.ErrUserNameIsExists,
		},
	}

	for _, testcase := range testcases {
		s.GetQuery(TransactionQuery, query, testcase, testcase.data.Name, testcase.data.Password, testcase.data.Username)
		res, err := s.repo.Create(context.Background(), &testcase.data)
		if testcase.want != nil {
			s.Nil(err)
			s.Equal(testcase.want, *res)
		} else {
			s.ErrorIs(err, testcase.wantErr)

		}
	}
}

func (s *Suite) Test_userRepository_Update() {
	var (
		query = `UPDATE "users"`
	)

	testcases := []testcase{
		{
			// success
			data:         mockData[0],
			exportResult: sqlmock.NewResult(0, 1),
			exportError:  nil,
			want:         mockData[0],
			wantErr:      nil,
		},
		{
			// fail because duplicate username
			data:        mockData[1],
			exportError: &pgconn.PgError{Code: "23505"},
			want:        nil,
			wantErr:     domain.ErrUserNameIsExists,
		},
	}

	for _, testcase := range testcases {
		s.GetQuery(TransactionExecute, query, testcase, testcase.data.Name, testcase.data.Password, testcase.data.Username, testcase.data.ID)
		res, err := s.repo.Update(context.Background(), testcase.data.ID, &testcase.data)
		if testcase.want != nil {
			s.Nil(err)
			s.Equal(testcase.want, *res)
		} else {
			s.ErrorIs(err, testcase.wantErr)

		}
	}
}

func (s *Suite) Test_userRepository_Delete() {
	var (
		query = `DELETE FROM "users"`
	)

	testcases := []testcase{
		{
			// success
			data:         mockData[0],
			exportResult: sqlmock.NewResult(0, 1),
			exportError:  nil,
			wantErr:      nil,
		},
		{
			// success with no row
			data:         mockData[1],
			exportResult: sqlmock.NewResult(0, 0),
			exportError:  nil,
			wantErr:      nil,
		},
	}

	for _, testcase := range testcases {
		s.GetQuery(TransactionExecute, query, testcase, testcase.data.ID)
		err := s.repo.Delete(context.Background(), testcase.data.ID)
		s.ErrorIs(err, testcase.wantErr)
	}
}
