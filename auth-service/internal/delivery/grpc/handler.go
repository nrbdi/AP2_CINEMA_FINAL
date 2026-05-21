package grpc

import (
	"context"

	"cinema/auth-service/internal/domain"
	"cinema/auth-service/internal/usecase"
	pb "cinema/auth-service/proto/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	uc *usecase.AuthUsecase
}

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func mapErr(err error) error {
	switch err {
	case domain.ErrNotFound:
		return status.Error(codes.NotFound, err.Error())
	case domain.ErrUnauthorized, domain.ErrInvalidToken:
		return status.Error(codes.Unauthenticated, err.Error())
	case domain.ErrEmailExists:
		return status.Error(codes.AlreadyExists, err.Error())
	case domain.ErrInvalidInput:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func userToPB(u *domain.User) *pb.User {
	return &pb.User{
		Id:        u.ID.String(),
		Email:     u.Email,
		FullName:  u.FullName,
		Phone:     u.Phone,
		Role:      string(u.Role),
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, err := h.uc.Register(ctx, req.Email, req.Password, req.FullName, req.Phone)
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.RegisterResponse{User: userToPB(user), Message: "Registration successful"}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	pair, user, err := h.uc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.LoginResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		User:         userToPB(user),
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if err := h.uc.Logout(ctx, req.AccessToken); err != nil {
		return nil, mapErr(err)
	}
	return &pb.LogoutResponse{Success: true}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.LoginResponse, error) {
	pair, user, err := h.uc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.LoginResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		User:         userToPB(user),
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := h.uc.ValidateToken(ctx, req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}
	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: claims.UserID,
		Role:   claims.Role,
	}, nil
}

func (h *AuthHandler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.UserResponse, error) {
	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	user, err := h.uc.GetProfile(ctx, id)
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.UserResponse{User: userToPB(user)}, nil
}

func (h *AuthHandler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UserResponse, error) {
	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	user, err := h.uc.UpdateProfile(ctx, id, req.FullName, req.Phone)
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.UserResponse{User: userToPB(user)}, nil
}

func (h *AuthHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	if err := h.uc.DeleteUser(ctx, id); err != nil {
		return nil, mapErr(err)
	}
	return &pb.DeleteUserResponse{Success: true}, nil
}
