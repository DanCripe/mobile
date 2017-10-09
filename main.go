package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"strconv"
	// "io"

	"github.com/golang/freetype"

	"github.com/dancripe/tmo/game"

	// "github.com/cryptix/wav"
	imgfont "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	// "golang.org/x/mobile/exp/audio/al"
	"golang.org/x/mobile/exp/font"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

var (
	scale        float32 = 100.0
	pxPerPt      float32 = 1.0
	pxScale      int     = 100
	program      gl.Program
	grid         *game.Grid
	offset       gl.Uniform
	position     gl.Attrib
	win          bool
	context      *freetype.Context
	sz           size.Event
	images       *glutil.Images
	atlasImage   *glutil.Image
	controlImage *glutil.Image
	checkImage   *glutil.Image
	// winImage       *glutil.Image
	countImage     *glutil.Image
	colorIndex     map[string]int
	textWidth      float32
	heightOffsetPt float32
	heightOffsetPx int
	currentIndex   int
	index3         int
	index4         int
	muted          bool
	GridSize       int = 3
	dpi            float32
	xOffset        float32
	// source         al.Source
	// noSound        bool
)

const (
	AllImages     string = "texatlas.png"
	ControlImages string = "control.png"
	CheckImage    string = "check.png"
	ClickWav      string = "shortclick.wav"
)

const (
	ControlIndexRestart = 0
	ControlIndexSound   = 1
	ControlIndexCount   = 2
	ControlIndexSize    = 3
	ControlIndexSkip    = 4
)

const (
	UnmutedIndex = 1
	MutedIndex   = 5
	Grid3Index   = 3
	Grid4Index   = 6
)

func main() {
	app.Main(func(a app.App) {
		colorIndex = game.GetColorIndex()

		/*
			err := al.OpenDevice()
			if err != nil {
				log.Printf("No sound: %v", err)
				noSound = true
			}
		*/
		win = false
		var glctx gl.Context

		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					if grid == nil {
						grid, _ = game.Random(GridSize, currentIndex)
					}
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop(glctx)
					glctx = nil
				}
			case size.Event:
				sz = e
				inches := float32(e.HeightPt) / 72.0
				dpi = float32(e.HeightPx) / inches
				dpScale := dpi / 160.0
				heightOffsetPx = 36 * int(dpScale)
				heightOffsetPt = 36.0 * dpScale * float32(e.HeightPt) / float32(e.HeightPx)

				pxPerPt = float32(e.HeightPx) / float32(e.HeightPt)

				HeightPt := float32(e.HeightPt) - heightOffsetPt
				HeightPx := e.HeightPx - heightOffsetPx
				scale = float32(e.WidthPt)
				if HeightPt < scale {
					scale = HeightPt
				}
				pxScale = e.WidthPx
				if HeightPx < pxScale {
					pxScale = HeightPx
				}
				onResize()
				a.Send(paint.Event{})
			case paint.Event:
				if glctx == nil || e.External {
					continue
				}
				onPaint(glctx)
				a.Publish()
			case touch.Event:
				if e.Type == touch.TypeBegin {
					if win {
						if GridSize == 3 {
							index3++
							app.SetKeyValue("index3", fmt.Sprintf("%d", index3))
							currentIndex = index3
						} else {
							index4++
							app.SetKeyValue("index4", fmt.Sprintf("%d", index4))
							currentIndex = index4
						}
						grid, _ = game.Random(GridSize, currentIndex)
						win = false
						a.Send(paint.Event{})
					} else {
						// TODO: this is the wrong way to do this.
						// CalculateLocation is useless
						// location, valid := grid.CalculateLocation(int(e.X), int(e.Y)-heightOffsetPx, pxScale)
						x, y, valid := gridClick(e.X, e.Y)
						if valid {
							// play(source)
							win = grid.Click(game.Location{x, y})
							a.Send(paint.Event{})
						} else {
							// it may be a control click
							if controlClick(e.X, e.Y) {
								a.Send(paint.Event{})
							}
						}
					}
				}
			}
		}
	})
}

/*
func play(source al.Source) {
	if noSound {
		return
	}
	al.PlaySources(source)
}
*/

