// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package models provides cluster template models (for each sub-document of multi-doc YAML).
package models

import (
	"fmt"

	"github.com/cosi-project/runtime/pkg/resource"
)

// Meta is embedded into all template objects.
type Meta struct {
	Kind string `yaml:"kind"`
}

// TranslateContext is a context for translation.
type TranslateContext struct {
	// ClusterName is the name of the cluster.
	ClusterName string
}

// Model is a base interface for cluster templates.
type Model interface {
	Validate() error
	Translate(TranslateContext) ([]resource.Resource, error)
}

var registeredModels = map[string]func() Model{}

type model[T any] interface {
	*T
	Model
}

func register[T any, P model[T]](kind string) {
	if _, ok := registeredModels[kind]; ok {
		panic(fmt.Sprintf("model %s already registered", kind))
	}

	registeredModels[kind] = func() Model {
		return P(new(T))
	}
}

// New creates a model by kind.
func New(kind string) (Model, error) {
	f, ok := registeredModels[kind]
	if !ok {
		return nil, fmt.Errorf("unknown model kind %q", kind)
	}

	return f(), nil
}
