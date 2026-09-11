package schema

import (
	"picup/core/domain/constants"
	"picup/core/domain/variantargs"
)

type Variant struct {
	Key      string
	Name     string
	Type     string
	Engine   string
	RawArgs  string
	Pipeline variantargs.Pipeline
}

func (v *Variant) IsMaster() bool {
	return v.Key == constants.RESERVED_VARIANT_MASTER_KEY
}

func (v *Variant) IsPreview() bool {
	return v.Key == constants.RESERVED_VARIANT_PREVIEW_KEY
}

func (v *Variant) IsDisplay() bool {
	return v.Key == constants.RESERVED_VARIANT_DISPLAY_KEY
}

func (v *Variant) IsReserved() bool {
	return v.IsMaster() || v.IsPreview() || v.IsDisplay()
}

type VariantMap map[string]*Variant

func (p VariantMap) GetMaster() *Variant {
	return p[constants.RESERVED_VARIANT_MASTER_KEY]
}

func (p VariantMap) GetDisplay() *Variant {
	return p[constants.RESERVED_VARIANT_DISPLAY_KEY]
}
