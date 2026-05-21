package usecase

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"cinema/auth-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo    domain.UserRepository
	tokenCache  domain.TokenCache
	emailSender domain.EmailSender
	jwtSecret   []byte
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

func NewAuthUsecase(
	userRepo domain.UserRepository,
	tokenCache domain.TokenCache,
	emailSender domain.EmailSender,
) *AuthUsecase {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}
	return &AuthUsecase{
		userRepo:    userRepo,
		tokenCache:  tokenCache,
		emailSender: emailSender,
		jwtSecret:   []byte(secret),
		accessTTL:   15 * time.Minute,
		refreshTTL:  7 * 24 * time.Hour,
	}
}

func (uc *AuthUsecase) Register(ctx context.Context, email, password, fullName, phone string) (*domain.User, error) {
	if email == "" || password == "" || fullName == "" {
		return nil, fmt.Errorf("%w: email, password and full_name are required", domain.ErrInvalidInput)
	}
	if len(password) < 8 {
		return nil, fmt.Errorf("%w: password must be at least 8 characters", domain.ErrInvalidInput)
	}

	existing, _ := uc.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, domain.ErrEmailExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
		FullName:     fullName,
		Phone:        phone,
		Role:         domain.RoleUser,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	
	go uc.emailSender.SendWelcome(email, fullName)

	return user, nil
}

func (uc *AuthUsecase) Login(ctx context.Context, email, password string) (*domain.TokenPair, *domain.User, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, domain.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, domain.ErrUnauthorized
	}

	pair, err := uc.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	return pair, user, nil
}

func (uc *AuthUsecase) Logout(ctx context.Context, accessToken string) error {
	claims, err := uc.parseToken(accessToken)
	if err != nil {
		return domain.ErrInvalidToken
	}
	ttl := time.Until(time.Unix(claims.Exp, 0))
	if ttl > 0 {
		_ = uc.tokenCache.BlacklistToken(ctx, accessToken, ttl)
	}
	return nil
}

func (uc *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, *domain.User, error) {
	userID, err := uc.tokenCache.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, nil, domain.ErrInvalidToken
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, nil, domain.ErrInvalidToken
	}

	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, domain.ErrNotFound
	}

	_ = uc.tokenCache.DeleteRefreshToken(ctx, refreshToken)

	pair, err := uc.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	return pair, user, nil
}

func (uc *AuthUsecase) ValidateToken(ctx context.Context, token string) (*domain.Claims, error) {
	blacklisted, _ := uc.tokenCache.IsBlacklisted(ctx, token)
	if blacklisted {
		return nil, domain.ErrInvalidToken
	}
	return uc.parseToken(token)
}

func (uc *AuthUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return uc.userRepo.GetByID(ctx, userID)
}

func (uc *AuthUsecase) UpdateProfile(ctx context.Context, userID uuid.UUID, fullName, phone string) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if fullName != "" {
		user.FullName = fullName
	}
	if phone != "" {
		user.Phone = phone
	}
	user.UpdatedAt = time.Now()
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *AuthUsecase) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return uc.userRepo.Delete(ctx, userID)
}

func (uc *AuthUsecase) generateTokenPair(ctx context.Context, user *domain.User) (*domain.TokenPair, error) {
	// Access token
	now := time.Now()
	accessClaims := jwt.MapClaims{
		"sub":  user.ID.String(),
		"role": string(user.Role),
		"exp":  now.Add(uc.accessTTL).Unix(),
		"iat":  now.Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(uc.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Refresh token (random UUID stored in Redis)
	refreshToken := uuid.New().String() + strconv.FormatInt(now.UnixNano(), 36)
	if err := uc.tokenCache.StoreRefreshToken(ctx, user.ID.String(), refreshToken, uc.refreshTTL); err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *AuthUsecase) parseToken(tokenStr string) (*domain.Claims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return uc.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, domain.ErrInvalidToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.ErrInvalidToken
	}
	exp, _ := claims["exp"].(float64)
	return &domain.Claims{
		UserID: claims["sub"].(string),
		Role:   claims["role"].(string),
		Exp:    int64(exp),
	}, nil
}
