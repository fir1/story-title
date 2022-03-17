package utilerror

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ConversionToGrpcError(err error) error {
	if _, ok := status.FromError(err); ok {
		return err
	}

	var (
		errArgument     ErrArgument
		errForbidden    ErrForbidden
		errPrecondition ErrPrecondition
		errUnauthorised ErrUnauthorised
	)

	switch {
	case errors.As(err, &errArgument):
		return status.Error(codes.InvalidArgument, fmt.Sprintf("%v (%v)", err, errArgument.Unwrap()))
	case errors.As(err, &errForbidden):
		return status.Error(codes.PermissionDenied, fmt.Sprintf("%v (%v)", err, errForbidden.Unwrap()))
	case errors.As(err, &errPrecondition):
		return status.Error(codes.FailedPrecondition, fmt.Sprintf("%v (%v)", err, errPrecondition.Unwrap()))
	case errors.As(err, &errUnauthorised):
		return status.Error(codes.Unauthenticated, fmt.Sprintf("%v (%v)", err, errUnauthorised.Unwrap()))
	}

	return status.Error(codes.Internal, "internal error")
}

func ConversionToGrpcDetailedError(err error) error {
	var (
		errArgument     ErrArgument
		errForbidden    ErrForbidden
		errPrecondition ErrPrecondition
		errUnauthorized ErrUnauthorised
	)

	switch {
	case errors.As(err, &errArgument):
		return status.Error(codes.InvalidArgument, FormatErrorForLog(err))
	case errors.As(err, &errForbidden):
		return status.Error(codes.PermissionDenied, FormatErrorForLog(err))
	case errors.As(err, &errPrecondition):
		return status.Error(codes.FailedPrecondition, FormatErrorForLog(err))
	case errors.As(err, &errUnauthorized):
		return status.Error(codes.Unauthenticated, FormatErrorForLog(err))
	}

	return status.Error(codes.Internal, fmt.Sprintf("internal error: %v", err))
}

func FormatErrorForLog(err error) string {
	var (
		errArgument     ErrArgument
		errForbidden    ErrForbidden
		errPrecondition ErrPrecondition
		errUnauthorized ErrUnauthorised
	)

	switch {
	case errors.As(err, &errArgument):
		return fmt.Sprintf("BLL error: %v (%v)", err, errArgument.Unwrap())
	case errors.As(err, &errForbidden):
		return fmt.Sprintf("BLL error: %v (%v)", err, errForbidden.Unwrap())
	case errors.As(err, &errPrecondition):
		return fmt.Sprintf("BLL error: %v (%v)", err, errPrecondition.Unwrap())
	case errors.As(err, &errUnauthorized):
		return fmt.Sprintf("BLL error: %v (%v)", err, errUnauthorized.Unwrap())
	}

	return err.Error()
}
