package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// Module holds the schema definition for the Module entity.
type Module struct {
	ent.Schema
}

// Fields of the Module.
func (Module) Fields() []ent.Field {
	return []ent.Field{
		field.String("owner"),
		field.String("namespace"),
		field.String("name"),
		field.String("provider"),
		field.String("description"),
		field.String("source"),
		field.Int64("downloads").Default(0),
		field.Time("published_at"),

		// For managing with Webhooks
		field.Int64("installation_id"),
		field.Int64("app_id"),
		field.String("repo_name"),
	}
}

// Edges of the Module.
func (Module) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("version", ModuleVersion.Type),
	}
}
