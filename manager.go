// Package dao provides a data access object library.
//
// Copyright 2016 Pedro Salgado
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package dao

import "fmt"

// Manager interface for data access object managers.
type Manager interface {
	// CommitTransaction commits the context database transaction.
	CommitTransaction(ctx *Context) error
	// CreateDAO returns a new data access object.
	CreateDAO(ctx *Context, nm string) (interface{}, error)
	// EndTransaction ends the context database transaction.
	EndTransaction(ctx *Context)
	// RegisterDataSource registers a new mapping between a name and a data source.
	RegisterDataSource(ds *DataSource)
	// RegisterFactory registers the given factory with the data access object names it can generate.
	RegisterFactory(f Factory)
	// RollbackTransaction rollback the context database transaction.
	RollbackTransaction(ctx *Context) error
	// Source returns the data source associated with the given name.
	Source(nm string) *DataSource
	// StartTransaction start the context database transaction.
	StartTransaction() (*Context, error)
}

// BaseManager base implementation of the Manager interface.
type BaseManager struct {
	Sources   map[string]*DataSource
	Factories map[string]Factory
}

// CommitTransaction commits the context database transaction.
func (m *BaseManager) CommitTransaction(ctx *Context) error {
	for _, dao := range ctx.daos {
		if err := dao.Tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

// CreateDAO returns a new data access object.
func (m *BaseManager) CreateDAO(ctx *Context, nm string) (interface{}, error) {
	// factory
	fty, found := m.Factories[nm]
	if !found {
		return nil, fmt.Errorf("no factory for data access object %s", nm)
	}
	return fty.NewDataAccessObject(ctx, nm)
}

// EndTransaction ends the context database transaction.
func (m *BaseManager) EndTransaction(ctx *Context) {
	ctx.daos = map[string]*DataAccessObject{}
}

// RegisterDataSource registers a new mapping between a name and a data source.
func (m *BaseManager) RegisterDataSource(ds *DataSource) {
	m.Sources[ds.Name] = ds
}

// RegisterFactory registers the given factory with the data access object names it can generate.
func (m *BaseManager) RegisterFactory(f Factory) {
	for _, nm := range f.DataAccessObjects() {
		m.Factories[nm] = f
	}
}

// RollbackTransaction rollback the context database transaction.
func (m *BaseManager) RollbackTransaction(ctx *Context) error {
	var errout error

	for _, dao := range ctx.daos {
		if err := dao.Tx.Rollback(); err != nil {
			errout = err
		}
	}
	return errout
}

// Source returns the data source associated with the given name.
func (m *BaseManager) Source(nm string) *DataSource {
	return m.Sources[nm]
}

// StartTransaction start the context database transaction.
func (m *BaseManager) StartTransaction() (*Context, error) {
	return NewContext(m), nil
}

// NewBaseManager returns a generic Manager implementation.
func NewBaseManager() *BaseManager {
	return &BaseManager{
		Sources:   map[string]*DataSource{},
		Factories: map[string]Factory{},
	}
}

var _ Manager = (*BaseManager)(nil)
