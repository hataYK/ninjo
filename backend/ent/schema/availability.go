package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Availability holds the schema definition for the Availability entity.
// 1ユーザーにつき1行。週の曜日ごとの可処分時間をまとめて管理する。
type Availability struct {
	ent.Schema
}

// Fields of the Availability.
func (Availability) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Float("sun_hours").Default(0).Min(0).Max(24),
		field.Float("mon_hours").Default(0).Min(0).Max(24),
		field.Float("tue_hours").Default(0).Min(0).Max(24),
		field.Float("wed_hours").Default(0).Min(0).Max(24),
		field.Float("thu_hours").Default(0).Min(0).Max(24),
		field.Float("fri_hours").Default(0).Min(0).Max(24),
		field.Float("sat_hours").Default(0).Min(0).Max(24),
	}
}

// Edges of the Availability.
func (Availability) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("availabilities").
			Required().
			Unique(),
	}
}
