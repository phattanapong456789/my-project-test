package handlers

import (
	"auth-app/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TableAdminHandler struct {
	db *gorm.DB
}

func NewTableAdminHandler(db *gorm.DB) *TableAdminHandler {
	return &TableAdminHandler{db: db}
}

// GET /api/admin/tables
func (h *TableAdminHandler) GetAllTables(c *gin.Context) {
	var tables []models.Table
	h.db.Order("number asc").Find(&tables)
	result := make([]models.TableResponse, len(tables))
	for i, t := range tables {
		result[i] = toTableResponse(t)
	}
	c.JSON(http.StatusOK, result)
}

// POST /api/admin/tables — เพิ่มโต๊ะ (auto-number)
func (h *TableAdminHandler) CreateTable(c *gin.Context) {
	var input struct {
		Seats int     `json:"seats"`
		Price float64 `json:"price"`
		PosX  int     `json:"pos_x"`
		PosY  int     `json:"pos_y"`
	}
	c.ShouldBindJSON(&input)

	// Auto-number: หา max number แล้ว +1
	var maxNumber int
	row := h.db.Model(&models.Table{}).Select("COALESCE(MAX(number),0)").Row()
	row.Scan(&maxNumber)
	nextNumber := maxNumber + 1

	seats := input.Seats
	if seats <= 0 {
		seats = 4
	}

	table := models.Table{
		Number: nextNumber, Seats: seats, Price: input.Price,
		PosX: input.PosX, PosY: input.PosY, IsActive: true,
	}
	h.db.Create(&table)
	c.JSON(http.StatusCreated, toTableResponse(table))
}

// PUT /api/admin/tables/:id
func (h *TableAdminHandler) UpdateTable(c *gin.Context) {
	var table models.Table
	if err := h.db.First(&table, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบโต๊ะ"})
		return
	}
	var input struct {
		Seats    *int     `json:"seats"`
		Price    *float64 `json:"price"`
		PosX     *int     `json:"pos_x"`
		PosY     *int     `json:"pos_y"`
		IsActive *bool    `json:"is_active"`
	}
	c.ShouldBindJSON(&input)
	updates := map[string]interface{}{}
	if input.Seats != nil {
		updates["seats"] = *input.Seats
	}
	if input.Price != nil {
		updates["price"] = *input.Price
	}
	if input.PosX != nil {
		updates["pos_x"] = *input.PosX
	}
	if input.PosY != nil {
		updates["pos_y"] = *input.PosY
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}
	h.db.Model(&table).Updates(updates)
	h.db.First(&table, table.ID)
	c.JSON(http.StatusOK, toTableResponse(table))
}

// DELETE /api/admin/tables/:id
func (h *TableAdminHandler) DeleteTable(c *gin.Context) {
	var table models.Table
	if err := h.db.First(&table, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบโต๊ะ"})
		return
	}
	h.db.Delete(&table)
	c.JSON(http.StatusOK, gin.H{"message": "ลบโต๊ะสำเร็จ"})
}

// ---- Floor Items ----

// GET /api/admin/floor-items
func (h *TableAdminHandler) GetFloorItems(c *gin.Context) {
	var items []models.FloorItem
	h.db.Find(&items)
	result := make([]models.FloorItemResponse, len(items))
	for i, fi := range items {
		result[i] = models.FloorItemResponse{
			ID: fi.ID, Type: fi.Type, Label: fi.Label,
			PosX: fi.PosX, PosY: fi.PosY, Width: fi.Width, Height: fi.Height,
		}
	}
	c.JSON(http.StatusOK, result)
}

