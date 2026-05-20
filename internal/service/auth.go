package service

import (
	"context"
	"time"

	"github.com/anhnmt/go-authxx/token"
	"github.com/google/uuid"

	pb "github.com/tencat-dev/go-api-base/api/auth/v1"
	"github.com/tencat-dev/go-api-base/internal/biz"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer

	authBiz    *biz.AuthBiz
	tokenMaker token.TokenMaker
}

func NewAuthService(authBiz *biz.AuthBiz, tokenMaker token.TokenMaker) pb.AuthServiceServer {
	return &AuthService{
		authBiz:    authBiz,
		tokenMaker: tokenMaker,
	}
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	user, err := s.authBiz.Login(ctx, &biz.AuthLogin{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.issueTokenPair(
		time.Now().UTC(),
		user.ID,
		uuid.Must(uuid.NewV7()),
	)
	if err != nil {
		return nil, err
	}

	return &pb.LoginReply{
		Id:           user.ID.String(),
		Email:        user.Email,
		Name:         user.Name,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) issueTokenPair(
	now time.Time,
	userID uuid.UUID,
	sessionID uuid.UUID,
) (string, string, error) {
	access, err := s.tokenMaker.CreateToken(now, token.TokenPayload{
		UserID:    userID,
		SessionID: sessionID,
		TokenID:   uuid.Must(uuid.NewV7()),
		Type:      token.AccessToken,
	})
	if err != nil {
		return "", "", err
	}

	refresh, err := s.tokenMaker.CreateToken(now, token.TokenPayload{
		UserID:    userID,
		SessionID: sessionID,
		TokenID:   uuid.Must(uuid.NewV7()),
		Type:      token.RefreshToken,
	})
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
