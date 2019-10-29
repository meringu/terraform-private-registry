package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// ModuleVersion holds the schema definition for the ModuleVersion entity.
type ModuleVersion struct {
	ent.Schema
}

// Fields of the ModuleVersion.
func (ModuleVersion) Fields() []ent.Field {
	return []ent.Field{
		field.Int("major").Min(0),
		field.Int("minor").Min(0),
		field.Int("patch").Min(0),
		field.String("tag"),
	}
}

// Edges of the ModuleVersion.
func (ModuleVersion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("module", Module.Type).
			Ref("version").
			Unique().
			Required(),
	}
}
