package schema

import "entgo.io/ent"

// Availability holds the schema definition for the Availability entity.
type Availability struct {
	ent.Schema
}

// Fields of the Availability.
func (Availability) Fields() []ent.Field {
	return nil
}

// Edges of the Availability.
func (Availability) Edges() []ent.Edge {
	return nil
}
