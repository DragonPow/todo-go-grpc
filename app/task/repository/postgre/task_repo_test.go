package postgre

import (
	"context"

	"github.com/DATA-DOG/go-sqlmock"
)

func (s *Suite) Test_taskRepository_GetByID() {
	var (
		query = `SELECT \* FROM "tasks"`
	)

	testcases := []testcase_task{
		{
			// success
			data: mockDataTask[0],
			exportRow: sqlmock.NewRows([]string{"id", "name", "description", "id_done", "created_at", "done_at"}).
				AddRow(mockDataTask[0].ID, mockDataTask[0].Name, mockDataTask[0].Description, mockDataTask[0].IsDone, mockDataTask[0].CreatedAt, mockDataTask[0].DoneAt),
			exportError: nil,
			want:        mockDataTask[0],
			wantErr:     nil,
		},
		// {
		// 	// fail because id not exists
		// 	data:        mockDataTask[0],
		// 	exportError: gorm.ErrRecordNotFound,
		// 	want:        nil,
		// 	wantErr:     domain.ErrtaskNotExists,
		// },
	}

	for _, testcase := range testcases {
		s.GetQuery(Query, query, testcase, testcase.data.ID)
		res, err := s.taskRepo.GetByID(context.Background(), testcase.data.ID)
		if testcase.want != nil {
			s.Nil(err)
			s.Equal(testcase.want, *res)
		} else {
			s.ErrorIs(err, testcase.wantErr)

		}
	}
}

func (s *Suite) Test_taskRepository_Create() {
	var (
		query            = `INSERT INTO "tasks"`
		creator_id int32 = 1
	)

	testcases := []testcase_task{
		{
			// success
			data:        mockDataTask[0],
			arg:         creator_id,
			exportRow:   sqlmock.NewRows([]string{}),
			exportError: nil,
			want:        mockDataTask[0],
			wantErr:     nil,
		},
		// {
		// 	// fail because duplicate taskname
		// 	data:        mockDataTask[1],
		// 	exportError: &pgconn.PgError{Code: "23505"},
		// 	want:        nil,
		// 	wantErr:     domain.ErrtaskNameIsExists,
		// },
	}

	for _, testcase := range testcases {
		s.GetQuery(TransactionQuery, query, testcase)
		res, err := s.taskRepo.Create(context.Background(), testcase.arg.(int32), &testcase.data)
		if testcase.want != nil {
			s.Nil(err)
			s.Equal(testcase.want, *res)
		} else {
			s.ErrorIs(err, testcase.wantErr)

		}
	}
}

func (s *Suite) Test_taskRepository_Update() {
	var (
		query = `UPDATE "tasks"`
	)

	testcases := []testcase_task{
		{
			// success
			data: mockDataTask[0],
			arg: [2][]int32{
				[]int32{mockDataTag[0].ID, mockDataTag[1].ID}, // added tag
				[]int32{mockDataTag[2].ID},                    // remove tag
			},
			exportResult: sqlmock.NewResult(0, 1),
			exportError:  nil,
			want:         mockDataTask[0],
			wantErr:      nil,
		},
		// {
		// 	// fail because duplicate taskname
		// 	data:        mockDataTask[1],
		// 	exportError: &pgconn.PgError{Code: "23505"},
		// 	want:        nil,
		// 	wantErr:     domain.ErrtaskNameIsExists,
		// },
	}

	for _, testcase := range testcases {
		s.GetQuery(TransactionExecute, query, testcase, testcase.data.Name, testcase.data.ID)
		res, err := s.taskRepo.Update(
			context.Background(),
			testcase.data.ID,
			&testcase.data,
			testcase.arg.([]interface{})[0].([]int32),
			testcase.arg.([]interface{})[1].([]int32),
		)
		if testcase.want != nil {
			s.Nil(err)
			s.Equal(testcase.want, *res)
		} else {
			s.ErrorIs(err, testcase.wantErr)

		}
	}
}

func (s *Suite) Test_taskRepository_Delete() {
	var (
		query = `DELETE FROM "tasks"`
	)

	testcases := []testcase_task{
		{
			// success
			// data:         mockDataTask[0],
			arg:          []int32{mockDataTask[0].ID, mockDataTask[1].ID},
			exportResult: sqlmock.NewResult(0, 1),
			exportError:  nil,
			wantErr:      nil,
		},
		{
			// success with no row
			// data:         mockDataTask[1],
			exportResult: sqlmock.NewResult(0, 0),
			exportError:  nil,
			wantErr:      nil,
		},
	}

	for _, testcase := range testcases {
		s.GetQuery(TransactionExecute, query, testcase, testcase.arg.([]int32))
		err := s.taskRepo.Delete(context.Background(), testcase.arg.([]int32))
		s.ErrorIs(err, testcase.wantErr)
	}
}
