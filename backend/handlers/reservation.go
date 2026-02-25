package handlers

import (
	"crypto/rand"
	"auth-app/models"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReservationHandler struct {
	db *gorm.DB
}

func (h *ReservationHandler) UploadSlip(c *gin.Context) {
	file, err := c.FormFile("slip")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ต้องแนบไฟล์ slip"})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".jpg"
	}

	b := make([]byte, 8)
	_, _ = rand.Read(b)
	name := fmt.Sprintf("%d-%x%s", time.Now().UnixNano(), b, ext)

	dir := filepath.Join("uploads", "slips")
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้างโฟลเดอร์อัปโหลด", "detail": err.Error()})
		return
	}

	path := filepath.Join(dir, name)
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "อัปโหลดไฟล์ไม่สำเร็จ", "detail": err.Error()})
		return
	}

	slipURL := "/uploads/slips/" + name
	c.JSON(http.StatusCreated, gin.H{"slip_url": slipURL})
}

func NewReservationHandler(db *gorm.DB) *ReservationHandler {
	return &ReservationHandler{db: db}
}

func toTableResponse(t models.Table) models.TableResponse {
	return models.TableResponse{
		ID: t.ID, Number: t.Number,
		Seats: t.Seats, Price: t.Price, PosX: t.PosX, PosY: t.PosY, IsActive: t.IsActive,
	}
}

func toReservationResponse(r models.Reservation) models.ReservationResponse {
	return models.ReservationResponse{
		ID: r.ID, ReservedAt: r.ReservedAt,
		Note: r.Note, SlipURL: r.SlipURL, BookingRef: r.BookingRef, Status: r.Status, AdminNote: r.AdminNote,
		User: toUserResponse(r.User), Table: toTableResponse(r.Table),
		CreatedAt: r.CreatedAt,
	}
}

// GET /api/tables?date=2024-12-25 — โต๊ะทั้งหมด + สถานะว่าง
func (h *ReservationHandler) GetTables(c *gin.Context) {
	var tables []models.Table
	h.db.Where("is_active = ?", true).Order("number asc").Find(&tables)

	// เช็คโต๊ะที่ถูกจองในวันนั้น
	bookedIDs := map[uint]bool{}
	if dateStr := c.Query("date"); dateStr != "" {
		var reservations []models.Reservation
		// ใช้ DATE() ของ PostgreSQL ใน Bangkok timezone
		h.db.Where(
			"DATE(reserved_at AT TIME ZONE 'Asia/Bangkok') = ? AND status IN (?)",
			dateStr,
			[]string{models.ReservationStatusPending, models.ReservationStatusApproved},
		).Find(&reservations)

		for _, r := range reservations {
			bookedIDs[r.TableID] = true
		}
	}

	// ดึง floor items ด้วย
	var floorItems []models.FloorItem
	h.db.Find(&floorItems)
	floorItemsResp := make([]models.FloorItemResponse, len(floorItems))
	for i, fi := range floorItems {
		floorItemsResp[i] = models.FloorItemResponse{
			ID: fi.ID, Type: fi.Type, Label: fi.Label,
			PosX: fi.PosX, PosY: fi.PosY, Width: fi.Width, Height: fi.Height,
		}
	}

	type TableWithStatus struct {
		models.TableResponse
		IsBooked bool `json:"is_booked"`
	}
	result := make([]TableWithStatus, len(tables))
	for i, t := range tables {
		result[i] = TableWithStatus{
			TableResponse: toTableResponse(t),
			IsBooked:      bookedIDs[t.ID],
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"tables":      result,
		"floor_items": floorItemsResp,
	})
}

// POST /api/reservations — จองโต๊ะ (ไม่ต้องเลือกเวลา, set 21:00 อัตโนมัติ)
func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	var input struct {
		TableID uint   `json:"table_id" binding:"required"`
		Date    string `json:"date" binding:"required"` // "2024-12-25"
		Note    string `json:"note"`
		SlipURL string `json:"slip_url"`
		BookingRef string `json:"booking_ref"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse วัน และ set เวลา 21:00 (3 ทุ่ม)
	loc, _ := time.LoadLocation("Asia/Bangkok")
	date, err := time.ParseInLocation("2006-01-02", input.Date, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบวันที่ไม่ถูกต้อง (ใช้ 2024-12-25)"})
		return
	}
	reservedAt := time.Date(date.Year(), date.Month(), date.Day(), 21, 0, 0, 0, loc)
	if reservedAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ไม่สามารถจองย้อนหลังได้"})
		return
	}

	// เช็คโต๊ะ
	var table models.Table
	if err := h.db.First(&table, input.TableID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบโต๊ะนี้"})
		return
	}
	if !table.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "โต๊ะนี้ไม่เปิดให้บริการ"})
		return
	}

	// เช็คโต๊ะซ้ำในวันเดียวกัน
	var conflict models.Reservation
	err = h.db.Where(
		"table_id = ? AND DATE(reserved_at AT TIME ZONE 'Asia/Bangkok') = ? AND status IN (?)",
		input.TableID, input.Date,
		[]string{models.ReservationStatusPending, models.ReservationStatusApproved},
	).First(&conflict).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "โต๊ะนี้ถูกจองในวันดังกล่าวแล้ว"})
		return
	}

	userID := c.GetUint("userID")
	reservation := models.Reservation{
		UserID:     userID,
		TableID:    input.TableID,
		ReservedAt: reservedAt,
		Note:       input.Note,
		SlipURL:    input.SlipURL,
		BookingRef: input.BookingRef,
		Status:     models.ReservationStatusPending,
	}
	h.db.Create(&reservation)
	h.db.Preload("User").Preload("Table").First(&reservation, reservation.ID)
	BroadcastReservationsChanged()
	c.JSON(http.StatusCreated, toReservationResponse(reservation))
}

// GET /api/reservations — ดูการจองของตัวเอง
func (h *ReservationHandler) GetMyReservations(c *gin.Context) {
	userID := c.GetUint("userID")
	loc, _ := time.LoadLocation("Asia/Bangkok")
	today := time.Now().In(loc)
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, loc)
	
	var reservations []models.Reservation
	// แสดงเฉพาะ reservations ที่ยังไม่ผ่านวันนี้ (รวมวันนี้ด้วย)
	h.db.Preload("User").Preload("Table").
		Where("user_id = ? AND DATE(reserved_at AT TIME ZONE 'Asia/Bangkok') >= ?", userID, todayStart.Format("2006-01-02")).
		Order("reserved_at desc").Find(&reservations)

	result := make([]models.ReservationResponse, len(reservations))
	for i, r := range reservations {
		result[i] = toReservationResponse(r)
	}
	c.JSON(http.StatusOK, result)
}

// DELETE /api/reservations/:id — ยกเลิก
func (h *ReservationHandler) CancelReservation(c *gin.Context) {
	userID := c.GetUint("userID")
	var reservation models.Reservation
	if err := h.db.First(&reservation, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบการจองนี้"})
		return
	}
	if reservation.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "ไม่มีสิทธิ์ยกเลิกการจองนี้"})
		return
	}
	if reservation.Status != models.ReservationStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ยกเลิกได้เฉพาะการจองที่รอยืนยันเท่านั้น"})
		return
	}
	h.db.Model(&reservation).Update("status", models.ReservationStatusCancelled)
	BroadcastReservationsChanged()
	c.JSON(http.StatusOK, gin.H{"message": "ยกเลิกการจองสำเร็จ"})
}