func onStart(glctx gl.Context) {
	var err error
	str := app.GetKeyValue("index3")
	if str == "" {
		index3 = 1
	} else {
		index3, _ = strconv.Atoi(str)
	}
	str = app.GetKeyValue("index4")
	if str == "" {
		index4 = 1
	} else {
		index4, _ = strconv.Atoi(str)
	}
	str = app.GetKeyValue("gridSize")
	if str == "4" {
		GridSize = 4
		currentIndex = index4
	} else {
		GridSize = 3
		currentIndex = index3
	}
	str = app.GetKeyValue("muted")
	if str == "true" {
		muted = true
	} else {
		muted = false
	}

	program, err = glutil.CreateProgram(glctx, vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	images = glutil.NewImages(glctx)
	atlasImage = images.NewImage(1975, 175)
	// read img from png
	f, err := asset.Open(AllImages)
	if err != nil {
		log.Printf("error opening images: %v", err)
		return
	}
	defer f.Close()

	imageData, err := png.Decode(f)
	if err != nil {
		log.Printf("unable to decode images file: %v", err)
		return
	}

	// write img to new image RGBA
	draw.Draw(atlasImage.RGBA, image.Rect(0, 0, 1975, 175), imageData, image.ZP, draw.Src)
	atlasImage.Upload()

	controlImage = images.NewImage(1375, 175)
	control, err := asset.Open(ControlImages)
	if err != nil {
		log.Printf("error opening control images: %v", err)
		return
	}
	defer control.Close()

	imageData, err = png.Decode(control)
	if err != nil {
		log.Printf("unable to decode control image file: %v", err)
		return
	}
	draw.Draw(controlImage.RGBA, image.Rect(0, 0, 1375, 175), imageData, image.ZP, draw.Src)
	controlImage.Upload()

	checkImage = images.NewImage(175, 175)
	check, err := asset.Open(CheckImage)
	if err != nil {
		log.Printf("error opening control images: %v", err)
		return
	}
	defer check.Close()
	imageData, err = png.Decode(check)
	if err != nil {
		log.Printf("unable to decode check image file: %v", err)
		return
	}
	draw.Draw(checkImage.RGBA, image.Rect(0, 0, 175, 175), imageData, image.ZP, draw.Src)
	checkImage.Upload()

	position = glctx.GetAttribLocation(program, "position")
	offset = glctx.GetUniformLocation(program, "offset")

	fontBytes := font.Default()
	ttfFont, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Printf("no font: %v\n", err)
		return
	}
	context = freetype.NewContext()
	context.SetFont(ttfFont)
	context.SetFontSize(float64(30))
	context.SetSrc(image.White)
	context.SetDPI(float64(dpi))
	context.SetHinting(imgfont.HintingNone)

	/*
		winImage = images.NewImage(200, 100)

		context.SetDst(winImage.RGBA)
		context.SetClip(winImage.RGBA.Bounds())

		textWidth = float32(200.0)
		p, err := context.DrawString("You win!", fixed.P(0, 100))
		if err == nil {
			textWidth = float32(int32(p.X)>>6) + float32(int32(p.X)&0x3f)/1000000
		}
		winImage.Upload()
	*/

	/*
		audio, err := asset.Open(ClickWav)
		if err != nil {
			noSound = true
			return
		}
		defer audio.Close()

		size, err := audio.Seek(0, io.SeekEnd)
		if err != nil {
			noSound = true
			return
		}
		_, err = audio.Seek(0, io.SeekStart)
		if err != nil {
			noSound = true
			return
		}

		reader, err := wav.NewReader(audio, size)
		if err != nil {
			noSound = true
			return
		}

		buffers := al.GenBuffers(1)
		buffer := buffers[0]

		count := reader.GetSampleCount()

		var i uint32

		fileData := reader.GetFile()
		var bufferData []byte
		for i = 0; i < count; i++ {
			data, err := reader.ReadRawSample()
			if err != nil {
				return
			}
			bufferData = append(bufferData, data...)
		}
		buffer.BufferData(uint32(fileData.AudioFormat), bufferData, int32(fileData.SampleRate))

		sources := al.GenSources(1)
		source = sources[0]

		source.QueueBuffers(buffer)
	*/
}

func onStop(glctx gl.Context) {
	atlasImage.Release()
	controlImage.Release()
	checkImage.Release()
	// winImage.Release()
	if countImage != nil {
		countImage.Release()
	}
	images.Release()
	glctx.DeleteProgram(program)
}

