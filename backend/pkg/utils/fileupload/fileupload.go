package fileupload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

const (
	UploadFolder = "./uploads"
	MaxFileSize  = 10 << 20 // 10 MB
)

func init() {
	// Crear carpeta de uploads si no existe
	if err := os.MkdirAll(UploadFolder, os.ModePerm); err != nil {
		fmt.Printf("Error creating uploads folder: %v\n", err)
	}
}

func SaveFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Generar nombre único
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, header.Filename)

	// Crear archivo destino
	dstPath := filepath.Join(UploadFolder, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %w", err)
	}
	defer dst.Close()

	// Copiar contenido
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("error copying file: %w", err)
	}

	return filename, nil
}

func GetFilePath(filename string) string {
	return filepath.Join(UploadFolder, filename)
}
