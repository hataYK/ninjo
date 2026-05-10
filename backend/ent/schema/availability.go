package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Availability holds the schema definition for the Availability entity.
type Availability struct {
	ent.Schema
}

// Fields of the Availability.
func (Availability) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Int8("day_of_week").
			Range(0, 6),
		field.Float("hours").
			Min(0).
			Max(24),
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

// Indexes of the Availability.
func (Availability) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("user").
			Fields("day_of_week").
			Unique(),
	}
}
