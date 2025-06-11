package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func mediaTypeToExt(mediaType string) string {
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}
	return "." + parts[1]
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	buffer := new(bytes.Buffer)
	cmd.Stdout = buffer
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	type videodata struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}
	type jsondata struct {
		Streams []videodata `json:"streams"`
	}
	var data jsondata
	err = json.Unmarshal(buffer.Bytes(), &data)
	if err != nil {
		return "", err
	}
	if (data.Streams[0].Width / data.Streams[0].Height) == 16/9 {
		return "landscape", nil
	} else if (data.Streams[0].Width / data.Streams[0].Height) == 9/16 {
		return "portrait", nil
	} else {
		return "other", nil
	}
}

func processVideoForFastStart(filePath string) (string, error) {
	output := "/tmp/output2.mp4"
	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", output)
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return output, nil
}
