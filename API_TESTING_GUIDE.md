# University ERP Backend - API Testing Guide

Base URL: `http://localhost:8080`

All protected endpoints require header: `Authorization: Bearer <token>`

---

## 1. Authentication

### Register
```bash
curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin1",
    "email": "admin1@university.edu",
    "password": "secret123",
    "role_name": "admin"
  }'
```
Response: `{ "user_id": 1, "username": "admin1", "email": "admin1@university.edu" }`

### Login
```bash
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin1",
    "password": "secret123"
  }'
```
Response: `{ "token": "eyJ...", "user_id": 1, "username": "admin1", "roles": ["admin"] }`

Save the token for all subsequent requests:
```bash
TOKEN="eyJ..."
```

### Get Profile
```bash
curl -s http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer $TOKEN"
```

---

## 2. Core (Lookups & Setup)

### Public Lookups (No Auth)
```bash
curl -s http://localhost:8080/api/v1/lookups/genders
curl -s http://localhost:8080/api/v1/lookups/categories
curl -s http://localhost:8080/api/v1/lookups/blood-groups
curl -s http://localhost:8080/api/v1/lookups/status-codes
```

### University
```bash
# Create
curl -s -X POST http://localhost:8080/api/v1/universities \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "State University",
    "short_name": "SU",
    "established_year": 1965,
    "address": "University Road",
    "city": "Mumbai",
    "state": "Maharashtra",
    "postal_code": "400001",
    "phone": "022-12345678",
    "email": "info@stateuni.edu",
    "website": "https://stateuni.edu",
    "is_active": true
  }'

# List
curl -s http://localhost:8080/api/v1/universities -H "Authorization: Bearer $TOKEN"

# Get by ID
curl -s http://localhost:8080/api/v1/universities/1 -H "Authorization: Bearer $TOKEN"

# Update
curl -s -X PUT http://localhost:8080/api/v1/universities/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "name": "State University Updated", "short_name": "SU", "is_active": true }'
```

### Campus
```bash
curl -s -X POST http://localhost:8080/api/v1/campuses \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "university_id": 1,
    "name": "Main Campus",
    "code": "MC",
    "address": "University Road",
    "city": "Mumbai",
    "state": "Maharashtra",
    "postal_code": "400001",
    "phone": "022-12345678",
    "is_main_campus": true,
    "is_active": true
  }'

curl -s http://localhost:8080/api/v1/campuses -H "Authorization: Bearer $TOKEN"
curl -s http://localhost:8080/api/v1/campuses/1 -H "Authorization: Bearer $TOKEN"
```

### Department
```bash
curl -s -X POST http://localhost:8080/api/v1/departments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "campus_id": 1,
    "name": "Computer Science",
    "code": "CSE",
    "established_year": 2000,
    "description": "Department of Computer Science and Engineering",
    "is_active": true
  }'

curl -s http://localhost:8080/api/v1/departments -H "Authorization: Bearer $TOKEN"
```

### Room (Facility)
```bash
curl -s -X POST http://localhost:8080/api/v1/rooms \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "campus_id": 1,
    "room_number": "LH-101",
    "room_type": "lecture_hall",
    "capacity": 120,
    "building": "Academic Block A",
    "floor": 1,
    "is_active": true
  }'

curl -s http://localhost:8080/api/v1/rooms -H "Authorization: Bearer $TOKEN"
```

---

## 3. Admissions

### Create Admission Cycle
```bash
curl -s -X POST http://localhost:8080/api/v1/admissions/cycles \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Fall 2026 Admissions",
    "academic_year": "2026-2027",
    "start_date": "2026-01-01T00:00:00Z",
    "end_date": "2026-06-30T00:00:00Z",
    "status_id": 1
  }'
```

### List Open Cycles (Public - No Auth)
```bash
curl -s http://localhost:8080/api/v1/admissions/cycles/open
```

### Submit Application
```bash
curl -s -X POST http://localhost:8080/api/v1/admissions/applicants \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "cycle_id": 1,
    "program_id": 1,
    "first_name": "Rahul",
    "last_name": "Sharma",
    "email": "rahul@example.com",
    "phone": "9876543210",
    "date_of_birth": "2004-05-15T00:00:00Z",
    "gender_id": 1,
    "category_id": 1,
    "address": "123 Main St, Mumbai",
    "previous_qualification": "12th Science",
    "previous_percentage": 92.5
  }'
```

### Update Application Status (Approve - triggers auto student creation)
```bash
curl -s -X PUT http://localhost:8080/api/v1/admissions/applicants/1/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "status_id": 2 }'
```
**Event**: `admission.application_approved` -> Student module auto-creates student profile

### Allocate Seat
```bash
curl -s -X POST http://localhost:8080/api/v1/admissions/seat-allocations \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "applicant_id": 1,
    "program_id": 1,
    "batch_id": 1,
    "seat_type": "merit"
  }'
```

### List Applicants
```bash
curl -s "http://localhost:8080/api/v1/admissions/applicants?cycle_id=1" \
  -H "Authorization: Bearer $TOKEN"
```

### Cycle Stats
```bash
curl -s http://localhost:8080/api/v1/admissions/cycles/1/stats \
  -H "Authorization: Bearer $TOKEN"
```

