package main

import (
	"image"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type (
	exitStatus struct {
		kind exitCode
	}

	exitCode int
)

const (
	exitNoError exitCode = iota
)

func (e exitStatus) Error() string {
	switch e.kind {
	case exitNoError:
		return "No Error"
	default:
		return ""
	}
}

type (
	rectangle struct {
		x, y   float64
		width  float64
		height float64
	}

	point [2]float64

	Color [4]uint8

	constraint struct {
		Left, Right float64
		Up, Down    float64
	}
)

var (
	White = Color{255, 255, 255, 255}
	Black = Color{0, 0, 0, 255}
)

func (c Color) RGBA() (r, g, b, a uint32) {
	r = uint32(c[0])
	r |= r << 8
	g = uint32(c[1])
	g |= g << 8
	b = uint32(c[2])
	b |= b << 8
	a = uint32(c[3])
	a |= a << 8
	return
}

func (r rectangle) boundCheck(p point) bool {
	return (p[0] >= r.x && p[0] <= r.x+r.width) && (p[1] >= r.y && p[1] <= r.y+r.height)
}

type Font struct {
	faces map[int]font.Face
}

func NewFont(path string, dpi float64, sizes []int) Font {
	f := Font{
		faces: make(map[int]font.Face, len(sizes)),
	}

	fontData, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	tt, err := opentype.Parse(fontData)
	if err != nil {
		panic(err)
	}

	for _, v := range sizes {
		face, err := opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    float64(v),
			DPI:     dpi,
			Hinting: font.HintingNone,
		})
		if err != nil {
			panic(err)
		}
		f.faces[v] = face
	}
	return f
}

func (f *Font) GlyphAdvance(r rune, size float64) float64 {
	x, _ := f.faces[int(size)].GlyphAdvance(r)
	return float64(x>>6) + float64(x&((1<<6)-1))/float64(1<<6)
}

func (f *Font) Ascent(size float64) float64 {
	x := f.faces[int(size)].Metrics().Ascent
	return float64(x>>6) + float64(x&((1<<6)-1))/float64(1<<6)
}

func (f *Font) MeasureText(t string, size float64) point {
	measure := point{}

	if v, exist := f.faces[int(size)]; !exist {
		panic("No face of size in given Font")
	} else {
		r := text.BoundString(v, t)
		measure[0] = float64(r.Dx())
		measure[1] = float64(r.Dy())
	}

	return measure
}

func drawText(dst *ebiten.Image, f *Font, txt string, p point, textSize float64, clr Color) {
	ascent := f.Ascent(textSize)
	text.Draw(
		dst,
		txt,
		f.faces[int(textSize)],
		int(p[0]),
		int(p[1]+ascent),
		clr,
	)
}

