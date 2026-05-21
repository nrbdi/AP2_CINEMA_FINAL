// Hand-written stub replacing protoc output. Replace with: make proto
package auth

import (
	"context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var _ = codes.OK
var _ = status.New

// ── Messages ──────────────────────────────────────────────

type User struct {
	Id        string `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	FullName  string `json:"full_name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Role      string `json:"role,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

type RegisterRequest struct{ Email, Password, FullName, Phone string }
type RegisterResponse struct{ User *User; Message string }
type LoginRequest struct{ Email, Password string }
type LoginResponse struct{ AccessToken, RefreshToken string; User *User }
type LogoutRequest struct{ AccessToken string }
type LogoutResponse struct{ Success bool }
type RefreshTokenRequest struct{ RefreshToken string }
type ValidateTokenRequest struct{ Token string }
type ValidateTokenResponse struct{ Valid bool; UserId, Role string }
type GetProfileRequest struct{ UserId string }
type UpdateProfileRequest struct{ UserId, FullName, Phone string }
type UserResponse struct{ User *User }
type DeleteUserRequest struct{ UserId string }
type DeleteUserResponse struct{ Success bool }

// ── Server ────────────────────────────────────────────────

type AuthServiceServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	Logout(context.Context, *LogoutRequest) (*LogoutResponse, error)
	RefreshToken(context.Context, *RefreshTokenRequest) (*LoginResponse, error)
	ValidateToken(context.Context, *ValidateTokenRequest) (*ValidateTokenResponse, error)
	GetProfile(context.Context, *GetProfileRequest) (*UserResponse, error)
	UpdateProfile(context.Context, *UpdateProfileRequest) (*UserResponse, error)
	DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error)
	mustEmbedUnimplementedAuthServiceServer()
}

type UnimplementedAuthServiceServer struct{}

func (UnimplementedAuthServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Register not implemented")
}
func (UnimplementedAuthServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Login not implemented")
}
func (UnimplementedAuthServiceServer) Logout(context.Context, *LogoutRequest) (*LogoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Logout not implemented")
}
func (UnimplementedAuthServiceServer) RefreshToken(context.Context, *RefreshTokenRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "RefreshToken not implemented")
}
func (UnimplementedAuthServiceServer) ValidateToken(context.Context, *ValidateTokenRequest) (*ValidateTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "ValidateToken not implemented")
}
func (UnimplementedAuthServiceServer) GetProfile(context.Context, *GetProfileRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "GetProfile not implemented")
}
func (UnimplementedAuthServiceServer) UpdateProfile(context.Context, *UpdateProfileRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "UpdateProfile not implemented")
}
func (UnimplementedAuthServiceServer) DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "DeleteUser not implemented")
}
func (UnimplementedAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {}

func RegisterAuthServiceServer(s *grpc.Server, srv AuthServiceServer) {
	s.RegisterService(&AuthService_ServiceDesc, srv)
}

var AuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Register", Handler: _AuthService_Register_Handler},
		{MethodName: "Login", Handler: _AuthService_Login_Handler},
		{MethodName: "Logout", Handler: _AuthService_Logout_Handler},
		{MethodName: "RefreshToken", Handler: _AuthService_RefreshToken_Handler},
		{MethodName: "ValidateToken", Handler: _AuthService_ValidateToken_Handler},
		{MethodName: "GetProfile", Handler: _AuthService_GetProfile_Handler},
		{MethodName: "UpdateProfile", Handler: _AuthService_UpdateProfile_Handler},
		{MethodName: "DeleteUser", Handler: _AuthService_DeleteUser_Handler},
	},
	Streams: []grpc.StreamDesc{},
}

