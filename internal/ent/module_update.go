// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/meringu/terraform-private-registry/internal/ent/module"
	"github.com/meringu/terraform-private-registry/internal/ent/moduleversion"
	"github.com/meringu/terraform-private-registry/internal/ent/predicate"
)

// ModuleUpdate is the builder for updating Module entities.
type ModuleUpdate struct {
	config
	owner              *string
	namespace          *string
	name               *string
	provider           *string
	description        *string
	source             *string
	downloads          *int64
	adddownloads       *int64
	published_at       *time.Time
	installation_id    *int64
	addinstallation_id *int64
	app_id             *int64
	addapp_id          *int64
	repo_name          *string
	version            map[int]struct{}
	removedVersion     map[int]struct{}
	predicates         []predicate.Module
}

// Where adds a new predicate for the builder.
func (mu *ModuleUpdate) Where(ps ...predicate.Module) *ModuleUpdate {
	mu.predicates = append(mu.predicates, ps...)
	return mu
}

// SetOwner sets the owner field.
func (mu *ModuleUpdate) SetOwner(s string) *ModuleUpdate {
	mu.owner = &s
	return mu
}

// SetNamespace sets the namespace field.
func (mu *ModuleUpdate) SetNamespace(s string) *ModuleUpdate {
	mu.namespace = &s
	return mu
}

// SetName sets the name field.
func (mu *ModuleUpdate) SetName(s string) *ModuleUpdate {
	mu.name = &s
	return mu
}

// SetProvider sets the provider field.
func (mu *ModuleUpdate) SetProvider(s string) *ModuleUpdate {
	mu.provider = &s
	return mu
}

// SetDescription sets the description field.
func (mu *ModuleUpdate) SetDescription(s string) *ModuleUpdate {
	mu.description = &s
	return mu
}

// SetSource sets the source field.
func (mu *ModuleUpdate) SetSource(s string) *ModuleUpdate {
	mu.source = &s
	return mu
}

// SetDownloads sets the downloads field.
func (mu *ModuleUpdate) SetDownloads(i int64) *ModuleUpdate {
	mu.downloads = &i
	mu.adddownloads = nil
	return mu
}

// SetNillableDownloads sets the downloads field if the given value is not nil.
func (mu *ModuleUpdate) SetNillableDownloads(i *int64) *ModuleUpdate {
	if i != nil {
		mu.SetDownloads(*i)
	}
	return mu
}

// AddDownloads adds i to downloads.
func (mu *ModuleUpdate) AddDownloads(i int64) *ModuleUpdate {
	if mu.adddownloads == nil {
		mu.adddownloads = &i
	} else {
		*mu.adddownloads += i
	}
	return mu
}

// SetPublishedAt sets the published_at field.
func (mu *ModuleUpdate) SetPublishedAt(t time.Time) *ModuleUpdate {
	mu.published_at = &t
	return mu
}

// SetInstallationID sets the installation_id field.
func (mu *ModuleUpdate) SetInstallationID(i int64) *ModuleUpdate {
	mu.installation_id = &i
	mu.addinstallation_id = nil
	return mu
}

// AddInstallationID adds i to installation_id.
func (mu *ModuleUpdate) AddInstallationID(i int64) *ModuleUpdate {
	if mu.addinstallation_id == nil {
		mu.addinstallation_id = &i
	} else {
		*mu.addinstallation_id += i
	}
	return mu
}

// SetAppID sets the app_id field.
func (mu *ModuleUpdate) SetAppID(i int64) *ModuleUpdate {
	mu.app_id = &i
	mu.addapp_id = nil
	return mu
}

// AddAppID adds i to app_id.
func (mu *ModuleUpdate) AddAppID(i int64) *ModuleUpdate {
	if mu.addapp_id == nil {
		mu.addapp_id = &i
	} else {
		*mu.addapp_id += i
	}
	return mu
}

// SetRepoName sets the repo_name field.
func (mu *ModuleUpdate) SetRepoName(s string) *ModuleUpdate {
	mu.repo_name = &s
	return mu
}