### Upload Document
```bash
curl -s -X POST http://localhost:8080/api/v1/admissions/applicants/1/documents \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "document_type": "marksheet",
    "file_path": "/uploads/marksheet_rahul.pdf",
    "verified": false
  }'
```

### Verify Document
```bash
curl -s -X POST http://localhost:8080/api/v1/admissions/documents/1/verify \
  -H "Authorization: Bearer $TOKEN"
```

### Add to Waitlist
```bash
curl -s -X POST http://localhost:8080/api/v1/admissions/waitlist \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "applicant_id": 2,
    "program_id": 1,
    "position": 1
  }'
```

---

## 4. Student

### Enroll Student (Manual)
```bash
curl -s -X POST http://localhost:8080/api/v1/students/enroll \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 2,
    "program_id": 1,
    "first_name": "Priya",
    "last_name": "Patel",
    "date_of_birth": "2004-03-20",
    "email": "priya@example.com",
    "phone": "9876543211",
    "gender_id": 2,
    "category_id": 1,
    "admission_year": 2026,
    "admission_quota": "merit"
  }'
```

### List Students
```bash
curl -s http://localhost:8080/api/v1/students -H "Authorization: Bearer $TOKEN"
```

### Get Student
```bash
curl -s http://localhost:8080/api/v1/students/1 -H "Authorization: Bearer $TOKEN"
```

### My Profile (Student's own)
```bash
curl -s http://localhost:8080/api/v1/students/me -H "Authorization: Bearer $TOKEN"
```

### My Dashboard
```bash
curl -s http://localhost:8080/api/v1/students/me/dashboard -H "Authorization: Bearer $TOKEN"
```

### Add Guardian
```bash
curl -s -X POST http://localhost:8080/api/v1/students/1/guardians \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Suresh Sharma",
    "relation": "father",
    "phone": "9876543212",
    "email": "suresh@example.com",
    "occupation": "Business",
    "is_primary": true
  }'
```

### Get Guardians
```bash
curl -s http://localhost:8080/api/v1/students/1/guardians -H "Authorization: Bearer $TOKEN"
```

### File Grievance
```bash
curl -s -X POST http://localhost:8080/api/v1/students/1/grievances \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "category": "academic",
    "description": "Issue with grade calculation in CS101"
  }'
```

---

## 5. Academic

### Create Academic Term
```bash
curl -s -X POST http://localhost:8080/api/v1/academic/terms \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Fall 2026 Semester",
    "term_number": 1,
    "start_date": "2026-08-01T00:00:00Z",
    "end_date": "2026-12-15T00:00:00Z",
    "academic_year": "2026-2027",
    "is_current": false
  }'
```

### Set Current Term
```bash
curl -s -X POST http://localhost:8080/api/v1/academic/terms/1/set-current \
  -H "Authorization: Bearer $TOKEN"
```

### Get Current Term
```bash
curl -s http://localhost:8080/api/v1/academic/terms/current -H "Authorization: Bearer $TOKEN"
```

### Create Program
```bash
curl -s -X POST http://localhost:8080/api/v1/academic/programs \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "B.Tech Computer Science",
    "code": "BT-CS",
    "department_id": 1,
    "duration_years": 4,
    "total_semesters": 8,
    "degree_level": "undergraduate",
    "is_active": true
  }'
```

### Create Subject
```bash
curl -s -X POST http://localhost:8080/api/v1/academic/subjects \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Data Structures and Algorithms",
    "code": "CS201",
    "credits": 4,
    "lecture_hours": 3,
    "lab_hours": 2,
    "type": "core"
  }'
```

### Add Subject to Program Curriculum
```bash
curl -s -X POST http://localhost:8080/api/v1/academic/programs/1/subjects \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "subject_id": 1,
    "semester_number": 3,
    "is_elective": false
  }'
```

### Get Curriculum
```bash
curl -s http://localhost:8080/api/v1/academic/programs/1/curriculum \
  -H "Authorization: Bearer $TOKEN"
```

### Create Batch
```bash
curl -s -X POST http://localhost:8080/api/v1/academic/batches \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "program_id": 1,
    "batch_year": 2026,
    "name": "2026-2030 Batch",
    "current_semester": 1
  }'
```

### Create Section
```bash
curl -s -X POST http://localhost:8080/api/v1/academic/batches/1/sections \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Section A",
    "capacity": 60
  }'
```

### Create Course Offering
```bash
curl -s -X POST http://localhost:8080/api/v1/academic/offerings \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "subject_id": 1,
    "term_id": 1,
    "employee_id": 1,
    "section_id": 1,
    "max_enrollment": 60
  }'
```

### Register Student for Term (triggers finance invoice)
```bash
curl -s -X POST http://localhost:8080/api/v1/academic/term-registrations \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 1,
    "academic_term_id": 1,
    "batch_id": 1
  }'
```
**Event**: `academic.term_registered` -> Finance module auto-generates invoice

### Register Student for Course
```bash
curl -s -X POST http://localhost:8080/api/v1/academic/course-registrations \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 1,
    "offering_id": 1
  }'
```

### Drop Course
```bash
curl -s -X POST http://localhost:8080/api/v1/academic/course-registrations/1/drop \
  -H "Authorization: Bearer $TOKEN"
```

