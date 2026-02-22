package utils

import (
	"auth-app/models"
	"log"
	"time"

	"gorm.io/gorm"
)

// StartCleanupScheduler เริ่มต้น scheduled task สำหรับลบ reservations ที่ผ่านมาแล้ว
// รันทุกวันเวลา 02:00 น. (หลังเวลารับโต๊ะ 21:00 น.)
// และรัน cleanup ทันทีเมื่อ server เริ่มทำงาน (สำหรับกรณีที่ server ปิดอยู่ตอนตี 2)
func StartCleanupScheduler(db *gorm.DB) {
	// รัน cleanup ทันทีเมื่อ server เริ่มทำงาน (สำหรับข้อมูลที่ผ่านไปแล้ว)
	log.Println("Running initial cleanup for old reservations...")
	cleanupOldReservations(db)
	
	// ตั้งค่า scheduled task สำหรับรันทุกวันเวลา 02:00 น.
	go func() {
		loc, _ := time.LoadLocation("Asia/Bangkok")
		
		for {
			// คำนวณเวลาถัดไปที่จะรัน (02:00 น. วันถัดไป)
			now := time.Now().In(loc)
			nextRun := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, loc)
			if nextRun.Before(now) || nextRun.Equal(now) {
				nextRun = nextRun.Add(24 * time.Hour)
			}
			
			// รอจนถึงเวลาที่กำหนด
			duration := nextRun.Sub(now)
			log.Printf("Cleanup scheduler will run next at %s (in %v)", nextRun.Format("2006-01-02 15:04:05"), duration)
			time.Sleep(duration)
			
			// รัน cleanup
			cleanupOldReservations(db)
		}
	}()
}

// cleanupOldReservations ลบ reservations ที่ผ่านมาแล้ว (วันที่จองผ่านไปแล้ว)
func cleanupOldReservations(db *gorm.DB) {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc)
	// ลบ reservations ที่วันที่จองผ่านไปแล้ว (ก่อนวันนี้)
	yesterday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).Add(-24 * time.Hour)
	
	result := db.Where("DATE(reserved_at AT TIME ZONE 'Asia/Bangkok') < ?", yesterday.Format("2006-01-02")).Delete(&models.Reservation{})
	
	if result.Error != nil {
		log.Printf("Error cleaning up old reservations: %v", result.Error)
	} else {
		log.Printf("Cleaned up %d old reservations (before %s)", result.RowsAffected, yesterday.Format("2006-01-02"))
	}
}
