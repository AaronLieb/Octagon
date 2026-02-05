package ocr

import (
	"image"
	"image/color"
	"math"

	"github.com/charmbracelet/log"
	"github.com/disintegration/imaging"
)

func colorToRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

func diff(c1 color.Color, c2 color.Color) uint32 {
	ca := colorToRGBA(c1)
	cb := colorToRGBA(c2)
	return uint32((math.Abs(float64(ca.R-cb.R)) + math.Abs(float64(ca.G-cb.G)) + math.Abs(float64(ca.B-cb.B))) / 3)
}

// ColorDifference calculates perceptual color difference using Delta E CIE76
func ColorDifference(c1, c2 color.Color) float64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()

	// Convert to 0-1 range
	rf1, gf1, bf1 := float64(r1)/65535, float64(g1)/65535, float64(b1)/65535
	rf2, gf2, bf2 := float64(r2)/65535, float64(g2)/65535, float64(b2)/65535

	// Convert RGB to XYZ
	x1, y1, z1 := rgbToXYZ(rf1, gf1, bf1)
	x2, y2, z2 := rgbToXYZ(rf2, gf2, bf2)

	// Convert XYZ to LAB
	l1, a1, b1Lab := xyzToLab(x1, y1, z1)
	l2, a2, b2Lab := xyzToLab(x2, y2, z2)

	// Calculate Delta E
	return math.Sqrt(math.Pow(l2-l1, 2) + math.Pow(a2-a1, 2) + math.Pow(b2Lab-b1Lab, 2))
}

func rgbToXYZ(r, g, b float64) (float64, float64, float64) {
	// Gamma correction
	if r > 0.04045 {
		r = math.Pow((r+0.055)/1.055, 2.4)
	} else {
		r = r / 12.92
	}
	if g > 0.04045 {
		g = math.Pow((g+0.055)/1.055, 2.4)
	} else {
		g = g / 12.92
	}
	if b > 0.04045 {
		b = math.Pow((b+0.055)/1.055, 2.4)
	} else {
		b = b / 12.92
	}

	// Observer = 2°, Illuminant = D65
	x := r*0.4124564 + g*0.3575761 + b*0.1804375
	y := r*0.2126729 + g*0.7151522 + b*0.0721750
	z := r*0.0193339 + g*0.1191920 + b*0.9503041

	return x, y, z
}

func xyzToLab(x, y, z float64) (float64, float64, float64) {
	// Reference white D65
	xn, yn, zn := 0.95047, 1.00000, 1.08883

	x, y, z = x/xn, y/yn, z/zn

	if x > 0.008856 {
		x = math.Pow(x, 1.0/3.0)
	} else {
		x = (7.787*x + 16.0/116.0)
	}
	if y > 0.008856 {
		y = math.Pow(y, 1.0/3.0)
	} else {
		y = (7.787*y + 16.0/116.0)
	}
	if z > 0.008856 {
		z = math.Pow(z, 1.0/3.0)
	} else {
		z = (7.787*z + 16.0/116.0)
	}

	l := 116*y - 16
	a := 500 * (x - y)
	b := 200 * (y - z)

	return l, a, b
}

// PreprocessPercent applies image processing optimized for Smash Bros percentages
func PreprocessPercent(img image.Image, textColor color.RGBA) image.Image {
	log.Debug("preprocess", "textColor", textColor)

	if textColor.R > 233 && textColor.G > 224 && textColor.B > 224 {
		img = imaging.Crop(img, image.Rect(597, 921, 652, 992))
	} else if textColor.R > 200 && textColor.G > 10 {
		img = imaging.Crop(img, image.Rect(552, 921, 652, 992))
	}

	// img = imaging.Blur(img, 1)
	// img = imaging.AdjustContrast(img, 70)
	// img = imaging.Sharpen(img, 20)
	// img = imaging.Grayscale(img)

	bounds := img.Bounds()
	processed := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if ColorDifference(img.At(x, y), textColor) < 25 {
				processed.Set(x, y, textColor)
			} else {
				processed.Set(x, y, color.RGBA{255, 255, 255, 0})
			}
		}
	}

	// ymin := 934 - 917
	// ymax := 955 - 917
	// xmin := 627 - 507
	// xmax := 641 - 507
	// for y := ymin; y < ymax; y++ {
	// 	for x := xmin; x < xmax; x++ {
	// 		processed.Set(x, y, color.RGBA{0, 0, 0, 255})
	// 	}
	// }
	// ratio := float64(img.Bounds().Dx()) / float64(img.Bounds().Dy())
	// img = imaging.Resize(processed, int(40*ratio), int(40), imaging.NearestNeighbor)

	// return img
	return processed
}

// GetAverageNonBlackColor removes colors close to black and returns average of remaining colors
func GetAverageNonBlackColor(img image.Image) color.RGBA {
	bounds := img.Bounds()
	var totalR, totalG, totalB, count int

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			// Skip colors very close to black
			if r8 > 60 || g8 > 60 || b8 > 60 {
				totalR += int(r8)
				totalG += int(g8)
				totalB += int(b8)
				count++
			}
		}
	}

	if count == 0 {
		return color.RGBA{255, 255, 255, 255} // default to white
	}

	return color.RGBA{
		uint8(totalR / count),
		uint8(totalG / count),
		uint8(totalB / count),
		255,
	}
}

func GetRed(img image.Image) float64 {
	bounds := img.Bounds()
	count := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Max.X / 2; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
			dr8 := r8 - max(g8, b8)
			if dr8 > 20 {
				count += int(dr8)
			}
		}
	}
	return float64(count) / float64(bounds.Dx()*bounds.Dy())
}
