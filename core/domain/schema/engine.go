package schema

import "picup/core/domain/constants"

type Engine struct {
	Key        string
	Provider   string
	Name       string
	Parameters map[string]string
}

type EngineMap map[string]*Engine

func (p EngineMap) GetDefault() *Engine {
	return p[constants.RESERVED_ENGINE_DEFAULT_KEY]
}
