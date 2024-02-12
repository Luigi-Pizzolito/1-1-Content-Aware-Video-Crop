package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/png"
	"fmt"
)

const (
	steps = 5
)

var (
	screenHeight    int
    screenWidth     int

	ppan            int
    pzoom           float64
    subj            int
    tframes         int
    cframes         int

	termUI			bool
)

func setupUI() {
	// Setup dimensions for player-only and algorithm show modes
	if !playOnlyMode {
		screenHeight = screenHeight*steps
	} else {
		screenWidth = screenHeight
		ebiten.SetWindowFloating(true)
		// ebiten.SetWindowDecorated(false)
		// ebiten.SetFullscreen(true)
		// ebiten.SetWindowMousePassthrough(true)
	}
	termUI = false
	fmt.Println("w_size", screenWidth,screenHeight)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("1:1 Content-Aware Video Crop")
	
}

type Game struct {
	
}

func (g *Game) Update() error {
	// called every tick
    if termUI {
        return ebiten.Termination
    }
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !playOnlyMode {
		inputFrameRW.Lock()
		screen.DrawImage(ebiten.NewImageFromImage(inputFrame), &ebiten.DrawImageOptions{})
		inputFrameRW.Unlock()

		op1 := &ebiten.DrawImageOptions{}
		op1.GeoM.Translate(0, float64(squareSize))
		grayFrameRW.Lock()
		screen.DrawImage(ebiten.NewImageFromImage(grayFrame), op1)
		grayFrameRW.Unlock()

		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Translate(0, float64(squareSize*2))
		binaryFrameRW.Lock()
		screen.DrawImage(ebiten.NewImageFromImage(binaryFrame), op2)
		binaryFrameRW.Unlock()

		op3 := &ebiten.DrawImageOptions{}
		op3.GeoM.Translate(0, float64(squareSize*3))
		calcFrameRW.Lock()
		screen.DrawImage(ebiten.NewImageFromImage(calcFrame), op3)
		calcFrameRW.Unlock()

		op4 := &ebiten.DrawImageOptions{}
		op4.GeoM.Translate(float64((screenWidth-squareSize))/2, float64(squareSize*4))
		croppedFrameRW.Lock()
		screen.DrawImage(ebiten.NewImageFromImage(croppedFrame), op4)
		croppedFrameRW.Unlock()

		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f\nPan: %dpx\nZoom: %.2f%%\nSubj: %d\nProg: %.2f%%", ebiten.ActualFPS(), ppan, pzoom, subj, float64(float64(cframes)/float64(tframes))*100))
	} else {
		// Player-only mode
		croppedFrameRW.Lock()
		screen.DrawImage(ebiten.NewImageFromImage(croppedFrame), nil)
		croppedFrameRW.Unlock()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func runUI() {
	if err := ebiten.RunGameWithOptions(&Game{}, &ebiten.RunGameOptions{ScreenTransparent: true}); err != nil {
		fmt.Println(err.Error())
		selfExit()
		return
	}
}