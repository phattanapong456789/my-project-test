package handlers

import (
	"auth-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

// GET /api/admin/users — ดูรายชื่อ user ทั้งหมด
func (h *AdminHandler) GetUsers(c *gin.Context) {
	var users []models.User
	h.db.Order("created_at desc").Find(&users)

	result := make([]models.UserResponse, len(users))
	for i, u := range users {
		result[i] = toUserResponse(u)
	}

	c.JSON(http.StatusOK, result)
}

// PUT /api/admin/users/:id/role — เปลี่ยน role user
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	var input struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Role != models.RoleAdmin && input.Role != models.RoleUser {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role must be 'admin' or 'user'"})
		return
	}

	// ป้องกัน admin ลด role ตัวเอง
	currentUserID := c.GetUint("userID")
	var targetUser models.User
	if err := h.db.First(&targetUser, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if targetUser.ID == currentUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot change your own role"})
		return
	}

	h.db.Model(&targetUser).Update("role", input.Role)
	c.JSON(http.StatusOK, gin.H{
		"message": "Role updated successfully",
		"user":    toUserResponse(targetUser),
	})
}

// DELETE /api/admin/users/:id — ลบ user
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	currentUserID := c.GetUint("userID")

	var user models.User
	if err := h.db.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.ID == currentUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete yourself"})
		return
	}

	h.db.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// POST /api/admin/users — สร้าง user โดย admin (กำหนด role ได้)
func (h *AdminHandler) CreateUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required,min=2"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.User
	if err := h.db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	role := models.RoleUser
	if input.Role == models.RoleAdmin {
		role = models.RoleAdmin
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashed),
		Role:     role,
	}
	h.db.Create(&user)

	c.JSON(http.StatusCreated, toUserResponse(user))
}
