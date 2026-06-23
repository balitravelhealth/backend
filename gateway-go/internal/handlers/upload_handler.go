package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	maxUploadSize = 5 << 20 // 5 MB
	uploadDir     = "./uploads"
)

var allowedMIME = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

func (h *Handler) UploadImage(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)
	if err := c.Request.ParseMultipartForm(maxUploadSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file terlalu besar (maks 5 MB)"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file tidak ditemukan dalam request"})
		return
	}
	defer file.Close()

	// detect MIME from first 512 bytes
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	mime := http.DetectContentType(buf[:n])
	// seek back
	if seeker, ok := file.(interface{ Seek(int64, int) (int64, error) }); ok {
		seeker.Seek(0, 0)
	}

	ext, ok := allowedMIME[mime]
	if !ok {
		// fallback: use extension from original filename
		origExt := strings.ToLower(filepath.Ext(header.Filename))
		allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
		if !allowed[origExt] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "hanya JPEG, PNG, WebP, GIF yang diizinkan"})
			return
		}
		ext = origExt
	}

	// unique filename: timestamp_random
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixMilli(), randomHex(6), ext)
	destPath := filepath.Join(uploadDir, filename)

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat direktori upload"})
		return
	}

	dst, err := os.Create(destPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan file"})
		return
	}
	defer dst.Close()

	written := int64(0)
	chunk := make([]byte, 32*1024)
	for {
		nr, rerr := file.Read(chunk)
		if nr > 0 {
			nw, werr := dst.Write(chunk[:nr])
			written += int64(nw)
			if werr != nil || written > maxUploadSize {
				dst.Close()
				os.Remove(destPath)
				c.JSON(http.StatusBadRequest, gin.H{"error": "file terlalu besar"})
				return
			}
		}
		if rerr != nil {
			break
		}
	}

	// return public URL
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := c.Request.Host
	url := fmt.Sprintf("%s://%s/uploads/%s", scheme, host, filename)
	c.JSON(http.StatusOK, gin.H{"url": url, "filename": filename})
}

func randomHex(n int) string {
	const chars = "abcdef0123456789"
	b := make([]byte, n)
	// use current nano time as simple entropy
	seed := time.Now().UnixNano()
	for i := range b {
		b[i] = chars[(seed>>uint(i*4))&0xf]
		seed = seed*6364136223846793005 + 1442695040888963407
	}
	return string(b)
}
