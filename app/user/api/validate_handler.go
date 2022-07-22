package api

import "errors"

var (
	ErrUsernameOrPaswordIsEmpty = errors.New("Username or Password must not be empty")
	ErrIdIsEmpty                = errors.New("Id must not be empty or zero")
	ErrCreateReqIsEmpty         = errors.New("Name, Username or Password must not be empty")
	ErrUpdateReqUpdateIsEmpty   = errors.New("NewUserInfor must not be empty")
)

func (req *LoginReq) Valid() error {
	if req.Username == "" || req.Password == "" {
		return ErrUsernameOrPaswordIsEmpty
	}
	return nil
}

func (req *GetReq) Valid() error {
	if req.Id == 0 {
		return ErrIdIsEmpty
	}
	return nil
}

func (req *CreateReq) Valid() error {
	if req.Name == "" || req.Username == "" || req.Password == "" {
		return ErrCreateReqIsEmpty
	}
	return nil
}

func (req *UpdateReq) Valid() error {
	if req.Id == 0 {
		return ErrIdIsEmpty
	}
	if req.NewUserInfor == nil {
		return ErrUpdateReqUpdateIsEmpty
	}
	return nil
}

func (req *DeleteReq) Valid() error {
	if req.Id == 0 {
		return ErrIdIsEmpty
	}
	return nil
}
