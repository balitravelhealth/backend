package middleware

import (
	"github.com/gin-gonic/gin"
)

const LanguageKey = "language"

// Language extracts the language preference from query parameter or header.
// Defaults to "id" (Indonesian) if not specified.
// Accepts: lang query param, Accept-Language header
func Language() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.DefaultQuery("lang", "")

		// Default to Indonesian
		if lang == "" {
			lang = "id"
		}

		// Validate language
		if lang != "en" && lang != "id" {
			lang = "id"
		}

		c.Set(LanguageKey, lang)
		c.Next()
	}
}

// GetLanguage retrieves the language from context, defaults to "id"
func GetLanguage(c *gin.Context) string {
	lang, exists := c.Get(LanguageKey)
	if !exists {
		return "id"
	}
	return lang.(string)
}