// AddVersionIDs adds the version edge to ModuleVersion by ids.
func (mu *ModuleUpdate) AddVersionIDs(ids ...int) *ModuleUpdate {
	if mu.version == nil {
		mu.version = make(map[int]struct{})
	}
	for i := range ids {
		mu.version[ids[i]] = struct{}{}
	}
	return mu
}

// AddVersion adds the version edges to ModuleVersion.
func (mu *ModuleUpdate) AddVersion(m ...*ModuleVersion) *ModuleUpdate {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return mu.AddVersionIDs(ids...)
}

// RemoveVersionIDs removes the version edge to ModuleVersion by ids.
func (mu *ModuleUpdate) RemoveVersionIDs(ids ...int) *ModuleUpdate {
	if mu.removedVersion == nil {
		mu.removedVersion = make(map[int]struct{})
	}
	for i := range ids {
		mu.removedVersion[ids[i]] = struct{}{}
	}
	return mu
}

// RemoveVersion removes version edges to ModuleVersion.
func (mu *ModuleUpdate) RemoveVersion(m ...*ModuleVersion) *ModuleUpdate {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return mu.RemoveVersionIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (mu *ModuleUpdate) Save(ctx context.Context) (int, error) {
	return mu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (mu *ModuleUpdate) SaveX(ctx context.Context) int {
	affected, err := mu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (mu *ModuleUpdate) Exec(ctx context.Context) error {
	_, err := mu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mu *ModuleUpdate) ExecX(ctx context.Context) {
	if err := mu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (mu *ModuleUpdate) sqlSave(ctx context.Context) (n int, err error) {
	selector := sql.Select(module.FieldID).From(sql.Table(module.Table))
	for _, p := range mu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = mu.driver.Query(ctx, query, args, rows); err != nil {
		return 0, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("ent: failed reading id: %v", err)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return 0, nil
	}

	tx, err := mu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		builder = sql.Update(module.Table).Where(sql.InInts(module.FieldID, ids...))
	)
	if value := mu.owner; value != nil {
		builder.Set(module.FieldOwner, *value)
	}
	if value := mu.namespace; value != nil {
		builder.Set(module.FieldNamespace, *value)
	}
	if value := mu.name; value != nil {
		builder.Set(module.FieldName, *value)
	}
	if value := mu.provider; value != nil {
		builder.Set(module.FieldProvider, *value)
	}
	if value := mu.description; value != nil {
		builder.Set(module.FieldDescription, *value)
	}
	if value := mu.source; value != nil {
		builder.Set(module.FieldSource, *value)
	}
	if value := mu.downloads; value != nil {
		builder.Set(module.FieldDownloads, *value)
	}
	if value := mu.adddownloads; value != nil {
		builder.Add(module.FieldDownloads, *value)
	}
	if value := mu.published_at; value != nil {
		builder.Set(module.FieldPublishedAt, *value)
	}
	if value := mu.installation_id; value != nil {
		builder.Set(module.FieldInstallationID, *value)
	}
	if value := mu.addinstallation_id; value != nil {
		builder.Add(module.FieldInstallationID, *value)
	}
	if value := mu.app_id; value != nil {
		builder.Set(module.FieldAppID, *value)
	}
	if value := mu.addapp_id; value != nil {
		builder.Add(module.FieldAppID, *value)
	}
	if value := mu.repo_name; value != nil {
		builder.Set(module.FieldRepoName, *value)
	}
	if !builder.Empty() {
		query, args := builder.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(mu.removedVersion) > 0 {
		eids := make([]int, len(mu.removedVersion))
		for eid := range mu.removedVersion {
			eids = append(eids, eid)
		}
		query, args := sql.Update(module.VersionTable).
			SetNull(module.VersionColumn).
			Where(sql.InInts(module.VersionColumn, ids...)).
			Where(sql.InInts(moduleversion.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(mu.version) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range mu.version {
				p.Or().EQ(moduleversion.FieldID, eid)
			}
			query, args := sql.Update(module.VersionTable).
				Set(module.VersionColumn, id).
				Where(sql.And(p, sql.IsNull(module.VersionColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(mu.version) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"version\" %v already connected to a different \"Module\"", keys(mu.version))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// ModuleUpdateOne is the builder for updating a single Module entity.
type ModuleUpdateOne struct {
	config
	id                 int
	owner              *string
	namespace          *string
	name               *string
	provider           *string
	description        *string
	source             *string
	downloads          *int64
	adddownloads       *int64
	published_at       *time.Time
	installation_id    *int64
	addinstallation_id *int64
	app_id             *int64
	addapp_id          *int64
	repo_name          *string
	version            map[int]struct{}
	removedVersion     map[int]struct{}
}

// SetOwner sets the owner field.
func (muo *ModuleUpdateOne) SetOwner(s string) *ModuleUpdateOne {
	muo.owner = &s
	return muo
}

// SetNamespace sets the namespace field.
func (muo *ModuleUpdateOne) SetNamespace(s string) *ModuleUpdateOne {
	muo.namespace = &s
	return muo
}

// SetName sets the name field.
func (muo *ModuleUpdateOne) SetName(s string) *ModuleUpdateOne {
	muo.name = &s
	return muo
}

// SetProvider sets the provider field.
func (muo *ModuleUpdateOne) SetProvider(s string) *ModuleUpdateOne {
	muo.provider = &s
	return muo
}

// SetDescription sets the description field.
func (muo *ModuleUpdateOne) SetDescription(s string) *ModuleUpdateOne {
	muo.description = &s
	return muo
}

// SetSource sets the source field.
func (muo *ModuleUpdateOne) SetSource(s string) *ModuleUpdateOne {
	muo.source = &s
	return muo
}

// SetDownloads sets the downloads field.
func (muo *ModuleUpdateOne) SetDownloads(i int64) *ModuleUpdateOne {
	muo.downloads = &i
	muo.adddownloads = nil
	return muo
}

// SetNillableDownloads sets the downloads field if the given value is not nil.
func (muo *ModuleUpdateOne) SetNillableDownloads(i *int64) *ModuleUpdateOne {
	if i != nil {
		muo.SetDownloads(*i)
	}
	return muo
}

// AddDownloads adds i to downloads.
func (muo *ModuleUpdateOne) AddDownloads(i int64) *ModuleUpdateOne {
	if muo.adddownloads == nil {
		muo.adddownloads = &i
	} else {
		*muo.adddownloads += i
	}
	return muo
}

// SetPublishedAt sets the published_at field.
func (muo *ModuleUpdateOne) SetPublishedAt(t time.Time) *ModuleUpdateOne {
	muo.published_at = &t
	return muo
}

// SetInstallationID sets the installation_id field.
func (muo *ModuleUpdateOne) SetInstallationID(i int64) *ModuleUpdateOne {
	muo.installation_id = &i
	muo.addinstallation_id = nil
	return muo
}

// AddInstallationID adds i to installation_id.
func (muo *ModuleUpdateOne) AddInstallationID(i int64) *ModuleUpdateOne {
	if muo.addinstallation_id == nil {
		muo.addinstallation_id = &i
	} else {
		*muo.addinstallation_id += i
	}
	return muo
}

// SetAppID sets the app_id field.
func (muo *ModuleUpdateOne) SetAppID(i int64) *ModuleUpdateOne {
	muo.app_id = &i
	muo.addapp_id = nil
	return muo
}

// AddAppID adds i to app_id.
func (muo *ModuleUpdateOne) AddAppID(i int64) *ModuleUpdateOne {
	if muo.addapp_id == nil {
		muo.addapp_id = &i
	} else {
		*muo.addapp_id += i
	}
	return muo
}

// SetRepoName sets the repo_name field.
func (muo *ModuleUpdateOne) SetRepoName(s string) *ModuleUpdateOne {
	muo.repo_name = &s
	return muo
}

// AddVersionIDs adds the version edge to ModuleVersion by ids.
func (muo *ModuleUpdateOne) AddVersionIDs(ids ...int) *ModuleUpdateOne {
	if muo.version == nil {
		muo.version = make(map[int]struct{})
	}
	for i := range ids {
		muo.version[ids[i]] = struct{}{}
	}
	return muo
}

// AddVersion adds the version edges to ModuleVersion.
func (muo *ModuleUpdateOne) AddVersion(m ...*ModuleVersion) *ModuleUpdateOne {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return muo.AddVersionIDs(ids...)
}

// RemoveVersionIDs removes the version edge to ModuleVersion by ids.
func (muo *ModuleUpdateOne) RemoveVersionIDs(ids ...int) *ModuleUpdateOne {
	if muo.removedVersion == nil {
		muo.removedVersion = make(map[int]struct{})
	}
	for i := range ids {
		muo.removedVersion[ids[i]] = struct{}{}
	}
	return muo
}

// RemoveVersion removes version edges to ModuleVersion.
func (muo *ModuleUpdateOne) RemoveVersion(m ...*ModuleVersion) *ModuleUpdateOne {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return muo.RemoveVersionIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (muo *ModuleUpdateOne) Save(ctx context.Context) (*Module, error) {
	return muo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (muo *ModuleUpdateOne) SaveX(ctx context.Context) *Module {
	m, err := muo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return m
}

// Exec executes the query on the entity.
func (muo *ModuleUpdateOne) Exec(ctx context.Context) error {
	_, err := muo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (muo *ModuleUpdateOne) ExecX(ctx context.Context) {
	if err := muo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (muo *ModuleUpdateOne) sqlSave(ctx context.Context) (m *Module, err error) {
	selector := sql.Select(module.Columns...).From(sql.Table(module.Table))
	module.ID(muo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = muo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		m = &Module{config: muo.config}
		if err := m.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into Module: %v", err)
		}
		id = m.ID
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("Module with id: %v", muo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one Module with the same id: %v", muo.id)
	}

	tx, err := muo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		builder = sql.Update(module.Table).Where(sql.InInts(module.FieldID, ids...))
	)
	if value := muo.owner; value != nil {
		builder.Set(module.FieldOwner, *value)
		m.Owner = *value
	}
	if value := muo.namespace; value != nil {
		builder.Set(module.FieldNamespace, *value)
		m.Namespace = *value
	}
	if value := muo.name; value != nil {
		builder.Set(module.FieldName, *value)
		m.Name = *value
	}
	if value := muo.provider; value != nil {
		builder.Set(module.FieldProvider, *value)
		m.Provider = *value
	}
	if value := muo.description; value != nil {
		builder.Set(module.FieldDescription, *value)
		m.Description = *value
	}
	if value := muo.source; value != nil {
		builder.Set(module.FieldSource, *value)
		m.Source = *value
	}
	if value := muo.downloads; value != nil {
		builder.Set(module.FieldDownloads, *value)
		m.Downloads = *value
	}
	if value := muo.adddownloads; value != nil {
		builder.Add(module.FieldDownloads, *value)
		m.Downloads += *value
	}
	if value := muo.published_at; value != nil {
		builder.Set(module.FieldPublishedAt, *value)
		m.PublishedAt = *value
	}
	if value := muo.installation_id; value != nil {
		builder.Set(module.FieldInstallationID, *value)
		m.InstallationID = *value
	}
	if value := muo.addinstallation_id; value != nil {
		builder.Add(module.FieldInstallationID, *value)
		m.InstallationID += *value
	}
	if value := muo.app_id; value != nil {
		builder.Set(module.FieldAppID, *value)
		m.AppID = *value
	}
	if value := muo.addapp_id; value != nil {
		builder.Add(module.FieldAppID, *value)
		m.AppID += *value
	}
	if value := muo.repo_name; value != nil {
		builder.Set(module.FieldRepoName, *value)
		m.RepoName = *value
	}
	if !builder.Empty() {
		query, args := builder.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(muo.removedVersion) > 0 {
		eids := make([]int, len(muo.removedVersion))
		for eid := range muo.removedVersion {
			eids = append(eids, eid)
		}
		query, args := sql.Update(module.VersionTable).
			SetNull(module.VersionColumn).
			Where(sql.InInts(module.VersionColumn, ids...)).
			Where(sql.InInts(moduleversion.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(muo.version) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range muo.version {
				p.Or().EQ(moduleversion.FieldID, eid)
			}
			query, args := sql.Update(module.VersionTable).
				Set(module.VersionColumn, id).
				Where(sql.And(p, sql.IsNull(module.VersionColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(muo.version) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"version\" %v already connected to a different \"Module\"", keys(muo.version))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return m, nil
}
