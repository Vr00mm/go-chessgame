package ui2d

import (
	"path/filepath"
	"strings"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func (ui *ui) GetSinglePixelTex(color sdl.Color) *sdl.Texture {
	tex, err := ui.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, 1, 1)
	if err != nil {
		panic(err)
	}
	pixels := make([]byte, 4)
	pixels[0] = color.R
	pixels[1] = color.G
	pixels[2] = color.B
	pixels[3] = color.A
	tex.Update(nil, pixels, 4)
	return tex
}

func (ui *ui) loadTextures() {

	filenames, err := filepath.Glob("imgs/*.png")
	if err != nil {
		panic(err)
	}

	ui.texturesIndex = make(map[string]*sdl.Texture, len(filenames))
	for _, filename := range filenames {

		tmp := strings.TrimSuffix(filename, filepath.Ext(filename))
		textureName := strings.Split(tmp, "\\")

		tex := ui.imgFileToTexture(filename)
		ui.texturesIndex[textureName[1]] = tex
	}

}

func (ui *ui) imgFileToTexture(filename string) *sdl.Texture {
	pngImage, err := img.Load(filename)
	if err != nil {
		panic(err)
	}

	tex, err := ui.renderer.CreateTextureFromSurface(pngImage)
	if err != nil {
		panic(err)
	}
	return tex
}
