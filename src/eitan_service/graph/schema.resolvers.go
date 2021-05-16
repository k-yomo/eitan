package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
	gql "github.com/k-yomo/eitan/src/eitan_service/graph/generated"
	"github.com/k-yomo/eitan/src/eitan_service/graph/model"
	"github.com/k-yomo/eitan/src/eitan_service/infra"
	"github.com/k-yomo/eitan/src/eitan_service/internal/auth"
	"github.com/k-yomo/eitan/src/eitan_service/internal/customerror"
	"github.com/k-yomo/eitan/src/eitan_service/internal/rediskeys"
	"github.com/k-yomo/eitan/src/internal/pb/eitan"
	"github.com/k-yomo/eitan/src/pkg/clock"
	"github.com/k-yomo/eitan/src/pkg/logging"
	"github.com/k-yomo/eitan/src/pkg/uuid"
	"go.uber.org/zap"
)

func (r *mutationResolver) UpdatePlayerID(ctx context.Context, input model.UpdatePlayerIDInput) (*model.Player, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CancelWaitingMatch(ctx context.Context) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateCurrentPlayerQuizRoomStatus(ctx context.Context, roomID string) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateAnswer(ctx context.Context, roomID string, quizID string, answer string) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *playerResolver) UserProfile(ctx context.Context, obj *model.Player) (*model.UserProfile, error) {
	res, err := r.accountServiceClient.GetUserProfile(ctx, &eitan.GetUserProfileRequest{
		UserId: obj.UserID,
	})
	if err != nil {
		return nil, customerror.New(err, customerror.ErrInternal)
	}
	return mapToGraphqlUserProfile(res.UserProfile), nil
}

func (r *queryResolver) Node(ctx context.Context, id string) (model.Node, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Nodes(ctx context.Context, ids []string) ([]model.Node, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) CurrentUserProfile(ctx context.Context) (*model.CurrentUserProfile, error) {
	res, err := r.accountServiceClient.GetCurrentUserProfile(ctx, &eitan.Empty{})
	if err != nil {
		return nil, customerror.New(err, customerror.ErrInternal)
	}

	return mapToGraphqlCurrentUserProfile(res.UserProfile), nil
}

func (r *queryResolver) CurrentPlayer(ctx context.Context) (*model.Player, error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, customerror.NewErrGetUserIDFailedInAuthRequiredMethod()
	}

	player, err := infra.PlayerByUserID(ctx, r.db, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerror.NewErrNotFound(err)
		}
		return nil, customerror.New(err, customerror.ErrInternal)
	}
	return mapToGraphqlPlayer(player), nil
}

func (r *subscriptionResolver) RandomMatchRoomDecided(ctx context.Context) (<-chan *model.QuizRoom, error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, customerror.NewErrGetUserIDFailedInAuthRequiredMethod()
	}

	player, err := infra.PlayerByUserID(ctx, r.db, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerror.NewErrNotFound(err)
		}
		return nil, customerror.New(err, customerror.ErrInternal)
	}

	now := clock.Now()

	roomChan := make(chan *model.QuizRoom, 1)
	var roomDecided bool
	var matchedPlayer *infra.Player
	err = r.txManager.RunInTx(ctx, func(ctx context.Context) error {
		// TODO: This getting waiting users logic is very simple, and should be fixed
		//   1. Longest waiting players should match first
		//   2. We shouldn't take all waiting players (it can be enormous)
		//   3. We shouldn't get too old waiting players (more than 30m ago)
		waitingPlayers, err := infra.GetAllMatchWaitingPlayers(ctx, r.db)
		if err != nil {
			return err
		}
		if len(waitingPlayers) == 0 {
			return nil
		}

		matchedWaitingPlayer := waitingPlayers[0]
		if err := matchedWaitingPlayer.Delete(ctx, r.db); err != nil {
			return err
		}

		matchedPlayer, err = infra.GetPlayer(ctx, r.db, matchedWaitingPlayer.PlayerID)
		if err != nil {
			return err
		}

		quizRoom := infra.QuizRoom{
			ID:        uuid.Generate(),
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := quizRoom.Insert(ctx, r.db); err != nil {
			return err
		}

		quizRoomUsers := []*infra.QuizRoomPlayer{
			{QuizRoomID: quizRoom.ID, PlayerID: player.ID, CreatedAt: now},
			{QuizRoomID: quizRoom.ID, PlayerID: matchedPlayer.ID, CreatedAt: now},
		}
		for _, qru := range quizRoomUsers {
			if err := qru.Insert(ctx, r.db); err != nil {
				return err
			}
		}

		room := &model.QuizRoom{
			ID: quizRoom.ID,
			Players: []*model.Player{
				mapToGraphqlPlayer(player),
				mapToGraphqlPlayer(matchedPlayer),
			},
		}
		roomJSON, err := json.Marshal(room)
		if err != nil {
			return  err
		}
		err = r.redisClient.Publish(ctx, rediskeys.NewWaitingRoomPlayerKey(matchedPlayer.ID).String(), roomJSON).Err()
		if err != nil {
			return  err
		}
		roomChan <- room
		roomDecided = true

		return nil
	})
	if err != nil {
		return nil, customerror.New(err, customerror.ErrInternal)
	}

	if roomDecided {
		return nil, err
	}

	// if not match, let user wait matching
	waitingPlayer := &infra.MatchWaitingPlayer{
		PlayerID:  player.ID,
		CreatedAt: clock.Now(),
	}
	if err := waitingPlayer.Insert(ctx, r.db); err != nil {
		return nil, err
	}

	go func() {
		pubsub := r.redisClient.Subscribe(ctx, rediskeys.NewWaitingRoomPlayerKey(player.ID).String())
		defer pubsub.Close()
		for {
			time.Sleep(1 * time.Second)
			res, _ := pubsub.Receive(ctx)
			switch msg := res.(type) {
			case *redis.Message:
				room := model.QuizRoom{}
				if err := json.Unmarshal([]byte(msg.Payload), &room); err != nil {
					logging.Logger(ctx).Error("json.Unmarshal", zap.Error(err))
					break
				}
				roomChan <- &room
				return
			default:
			}
		}
	}()
	return roomChan, nil
}

func (r *subscriptionResolver) QuizPosted(ctx context.Context, roomID string) (<-chan model.Quiz, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *subscriptionResolver) QuizAnswered(ctx context.Context, roomID string) (<-chan model.QuizAnswer, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns gql.MutationResolver implementation.
func (r *Resolver) Mutation() gql.MutationResolver { return &mutationResolver{r} }

// Player returns gql.PlayerResolver implementation.
func (r *Resolver) Player() gql.PlayerResolver { return &playerResolver{r} }

// Query returns gql.QueryResolver implementation.
func (r *Resolver) Query() gql.QueryResolver { return &queryResolver{r} }

// Subscription returns gql.SubscriptionResolver implementation.
func (r *Resolver) Subscription() gql.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type playerResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
