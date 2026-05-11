package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"time"
)

// Skill holds the schema definition for the Skill entity.
type Skill struct {
	ent.Schema
}

// Fields of the Skill.
func (Skill) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("name").
			MaxLen(200).
			NotEmpty(),
		field.String("category").
			MaxLen(100).
			NotEmpty(),
		field.Enum("source").
			Values("ai", "manual").
			Default("manual"),
		field.UUID("task_id", uuid.UUID{}).
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Skill.
func (Skill) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("skills").
			Required().
			Unique(),
	}
}

// Indexes of the Skill.
func (Skill) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("user"),
		index.Edges("user").
			Fields("category"),
		index.Edges("user").
			Fields("name").
			Unique(),
	}
}
