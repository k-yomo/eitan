package graph

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/k-yomo/eitan/src/eitan_service/graph/model"
	"github.com/k-yomo/eitan/src/eitan_service/internal/customerror"
	"github.com/k-yomo/eitan/src/pkg/logging"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// TODO: Don't return raw error message, handle gracefully
func NewErrorPresenter() graphql.ErrorPresenterFunc {
	return func(ctx context.Context, err error) *gqlerror.Error {
		gqlErr := graphql.DefaultErrorPresenter(ctx, err)
		code := mapFromCustomErrorType(customerror.Type(gqlErr.Unwrap()))
		if code == model.ErrorCodeInternal {
			logging.Logger(ctx).Error(err.Error())
		}
		gqlErr.Extensions = map[string]interface{}{"code": code}
		return gqlErr
	}
}

func mapFromCustomErrorType(errType customerror.ErrType) model.ErrorCode {
	switch errType {
	case customerror.ErrUnauthenticated:
		return model.ErrorCodeUnauthenticated
	case customerror.ErrAlreadyExist:
		return model.ErrorCodeAlreadyExist
	case customerror.ErrNotFound:
		return model.ErrorCodeNotFound
	case customerror.ErrInvalidArgument:
		return model.ErrorCodeUnauthenticated
	default:
		return model.ErrorCodeInternal
	}
}