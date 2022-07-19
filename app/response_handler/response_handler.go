package response_handler

import (
	"google.golang.org/grpc/codes"
	grpc_status "google.golang.org/grpc/status"
)

func ResponseErrorNotFound(err error) error {
	return grpc_status.Error(codes.NotFound, err.Error())
}

func ResponseErrorInvalidArgument(err error) error {
	return grpc_status.Error(codes.InvalidArgument, err.Error())
}

func ResponseErrorUnknown(err error) error {
	return grpc_status.Error(codes.Unknown, err.Error())
}

func ResponseErrorPermissionDenied(err error) error {
	return grpc_status.Error(codes.PermissionDenied, err.Error())
}

func ResponseErrorAlreadyExists(err error) error {
	return grpc_status.Error(codes.AlreadyExists, err.Error())
}

func ResponseErrorUnauthenticated(err error) error {
	return grpc_status.Error(codes.Unauthenticated, err.Error())
}