func _AuthService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(AuthServiceServer).Register(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/auth.AuthService/Register"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(AuthServiceServer).Register(ctx, req.(*RegisterRequest)) }
	return interceptor(ctx, in, info, handler)
}
func _AuthService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(AuthServiceServer).Login(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/auth.AuthService/Login"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(AuthServiceServer).Login(ctx, req.(*LoginRequest)) }
	return interceptor(ctx, in, info, handler)
}
func _AuthService_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(AuthServiceServer).Logout(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/auth.AuthService/Logout"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(AuthServiceServer).Logout(ctx, req.(*LogoutRequest)) }
	return interceptor(ctx, in, info, handler)
}
func _AuthService_RefreshToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshTokenRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(AuthServiceServer).RefreshToken(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/auth.AuthService/RefreshToken"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(AuthServiceServer).RefreshToken(ctx, req.(*RefreshTokenRequest)) }
	return interceptor(ctx, in, info, handler)
}
func _AuthService_ValidateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateTokenRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(AuthServiceServer).ValidateToken(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/auth.AuthService/ValidateToken"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(AuthServiceServer).ValidateToken(ctx, req.(*ValidateTokenRequest)) }
	return interceptor(ctx, in, info, handler)
}
func _AuthService_GetProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProfileRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(AuthServiceServer).GetProfile(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/auth.AuthService/GetProfile"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(AuthServiceServer).GetProfile(ctx, req.(*GetProfileRequest)) }
	return interceptor(ctx, in, info, handler)
}
func _AuthService_UpdateProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProfileRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(AuthServiceServer).UpdateProfile(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/auth.AuthService/UpdateProfile"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(AuthServiceServer).UpdateProfile(ctx, req.(*UpdateProfileRequest)) }
	return interceptor(ctx, in, info, handler)
}
func _AuthService_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(AuthServiceServer).DeleteUser(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/auth.AuthService/DeleteUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(AuthServiceServer).DeleteUser(ctx, req.(*DeleteUserRequest)) }
	return interceptor(ctx, in, info, handler)
}

// ── Client ────────────────────────────────────────────────

type AuthServiceClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error)
	RefreshToken(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	ValidateToken(ctx context.Context, in *ValidateTokenRequest, opts ...grpc.CallOption) (*ValidateTokenResponse, error)
	GetProfile(ctx context.Context, in *GetProfileRequest, opts ...grpc.CallOption) (*UserResponse, error)
	UpdateProfile(ctx context.Context, in *UpdateProfileRequest, opts ...grpc.CallOption) (*UserResponse, error)
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error)
}

type authServiceClient struct{ cc grpc.ClientConnInterface }

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient { return &authServiceClient{cc} }

func (c *authServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse); err := c.cc.Invoke(ctx, "/auth.AuthService/Register", in, out, opts...); return out, err
}
func (c *authServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse); err := c.cc.Invoke(ctx, "/auth.AuthService/Login", in, out, opts...); return out, err
}
func (c *authServiceClient) Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error) {
	out := new(LogoutResponse); err := c.cc.Invoke(ctx, "/auth.AuthService/Logout", in, out, opts...); return out, err
}
func (c *authServiceClient) RefreshToken(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse); err := c.cc.Invoke(ctx, "/auth.AuthService/RefreshToken", in, out, opts...); return out, err
}
func (c *authServiceClient) ValidateToken(ctx context.Context, in *ValidateTokenRequest, opts ...grpc.CallOption) (*ValidateTokenResponse, error) {
	out := new(ValidateTokenResponse); err := c.cc.Invoke(ctx, "/auth.AuthService/ValidateToken", in, out, opts...); return out, err
}
func (c *authServiceClient) GetProfile(ctx context.Context, in *GetProfileRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse); err := c.cc.Invoke(ctx, "/auth.AuthService/GetProfile", in, out, opts...); return out, err
}
func (c *authServiceClient) UpdateProfile(ctx context.Context, in *UpdateProfileRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse); err := c.cc.Invoke(ctx, "/auth.AuthService/UpdateProfile", in, out, opts...); return out, err
}
func (c *authServiceClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error) {
	out := new(DeleteUserResponse); err := c.cc.Invoke(ctx, "/auth.AuthService/DeleteUser", in, out, opts...); return out, err
}
