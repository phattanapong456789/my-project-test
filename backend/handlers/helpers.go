package handlers

import "auth-app/models"

// แปลง User model → UserResponse (ตัดข้อมูลที่ไม่ควรส่งออกไป เช่น password)
func toUserResponse(u models.User) models.UserResponse {
	return models.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}
