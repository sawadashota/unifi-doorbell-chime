package assets

import "github.com/gobuffalo/packr/v2"

type Assets struct {
	assets *packr.Box
}

func New() *Assets {
	return &Assets{
		assets: packr.New("assets", "./data"),
	}
}

func (a *Assets) AppIcon() []byte {
	bin, _ := a.assets.Find("AppIcon.png")
	return bin
}
