package postgre

import (
	"context"
	"todo-go-grpc/app/task/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgconn"
)

func (s *Suite) Test_taskRepository_GetByID() {
	var (
		query = `SELECT \* FROM "tasks"`
	)

	testcases := []testcase_task{
		{
			// success
			Data: mockDataTask[0],
			TestcaseTemplate: TestcaseTemplate{
				ExportRow: sqlmock.NewRows([]string{"id", "name", "description", "id_done", "created_at", "done_at"}).
					AddRow(mockDataTask[0].ID, mockDataTask[0].Name, mockDataTask[0].Description, mockDataTask[0].IsDone, mockDataTask[0].CreatedAt, mockDataTask[0].DoneAt),
				ExportError: nil,
				Want:        mockDataTask[0],
				WantErr:     nil,
			},
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
		s.GetQuery(Query, query, testcase.TestcaseTemplate, testcase.Data.ID)
		res, err := s.taskRepo.GetByID(context.Background(), testcase.Data.ID)
		if testcase.Want != nil {
			s.Nil(err)
			s.Equal(testcase.Want, *res)
		} else {
			s.ErrorIs(err, testcase.WantErr)

		}
	}
}

func (s *Suite) Test_taskRepository_Create() {
	var (
		query_insert = `INSERT INTO "tasks"`
		query_add    = `INSERT INTO "task_tags"`
	)

	testcases := []testcase_task{
		{
			// success
			Data: mockDataTask[0],
			TestcaseTemplate: TestcaseTemplate{
				ExportRow:   sqlmock.NewRows([]string{"id"}).AddRow(1),
				ExportError: nil,
				Want:        mockDataTask[0],
				WantErr:     nil,
			},
		},
		{
			// fail because duplicate taskname
			Data: mockDataTask[1],
			TestcaseTemplate: TestcaseTemplate{
				ExportError: &pgconn.PgError{Code: "23505"},
				Want:        nil,
				WantErr:     domain.ErrTaskExists,
			},
		},
	}

	for _, testcase := range testcases {
		// s.BeginQuery()
		s.GetQuery(TransactionQuery, query_insert, testcase.TestcaseTemplate,
			testcase.Data.CreatorId,
			testcase.Data.Description,
			testcase.Data.DoneAt,
			testcase.Data.IsDone,
			testcase.Data.Name,
		)
		tags_id := []int32{}
		for _, tag := range testcase.Data.Tags {
			tags_id = append(tags_id, tag.ID)
		}
		s.GetQuery(TransactionQuery, query_add, testcase.TestcaseTemplate, testcase.Data.ID, tags_id)
		// s.EndQuery(testcase.TestcaseTemplate)

		res, err := s.taskRepo.Create(context.Background(), testcase.Data.CreatorId, &testcase.Data)
		if testcase.Want != nil {
			s.Nil(err)
			s.Equal(testcase.Want, *res)
		} else {
			s.ErrorIs(err, testcase.WantErr)
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
			Data: mockDataTask[0],
			Arg: [2][]int32{
				[]int32{mockDataTag[0].ID, mockDataTag[1].ID}, // added tag
				[]int32{mockDataTag[2].ID},                    // remove tag
			},
			TestcaseTemplate: TestcaseTemplate{
				ExportResult: sqlmock.NewResult(0, 1),
				ExportError:  nil,
				Want:         mockDataTask[0],
				WantErr:      nil,
			},
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
		s.GetQuery(TransactionExecute, query, testcase.TestcaseTemplate, testcase.Data.Name, testcase.Data.ID)
		res, err := s.taskRepo.Update(
			context.Background(),
			testcase.Data.ID,
			&testcase.Data,
			testcase.Arg.([]interface{})[0].([]int32),
			testcase.Arg.([]interface{})[1].([]int32),
		)
		if testcase.Want != nil {
			s.Nil(err)
			s.Equal(testcase.Want, *res)
		} else {
			s.ErrorIs(err, testcase.WantErr)

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
			Arg: []int32{mockDataTask[0].ID, mockDataTask[1].ID},
			TestcaseTemplate: TestcaseTemplate{
				ExportResult: sqlmock.NewResult(0, 1),
				ExportError:  nil,
				WantErr:      nil,
			},
		},
		{
			// success with no row
			// data:         mockDataTask[1],
			TestcaseTemplate: TestcaseTemplate{
				ExportResult: sqlmock.NewResult(0, 0),
				ExportError:  nil,
				WantErr:      nil,
			},
		},
	}

	for _, testcase := range testcases {
		s.GetQuery(TransactionExecute, query, testcase.TestcaseTemplate, testcase.Arg.([]int32))
		err := s.taskRepo.Delete(context.Background(), testcase.Arg.([]int32))
		s.ErrorIs(err, testcase.WantErr)
	}
}
