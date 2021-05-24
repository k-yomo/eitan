package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/k-yomo/eitan/src/account_service/internal/infra"
	"github.com/k-yomo/eitan/src/account_service/internal/sessionmanager"
	"github.com/k-yomo/eitan/src/internal/pb/eitan"
	"github.com/k-yomo/eitan/src/internal/sharedctx"
	"github.com/k-yomo/eitan/src/pkg/sqlutil"
	"github.com/k-yomo/eitan/src/pkg/tx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountServiceServer struct {
	db             *sqlx.DB
	txManager      tx.Manager
	sessionManager sessionmanager.SessionManager
}

func NewAccountServiceServer(db *sqlx.DB, sessionManager sessionmanager.SessionManager) eitan.AccountServiceServer {
	return &AccountServiceServer{
		db:             db,
		txManager:      tx.NewManager(db),
		sessionManager: sessionManager,
	}
}

func (a *AccountServiceServer) Authenticate(ctx context.Context, req *eitan.AuthenticateRequest) (*eitan.AuthenticateResponse, error) {
	userID, err := a.sessionManager.Authenticate(req.SessionId)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	userProfile, err := infra.UserProfileByUserID(ctx, a.db, userID)
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("userProfile with id: %s not found", userID))
	}
	if err != nil {
		return nil, err
	}

	return &eitan.AuthenticateResponse{UserProfile: mapToGRPCCurrentUserProfile(userProfile)}, nil
}

func (a *AccountServiceServer) GetCurrentUserProfile(ctx context.Context, _ *eitan.Empty) (*eitan.GetCurrentUserProfileResponse, error) {
	userID, ok := sharedctx.GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	userProfile, err := infra.UserProfileByUserID(ctx, a.db, userID)
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("userProfile with id: '%s' not found", userID))
	}
	if err != nil {
		return nil, err
	}

	return &eitan.GetCurrentUserProfileResponse{UserProfile: mapToGRPCCurrentUserProfile(userProfile)}, nil
}

func (a *AccountServiceServer) GetUserProfile(ctx context.Context, req *eitan.GetUserProfileRequest) (*eitan.GetUserProfileResponse, error) {
	userProfile, err := infra.UserProfileByUserID(ctx, a.db, req.UserId)
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("userProfile with id: '%s' not found", req.UserId))
	}
	if err != nil {
		return nil, err
	}

	return &eitan.GetUserProfileResponse{UserProfile: mapToGRPCUserProfile(userProfile)}, nil
}

func mapToGRPCCurrentUserProfile(up *infra.UserProfile) *eitan.CurrentUserProfile {
	return &eitan.CurrentUserProfile{
		UserId:       up.UserID,
		Email:        up.Email,
		DisplayName:  up.DisplayName,
		ScreenImgUrl: sqlutil.NullStrToPtr(up.ScreenImgURL),
	}
}

func mapToGRPCUserProfile(up *infra.UserProfile) *eitan.UserProfile {
	return &eitan.UserProfile{
		UserId:       up.UserID,
		DisplayName:  up.DisplayName,
		ScreenImgUrl: sqlutil.NullStrToPtr(up.ScreenImgURL),
	}
}
