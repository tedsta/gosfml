package sf

import (
	"errors"
	"github.com/go-gl/gl"
	"image"
	"image/png"
	"os"
)

type Texture struct {
	t             gl.Texture
	size          Vector2
	isSmooth      bool   // Status of the smooth filter
	isRepeated    bool   // Is the texture in repeat mode?
	pixelsFlipped bool   // To work around the inconsistency in Y orientation
	cacheId       uint64 // Unique number that identifies the texture to the render target's cache
}

func NewTextureFromFile(fname string) *Texture {
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}

	t, err := CreateTexture(img)
	if err != nil {
		panic(err)
	}

	return t
}

func (t *Texture) Size() Vector2 {
	return t.size
}

type CoordType uint8

const (
	CoordNormalized CoordType = iota
	CoordPixels
)

// Bind binds the texture
func (t *Texture) Bind(coordType CoordType) {
	// ensureGlContext()

	if t != nil && t.t != 0 {
		// Bind the texture
		t.t.Bind(gl.TEXTURE_2D)

		// Check if we need to define a special texture matrix
		if coordType == CoordPixels || t.pixelsFlipped {
			matrix := [16]float32{1, 0, 0, 0,
				0, 1, 0, 0,
				0, 0, 1, 0,
				0, 0, 0, 1}

			// If non-normalized coordinates (= pixels) are requested, we need to
			// setup scale factors that convert the range [0 .. size] to [0 .. 1]
			if coordType == CoordPixels {
				matrix[0] = 1.0 / t.size.X
				matrix[5] = 1.0 / t.size.Y
			}

			// If pixels are flipped we must invert the Y axis
			if t.pixelsFlipped {
				matrix[5] = -matrix[5]
				matrix[13] = 1.0
			}

			// Load the matrix
			gl.MatrixMode(gl.TEXTURE)
			gl.LoadMatrixf(&matrix)

			// Go back to model-view mode (sf::RenderTarget relies on it)
			gl.MatrixMode(gl.MODELVIEW)
		}
	} else {
		// Bind no texture
		gl.Texture(0).Unbind(gl.TEXTURE_2D)

		// Reset the texture matrix
		gl.MatrixMode(gl.TEXTURE)
		gl.LoadIdentity()

		// Go back to model-view mode (sf::RenderTarget relies on it)
		gl.MatrixMode(gl.MODELVIEW)
	}
}

// Utilities ###################################################################

func CreateTexture(img image.Image) (*Texture, error) {
	imgW, imgH := img.Bounds().Dx(), img.Bounds().Dy()
	imgDim := Vector2{float32(imgW), float32(imgH)}

	rgbaImg, ok := img.(*image.NRGBA)
	if !ok {
		return nil, errors.New("texture must be an NRGBA image")
	}

	textureId := gl.GenTexture()
	textureId.Bind(gl.TEXTURE_2D)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)

	gl.TexImage2D(gl.TEXTURE_2D, 0, 4, imgW, imgH, 0, gl.RGBA, gl.UNSIGNED_BYTE, rgbaImg.Pix)

	return &Texture{textureId, imgDim, false, false, false, nextTextureCacheId()}, nil
}

// Unique cache id generator
// Thread-safe unique identifier generator,
// is used for states cache (see RenderTarget)
var nextTexCacheId uint64 = 0

func nextTextureCacheId() uint64 {
	nextTexCacheId++
	return nextTexCacheId
}
