package api_tag

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
	if req.Value == "" {
		return errors.New("Value must not be empty")
	}
	return nil
}

func (req *UpdateReq) Valid() error {
	if req.Id == 0 {
		return errors.New("Id must not be empty or zero")
	}
	if req.NewTagInfo == nil {
		return errors.New("NewTagInfo must not be empty")
	} else if req.NewTagInfo.Value == "" {
		return errors.New("Value of new update tag must not be empty")
	}
	return nil
}

func (req *DeleteReq) Valid() error {
	if req.Id == 0 {
		return errors.New("Id must not be empty or zero")
	}
	return nil
}
