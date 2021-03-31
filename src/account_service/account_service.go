package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/k-yomo/eitan/src/account_service/infra"
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
	accountID, err := a.sessionManager.Authenticate(req.SessionId)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	account, err := infra.AccountByID(ctx, a.db, accountID)
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("account with id: %s not found", accountID))
	}
	if err != nil {
		return nil, err
	}

	return &eitan.AuthenticateResponse{Account: mapToGRPCAccount(account)}, nil
}

func (a *AccountServiceServer) GetCurrentAccount(ctx context.Context, req *eitan.Empty) (*eitan.GetCurrentAccountResponse, error) {
	accountID, ok := sharedctx.GetAccountID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	account, err := infra.AccountByID(ctx, a.db, accountID)
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("account with id: %s not found", accountID))
	}
	if err != nil {
		return nil, err
	}

	return &eitan.GetCurrentAccountResponse{Account: mapToGRPCAccount(account)}, nil
}

func mapToGRPCAccount(a *infra.Account) *eitan.Account {
	return &eitan.Account{
		Id:           a.ID,
		Provider:     a.Provider,
		Email:        a.Email,
		DisplayName:  a.DisplayName,
		ScreenImgUrl: sqlutil.NullStrToPtr(a.ScreenImgURL),
	}
}
