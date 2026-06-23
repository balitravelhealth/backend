package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RequireRole returns a middleware that checks whether the authenticated user
// has at least one of the given role names. Must be used after Auth().
func RequireRole(db *pgxpool.Pool, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get(UserIDKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		rows, err := db.Query(c.Request.Context(),
			`SELECT r.nama_role
			 FROM user_roles ur
			 JOIN roles r ON r.role_id = ur.role_id
			 WHERE ur.user_id = $1`,
			userID,
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		defer rows.Close()

		allowed := make(map[string]struct{}, len(roles))
		for _, r := range roles {
			allowed[r] = struct{}{}
		}

		for rows.Next() {
			var roleName string
			if err := rows.Scan(&roleName); err != nil {
				continue
			}
			if _, ok := allowed[roleName]; ok {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	}
}
