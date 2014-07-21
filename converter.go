package uitest

import (
	"errors"
	"image"
	"image/color"
	"sync"

	"github.com/mitchellh/go-vnc"
)

var (
	// Use this map to add other encodings that the converter can understand
	TypeColors = make(map[int32]func(vnc.Encoding) []vnc.Color)

	// I don't know this encoding.
	ErrUnknownEncoding = errors.New("uitest: unknown encoding")
)

func init() {
	rawencode := &vnc.RawEncoding{}
	TypeColors[rawencode.Type()] = func(v vnc.Encoding) []vnc.Color {
		r := v.(*vnc.RawEncoding)
		return r.Colors
	}
}

type Converter struct {
	sync.Mutex

	conn    *vnc.ClientConn
	lastimg *image.NRGBA64
}

func NewConverter(conn *vnc.ClientConn) *Converter {
	return &Converter{
		conn:    conn,
		lastimg: image.NewNRGBA64(image.Rect(0, 0, int(conn.FrameBufferWidth), int(conn.FrameBufferHeight))),
	}
}

func (c *Converter) Process(m *vnc.FramebufferUpdateMessage) (image.Image, error) {
	c.Lock()
	defer c.Unlock()

	for _, v := range m.Rectangles {
		f, ok := TypeColors[v.Enc.Type()]
		if !ok {
			return nil, ErrUnknownEncoding
		}

		colors := f(v.Enc)

		index := 0
		for y := v.Y; y < v.Height; y++ {
			for x := v.X; x < v.Width; x++ {
				vcolor := colors[index]

				c.lastimg.SetNRGBA64(int(x), int(y), color.NRGBA64{
					R: vcolor.R,
					G: vcolor.G,
					B: vcolor.B,
					A: 255,
				})
			}
		}
	}

	// copy the image
	copy := image.NewNRGBA64(c.lastimg.Rect)
	for k, v := range c.lastimg.Pix {
		copy.Pix[k] = v
	}

	return copy, nil
}