func onPaint(glctx gl.Context) {
	glctx.ClearColor(0, 0, 0, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)

	glctx.UseProgram(program)

	size := grid.Size
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			color := grid.GetDotColor(x, y)
			drawDot(glctx, x, y, color)
		}
	}
	for i := 0; i < 5; i++ {
		drawControl(glctx, i)
	}
	if win {
		gridSize := float32(grid.Size) * dotSize
		checkImage.Draw(sz,
			geom.Point{geom.Pt(xOffset + gridSize/5.0), geom.Pt(heightOffsetPt + gridSize/5.0)},
			geom.Point{geom.Pt(xOffset + gridSize*4.0/5.0), geom.Pt(heightOffsetPt + gridSize/5.0)},
			geom.Point{geom.Pt(xOffset + gridSize/5.0), geom.Pt(heightOffsetPt + gridSize*4.0/5.0)},
			checkImage.RGBA.Bounds())
		/*
			start := float32(sz.WidthPt)/2.0 - textWidth/2.0
			winImage.Draw(sz, geom.Point{geom.Pt(start), sz.HeightPt - 100},
				geom.Point{geom.Pt(start + 200), sz.HeightPt - 100},
				geom.Point{geom.Pt(start), sz.HeightPt},
				winImage.RGBA.Bounds())
		*/
	}
}

var (
	landscape       bool
	dotSize         float32
	controlSize     float32
	controlLocation [5]geom.Point
)

func gridClick(x, y float32) (int, int, bool) {
	startX := x/pxPerPt - xOffset
	startY := y/pxPerPt - heightOffsetPt
	if startX < 0.0 || startY < 0.0 {
		return 0, 0, false
	}
	gridX := int(startX / dotSize)
	gridY := int(startY / dotSize)
	if gridX < grid.Size && gridY < grid.Size {
		return gridX, gridY, true
	}
	return 0, 0, false
}

func controlClick(x, y float32) bool {
	for i := 0; i < 5; i++ {
		startX := float32(controlLocation[i].X) * pxPerPt
		startY := (float32(controlLocation[i].Y) + heightOffsetPt) * pxPerPt
		endX := startX + controlSize*pxPerPt
		endY := startY + controlSize*pxPerPt
		if x > startX && x < endX && y > startY && y < endY {
			switch i {
			case ControlIndexRestart:
				grid.Reset()
				return true
			case ControlIndexSound:
				muted = !muted
				if muted {
					app.SetKeyValue("muted", "true")
				} else {
					app.SetKeyValue("muted", "false")
				}
				return true
			case ControlIndexCount:
				// nothing to do
				return false
			case ControlIndexSize:
				if GridSize == 3 {
					GridSize = 4
					currentIndex = index4
				} else {
					GridSize = 3
					currentIndex = index3
				}
				onResize()
				app.SetKeyValue("gridSize", fmt.Sprintf("%d", GridSize))
				grid, _ = game.Random(GridSize, currentIndex)
				return true
			case ControlIndexSkip:
				if GridSize == 3 {
					index3++
					app.SetKeyValue("index3", fmt.Sprintf("%d", index3))
					currentIndex = index3
				} else {
					index4++
					app.SetKeyValue("index4", fmt.Sprintf("%d", index4))
					currentIndex = index4
				}
				grid, _ = game.Random(GridSize, currentIndex)
				return true
			}
		}
	}
	return false
}

func onResize() {
	if sz.WidthPx > sz.HeightPx-36 {
		landscape = true
	} else {
		landscape = false
	}
	if landscape {
		remainingWidth := float32(sz.WidthPt) * 0.8
		if remainingWidth < float32(sz.HeightPt)-heightOffsetPt {
			dotSize = remainingWidth / float32(GridSize)
		} else {
			dotSize = (float32(sz.HeightPt) - heightOffsetPt) / float32(GridSize)
		}
		controlSize = (float32(sz.HeightPt) - heightOffsetPt) / 5.0
		for i := 0; i < 5; i++ {
			controlLocation[i] = geom.Point{
				geom.Pt(float32(sz.WidthPt) - (float32(sz.HeightPt)-heightOffsetPt)*0.2),
				geom.Pt(float32(i) * controlSize), // relative to height offset for notification bar
			}
		}
	} else {
		remainingHeight := float32(sz.HeightPt) - heightOffsetPt - float32(sz.WidthPt)*0.2
		if remainingHeight < float32(sz.WidthPt) {
			dotSize = remainingHeight / float32(GridSize)
		} else {
			dotSize = float32(sz.WidthPt) / float32(GridSize)
		}
		controlSize = float32(sz.WidthPt) / 5.0
		for i := 0; i < 5; i++ {
			controlLocation[i] = geom.Point{
				geom.Pt(float32(i) * controlSize),
				geom.Pt(dotSize * float32(GridSize)), // relative to height offset for notification bar
			}
		}
	}
}

