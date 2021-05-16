package customerror

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrType int

const (
	ErrUnauthenticated ErrType = iota + 1
	ErrNotFound
	ErrAlreadyExist
	ErrInvalidArgument
	ErrExpired
	ErrUnknown
	ErrInternal
	// charge error
	ErrCardDeclined
	ErrBalanceInsufficient
	ErrIncorrectCVC
	ErrExpiredCard
	ErrChargeProcessingError
)

type customError struct {
	ErrType
	err error
}

func New(err error, errType ErrType) error {
	return &customError{ErrType: errType, err: err}
}

func FromGrpcError(grpcErr *status.Status) error {
	var errType ErrType
	// use unknown error for commented out error
	switch grpcErr.Code() {
	// case codes.Canceled:
	case codes.Unknown:
		errType = ErrUnknown
	case codes.InvalidArgument:
		errType = ErrInvalidArgument
	// case codes.DeadlineExceeded:
	case codes.NotFound:
		errType = ErrNotFound
	case codes.AlreadyExists:
		errType = ErrAlreadyExist
	// case codes.ResourceExhausted:
	case codes.FailedPrecondition:
		errType = ErrInvalidArgument
	// case codes.Aborted:
	case codes.OutOfRange:
		errType = ErrInvalidArgument
	// case codes.Unimplemented:
	case codes.Internal:
		errType = ErrInternal
	// case codes.Unavailable:
	// case codes.DataLoss:
	case codes.Unauthenticated:
		errType = ErrUnauthenticated
	default:
		errType = ErrUnknown
	}
	return &customError{ErrType: errType, err: grpcErr.Err()}
}

func (ce *customError) Error() string {
	return ce.err.Error()
}

func NewErrUnauthenticated() error {
	return New(errors.New("unauthenticated"), ErrUnauthenticated)
}

func NewErrNotFound(err error) error {
	return New(err, ErrNotFound)
}

func NewErrInvalidArgument(err error) error {
	return New(err, ErrInvalidArgument)
}

func NewErrGetUserIDFailedInAuthRequiredMethod() error {
	return New(errors.New("failed to get user id"), ErrInvalidArgument)
}

func Type(err error) ErrType {
	e, ok := err.(*customError)
	if !ok {
		return ErrUnknown
	}
	return e.ErrType
}