### Student Courses
```bash
curl -s http://localhost:8080/api/v1/academic/students/1/courses \
  -H "Authorization: Bearer $TOKEN"
```

### Student Timetable
```bash
curl -s http://localhost:8080/api/v1/academic/students/1/timetable \
  -H "Authorization: Bearer $TOKEN"
```

### Create Timetable Entry
```bash
curl -s -X POST http://localhost:8080/api/v1/academic/offerings/1/timetable \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "day_of_week": "monday",
    "start_time": "09:00",
    "end_time": "10:00",
    "room_id": 1
  }'
```

### Mark Attendance
```bash
curl -s -X POST http://localhost:8080/api/v1/academic/attendance \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": 1,
    "student_id": 1,
    "status_id": 1
  }'
```

### Get Student Attendance
```bash
curl -s http://localhost:8080/api/v1/academic/students/1/attendance \
  -H "Authorization: Bearer $TOKEN"
```

### Get Attendance Summary
```bash
curl -s http://localhost:8080/api/v1/academic/students/1/attendance/summary \
  -H "Authorization: Bearer $TOKEN"
```

### Academic Calendar
```bash
# List
curl -s http://localhost:8080/api/v1/academic/calendar -H "Authorization: Bearer $TOKEN"

# Create event
curl -s -X POST http://localhost:8080/api/v1/academic/calendar \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Mid-Semester Break",
    "event_type": "holiday",
    "start_date": "2026-10-10T00:00:00Z",
    "end_date": "2026-10-15T00:00:00Z"
  }'
```

---

## 6. Finance

### Finance Summary
```bash
curl -s http://localhost:8080/api/v1/finance/summary -H "Authorization: Bearer $TOKEN"
```

### Fee Heads
```bash
# Create
curl -s -X POST http://localhost:8080/api/v1/finance/fee-heads \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Tuition Fee",
    "code": "TUITION",
    "description": "Annual tuition fee",
    "is_mandatory": true
  }'

# List
curl -s http://localhost:8080/api/v1/finance/fee-heads -H "Authorization: Bearer $TOKEN"

# Update
curl -s -X PUT http://localhost:8080/api/v1/finance/fee-heads/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "name": "Tuition Fee Updated", "code": "TUITION", "is_mandatory": true }'
```

### Fee Structures
```bash
# Create
curl -s -X POST http://localhost:8080/api/v1/finance/fee-structures \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "program_id": 1,
    "semester_number": 1,
    "fee_head_id": 1,
    "amount": 75000,
    "academic_year": "2026-2027",
    "is_active": true
  }'

# List
curl -s http://localhost:8080/api/v1/finance/fee-structures -H "Authorization: Bearer $TOKEN"
```

### Generate Invoice (Manual)
```bash
curl -s -X POST http://localhost:8080/api/v1/finance/invoices/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 1,
    "term_id": 1,
    "program_id": 1
  }'
```

### List Invoices
```bash
curl -s http://localhost:8080/api/v1/finance/invoices -H "Authorization: Bearer $TOKEN"
curl -s http://localhost:8080/api/v1/finance/invoices/student/1 -H "Authorization: Bearer $TOKEN"
curl -s http://localhost:8080/api/v1/finance/invoices/me -H "Authorization: Bearer $TOKEN"
```

### Process Payment
```bash
curl -s -X POST http://localhost:8080/api/v1/finance/payments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "invoice_id": 1,
    "student_id": 1,
    "amount": 75000,
    "transaction_id": "TXN-20260520-001",
    "payment_mode_id": 1
  }'
```

### List Payments
```bash
curl -s http://localhost:8080/api/v1/finance/payments -H "Authorization: Bearer $TOKEN"
```

### Scholarships
```bash
# Create
curl -s -X POST http://localhost:8080/api/v1/finance/scholarships \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Merit Scholarship",
    "description": "For students with >90% in previous qualification",
    "eligibility_criteria": {"min_percentage": 90},
    "amount": 25000,
    "renewable": true
  }'

# Assign to student
curl -s -X POST http://localhost:8080/api/v1/finance/scholarships/assign \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 1,
    "scholarship_id": 1,
    "academic_year": "2026-2027",
    "amount_awarded": 25000
  }'

# List student scholarships
curl -s http://localhost:8080/api/v1/finance/scholarships/student/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Discounts
```bash
curl -s -X POST http://localhost:8080/api/v1/finance/discounts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 1,
    "fee_head_id": 1,
    "academic_year": "2026-2027",
    "amount": 10000,
    "reason": "Staff ward concession"
  }'

curl -s http://localhost:8080/api/v1/finance/discounts/student/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Installments
```bash
curl -s -X POST http://localhost:8080/api/v1/finance/installments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 1,
    "academic_term_id": 1,
    "due_date": "2026-09-01T00:00:00Z",
    "amount": 37500
  }'

curl -s http://localhost:8080/api/v1/finance/installments/student/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Refunds
```bash
# Request
curl -s -X POST http://localhost:8080/api/v1/finance/refunds \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "payment_id": 1,
    "student_id": 1,
    "amount": 5000,
    "reason": "Course drop - partial refund"
  }'

# Approve
curl -s -X POST http://localhost:8080/api/v1/finance/refunds/1/approve \
  -H "Authorization: Bearer $TOKEN"

