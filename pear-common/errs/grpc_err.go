package errs

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcError (e *BError) error {
	return status.Error(codes.Code(e.Code), e.Msg)
}