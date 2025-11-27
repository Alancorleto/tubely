package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
)

func getVideoAspectRatio(videoFilePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", videoFilePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		return "", err
	}

	const aspectRatio16_9 = 16.0 / 9.0
	const aspectRatio9_16 = 9.0 / 16.0
	const toleranceRange = 0.1

	streams, ok := result["streams"].([]interface{})
	if !ok || len(streams) == 0 {
		return "", fmt.Errorf("no streams found in ffprobe output")
	}

	for _, s := range streams {
		stream, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		// prefer explicit display_aspect_ratio if present
		if dar, ok := stream["display_aspect_ratio"].(string); ok && dar != "" {
			return dar, nil
		}
		// fallback: compute from width/height
		width, wok := stream["width"].(float64)
		height, hok := stream["height"].(float64)
		if wok && hok && height != 0 {
			videoAspectRatio := width / height
			if math.Abs(videoAspectRatio-aspectRatio16_9) < toleranceRange {
				return "16:9", nil
			} else if math.Abs(videoAspectRatio-aspectRatio9_16) < toleranceRange {
				return "9:16", nil
			} else {
				return "other", nil
			}
		}
	}
	return "", fmt.Errorf("unable to find video aspect ratio")
}