# List
curl -s http://localhost:8080/api/v1/finance/refunds -H "Authorization: Bearer $TOKEN"
```

---

## 7. HR

### HR Stats
```bash
curl -s http://localhost:8080/api/v1/hr/stats -H "Authorization: Bearer $TOKEN"
```

### Designations
```bash
curl -s -X POST http://localhost:8080/api/v1/hr/designations \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "code": "PROF", "name": "Professor", "is_active": true }'

curl -s http://localhost:8080/api/v1/hr/designations -H "Authorization: Bearer $TOKEN"
```

### Employment Types
```bash
curl -s http://localhost:8080/api/v1/hr/employment-types -H "Authorization: Bearer $TOKEN"
```

### Leave Types
```bash
curl -s -X POST http://localhost:8080/api/v1/hr/leave-types \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "CL",
    "name": "Casual Leave",
    "max_days": 12,
    "paid": true,
    "is_active": true
  }'

curl -s http://localhost:8080/api/v1/hr/leave-types -H "Authorization: Bearer $TOKEN"
```

### Salary Components
```bash
curl -s -X POST http://localhost:8080/api/v1/hr/salary-components \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "code": "HRA", "name": "House Rent Allowance", "type": "earning", "is_active": true }'

curl -s http://localhost:8080/api/v1/hr/salary-components -H "Authorization: Bearer $TOKEN"
```

### Create Employee (triggers onboarding event)
```bash
curl -s -X POST http://localhost:8080/api/v1/hr/employees \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 3,
    "employee_code": "EMP-001",
    "first_name": "Dr. Anil",
    "last_name": "Kumar",
    "phone": "9876543213",
    "address": "456 Faculty Quarters",
    "joining_date": "2020-07-01T00:00:00Z",
    "employment_type_id": 1,
    "department_id": 1,
    "designation_id": 1,
    "is_active": true
  }'
```
**Event**: `hr.employee_onboarded`

### List Employees
```bash
curl -s http://localhost:8080/api/v1/hr/employees -H "Authorization: Bearer $TOKEN"
curl -s http://localhost:8080/api/v1/hr/employees/1 -H "Authorization: Bearer $TOKEN"
curl -s http://localhost:8080/api/v1/hr/employees/me -H "Authorization: Bearer $TOKEN"
```

### Faculty Profile
```bash
curl -s -X PUT http://localhost:8080/api/v1/hr/employees/1/faculty-profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "specialization": "Machine Learning",
    "qualification": "Ph.D. Computer Science",
    "research_area": "Artificial Intelligence",
    "office_hours": "Mon-Fri 2-4 PM",
    "max_load_credits": 16
  }'

curl -s http://localhost:8080/api/v1/hr/employees/1/faculty-profile \
  -H "Authorization: Bearer $TOKEN"

curl -s http://localhost:8080/api/v1/hr/faculty -H "Authorization: Bearer $TOKEN"
```

### Assign Salary
```bash
curl -s -X POST http://localhost:8080/api/v1/hr/employees/1/salary \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "effective_from": "2026-01-01T00:00:00Z",
    "base_pay": 120000,
    "net_salary": 95000,
    "is_active": true
  }'
```

### Run Payroll (triggers finance voucher)
```bash
curl -s -X POST http://localhost:8080/api/v1/hr/payroll/run \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": 1,
    "month": "2026-05"
  }'
```
**Event**: `hr.payroll_processed` -> Finance module records payroll voucher

### Leave Management
```bash
# Request leave
curl -s -X POST http://localhost:8080/api/v1/hr/leave-requests \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": 1,
    "leave_type_id": 1,
    "start_date": "2026-06-01T00:00:00Z",
    "end_date": "2026-06-03T00:00:00Z",
    "reason": "Personal work"
  }'

# List leave requests
curl -s http://localhost:8080/api/v1/hr/leave-requests -H "Authorization: Bearer $TOKEN"

# Approve
curl -s -X POST http://localhost:8080/api/v1/hr/leave-requests/1/approve \
  -H "Authorization: Bearer $TOKEN"

# Reject
curl -s -X POST http://localhost:8080/api/v1/hr/leave-requests/1/reject \
  -H "Authorization: Bearer $TOKEN"

# Leave balances
curl -s http://localhost:8080/api/v1/hr/employees/1/leave-balances \
  -H "Authorization: Bearer $TOKEN"
```

### HR Attendance
```bash
curl -s -X POST http://localhost:8080/api/v1/hr/attendance \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": 1,
    "attendance_date": "2026-05-20T00:00:00Z",
    "check_in": "2026-05-20T09:00:00Z",
    "status_id": 1
  }'

curl -s http://localhost:8080/api/v1/hr/employees/1/attendance \
  -H "Authorization: Bearer $TOKEN"
```

### Recruitment
```bash
# Post job
curl -s -X POST http://localhost:8080/api/v1/hr/jobs \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Assistant Professor - CSE",
    "department_id": 1,
    "employment_type_id": 1,
    "vacancies": 2,
    "posted_date": "2026-05-20T00:00:00Z",
    "closing_date": "2026-06-30T00:00:00Z",
    "description": "Looking for candidates with Ph.D. in CS"
  }'

# List jobs
curl -s http://localhost:8080/api/v1/hr/jobs -H "Authorization: Bearer $TOKEN"

