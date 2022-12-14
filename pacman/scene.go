package pacman

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten"
	pacimages "github.com/kgosse/pacmanresources/images"
)

type scene struct {
	matrix      [][]elem
	wallSurface *ebiten.Image
	images      map[elem]*ebiten.Image
	stage       *stage
}

func newScene(st *stage) *scene {
	s := &scene{}
	s.stage = st
	if s.stage == nil {
		s.stage = defaultStage
	}
	s.images = make(map[elem]*ebiten.Image)
	s.loadImages()
	s.createStage()
	s.buildWallSurface()
	return s
}

func (s *scene) createStage() {
	h := len(s.stage.matrix)
	w := len(s.stage.matrix[0])
	s.matrix = make([][]elem, h)
	for i := 0; i < h; i++ {
		s.matrix[i] = make([]elem, w)
		for j := 0; j < w; j++ {
			c := s.stage.matrix[i][j] - '0'
			if c <= 9 {
				s.matrix[i][j] = elem(c)
			} else {
				s.matrix[i][j] = elem(s.stage.matrix[i][j] - 'a' + 10)
			}
		}

	}
}

func (s *scene) screenWidth() int {
	w := len(s.stage.matrix[0])
	return w * stageBlocSize
}

func (s *scene) screenHeight() int {
	h := len(s.stage.matrix)
	return h * stageBlocSize
}

func (s *scene) buildWallSurface() {
	h := len(s.stage.matrix)
	w := len(s.stage.matrix[0])

	sizeW := ((w*stageBlocSize)/backgroundImageSize + 1) * backgroundImageSize
	sizeH := ((w*stageBlocSize)/backgroundImageSize + 1) * backgroundImageSize
	s.wallSurface, _ = ebiten.NewImage(sizeW, sizeH, ebiten.FilterDefault)

	for i := 0; i < sizeH/backgroundImageSize; i++ {
		y := float64(i * backgroundImageSize)
		for j := 0; j < sizeW/backgroundImageSize; j++ {
			op := &ebiten.DrawImageOptions{}
			x := float64(j * backgroundImageSize)
			op.GeoM.Translate(x, y)
			s.wallSurface.DrawImage(s.images[backgroundElem], op)
		}
	}

	for i := 0; i < h; i++ {
		y := float64(i * stageBlocSize)
		for j := 0; j < w; j++ {
			if !isWall(s.matrix[i][j]) {
				continue
			}
			op := &ebiten.DrawImageOptions{}
			x := float64(j * stageBlocSize)
			op.GeoM.Translate(x, y)
			s.wallSurface.DrawImage(s.matrix[i][j], op)
		}
	}
}

func (s *scene) loadImages() {
	for i := w0; i <= w24; i++ {
		img, _, err := image.Decode(bytes.NewReader(pacimages.Images[i]))
		handleError(err)
		s.images[i], err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
		handleError(err)
	}

	img, _, err := image.Decode(bytes.NewReader(pacimages.Background_png))
	handleError(err)
	s.images[backgroundElem], err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	handleError(err)
}
