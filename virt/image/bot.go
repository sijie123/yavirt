package image

import "fmt"

func (img *Image) Filepath() string {
	return img.JoinVirtPath(img.Filename())
}

func (img *Image) Filename() string {
	return fmt.Sprintf("%s.img", img.Name)
}

func (img *Image) Cache() error {
	// TODO
	return nil
}
