package variant

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"picup/core/domain/schema"
	"picup/core/domain/variantargs"
	"picup/core/infrastructure/config"
	"slices"
	"strings"
	"sync"
)

type FactoryContext struct {
	Provider   string
	Name       string
	Parameters map[string]string
	WorkDir    string
}

// Factory constructs a VariantProcessor using raw configuration data.
type Factory func(rawConfig *FactoryContext) (VariantProcessor, error)

type Engine struct {
	mu        sync.RWMutex
	factories map[string]Factory
	instances map[string]VariantProcessor

	variants map[string]*schema.Variant
}

func (e *Engine) RegisterFactory(name string, factory Factory) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.factories[name] = factory
}

func (e *Engine) InitProcessors() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	var errs []error

	providerNames := slices.Collect(maps.Keys(e.factories))

	// 1. Initialize engine processors
	for ecname, ec := range config.Engine {
		factory, ok := e.factories[ec.Provider]
		if !ok {
			errs = append(errs, fmt.Errorf(
				"engine %q: unknown provider %q (available providers: %q)",
				ecname, ec.Provider, providerNames,
			))
			continue
		}

		processor, err := factory(&FactoryContext{
			Provider:   ec.Provider,
			Name:       ec.Name,
			Parameters: ec.Parameters,
			WorkDir:    config.Server.TempDir,
		})

		if err != nil {
			errs = append(errs, fmt.Errorf("engine %q: %w", ecname, err))
			continue
		}
		e.instances[ecname] = processor
	}

	if len(errs) > 0 {
		return fmt.Errorf("engine configuration failed:\n%w", errors.Join(errs...))
	}

	instanceNames := slices.Collect(maps.Keys(e.instances))

	// 2. Validate variants against initialized engines
	for vcname, vc := range config.Variant {
		if vc.IsMaster() {
			continue
		}

		if vc == nil {
			errs = append(errs, fmt.Errorf("variant %q: configuration is nil", vcname))
			continue
		}

		if vc.Engine == "" {
			errs = append(errs, fmt.Errorf(
				"variant %q: missing engine name (available engines: %q)",
				vcname, instanceNames,
			))
			continue
		}

		processor := e.instances[vc.Engine]
		if processor == nil {
			errs = append(errs, fmt.Errorf(
				"variant %q: references unknown engine %q (available engines: %q)",
				vcname, vc.Engine, instanceNames,
			))
			continue
		}

		if err := processor.Supports(vc); err != nil {
			errs = append(errs, fmt.Errorf("variant %q: %w", vcname, err))
			continue
		}

		e.variants[vc.Name] = vc
	}

	if len(errs) > 0 {
		return fmt.Errorf("variant validation failed:\n%w", errors.Join(errs...))
	}

	return nil
}

type processctx struct {
	context  context.Context
	pipeline variantargs.Pipeline
	workPath string
	workType string
	destPath string
	destType string

	sourceWidth  int
	sourceHeight int
}

func (c *processctx) Context() context.Context {
	return c.context
}

func (c *processctx) Pipeline() variantargs.Pipeline {
	return c.pipeline
}

func (c *processctx) WorkFileType() string {
	return strings.TrimPrefix(c.workType, ".")
}

func (c *processctx) DestFileType() string {
	return strings.TrimPrefix(c.destType, ".")
}

func (c *processctx) WorkFilePath() string {
	return c.workPath
}

func (c *processctx) DestFilePath() string {
	return c.destPath
}

func (c *processctx) SourceWidth() int {
	return c.sourceWidth
}

func (c *processctx) SourceHeight() int {
	return c.sourceHeight
}

func (e *Engine) Process(ctx context.Context, request ProcessRequest) error {
	e.mu.RLock()
	v, exists := e.variants[request.VariantName]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("varian didnt exist: %s", request.VariantName)
	}

	e.mu.RLock()
	p, exists := e.instances[v.Engine]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no processor configured for file extension: %s", v.Engine)
	}

	return p.Process(&processctx{
		context:  ctx,
		pipeline: v.Pipeline,
		workPath: request.WorkFilePath,
		workType: request.WorkFileType,
		destPath: request.DestFilePath,
		destType: request.DestFileType,

		sourceWidth:  request.SourceWidth,
		sourceHeight: request.SourceHeight,
	})
}

var globalEngine = NewEngine()

func NewEngine() *Engine {
	return &Engine{
		factories: make(map[string]Factory),
		instances: make(map[string]VariantProcessor),

		variants: make(map[string]*schema.Variant),
	}
}

// RegisterFactory is called inside each processor's init() function.
func RegisterFactory(name string, factory Factory) {
	globalEngine.RegisterFactory(name, factory)
}

// InitProcessors instantiates processors using the provided configuration map.
func InitProcessors() error {
	return globalEngine.InitProcessors()
}

func Process(ctx context.Context, request ProcessRequest) error {
	return globalEngine.Process(ctx, request)
}
