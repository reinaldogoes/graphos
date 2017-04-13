package coreGraphos

import "math"

func Distance(x1, y1, x2, y2 int) int {
	first := math.Pow(float64(x2-x1), 2)
	second := math.Pow(float64(y2-y1), 2)
	return int(math.Sqrt(first + second))
}

func bLine(x1, y1, x2, y2 int) {
	var x, y, dx, dy, dx1, dy1, px, py, xe, ye, i int
	dx = x2 - x1
	dy = y2 - y1
	if dx < 0 {
		dx1 = -dx
	} else {
		dx1 = dx
	}

	if dy < 0 {
		dy1 = -dy
	} else {
		dy1 = dy
	}
	px = 2*dy1 - dx1
	py = 2*dx1 - dy1
	if dy1 <= dx1 {
		if dx >= 0 {
			x = x1
			y = y1
			xe = x2
		} else {
			x = x2
			y = y2
			xe = x1
		}
		drawPix(x, y, 0xf)
		for i = 0; x < xe; i++ {
			x = x + 1
			if px < 0 {
				px = px + 2*dy1
			} else {
				if (dx < 0 && dy < 0) || (dx > 0 && dy > 0) {
					y = y + 1
				} else {
					y = y - 1
				}
				px = px + 2*(dy1-dx1)
			}
			drawPix(x, y, 0xf)
		}
	} else {
		if dy >= 0 {
			x = x1
			y = y1
			ye = y2
		} else {
			x = x2
			y = y2
			ye = y1
		}
		drawPix(x, y, 0xf)
		for i = 0; y < ye; i++ {
			y = y + 1
			if py <= 0 {
				py = py + 2*dx1
			} else {
				if (dx < 0 && dy < 0) || (dx > 0 && dy > 0) {
					x = x + 1
				} else {
					x = x - 1
				}
				py = py + 2*(dx1-dy1)
			}
			drawPix(x, y, 0xf)
		}
	}
}

func bBox(x1, y1, x2, y2 int) {
	for y := y1; y <= y2; y++ {
		drawPix(x1, y, 0xf)
		drawPix(x2, y, 0xf)
	}
	for x := x1; x <= x2; x++ {
		drawPix(x, y1, 0xf)
		drawPix(x, y2, 0xf)
	}
}

func bCircle(x0, y0, radius int) {
	x := radius
	y := 0
	e := 0

	for x >= y {
		drawPix(x0+x, y0+y, 0xf)
		drawPix(x0+y, y0+x, 0xf)
		drawPix(x0-y, y0+x, 0xf)
		drawPix(x0-x, y0+y, 0xf)
		drawPix(x0-x, y0-y, 0xf)
		drawPix(x0-y, y0-x, 0xf)
		drawPix(x0+y, y0-x, 0xf)
		drawPix(x0+x, y0-y, 0xf)

		if e <= 0 {
			y += 1
			e += 2*y + 1
		}
		if e > 0 {
			x -= 1
			e -= 2*x + 1
		}
	}
}

func bFilledCircle(x0, y0, radius int) {
	x := radius
	y := 0
	xChange := 1 - (radius << 1)
	yChange := 0
	radiusError := 0

	for x >= y {
		for i := x0 - x; i <= x0+x; i++ {
			drawPix(i, y0+y, 0xf)
			drawPix(i, y0-y, 0xf)
		}
		for i := x0 - y; i <= x0+y; i++ {
			drawPix(i, y0+x, 0xf)
			drawPix(i, y0-x, 0xf)
		}

		y++
		radiusError += yChange
		yChange += 2
		if ((radiusError << 1) + xChange) > 0 {
			x--
			radiusError += xChange
			xChange += 2
		}
	}
}

func drawChar(index, fgColor, bgColor byte, x, y int) {
	var a, b uint64
	for a = 0; a < 8; a++ {
		for b = 0; b < 8; b++ {
			if font.Bitmap[index][b]&(0x80>>a) != 0 {
				drawPix(int(a)+x, int(b)+y, fgColor)
			} else {
				drawPix(int(a)+x, int(b)+y, bgColor)
			}
		}
	}
}

func drawCursor(index, fgColor, bgColor byte, x, y int) {
	if cursorSetBlink {
		if cursorBlinkTimer < 15 {
			drawChar(index, fgColor, bgColor, x, y)
		} else {
			drawChar(index, bgColor, fgColor, x, y)
		}
		cursorBlinkTimer++
		if cursorBlinkTimer > 30 {
			cursorBlinkTimer = 0
		}
	} else {
		drawChar(index, bgColor, fgColor, x, y)
	}
}

func drawVideoTextMode() {
	i := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < columns; c++ {
			color := videoTextMemory[i]
			f := color & 0x0f
			b := color & 0xf0 >> 4
			i++
			if i-1 == cursor {
				drawCursor(videoTextMemory[i], f, b, c*8, r*8)
			} else {
				drawChar(videoTextMemory[i], f, b, c*8, r*8)
			}
			i++
		}
	}
}
