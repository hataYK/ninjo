package schema

import "entgo.io/ent"

// Plan holds the schema definition for the Plan entity.
type Plan struct {
	ent.Schema
}

// Fields of the Plan.
func (Plan) Fields() []ent.Field {
	return nil
}

// Edges of the Plan.
func (Plan) Edges() []ent.Edge {
	return nil
}
