package service

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/tencat-dev/go-api-base/api/authz/v1"
	"github.com/tencat-dev/go-api-base/internal/biz"
)

type AuthzService struct {
	pb.UnimplementedAuthzServiceServer

	authzBiz *biz.AuthzBiz
}

func NewAuthzService(authzBiz *biz.AuthzBiz) pb.AuthzServiceServer {
	return &AuthzService{
		authzBiz: authzBiz,
	}
}

func (s *AuthzService) GrantRole(ctx context.Context, req *pb.GrantRoleRequest) (*emptypb.Empty, error) {
	if err := s.authzBiz.GrantRole(ctx, uuid.MustParse(req.Id), req.Role); err != nil {
		return nil, status.Errorf(codes.Internal, "grant role failed: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *AuthzService) RevokeRole(ctx context.Context, req *pb.RevokeRoleRequest) (*emptypb.Empty, error) {
	if err := s.authzBiz.RevokeRole(ctx, uuid.MustParse(req.Id), req.Role); err != nil {
		return nil, status.Errorf(codes.Internal, "revoke role failed: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *AuthzService) GrantPermission(ctx context.Context, req *pb.GrantPermissionRequest) (*emptypb.Empty, error) {
	if err := s.authzBiz.GrantPermission(ctx,
		req.Subject,
		req.Object,
		req.Action,
	); err != nil {
		return nil, status.Errorf(codes.Internal, "grant permission failed: %v", err)
	}

	return &emptypb.Empty{}, nil
}
