package uitest

import "image"

const (
	MATCH NeedleType = iota
	OCR
	EXCLUDE
)

type Comparison struct {
	BaseImage image.Image
	Needles   []*Needle
}

type NeedleType int

type Needle struct {
	X      int
	Y      int
	Width  int
	Height int
	Type   NeedleType
	Match  int
}

func NewComparison(img image.Image) *Comparison {
	return &Comparison{
		BaseImage: img,
		Needles:   make([]*Needle, 0),
	}
}

func (c *Comparison) Compare(with image.Image) bool {

}
