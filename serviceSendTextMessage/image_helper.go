package main

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func GetImageFromUrl(url string) image.Image {
	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		fmt.Println("Error: ", err)
	}

	defer res.Body.Close()
	m, _, err := image.Decode(res.Body)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	return m
}

func SaveImage(i image.Image) (path string) {
	path = os.TempDir() + "/" + string(time.Now().UnixNano())

	f, err := os.Create(path)
	failOnError(err, "SaveImage")

	err = jpeg.Encode(f, i, nil)
	failOnError(err, "SaveImage")

	err = f.Close()
	failOnError(err, "SaveImage")

	return
}

func GetReaderFromImage(url string) (out io.Reader) {
	b, err := ioutil.ReadFile(url)
	failOnError(err, "GetReaderFromImage")

	out = bytes.NewReader(b)

	return
}

func GetRenderFromUrl(url string) io.Reader {
	img := GetImageFromUrl(url)
	path := SaveImage(img)
	return GetReaderFromImage(path)
}

func getThumbnail(image image.Image) ([]byte, error) {

	b := image.Bounds()
	imgWidth := b.Max.X
	imgHeight := b.Max.Y

	thumbWidth := 100
	thumbHeight := 100

	if imgWidth > imgHeight {
		thumbHeight = 56
	} else {
		thumbWidth = 56
	}

	thumb := imaging.Thumbnail(image, thumbWidth, thumbHeight, imaging.CatmullRom)
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, thumb, nil)
	failOnError(err, "Error create thumbnail")
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
