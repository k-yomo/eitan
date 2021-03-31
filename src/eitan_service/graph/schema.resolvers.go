package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	gql "github.com/k-yomo/eitan/src/eitan_service/graph/generated"
	"github.com/k-yomo/eitan/src/eitan_service/graph/model"
	"github.com/k-yomo/eitan/src/eitan_service/internal/customerror"
	"github.com/k-yomo/eitan/src/internal/pb/eitan"
)

func (r *queryResolver) Node(ctx context.Context, id string) (model.Node, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Nodes(ctx context.Context, ids []string) ([]model.Node, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) CurrentAccount(ctx context.Context) (*model.Account, error) {
	res, err := r.accountServiceClient.GetCurrentAccount(ctx, &eitan.Empty{})
	if err != nil {
		return nil, customerror.New(err, customerror.ErrInternal)
	}

	return &model.Account{
		ID:           res.Account.Id,
		Email:        res.Account.Email,
		DisplayName:  res.Account.DisplayName,
		ScreenImgURL: res.Account.ScreenImgUrl,
	}, nil
}

// Query returns gql.QueryResolver implementation.
func (r *Resolver) Query() gql.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
