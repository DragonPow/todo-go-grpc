package postgre

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"
	"todo-go-grpc/app/dbservice"
	"todo-go-grpc/app/task/domain"

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

var (
	mockData = []domain.Tag{
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
)

type testcase struct {
	data         domain.Tag
	exportRow    *sqlmock.Rows
	exportResult sql.Result
	exportError  error
	want         any
	wantErr      error
}

func (s *Suite) GetQuery(isQuery bool, query string, testcase testcase, args ...driver.Value) {
	s.mock.ExpectBegin()

	if isQuery {
		mockQuery := s.mock.ExpectQuery(query).WithArgs(args...)
		if testcase.exportRow != nil {
			mockQuery.WillReturnRows(testcase.exportRow)
		} else {
			mockQuery.WillReturnError(testcase.exportError)
		}
	} else {
		mockQuery := s.mock.ExpectExec(query).WithArgs(args...)
		if testcase.exportResult != nil {
			mockQuery.WillReturnResult(testcase.exportResult)
		} else {
			mockQuery.WillReturnError(testcase.exportError)
		}
	}

	if testcase.exportError == nil {
		s.mock.ExpectCommit()
	} else {
		s.mock.ExpectRollback()
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_tagRepository_GetByID() {
	case1 := mockData[0]
	case2 := mockData[1]
	var (
		query         = `SELECT \* FROM "tags"`
		expecRowCase1 = sqlmock.NewRows([]string{"id", "value", "description", "created_at"}).AddRow(case1.ID, case1.Value, case1.Description, case1.CreatedAt)
		emptyRows     = sqlmock.NewRows([]string{})
	)
	s.mock.ExpectQuery(query).WithArgs(case1.ID).WillReturnRows(expecRowCase1)
	s.mock.ExpectQuery(query).WithArgs(case2.ID).WillReturnRows(emptyRows)

	// Test found
	res, err := s.repo.GetByID(context.Background(), case1.ID)
	require.NoError(s.T(), err)

	s.Equal(case1, *res)

	// Test not found
	_, err = s.repo.GetByID(context.Background(), case2.ID)
	require.Error(s.T(), err)

	s.Error(err, domain.ErrTagNotExists)
}

func (s *Suite) Test_tagRepository_Create() {
	testcases := []testcase{
		{
			data:         mockData[0],
			exportRow:    sqlmock.NewRows([]string{}),
			exportResult: nil,
			exportError:  nil,
			want:         mockData[0],
			wantErr:      nil,
		},
		{
			data:         mockData[1],
			exportRow:    nil,
			exportResult: nil,
			exportError:  &pgconn.PgError{Code: "23505"},
			want:         nil,
			wantErr:      domain.ErrTagIsExists,
		},
	}

	testcases[1].data.Value = mockData[1].Value

	var (
		query = `INSERT INTO "tags"`
	)

	isQuery := true
	for _, testcase := range testcases {
		s.GetQuery(isQuery, query, testcase, testcase.data.Value, testcase.data.Description, testcase.data.CreatedAt, testcase.data.ID)

		// Do test here
		rs, err := s.repo.Create(context.Background(), &testcase.data)
		if testcase.wantErr != nil {
			s.Nil(rs)
			s.ErrorIs(err, testcase.wantErr)
		} else {
			s.Equal(*rs, testcase.want)
		}
	}
}

func (s *Suite) Test_tagRepository_Update() {
	testcases := []testcase{
		{
			// success
			data:         mockData[0],
			exportRow:    nil,
			exportResult: sqlmock.NewResult(0, 1),
			exportError:  nil,
			want:         mockData[0],
			wantErr:      nil,
		},
		{
			// failse beucase id not exists
			data:         mockData[1],
			exportRow:    nil,
			exportResult: nil,
			exportError:  &pgconn.PgError{Code: "23505"},
			want:         nil,
			wantErr:      domain.ErrTagIsExists,
		},
	}

	// testcases[1].data.Value = mockData[1].Value //

	var (
		// select_query = `SELECT \* FROM "tags"`
		update_query = `UPDATE "tags"`
	)

	isQuery := false
	for _, testcase := range testcases {
		s.GetQuery(isQuery, update_query, testcase, testcase.data.Description, testcase.data.Value, testcase.data.ID)

		// Do test here
		rs, err := s.repo.Update(context.Background(), testcase.data.ID, &testcase.data)
		if testcase.wantErr != nil {
			s.Nil(rs)
			s.ErrorIs(err, testcase.wantErr)
		} else {
			s.Equal(*rs, testcase.want)
		}
	}
}

func (s *Suite) Test_tagRepository_Delete() {
	testcases := []testcase{
		{
			// success
			data:         mockData[0],
			exportResult: sqlmock.NewResult(0, 1),
			exportError:  nil,
			wantErr:      nil,
		},
		{
			// fail because task is link to tag
			data:         mockData[1],
			exportResult: nil,
			exportError:  &pgconn.PgError{Code: "23503"},
			wantErr:      domain.ErrTagStillReference,
		},
	}

	var (
		// select_query = `SELECT \* FROM "tags"`
		update_query = `DELETE FROM "tags"`
	)

	isQuery := false
	for _, testcase := range testcases {
		s.GetQuery(isQuery, update_query, testcase, testcase.data.ID)

		// Do test here
		err := s.repo.Delete(context.Background(), testcase.data.ID)
		s.ErrorIs(err, testcase.wantErr)
	}
}

func (s *Suite) Test_tagRepository_DeleteAll() {
	testcases := []testcase{
		{
			// success
			exportResult: sqlmock.NewResult(0, 1),
			exportError:  nil,
			wantErr:      nil,
		},
		{
			// fail because task is link to tag
			exportResult: nil,
			exportError:  &pgconn.PgError{Code: "23503"},
			wantErr:      domain.ErrTagStillReference,
		},
	}

	var (
		query = `DELETE FROM "tags"`
	)

	isQuery := false
	for _, testcase := range testcases {
		s.GetQuery(isQuery, query, testcase)

		// Do test here
		err := s.repo.DeleteAll(context.Background())
		s.ErrorIs(err, testcase.wantErr)
	}
}

func (s *Suite) Test_tagRepository_FetchAll() {
	testcases := []testcase{
		{
			// success with 3 rows
			exportRow: sqlmock.NewRows([]string{"id", "value", "description", "created_at"}).
				AddRow(mockData[0].ID, mockData[0].Value, mockData[0].Description, mockData[0].CreatedAt).
				AddRow(mockData[1].ID, mockData[1].Value, mockData[1].Description, mockData[1].CreatedAt).
				AddRow(mockData[2].ID, mockData[2].Value, mockData[2].Description, mockData[2].CreatedAt),
			exportError: nil,
			want:        []domain.Tag{mockData[0], mockData[1], mockData[2]},
			wantErr:     nil,
		},
		{
			// success with no rows
			data:        mockData[0],
			exportRow:   sqlmock.NewRows([]string{}),
			exportError: nil,
			want:        []domain.Tag{},
			wantErr:     nil,
		},
	}

	var (
		query = `SELECT \* FROM "tags"`
	)

	for _, testcase := range testcases {
		s.mock.ExpectQuery(query).WillReturnRows(testcase.exportRow)

		// Do test here
		rs, err := s.repo.FetchAll(context.Background())
		require.Nil(s.T(), err)
		s.Equal(rs, testcase.want)
	}
}
