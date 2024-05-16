package auth

import (
	"context"
	sso "github.com/Rasikrr/protobuff/protos/gen/go/sso"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyValue = 0
)

type Auth interface {
	Login(ctx context.Context, email, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	sso.UnimplementedAuthServer
	auth      Auth
	validator *validator.Validate
}

func Register(gRPC *grpc.Server, auth Auth, val *validator.Validate) {
	sso.RegisterAuthServer(gRPC,
		&serverAPI{
			validator: val,
			auth:      auth,
		})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *sso.LoginRequest,
) (*sso.LoginResponse, error) {

	if err := s.validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		// TODO...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &sso.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *sso.RegisterRequest,
) (*sso.RegisterResponse, error) {
	if err := s.validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		// TODO...
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return &sso.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *sso.IsAdminRequest,
) (*sso.IsAdminResponse, error) {
	if err := s.validateIsAdmin(req); err != nil {
		return nil, err
	}
	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		// TODO...
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &sso.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func (s *serverAPI) validateLogin(req *sso.LoginRequest) error {
	if err := s.validator.Var(req.GetEmail(), "required,email"); err != nil {
		return status.Error(codes.InvalidArgument, "invalid email")
	}
	if err := s.validator.Var(
		req.GetPassword(),
		"required"); err != nil {
		return status.Error(
			codes.InvalidArgument,
			"password must contain at least one upper, one digit and one spec.symbol")
	}
	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}
	return nil
}

func (s *serverAPI) validateRegister(req *sso.RegisterRequest) error {
	if err := s.validator.Var(req.GetEmail(), "required,email"); err != nil {
		return status.Error(codes.InvalidArgument, "invalid email")
	}
	if err := s.validator.Var(
		req.GetPassword(),
		"required,min=8,contains_uppercase,contains_special"); err != nil {
		return status.Error(
			codes.InvalidArgument,
			"password must contain at least one upper, one digit and one spec.symbol")
	}
	return nil
}

func (s *serverAPI) validateIsAdmin(req *sso.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "invalid userID")
	}
	return nil
}
