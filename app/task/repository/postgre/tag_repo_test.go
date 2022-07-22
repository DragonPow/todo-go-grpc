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
			// success
			Data: mockDataTag[0],
			TestcaseTemplate: TestcaseTemplate{
				ExportRow:    sqlmock.NewRows([]string{}),
				ExportResult: nil,
				ExportError:  nil,
				Want:         mockDataTag[0],
				WantErr:      nil,
			},
		},
		{
			// fail because tag value is exists
			Data: mockDataTag[1],
			TestcaseTemplate: TestcaseTemplate{
				ExportError: &pgconn.PgError{Code: "23505"},
				Want:        nil,
				WantErr:     domain.ErrTagIsExists,
			},
		},
	}

	testcases[1].Data.Value = mockDataTag[1].Value

	var (
		query = `INSERT INTO "tags"`
	)

	for _, testcase := range testcases {
		s.GetQuery(TransactionQuery, query, testcase.TestcaseTemplate, testcase.Data.Value, testcase.Data.Description, testcase.Data.CreatedAt, testcase.Data.ID)

		// Do test here
		rs, err := s.tagRepo.Create(context.Background(), &testcase.Data)
		if testcase.WantErr != nil {
			s.Nil(rs)
			s.ErrorIs(err, testcase.WantErr)
		} else {
			s.Equal(*rs, testcase.Want)
		}
	}
}

func (s *Suite) Test_tagRepository_Update() {
	testcases := []testcase_tag{
		{
			// success
			Data: mockDataTag[0],
			TestcaseTemplate: TestcaseTemplate{
				ExportResult: sqlmock.NewResult(0, 1),
				Want:         mockDataTag[0],
				WantErr:      nil,
			},
		},
		{
			// failse beucase id not exists
			Data: mockDataTag[1],
			TestcaseTemplate: TestcaseTemplate{
				ExportError: &pgconn.PgError{Code: "23505"},
				Want:        nil,
				WantErr:     domain.ErrTagIsExists,
			},
		},
	}

	var (
		update_query = `UPDATE "tags"`
	)

	for _, testcase := range testcases {
		s.GetQuery(TransactionExecute, update_query, testcase.TestcaseTemplate, testcase.Data.Description, testcase.Data.Value, testcase.Data.ID)

		// Do test here
		rs, err := s.tagRepo.Update(context.Background(), testcase.Data.ID, &testcase.Data)
		if testcase.WantErr != nil {
			s.Nil(rs)
			s.ErrorIs(err, testcase.WantErr)
		} else {
			s.Equal(*rs, testcase.Want)
		}
	}
}

func (s *Suite) Test_tagRepository_Delete() {
	testcases := []testcase_tag{
		{
			// success
			Data: mockDataTag[0],
			TestcaseTemplate: TestcaseTemplate{
				ExportResult: sqlmock.NewResult(0, 1),
				ExportError:  nil,
				WantErr:      nil,
			},
		},
		{
			// fail because task is link to tag
			Data: mockDataTag[1],
			TestcaseTemplate: TestcaseTemplate{
				ExportResult: nil,
				ExportError:  &pgconn.PgError{Code: "23503"},
				WantErr:      domain.ErrTagStillReference,
			},
		},
	}

	var (
		// select_query = `SELECT \* FROM "tags"`
		update_query = `DELETE FROM "tags"`
	)

	for _, testcase := range testcases {
		s.GetQuery(TransactionExecute, update_query, testcase.TestcaseTemplate, testcase.Data.ID)

		// Do test here
		err := s.tagRepo.Delete(context.Background(), testcase.Data.ID)
		s.ErrorIs(err, testcase.WantErr)
	}
}

func (s *Suite) Test_tagRepository_DeleteAll() {
	testcases := []testcase_tag{
		{
			// success
			TestcaseTemplate: TestcaseTemplate{
				ExportResult: sqlmock.NewResult(0, 1),
				ExportError:  nil,
				WantErr:      nil,
			},
		},
		{
			// fail because task is link to tag
			TestcaseTemplate: TestcaseTemplate{
				ExportResult: nil,
				ExportError:  &pgconn.PgError{Code: "23503"},
				WantErr:      domain.ErrTagStillReference,
			},
		},
	}

	var (
		query = `DELETE FROM "tags"`
	)

	for _, testcase := range testcases {
		s.GetQuery(TransactionExecute, query, testcase.TestcaseTemplate)

		// Do test here
		err := s.tagRepo.DeleteAll(context.Background())
		s.ErrorIs(err, testcase.WantErr)
	}
}

func (s *Suite) Test_tagRepository_FetchAll() {
	testcases := []testcase_tag{
		{
			// success with 3 rows
			TestcaseTemplate: TestcaseTemplate{
				ExportRow: sqlmock.NewRows([]string{"id", "value", "description", "created_at"}).
					AddRow(mockDataTag[0].ID, mockDataTag[0].Value, mockDataTag[0].Description, mockDataTag[0].CreatedAt).
					AddRow(mockDataTag[1].ID, mockDataTag[1].Value, mockDataTag[1].Description, mockDataTag[1].CreatedAt).
					AddRow(mockDataTag[2].ID, mockDataTag[2].Value, mockDataTag[2].Description, mockDataTag[2].CreatedAt),
				ExportError: nil,
				Want:        []domain.Tag{mockDataTag[0], mockDataTag[1], mockDataTag[2]},
				WantErr:     nil,
			},
		},
		{
			// success with no rows
			Data: mockDataTag[0],
			TestcaseTemplate: TestcaseTemplate{
				ExportRow:   sqlmock.NewRows([]string{}),
				ExportError: nil,
				Want:        []domain.Tag{},
				WantErr:     nil,
			},
		},
	}

	var (
		query = `SELECT \* FROM "tags"`
	)

	for _, testcase := range testcases {
		s.mock.ExpectQuery(query).WillReturnRows(testcase.ExportRow)

		// Do test here
		rs, err := s.tagRepo.FetchAll(context.Background())
		require.Nil(s.T(), err)
		s.Equal(rs, testcase.Want)
	}
}
