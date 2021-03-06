package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type spriteRenderer struct {
	renderer      *sdl.Renderer
	tex           *sdl.Texture
	width, height float64
}

type drawParameters interface {
	getPosition() vector
	getRotation() float64
}

type spriteDrawParameters struct {
	position vector
	rotation float64
}

func (parameters *spriteDrawParameters) getPosition() vector {
	return parameters.position
}

func (parameters *spriteDrawParameters) getRotation() float64 {
	return parameters.rotation
}

func newSpriteRenderer(renderer *sdl.Renderer, filename string) *spriteRenderer {
	tex := textureFromBMP(renderer, filename)
	_, _, width, height, err := tex.Query()
	if err != nil {
		panic(fmt.Errorf("querying texture: %v", err))
	}
	return newSpriteRendererWithCustomSize(renderer, filename, float64(width), float64(height))
}

func newSpriteRendererWithCustomSize(renderer *sdl.Renderer, filename string, width float64, height float64) *spriteRenderer {
	return &spriteRenderer{
		renderer: renderer,
		tex:      textureFromBMP(renderer, filename),
		width:    width,
		height:   height,
	}
}

func (sr *spriteRenderer) onDraw(parameters drawParameters) error {
	// Converting coordinates to top left of sprite
	x := parameters.getPosition().x - sr.width/2.0
	y := parameters.getPosition().y - sr.height/2.0

	sr.renderer.CopyEx(
		sr.tex,
		&sdl.Rect{X: 0, Y: 0, W: int32(sr.width), H: int32(sr.height)},
		&sdl.Rect{X: int32(x), Y: int32(y), W: int32(sr.width), H: int32(sr.height)},
		parameters.getRotation(),
		&sdl.Point{X: int32(sr.width) / 2, Y: int32(sr.height) / 2},
		sdl.FLIP_NONE)

	return nil
}

func textureFromPNG(renderer *sdl.Renderer, filename string) *sdl.Texture {
	texture, err := img.LoadTexture(renderer, filename)
	if err != nil {
		panic(fmt.Errorf("loading %v: %v", filename, err))
	}
	return texture
}

func textureFromBMP(renderer *sdl.Renderer, filename string) *sdl.Texture {
	bitmap, err := sdl.LoadBMP(filename)
	if err != nil {
		panic(fmt.Errorf("loading %v: %v", filename, err))
	}
	defer bitmap.Free()
	tex, err := renderer.CreateTextureFromSurface(bitmap)
	if err != nil {
		panic(fmt.Errorf("creating texture from %v: %v", filename, err))
	}
	return tex
}
