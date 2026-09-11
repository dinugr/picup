package config

import (
	"picup/core/domain/constants"
	"picup/core/domain/schema"
)

func registerReservedEngine(ecfg schema.EngineMap) {
	engine := &schema.Engine{
		Key:        constants.RESERVED_ENGINE_DEFAULT_KEY,
		Name:       constants.RESERVED_ENGINE_DEFAULT_KEY,
		Provider:   "ffmpeg",
		Parameters: map[string]string{},
	}

	val, exist := ecfg[constants.RESERVED_ENGINE_DEFAULT_KEY]
	if exist {
		engine.Name = fallbackValue(val.Name, engine.Name)
	}

	ecfg[constants.RESERVED_ENGINE_DEFAULT_KEY] = engine
}

func registerReservedVariant(vcfg schema.VariantMap) {

	// -- master ----------------------------
	master := newVariantHelper(&schema.Variant{
		Key:  constants.RESERVED_VARIANT_MASTER_KEY,
		Name: constants.RESERVED_VARIANT_MASTER_KEY,
	})

	master.applyName()
	master.registerTo(vcfg)

	// -- display ---------------------------

	display := newVariantHelper(&schema.Variant{
		Key:    constants.RESERVED_VARIANT_DISPLAY_KEY,
		Name:   constants.RESERVED_VARIANT_DISPLAY_KEY,
		Engine: constants.RESERVED_ENGINE_DEFAULT_KEY,
		Type:   "webp",
	})

	display.applyName()
	display.applyType()
	display.applyArguments()
	display.setArgs("scale=256::lrext,crop=256")

	display.registerTo(vcfg)

	// -- preview ---------------------------

	preview := newVariantHelper(&schema.Variant{
		Key:    constants.RESERVED_VARIANT_PREVIEW_KEY,
		Name:   constants.RESERVED_VARIANT_PREVIEW_KEY,
		Engine: constants.RESERVED_ENGINE_DEFAULT_KEY,
		Type:   "webp",
	})

	preview.applyName()
	preview.applyType()
	preview.applyArguments()
	preview.setArgs("scale=900:72:lrshr")
	preview.registerTo(vcfg)

}