# Apply for job
curl -s -X POST http://localhost:8080/api/v1/hr/jobs/1/applications \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "applicant_name": "Dr. Meera Joshi",
    "email": "meera@example.com",
    "phone": "9876543214",
    "resume_path": "/uploads/resume_meera.pdf"
  }'

# Update application status
curl -s -X PUT http://localhost:8080/api/v1/hr/job-applications/1/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "status_id": 2 }'
```

### Transfer Department
```bash
curl -s -X POST http://localhost:8080/api/v1/hr/employees/1/transfer \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "department_id": 2 }'
```

### Deactivate Employee
```bash
curl -s -X POST http://localhost:8080/api/v1/hr/employees/1/deactivate \
  -H "Authorization: Bearer $TOKEN"
```

---

## 8. Exam

### Exam Components
```bash
curl -s -X POST http://localhost:8080/api/v1/exam/components \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mid Semester",
    "max_marks": 30,
    "weightage_percentage": 30
  }'

curl -s http://localhost:8080/api/v1/exam/components -H "Authorization: Bearer $TOKEN"
```

### Exam Schedules
```bash
curl -s -X POST http://localhost:8080/api/v1/exam/schedules \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "term_id": 1,
    "subject_id": 1,
    "exam_type": "mid_sem",
    "start_date": "2026-09-15T09:00:00Z",
    "end_date": "2026-09-15T12:00:00Z",
    "room_id": 1
  }'

curl -s http://localhost:8080/api/v1/exam/schedules -H "Authorization: Bearer $TOKEN"
```

### Enter Result
```bash
curl -s -X POST http://localhost:8080/api/v1/exam/results \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 1,
    "subject_id": 1,
    "schedule_id": 1,
    "marks_obtained": 78,
    "total_marks": 100,
    "grade": "A",
    "grade_point": 8
  }'
```

### Bulk Enter Results
```bash
curl -s -X POST http://localhost:8080/api/v1/exam/results/bulk \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '[
    { "student_id": 1, "subject_id": 1, "schedule_id": 1, "marks_obtained": 78, "total_marks": 100, "grade": "A", "grade_point": 8 },
    { "student_id": 2, "subject_id": 1, "schedule_id": 1, "marks_obtained": 65, "total_marks": 100, "grade": "B", "grade_point": 6 }
  ]'
```

### Publish Results (triggers event)
```bash
curl -s -X POST http://localhost:8080/api/v1/exam/results/publish \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "term_id": 1 }'
```
**Event**: `exam.result_published`

### Student Results & Transcript
```bash
curl -s http://localhost:8080/api/v1/exam/students/1/results -H "Authorization: Bearer $TOKEN"
curl -s http://localhost:8080/api/v1/exam/students/1/transcript -H "Authorization: Bearer $TOKEN"
curl -s http://localhost:8080/api/v1/exam/students/1/sgpa -H "Authorization: Bearer $TOKEN"
curl -s http://localhost:8080/api/v1/exam/students/1/cgpa -H "Authorization: Bearer $TOKEN"
curl -s http://localhost:8080/api/v1/exam/students/me/results -H "Authorization: Bearer $TOKEN"
curl -s http://localhost:8080/api/v1/exam/students/me/transcript -H "Authorization: Bearer $TOKEN"
curl -s http://localhost:8080/api/v1/exam/students/me/cgpa -H "Authorization: Bearer $TOKEN"
```

### Component Marks
```bash
curl -s -X POST http://localhost:8080/api/v1/exam/results/1/component-marks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "component_id": 1,
    "marks_obtained": 25,
    "max_marks": 30
  }'

curl -s http://localhost:8080/api/v1/exam/results/1/component-marks \
  -H "Authorization: Bearer $TOKEN"
```

### Revaluation
```bash
curl -s -X POST http://localhost:8080/api/v1/exam/revaluations \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "result_id": 1,
    "student_id": 1,
    "reason": "Marks seem incorrect"
  }'

curl -s http://localhost:8080/api/v1/exam/revaluations -H "Authorization: Bearer $TOKEN"

# Process revaluation
curl -s -X POST http://localhost:8080/api/v1/exam/revaluations/1/process \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "reviewed_marks": 82,
    "grade": "A+",
    "remarks": "Marks updated after revaluation"
  }'
```

### Supplementary Exam
```bash
curl -s -X POST http://localhost:8080/api/v1/exam/supplementary \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "term_id": 1,
    "subject_id": 1,
    "exam_date": "2027-01-15T09:00:00Z"
  }'

curl -s http://localhost:8080/api/v1/exam/supplementary -H "Authorization: Bearer $TOKEN"
```

---

## 9. Library

### Stats
```bash
curl -s http://localhost:8080/api/v1/library/stats -H "Authorization: Bearer $TOKEN"
```

### Authors
```bash
curl -s -X POST http://localhost:8080/api/v1/library/authors \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "name": "Thomas Cormen", "biography": "Author of Introduction to Algorithms" }'

curl -s http://localhost:8080/api/v1/library/authors -H "Authorization: Bearer $TOKEN"
curl -s http://localhost:8080/api/v1/library/authors/1 -H "Authorization: Bearer $TOKEN"

