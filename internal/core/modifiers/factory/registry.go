package factory

import (
	"fmt"
	"sync"

	"github.com/SinfulDragon/dice-calculator/internal/core/common"
)

var GlobalRegistry = NewBuilderRegistry()

type BuilderRegistry struct {
	builders map[string]Builder
	mutex    sync.RWMutex
}

func NewBuilderRegistry() *BuilderRegistry {
	return &BuilderRegistry{
		builders: make(map[string]Builder),
	}
}

func (r *BuilderRegistry) Register(name string, builder Builder) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.builders[name] = builder
}

func (r *BuilderRegistry) Create(name string) (Builder, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	builder, ok := r.builders[name]
	return builder, ok
}

func (r *BuilderRegistry) Build(name string, args ModifierArgs) (common.Modifier, error) {
	builder, ok := r.Create(name)
	if !ok {
		return nil, fmt.Errorf("unknown modifier: %s", name)
	}
	return builder.Build(args)
}

func (r *BuilderRegistry) Schema(name string) (ModifierSchema, bool) {
	b, ok := r.Create(name)
	if !ok {
		return ModifierSchema{}, false
	}
	return b.Schema(), true
}
