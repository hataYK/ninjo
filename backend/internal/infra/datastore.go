package infra

import (
	"github.com/hatamotoyuki/ninjo/backend/ent"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/repository"
	"github.com/hatamotoyuki/ninjo/backend/internal/infra/persistence"
)

// DataStore はリポジトリ層のファサード。
// 全リポジトリへのアクセスを一元管理する。
// 機能追加時はここにメソッドを足すだけでよい。
type DataStore struct {
	client *ent.Client
}

func NewDataStore(client *ent.Client) *DataStore {
	return &DataStore{client: client}
}

func (ds *DataStore) User() repository.UserRepository {
	return persistence.NewUserRepository(ds.client)
}

func (ds *DataStore) Availability() repository.AvailabilityRepository {
	return persistence.NewAvailabilityRepository(ds.client)
}

func (ds *DataStore) Plan() repository.PlanRepository {
	return persistence.NewPlanRepository(ds.client)
}

func (ds *DataStore) Skill() repository.SkillRepository {
	return persistence.NewSkillRepository(ds.client)
}
