package config

import "picup/core/domain/schema"

func fallbackValue[T comparable](value, fallback T) T {
	var zero T
	if value == zero {
		return fallback
	}
	return value
}

type variantHelper struct {
	variant   *schema.Variant
	applylist []func(v *schema.Variant)
}

func newVariantHelper(variant *schema.Variant) variantHelper {
	return variantHelper{
		variant: variant,
	}
}

func (c *variantHelper) applyType() {
	c.applylist = append(c.applylist, func(val *schema.Variant) {
		c.variant.Type = fallbackValue(val.Type, c.variant.Type)
	})
}

func (c *variantHelper) applyName() {
	c.applylist = append(c.applylist, func(val *schema.Variant) {
		c.variant.Name = fallbackValue(val.Name, c.variant.Name)
	})
}

func (c *variantHelper) applyArguments() {
	c.applylist = append(c.applylist, func(val *schema.Variant) {
		c.variant.RawArgs = fallbackValue(val.RawArgs, c.variant.RawArgs)
	})
}

func (c *variantHelper) setArgs(value string) {
	c.variant.RawArgs = value
}

func (c *variantHelper) registerTo(vcfg schema.VariantMap) {
	val, exist := vcfg[c.variant.Key]
	if exist {
		for _, apply := range c.applylist {
			apply(val)
		}
	}
	vcfg[c.variant.Key] = c.variant
}
