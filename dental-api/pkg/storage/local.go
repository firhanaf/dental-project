package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var allowedMIME = map[string]string{
	"application/pdf": ".pdf",
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"image/webp":      ".webp",
}

type Local struct {
	baseDir   string
	maxSizeMB int64
}

func NewLocal(baseDir string, maxSizeMB int64) *Local {
	return &Local{baseDir: baseDir, maxSizeMB: maxSizeMB}
}

// Save menyimpan file ke disk. Returns storedName (uuid.ext), filePath (relatif), mimeType, error.
func (s *Local) Save(file multipart.File, header *multipart.FileHeader, patientID, visitID string) (storedName, filePath, mimeType string, err error) {
	if header.Size > s.maxSizeMB*1024*1024 {
		return "", "", "", fmt.Errorf("ukuran file melebihi batas %d MB", s.maxSizeMB)
	}

	// Detect MIME dari konten aktual, bukan ekstensi
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", "", "", fmt.Errorf("baca file: %w", err)
	}
	detectedMime := http.DetectContentType(buf[:n])
	// Strip parameter (e.g. "image/jpeg; charset=...")
	detectedMime = strings.Split(detectedMime, ";")[0]
	detectedMime = strings.TrimSpace(detectedMime)

	ext, ok := allowedMIME[detectedMime]
	if !ok {
		return "", "", "", fmt.Errorf("tipe file tidak diizinkan: %s", detectedMime)
	}

	// Seek kembali ke awal
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", "", "", fmt.Errorf("seek file: %w", err)
	}

	storedName = uuid.New().String() + ext
	relPath := filepath.Join(patientID, visitID, storedName)
	absDir := filepath.Join(s.baseDir, patientID, visitID)

	if err := os.MkdirAll(absDir, 0755); err != nil {
		return "", "", "", fmt.Errorf("buat direktori: %w", err)
	}

	dst, err := os.Create(filepath.Join(s.baseDir, relPath))
	if err != nil {
		return "", "", "", fmt.Errorf("buat file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", "", "", fmt.Errorf("tulis file: %w", err)
	}

	return storedName, relPath, detectedMime, nil
}

// Delete menghapus file dari disk.
func (s *Local) Delete(filePath string) error {
	absPath := filepath.Join(s.baseDir, filePath)
	if err := os.Remove(absPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("hapus file: %w", err)
	}
	return nil
}

// ServeFile stream file ke response dengan header yang sesuai.
func (s *Local) ServeFile(w http.ResponseWriter, filePath, mimeType, originalName string) error {
	absPath := filepath.Join(s.baseDir, filePath)
	f, err := os.Open(absPath)
	if err != nil {
		return fmt.Errorf("buka file: %w", err)
	}
	defer f.Close()

	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, originalName))
	_, err = io.Copy(w, f)
	return err
}
