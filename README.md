# Auth App — Vue + Golang + PostgreSQL

ระบบ Login/Register พร้อม JWT Authentication

## Features
- สมัครสมาชิก (Register) พร้อม validation
- เข้าสู่ระบบ (Login) ด้วย email/password
- JWT Token authentication
- Route guard (กันเข้าหน้า dashboard โดยไม่ login)
- Dashboard แสดงข้อมูล user
- Logout

## API Endpoints

| Method | Endpoint          | Auth | Description      |
|--------|-------------------|------|------------------|
| POST   | /api/auth/register | ❌  | สมัครสมาชิก      |
| POST   | /api/auth/login    | ❌  | เข้าสู่ระบบ       |
| GET    | /api/auth/me       | ✅  | ดูข้อมูลตัวเอง   |

## วิธีรัน

1. ติดตั้ง Go, Node.js, PostgreSQL (ดู README_SETUP.md)
2. ดับเบิลคลิก `setup-db.bat`
3. ดับเบิลคลิก `start-backend.bat`
4. ดับเบิลคลิก `start-frontend.bat`
5. เปิด http://localhost:5173

## Git

```bash
git remote add origin https://github.com/your-username/auth-app.git
git push -u origin master
```
