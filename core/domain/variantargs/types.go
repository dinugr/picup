package variantargs

type Pipeline struct {
	Operations []Operation
}

type Operation interface {
	opKind() string
}

type Dimension struct {
	Width  *int
	Height *int
}

type Scale struct {
	Size          Dimension
	Quality       int // 1..100
	Mode          ScaleMode
	Interpolation string
}

func (Scale) opKind() string { return "scale" }

type ScaleMode string

const (
	ScaleModeFixed           ScaleMode = "fixed"
	ScaleModeLockRatioExtend ScaleMode = "lrext"
	ScaleModeLockRatioShrink ScaleMode = "lrshr"
)

type Crop struct {
	Size    Dimension
	Gravity CropDirection
}

func (Crop) opKind() string { return "crop" }

type CropDirection struct {
	Vertical   Direction
	Horizontal Direction
}

type Direction string

const (
	DirectionNorth  Direction = "n"
	DirectionSouth  Direction = "s"
	DirectionCenter Direction = "c"
	DirectionWest   Direction = "w"
	DirectionEast   Direction = "e"
)
