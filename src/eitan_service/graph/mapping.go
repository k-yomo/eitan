package graph

import (
	"github.com/k-yomo/eitan/src/eitan_service/graph/model"
	"github.com/k-yomo/eitan/src/eitan_service/infra"
	"github.com/k-yomo/eitan/src/internal/pb/eitan"
)

func mapToGraphqlCurrentUserProfile(up *eitan.CurrentUserProfile) *model.CurrentUserProfile {
	return &model.CurrentUserProfile{
		ID:           up.UserId,
		Email:        up.Email,
		DisplayName:  up.DisplayName,
		ScreenImgURL: up.ScreenImgUrl,
	}
}

func mapToGraphqlUserProfile(up *eitan.UserProfile) *model.UserProfile {
	return &model.UserProfile{
		ID:           up.UserId,
		DisplayName:  up.DisplayName,
		ScreenImgURL: up.ScreenImgUrl,
	}
}

func mapToGraphqlPlayer(p *infra.Player) *model.Player {
	return &model.Player{
		ID: p.ID,
		UserID: p.UserID,
	}
}

