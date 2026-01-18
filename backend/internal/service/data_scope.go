package service

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// WithDataScope returns a GORM scope for data permission filtering
// usage: db.Scopes(s.WithDataScope(c)).Find(&users)
func (s *userService) WithDataScope(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var userID uint
		// Try to get user_id from context (could be string or uint depending on middleware)
		if val, exists := c.Get("user_id"); exists {
			switch v := val.(type) {
			case uint:
				userID = v
			case string:
				if id, err := strconv.ParseUint(v, 10, 64); err == nil {
					userID = uint(id)
				}
			}
		}
		username := c.GetString("user_name")

		scope, err := s.GetUserMaxDataScope(c, userID)
		if err != nil {
			s.logger.Error("获取数据权限失败", "err", err)
			// Default to Self if error occurs
			return db.Where("created_by = ?", username)
		}

		switch scope {
		case "Self":
			return db.Where("created_by = ?", username)
		case "Subordinate":
			subordinates, _ := s.GetSubordinateUsernames(c, username)
			subordinates = append(subordinates, username)
			return db.Where("created_by IN ?", subordinates)
		case "All":
			return db
		default:
			// Default to Self for unknown scope
			return db.Where("created_by = ?", username)
		}
	}
}
