package utils

import (
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/imroc/req"
	uuid "github.com/satori/go.uuid"
)

// File interface
type File interface {
	io.Reader
	io.Seeker
}

// SaveImage save image with a specific scale, depend on ffmpeg
func SaveImage(file File, scale string, path string) (string, error) {
	buffer := make([]byte, 512)
	file.Read(buffer)
	filetype := http.DetectContentType(buffer)
	var ext string
	switch filetype {
	case "image/jpeg", "image/jpg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	default:
		return "", errors.New("illegal image type")
	}

	// uuid
	id := uuid.Must(uuid.NewV4(), nil)
	uuidStr := id.String()

	// save origin file
	originName := "o-" + uuidStr + ext
	originPath := path + originName
	targetName := uuidStr + ext
	targetPath := path + targetName

	originFile, err := os.Create(originPath)
	if err != nil {
		return "", err
	}
	file.Seek(0, 0)
	if _, err := io.Copy(originFile, file); err != nil {
		return "", err
	}
	originFile.Close()
	defer os.Remove(originPath)

	// ffmpeg process
	cmd := exec.Command(
		"ffmpeg",
		"-i", originPath,
		"-y", "-strict", "-2",
		"-vf", "scale="+scale+",setdar=1:1",
		targetPath,
	)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return targetName, nil
}

// DownloadImage download image with a specific scale, depend on ffmpeg
func DownloadImage(url string, scale string, path string) (string, error) {
	uuidStr := uuid.Must(uuid.NewV4(), nil).String()
	tmpFile := path + uuidStr

	ri, err := req.Get(url)
	if err != nil {
		return "", err
	}
	if ri.Response().StatusCode != 200 {
		return "", errors.New("image download error")
	}
	defer os.Remove(tmpFile)
	err = ri.ToFile(tmpFile)
	if err != nil {
		return "", err
	}

	file, err := os.Open(tmpFile)
	if err != nil {
		return "", err
	}

	targetName, err := SaveImage(file, scale, path)
	if err != nil {
		return "", err
	}

	return targetName, nil
}

// OptimizeImage make image limited in a specific scaling
func OptimizeImage(path string, scale string) (string, error) {
	dir := filepath.Dir(path)
	target := dir + "/o-" + filepath.Base(path)

	// target already exists, return
	if _, err := os.Stat(target); err == nil {
		return target, nil
	}

	wh := strings.Split(scale, ":")

	// ffmpeg process
	cmd := exec.Command(
		"ffmpeg",
		"-i", path,
		"-y", "-strict", "-2",
		"-vf", "scale=w="+wh[0]+":h="+wh[1]+":force_original_aspect_ratio=decrease",
		target,
	)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return target, nil
}
