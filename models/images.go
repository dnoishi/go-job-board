package models

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Image is NOT stored in the DB
type Image struct {
	GalleryID uint
	Filename  string
}

func (i *Image) Path() string {
	pathUrl := url.URL{
		Path: "/" + i.RelativePath(),
	}
	return pathUrl.String()
}
func (i *Image) RelativePath() string {
	return fmt.Sprintf("images/galleries/%v/%v", i.GalleryID, i.Filename)
}

type ImageService interface {
	Create(galleryID uint, r io.ReadCloser, filename string) error
	ByGalleryID(galleryID uint) ([]Image, error)
	Delete(i *Image) error
}

func NewImageService() ImageService {

	return &imageService{}
}

type imageService struct {
}

func (is *imageService) Create(galleryID uint, r io.ReadCloser, filename string) error {
	defer r.Close()
	path, err := is.mkImagePath(galleryID)
	if err != nil {
		return err
	}

	// Create a destination file
	dst, err := os.Create(path + filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy reader data to the destination file
	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}

	return nil
}
func (is *imageService) Delete(image *Image) error {

	return os.Remove(image.RelativePath())
}

func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {
	gallerypath := is.galleryPath(galleryID)
	imagesPaths, err := filepath.Glob(gallerypath + "*")
	if err != nil {
		return nil, err
	}
	result := make([]Image, len(imagesPaths))
	for i := range imagesPaths {

		imagesPaths[i] = strings.Replace(imagesPaths[i], gallerypath, "", 1)
		result[i] = Image{
			Filename:  imagesPaths[i],
			GalleryID: galleryID,
		}
	}
	return result, nil
}

func (is imageService) galleryPath(galleryID uint) string {
	return fmt.Sprintf("images/galleries/%v/", galleryID)
}

func (is *imageService) mkImagePath(galleryID uint) (string, error) {
	galleryPath := is.galleryPath(galleryID)
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}

	return galleryPath, nil
}