curl -s -X PUT http://localhost:8080/api/v1/library/authors/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "name": "Thomas H. Cormen", "biography": "Updated bio" }'
```

### Books
```bash
curl -s -X POST http://localhost:8080/api/v1/library/books \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Introduction to Algorithms",
    "isbn": "978-0262033848",
    "publisher": "MIT Press",
    "publication_year": 2009,
    "edition": "3rd",
    "total_copies": 5,
    "available_copies": 5,
    "location": "CS-Section-A3"
  }'

curl -s "http://localhost:8080/api/v1/library/books?search=algorithms" \
  -H "Authorization: Bearer $TOKEN"
curl -s http://localhost:8080/api/v1/library/books/1 -H "Authorization: Bearer $TOKEN"
```

### Book Authors (Many-to-Many)
```bash
curl -s -X POST http://localhost:8080/api/v1/library/books/1/authors \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "author_id": 1 }'

curl -s http://localhost:8080/api/v1/library/books/1/authors \
  -H "Authorization: Bearer $TOKEN"

curl -s -X DELETE http://localhost:8080/api/v1/library/books/1/authors/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Book Copies
```bash
curl -s -X POST http://localhost:8080/api/v1/library/book-copies \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "book_id": 1,
    "barcode": "LIB-CS-001-01",
    "copy_number": 1,
    "condition": "good",
    "shelf_location": "A3-R2-S1"
  }'

curl -s "http://localhost:8080/api/v1/library/book-copies?book_id=1" \
  -H "Authorization: Bearer $TOKEN"
```

### Digital Resources
```bash
curl -s -X POST http://localhost:8080/api/v1/library/digital-resources \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "IEEE Xplore Digital Library",
    "resource_type": "journal_database",
    "url": "https://ieeexplore.ieee.org",
    "access_link": "https://library.stateuni.edu/ieee",
    "publisher": "IEEE"
  }'

curl -s "http://localhost:8080/api/v1/library/digital-resources?search=ieee" \
  -H "Authorization: Bearer $TOKEN"
```

### Issue Book (triggers event)
```bash
curl -s -X POST http://localhost:8080/api/v1/library/circulations/issue \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "book_copy_id": 1,
    "student_id": 1
  }'
```
**Event**: `library.book_issued`

### Return Book (calculates fines, triggers overdue event if late)
```bash
curl -s -X POST http://localhost:8080/api/v1/library/circulations/1/return \
  -H "Authorization: Bearer $TOKEN"
```
**Event**: `library.book_returned` (and `library.book_overdue` if late, which triggers Finance fine invoice)

### Active Circulations
```bash
curl -s "http://localhost:8080/api/v1/library/circulations/active?student_id=1" \
  -H "Authorization: Bearer $TOKEN"
```

### Mark Overdue (scan all overdue books)
```bash
curl -s -X POST http://localhost:8080/api/v1/library/circulations/mark-overdue \
  -H "Authorization: Bearer $TOKEN"
```

### Reservations
```bash
curl -s -X POST http://localhost:8080/api/v1/library/reservations \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "book_id": 1,
    "student_id": 1
  }'

curl -s "http://localhost:8080/api/v1/library/reservations?book_id=1" \
  -H "Authorization: Bearer $TOKEN"
```

### Fines
```bash
curl -s "http://localhost:8080/api/v1/library/fines?student_id=1" \
  -H "Authorization: Bearer $TOKEN"

curl -s -X POST http://localhost:8080/api/v1/library/fines/1/pay \
  -H "Authorization: Bearer $TOKEN"
```

### Purchase Requests
```bash
curl -s -X POST http://localhost:8080/api/v1/library/purchase-requests \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "requested_by": 1,
    "title": "Design Patterns",
    "author": "Gang of Four",
    "isbn": "978-0201633610",
    "reason": "Needed for SE course"
  }'

curl -s http://localhost:8080/api/v1/library/purchase-requests \
  -H "Authorization: Bearer $TOKEN"

curl -s -X POST http://localhost:8080/api/v1/library/purchase-requests/1/approve \
  -H "Authorization: Bearer $TOKEN"
```

---

## 10. Hostel

### Stats
```bash
curl -s http://localhost:8080/api/v1/hostel/stats -H "Authorization: Bearer $TOKEN"
```

### Hostels
```bash
curl -s -X POST http://localhost:8080/api/v1/hostel/hostels \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Boys Hostel A",
    "code": "BH-A",
    "campus_id": 1,
    "total_rooms": 50,
    "contact_number": "022-87654321",
    "address": "Campus Road, Mumbai",
    "is_active": true
  }'

curl -s "http://localhost:8080/api/v1/hostel/hostels?campus_id=1" \
  -H "Authorization: Bearer $TOKEN"
curl -s http://localhost:8080/api/v1/hostel/hostels/1 -H "Authorization: Bearer $TOKEN"
```

### Rooms
```bash
curl -s -X POST http://localhost:8080/api/v1/hostel/rooms \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "hostel_id": 1,
    "room_number": "101",
    "room_type": "double",
    "capacity": 2,
    "monthly_rent": 5000
  }'

curl -s "http://localhost:8080/api/v1/hostel/rooms?hostel_id=1" \
  -H "Authorization: Bearer $TOKEN"
```

### Beds
```bash
curl -s -X POST http://localhost:8080/api/v1/hostel/beds \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "room_id": 1,
    "bed_number": "A"
  }'

curl -s "http://localhost:8080/api/v1/hostel/beds?room_id=1" \
  -H "Authorization: Bearer $TOKEN"
```

