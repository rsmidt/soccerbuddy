package authz

import "context"

type RelationStore interface {
	AddRelations(ctx context.Context, relations []Relation) error
	RemoveRelations(ctx context.Context, relations []Relation) error
}

type Relation struct {
	SubjectType string
	SubjectID   string
	EntityType  string
	EntityID    string
	Relation    string
}

func (r Relation) String() string {
	return r.SubjectType + ":" + r.SubjectID + "#" + r.Relation + "@" + r.EntityType + ":" + r.EntityID
}

type RelationBuilder struct {
	relations []Relation
}

func (b *RelationBuilder) Entity(typ, id string) *RelationBuilder {
	b.relations = append(b.relations, Relation{EntityType: typ, EntityID: id})
	return b
}

func (b *RelationBuilder) Subject(typ, id string) *RelationBuilder {
	b.relations[len(b.relations)-1].SubjectType = typ
	b.relations[len(b.relations)-1].SubjectID = id
	return b
}

func (b *RelationBuilder) Relate(relation string) *RelationBuilder {
	b.relations[len(b.relations)-1].Relation = relation
	return b
}

func (b *RelationBuilder) And() *RelationBuilder {
	return b
}

func (b *RelationBuilder) Build() []Relation {
	return b.relations
}
