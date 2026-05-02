package schema

import "entgo.io/ent"

// DailyTask holds the schema definition for the DailyTask entity.
type DailyTask struct {
	ent.Schema
}

// Fields of the DailyTask.
func (DailyTask) Fields() []ent.Field {
	return nil
}

// Edges of the DailyTask.
func (DailyTask) Edges() []ent.Edge {
	return nil
}
