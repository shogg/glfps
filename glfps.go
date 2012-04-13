// Simple fps overlay for opengl.
package glfps

import (
	"github.com/banthar/Go-SDL/sdl"
	"github.com/banthar/gl"
	"image"
	"image/png"
	"bytes"
)

var (
	digitImages []*image.NRGBA
	t0, frames uint32
	digitStore [5]int8
)

func init() {
	img, err := readImage()
	if err != nil { panic(err) }

	digitImages = cutImage(img)
}

// Draws the average fps (frames per second).
// (x, y) sets the right top corner of the drawing area.
// (0, 0) is left top on the screen.
func Draw(x, y int) {
	drawNumber(x, y, int(fps()))
}

func fps() uint32 {

	t := sdl.GetTicks()
	frames++

	if t0 == 0 { t0 = t; frames = 0; return 0 }
	if t - t0 > 2000 {
		t0 += (t - t0) / 2; frames /= 2
	}
	if t - t0 == 0 { return 0 }
	fps := 1000 * frames / (t - t0)
	return fps
}

func readImage() (rgb *image.NRGBA, err error) {

	b := bytes.NewBuffer(digits_png[:])
	img, err := png.Decode(b)
	if err != nil { return }

	rgb = img.(*image.NRGBA)
	return
}

func cutImage(img *image.NRGBA) []*image.NRGBA {
	r := image.Rect(0, 0, 6, 11)
	cuts := make([]*image.NRGBA, 10)
	for i := 0; i < 10; i++ {
		cuts[i] = img.SubImage(r).(*image.NRGBA)
		r = r.Add(image.Pt(7, 0))
	}

	return cuts
}

func countDigits(n int) (d int) {
	d = 0
	for n != 0 { n /= 10; d++ }
	return
}

func getDigits(n int) (digits []int8) {

	count := countDigits(n)
	digits = digitStore[:count]

	for i := 0; i < count; i++ {
		digits[i] = int8(n % 10)
		n /= 10
	}

	return
}

func drawNumber(x, y int, n int) {
	digits := getDigits(n)
	drawDigits(x, y, digits)
}

func drawDigits(x, y int, digits []int8) {

	s := sdl.GetVideoSurface()
	lft, rgt, btm, top := 0.0, float64(s.W), float64(s.H), 0.0

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.MatrixMode(gl.PROJECTION)
	gl.PushMatrix()
	gl.LoadIdentity()
	gl.Ortho(lft, rgt, btm, top, 0.0, 1.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	gl.LoadIdentity()

	w := digitImages[0].Rect.Dx()
	for i := 0; i < len(digits); i++ {
		img := digitImages[digits[i]]
		drawImage(x - w*i-i, y, img)
	}

	gl.MatrixMode(gl.MODELVIEW)
	gl.PopMatrix()

	gl.MatrixMode(gl.PROJECTION)
	gl.PopMatrix()
}

func drawImage(x, y int, img *image.NRGBA) {

	w, h := img.Rect.Dx(), img.Rect.Dy()
	stride := img.Stride

	for line := 0; line < h; line++ {
		gl.RasterPos2i(x, y + line)
		gl.DrawPixels(w, 1, gl.RGBA, gl.UNSIGNED_BYTE, &img.Pix[line*stride])
	}
}