### Allocate Room (triggers event)
```bash
curl -s -X POST http://localhost:8080/api/v1/hostel/allocations/allocate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 1,
    "room_id": 1,
    "bed_id": 1
  }'
```
**Event**: `hostel.allocated`

### List Allocations
```bash
curl -s "http://localhost:8080/api/v1/hostel/allocations?hostel_id=1" \
  -H "Authorization: Bearer $TOKEN"
curl -s "http://localhost:8080/api/v1/hostel/allocations?student_id=1" \
  -H "Authorization: Bearer $TOKEN"
```

### Deallocate Room
```bash
curl -s -X POST http://localhost:8080/api/v1/hostel/allocations/1/deallocate \
  -H "Authorization: Bearer $TOKEN"
```

### Mess Bills
```bash
curl -s -X POST http://localhost:8080/api/v1/hostel/mess-bills \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 1,
    "month": "2026-05-01T00:00:00Z",
    "amount": 3500
  }'

curl -s "http://localhost:8080/api/v1/hostel/mess-bills?student_id=1" \
  -H "Authorization: Bearer $TOKEN"

curl -s -X POST http://localhost:8080/api/v1/hostel/mess-bills/1/pay \
  -H "Authorization: Bearer $TOKEN"
```

### Maintenance Requests (triggers event)
```bash
curl -s -X POST http://localhost:8080/api/v1/hostel/maintenance-requests \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 1,
    "room_id": 1,
    "category": "plumbing",
    "description": "Water tap leaking in bathroom"
  }'
```
**Event**: `hostel.maintenance_requested`

```bash
curl -s "http://localhost:8080/api/v1/hostel/maintenance-requests?room_id=1" \
  -H "Authorization: Bearer $TOKEN"

curl -s -X POST http://localhost:8080/api/v1/hostel/maintenance-requests/1/resolve \
  -H "Authorization: Bearer $TOKEN"
```

### Visitor Logs
```bash
curl -s -X POST http://localhost:8080/api/v1/hostel/visitor-logs \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "hostel_id": 1,
    "visitor_name": "Ramesh Sharma",
    "student_id": 1,
    "purpose": "Parent visit",
    "id_proof": "Aadhar - 1234-5678-9012"
  }'

curl -s "http://localhost:8080/api/v1/hostel/visitor-logs?hostel_id=1" \
  -H "Authorization: Bearer $TOKEN"

curl -s -X POST http://localhost:8080/api/v1/hostel/visitor-logs/1/exit \
  -H "Authorization: Bearer $TOKEN"
```

---

## 11. Transport

### Stats
```bash
curl -s http://localhost:8080/api/v1/transport/stats -H "Authorization: Bearer $TOKEN"
```

### Buses
```bash
curl -s -X POST http://localhost:8080/api/v1/transport/buses \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bus_number": "BUS-001",
    "registration_no": "MH-01-AB-1234",
    "capacity": 52,
    "is_active": true
  }'

curl -s "http://localhost:8080/api/v1/transport/buses?active=true" \
  -H "Authorization: Bearer $TOKEN"
curl -s http://localhost:8080/api/v1/transport/buses/1 -H "Authorization: Bearer $TOKEN"
```

### Routes
```bash
curl -s -X POST http://localhost:8080/api/v1/transport/routes \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "route_name": "Route 1 - Andheri to Campus",
    "description": "Morning pickup from Andheri station",
    "distance_km": 15.5,
    "estimated_time": "45 mins",
    "is_active": true
  }'

curl -s "http://localhost:8080/api/v1/transport/routes?active=true" \
  -H "Authorization: Bearer $TOKEN"
```

### Stops
```bash
curl -s -X POST http://localhost:8080/api/v1/transport/stops \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "route_id": 1,
    "stop_name": "Andheri Station",
    "stop_order": 1,
    "latitude": 19.1197,
    "longitude": 72.8464,
    "arrival_time": "07:30",
    "departure_time": "07:35"
  }'

curl -s "http://localhost:8080/api/v1/transport/stops?route_id=1" \
  -H "Authorization: Bearer $TOKEN"
```

### Bus Assignments
```bash
curl -s -X POST http://localhost:8080/api/v1/transport/bus-assignments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bus_id": 1,
    "route_id": 1
  }'

curl -s "http://localhost:8080/api/v1/transport/bus-assignments?bus_id=1" \
  -H "Authorization: Bearer $TOKEN"

# End assignment
curl -s -X POST http://localhost:8080/api/v1/transport/bus-assignments/1/end \
  -H "Authorization: Bearer $TOKEN"
```

### Issue Student Pass (triggers event)
```bash
curl -s -X POST http://localhost:8080/api/v1/transport/passes/issue \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 1,
    "route_id": 1,
    "pickup_stop_id": 1,
    "drop_stop_id": 3,
    "valid_from": "2026-06-01T00:00:00Z",
    "valid_to": "2027-05-31T00:00:00Z",
    "fee_paid": 15000
  }'
```
**Event**: `transport.pass_issued`

```bash
curl -s "http://localhost:8080/api/v1/transport/passes?student_id=1" \
  -H "Authorization: Bearer $TOKEN"
```

