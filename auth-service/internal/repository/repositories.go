package postgres

import (
	"context"
	"time"

	"cinema/auth-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)


type UserRepo struct{ db *pgxpool.Pool }

func NewUserRepo(db *pgxpool.Pool) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (id,email,password_hash,full_name,phone,role,created_at,updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		u.ID, u.Email, u.PasswordHash, u.FullName, u.Phone, u.Role, u.CreatedAt, u.UpdatedAt,
	)
	return err
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u := &domain.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id,email,password_hash,full_name,phone,role,created_at,updated_at FROM users WHERE id=$1`, id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Phone, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return u, err
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id,email,password_hash,full_name,phone,role,created_at,updated_at FROM users WHERE email=$1`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Phone, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return u, err
}

func (r *UserRepo) Update(ctx context.Context, u *domain.User) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET full_name=$1, phone=$2, updated_at=$3 WHERE id=$4`,
		u.FullName, u.Phone, u.UpdatedAt, u.ID,
	)
	return err
}

func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	return err
}

type TokenCacheRedis struct{ client *redis.Client }

func NewTokenCache(client *redis.Client) *TokenCacheRedis {
	return &TokenCacheRedis{client: client}
}

func (c *TokenCacheRedis) BlacklistToken(ctx context.Context, token string, ttl time.Duration) error {
	return c.client.Set(ctx, "bl:"+token, "1", ttl).Err()
}

func (c *TokenCacheRedis) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	val, err := c.client.Exists(ctx, "bl:"+token).Result()
	return val > 0, err
}

func (c *TokenCacheRedis) StoreRefreshToken(ctx context.Context, userID, token string, ttl time.Duration) error {
	return c.client.Set(ctx, "rt:"+token, userID, ttl).Err()
}

func (c *TokenCacheRedis) GetRefreshToken(ctx context.Context, token string) (string, error) {
	return c.client.Get(ctx, "rt:"+token).Result()
}

func (c *TokenCacheRedis) DeleteRefreshToken(ctx context.Context, token string) error {
	return c.client.Del(ctx, "rt:"+token).Err()
}