func drawDot(glctx gl.Context, x, y int, color string) {
	idx := colorIndex[color]
	// determine start location on screen
	xOffset = float32(0.0)
	if landscape {
		if dotSize*float32(grid.Size) < float32(sz.WidthPt)-(float32(sz.HeightPt)-heightOffsetPt)*0.2 {
			xOffset = (float32(sz.WidthPt) - (float32(sz.HeightPt)-heightOffsetPt)*0.2 - (dotSize * float32(grid.Size))) / 2.0
		}
	} else {
		if dotSize*float32(grid.Size) < float32(sz.WidthPt) {
			xOffset = (float32(sz.WidthPt) - (dotSize * float32(grid.Size))) / 2.0
		}
	}
	width := dotSize
	topLeft := geom.Point{geom.Pt(width*float32(x) + xOffset), geom.Pt(width*float32(y) + heightOffsetPt)}
	topRight := geom.Point{geom.Pt(width*float32(x+1) + xOffset), geom.Pt(width*float32(y) + heightOffsetPt)}
	bottomLeft := geom.Point{geom.Pt(width*float32(x) + xOffset), geom.Pt(width*float32(y+1) + heightOffsetPt)}
	// determine offset to correct sub-image
	srcBounds := image.Rect(200*idx, 0, 200*idx+175, 175)
	// draw sub-image
	atlasImage.Draw(sz, topLeft, topRight, bottomLeft, srcBounds)
}

func drawControl(glctx gl.Context, idx int) {
	location := controlLocation[idx]
	switch idx {
	case ControlIndexRestart:
	case ControlIndexSound:
		if muted {
			idx = MutedIndex
		} else {
			idx = UnmutedIndex
		}
	case ControlIndexCount:
	case ControlIndexSize:
		if GridSize == 3 {
			idx = Grid3Index
		} else {
			idx = Grid4Index
		}
	case ControlIndexSkip:
	default:
		return
	}
	// draw image index at location
	topLeft := geom.Point{geom.Pt(location.X), geom.Pt(float32(location.Y) + heightOffsetPt)}
	topRight := geom.Point{geom.Pt(float32(location.X) + controlSize), geom.Pt(float32(location.Y) + heightOffsetPt)}
	bottomLeft := geom.Point{geom.Pt(location.X), geom.Pt(float32(location.Y) + controlSize + heightOffsetPt)}
	// determine offset to correct sub-image
	srcBounds := image.Rect(200*idx, 0, 200*idx+175, 175)
	controlImage.Draw(sz, topLeft, topRight, bottomLeft, srcBounds)

	if idx == ControlIndexCount {
		// display text
		str := fmt.Sprintf("%d", currentIndex)
		if countImage != nil {
			countImage.Release()
		}
		countImage = images.NewImage(100, 25)
		context.SetDst(countImage.RGBA)
		context.SetClip(countImage.RGBA.Bounds())
		context.SetFontSize(24.0)
		context.SetSrc(image.White)
		context.SetDPI(float64(72.0))
		context.SetHinting(imgfont.HintingNone)
		p, err := context.DrawString(str, fixed.P(0, 24))
		if err != nil {
			countImage.Release()
			return
		}

		countImage.Upload()

		textWidth := float32(int32(p.X)>>6) + float32(int32(p.X)&0x3f)/1000000
		textHeight := controlSize / 5.0
		textWidth = textWidth * (controlSize / 5.0) / 24.0
		// actual height: 1/5th controlSize
		// width is proportional: (4x)
		xStart := float32(location.X) + controlSize/2.0 - (textWidth)/2.0
		yStart := float32(location.Y) + controlSize/2.0 - (textHeight / 2.0)
		xEnd := xStart + 4.0*(controlSize/5.0)
		yEnd := yStart + textHeight

		countImage.Draw(
			sz,
			geom.Point{geom.Pt(xStart), geom.Pt(yStart + heightOffsetPt)},
			geom.Point{geom.Pt(xEnd), geom.Pt(yStart + heightOffsetPt)},
			geom.Point{geom.Pt(xStart), geom.Pt(yEnd + heightOffsetPt)},
			countImage.RGBA.Bounds())
	}
}

const vertexShader = `#version 100
uniform vec2 offset;

attribute vec4 position;
void main() {
	// offset comes in with x/y values between 0 and 1.
	// position bounds are -1 to 1.
	vec4 offset4 = vec4(2.0*offset.x-1.0, 1.0-2.0*offset.y, 0, 0);
	gl_Position = position + offset4;
}`

const fragmentShader = `#version 100
precision mediump float;
uniform vec4 color;
void main() {
	gl_FragColor = color;
}`
