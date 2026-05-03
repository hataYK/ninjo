package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

// Plan holds the schema definition for the Plan entity.
type Plan struct {
	ent.Schema
}

// Fields of the Plan.
func (Plan) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("title").
			MaxLen(200).
			NotEmpty(),
		field.Int("total_pages").
			Positive(),
		field.Time("start_date"),
		field.Time("target_date"),
		field.Enum("status").
			Values("active", "completed", "paused").
			Default("active"),
		field.Text("ai_review").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Plan.
func (Plan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("plans").
			Required().
			Unique(),
		edge.To("daily_tasks", DailyTask.Type),
	}
}
