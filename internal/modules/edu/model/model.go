package model

import "time"

type Teacher struct {
	ID             uint64    `gorm:"primaryKey"`
	Name           string    `gorm:"size:64;not null"`
	Mobile         string    `gorm:"size:32"`
	Title          string    `gorm:"size:64"`
	MainSubject    string    `gorm:"column:main_subject;size:64"`
	EmploymentType string    `gorm:"column:employment_type;size:32"`
	WeeklyHours    int       `gorm:"column:weekly_hours;not null;default:0"`
	Campus         string    `gorm:"size:64"`
	Status         string    `gorm:"size:32;not null;default:'在职'"`
	Remark         string    `gorm:"size:255"`
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
}

type Student struct {
	ID             uint64    `gorm:"primaryKey"`
	Name           string    `gorm:"size:64;not null"`
	Gender         string    `gorm:"size:16"`
	SchoolName     string    `gorm:"column:school_name;size:128"`
	GradeName      string    `gorm:"column:grade_name;size:64"`
	Campus         string    `gorm:"size:64"`
	RemainingHours int       `gorm:"column:remaining_hours;not null;default:0"`
	Status         string    `gorm:"size:32;not null;default:'在读'"`
	Remark         string    `gorm:"size:255"`
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
}

type StudentGuardian struct {
	ID        uint64    `gorm:"primaryKey"`
	StudentID uint64    `gorm:"column:student_id;not null;index"`
	Name      string    `gorm:"size:64;not null"`
	Relation  string    `gorm:"size:32"`
	Mobile    string    `gorm:"size:32;not null"`
	IsPrimary bool      `gorm:"column:is_primary;not null;default:false"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type Course struct {
	ID                    uint64    `gorm:"primaryKey"`
	Name                  string    `gorm:"size:128;not null"`
	Category              string    `gorm:"size:64"`
	Description           string    `gorm:"type:text"`
	AgeRange              string    `gorm:"column:age_range;size:64"`
	LessonDurationMinutes int       `gorm:"column:lesson_duration_minutes;not null;default:0"`
	TotalLessons          int       `gorm:"column:total_lessons;not null;default:0"`
	DeliveryType          string    `gorm:"column:delivery_type;size:32"`
	Status                string    `gorm:"size:32;not null;default:'启用'"`
	CreatedAt             time.Time `gorm:"not null"`
	UpdatedAt             time.Time `gorm:"not null"`
}

type Class struct {
	ID             uint64     `gorm:"primaryKey"`
	Name           string     `gorm:"size:128;not null"`
	CourseID       uint64     `gorm:"column:course_id;not null;index"`
	TeacherID      uint64     `gorm:"column:teacher_id;not null;index"`
	Campus         string     `gorm:"size:64"`
	Capacity       int        `gorm:"not null;default:0"`
	WeeklySchedule string     `gorm:"column:weekly_schedule;size:64"`
	StartDate      *time.Time `gorm:"column:start_date;type:date"`
	EndDate        *time.Time `gorm:"column:end_date;type:date"`
	Status         string     `gorm:"size:32;not null;default:'开班中'"`
	Remark         string     `gorm:"size:255"`
	CreatedAt      time.Time  `gorm:"not null"`
	UpdatedAt      time.Time  `gorm:"not null"`
}

type ClassStudent struct {
	ID        uint64     `gorm:"primaryKey"`
	ClassID   uint64     `gorm:"column:class_id;not null;index"`
	StudentID uint64     `gorm:"column:student_id;not null;index"`
	JoinDate  *time.Time `gorm:"column:join_date;type:date"`
	LeaveDate *time.Time `gorm:"column:leave_date;type:date"`
	Status    string     `gorm:"size:32;not null;default:'在读'"`
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
}

type ClassSchedule struct {
	ID           uint64    `gorm:"primaryKey"`
	ClassID      uint64    `gorm:"column:class_id;not null;index"`
	CourseID     uint64    `gorm:"column:course_id;not null;index"`
	TeacherID    uint64    `gorm:"column:teacher_id;not null;index"`
	ScheduleType string    `gorm:"column:schedule_type;size:32;not null;default:'常规课'"`
	ScheduleDate time.Time `gorm:"column:schedule_date;type:date;not null;index"`
	StartTime    string    `gorm:"column:start_time;size:8;not null"`
	EndTime      string    `gorm:"column:end_time;size:8;not null"`
	Location     string    `gorm:"size:128"`
	Status       string    `gorm:"size:32;not null;default:'待上课'"`
	Remark       string    `gorm:"size:255"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}

type Notice struct {
	ID             uint64     `gorm:"primaryKey"`
	Title          string     `gorm:"size:128;not null"`
	Content        string     `gorm:"type:text;not null"`
	NoticeType     string     `gorm:"column:notice_type;size:32;not null;default:'校区通知'"`
	TargetScope    string     `gorm:"column:target_scope;size:128;not null"`
	AuthorName     string     `gorm:"column:author_name;size:64;not null"`
	RelatedClassID *uint64    `gorm:"column:related_class_id"`
	Status         string     `gorm:"size:32;not null;default:'草稿'"`
	PublishAt      *time.Time `gorm:"column:publish_at"`
	CreatedAt      time.Time  `gorm:"not null"`
	UpdatedAt      time.Time  `gorm:"not null"`
}

type AttendanceRecord struct {
	ID         uint64     `gorm:"primaryKey"`
	ScheduleID uint64     `gorm:"column:schedule_id;not null;index;uniqueIndex:uk_schedule_student"`
	StudentID  uint64     `gorm:"column:student_id;not null;index;uniqueIndex:uk_schedule_student"`
	Status     string     `gorm:"size:32;not null;default:'待确认'"`
	Remark     string     `gorm:"size:255"`
	CheckedAt  *time.Time `gorm:"column:checked_at"`
	UpdatedBy  string     `gorm:"column:updated_by;size:64"`
	CreatedAt  time.Time  `gorm:"not null"`
	UpdatedAt  time.Time  `gorm:"not null"`
}
