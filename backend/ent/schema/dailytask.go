package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"time"
)

// DailyTask holds the schema definition for the DailyTask entity.
type DailyTask struct {
	ent.Schema
}

// Fields of the DailyTask.
func (DailyTask) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Time("date"),
		field.Int("start_page").
			Positive(),
		field.Int("end_page").
			Positive(),
		field.Int("actual_end_page").
			Optional().
			Nillable().
			Positive(),
		field.Bool("is_completed").
			Default(false),
		field.Text("memo").
			Optional().
			Nillable(),
		field.Time("completed_at").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the DailyTask.
func (DailyTask) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("plan", Plan.Type).
			Ref("daily_tasks").
			Required().
			Unique(),
	}
}

// Indexes of the DailyTask.
func (DailyTask) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("plan").
			Fields("date"),
	}
}
