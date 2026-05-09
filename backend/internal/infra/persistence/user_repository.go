package persistence

import (
	"context"

	"github.com/google/uuid"

	"github.com/hatamotoyuki/ninjo/backend/ent"
	entUser "github.com/hatamotoyuki/ninjo/backend/ent/user"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/model"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/repository"
)

// userRepository は ent を使った UserRepository の実装。
// domain の interface を満たす。
type userRepository struct {
	client *ent.Client
}

// NewUserRepository は UserRepository の実装を返す。
func NewUserRepository(client *ent.Client) repository.UserRepository {
	return &userRepository{client: client}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	created, err := r.client.User.
		Create().
		SetID(user.ID).
		SetEmail(user.Email).
		SetPasswordHash(user.PasswordHash).
		SetDisplayName(user.DisplayName).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toUserModel(created), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	found, err := r.client.User.
		Query().
		Where(entUser.EmailEQ(email)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return toUserModel(found), nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	found, err := r.client.User.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toUserModel(found), nil
}

// toUserModel は ent のエンティティを domain のモデルに変換する。
// ent の型が infra 層の外に漏れないようにする。
func toUserModel(e *ent.User) *model.User {
	return &model.User{
		ID:           e.ID,
		Email:        e.Email,
		PasswordHash: e.PasswordHash,
		DisplayName:  e.DisplayName,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}
