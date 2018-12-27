// Copyright 2016 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package graphicscommand

import (
	"github.com/hajimehoshi/ebiten/v2/internal/affine"
	"github.com/hajimehoshi/ebiten/v2/internal/graphics"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver"
)

// Image represents an image that is implemented with OpenGL.
type Image struct {
	image  graphicsdriver.Image
	width  int
	height int
}

// NewImage returns a new image.
//
// Note that the image is not initialized yet.
func NewImage(width, height int) *Image {
	i := &Image{
		width:  width,
		height: height,
	}
	c := &newImageCommand{
		result: i,
		width:  width,
		height: height,
	}
	theCommandQueue.Enqueue(c)
	return i
}

func NewScreenFramebufferImage(width, height int) *Image {
	i := &Image{
		width:  width,
		height: height,
	}
	c := &newScreenFramebufferImageCommand{
		result: i,
		width:  width,
		height: height,
	}
	theCommandQueue.Enqueue(c)
	return i
}

func (i *Image) Dispose() {
	c := &disposeCommand{
		target: i,
	}
	theCommandQueue.Enqueue(c)
}

func (i *Image) Size() (int, int) {
	// i.image can be nil before initializing.
	return i.width, i.height
}

func (i *Image) DrawImage(src *Image, vertices []float32, indices []uint16, clr *affine.ColorM, mode graphics.CompositeMode, filter graphics.Filter, address graphics.Address) {
	theCommandQueue.EnqueueDrawImageCommand(i, src, vertices, indices, clr, mode, filter, address)
}

// Pixels returns the image's pixels.
// Pixels might return nil when OpenGL error happens.
func (i *Image) Pixels() []byte {
	c := &pixelsCommand{
		result: nil,
		img:    i,
	}
	theCommandQueue.Enqueue(c)
	theCommandQueue.Flush()
	return c.result
}

func (i *Image) ReplacePixels(p []byte, x, y, width, height int) {
	pixels := make([]byte, len(p))
	copy(pixels, p)
	c := &replacePixelsCommand{
		dst:    i,
		pixels: pixels,
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}
	theCommandQueue.Enqueue(c)
}

func (i *Image) IsInvalidated() bool {
	// i.image can be nil before initializing.
	if i.image == nil {
		return false
	}
	return i.image.IsInvalidated()
}