// POST /api/admin/floor-items
func (h *TableAdminHandler) CreateFloorItem(c *gin.Context) {
	var input struct {
		Type   string `json:"type" binding:"required"`
		Label  string `json:"label"`
		PosX   int    `json:"pos_x"`
		PosY   int    `json:"pos_y"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Width == 0 {
		input.Width = 120
	}
	if input.Height == 0 {
		input.Height = 80
	}
	label := input.Label
	if label == "" {
		labels := map[string]string{
			"stage": "🎤 เวที",
		}
		if l, ok := labels[input.Type]; ok {
			label = l
		}
	}
	item := models.FloorItem{
		Type: input.Type, Label: label,
		PosX: input.PosX, PosY: input.PosY,
		Width: input.Width, Height: input.Height,
	}
	h.db.Create(&item)
	c.JSON(http.StatusCreated, models.FloorItemResponse{
		ID: item.ID, Type: item.Type, Label: item.Label,
		PosX: item.PosX, PosY: item.PosY, Width: item.Width, Height: item.Height,
	})
}

// PUT /api/admin/floor-items/:id
func (h *TableAdminHandler) UpdateFloorItem(c *gin.Context) {
	var item models.FloorItem
	if err := h.db.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบ item"})
		return
	}
	var input struct {
		Label  *string `json:"label"`
		PosX   *int    `json:"pos_x"`
		PosY   *int    `json:"pos_y"`
		Width  *int    `json:"width"`
		Height *int    `json:"height"`
	}
	c.ShouldBindJSON(&input)
	updates := map[string]interface{}{}
	if input.Label != nil {
		updates["label"] = *input.Label
	}
	if input.PosX != nil {
		updates["pos_x"] = *input.PosX
	}
	if input.PosY != nil {
		updates["pos_y"] = *input.PosY
	}
	if input.Width != nil {
		updates["width"] = *input.Width
	}
	if input.Height != nil {
		updates["height"] = *input.Height
	}
	h.db.Model(&item).Updates(updates)
	h.db.First(&item, item.ID)
	c.JSON(http.StatusOK, models.FloorItemResponse{
		ID: item.ID, Type: item.Type, Label: item.Label,
		PosX: item.PosX, PosY: item.PosY, Width: item.Width, Height: item.Height,
	})
}

// DELETE /api/admin/floor-items/:id
func (h *TableAdminHandler) DeleteFloorItem(c *gin.Context) {
	var item models.FloorItem
	if err := h.db.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบ item"})
		return
	}
	h.db.Delete(&item)
	c.JSON(http.StatusOK, gin.H{"message": "ลบสำเร็จ"})
}

// ---- Reservations ----

func (h *TableAdminHandler) GetAllReservations(c *gin.Context) {
	query := h.db.Preload("User").Preload("Table").Order("reserved_at asc")
	if date := c.Query("date"); date != "" {
		query = query.Where("DATE(reserved_at) = ?", date)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	var reservations []models.Reservation
	query.Find(&reservations)

	result := make([]models.ReservationResponse, len(reservations))
	for i, r := range reservations {
		result[i] = toReservationResponse(r)
	}
	c.JSON(http.StatusOK, result)
}

func (h *TableAdminHandler) UpdateReservationStatus(c *gin.Context) {
	var input struct {
		Status    string `json:"status" binding:"required"`
		AdminNote string `json:"admin_note"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	valid := map[string]bool{
		models.ReservationStatusApproved: true,
		models.ReservationStatusRejected: true,
	}
	if !valid[input.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status ต้องเป็น approved หรือ rejected"})
		return
	}
	var reservation models.Reservation
	if err := h.db.Preload("User").Preload("Table").First(&reservation, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบการจองนี้"})
		return
	}
	if reservation.Status != models.ReservationStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "เปลี่ยนสถานะได้เฉพาะ pending เท่านั้น"})
		return
	}
	h.db.Model(&reservation).Updates(map[string]interface{}{
		"status": input.Status, "admin_note": input.AdminNote,
	})
	h.db.Preload("User").Preload("Table").First(&reservation, reservation.ID)
	c.JSON(http.StatusOK, toReservationResponse(reservation))
}

func (h *TableAdminHandler) GetSummary(c *gin.Context) {
	today := "NOW()::date"
	var totalToday, pendingToday, approvedToday int64
	h.db.Model(&models.Reservation{}).Where("DATE(reserved_at) = " + today).Count(&totalToday)
	h.db.Model(&models.Reservation{}).Where("DATE(reserved_at) = "+today+" AND status = ?", models.ReservationStatusPending).Count(&pendingToday)
	h.db.Model(&models.Reservation{}).Where("DATE(reserved_at) = "+today+" AND status = ?", models.ReservationStatusApproved).Count(&approvedToday)
	var activeTables int64
	h.db.Model(&models.Table{}).Where("is_active = ?", true).Count(&activeTables)
	c.JSON(http.StatusOK, gin.H{
		"today":  gin.H{"total": totalToday, "pending": pendingToday, "approved": approvedToday},
		"tables": gin.H{"active": activeTables},
	})
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