func drawImageSlice(dst *ebiten.Image, dstRect rectangle, src *ebiten.Image, c constraint, clr Color) {
	dstRects := [9]rectangle{}
	srcRects := [9]rectangle{}

	l := c.Left
	r := c.Right
	u := c.Up
	d := c.Down

	imgW := float64(src.Bounds().Dx())
	imgH := float64(src.Bounds().Dy())

	srcX0 := float64(0)
	srcX1 := l
	srcX2 := imgW - r

	srcY0 := float64(0)
	srcY1 := u
	srcY2 := imgH - d

	dstL := l
	dstR := r
	dstU := u
	dstD := d

	// if scale > 0 {
	// 	dstL *= scale
	// 	dstR *= scale
	// 	dstU *= scale
	// 	dstD *= scale
	// }

	dstX0 := dstRect.x
	dstX1 := dstRect.x + dstL
	dstX2 := dstRect.x + dstRect.width - dstR

	dstY0 := dstRect.y
	dstY1 := dstRect.y + dstU
	dstY2 := dstRect.y + dstRect.height - dstD

	// TOP
	dstRects[0] = rectangle{x: dstX0, y: dstY0, width: dstL, height: dstU}
	srcRects[0] = rectangle{x: srcX0, y: srcY0, width: l, height: u}
	//
	dstRects[1] = rectangle{x: dstX1, y: dstY0, width: dstRect.width - (dstL + dstR), height: dstU}
	srcRects[1] = rectangle{x: srcX1, y: srcY0, width: imgW - (l + r), height: u}
	//
	dstRects[2] = rectangle{x: dstX2, y: dstY0, width: dstR, height: dstU}
	srcRects[2] = rectangle{x: srcX2, y: srcY0, width: r, height: u}
	//
	// MIDDLE
	dstRects[3] = rectangle{x: dstX0, y: dstY1, width: dstL, height: dstRect.height - (dstU + dstD)}
	srcRects[3] = rectangle{x: srcX0, y: srcY1, width: l, height: imgH - (u + d)}
	//
	dstRects[4] = rectangle{x: dstX1, y: dstY1, width: dstRect.width - (dstL + dstR), height: dstRect.height - (dstU + dstD)}
	srcRects[4] = rectangle{x: srcX1, y: srcY1, width: imgW - (l + r), height: imgH - (u + d)}
	//
	dstRects[5] = rectangle{x: dstX2, y: dstY1, width: dstR, height: dstRect.height - (dstU + dstD)}
	srcRects[5] = rectangle{x: srcX2, y: srcY1, width: r, height: imgH - (u + d)}
	//
	// BOTTOM
	dstRects[6] = rectangle{x: dstX0, y: dstY2, width: dstL, height: dstD}
	srcRects[6] = rectangle{x: srcX0, y: srcY2, width: l, height: d}
	//
	dstRects[7] = rectangle{x: dstX1, y: dstY2, width: dstRect.width - (dstL + dstR), height: dstD}
	srcRects[7] = rectangle{x: srcX1, y: srcY2, width: imgW - (l + r), height: d}
	//
	dstRects[8] = rectangle{x: dstX2, y: dstY2, width: dstR, height: dstD}
	srcRects[8] = rectangle{x: srcX2, y: srcY2, width: r, height: d}

	img := src
	for i := 0; i < 9; i += 1 {
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Scale(
			dstRects[i].width/srcRects[i].width,
			dstRects[i].height/srcRects[i].height,
		)
		opt.GeoM.Translate(dstRects[i].x, dstRects[i].y)
		if clr[3] != 0 {
			r, g, b, a := clr.RGBA()
			opt.ColorM.Scale(
				float64(r)/float64(a),
				float64(g)/float64(a),
				float64(b)/float64(a),
				float64(a)/0xffff,
			)
		}
		r := image.Rect(
			int(srcRects[i].x), int(srcRects[i].y),
			int(srcRects[i].x+srcRects[i].width), int(srcRects[i].y+srcRects[i].height),
		)
		dst.DrawImage(img.SubImage(r).(*ebiten.Image), opt)
	}
}

func drawRect(dst *ebiten.Image, r rectangle, clr Color) {
	ebitenutil.DrawRect(dst, r.x, r.y, r.width, r.height, clr)
}

type (
	minute  int
	seconds int

	timer struct {
		running bool
		min     minute
		sec     seconds
		timer   seconds
		buf     [5]rune
	}
)

func (t *timer) setDuration(m minute, s seconds) {
	t.min = m
	t.sec = s
	t.timer = 0
	t.updateString()
}

func (t *timer) start() {
	t.running = true
}

func (t *timer) advance() (finished bool) {
	const tps = 5
	finished = false
	if t.running {
		t.timer += 1
		if t.timer == tps {
			t.timer = 0
			if t.sec == 0 {
				if t.min == 0 && t.sec == 0 {
					t.running = false
					finished = true
				} else {
					t.min -= 1
					t.sec = 59
				}
			} else {
				t.sec -= 1
			}
			t.updateString()

		}
	}
	return
}

func (t *timer) updateString() []rune {
	toChar := func(n int) (r rune) {
		if n <= 10 {
			return rune(n + 48)
		} else {
			r = 0
		}
		return
	}

	if t.min >= 10 {
		rem := int(t.min)
		digit := rem % 10
		t.buf[1] = toChar(digit)
		rem /= 10
		digit = rem % 10
		t.buf[0] = toChar(digit)
	} else {
		t.buf[0] = '0'
		t.buf[1] = toChar(int(t.min))
	}

	if t.sec >= 10 {
		rem := int(t.sec)
		digit := rem % 10
		t.buf[4] = toChar(digit)
		rem /= 10
		digit = rem % 10
		t.buf[3] = toChar(digit)
	} else {
		t.buf[3] = '0'
		t.buf[4] = toChar(int(t.sec))
	}
	t.buf[2] = ':'
	return t.buf[:]
}

func (t timer) toString() []rune {
	return t.buf[:]
}
