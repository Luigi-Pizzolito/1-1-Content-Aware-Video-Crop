package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
)

func setupUI() {
	screenHeight = screenHeight*steps

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("1:1 Content-Aware Video Crop")
}

type Game struct {
	
}

func (g *Game) Update() error {
	// called every tick
    if checkPipelineDone() {
        return ebiten.Termination
    }
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
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

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nPan: %dpx\nZoom: %.2f%%\nSubj: %d\nProg: %.2f%%", ebiten.ActualTPS(), ebiten.ActualFPS(), ppan, pzoom, subj, float64(float64(cframes)/float64(tframes))*100))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func runUI() {
	if err := ebiten.RunGameWithOptions(&Game{}, &ebiten.RunGameOptions{ScreenTransparent: true}); err != nil {
		panic(err)
	}
}