### Renew Pass
```bash
curl -s -X POST http://localhost:8080/api/v1/transport/passes/1/renew \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "valid_to": "2028-05-31",
    "fee_paid": 15000
  }'
```

### Vehicle Maintenance
```bash
curl -s -X POST http://localhost:8080/api/v1/transport/maintenance \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bus_id": 1,
    "description": "Engine oil change and tire rotation",
    "cost": 8500,
    "next_due_date": "2026-11-20T00:00:00Z"
  }'

curl -s "http://localhost:8080/api/v1/transport/maintenance?bus_id=1" \
  -H "Authorization: Bearer $TOKEN"
```

---

## Cross-Module Event Workflow Tests

These test the event-driven architecture. Perform steps in order.

### Workflow 1: Admission -> Student Auto-Creation
```bash
# Step 1: Register a user for the applicant
curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{ "username": "applicant1", "email": "applicant1@test.com", "password": "test123", "role_name": "student" }'

# Step 2: Create admission cycle + program (prerequisites)
# (see Core and Admissions sections above)

# Step 3: Submit application
curl -s -X POST http://localhost:8080/api/v1/admissions/applicants \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "cycle_id": 1, "program_id": 1, "first_name": "Test", "last_name": "Student", "email": "test@student.com", "phone": "9999999999" }'

# Step 4: Approve application (this triggers auto student creation)
curl -s -X PUT http://localhost:8080/api/v1/admissions/applicants/1/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "status_id": 2 }'

# Step 5: Verify student was auto-created
curl -s http://localhost:8080/api/v1/students -H "Authorization: Bearer $TOKEN"
```

### Workflow 2: Term Registration -> Finance Invoice
```bash
# Step 1: Register student for term
curl -s -X POST http://localhost:8080/api/v1/academic/term-registrations \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "student_id": 1, "academic_term_id": 1, "batch_id": 1 }'

# Step 2: Check if invoice was auto-generated
curl -s http://localhost:8080/api/v1/finance/invoices/student/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Workflow 3: Library Overdue -> Finance Fine
```bash
# Step 1: Issue a book
curl -s -X POST http://localhost:8080/api/v1/library/circulations/issue \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "book_copy_id": 1, "student_id": 1 }'

# Step 2: Return the book (if overdue, fine is calculated and event fires)
curl -s -X POST http://localhost:8080/api/v1/library/circulations/1/return \
  -H "Authorization: Bearer $TOKEN"

# Step 3: Check library fines
curl -s "http://localhost:8080/api/v1/library/fines?student_id=1" \
  -H "Authorization: Bearer $TOKEN"

# Step 4: Check finance invoices for the fine
curl -s http://localhost:8080/api/v1/finance/invoices/student/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Workflow 4: HR Payroll -> Finance Voucher
```bash
# Step 1: Run payroll
curl -s -X POST http://localhost:8080/api/v1/hr/payroll/run \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "employee_id": 1, "month": "2026-05" }'

# Step 2: Check finance for payroll voucher (via event)
# The finance module subscribes to hr.payroll_processed and records a voucher
```

---

## Common Error Responses

| Status | Code | Meaning |
|--------|------|---------|
| 400 | `BAD_REQUEST` | Missing or invalid input fields |
| 401 | `UNAUTHORIZED` | Missing or invalid JWT token |
| 404 | `NOT_FOUND` | Resource does not exist |
| 409 | `CONFLICT` | Duplicate entry or state conflict |
| 500 | `INTERNAL_ERROR` | Server-side error |

Error format:
```json
{
  "error": {
    "code": "BAD_REQUEST",
    "message": "hostel name and code are required"
  }
}
```

---

## Quick Smoke Test Script

```bash
#!/bin/bash
BASE="http://localhost:8080"

echo "=== 1. Register ==="
curl -s -X POST $BASE/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"smoketest","email":"smoke@test.com","password":"test123","role_name":"admin"}'

echo -e "\n=== 2. Login ==="
TOKEN=$(curl -s -X POST $BASE/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"smoketest","password":"test123"}' | jq -r '.token')
echo "Token: ${TOKEN:0:20}..."

echo -e "\n=== 3. Profile ==="
curl -s $BASE/api/v1/auth/profile -H "Authorization: Bearer $TOKEN" | jq .

echo -e "\n=== 4. Lookups (public) ==="
curl -s $BASE/api/v1/lookups/genders | jq '.[0]'

echo -e "\n=== 5. Finance Summary ==="
curl -s $BASE/api/v1/finance/summary -H "Authorization: Bearer $TOKEN" | jq .

echo -e "\n=== 6. Library Stats ==="
curl -s $BASE/api/v1/library/stats -H "Authorization: Bearer $TOKEN" | jq .

echo -e "\n=== 7. Hostel Stats ==="
curl -s $BASE/api/v1/hostel/stats -H "Authorization: Bearer $TOKEN" | jq .

echo -e "\n=== 8. Transport Stats ==="
curl -s $BASE/api/v1/transport/stats -H "Authorization: Bearer $TOKEN" | jq .

echo -e "\n=== 9. HR Stats ==="
curl -s $BASE/api/v1/hr/stats -H "Authorization: Bearer $TOKEN" | jq .

echo -e "\n=== All smoke tests passed ==="
```
