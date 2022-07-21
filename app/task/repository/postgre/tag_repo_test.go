package postgre

import (
	"context"
	"todo-go-grpc/app/task/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/require"
)

func (s *Suite) Test_tagRepository_GetByID() {
	case1 := mockDataTag[0]
	case2 := mockDataTag[1]
	var (
		query         = `SELECT \* FROM "tags"`
		expecRowCase1 = sqlmock.NewRows([]string{"id", "value", "description", "created_at"}).AddRow(case1.ID, case1.Value, case1.Description, case1.CreatedAt)
		emptyRows     = sqlmock.NewRows([]string{})
	)
	s.mock.ExpectQuery(query).WithArgs(case1.ID).WillReturnRows(expecRowCase1)
	s.mock.ExpectQuery(query).WithArgs(case2.ID).WillReturnRows(emptyRows)

	// Test found
	res, err := s.tagRepo.GetByID(context.Background(), case1.ID)
	require.NoError(s.T(), err)

	s.Equal(case1, *res)

	// Test not found
	_, err = s.tagRepo.GetByID(context.Background(), case2.ID)
	require.Error(s.T(), err)

	s.Error(err, domain.ErrTagNotExists)
}

func (s *Suite) Test_tagRepository_Create() {
	testcases := []testcase_tag{
		{
			data:         mockDataTag[0],
			exportRow:    sqlmock.NewRows([]string{}),
			exportResult: nil,
			exportError:  nil,
			want:         mockDataTag[0],
			wantErr:      nil,
		},
		{
			data:         mockDataTag[1],
			exportRow:    nil,
			exportResult: nil,
			exportError:  &pgconn.PgError{Code: "23505"},
			want:         nil,
			wantErr:      domain.ErrTagIsExists,
		},
	}

	testcases[1].data.Value = mockDataTag[1].Value

	var (
		query = `INSERT INTO "tags"`
	)

	for _, testcase := range testcases {
		s.GetQuery(TransactionQuery, query, testcase, testcase.data.Value, testcase.data.Description, testcase.data.CreatedAt, testcase.data.ID)

		// Do test here
		rs, err := s.tagRepo.Create(context.Background(), &testcase.data)
		if testcase.wantErr != nil {
			s.Nil(rs)
			s.ErrorIs(err, testcase.wantErr)
		} else {
			s.Equal(*rs, testcase.want)
		}
	}
}

func (s *Suite) Test_tagRepository_Update() {
	testcases := []testcase_tag{
		{
			// success
			data:         mockDataTag[0],
			exportRow:    nil,
			exportResult: sqlmock.NewResult(0, 1),
			exportError:  nil,
			want:         mockDataTag[0],
			wantErr:      nil,
		},
		{
			// failse beucase id not exists
			data:         mockDataTag[1],
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

	for _, testcase := range testcases {
		s.GetQuery(TransactionExecute, update_query, testcase, testcase.data.Description, testcase.data.Value, testcase.data.ID)

		// Do test here
		rs, err := s.tagRepo.Update(context.Background(), testcase.data.ID, &testcase.data)
		if testcase.wantErr != nil {
			s.Nil(rs)
			s.ErrorIs(err, testcase.wantErr)
		} else {
			s.Equal(*rs, testcase.want)
		}
	}
}

func (s *Suite) Test_tagRepository_Delete() {
	testcases := []testcase_tag{
		{
			// success
			data:         mockDataTag[0],
			exportResult: sqlmock.NewResult(0, 1),
			exportError:  nil,
			wantErr:      nil,
		},
		{
			// fail because task is link to tag
			data:         mockDataTag[1],
			exportResult: nil,
			exportError:  &pgconn.PgError{Code: "23503"},
			wantErr:      domain.ErrTagStillReference,
		},
	}

	var (
		// select_query = `SELECT \* FROM "tags"`
		update_query = `DELETE FROM "tags"`
	)

	for _, testcase := range testcases {
		s.GetQuery(TransactionExecute, update_query, testcase, testcase.data.ID)

		// Do test here
		err := s.tagRepo.Delete(context.Background(), testcase.data.ID)
		s.ErrorIs(err, testcase.wantErr)
	}
}

func (s *Suite) Test_tagRepository_DeleteAll() {
	testcases := []testcase_tag{
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

	for _, testcase := range testcases {
		s.GetQuery(TransactionExecute, query, testcase)

		// Do test here
		err := s.tagRepo.DeleteAll(context.Background())
		s.ErrorIs(err, testcase.wantErr)
	}
}

func (s *Suite) Test_tagRepository_FetchAll() {
	testcases := []testcase_tag{
		{
			// success with 3 rows
			exportRow: sqlmock.NewRows([]string{"id", "value", "description", "created_at"}).
				AddRow(mockDataTag[0].ID, mockDataTag[0].Value, mockDataTag[0].Description, mockDataTag[0].CreatedAt).
				AddRow(mockDataTag[1].ID, mockDataTag[1].Value, mockDataTag[1].Description, mockDataTag[1].CreatedAt).
				AddRow(mockDataTag[2].ID, mockDataTag[2].Value, mockDataTag[2].Description, mockDataTag[2].CreatedAt),
			exportError: nil,
			want:        []domain.Tag{mockDataTag[0], mockDataTag[1], mockDataTag[2]},
			wantErr:     nil,
		},
		{
			// success with no rows
			data:        mockDataTag[0],
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
		rs, err := s.tagRepo.FetchAll(context.Background())
		require.Nil(s.T(), err)
		s.Equal(rs, testcase.want)
	}
}
