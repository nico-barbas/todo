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
	White     = Color{255, 255, 255, 255}
	WhiteA125 = Color{255, 255, 255, 125}
	// Black     = Color{0, 0, 0, 255}
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

func (r rectangle) addPoint(p point) rectangle {
	return rectangle{
		x: r.x + p[0], y: r.y + p[1],
		width: r.width, height: r.height,
	}
}

func (r rectangle) boundCheck(p point) bool {
	return (p[0] >= r.x && p[0] <= r.x+r.width) && (p[1] >= r.y && p[1] <= r.y+r.height)
}

func (p point) sub(p2 point) point {
	return point{p[0] - p2[0], p[1] - p2[1]}
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

type textOptions struct {
	font   *Font
	text   string
	pos    point
	bounds rectangle
	size   float64
	clr    Color
}

func drawText(dst *ebiten.Image, opt textOptions) {
	ascent := opt.font.Ascent(opt.size)
	text.Draw(
		dst,
		opt.text,
		opt.font.faces[int(opt.size)],
		int(opt.pos[0]),
		int(opt.pos[1]+ascent),
		opt.clr,
	)
}

func drawTextCenter(dst *ebiten.Image, opt textOptions) {
	const yOffset = -2
	ascent := opt.font.Ascent(opt.size)
	textWidth := opt.font.MeasureText(opt.text, opt.size)[0]

	textPos := point{
		opt.bounds.x + (opt.bounds.width/2 - textWidth/2),
		opt.bounds.y + (opt.bounds.height/2 - ascent/2) + yOffset,
	}
	text.Draw(
		dst,
		opt.text,
		opt.font.faces[int(opt.size)],
		int(textPos[0]),
		int(textPos[1]+ascent),
		opt.clr,
	)
}

func drawImage(dst, src *ebiten.Image, p point) {
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(p[0], p[1])
	dst.DrawImage(src, &opt)
}

func drawImageCentered(dst, src *ebiten.Image, bounds rectangle, scale float64, clr Color) {
	opt := ebiten.DrawImageOptions{}
	// Color
	r, g, b, a := clr.RGBA()
	opt.ColorM.Scale(
		float64(r)/float64(a),
		float64(g)/float64(a),
		float64(b)/float64(a),
		float64(a)/0xffff,
	)

	opt.GeoM.Scale(scale, scale)
	// Position
	w := float64(src.Bounds().Dx()) * scale
	h := float64(src.Bounds().Dy()) * scale
	opt.GeoM.Translate(
		bounds.x+(bounds.width/2-w/2),
		bounds.y+(bounds.height/2-h/2),
	)

	dst.DrawImage(src, &opt)
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

func (t *timer) stop() {
	t.running = false
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

func (t *timer) updateString() {
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
}

// func (t timer) toString() []rune {
// 	return t.buf[:]
// }

func numberToString(n int, buf []rune) (last int) {
	toChar := func(n int) (r rune) {
		if n <= 10 {
			return rune(n + 48)
		} else {
			r = 0
		}
		return
	}

	rem := n
	if n < 10 {
		buf[0] = '0'
		buf[1] = toChar(n)
		last = 2
	} else {
		for rem > 0 {
			digit := rem % 10
			buf[last] = toChar(digit)
			rem /= 10
			last += 1
		}
		for i, j := 0, last-1; i < j; i, j = i+1, j-1 {
			buf[i], buf[j] = buf[j], buf[i]
		}
	}
	return
}

////////////////
////////////////
////////////////
const defaultTextBoxCap = 100

type textBox struct {
	font      *Font
	fontSize  float64
	charBuf   []rune
	charCount int
	cursor    rectangle
}

func (t *textBox) init(font *Font, fontSize float64) {
	t.font = font
	t.fontSize = fontSize
	t.cursor = rectangle{
		width:  2,
		height: textSize,
	}
	t.charBuf = make([]rune, defaultTextBoxCap)
}

func (t *textBox) update() {

}

func (t *textBox) AppendChar(r rune) {
	t.charBuf[t.charCount] = r
	t.charCount += 1
	t.cursor.x += t.font.GlyphAdvance(r, t.fontSize)
}

func (t *textBox) DeleteChar() {
	if t.charCount > 0 {
		r := t.charBuf[t.charCount-1]
		t.charCount -= 1
		t.cursor.x -= t.font.GlyphAdvance(r, t.fontSize)
	}
}

func (t *textBox) Clear() {
	t.cursor = rectangle{
		width:  2,
		height: textSize,
	}
	t.charCount = 0
}

func (t *textBox) GetText() []rune {
	return t.charBuf[:t.charCount]
}

////////////////
////////////////
////////////////

// Main downside is not being able to pop children rectLayout
// FIXME: Add checks when cutting if length isn't too big for
// the remaining rect size
type (
	rectLayout struct {
		full      rectangle
		remaining rectangle
	}

	rectCutKind int
)

const (
	rectCutUp rectCutKind = iota
	rectCutLeft
	rectCutDown
	rectCutRight
)

func newRectLayout(rect rectangle) rectLayout {
	return rectLayout{
		full:      rect,
		remaining: rect,
	}
}

func (r *rectLayout) x() float64 {
	return r.remaining.x
}

func (r *rectLayout) y() float64 {
	return r.remaining.y
}

func (r *rectLayout) width() float64 {
	return r.remaining.width
}

func (r *rectLayout) height() float64 {
	return r.remaining.height
}

func (r *rectLayout) cut(cutKind rectCutKind, length float64, padding float64) rectLayout {
	var result rectangle
	switch cutKind {
	case rectCutUp:
		result = rectangle{
			x: r.remaining.x, y: r.remaining.y,
			width: r.remaining.width, height: length,
		}
		r.remaining.y += length + padding
		r.remaining.height -= length + padding

	case rectCutLeft:
		result = rectangle{
			x: r.remaining.x, y: r.remaining.y,
			width: length, height: r.remaining.height,
		}
		r.remaining.x += length + padding
		r.remaining.width -= length + padding

	case rectCutDown:
		result = rectangle{
			x: r.remaining.x, y: r.remaining.y + (r.remaining.height - length),
			width: r.remaining.width, height: length,
		}
		r.remaining.height -= length + padding

	case rectCutRight:
		result = rectangle{
			x: r.remaining.x + (r.remaining.width - length), y: r.remaining.y,
			width: length, height: r.remaining.height,
		}
		r.remaining.width -= length + padding
	}
	return newRectLayout(result)
}

////////////////
////////////////
////////////////

type (
	rectArray struct {
		offset point
		rects  []rectElement
		focus  struct {
			active bool
			rect   rectangle
		}
		receiver rectReceiver
	}

	rectElement struct {
		userID rectID
		bounds rectangle
	}

	rectReceiver interface {
		onClick(userID rectID)
	}

	rectID int
)

func (r *rectArray) init(receiver rectReceiver, cap int) {
	r.receiver = receiver
	r.rects = make([]rectElement, 0, cap)
}

func (r *rectArray) setOffset(p point) {
	r.offset = p
}

func (r *rectArray) add(rect rectangle, userID rectID) {
	r.rects = append(r.rects, rectElement{userID, rect})
}

func (r *rectArray) update(mPos point, mLeft bool) {
	r.focus.active = false
	relPos := mPos.sub(r.offset)
	for _, rect := range r.rects {
		if rect.bounds.boundCheck(relPos) {
			r.focus.active = true
			r.focus.rect = rect.bounds
			if mLeft {
				r.receiver.onClick(rect.userID)
			}
			break
		}
	}
}

func (r *rectArray) highlight(dst *ebiten.Image) {
	if r.focus.active {
		drawRect(dst, r.focus.rect.addPoint(r.offset), WhiteA125)
	}
}
