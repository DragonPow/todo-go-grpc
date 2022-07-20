package api_task

import "errors"

func (req *ListReq) Valid() error {
	return nil
}

func (req *GetReq) Valid() error {
	if req.Id == 0 {
		return errors.New("Id must not be empty or zero")
	}
	return nil
}

func (req *CreateReq) Valid() error {
	if req.Name == "" {
		return errors.New("Name of task must not be empty")
	}
	return nil
}

func (req *UpdateReq) Valid() error {
	if req.Id == 0 {
		return errors.New("Id must not be empty or zero")
	}
	if req.NewTaskInfo == nil {
		return errors.New("NewTagInfo must not be empty")
	}
	if req.NewTaskInfo.Name == "" {
		return errors.New("Name of task must not be empty")
	}
	return nil
}

func (req *DeleteMultipleReq) Valid() error {
	if req.TasksId == nil || len(req.TasksId) == 0 {
		return errors.New("Tasks id must not be empty")
	}
	return nil
}
