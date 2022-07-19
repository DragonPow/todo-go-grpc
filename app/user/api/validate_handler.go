package api

import "errors"

func (req *LoginReq) Valid() error {
	if req.Username == "" || req.Password == "" {
		return errors.New("Username or Password must not be empty")
	}
	return nil
}

func (req *GetReq) Valid() error {
	if req.Id == 0 {
		return errors.New("Id must not be empty or zero")
	}
	return nil
}

func (req *CreateReq) Valid() error {
	if req.Name == "" || req.Username == "" || req.Password == "" {
		return errors.New("Name, Username or Password must not be empty")
	}
	return nil
}

func (req *UpdateReq) Valid() error {
	if req.Id == 0 {
		return errors.New("Id must not be empty or zero")
	}
	if req.NewUserInfor == nil {
		return errors.New("NewUserInfor must not be empty")
	}
	return nil
}

func (req *DeleteReq) Valid() error {
	if req.Id == 0 {
		return errors.New("Id must not be empty or zero")
	}
	return nil
}
