package coreGame

import "math"

func Distance(x1, y1, x2, y2 int) int {
	first := math.Pow(float64(x2-x1), 2)
	second := math.Pow(float64(y2-y1), 2)
	return int(math.Sqrt(first + second))
}

func (p *Instance) Line(x1, y1, x2, y2 int) {
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
		p.DrawPix(x, y)
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
			p.DrawPix(x, y)
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
		p.DrawPix(x, y)
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
			p.DrawPix(x, y)
		}
	}
}

func (p *Instance) Box(x1, y1, x2, y2 int) {
	for y := y1; y <= y2; y++ {
		p.DrawPix(x1, y)
		p.DrawPix(x2, y)
	}
	for x := x1; x <= x2; x++ {
		p.DrawPix(x, y1)
		p.DrawPix(x, y2)
	}
}

func (p *Instance) Circle(x0, y0, radius int) {
	x := radius
	y := 0
	e := 0

	for x >= y {
		p.DrawPix(x0+x, y0+y)
		p.DrawPix(x0+y, y0+x)
		p.DrawPix(x0-y, y0+x)
		p.DrawPix(x0-x, y0+y)
		p.DrawPix(x0-x, y0-y)
		p.DrawPix(x0-y, y0-x)
		p.DrawPix(x0+y, y0-x)
		p.DrawPix(x0+x, y0-y)

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

func (p *Instance) FilledCircle(x0, y0, radius int) {
	x := radius
	y := 0
	xChange := 1 - (radius << 1)
	yChange := 0
	radiusError := 0

	for x >= y {
		for i := x0 - x; i <= x0+x; i++ {
			p.DrawPix(i, y0+y)
			p.DrawPix(i, y0-y)
		}
		for i := x0 - y; i <= x0+y; i++ {
			p.DrawPix(i, y0+x)
			p.DrawPix(i, y0-x)
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

func (p *Instance) DrawChar(index, fgColor, bgColor byte, x, y int) {
	var a, b uint64
	for a = 0; a < 8; a++ {
		for b = 0; b < 8; b++ {
			if font.Bitmap[index][b]&(0x80>>a) != 0 {
				p.CurrentColor = fgColor
				p.DrawPix(int(a)+x, int(b)+y)
			} else {
				p.CurrentColor = bgColor
				p.DrawPix(int(a)+x, int(b)+y)
			}
		}
	}
}

func (p *Instance) DrawCursor(index, fgColor, bgColor byte, x, y int) {
	if cursorSetBlink {
		if cursorBlinkTimer < 15 {
			p.DrawChar(index, fgColor, bgColor, x, y)
		} else {
			p.DrawChar(index, bgColor, fgColor, x, y)
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
