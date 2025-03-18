package utils

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func ExtractMetadata(filePath string) (map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	metadata := map[string]interface{}{
		"width":  bounds.Max.X,
		"height": bounds.Max.Y,
	}

	return metadata, nil
}

func RenameFile(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}
