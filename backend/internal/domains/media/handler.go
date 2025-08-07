package media

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"social-app/internal/models"
	"social-app/pkg/middleware"
)

type Handler struct {
	usecase UseCase
}

func NewHandler(uc UseCase) Handler {
	return Handler{
		usecase: uc,
	}
}

// UploadImage godoc
// @Summary Upload image
// @Description Upload an image file
// @Tags Media
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file to upload"
// @Success 200 {object} models.Media
// @Failure 400 {object} map[string]string "No file is received"
// @Failure 500 {object} map[string]string "Failed to upload image"
// @Router /media/upload [post]
// @Security BearerAuth
func (h Handler) UploadImage(c *middleware.Context) {
	mf, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}
	formFiles := mf.File["images"]

	if len(formFiles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	ms := make([]models.Media, len(formFiles))
	wg := sync.WaitGroup{}
	for i := range formFiles {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()
			m, err := h.usecase.UploadImage(c.Request.Context(), file, c.User.ID)
			if err != nil {
				fmt.Println("Failed to upload image:", err)
				return
			}
			ms[i] = m
		}(formFiles[i])
	}
	wg.Wait()

	c.JSON(http.StatusOK, gin.H{
		"medias": ms,
	})
}
