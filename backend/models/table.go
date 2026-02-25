package models

import "time"

type Table struct {
	ID       uint    `json:"id" gorm:"primarykey"`
	Number   int     `json:"number" gorm:"uniqueIndex;not null"` // เลขโต๊ะ auto
	Seats    int     `json:"seats" gorm:"default:4;not null"`    // จำนวนที่นั่ง default 4
	Price    float64 `json:"price" gorm:"default:0"`             // ราคาต่อโต๊ะ
	PosX     int     `json:"pos_x" gorm:"default:100"`          // ตำแหน่ง X บน canvas
	PosY     int     `json:"pos_y" gorm:"default:100"`          // ตำแหน่ง Y บน canvas
	IsActive bool    `json:"is_active" gorm:"default:true"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FloorItem — elements บน canvas ที่ไม่ใช่โต๊ะ (เวที, บาร์, ห้องน้ำ, ทางเข้า)
type FloorItem struct {
	ID    uint   `json:"id" gorm:"primarykey"`
	Type  string `json:"type" gorm:"not null"` // stage, bar, restroom, entrance
	Label string `json:"label"`                // ชื่อที่แสดง เช่น "เวที", "บาร์"
	PosX  int    `json:"pos_x" gorm:"default:0"`
	PosY  int    `json:"pos_y" gorm:"default:0"`
	Width int    `json:"width" gorm:"default:120"`
	Height int   `json:"height" gorm:"default:80"`
}

// Reservation status
const (
	ReservationStatusPending   = "pending"
	ReservationStatusApproved  = "approved"
	ReservationStatusRejected  = "rejected"
	ReservationStatusCancelled = "cancelled"
)

type Reservation struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	UserID     uint      `json:"user_id" gorm:"not null;index"`
	TableID    uint      `json:"table_id" gorm:"not null;index"`
	ReservedAt time.Time `json:"reserved_at" gorm:"not null"` // วันที่จอง (เวลา 21:00 เสมอ)
	Note       string    `json:"note"`
	SlipURL    string    `json:"slip_url"`
	BookingRef string    `json:"booking_ref" gorm:"index"`
	Status     string    `json:"status" gorm:"default:'pending'"`
	AdminNote  string    `json:"admin_note"`

	User  User  `json:"user" gorm:"foreignKey:UserID"`
	Table Table `json:"table" gorm:"foreignKey:TableID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TableResponse struct {
	ID       uint    `json:"id"`
	Number   int     `json:"number"`
	Seats    int     `json:"seats"`
	Price    float64 `json:"price"`
	PosX     int     `json:"pos_x"`
	PosY     int     `json:"pos_y"`
	IsActive bool    `json:"is_active"`
}

type FloorItemResponse struct {
	ID     uint   `json:"id"`
	Type   string `json:"type"`
	Label  string `json:"label"`
	PosX   int    `json:"pos_x"`
	PosY   int    `json:"pos_y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type ReservationResponse struct {
	ID         uint          `json:"id"`
	ReservedAt time.Time     `json:"reserved_at"`
	Note       string        `json:"note"`
	SlipURL    string        `json:"slip_url"`
	BookingRef string        `json:"booking_ref"`
	Status     string        `json:"status"`
	AdminNote  string        `json:"admin_note"`
	User       UserResponse  `json:"user"`
	Table      TableResponse `json:"table"`
	CreatedAt  time.Time     `json:"created_at"`
}
