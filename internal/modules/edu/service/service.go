package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"edu-admin/internal/modules/demo"
	edumodel "edu-admin/internal/modules/edu/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	activeClassStudentStatus = "在读"
	dateLayout               = "2006-01-02"
	dateTimeLayout           = "2006-01-02 15:04"
)

type Service struct {
	db *gorm.DB
}

type Option struct {
	Value uint64 `json:"value"`
	Label string `json:"label"`
}

type TeacherItem struct {
	ID             uint64 `json:"id"`
	Name           string `json:"name"`
	Mobile         string `json:"mobile"`
	Title          string `json:"title"`
	MainSubject    string `json:"mainSubject" gorm:"column:main_subject"`
	EmploymentType string `json:"employmentType" gorm:"column:employment_type"`
	WeeklyHours    int    `json:"weeklyHours" gorm:"column:weekly_hours"`
	Campus         string `json:"campus"`
	Status         string `json:"status"`
	Remark         string `json:"remark"`
}

type CourseFilter struct {
	Keyword  string
	Category string
	Status   string
}

type TeacherPayload struct {
	Name           string
	Mobile         string
	Title          string
	MainSubject    string
	EmploymentType string
	WeeklyHours    int
	Campus         string
	Status         string
	Remark         string
}

type CoursePayload struct {
	Name                  string
	Category              string
	Description           string
	AgeRange              string
	LessonDurationMinutes int
	TotalLessons          int
	DeliveryType          string
	Status                string
}

type CourseItem struct {
	ID                    uint64 `json:"id"`
	Name                  string `json:"name"`
	Category              string `json:"category"`
	Description           string `json:"description"`
	AgeRange              string `json:"ageRange" gorm:"column:age_range"`
	LessonDurationMinutes int    `json:"lessonDurationMinutes" gorm:"column:lesson_duration_minutes"`
	TotalLessons          int    `json:"totalLessons" gorm:"column:total_lessons"`
	DeliveryType          string `json:"deliveryType" gorm:"column:delivery_type"`
	Status                string `json:"status"`
	ClassCount            int    `json:"classCount" gorm:"column:class_count"`
}

type StudentItem struct {
	ID             uint64 `json:"id"`
	Name           string `json:"name"`
	Grade          string `json:"grade"`
	ParentName     string `json:"parentName" gorm:"column:parent_name"`
	ParentMobile   string `json:"parentMobile" gorm:"column:parent_mobile"`
	Campus         string `json:"campus"`
	ClassID        uint64 `json:"classId" gorm:"column:class_id"`
	ClassName      string `json:"className" gorm:"column:class_name"`
	RemainingHours int    `json:"remainingHours" gorm:"column:remaining_hours"`
	Status         string `json:"status"`
}

type StudentPayload struct {
	Name             string
	Gender           string
	SchoolName       string
	GradeName        string
	Campus           string
	RemainingHours   int
	Status           string
	Remark           string
	GuardianName     string
	GuardianMobile   string
	GuardianRelation string
}

type StudentProfile struct {
	ID             uint64 `json:"id"`
	Name           string `json:"name"`
	Gender         string `json:"gender"`
	SchoolName     string `json:"schoolName"`
	Grade          string `json:"grade"`
	ParentName     string `json:"parentName"`
	ParentMobile   string `json:"parentMobile"`
	Campus         string `json:"campus"`
	RemainingHours int    `json:"remainingHours"`
	Status         string `json:"status"`
	Remark         string `json:"remark"`
}

type StudentGuardianItem struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Relation  string `json:"relation"`
	Mobile    string `json:"mobile"`
	IsPrimary bool   `json:"isPrimary"`
}

type StudentGuardianPayload struct {
	Name      string
	Relation  string
	Mobile    string
	IsPrimary bool
}

type StudentAttendanceRecord struct {
	ScheduleID  uint64 `json:"scheduleId"`
	ClassID     uint64 `json:"classId"`
	ClassName   string `json:"className"`
	CourseName  string `json:"courseName"`
	TeacherName string `json:"teacherName"`
	Campus      string `json:"campus"`
	Classroom   string `json:"classroom"`
	LessonDate  string `json:"lessonDate"`
	LessonTime  string `json:"lessonTime"`
	Status      string `json:"status"`
	Remark      string `json:"remark"`
}

type ClassItem struct {
	ID             uint64 `json:"id"`
	CourseID       uint64 `json:"courseId" gorm:"column:course_id"`
	Name           string `json:"name"`
	CourseName     string `json:"courseName" gorm:"column:course_name"`
	TeacherID      uint64 `json:"teacherId" gorm:"column:teacher_id"`
	TeacherName    string `json:"teacherName" gorm:"column:teacher_name"`
	Campus         string `json:"campus"`
	StudentCount   int    `json:"studentCount" gorm:"column:student_count"`
	Capacity       int    `json:"capacity"`
	WeeklySchedule string `json:"weeklySchedule" gorm:"column:weekly_schedule"`
	StartDate      string `json:"startDate" gorm:"column:start_date"`
	EndDate        string `json:"endDate" gorm:"column:end_date"`
	Status         string `json:"status"`
	Remark         string `json:"remark"`
}

type TeacherDetail struct {
	Teacher         TeacherItem    `json:"teacher"`
	Classes         []ClassItem    `json:"classes"`
	RecentSchedules []ScheduleItem `json:"recentSchedules"`
}

type StudentDetail struct {
	Student          StudentProfile            `json:"student"`
	Guardians        []StudentGuardianItem     `json:"guardians"`
	Classes          []ClassItem               `json:"classes"`
	RecentSchedules  []ScheduleItem            `json:"recentSchedules"`
	RecentAttendance []StudentAttendanceRecord `json:"recentAttendance"`
	RecentHomeworks  []HomeworkItem            `json:"recentHomeworks"`
	RecentFeedbacks  []FeedbackItem            `json:"recentFeedbacks"`
}

type ClassDetail struct {
	Class             ClassItem               `json:"class"`
	Students          []StudentItem           `json:"students"`
	UpcomingSchedules []ScheduleItem          `json:"upcomingSchedules"`
	RecentAttendance  []AttendanceSessionItem `json:"recentAttendance"`
	RecentHomeworks   []HomeworkItem          `json:"recentHomeworks"`
	RecentNotices     []NoticeItem            `json:"recentNotices"`
}

type ClassStudentPayload struct {
	StudentIDs []uint64 `json:"studentIds"`
}

type ClassPayload struct {
	Name           string
	CourseID       uint64
	TeacherID      uint64
	Campus         string
	Capacity       int
	WeeklySchedule string
	StartDate      string
	EndDate        string
	Status         string
	Remark         string
}

type ScheduleItem struct {
	ID               uint64 `json:"id"`
	ClassID          uint64 `json:"classId" gorm:"column:class_id"`
	SourceScheduleID uint64 `json:"sourceScheduleId" gorm:"column:source_schedule_id"`
	ClassName        string `json:"className" gorm:"column:class_name"`
	CourseName       string `json:"courseName" gorm:"column:course_name"`
	TeacherName      string `json:"teacherName" gorm:"column:teacher_name"`
	TeacherID        uint64 `json:"teacherId"`
	Campus           string `json:"campus"`
	Classroom        string `json:"classroom"`
	ScheduleType     string `json:"scheduleType"`
	LessonDate       string `json:"lessonDate"`
	LessonTime       string `json:"lessonTime"`
	AttendanceStatus string `json:"attendanceStatus"`
	Remark           string `json:"remark"`
}

type SchedulePayload struct {
	ClassID      uint64 `json:"classId"`
	ScheduleType string `json:"scheduleType"`
	LessonDate   string `json:"lessonDate"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
	Classroom    string `json:"classroom"`
	Remark       string `json:"remark"`
}

type ScheduleActionPayload struct {
	LessonDate string `json:"lessonDate"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	Classroom  string `json:"classroom"`
	Remark     string `json:"remark"`
}

type ScheduleDetail struct {
	Schedule       ScheduleItem          `json:"schedule"`
	Class          ClassItem             `json:"class"`
	Students       []StudentItem         `json:"students"`
	Attendance     AttendanceSessionItem `json:"attendance"`
	Homework       *HomeworkItem         `json:"homework"`
	Feedback       *FeedbackItem         `json:"feedback"`
	RelatedNotices []NoticeItem          `json:"relatedNotices"`
}

type NoticeItem struct {
	ID                uint64 `json:"id"`
	Title             string `json:"title"`
	Content           string `json:"content"`
	Category          string `json:"category"`
	TargetScope       string `json:"targetScope" gorm:"column:target_scope"`
	RelatedClassID    uint64 `json:"relatedClassId" gorm:"column:related_class_id"`
	RelatedScheduleID uint64 `json:"relatedScheduleId" gorm:"column:related_schedule_id"`
	Status            string `json:"status"`
	PublishAt         string `json:"publishAt"`
	Author            string `json:"author"`
}

type NoticeTargetItem struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Campus string `json:"campus"`
}

type NoticePayload struct {
	Title             string `json:"title"`
	Content           string `json:"content"`
	Category          string `json:"category"`
	TargetScope       string `json:"targetScope"`
	RelatedClassID    uint64
	RelatedScheduleID uint64
	Status            string `json:"status"`
	Author            string `json:"author"`
}

type NoticeFilter struct {
	ClassID uint64
	Status  string
	Date    string
}

type AttendanceItem struct {
	StudentID    uint64 `json:"studentId"`
	StudentName  string `json:"studentName"`
	Grade        string `json:"grade"`
	ParentMobile string `json:"parentMobile"`
	Status       string `json:"status"`
	Remark       string `json:"remark"`
	UpdatedBy    string `json:"updatedBy"`
	UpdatedAt    string `json:"updatedAt"`
}

type AttendanceSessionItem struct {
	ID               uint64 `json:"id"`
	ClassID          uint64 `json:"classId"`
	ClassName        string `json:"className"`
	CourseName       string `json:"courseName"`
	TeacherName      string `json:"teacherName"`
	Campus           string `json:"campus"`
	Classroom        string `json:"classroom"`
	LessonDate       string `json:"lessonDate"`
	LessonTime       string `json:"lessonTime"`
	AttendanceStatus string `json:"attendanceStatus"`
	StudentCount     int    `json:"studentCount"`
	PresentCount     int    `json:"presentCount"`
	LeaveCount       int    `json:"leaveCount"`
	AbsentCount      int    `json:"absentCount"`
	PendingCount     int    `json:"pendingCount"`
}

type AttendanceRecordFilter struct {
	ClassID   uint64
	StudentID uint64
	Date      string
	Status    string
}

type AttendanceRecordItem struct {
	ScheduleID   uint64 `json:"scheduleId"`
	ClassID      uint64 `json:"classId"`
	ClassName    string `json:"className"`
	StudentID    uint64 `json:"studentId"`
	StudentName  string `json:"studentName"`
	TeacherName  string `json:"teacherName"`
	LessonDate   string `json:"lessonDate"`
	LessonTime   string `json:"lessonTime"`
	Status       string `json:"status"`
	Remark       string `json:"remark"`
	UpdatedBy    string `json:"updatedBy"`
	UpdatedAt    string `json:"updatedAt"`
	ParentMobile string `json:"parentMobile"`
}

type AttendanceSaveItem struct {
	StudentID uint64 `json:"studentId"`
	Status    string `json:"status"`
	Remark    string `json:"remark"`
}

type AttendanceSavePayload struct {
	Items []AttendanceSaveItem `json:"items"`
}

type HomeworkItem struct {
	ID             uint64 `json:"id"`
	ScheduleID     uint64 `json:"scheduleId"`
	ClassID        uint64 `json:"classId"`
	ClassName      string `json:"className"`
	CourseName     string `json:"courseName"`
	TeacherName    string `json:"teacherName"`
	LessonDate     string `json:"lessonDate"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	SubmissionNote string `json:"submissionNote"`
	Status         string `json:"status"`
}

type HomeworkPayload struct {
	Title          string `json:"title"`
	Content        string `json:"content"`
	SubmissionNote string `json:"submissionNote"`
	Status         string `json:"status"`
}

type FeedbackItem struct {
	ID             uint64 `json:"id"`
	ScheduleID     uint64 `json:"scheduleId"`
	ClassID        uint64 `json:"classId"`
	ClassName      string `json:"className"`
	CourseName     string `json:"courseName"`
	TeacherName    string `json:"teacherName"`
	LessonDate     string `json:"lessonDate"`
	Summary        string `json:"summary"`
	LearningStatus string `json:"learningStatus"`
	NextSuggestion string `json:"nextSuggestion"`
	ParentNotice   string `json:"parentNotice"`
}

type FeedbackPayload struct {
	Summary        string `json:"summary"`
	LearningStatus string `json:"learningStatus"`
	NextSuggestion string `json:"nextSuggestion"`
	ParentNotice   string `json:"parentNotice"`
}

func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Bootstrap(autoSeed bool) error {
	if s.db == nil {
		return nil
	}

	migrateErr := s.db.AutoMigrate(
		&edumodel.User{},
		&edumodel.Role{},
		&edumodel.UserRole{},
		&edumodel.OperationLog{},
		&edumodel.Teacher{},
		&edumodel.Student{},
		&edumodel.StudentGuardian{},
		&edumodel.Course{},
		&edumodel.Class{},
		&edumodel.ClassStudent{},
		&edumodel.ClassSchedule{},
		&edumodel.AttendanceRecord{},
		&edumodel.Homework{},
		&edumodel.ClassFeedback{},
		&edumodel.Notice{},
	)
	if migrateErr != nil {
		return migrateErr
	}

	if !autoSeed {
		return nil
	}

	return s.seedIfEmpty()
}

func (s *Service) Overview() (map[string]any, error) {
	if s.db == nil {
		return demo.Overview(), nil
	}

	var studentCount int64
	studentCountErr := s.db.Model(&edumodel.Student{}).Count(&studentCount).Error
	if studentCountErr != nil {
		return nil, studentCountErr
	}

	var classCount int64
	classCountErr := s.db.Model(&edumodel.Class{}).Count(&classCount).Error
	if classCountErr != nil {
		return nil, classCountErr
	}

	attendanceSessions, sessionErr := s.AttendanceSessions()
	if sessionErr != nil {
		return nil, sessionErr
	}

	today := time.Now().Format(dateLayout)
	todayCourses := 0
	todayPendingCheck := 0
	todayLeaveCount := 0
	todayAbsentCount := 0

	for _, item := range attendanceSessions {
		if item.LessonDate != today {
			continue
		}

		todayCourses++
		todayLeaveCount += item.LeaveCount
		todayAbsentCount += item.AbsentCount

		if item.AttendanceStatus == "待签到" {
			todayPendingCheck++
		}
	}

	var pendingActionCount int64
	pendingActionErr := s.db.Model(&edumodel.Notice{}).
		Where("status IN ?", []string{"草稿", "待发送"}).
		Count(&pendingActionCount).Error
	if pendingActionErr != nil {
		return nil, pendingActionErr
	}

	upcomingLessons, scheduleErr := s.Schedules()
	if scheduleErr != nil {
		return nil, scheduleErr
	}

	latestNotices, noticeErr := s.Notices()
	if noticeErr != nil {
		return nil, noticeErr
	}

	if len(upcomingLessons) > 5 {
		upcomingLessons = upcomingLessons[:5]
	}

	if len(latestNotices) > 5 {
		latestNotices = latestNotices[:5]
	}

	return map[string]any{
		"todayCourses":       todayCourses,
		"todayPendingCheck":  todayPendingCheck,
		"todayLeaveCount":    todayLeaveCount,
		"todayAbsentCount":   todayAbsentCount,
		"studentCount":       studentCount,
		"classCount":         classCount,
		"pendingActionCount": pendingActionCount,
		"upcomingLessons":    upcomingLessons,
		"latestNotices":      latestNotices,
	}, nil
}

func (s *Service) Teachers() ([]TeacherItem, error) {
	if s.db == nil {
		return teacherItemsFromDemo(), nil
	}

	var items []TeacherItem
	listErr := s.db.Model(&edumodel.Teacher{}).
		Select("id, name, mobile, title, main_subject, employment_type, weekly_hours, campus, status, remark").
		Order("id ASC").
		Scan(&items).Error
	if listErr != nil {
		return nil, listErr
	}

	return items, nil
}

func (s *Service) Teacher(rawID string) (TeacherItem, bool, error) {
	if s.db == nil {
		return teacherItemFromDemo(rawID)
	}

	var item TeacherItem
	findErr := s.db.Model(&edumodel.Teacher{}).
		Select("id, name, mobile, title, main_subject, employment_type, weekly_hours, campus, status, remark").
		Where("id = ?", rawID).
		Scan(&item).Error
	if findErr != nil {
		return TeacherItem{}, false, findErr
	}

	if item.ID == 0 {
		return TeacherItem{}, false, nil
	}

	return item, true, nil
}

func (s *Service) CreateTeacher(input TeacherPayload) (TeacherItem, error) {
	if s.db == nil {
		return TeacherItem{
			ID:             1,
			Name:           input.Name,
			Mobile:         input.Mobile,
			Title:          input.Title,
			MainSubject:    input.MainSubject,
			EmploymentType: input.EmploymentType,
			WeeklyHours:    input.WeeklyHours,
			Campus:         input.Campus,
			Status:         input.Status,
			Remark:         input.Remark,
		}, nil
	}

	now := time.Now()
	teacher := edumodel.Teacher{
		Name:           input.Name,
		Mobile:         input.Mobile,
		Title:          input.Title,
		MainSubject:    input.MainSubject,
		EmploymentType: input.EmploymentType,
		WeeklyHours:    input.WeeklyHours,
		Campus:         input.Campus,
		Status:         input.Status,
		Remark:         input.Remark,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	createErr := s.db.Create(&teacher).Error
	if createErr != nil {
		return TeacherItem{}, createErr
	}

	createdItem, found, teacherErr := s.Teacher(fmt.Sprintf("%d", teacher.ID))
	if teacherErr != nil {
		return TeacherItem{}, teacherErr
	}
	if !found {
		return TeacherItem{
			ID:             teacher.ID,
			Name:           teacher.Name,
			Mobile:         teacher.Mobile,
			Title:          teacher.Title,
			MainSubject:    teacher.MainSubject,
			EmploymentType: teacher.EmploymentType,
			WeeklyHours:    teacher.WeeklyHours,
			Campus:         teacher.Campus,
			Status:         teacher.Status,
			Remark:         teacher.Remark,
		}, nil
	}

	return createdItem, nil
}

func (s *Service) UpdateTeacher(rawID string, input TeacherPayload) (TeacherItem, bool, error) {
	if s.db == nil {
		item, found, teacherErr := s.Teacher(rawID)
		if teacherErr != nil || !found {
			return item, found, teacherErr
		}

		item.Name = input.Name
		item.Mobile = input.Mobile
		item.Title = input.Title
		item.MainSubject = input.MainSubject
		item.EmploymentType = input.EmploymentType
		item.WeeklyHours = input.WeeklyHours
		item.Campus = input.Campus
		item.Status = input.Status
		item.Remark = input.Remark

		return item, true, nil
	}

	updateResult := s.db.Model(&edumodel.Teacher{}).
		Where("id = ?", rawID).
		Updates(map[string]any{
			"name":            input.Name,
			"mobile":          input.Mobile,
			"title":           input.Title,
			"main_subject":    input.MainSubject,
			"employment_type": input.EmploymentType,
			"weekly_hours":    input.WeeklyHours,
			"campus":          input.Campus,
			"status":          input.Status,
			"remark":          input.Remark,
			"updated_at":      time.Now(),
		})
	if updateResult.Error != nil {
		return TeacherItem{}, false, updateResult.Error
	}
	if updateResult.RowsAffected == 0 {
		return TeacherItem{}, false, nil
	}

	updatedItem, found, teacherErr := s.Teacher(rawID)
	if teacherErr != nil {
		return TeacherItem{}, false, teacherErr
	}

	return updatedItem, found, nil
}

func (s *Service) TeacherOptions() ([]Option, error) {
	if s.db == nil {
		options := make([]Option, 0, len(demo.TeacherOptions()))
		for _, option := range demo.TeacherOptions() {
			options = append(options, Option{
				Value: uint64(option.Value),
				Label: option.Label,
			})
		}
		return options, nil
	}

	teachers, teacherErr := s.Teachers()
	if teacherErr != nil {
		return nil, teacherErr
	}

	options := make([]Option, 0, len(teachers))
	for _, teacher := range teachers {
		options = append(options, Option{
			Value: teacher.ID,
			Label: teacher.Name,
		})
	}

	return options, nil
}

func (s *Service) TeacherDetail(rawID string) (TeacherDetail, bool, error) {
	const recentItemLimit = 6

	teacherItem, found, teacherErr := s.Teacher(rawID)
	if teacherErr != nil {
		return TeacherDetail{}, false, teacherErr
	}
	if !found {
		return TeacherDetail{}, false, nil
	}

	allClasses, classErr := s.Classes()
	if classErr != nil {
		return TeacherDetail{}, false, classErr
	}

	teacherClasses := make([]ClassItem, 0, len(allClasses))
	for _, classItem := range allClasses {
		if classItem.TeacherID != teacherItem.ID {
			continue
		}

		teacherClasses = append(teacherClasses, classItem)
	}

	allSchedules, scheduleErr := s.Schedules()
	if scheduleErr != nil {
		return TeacherDetail{}, false, scheduleErr
	}

	teacherSchedules := make([]ScheduleItem, 0, len(allSchedules))
	for _, scheduleItem := range allSchedules {
		if scheduleItem.TeacherID != teacherItem.ID {
			continue
		}

		teacherSchedules = append(teacherSchedules, scheduleItem)
	}

	return TeacherDetail{
		Teacher:         teacherItem,
		Classes:         teacherClasses,
		RecentSchedules: teacherRecentSchedules(teacherSchedules, recentItemLimit),
	}, true, nil
}

func (s *Service) Courses(filter CourseFilter) ([]CourseItem, error) {
	if s.db == nil {
		return courseItemsFromDemo(demo.Courses()), nil
	}

	courseQuery := s.courseQuery(filter)

	var items []CourseItem
	listErr := courseQuery.Order("c.id ASC").Scan(&items).Error
	if listErr != nil {
		return nil, listErr
	}

	return items, nil
}

func (s *Service) Course(rawID string) (CourseItem, bool, error) {
	if s.db == nil {
		return courseItemFromDemo(rawID)
	}

	courseQuery := s.courseQuery(CourseFilter{}).Where("c.id = ?", rawID)

	var item CourseItem
	findErr := courseQuery.Limit(1).Scan(&item).Error
	if findErr != nil {
		return CourseItem{}, false, findErr
	}

	if item.ID == 0 {
		return CourseItem{}, false, nil
	}

	return item, true, nil
}

func (s *Service) CourseOptions() ([]Option, error) {
	if s.db == nil {
		options := make([]Option, 0, len(demo.CourseOptions()))
		for _, option := range demo.CourseOptions() {
			options = append(options, Option{
				Value: uint64(option.Value),
				Label: option.Label,
			})
		}
		return options, nil
	}

	var options []Option
	listErr := s.db.Model(&edumodel.Course{}).
		Select("id AS value, name AS label").
		Order("id ASC").
		Scan(&options).Error
	if listErr != nil {
		return nil, listErr
	}

	return options, nil
}

func (s *Service) CreateCourse(input CoursePayload) (CourseItem, error) {
	if s.db == nil {
		return courseItemFromPayload(1, input), nil
	}

	course := edumodel.Course{
		Name:                  input.Name,
		Category:              input.Category,
		Description:           input.Description,
		AgeRange:              input.AgeRange,
		LessonDurationMinutes: input.LessonDurationMinutes,
		TotalLessons:          input.TotalLessons,
		DeliveryType:          input.DeliveryType,
		Status:                input.Status,
	}

	createErr := s.db.Create(&course).Error
	if createErr != nil {
		return CourseItem{}, createErr
	}

	createdItem, found, detailErr := s.Course(fmt.Sprintf("%d", course.ID))
	if detailErr != nil {
		return CourseItem{}, detailErr
	}
	if !found {
		return courseItemFromPayload(course.ID, input), nil
	}

	return createdItem, nil
}

func (s *Service) UpdateCourse(rawID string, input CoursePayload) (CourseItem, bool, error) {
	if s.db == nil {
		item, found, itemErr := courseItemFromDemo(rawID)
		if itemErr != nil || !found {
			return item, found, itemErr
		}

		item.Name = input.Name
		item.Category = input.Category
		item.Description = input.Description
		item.AgeRange = input.AgeRange
		item.LessonDurationMinutes = input.LessonDurationMinutes
		item.TotalLessons = input.TotalLessons
		item.DeliveryType = input.DeliveryType
		item.Status = input.Status

		return item, true, nil
	}

	updateValues := map[string]any{
		"name":                    input.Name,
		"category":                input.Category,
		"description":             input.Description,
		"age_range":               input.AgeRange,
		"lesson_duration_minutes": input.LessonDurationMinutes,
		"total_lessons":           input.TotalLessons,
		"delivery_type":           input.DeliveryType,
		"status":                  input.Status,
	}

	updateResult := s.db.Model(&edumodel.Course{}).
		Where("id = ?", rawID).
		Updates(updateValues)
	if updateResult.Error != nil {
		return CourseItem{}, false, updateResult.Error
	}
	if updateResult.RowsAffected == 0 {
		return CourseItem{}, false, nil
	}

	updatedItem, found, detailErr := s.Course(rawID)
	if detailErr != nil {
		return CourseItem{}, false, detailErr
	}

	return updatedItem, found, nil
}

func (s *Service) Students() ([]StudentItem, error) {
	if s.db == nil {
		return studentItemsFromDemo(), nil
	}

	query := `
SELECT
  s.id,
  s.name,
  s.grade_name AS grade,
  COALESCE(g.name, '') AS parent_name,
  COALESCE(g.mobile, '') AS parent_mobile,
  s.campus,
  COALESCE(cs.class_id, 0) AS class_id,
  COALESCE(c.name, '') AS class_name,
  s.remaining_hours,
  s.status
FROM students AS s
LEFT JOIN student_guardians AS g
  ON g.student_id = s.id AND g.is_primary = 1
LEFT JOIN (
  SELECT student_id, MIN(class_id) AS class_id
  FROM class_students
  WHERE status = ?
  GROUP BY student_id
) AS cs
  ON cs.student_id = s.id
LEFT JOIN classes AS c
  ON c.id = cs.class_id
ORDER BY s.id ASC
`

	var items []StudentItem
	listErr := s.db.Raw(query, activeClassStudentStatus).Scan(&items).Error
	if listErr != nil {
		return nil, listErr
	}

	return items, nil
}

func (s *Service) Student(rawID string) (StudentItem, bool, error) {
	if s.db == nil {
		return studentItemFromDemo(rawID)
	}

	query := `
SELECT
  s.id,
  s.name,
  s.grade_name AS grade,
  COALESCE(g.name, '') AS parent_name,
  COALESCE(g.mobile, '') AS parent_mobile,
  s.campus,
  COALESCE(cs.class_id, 0) AS class_id,
  COALESCE(c.name, '') AS class_name,
  s.remaining_hours,
  s.status
FROM students AS s
LEFT JOIN student_guardians AS g
  ON g.student_id = s.id AND g.is_primary = 1
LEFT JOIN (
  SELECT student_id, MIN(class_id) AS class_id
  FROM class_students
  WHERE status = ?
  GROUP BY student_id
) AS cs
  ON cs.student_id = s.id
LEFT JOIN classes AS c
  ON c.id = cs.class_id
WHERE s.id = ?
LIMIT 1
`

	var item StudentItem
	findErr := s.db.Raw(query, activeClassStudentStatus, rawID).Scan(&item).Error
	if findErr != nil {
		return StudentItem{}, false, findErr
	}

	if item.ID == 0 {
		return StudentItem{}, false, nil
	}

	return item, true, nil
}

func (s *Service) StudentGuardians(rawStudentID string) ([]StudentGuardianItem, error) {
	if s.db == nil {
		studentItem, found, studentErr := s.Student(rawStudentID)
		if studentErr != nil {
			return nil, studentErr
		}
		if !found || (studentItem.ParentName == "" && studentItem.ParentMobile == "") {
			return []StudentGuardianItem{}, nil
		}

		return []StudentGuardianItem{
			{
				ID:        studentItem.ID,
				Name:      studentItem.ParentName,
				Relation:  guardianRelationFromName(studentItem.ParentName),
				Mobile:    studentItem.ParentMobile,
				IsPrimary: true,
			},
		}, nil
	}

	query := `
SELECT
  id,
  name,
  COALESCE(relation, '') AS relation,
  mobile,
  is_primary
FROM student_guardians
WHERE student_id = ?
ORDER BY is_primary DESC, id ASC
`

	var items []StudentGuardianItem
	listErr := s.db.Raw(query, rawStudentID).Scan(&items).Error
	if listErr != nil {
		return nil, listErr
	}

	return items, nil
}

func (s *Service) CreateStudent(input StudentPayload) (StudentDetail, error) {
	if s.db == nil {
		return StudentDetail{}, nil
	}

	now := time.Now()
	var createdStudentID uint64

	createErr := s.db.Transaction(func(tx *gorm.DB) error {
		student := edumodel.Student{
			Name:           input.Name,
			Gender:         input.Gender,
			SchoolName:     input.SchoolName,
			GradeName:      input.GradeName,
			Campus:         input.Campus,
			RemainingHours: input.RemainingHours,
			Status:         input.Status,
			Remark:         input.Remark,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		createStudentErr := tx.Create(&student).Error
		if createStudentErr != nil {
			return createStudentErr
		}

		guardian := edumodel.StudentGuardian{
			StudentID: student.ID,
			Name:      input.GuardianName,
			Relation:  input.GuardianRelation,
			Mobile:    input.GuardianMobile,
			IsPrimary: true,
			CreatedAt: now,
			UpdatedAt: now,
		}

		createGuardianErr := tx.Create(&guardian).Error
		if createGuardianErr != nil {
			return createGuardianErr
		}

		createdStudentID = student.ID
		return nil
	})
	if createErr != nil {
		return StudentDetail{}, createErr
	}

	detail, found, detailErr := s.StudentDetail(fmt.Sprintf("%d", createdStudentID))
	if detailErr != nil {
		return StudentDetail{}, detailErr
	}
	if !found {
		return StudentDetail{}, nil
	}

	return detail, nil
}

func (s *Service) UpdateStudent(rawID string, input StudentPayload) (StudentDetail, bool, error) {
	if s.db == nil {
		detail, found, detailErr := s.StudentDetail(rawID)
		if detailErr != nil || !found {
			return detail, found, detailErr
		}

		return detail, true, nil
	}

	studentID, parseErr := strconv.ParseUint(rawID, 10, 64)
	if parseErr != nil {
		return StudentDetail{}, false, nil
	}

	updateErr := s.db.Transaction(func(tx *gorm.DB) error {
		updateResult := tx.Model(&edumodel.Student{}).
			Where("id = ?", studentID).
			Updates(map[string]any{
				"name":            input.Name,
				"gender":          input.Gender,
				"school_name":     input.SchoolName,
				"grade_name":      input.GradeName,
				"campus":          input.Campus,
				"remaining_hours": input.RemainingHours,
				"status":          input.Status,
				"remark":          input.Remark,
				"updated_at":      time.Now(),
			})
		if updateResult.Error != nil {
			return updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		var primaryGuardian edumodel.StudentGuardian
		findGuardianErr := tx.
			Where("student_id = ? AND is_primary = ?", studentID, true).
			Order("id ASC").
			Take(&primaryGuardian).Error
		if errors.Is(findGuardianErr, gorm.ErrRecordNotFound) {
			guardian := edumodel.StudentGuardian{
				StudentID: studentID,
				Name:      input.GuardianName,
				Relation:  input.GuardianRelation,
				Mobile:    input.GuardianMobile,
				IsPrimary: true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			return tx.Create(&guardian).Error
		}
		if findGuardianErr != nil {
			return findGuardianErr
		}

		return tx.Model(&edumodel.StudentGuardian{}).
			Where("id = ?", primaryGuardian.ID).
			Updates(map[string]any{
				"name":       input.GuardianName,
				"relation":   input.GuardianRelation,
				"mobile":     input.GuardianMobile,
				"is_primary": true,
				"updated_at": time.Now(),
			}).Error
	})
	if errors.Is(updateErr, gorm.ErrRecordNotFound) {
		return StudentDetail{}, false, nil
	}
	if updateErr != nil {
		return StudentDetail{}, false, updateErr
	}

	detail, found, detailErr := s.StudentDetail(rawID)
	if detailErr != nil {
		return StudentDetail{}, false, detailErr
	}

	return detail, found, nil
}

func (s *Service) CreateStudentGuardian(rawStudentID string, input StudentGuardianPayload) (StudentGuardianItem, bool, error) {
	if s.db == nil {
		return StudentGuardianItem{}, false, nil
	}

	studentID, parseErr := strconv.ParseUint(rawStudentID, 10, 64)
	if parseErr != nil {
		return StudentGuardianItem{}, false, nil
	}

	var createdGuardian edumodel.StudentGuardian
	createErr := s.db.Transaction(func(tx *gorm.DB) error {
		var student edumodel.Student
		findStudentErr := tx.Where("id = ?", studentID).Take(&student).Error
		if findStudentErr != nil {
			return findStudentErr
		}

		if input.IsPrimary {
			resetPrimaryErr := tx.Model(&edumodel.StudentGuardian{}).
				Where("student_id = ?", studentID).
				Update("is_primary", false).Error
			if resetPrimaryErr != nil {
				return resetPrimaryErr
			}
		}

		now := time.Now()
		createdGuardian = edumodel.StudentGuardian{
			StudentID: studentID,
			Name:      input.Name,
			Relation:  input.Relation,
			Mobile:    input.Mobile,
			IsPrimary: input.IsPrimary,
			CreatedAt: now,
			UpdatedAt: now,
		}

		return tx.Create(&createdGuardian).Error
	})
	if errors.Is(createErr, gorm.ErrRecordNotFound) {
		return StudentGuardianItem{}, false, nil
	}
	if createErr != nil {
		return StudentGuardianItem{}, false, createErr
	}

	return StudentGuardianItem{
		ID:        createdGuardian.ID,
		Name:      createdGuardian.Name,
		Relation:  createdGuardian.Relation,
		Mobile:    createdGuardian.Mobile,
		IsPrimary: createdGuardian.IsPrimary,
	}, true, nil
}

func (s *Service) UpdateStudentGuardian(rawStudentID string, rawGuardianID string, input StudentGuardianPayload) (StudentGuardianItem, bool, error) {
	if s.db == nil {
		return StudentGuardianItem{}, false, nil
	}

	studentID, studentParseErr := strconv.ParseUint(rawStudentID, 10, 64)
	if studentParseErr != nil {
		return StudentGuardianItem{}, false, nil
	}

	guardianID, guardianParseErr := strconv.ParseUint(rawGuardianID, 10, 64)
	if guardianParseErr != nil {
		return StudentGuardianItem{}, false, nil
	}

	updateErr := s.db.Transaction(func(tx *gorm.DB) error {
		var guardian edumodel.StudentGuardian
		findGuardianErr := tx.Where("id = ? AND student_id = ?", guardianID, studentID).Take(&guardian).Error
		if findGuardianErr != nil {
			return findGuardianErr
		}

		if input.IsPrimary {
			resetPrimaryErr := tx.Model(&edumodel.StudentGuardian{}).
				Where("student_id = ?", studentID).
				Update("is_primary", false).Error
			if resetPrimaryErr != nil {
				return resetPrimaryErr
			}
		}

		return tx.Model(&edumodel.StudentGuardian{}).
			Where("id = ?", guardianID).
			Updates(map[string]any{
				"name":       input.Name,
				"relation":   input.Relation,
				"mobile":     input.Mobile,
				"is_primary": input.IsPrimary,
				"updated_at": time.Now(),
			}).Error
	})
	if errors.Is(updateErr, gorm.ErrRecordNotFound) {
		return StudentGuardianItem{}, false, nil
	}
	if updateErr != nil {
		return StudentGuardianItem{}, false, updateErr
	}

	guardians, guardianErr := s.StudentGuardians(rawStudentID)
	if guardianErr != nil {
		return StudentGuardianItem{}, false, guardianErr
	}

	for _, guardian := range guardians {
		if guardian.ID == guardianID {
			return guardian, true, nil
		}
	}

	return StudentGuardianItem{}, false, nil
}

func (s *Service) DeleteStudentGuardian(rawStudentID string, rawGuardianID string) (bool, error) {
	if s.db == nil {
		return false, nil
	}

	studentID, studentParseErr := strconv.ParseUint(rawStudentID, 10, 64)
	if studentParseErr != nil {
		return false, nil
	}

	guardianID, guardianParseErr := strconv.ParseUint(rawGuardianID, 10, 64)
	if guardianParseErr != nil {
		return false, nil
	}

	deleteErr := s.db.Transaction(func(tx *gorm.DB) error {
		var guardians []edumodel.StudentGuardian
		listErr := tx.Where("student_id = ?", studentID).Order("id ASC").Find(&guardians).Error
		if listErr != nil {
			return listErr
		}
		if len(guardians) == 0 {
			return gorm.ErrRecordNotFound
		}

		var targetFound bool
		var removedPrimary bool
		var fallbackGuardianID uint64
		for _, guardian := range guardians {
			if guardian.ID == guardianID {
				targetFound = true
				removedPrimary = guardian.IsPrimary
				continue
			}
			if fallbackGuardianID == 0 {
				fallbackGuardianID = guardian.ID
			}
		}
		if !targetFound {
			return gorm.ErrRecordNotFound
		}

		deleteResult := tx.Where("id = ? AND student_id = ?", guardianID, studentID).Delete(&edumodel.StudentGuardian{})
		if deleteResult.Error != nil {
			return deleteResult.Error
		}
		if deleteResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		if removedPrimary && fallbackGuardianID > 0 {
			return tx.Model(&edumodel.StudentGuardian{}).
				Where("id = ?", fallbackGuardianID).
				Update("is_primary", true).Error
		}

		return nil
	})
	if errors.Is(deleteErr, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if deleteErr != nil {
		return false, deleteErr
	}

	return true, nil
}

func (s *Service) StudentDetail(rawStudentID string) (StudentDetail, bool, error) {
	const recentItemLimit = 4

	studentProfile, found, profileErr := s.studentProfile(rawStudentID)
	if profileErr != nil {
		return StudentDetail{}, false, profileErr
	}
	if !found {
		return StudentDetail{}, false, nil
	}

	guardians, guardianErr := s.StudentGuardians(rawStudentID)
	if guardianErr != nil {
		return StudentDetail{}, false, guardianErr
	}

	classes, classErr := s.StudentClasses(rawStudentID)
	if classErr != nil {
		return StudentDetail{}, false, classErr
	}

	classIDs := make(map[uint64]bool, len(classes))
	for _, classItem := range classes {
		classIDs[classItem.ID] = true
	}

	recentSchedules := make([]ScheduleItem, 0, recentItemLimit)
	recentAttendance := make([]StudentAttendanceRecord, 0, recentItemLimit)
	if len(classIDs) > 0 {
		schedules, scheduleErr := s.Schedules()
		if scheduleErr != nil {
			return StudentDetail{}, false, scheduleErr
		}

		for index := len(schedules) - 1; index >= 0; index-- {
			scheduleItem := schedules[index]
			if !classIDs[scheduleItem.ClassID] {
				continue
			}

			recentSchedules = append(recentSchedules, scheduleItem)

			attendanceItems, attendanceErr := s.attendanceFromSchedule(scheduleItem)
			if attendanceErr != nil {
				return StudentDetail{}, false, attendanceErr
			}

			for _, attendanceItem := range attendanceItems {
				if attendanceItem.StudentID != studentProfile.ID {
					continue
				}

				recentAttendance = append(recentAttendance, StudentAttendanceRecord{
					ScheduleID:  scheduleItem.ID,
					ClassID:     scheduleItem.ClassID,
					ClassName:   scheduleItem.ClassName,
					CourseName:  scheduleItem.CourseName,
					TeacherName: scheduleItem.TeacherName,
					Campus:      scheduleItem.Campus,
					Classroom:   scheduleItem.Classroom,
					LessonDate:  scheduleItem.LessonDate,
					LessonTime:  scheduleItem.LessonTime,
					Status:      attendanceItem.Status,
					Remark:      attendanceItem.Remark,
				})
				break
			}

			if len(recentSchedules) >= recentItemLimit {
				break
			}
		}
	}

	homeworks, homeworkErr := s.Homeworks()
	if homeworkErr != nil {
		return StudentDetail{}, false, homeworkErr
	}

	recentHomeworks := make([]HomeworkItem, 0, recentItemLimit)
	for _, item := range homeworks {
		if !classIDs[item.ClassID] {
			continue
		}

		recentHomeworks = append(recentHomeworks, item)
		if len(recentHomeworks) >= recentItemLimit {
			break
		}
	}

	feedbacks, feedbackErr := s.Feedbacks()
	if feedbackErr != nil {
		return StudentDetail{}, false, feedbackErr
	}

	recentFeedbacks := make([]FeedbackItem, 0, recentItemLimit)
	for _, item := range feedbacks {
		if !classIDs[item.ClassID] {
			continue
		}

		recentFeedbacks = append(recentFeedbacks, item)
		if len(recentFeedbacks) >= recentItemLimit {
			break
		}
	}

	return StudentDetail{
		Student:          studentProfile,
		Guardians:        guardians,
		Classes:          classes,
		RecentSchedules:  recentSchedules,
		RecentAttendance: recentAttendance,
		RecentHomeworks:  recentHomeworks,
		RecentFeedbacks:  recentFeedbacks,
	}, true, nil
}

func (s *Service) studentProfile(rawStudentID string) (StudentProfile, bool, error) {
	studentItem, found, studentErr := s.Student(rawStudentID)
	if studentErr != nil {
		return StudentProfile{}, false, studentErr
	}
	if !found {
		return StudentProfile{}, false, nil
	}

	if s.db == nil {
		return studentProfileFromStudentItem(studentItem, "", "", ""), true, nil
	}

	query := `
SELECT
  COALESCE(gender, '') AS gender,
  COALESCE(school_name, '') AS school_name,
  COALESCE(remark, '') AS remark
FROM students
WHERE id = ?
LIMIT 1
`

	type studentProfileRow struct {
		Gender     string `gorm:"column:gender"`
		SchoolName string `gorm:"column:school_name"`
		Remark     string `gorm:"column:remark"`
	}

	var row studentProfileRow
	profileErr := s.db.Raw(query, rawStudentID).Scan(&row).Error
	if profileErr != nil {
		return StudentProfile{}, false, profileErr
	}

	return studentProfileFromStudentItem(studentItem, row.SchoolName, row.Gender, row.Remark), true, nil
}

func (s *Service) StudentClasses(rawStudentID string) ([]ClassItem, error) {
	if s.db == nil {
		return classItemsFromDemo(demo.StudentClasses(rawStudentID)), nil
	}

	query := `
SELECT
  c.id,
  c.name,
  COALESCE(co.name, '') AS course_name,
  c.teacher_id,
  COALESCE(t.name, '') AS teacher_name,
  c.campus,
  COUNT(cs2.id) AS student_count,
  c.capacity,
  c.weekly_schedule,
  c.status
FROM class_students AS cs
JOIN classes AS c
  ON c.id = cs.class_id
LEFT JOIN courses AS co
  ON co.id = c.course_id
LEFT JOIN teachers AS t
  ON t.id = c.teacher_id
LEFT JOIN class_students AS cs2
  ON cs2.class_id = c.id AND cs2.status = ?
WHERE cs.student_id = ? AND cs.status = ?
GROUP BY
  c.id,
  c.name,
  co.name,
  c.teacher_id,
  t.name,
  c.campus,
  c.capacity,
  c.weekly_schedule,
  c.status
ORDER BY c.id ASC
`

	var items []ClassItem
	listErr := s.db.Raw(query, activeClassStudentStatus, rawStudentID, activeClassStudentStatus).Scan(&items).Error
	if listErr != nil {
		return nil, listErr
	}

	return items, nil
}

func (s *Service) Classes() ([]ClassItem, error) {
	if s.db == nil {
		return classItemsFromDemo(demo.Classes()), nil
	}

	query := `
SELECT
  c.id,
  c.course_id,
  c.name,
  COALESCE(co.name, '') AS course_name,
  c.teacher_id,
  COALESCE(t.name, '') AS teacher_name,
  c.campus,
  COUNT(cs.id) AS student_count,
  c.capacity,
  c.weekly_schedule,
  COALESCE(DATE_FORMAT(c.start_date, '%Y-%m-%d'), '') AS start_date,
  COALESCE(DATE_FORMAT(c.end_date, '%Y-%m-%d'), '') AS end_date,
  c.status
  ,
  COALESCE(c.remark, '') AS remark
FROM classes AS c
LEFT JOIN courses AS co
  ON co.id = c.course_id
LEFT JOIN teachers AS t
  ON t.id = c.teacher_id
LEFT JOIN class_students AS cs
  ON cs.class_id = c.id AND cs.status = ?
GROUP BY
  c.id,
  c.course_id,
  c.name,
  co.name,
  c.teacher_id,
  t.name,
  c.campus,
  c.capacity,
  c.weekly_schedule,
  c.start_date,
  c.end_date,
  c.status
  ,
  c.remark
ORDER BY c.id ASC
`

	var items []ClassItem
	listErr := s.db.Raw(query, activeClassStudentStatus).Scan(&items).Error
	if listErr != nil {
		return nil, listErr
	}

	return items, nil
}

func (s *Service) Class(rawID string) (ClassItem, bool, error) {
	if s.db == nil {
		return classItemFromDemo(rawID)
	}

	query := `
SELECT
  c.id,
  c.course_id,
  c.name,
  COALESCE(co.name, '') AS course_name,
  c.teacher_id,
  COALESCE(t.name, '') AS teacher_name,
  c.campus,
  COUNT(cs.id) AS student_count,
  c.capacity,
  c.weekly_schedule,
  COALESCE(DATE_FORMAT(c.start_date, '%Y-%m-%d'), '') AS start_date,
  COALESCE(DATE_FORMAT(c.end_date, '%Y-%m-%d'), '') AS end_date,
  c.status
  ,
  COALESCE(c.remark, '') AS remark
FROM classes AS c
LEFT JOIN courses AS co
  ON co.id = c.course_id
LEFT JOIN teachers AS t
  ON t.id = c.teacher_id
LEFT JOIN class_students AS cs
  ON cs.class_id = c.id AND cs.status = ?
WHERE c.id = ?
GROUP BY
  c.id,
  c.course_id,
  c.name,
  co.name,
  c.teacher_id,
  t.name,
  c.campus,
  c.capacity,
  c.weekly_schedule,
  c.start_date,
  c.end_date,
  c.status
  ,
  c.remark
LIMIT 1
`

	var item ClassItem
	findErr := s.db.Raw(query, activeClassStudentStatus, rawID).Scan(&item).Error
	if findErr != nil {
		return ClassItem{}, false, findErr
	}

	if item.ID == 0 {
		return ClassItem{}, false, nil
	}

	return item, true, nil
}

func (s *Service) CreateClass(input ClassPayload) (ClassItem, error) {
	if s.db == nil {
		return ClassItem{}, nil
	}

	startDate, endDate, parseErr := parseClassDates(input.StartDate, input.EndDate)
	if parseErr != nil {
		return ClassItem{}, parseErr
	}

	now := time.Now()
	classItem := edumodel.Class{
		Name:           input.Name,
		CourseID:       input.CourseID,
		TeacherID:      input.TeacherID,
		Campus:         input.Campus,
		Capacity:       input.Capacity,
		WeeklySchedule: input.WeeklySchedule,
		StartDate:      startDate,
		EndDate:        endDate,
		Status:         input.Status,
		Remark:         input.Remark,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	createErr := s.db.Create(&classItem).Error
	if createErr != nil {
		return ClassItem{}, createErr
	}

	createdItem, found, classErr := s.Class(fmt.Sprintf("%d", classItem.ID))
	if classErr != nil {
		return ClassItem{}, classErr
	}
	if !found {
		return ClassItem{}, nil
	}

	return createdItem, nil
}

func (s *Service) UpdateClass(rawID string, input ClassPayload) (ClassItem, bool, error) {
	if s.db == nil {
		return ClassItem{}, false, nil
	}

	startDate, endDate, parseErr := parseClassDates(input.StartDate, input.EndDate)
	if parseErr != nil {
		return ClassItem{}, false, parseErr
	}

	updateResult := s.db.Model(&edumodel.Class{}).
		Where("id = ?", rawID).
		Updates(map[string]any{
			"name":            input.Name,
			"course_id":       input.CourseID,
			"teacher_id":      input.TeacherID,
			"campus":          input.Campus,
			"capacity":        input.Capacity,
			"weekly_schedule": input.WeeklySchedule,
			"start_date":      startDate,
			"end_date":        endDate,
			"status":          input.Status,
			"remark":          input.Remark,
			"updated_at":      time.Now(),
		})
	if updateResult.Error != nil {
		return ClassItem{}, false, updateResult.Error
	}
	if updateResult.RowsAffected == 0 {
		return ClassItem{}, false, nil
	}

	updatedItem, found, classErr := s.Class(rawID)
	if classErr != nil {
		return ClassItem{}, false, classErr
	}

	return updatedItem, found, nil
}

func (s *Service) ClassStudents(rawClassID string) ([]StudentItem, error) {
	if s.db == nil {
		return studentItemsFromDemoWithSlice(demo.ClassStudents(rawClassID)), nil
	}

	query := `
SELECT
  s.id,
  s.name,
  s.grade_name AS grade,
  COALESCE(g.name, '') AS parent_name,
  COALESCE(g.mobile, '') AS parent_mobile,
  s.campus,
  COALESCE(cs.class_id, 0) AS class_id,
  COALESCE(c.name, '') AS class_name,
  s.remaining_hours,
  s.status
FROM class_students AS cs
JOIN students AS s
  ON s.id = cs.student_id
LEFT JOIN student_guardians AS g
  ON g.student_id = s.id AND g.is_primary = 1
LEFT JOIN classes AS c
  ON c.id = cs.class_id
WHERE cs.class_id = ? AND cs.status = ?
ORDER BY s.id ASC
`

	var items []StudentItem
	listErr := s.db.Raw(query, rawClassID, activeClassStudentStatus).Scan(&items).Error
	if listErr != nil {
		return nil, listErr
	}

	return items, nil
}

func (s *Service) ClassDetail(rawClassID string) (ClassDetail, bool, error) {
	classItem, found, classErr := s.Class(rawClassID)
	if classErr != nil {
		return ClassDetail{}, false, classErr
	}
	if !found {
		return ClassDetail{}, false, nil
	}

	students, studentErr := s.ClassStudents(rawClassID)
	if studentErr != nil {
		return ClassDetail{}, false, studentErr
	}

	upcomingSchedules, scheduleErr := s.UpcomingSchedules(rawClassID)
	if scheduleErr != nil {
		return ClassDetail{}, false, scheduleErr
	}

	recentAttendance := make([]AttendanceSessionItem, 0, len(upcomingSchedules))
	for _, scheduleItem := range upcomingSchedules {
		attendanceItems, attendanceErr := s.attendanceFromSchedule(scheduleItem)
		if attendanceErr != nil {
			return ClassDetail{}, false, attendanceErr
		}

		presentCount, leaveCount, absentCount, pendingCount := summarizeAttendanceItems(attendanceItems)
		recentAttendance = append(recentAttendance, AttendanceSessionItem{
			ID:               scheduleItem.ID,
			ClassID:          scheduleItem.ClassID,
			ClassName:        scheduleItem.ClassName,
			CourseName:       scheduleItem.CourseName,
			TeacherName:      scheduleItem.TeacherName,
			Campus:           scheduleItem.Campus,
			Classroom:        scheduleItem.Classroom,
			LessonDate:       scheduleItem.LessonDate,
			LessonTime:       scheduleItem.LessonTime,
			AttendanceStatus: sessionAttendanceStatus(scheduleItem.AttendanceStatus, pendingCount),
			StudentCount:     len(attendanceItems),
			PresentCount:     presentCount,
			LeaveCount:       leaveCount,
			AbsentCount:      absentCount,
			PendingCount:     pendingCount,
		})
	}

	homeworks, homeworkErr := s.Homeworks()
	if homeworkErr != nil {
		return ClassDetail{}, false, homeworkErr
	}

	recentHomeworks := make([]HomeworkItem, 0, len(homeworks))
	for _, item := range homeworks {
		if item.ClassID == classItem.ID {
			recentHomeworks = append(recentHomeworks, item)
		}
	}

	notices, noticeErr := s.Notices()
	if noticeErr != nil {
		return ClassDetail{}, false, noticeErr
	}

	recentNotices := make([]NoticeItem, 0, len(notices))
	for _, item := range notices {
		if item.RelatedClassID == classItem.ID {
			recentNotices = append(recentNotices, item)
		}
	}

	return ClassDetail{
		Class:             classItem,
		Students:          students,
		UpcomingSchedules: upcomingSchedules,
		RecentAttendance:  recentAttendance,
		RecentHomeworks:   recentHomeworks,
		RecentNotices:     recentNotices,
	}, true, nil
}

func (s *Service) AddStudentsToClass(rawClassID string, studentIDs []uint64) (bool, error) {
	classItem, found, classErr := s.Class(rawClassID)
	if classErr != nil {
		return false, classErr
	}
	if !found {
		return false, nil
	}

	if len(studentIDs) == 0 {
		return true, nil
	}

	if s.db == nil {
		return true, nil
	}

	now := time.Now()
	records := make([]edumodel.ClassStudent, 0, len(studentIDs))
	for _, studentID := range studentIDs {
		records = append(records, edumodel.ClassStudent{
			ClassID:   classItem.ID,
			StudentID: studentID,
			JoinDate:  datePointer(startOfDay(now)),
			Status:    activeClassStudentStatus,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	saveErr := s.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "class_id"},
			{Name: "student_id"},
		},
		DoUpdates: clause.Assignments(map[string]any{
			"status":     activeClassStudentStatus,
			"leave_date": nil,
			"updated_at": now,
		}),
	}).Create(&records).Error
	if saveErr != nil {
		return false, saveErr
	}

	return true, nil
}

func (s *Service) RemoveStudentFromClass(rawClassID string, rawStudentID string) (bool, error) {
	if s.db == nil {
		return true, nil
	}

	now := time.Now()
	updateResult := s.db.Model(&edumodel.ClassStudent{}).
		Where("class_id = ? AND student_id = ?", rawClassID, rawStudentID).
		Updates(map[string]any{
			"status":     "已移出",
			"leave_date": datePointer(startOfDay(now)),
			"updated_at": now,
		})
	if updateResult.Error != nil {
		return false, updateResult.Error
	}
	if updateResult.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (s *Service) ScheduleDetail(rawID string) (ScheduleDetail, bool, error) {
	scheduleItem, found, scheduleErr := s.Schedule(rawID)
	if scheduleErr != nil {
		return ScheduleDetail{}, false, scheduleErr
	}
	if !found {
		return ScheduleDetail{}, false, nil
	}

	classItem, classFound, classErr := s.Class(fmt.Sprintf("%d", scheduleItem.ClassID))
	if classErr != nil {
		return ScheduleDetail{}, false, classErr
	}
	if !classFound {
		return ScheduleDetail{}, false, nil
	}

	students, studentErr := s.ClassStudents(fmt.Sprintf("%d", scheduleItem.ClassID))
	if studentErr != nil {
		return ScheduleDetail{}, false, studentErr
	}

	attendanceItems, attendanceErr := s.attendanceFromSchedule(scheduleItem)
	if attendanceErr != nil {
		return ScheduleDetail{}, false, attendanceErr
	}

	presentCount, leaveCount, absentCount, pendingCount := summarizeAttendanceItems(attendanceItems)
	attendanceSummary := AttendanceSessionItem{
		ID:               scheduleItem.ID,
		ClassID:          scheduleItem.ClassID,
		ClassName:        scheduleItem.ClassName,
		CourseName:       scheduleItem.CourseName,
		TeacherName:      scheduleItem.TeacherName,
		Campus:           scheduleItem.Campus,
		Classroom:        scheduleItem.Classroom,
		LessonDate:       scheduleItem.LessonDate,
		LessonTime:       scheduleItem.LessonTime,
		AttendanceStatus: sessionAttendanceStatus(scheduleItem.AttendanceStatus, pendingCount),
		StudentCount:     len(attendanceItems),
		PresentCount:     presentCount,
		LeaveCount:       leaveCount,
		AbsentCount:      absentCount,
		PendingCount:     pendingCount,
	}

	var homeworkItem *HomeworkItem
	loadedHomework, homeworkFound, homeworkErr := s.Homework(rawID)
	if homeworkErr != nil {
		return ScheduleDetail{}, false, homeworkErr
	}
	if homeworkFound {
		homeworkItem = &loadedHomework
	}

	var feedbackItem *FeedbackItem
	loadedFeedback, feedbackFound, feedbackErr := s.Feedback(rawID)
	if feedbackErr != nil {
		return ScheduleDetail{}, false, feedbackErr
	}
	if feedbackFound {
		feedbackItem = &loadedFeedback
	}

	notices, noticeErr := s.Notices()
	if noticeErr != nil {
		return ScheduleDetail{}, false, noticeErr
	}

	relatedNotices := make([]NoticeItem, 0, len(notices))
	for _, item := range notices {
		if item.RelatedScheduleID == scheduleItem.ID || item.RelatedClassID == scheduleItem.ClassID {
			relatedNotices = append(relatedNotices, item)
		}
	}

	return ScheduleDetail{
		Schedule:       scheduleItem,
		Class:          classItem,
		Students:       students,
		Attendance:     attendanceSummary,
		Homework:       homeworkItem,
		Feedback:       feedbackItem,
		RelatedNotices: relatedNotices,
	}, true, nil
}

func (s *Service) Schedules() ([]ScheduleItem, error) {
	if s.db == nil {
		return scheduleItemsFromDemo(demo.Schedules()), nil
	}

	query := `
SELECT
  s.id,
  s.class_id,
  COALESCE(s.source_schedule_id, 0) AS source_schedule_id,
  COALESCE(c.name, '') AS class_name,
  COALESCE(co.name, '') AS course_name,
  COALESCE(s.teacher_id, 0) AS teacher_id,
  COALESCE(t.name, '') AS teacher_name,
  c.campus,
  COALESCE(s.location, '') AS classroom,
  COALESCE(s.schedule_type, '') AS schedule_type,
  s.schedule_date,
  s.start_time,
  s.end_time,
  s.status AS attendance_status,
  COALESCE(s.remark, '') AS remark
FROM class_schedules AS s
LEFT JOIN classes AS c
  ON c.id = s.class_id
LEFT JOIN courses AS co
  ON co.id = s.course_id
LEFT JOIN teachers AS t
  ON t.id = s.teacher_id
ORDER BY s.schedule_date ASC, s.start_time ASC, s.id ASC
`

	type scheduleRow struct {
		ID               uint64    `gorm:"column:id"`
		ClassID          uint64    `gorm:"column:class_id"`
		SourceScheduleID uint64    `gorm:"column:source_schedule_id"`
		ClassName        string    `gorm:"column:class_name"`
		CourseName       string    `gorm:"column:course_name"`
		TeacherID        uint64    `gorm:"column:teacher_id"`
		TeacherName      string    `gorm:"column:teacher_name"`
		Campus           string    `gorm:"column:campus"`
		Classroom        string    `gorm:"column:classroom"`
		ScheduleType     string    `gorm:"column:schedule_type"`
		ScheduleDate     time.Time `gorm:"column:schedule_date"`
		StartTime        string    `gorm:"column:start_time"`
		EndTime          string    `gorm:"column:end_time"`
		AttendanceStatus string    `gorm:"column:attendance_status"`
		Remark           string    `gorm:"column:remark"`
	}

	var rows []scheduleRow
	listErr := s.db.Raw(query).Scan(&rows).Error
	if listErr != nil {
		return nil, listErr
	}

	items := make([]ScheduleItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, ScheduleItem{
			ID:               row.ID,
			ClassID:          row.ClassID,
			SourceScheduleID: row.SourceScheduleID,
			ClassName:        row.ClassName,
			CourseName:       row.CourseName,
			TeacherID:        row.TeacherID,
			TeacherName:      row.TeacherName,
			Campus:           row.Campus,
			Classroom:        row.Classroom,
			ScheduleType:     row.ScheduleType,
			LessonDate:       row.ScheduleDate.Format(dateLayout),
			LessonTime:       formatLessonTime(row.StartTime, row.EndTime),
			AttendanceStatus: row.AttendanceStatus,
			Remark:           row.Remark,
		})
	}

	return items, nil
}

func (s *Service) Schedule(rawID string) (ScheduleItem, bool, error) {
	items, listErr := s.Schedules()
	if listErr != nil {
		return ScheduleItem{}, false, listErr
	}

	for _, item := range items {
		if fmt.Sprintf("%d", item.ID) == rawID {
			return item, true, nil
		}
	}

	return ScheduleItem{}, false, nil
}

func (s *Service) CreateSchedule(payload SchedulePayload) (ScheduleItem, error) {
	classItem, found, classErr := s.Class(fmt.Sprintf("%d", payload.ClassID))
	if classErr != nil {
		return ScheduleItem{}, classErr
	}
	if !found {
		return ScheduleItem{}, nil
	}

	lessonDate, parseErr := time.ParseInLocation(dateLayout, payload.LessonDate, time.Local)
	if parseErr != nil {
		return ScheduleItem{}, nil
	}

	scheduleType := strings.TrimSpace(payload.ScheduleType)
	if scheduleType == "" {
		scheduleType = "常规课"
	}

	if s.db == nil {
		nextID := scheduleItemsFromDemo(demo.Schedules())
		createdID := uint64(len(nextID) + 1)
		return ScheduleItem{
			ID:               createdID,
			ClassID:          classItem.ID,
			SourceScheduleID: 0,
			ClassName:        classItem.Name,
			CourseName:       classItem.CourseName,
			TeacherID:        classItem.TeacherID,
			TeacherName:      classItem.TeacherName,
			Campus:           classItem.Campus,
			Classroom:        strings.TrimSpace(payload.Classroom),
			ScheduleType:     scheduleType,
			LessonDate:       lessonDate.Format(dateLayout),
			LessonTime:       formatLessonTime(payload.StartTime, payload.EndTime),
			AttendanceStatus: "待上课",
			Remark:           strings.TrimSpace(payload.Remark),
		}, nil
	}

	courseID, courseFound, courseErr := s.classCourseID(payload.ClassID)
	if courseErr != nil {
		return ScheduleItem{}, courseErr
	}
	if !courseFound {
		return ScheduleItem{}, nil
	}

	now := time.Now()
	record := edumodel.ClassSchedule{
		ClassID:          classItem.ID,
		CourseID:         courseID,
		TeacherID:        classItem.TeacherID,
		SourceScheduleID: nil,
		ScheduleType:     scheduleType,
		ScheduleDate:     lessonDate,
		StartTime:        strings.TrimSpace(payload.StartTime),
		EndTime:          strings.TrimSpace(payload.EndTime),
		Location:         strings.TrimSpace(payload.Classroom),
		Status:           "待上课",
		Remark:           strings.TrimSpace(payload.Remark),
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	createErr := s.db.Create(&record).Error
	if createErr != nil {
		return ScheduleItem{}, createErr
	}

	createdItem, itemFound, detailErr := s.Schedule(fmt.Sprintf("%d", record.ID))
	if detailErr != nil {
		return ScheduleItem{}, detailErr
	}
	if !itemFound {
		return ScheduleItem{}, nil
	}

	return createdItem, nil
}

func (s *Service) UpdateSchedule(rawID string, payload SchedulePayload) (ScheduleItem, bool, error) {
	scheduleItem, found, scheduleErr := s.Schedule(rawID)
	if scheduleErr != nil {
		return ScheduleItem{}, false, scheduleErr
	}
	if !found {
		return ScheduleItem{}, false, nil
	}

	classItem, classFound, classErr := s.Class(fmt.Sprintf("%d", payload.ClassID))
	if classErr != nil {
		return ScheduleItem{}, false, classErr
	}
	if !classFound {
		return ScheduleItem{}, false, nil
	}

	lessonDate, parseErr := time.ParseInLocation(dateLayout, payload.LessonDate, time.Local)
	if parseErr != nil {
		return ScheduleItem{}, false, nil
	}

	scheduleType := strings.TrimSpace(payload.ScheduleType)
	if scheduleType == "" {
		scheduleType = scheduleItem.ScheduleType
		if scheduleType == "" {
			scheduleType = "常规课"
		}
	}

	if s.db == nil {
		return ScheduleItem{
			ID:               scheduleItem.ID,
			ClassID:          classItem.ID,
			SourceScheduleID: scheduleItem.SourceScheduleID,
			ClassName:        classItem.Name,
			CourseName:       classItem.CourseName,
			TeacherID:        classItem.TeacherID,
			TeacherName:      classItem.TeacherName,
			Campus:           classItem.Campus,
			Classroom:        strings.TrimSpace(payload.Classroom),
			ScheduleType:     scheduleType,
			LessonDate:       lessonDate.Format(dateLayout),
			LessonTime:       formatLessonTime(payload.StartTime, payload.EndTime),
			AttendanceStatus: scheduleItem.AttendanceStatus,
			Remark:           strings.TrimSpace(payload.Remark),
		}, true, nil
	}

	courseID, courseFound, courseErr := s.classCourseID(payload.ClassID)
	if courseErr != nil {
		return ScheduleItem{}, false, courseErr
	}
	if !courseFound {
		return ScheduleItem{}, false, nil
	}

	updateValues := map[string]any{
		"class_id":      classItem.ID,
		"course_id":     courseID,
		"teacher_id":    classItem.TeacherID,
		"schedule_type": scheduleType,
		"schedule_date": lessonDate,
		"start_time":    strings.TrimSpace(payload.StartTime),
		"end_time":      strings.TrimSpace(payload.EndTime),
		"location":      strings.TrimSpace(payload.Classroom),
		"remark":        strings.TrimSpace(payload.Remark),
		"updated_at":    time.Now(),
	}

	updateResult := s.db.Model(&edumodel.ClassSchedule{}).
		Where("id = ?", rawID).
		Updates(updateValues)
	if updateResult.Error != nil {
		return ScheduleItem{}, false, updateResult.Error
	}
	if updateResult.RowsAffected == 0 {
		return ScheduleItem{}, false, nil
	}

	updatedItem, itemFound, detailErr := s.Schedule(rawID)
	if detailErr != nil {
		return ScheduleItem{}, false, detailErr
	}

	return updatedItem, itemFound, nil
}

func (s *Service) Reschedule(rawID string, payload ScheduleActionPayload) (ScheduleItem, bool, error) {
	scheduleItem, found, scheduleErr := s.Schedule(rawID)
	if scheduleErr != nil {
		return ScheduleItem{}, false, scheduleErr
	}
	if !found {
		return ScheduleItem{}, false, nil
	}

	if s.db == nil {
		scheduleItem.AttendanceStatus = "已调课"
		return scheduleItem, true, nil
	}

	replacementPayload := SchedulePayload{
		ClassID:      scheduleItem.ClassID,
		ScheduleType: "调课",
		LessonDate:   payload.LessonDate,
		StartTime:    payload.StartTime,
		EndTime:      payload.EndTime,
		Classroom:    payload.Classroom,
		Remark:       strings.TrimSpace(payload.Remark),
	}

	replacementItem, createErr := s.CreateSchedule(replacementPayload)
	if createErr != nil {
		return ScheduleItem{}, false, createErr
	}
	if replacementItem.ID == 0 {
		return ScheduleItem{}, false, nil
	}

	now := time.Now()
	transactionErr := s.db.Transaction(func(tx *gorm.DB) error {
		sourceScheduleID := scheduleItem.ID
		updateReplacementErr := tx.Model(&edumodel.ClassSchedule{}).
			Where("id = ?", replacementItem.ID).
			Updates(map[string]any{
				"source_schedule_id": sourceScheduleID,
				"status":             "待上课",
				"updated_at":         now,
			}).Error
		if updateReplacementErr != nil {
			return updateReplacementErr
		}

		updateOriginalErr := tx.Model(&edumodel.ClassSchedule{}).
			Where("id = ?", rawID).
			Updates(map[string]any{
				"status":     "已调课",
				"updated_at": now,
			}).Error
		if updateOriginalErr != nil {
			return updateOriginalErr
		}

		return nil
	})
	if transactionErr != nil {
		return ScheduleItem{}, false, transactionErr
	}

	noticeErr := s.createScheduleNotice(rawID, NoticePayload{
		Title:             fmt.Sprintf("%s调课通知", scheduleItem.ClassName),
		Content:           buildRescheduleNoticeContent(scheduleItem, replacementItem, strings.TrimSpace(payload.Remark)),
		Category:          "调课通知",
		TargetScope:       buildClassTargetScope(scheduleItem.ClassName),
		RelatedClassID:    scheduleItem.ClassID,
		RelatedScheduleID: replacementItem.ID,
		Status:            "草稿",
		Author:            "教务老师",
	})
	if noticeErr != nil {
		return ScheduleItem{}, false, noticeErr
	}

	return s.Schedule(rawID)
}

func (s *Service) CancelSchedule(rawID string, payload ScheduleActionPayload) (ScheduleItem, bool, error) {
	scheduleItem, found, scheduleErr := s.Schedule(rawID)
	if scheduleErr != nil {
		return ScheduleItem{}, false, scheduleErr
	}
	if !found {
		return ScheduleItem{}, false, nil
	}

	if s.db == nil {
		scheduleItem.AttendanceStatus = "已停课"
		scheduleItem.Remark = strings.TrimSpace(payload.Remark)
		return scheduleItem, true, nil
	}

	updateResult := s.db.Model(&edumodel.ClassSchedule{}).
		Where("id = ?", rawID).
		Updates(map[string]any{
			"status":     "已停课",
			"remark":     strings.TrimSpace(payload.Remark),
			"updated_at": time.Now(),
		})
	if updateResult.Error != nil {
		return ScheduleItem{}, false, updateResult.Error
	}
	if updateResult.RowsAffected == 0 {
		return ScheduleItem{}, false, nil
	}

	noticeErr := s.createScheduleNotice(rawID, NoticePayload{
		Title:             fmt.Sprintf("%s停课通知", scheduleItem.ClassName),
		Content:           buildCancelNoticeContent(scheduleItem, strings.TrimSpace(payload.Remark)),
		Category:          "停课提醒",
		TargetScope:       buildClassTargetScope(scheduleItem.ClassName),
		RelatedClassID:    scheduleItem.ClassID,
		RelatedScheduleID: scheduleItem.ID,
		Status:            "草稿",
		Author:            "教务老师",
	})
	if noticeErr != nil {
		return ScheduleItem{}, false, noticeErr
	}

	return s.Schedule(rawID)
}

func (s *Service) CreateMakeupSchedule(rawID string, payload ScheduleActionPayload) (ScheduleItem, bool, error) {
	originalItem, found, originalErr := s.Schedule(rawID)
	if originalErr != nil {
		return ScheduleItem{}, false, originalErr
	}
	if !found {
		return ScheduleItem{}, false, nil
	}

	makeupPayload := SchedulePayload{
		ClassID:      originalItem.ClassID,
		ScheduleType: "补课",
		LessonDate:   payload.LessonDate,
		StartTime:    payload.StartTime,
		EndTime:      payload.EndTime,
		Classroom:    payload.Classroom,
		Remark:       payload.Remark,
	}

	createdItem, createErr := s.CreateSchedule(makeupPayload)
	if createErr != nil {
		return ScheduleItem{}, false, createErr
	}
	if createdItem.ID == 0 {
		return ScheduleItem{}, false, nil
	}

	if s.db == nil {
		createdItem.AttendanceStatus = "待上课"
		createdItem.SourceScheduleID = originalItem.ID
		return createdItem, true, nil
	}

	sourceScheduleID := originalItem.ID
	sourceErr := s.db.Model(&edumodel.ClassSchedule{}).
		Where("id = ?", createdItem.ID).
		Updates(map[string]any{
			"source_schedule_id": sourceScheduleID,
			"status":             "待上课",
			"updated_at":         time.Now(),
		}).Error
	if sourceErr != nil {
		return ScheduleItem{}, false, sourceErr
	}

	noticeErr := s.createScheduleNotice(rawID, NoticePayload{
		Title:             fmt.Sprintf("%s补课通知", originalItem.ClassName),
		Content:           buildMakeupNoticeContent(originalItem, createdItem, strings.TrimSpace(payload.Remark)),
		Category:          "调课通知",
		TargetScope:       buildClassTargetScope(originalItem.ClassName),
		RelatedClassID:    originalItem.ClassID,
		RelatedScheduleID: createdItem.ID,
		Status:            "草稿",
		Author:            "教务老师",
	})
	if noticeErr != nil {
		return ScheduleItem{}, false, noticeErr
	}

	return s.Schedule(fmt.Sprintf("%d", createdItem.ID))
}

func (s *Service) classCourseID(classID uint64) (uint64, bool, error) {
	if s.db == nil {
		return classID, true, nil
	}

	type classCourseRow struct {
		CourseID uint64 `gorm:"column:course_id"`
	}

	var row classCourseRow
	findErr := s.db.Model(&edumodel.Class{}).
		Select("course_id").
		Where("id = ?", classID).
		Scan(&row).Error
	if findErr != nil {
		return 0, false, findErr
	}
	if row.CourseID == 0 {
		return 0, false, nil
	}

	return row.CourseID, true, nil
}

func (s *Service) UpcomingSchedules(rawClassID string) ([]ScheduleItem, error) {
	items, listErr := s.Schedules()
	if listErr != nil {
		return nil, listErr
	}

	filteredItems := make([]ScheduleItem, 0, len(items))
	for _, item := range items {
		if fmt.Sprintf("%d", item.ClassID) == rawClassID {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems, nil
}

func (s *Service) Attendance(rawScheduleID string) ([]AttendanceItem, error) {
	if s.db == nil {
		return attendanceItemsFromDemo(demo.Attendance(rawScheduleID)), nil
	}

	scheduleItem, found, scheduleErr := s.Schedule(rawScheduleID)
	if scheduleErr != nil {
		return nil, scheduleErr
	}
	if !found {
		return []AttendanceItem{}, nil
	}

	return s.attendanceFromSchedule(scheduleItem)
}

func (s *Service) AttendanceSessions() ([]AttendanceSessionItem, error) {
	schedules, scheduleErr := s.Schedules()
	if scheduleErr != nil {
		return nil, scheduleErr
	}

	items := make([]AttendanceSessionItem, 0, len(schedules))
	for _, scheduleItem := range schedules {
		attendanceItems, attendanceErr := s.attendanceFromSchedule(scheduleItem)
		if attendanceErr != nil {
			return nil, attendanceErr
		}

		presentCount, leaveCount, absentCount, pendingCount := summarizeAttendanceItems(attendanceItems)

		items = append(items, AttendanceSessionItem{
			ID:               scheduleItem.ID,
			ClassID:          scheduleItem.ClassID,
			ClassName:        scheduleItem.ClassName,
			CourseName:       scheduleItem.CourseName,
			TeacherName:      scheduleItem.TeacherName,
			Campus:           scheduleItem.Campus,
			Classroom:        scheduleItem.Classroom,
			LessonDate:       scheduleItem.LessonDate,
			LessonTime:       scheduleItem.LessonTime,
			AttendanceStatus: sessionAttendanceStatus(scheduleItem.AttendanceStatus, pendingCount),
			StudentCount:     len(attendanceItems),
			PresentCount:     presentCount,
			LeaveCount:       leaveCount,
			AbsentCount:      absentCount,
			PendingCount:     pendingCount,
		})
	}

	return items, nil
}

func (s *Service) AttendanceRecords(filter AttendanceRecordFilter) ([]AttendanceRecordItem, error) {
	schedules, scheduleErr := s.Schedules()
	if scheduleErr != nil {
		return nil, scheduleErr
	}

	filteredStatus := normalizeAttendanceItemStatus(filter.Status)
	filterDate := strings.TrimSpace(filter.Date)

	items := make([]AttendanceRecordItem, 0)
	for _, scheduleItem := range schedules {
		if filter.ClassID > 0 && scheduleItem.ClassID != filter.ClassID {
			continue
		}
		if filterDate != "" && scheduleItem.LessonDate != filterDate {
			continue
		}

		attendanceItems, attendanceErr := s.attendanceFromSchedule(scheduleItem)
		if attendanceErr != nil {
			return nil, attendanceErr
		}

		for _, attendanceItem := range attendanceItems {
			if filter.StudentID > 0 && attendanceItem.StudentID != filter.StudentID {
				continue
			}
			if filteredStatus != "" && attendanceItem.Status != filteredStatus {
				continue
			}

			items = append(items, AttendanceRecordItem{
				ScheduleID:   scheduleItem.ID,
				ClassID:      scheduleItem.ClassID,
				ClassName:    scheduleItem.ClassName,
				StudentID:    attendanceItem.StudentID,
				StudentName:  attendanceItem.StudentName,
				TeacherName:  scheduleItem.TeacherName,
				LessonDate:   scheduleItem.LessonDate,
				LessonTime:   scheduleItem.LessonTime,
				Status:       attendanceItem.Status,
				Remark:       attendanceItem.Remark,
				UpdatedBy:    attendanceItem.UpdatedBy,
				UpdatedAt:    attendanceItem.UpdatedAt,
				ParentMobile: attendanceItem.ParentMobile,
			})
		}
	}

	return items, nil
}

func (s *Service) SaveAttendance(rawScheduleID string, payload AttendanceSavePayload, updatedBy string) (bool, error) {
	scheduleItem, found, scheduleErr := s.Schedule(rawScheduleID)
	if scheduleErr != nil {
		return false, scheduleErr
	}
	if !found {
		return false, nil
	}

	currentItems, currentErr := s.attendanceFromSchedule(scheduleItem)
	if currentErr != nil {
		return false, currentErr
	}

	payloadMap := make(map[uint64]AttendanceSaveItem, len(payload.Items))
	for _, item := range payload.Items {
		payloadMap[item.StudentID] = AttendanceSaveItem{
			StudentID: item.StudentID,
			Status:    normalizeAttendanceItemStatus(item.Status),
			Remark:    strings.TrimSpace(item.Remark),
		}
	}

	finalItems := make([]AttendanceItem, 0, len(currentItems))
	for _, item := range currentItems {
		nextItem := AttendanceItem{
			StudentID:    item.StudentID,
			StudentName:  item.StudentName,
			Grade:        item.Grade,
			ParentMobile: item.ParentMobile,
			Status:       normalizeAttendanceItemStatus(item.Status),
			Remark:       strings.TrimSpace(item.Remark),
		}

		if payloadItem, found := payloadMap[item.StudentID]; found {
			nextItem.Status = normalizeAttendanceItemStatus(payloadItem.Status)
			nextItem.Remark = strings.TrimSpace(payloadItem.Remark)
		}

		finalItems = append(finalItems, nextItem)
	}

	_, _, _, pendingCount := summarizeAttendanceItems(finalItems)
	nextStatus := nextSavedAttendanceStatus(pendingCount)
	now := time.Now()

	if s.db == nil {
		demoItems := make([]demo.AttendanceItem, 0, len(finalItems))
		for _, item := range finalItems {
			demoItems = append(demoItems, demo.AttendanceItem{
				StudentID:    int(item.StudentID),
				StudentName:  item.StudentName,
				Grade:        item.Grade,
				ParentMobile: item.ParentMobile,
				Status:       item.Status,
				Remark:       item.Remark,
				UpdatedBy:    strings.TrimSpace(updatedBy),
				UpdatedAt:    now.Format(dateTimeLayout),
			})
		}

		return demo.SaveAttendance(rawScheduleID, demoItems), nil
	}

	scheduleID, parseErr := strconv.ParseUint(rawScheduleID, 10, 64)
	if parseErr != nil {
		return false, nil
	}

	records := make([]edumodel.AttendanceRecord, 0, len(finalItems))
	editorName := strings.TrimSpace(updatedBy)
	if editorName == "" {
		editorName = "系统管理员"
	}
	for _, item := range finalItems {
		checkedAt := now
		records = append(records, edumodel.AttendanceRecord{
			ScheduleID: scheduleID,
			StudentID:  item.StudentID,
			Status:     item.Status,
			Remark:     item.Remark,
			CheckedAt:  &checkedAt,
			UpdatedBy:  editorName,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
	}

	transactionErr := s.db.Transaction(func(tx *gorm.DB) error {
		if len(records) > 0 {
			createErr := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "schedule_id"},
					{Name: "student_id"},
				},
				DoUpdates: clause.AssignmentColumns([]string{
					"status",
					"remark",
					"checked_at",
					"updated_by",
					"updated_at",
				}),
			}).Create(&records).Error
			if createErr != nil {
				return createErr
			}
		}

		updateErr := tx.Model(&edumodel.ClassSchedule{}).
			Where("id = ?", scheduleID).
			Updates(map[string]any{
				"status":     nextStatus,
				"updated_at": now,
			}).Error
		if updateErr != nil {
			return updateErr
		}

		return nil
	})
	if transactionErr != nil {
		return false, transactionErr
	}

	return true, nil
}

func (s *Service) attendanceFromSchedule(scheduleItem ScheduleItem) ([]AttendanceItem, error) {
	if s.db == nil {
		return attendanceItemsFromDemo(demo.Attendance(fmt.Sprintf("%d", scheduleItem.ID))), nil
	}

	query := `
SELECT
  st.id AS student_id,
  st.name AS student_name,
  st.grade_name AS grade,
  COALESCE(g.mobile, '') AS parent_mobile,
  COALESCE(ar.status, '') AS attendance_status,
  COALESCE(ar.remark, '') AS attendance_remark,
  COALESCE(ar.updated_by, '') AS updated_by,
  ar.checked_at AS checked_at
FROM class_students AS cs
JOIN students AS st
  ON st.id = cs.student_id
LEFT JOIN student_guardians AS g
  ON g.student_id = st.id AND g.is_primary = 1
LEFT JOIN attendance_records AS ar
  ON ar.schedule_id = ? AND ar.student_id = st.id
WHERE cs.class_id = ? AND cs.status = ?
ORDER BY st.id ASC
`

	type attendanceRow struct {
		StudentID        uint64     `gorm:"column:student_id"`
		StudentName      string     `gorm:"column:student_name"`
		Grade            string     `gorm:"column:grade"`
		ParentMobile     string     `gorm:"column:parent_mobile"`
		AttendanceStatus string     `gorm:"column:attendance_status"`
		AttendanceRemark string     `gorm:"column:attendance_remark"`
		UpdatedBy        string     `gorm:"column:updated_by"`
		CheckedAt        *time.Time `gorm:"column:checked_at"`
	}

	var rows []attendanceRow
	listErr := s.db.Raw(query, scheduleItem.ID, scheduleItem.ClassID, activeClassStudentStatus).Scan(&rows).Error
	if listErr != nil {
		return nil, listErr
	}

	items := make([]AttendanceItem, 0, len(rows))
	for index, row := range rows {
		status := normalizeAttendanceItemStatus(row.AttendanceStatus)
		if status == "" {
			status = defaultAttendanceItemStatus(scheduleItem.AttendanceStatus, index, len(rows))
		}

		items = append(items, AttendanceItem{
			StudentID:    row.StudentID,
			StudentName:  row.StudentName,
			Grade:        row.Grade,
			ParentMobile: row.ParentMobile,
			Status:       status,
			Remark:       row.AttendanceRemark,
			UpdatedBy:    row.UpdatedBy,
			UpdatedAt:    formatDateTime(row.CheckedAt),
		})
	}

	return items, nil
}

func (s *Service) Homeworks() ([]HomeworkItem, error) {
	if s.db == nil {
		return homeworkItemsFromDemo(demo.Homeworks()), nil
	}

	schedules, scheduleErr := s.Schedules()
	if scheduleErr != nil {
		return nil, scheduleErr
	}

	query := `
SELECT
  h.id,
  h.schedule_id,
  h.class_id,
  COALESCE(c.name, '') AS class_name,
  COALESCE(co.name, '') AS course_name,
  COALESCE(t.name, '') AS teacher_name,
  s.schedule_date,
  h.title,
  h.content,
  COALESCE(h.submission_note, '') AS submission_note,
  h.status
FROM homeworks AS h
LEFT JOIN class_schedules AS s
  ON s.id = h.schedule_id
LEFT JOIN classes AS c
  ON c.id = h.class_id
LEFT JOIN courses AS co
  ON co.id = s.course_id
LEFT JOIN teachers AS t
  ON t.id = s.teacher_id
ORDER BY s.schedule_date DESC, h.id DESC
`

	type homeworkRow struct {
		ID             uint64    `gorm:"column:id"`
		ScheduleID     uint64    `gorm:"column:schedule_id"`
		ClassID        uint64    `gorm:"column:class_id"`
		ClassName      string    `gorm:"column:class_name"`
		CourseName     string    `gorm:"column:course_name"`
		TeacherName    string    `gorm:"column:teacher_name"`
		ScheduleDate   time.Time `gorm:"column:schedule_date"`
		Title          string    `gorm:"column:title"`
		Content        string    `gorm:"column:content"`
		SubmissionNote string    `gorm:"column:submission_note"`
		Status         string    `gorm:"column:status"`
	}

	var rows []homeworkRow
	listErr := s.db.Raw(query).Scan(&rows).Error
	if listErr != nil {
		return nil, listErr
	}

	itemMap := make(map[uint64]HomeworkItem, len(rows))
	for _, row := range rows {
		itemMap[row.ScheduleID] = HomeworkItem{
			ID:             row.ID,
			ScheduleID:     row.ScheduleID,
			ClassID:        row.ClassID,
			ClassName:      row.ClassName,
			CourseName:     row.CourseName,
			TeacherName:    row.TeacherName,
			LessonDate:     row.ScheduleDate.Format(dateLayout),
			Title:          row.Title,
			Content:        row.Content,
			SubmissionNote: row.SubmissionNote,
			Status:         row.Status,
		}
	}

	items := make([]HomeworkItem, 0, len(schedules))
	for _, scheduleItem := range schedules {
		if item, found := itemMap[scheduleItem.ID]; found {
			items = append(items, item)
			continue
		}

		fallbackItem, fallbackFound := homeworkItemFromDemoWithSchedule(scheduleItem)
		if fallbackFound {
			items = append(items, fallbackItem)
		}
	}

	return items, nil
}

func (s *Service) Homework(rawScheduleID string) (HomeworkItem, bool, error) {
	if s.db == nil {
		return homeworkItemFromDemo(rawScheduleID)
	}

	scheduleItem, found, scheduleErr := s.Schedule(rawScheduleID)
	if scheduleErr != nil {
		return HomeworkItem{}, false, scheduleErr
	}
	if !found {
		return HomeworkItem{}, false, nil
	}

	query := `
SELECT
  h.id,
  h.schedule_id,
  h.class_id,
  COALESCE(c.name, '') AS class_name,
  COALESCE(co.name, '') AS course_name,
  COALESCE(t.name, '') AS teacher_name,
  s.schedule_date,
  h.title,
  h.content,
  COALESCE(h.submission_note, '') AS submission_note,
  h.status
FROM homeworks AS h
LEFT JOIN class_schedules AS s
  ON s.id = h.schedule_id
LEFT JOIN classes AS c
  ON c.id = h.class_id
LEFT JOIN courses AS co
  ON co.id = s.course_id
LEFT JOIN teachers AS t
  ON t.id = s.teacher_id
WHERE h.schedule_id = ?
LIMIT 1
`

	type homeworkRow struct {
		ID             uint64    `gorm:"column:id"`
		ScheduleID     uint64    `gorm:"column:schedule_id"`
		ClassID        uint64    `gorm:"column:class_id"`
		ClassName      string    `gorm:"column:class_name"`
		CourseName     string    `gorm:"column:course_name"`
		TeacherName    string    `gorm:"column:teacher_name"`
		ScheduleDate   time.Time `gorm:"column:schedule_date"`
		Title          string    `gorm:"column:title"`
		Content        string    `gorm:"column:content"`
		SubmissionNote string    `gorm:"column:submission_note"`
		Status         string    `gorm:"column:status"`
	}

	var row homeworkRow
	findErr := s.db.Raw(query, rawScheduleID).Scan(&row).Error
	if findErr != nil {
		return HomeworkItem{}, false, findErr
	}

	if row.ScheduleID == 0 {
		fallbackItem, fallbackFound := homeworkItemFromDemoWithSchedule(scheduleItem)
		if !fallbackFound {
			return HomeworkItem{}, false, nil
		}

		return fallbackItem, true, nil
	}

	return HomeworkItem{
		ID:             row.ID,
		ScheduleID:     row.ScheduleID,
		ClassID:        row.ClassID,
		ClassName:      row.ClassName,
		CourseName:     row.CourseName,
		TeacherName:    row.TeacherName,
		LessonDate:     row.ScheduleDate.Format(dateLayout),
		Title:          row.Title,
		Content:        row.Content,
		SubmissionNote: row.SubmissionNote,
		Status:         row.Status,
	}, true, nil
}

func (s *Service) SaveHomework(rawScheduleID string, payload HomeworkPayload) (HomeworkItem, bool, error) {
	scheduleItem, found, scheduleErr := s.Schedule(rawScheduleID)
	if scheduleErr != nil {
		return HomeworkItem{}, false, scheduleErr
	}
	if !found {
		return HomeworkItem{}, false, nil
	}

	normalizedStatus := normalizeHomeworkStatus(payload.Status)
	if normalizedStatus == "" {
		normalizedStatus = "published"
	}

	trimmedTitle := strings.TrimSpace(payload.Title)
	if trimmedTitle == "" {
		trimmedTitle = fmt.Sprintf("%s课后作业", scheduleItem.CourseName)
	}

	if s.db == nil {
		savedItem := HomeworkItem{
			ID:             uint64(scheduleItem.ID),
			ScheduleID:     scheduleItem.ID,
			ClassID:        scheduleItem.ClassID,
			ClassName:      scheduleItem.ClassName,
			CourseName:     scheduleItem.CourseName,
			TeacherName:    scheduleItem.TeacherName,
			LessonDate:     scheduleItem.LessonDate,
			Title:          trimmedTitle,
			Content:        strings.TrimSpace(payload.Content),
			SubmissionNote: strings.TrimSpace(payload.SubmissionNote),
			Status:         normalizedStatus,
		}

		demoSaved := demo.SaveHomework(demo.Homework{
			ID:             int(savedItem.ID),
			ScheduleID:     int(savedItem.ScheduleID),
			ClassID:        int(savedItem.ClassID),
			ClassName:      savedItem.ClassName,
			CourseName:     savedItem.CourseName,
			TeacherName:    savedItem.TeacherName,
			LessonDate:     savedItem.LessonDate,
			Title:          savedItem.Title,
			Content:        savedItem.Content,
			SubmissionNote: savedItem.SubmissionNote,
			Status:         savedItem.Status,
		})
		if !demoSaved {
			return HomeworkItem{}, false, nil
		}

		return savedItem, true, nil
	}

	scheduleID, parseErr := strconv.ParseUint(rawScheduleID, 10, 64)
	if parseErr != nil {
		return HomeworkItem{}, false, nil
	}

	now := time.Now()
	record := edumodel.Homework{
		ScheduleID:      scheduleID,
		ClassID:         scheduleItem.ClassID,
		Title:           trimmedTitle,
		Content:         strings.TrimSpace(payload.Content),
		SubmissionNote:  strings.TrimSpace(payload.SubmissionNote),
		CreatedByUserID: 1,
		Status:          normalizedStatus,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	saveErr := s.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "schedule_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"title",
			"content",
			"submission_note",
			"status",
			"updated_at",
		}),
	}).Create(&record).Error
	if saveErr != nil {
		return HomeworkItem{}, false, saveErr
	}

	savedItem, detailFound, detailErr := s.Homework(rawScheduleID)
	if detailErr != nil {
		return HomeworkItem{}, false, detailErr
	}

	return savedItem, detailFound, nil
}

func (s *Service) Feedbacks() ([]FeedbackItem, error) {
	if s.db == nil {
		return feedbackItemsFromDemo(demo.Feedbacks()), nil
	}

	schedules, scheduleErr := s.Schedules()
	if scheduleErr != nil {
		return nil, scheduleErr
	}

	query := `
SELECT
  f.id,
  f.schedule_id,
  f.class_id,
  COALESCE(c.name, '') AS class_name,
  COALESCE(co.name, '') AS course_name,
  COALESCE(t.name, '') AS teacher_name,
  s.schedule_date,
  COALESCE(f.summary, '') AS summary,
  COALESCE(f.learning_status, '') AS learning_status,
  COALESCE(f.next_suggestion, '') AS next_suggestion,
  COALESCE(f.parent_notice, '') AS parent_notice
FROM class_feedbacks AS f
LEFT JOIN class_schedules AS s
  ON s.id = f.schedule_id
LEFT JOIN classes AS c
  ON c.id = f.class_id
LEFT JOIN courses AS co
  ON co.id = s.course_id
LEFT JOIN teachers AS t
  ON t.id = s.teacher_id
ORDER BY s.schedule_date DESC, f.id DESC
`

	type feedbackRow struct {
		ID             uint64    `gorm:"column:id"`
		ScheduleID     uint64    `gorm:"column:schedule_id"`
		ClassID        uint64    `gorm:"column:class_id"`
		ClassName      string    `gorm:"column:class_name"`
		CourseName     string    `gorm:"column:course_name"`
		TeacherName    string    `gorm:"column:teacher_name"`
		ScheduleDate   time.Time `gorm:"column:schedule_date"`
		Summary        string    `gorm:"column:summary"`
		LearningStatus string    `gorm:"column:learning_status"`
		NextSuggestion string    `gorm:"column:next_suggestion"`
		ParentNotice   string    `gorm:"column:parent_notice"`
	}

	var rows []feedbackRow
	listErr := s.db.Raw(query).Scan(&rows).Error
	if listErr != nil {
		return nil, listErr
	}

	itemMap := make(map[uint64]FeedbackItem, len(rows))
	for _, row := range rows {
		itemMap[row.ScheduleID] = FeedbackItem{
			ID:             row.ID,
			ScheduleID:     row.ScheduleID,
			ClassID:        row.ClassID,
			ClassName:      row.ClassName,
			CourseName:     row.CourseName,
			TeacherName:    row.TeacherName,
			LessonDate:     row.ScheduleDate.Format(dateLayout),
			Summary:        row.Summary,
			LearningStatus: row.LearningStatus,
			NextSuggestion: row.NextSuggestion,
			ParentNotice:   row.ParentNotice,
		}
	}

	items := make([]FeedbackItem, 0, len(schedules))
	for _, scheduleItem := range schedules {
		if item, found := itemMap[scheduleItem.ID]; found {
			items = append(items, item)
			continue
		}

		fallbackItem, fallbackFound := feedbackItemFromDemoWithSchedule(scheduleItem)
		if fallbackFound {
			items = append(items, fallbackItem)
		}
	}

	return items, nil
}

func (s *Service) Feedback(rawScheduleID string) (FeedbackItem, bool, error) {
	if s.db == nil {
		return feedbackItemFromDemo(rawScheduleID)
	}

	scheduleItem, found, scheduleErr := s.Schedule(rawScheduleID)
	if scheduleErr != nil {
		return FeedbackItem{}, false, scheduleErr
	}
	if !found {
		return FeedbackItem{}, false, nil
	}

	query := `
SELECT
  f.id,
  f.schedule_id,
  f.class_id,
  COALESCE(c.name, '') AS class_name,
  COALESCE(co.name, '') AS course_name,
  COALESCE(t.name, '') AS teacher_name,
  s.schedule_date,
  COALESCE(f.summary, '') AS summary,
  COALESCE(f.learning_status, '') AS learning_status,
  COALESCE(f.next_suggestion, '') AS next_suggestion,
  COALESCE(f.parent_notice, '') AS parent_notice
FROM class_feedbacks AS f
LEFT JOIN class_schedules AS s
  ON s.id = f.schedule_id
LEFT JOIN classes AS c
  ON c.id = f.class_id
LEFT JOIN courses AS co
  ON co.id = s.course_id
LEFT JOIN teachers AS t
  ON t.id = s.teacher_id
WHERE f.schedule_id = ?
LIMIT 1
`

	type feedbackRow struct {
		ID             uint64    `gorm:"column:id"`
		ScheduleID     uint64    `gorm:"column:schedule_id"`
		ClassID        uint64    `gorm:"column:class_id"`
		ClassName      string    `gorm:"column:class_name"`
		CourseName     string    `gorm:"column:course_name"`
		TeacherName    string    `gorm:"column:teacher_name"`
		ScheduleDate   time.Time `gorm:"column:schedule_date"`
		Summary        string    `gorm:"column:summary"`
		LearningStatus string    `gorm:"column:learning_status"`
		NextSuggestion string    `gorm:"column:next_suggestion"`
		ParentNotice   string    `gorm:"column:parent_notice"`
	}

	var row feedbackRow
	findErr := s.db.Raw(query, rawScheduleID).Scan(&row).Error
	if findErr != nil {
		return FeedbackItem{}, false, findErr
	}

	if row.ScheduleID == 0 {
		fallbackItem, fallbackFound := feedbackItemFromDemoWithSchedule(scheduleItem)
		if !fallbackFound {
			return FeedbackItem{}, false, nil
		}

		return fallbackItem, true, nil
	}

	return FeedbackItem{
		ID:             row.ID,
		ScheduleID:     row.ScheduleID,
		ClassID:        row.ClassID,
		ClassName:      row.ClassName,
		CourseName:     row.CourseName,
		TeacherName:    row.TeacherName,
		LessonDate:     row.ScheduleDate.Format(dateLayout),
		Summary:        row.Summary,
		LearningStatus: row.LearningStatus,
		NextSuggestion: row.NextSuggestion,
		ParentNotice:   row.ParentNotice,
	}, true, nil
}

func (s *Service) SaveFeedback(rawScheduleID string, payload FeedbackPayload) (FeedbackItem, bool, error) {
	scheduleItem, found, scheduleErr := s.Schedule(rawScheduleID)
	if scheduleErr != nil {
		return FeedbackItem{}, false, scheduleErr
	}
	if !found {
		return FeedbackItem{}, false, nil
	}

	if s.db == nil {
		savedItem := FeedbackItem{
			ID:             uint64(scheduleItem.ID),
			ScheduleID:     scheduleItem.ID,
			ClassID:        scheduleItem.ClassID,
			ClassName:      scheduleItem.ClassName,
			CourseName:     scheduleItem.CourseName,
			TeacherName:    scheduleItem.TeacherName,
			LessonDate:     scheduleItem.LessonDate,
			Summary:        strings.TrimSpace(payload.Summary),
			LearningStatus: strings.TrimSpace(payload.LearningStatus),
			NextSuggestion: strings.TrimSpace(payload.NextSuggestion),
			ParentNotice:   strings.TrimSpace(payload.ParentNotice),
		}

		demoSaved := demo.SaveFeedback(demo.Feedback{
			ID:             int(savedItem.ID),
			ScheduleID:     int(savedItem.ScheduleID),
			ClassID:        int(savedItem.ClassID),
			ClassName:      savedItem.ClassName,
			CourseName:     savedItem.CourseName,
			TeacherName:    savedItem.TeacherName,
			LessonDate:     savedItem.LessonDate,
			Summary:        savedItem.Summary,
			LearningStatus: savedItem.LearningStatus,
			NextSuggestion: savedItem.NextSuggestion,
			ParentNotice:   savedItem.ParentNotice,
		})
		if !demoSaved {
			return FeedbackItem{}, false, nil
		}

		return savedItem, true, nil
	}

	scheduleID, parseErr := strconv.ParseUint(rawScheduleID, 10, 64)
	if parseErr != nil {
		return FeedbackItem{}, false, nil
	}

	now := time.Now()
	record := edumodel.ClassFeedback{
		ScheduleID:      scheduleID,
		ClassID:         scheduleItem.ClassID,
		Summary:         strings.TrimSpace(payload.Summary),
		LearningStatus:  strings.TrimSpace(payload.LearningStatus),
		NextSuggestion:  strings.TrimSpace(payload.NextSuggestion),
		ParentNotice:    strings.TrimSpace(payload.ParentNotice),
		CreatedByUserID: 1,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	saveErr := s.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "schedule_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"summary",
			"learning_status",
			"next_suggestion",
			"parent_notice",
			"updated_at",
		}),
	}).Create(&record).Error
	if saveErr != nil {
		return FeedbackItem{}, false, saveErr
	}

	savedItem, detailFound, detailErr := s.Feedback(rawScheduleID)
	if detailErr != nil {
		return FeedbackItem{}, false, detailErr
	}

	return savedItem, detailFound, nil
}

func (s *Service) Notices() ([]NoticeItem, error) {
	return s.NoticesWithFilter(NoticeFilter{})
}

func (s *Service) NoticesWithFilter(filter NoticeFilter) ([]NoticeItem, error) {
	if s.db == nil {
		return s.noticeItemsFromDemoWithFilter(filter), nil
	}

	query := `
SELECT
  id,
  title,
  content,
  notice_type AS category,
  target_scope,
  COALESCE(related_class_id, 0) AS related_class_id,
  COALESCE(related_schedule_id, 0) AS related_schedule_id,
  status,
  publish_at,
  created_at,
  author_name AS author
FROM notices
WHERE (? = 0 OR COALESCE(related_class_id, 0) = ?)
  AND (? = '' OR status = ?)
  AND (? = '' OR DATE(COALESCE(publish_at, created_at)) = ?)
ORDER BY COALESCE(publish_at, created_at) DESC, id DESC
`

	type noticeRow struct {
		ID                uint64     `gorm:"column:id"`
		Title             string     `gorm:"column:title"`
		Content           string     `gorm:"column:content"`
		Category          string     `gorm:"column:category"`
		TargetScope       string     `gorm:"column:target_scope"`
		RelatedClassID    uint64     `gorm:"column:related_class_id"`
		RelatedScheduleID uint64     `gorm:"column:related_schedule_id"`
		Status            string     `gorm:"column:status"`
		PublishAt         *time.Time `gorm:"column:publish_at"`
		CreatedAt         time.Time  `gorm:"column:created_at"`
		Author            string     `gorm:"column:author"`
	}

	var rows []noticeRow
	filterDate := strings.TrimSpace(filter.Date)
	filterStatus := strings.TrimSpace(filter.Status)
	listErr := s.db.Raw(
		query,
		filter.ClassID,
		filter.ClassID,
		filterStatus,
		filterStatus,
		filterDate,
		filterDate,
	).Scan(&rows).Error
	if listErr != nil {
		return nil, listErr
	}

	items := make([]NoticeItem, 0, len(rows))
	for _, row := range rows {
		publishAt := formatDateTime(row.PublishAt)
		if publishAt == "" {
			publishAt = row.CreatedAt.Format(dateTimeLayout)
		}

		items = append(items, NoticeItem{
			ID:                row.ID,
			Title:             row.Title,
			Content:           row.Content,
			Category:          row.Category,
			TargetScope:       row.TargetScope,
			RelatedClassID:    row.RelatedClassID,
			RelatedScheduleID: row.RelatedScheduleID,
			Status:            row.Status,
			PublishAt:         publishAt,
			Author:            row.Author,
		})
	}

	return items, nil
}

func (s *Service) Notice(rawID string) (NoticeItem, bool, error) {
	items, listErr := s.Notices()
	if listErr != nil {
		return NoticeItem{}, false, listErr
	}

	for _, item := range items {
		if fmt.Sprintf("%d", item.ID) == rawID {
			return item, true, nil
		}
	}

	return NoticeItem{}, false, nil
}

func (s *Service) NoticeTargets(rawNoticeID string) ([]NoticeTargetItem, error) {
	if s.db == nil {
		return noticeTargetItemsFromDemo(demo.NoticeTargets(rawNoticeID)), nil
	}

	query := `
SELECT
  n.target_scope AS scope_name,
  COALESCE(c.name, '') AS class_name,
  COALESCE(c.campus, '') AS campus,
  COALESCE(cs.id, 0) AS schedule_id,
  COALESCE(cs.schedule_date, NULL) AS schedule_date,
  COALESCE(cs.start_time, '') AS start_time,
  COALESCE(cs.end_time, '') AS end_time
FROM notices AS n
LEFT JOIN classes AS c
  ON c.id = n.related_class_id
LEFT JOIN class_schedules AS cs
  ON cs.id = n.related_schedule_id
WHERE n.id = ?
LIMIT 1
`

	type noticeTargetRow struct {
		ScopeName    string     `gorm:"column:scope_name"`
		ClassName    string     `gorm:"column:class_name"`
		Campus       string     `gorm:"column:campus"`
		ScheduleID   uint64     `gorm:"column:schedule_id"`
		ScheduleDate *time.Time `gorm:"column:schedule_date"`
		StartTime    string     `gorm:"column:start_time"`
		EndTime      string     `gorm:"column:end_time"`
	}

	var row noticeTargetRow
	findErr := s.db.Raw(query, rawNoticeID).Scan(&row).Error
	if findErr != nil {
		return nil, findErr
	}

	if row.ScopeName == "" && row.ClassName == "" {
		return []NoticeTargetItem{}, nil
	}

	items := make([]NoticeTargetItem, 0, 3)
	if row.ClassName != "" {
		items = append(items, NoticeTargetItem{
			Name:   row.ClassName,
			Type:   "关联班级",
			Campus: row.Campus,
		})
	}
	if row.ScopeName != "" {
		items = append(items, NoticeTargetItem{
			Name:   row.ScopeName,
			Type:   "通知范围",
			Campus: row.Campus,
		})
	}
	if row.ScheduleID > 0 && row.ScheduleDate != nil {
		items = append(items, NoticeTargetItem{
			Name:   fmt.Sprintf("%s %s", row.ScheduleDate.Format(dateLayout), formatLessonTime(row.StartTime, row.EndTime)),
			Type:   "关联课程安排",
			Campus: row.Campus,
		})
	}

	return items, nil
}

func (s *Service) CreateNotice(input NoticePayload) (NoticeItem, error) {
	if s.db == nil {
		notices := demo.Notices()
		nextID := len(notices) + 1

		now := time.Now()
		publishAt := now.Format(dateTimeLayout)
		if input.Status == "草稿" {
			publishAt = ""
		}

		item := NoticeItem{
			ID:                uint64(nextID),
			Title:             input.Title,
			Content:           input.Content,
			Category:          input.Category,
			TargetScope:       input.TargetScope,
			RelatedClassID:    input.RelatedClassID,
			RelatedScheduleID: input.RelatedScheduleID,
			Status:            input.Status,
			PublishAt:         publishAt,
			Author:            input.Author,
		}

		demo.SaveNotice(demo.Notice{
			ID:                nextID,
			Title:             item.Title,
			Content:           item.Content,
			Category:          item.Category,
			TargetScope:       item.TargetScope,
			RelatedClassID:    int(item.RelatedClassID),
			RelatedScheduleID: int(item.RelatedScheduleID),
			Status:            item.Status,
			PublishAt:         item.PublishAt,
			Author:            item.Author,
		})

		return item, nil
	}

	now := time.Now()
	record := edumodel.Notice{
		Title:       input.Title,
		Content:     input.Content,
		NoticeType:  input.Category,
		TargetScope: input.TargetScope,
		AuthorName:  input.Author,
		Status:      input.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if input.RelatedClassID > 0 {
		record.RelatedClassID = &input.RelatedClassID
	}
	if input.RelatedScheduleID > 0 {
		record.RelatedScheduleID = &input.RelatedScheduleID
	}
	if input.Status == "已发送" || input.Status == "待发送" {
		record.PublishAt = &now
	}

	createErr := s.db.Create(&record).Error
	if createErr != nil {
		return NoticeItem{}, createErr
	}

	createdItem, found, detailErr := s.Notice(fmt.Sprintf("%d", record.ID))
	if detailErr != nil {
		return NoticeItem{}, detailErr
	}
	if !found {
		return NoticeItem{
			ID:                record.ID,
			Title:             input.Title,
			Content:           input.Content,
			Category:          input.Category,
			TargetScope:       input.TargetScope,
			RelatedClassID:    input.RelatedClassID,
			RelatedScheduleID: input.RelatedScheduleID,
			Status:            input.Status,
			PublishAt:         formatDateTime(record.PublishAt),
			Author:            input.Author,
		}, nil
	}

	return createdItem, nil
}

func (s *Service) UpdateNotice(rawID string, input NoticePayload) (NoticeItem, bool, error) {
	if s.db == nil {
		rawItem, found := demo.FindNotice(rawID)
		if !found {
			return NoticeItem{}, false, nil
		}

		publishAt := rawItem.PublishAt
		if input.Status == "草稿" {
			publishAt = ""
		}
		if publishAt == "" && (input.Status == "待发送" || input.Status == "已发送") {
			publishAt = time.Now().Format(dateTimeLayout)
		}

		item := NoticeItem{
			ID:                uint64(rawItem.ID),
			Title:             input.Title,
			Content:           input.Content,
			Category:          input.Category,
			TargetScope:       input.TargetScope,
			RelatedClassID:    input.RelatedClassID,
			RelatedScheduleID: input.RelatedScheduleID,
			Status:            input.Status,
			PublishAt:         publishAt,
			Author:            input.Author,
		}

		demo.SaveNotice(demo.Notice{
			ID:                rawItem.ID,
			Title:             item.Title,
			Content:           item.Content,
			Category:          item.Category,
			TargetScope:       item.TargetScope,
			RelatedClassID:    int(item.RelatedClassID),
			RelatedScheduleID: int(item.RelatedScheduleID),
			Status:            item.Status,
			PublishAt:         item.PublishAt,
			Author:            item.Author,
		})

		return item, true, nil
	}

	updateValues := map[string]any{
		"title":               input.Title,
		"content":             input.Content,
		"notice_type":         input.Category,
		"target_scope":        input.TargetScope,
		"related_class_id":    nil,
		"related_schedule_id": nil,
		"status":              input.Status,
		"author_name":         input.Author,
		"updated_at":          time.Now(),
	}
	if input.RelatedClassID > 0 {
		updateValues["related_class_id"] = input.RelatedClassID
	}
	if input.RelatedScheduleID > 0 {
		updateValues["related_schedule_id"] = input.RelatedScheduleID
	}
	if input.Status == "草稿" {
		updateValues["publish_at"] = nil
	} else {
		updateValues["publish_at"] = time.Now()
	}

	updateResult := s.db.Model(&edumodel.Notice{}).
		Where("id = ?", rawID).
		Updates(updateValues)
	if updateResult.Error != nil {
		return NoticeItem{}, false, updateResult.Error
	}
	if updateResult.RowsAffected == 0 {
		return NoticeItem{}, false, nil
	}

	updatedItem, found, detailErr := s.Notice(rawID)
	if detailErr != nil {
		return NoticeItem{}, false, detailErr
	}

	return updatedItem, found, nil
}

func (s *Service) SendNotice(rawID string) (NoticeItem, bool, error) {
	item, found, itemErr := s.Notice(rawID)
	if itemErr != nil || !found {
		return item, found, itemErr
	}

	input := NoticePayload{
		Title:             item.Title,
		Content:           item.Content,
		Category:          item.Category,
		TargetScope:       item.TargetScope,
		RelatedClassID:    item.RelatedClassID,
		RelatedScheduleID: item.RelatedScheduleID,
		Status:            "已发送",
		Author:            item.Author,
	}

	return s.UpdateNotice(rawID, input)
}

func (s *Service) createScheduleNotice(rawScheduleID string, input NoticePayload) error {
	noticeItem, noticeFound, noticeErr := s.findNoticeByScheduleAndCategory(rawScheduleID, input.Category)
	if noticeErr != nil {
		return noticeErr
	}
	if noticeFound {
		_, _, updateErr := s.UpdateNotice(fmt.Sprintf("%d", noticeItem.ID), input)
		return updateErr
	}

	_, createErr := s.CreateNotice(input)
	return createErr
}

func (s *Service) findNoticeByScheduleAndCategory(rawScheduleID string, category string) (NoticeItem, bool, error) {
	notices, noticeErr := s.Notices()
	if noticeErr != nil {
		return NoticeItem{}, false, noticeErr
	}

	scheduleID, parseErr := strconv.ParseUint(rawScheduleID, 10, 64)
	if parseErr != nil {
		return NoticeItem{}, false, nil
	}

	for _, item := range notices {
		if item.RelatedScheduleID == scheduleID && item.Category == category {
			return item, true, nil
		}
	}

	return NoticeItem{}, false, nil
}

func buildClassTargetScope(className string) string {
	trimmedClassName := strings.TrimSpace(className)
	if trimmedClassName == "" {
		return "相关学员家长"
	}

	return fmt.Sprintf("%s家长群", trimmedClassName)
}

func buildRescheduleNoticeContent(originalItem ScheduleItem, replacementItem ScheduleItem, remark string) string {
	message := fmt.Sprintf(
		"%s原定于%s %s在%s上课，现调整为%s %s在%s上课，请家长留意时间变化。",
		originalItem.ClassName,
		originalItem.LessonDate,
		originalItem.LessonTime,
		originalItem.Classroom,
		replacementItem.LessonDate,
		replacementItem.LessonTime,
		replacementItem.Classroom,
	)
	if remark != "" {
		return fmt.Sprintf("%s 备注：%s", message, remark)
	}

	return message
}

func buildCancelNoticeContent(scheduleItem ScheduleItem, remark string) string {
	message := fmt.Sprintf(
		"%s原定于%s %s在%s上课，本次课程已停课，请家长留意后续安排。",
		scheduleItem.ClassName,
		scheduleItem.LessonDate,
		scheduleItem.LessonTime,
		scheduleItem.Classroom,
	)
	if remark != "" {
		return fmt.Sprintf("%s 原因：%s", message, remark)
	}

	return message
}

func buildMakeupNoticeContent(originalItem ScheduleItem, createdItem ScheduleItem, remark string) string {
	message := fmt.Sprintf(
		"%s已新增补课安排，补课时间为%s %s，地点为%s，请家长提前做好到课安排。",
		originalItem.ClassName,
		createdItem.LessonDate,
		createdItem.LessonTime,
		createdItem.Classroom,
	)
	if remark != "" {
		return fmt.Sprintf("%s 备注：%s", message, remark)
	}

	return message
}

func (s *Service) seedIfEmpty() error {
	var teacherCount int64
	countErr := s.db.Model(&edumodel.Teacher{}).Count(&teacherCount).Error
	if countErr != nil {
		return countErr
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		seedAccessErr := s.seedAccessControlIfEmpty(tx, now)
		if seedAccessErr != nil {
			return seedAccessErr
		}

		if teacherCount > 0 {
			return nil
		}

		today := startOfDay(now)
		tomorrow := today.AddDate(0, 0, 1)

		teachers := []edumodel.Teacher{
			{ID: 1, Name: "周老师", Mobile: "13800000001", MainSubject: "数学思维", EmploymentType: "全职", WeeklyHours: 18, Campus: "明发校区", Status: "在职", CreatedAt: now, UpdatedAt: now},
			{ID: 2, Name: "林老师", Mobile: "13800000002", MainSubject: "英语阅读", EmploymentType: "全职", WeeklyHours: 16, Campus: "百汇校区", Status: "在职", CreatedAt: now, UpdatedAt: now},
			{ID: 3, Name: "陈老师", Mobile: "13800000003", MainSubject: "创意美术", EmploymentType: "兼职", WeeklyHours: 12, Campus: "明发校区", Status: "排课中", CreatedAt: now, UpdatedAt: now},
		}

		teacherErr := tx.Create(&teachers).Error
		if teacherErr != nil {
			return teacherErr
		}

		courses := []edumodel.Course{
			{ID: 1, Name: "数学思维", Category: "数学", Description: "围绕数感、推理和应用题训练，适合周末持续提升。", AgeRange: "8-10岁", LessonDurationMinutes: 90, TotalLessons: 24, DeliveryType: "线下", Status: "启用", CreatedAt: now, UpdatedAt: now},
			{ID: 2, Name: "英语阅读", Category: "英语", Description: "通过分级阅读和表达训练，帮助孩子建立英文阅读习惯。", AgeRange: "9-11岁", LessonDurationMinutes: 90, TotalLessons: 24, DeliveryType: "线下", Status: "启用", CreatedAt: now, UpdatedAt: now},
			{ID: 3, Name: "创意美术", Category: "美术", Description: "以主题创作为主，兼顾色彩感受和动手表达。", AgeRange: "6-8岁", LessonDurationMinutes: 90, TotalLessons: 16, DeliveryType: "线下", Status: "启用", CreatedAt: now, UpdatedAt: now},
		}

		courseErr := tx.Create(&courses).Error
		if courseErr != nil {
			return courseErr
		}

		classes := []edumodel.Class{
			{ID: 1, Name: "周末奥数提高班", CourseID: 1, TeacherID: 1, Campus: "明发校区", Capacity: 16, WeeklySchedule: "周六 09:00-10:30", StartDate: datePointer(today), EndDate: datePointer(today.AddDate(0, 3, 0)), Status: "开班中", CreatedAt: now, UpdatedAt: now},
			{ID: 2, Name: "英语阅读进阶班", CourseID: 2, TeacherID: 2, Campus: "百汇校区", Capacity: 12, WeeklySchedule: "周六 14:00-15:30", StartDate: datePointer(today), EndDate: datePointer(today.AddDate(0, 3, 0)), Status: "开班中", CreatedAt: now, UpdatedAt: now},
			{ID: 3, Name: "少儿创意美术班", CourseID: 3, TeacherID: 3, Campus: "明发校区", Capacity: 10, WeeklySchedule: "周日 10:00-11:30", StartDate: datePointer(today), EndDate: datePointer(today.AddDate(0, 2, 0)), Status: "待满班", CreatedAt: now, UpdatedAt: now},
		}

		classErr := tx.Create(&classes).Error
		if classErr != nil {
			return classErr
		}

		students := []edumodel.Student{
			{ID: 1, Name: "李一诺", SchoolName: "实验小学", GradeName: "三年级", Campus: "明发校区", RemainingHours: 18, Status: "在读", CreatedAt: now, UpdatedAt: now},
			{ID: 2, Name: "王梓涵", SchoolName: "实验小学", GradeName: "四年级", Campus: "明发校区", RemainingHours: 12, Status: "在读", CreatedAt: now, UpdatedAt: now},
			{ID: 3, Name: "陈可欣", SchoolName: "百汇小学", GradeName: "五年级", Campus: "百汇校区", RemainingHours: 24, Status: "在读", CreatedAt: now, UpdatedAt: now},
			{ID: 4, Name: "张沐阳", SchoolName: "实验小学", GradeName: "二年级", Campus: "明发校区", RemainingHours: 10, Status: "待续费", CreatedAt: now, UpdatedAt: now},
		}

		studentErr := tx.Create(&students).Error
		if studentErr != nil {
			return studentErr
		}

		guardians := []edumodel.StudentGuardian{
			{ID: 1, StudentID: 1, Name: "李女士", Relation: "母亲", Mobile: "13900000001", IsPrimary: true, CreatedAt: now, UpdatedAt: now},
			{ID: 2, StudentID: 2, Name: "王先生", Relation: "父亲", Mobile: "13900000002", IsPrimary: true, CreatedAt: now, UpdatedAt: now},
			{ID: 3, StudentID: 3, Name: "陈女士", Relation: "母亲", Mobile: "13900000003", IsPrimary: true, CreatedAt: now, UpdatedAt: now},
			{ID: 4, StudentID: 4, Name: "张女士", Relation: "母亲", Mobile: "13900000004", IsPrimary: true, CreatedAt: now, UpdatedAt: now},
		}

		guardianErr := tx.Create(&guardians).Error
		if guardianErr != nil {
			return guardianErr
		}

		classStudents := []edumodel.ClassStudent{
			{ID: 1, ClassID: 1, StudentID: 1, JoinDate: datePointer(today), Status: activeClassStudentStatus, CreatedAt: now, UpdatedAt: now},
			{ID: 2, ClassID: 1, StudentID: 2, JoinDate: datePointer(today), Status: activeClassStudentStatus, CreatedAt: now, UpdatedAt: now},
			{ID: 3, ClassID: 2, StudentID: 3, JoinDate: datePointer(today), Status: activeClassStudentStatus, CreatedAt: now, UpdatedAt: now},
			{ID: 4, ClassID: 3, StudentID: 4, JoinDate: datePointer(today), Status: activeClassStudentStatus, CreatedAt: now, UpdatedAt: now},
		}

		classStudentErr := tx.Create(&classStudents).Error
		if classStudentErr != nil {
			return classStudentErr
		}

		schedules := []edumodel.ClassSchedule{
			{ID: 1, ClassID: 1, CourseID: 1, TeacherID: 1, ScheduleType: "常规课", ScheduleDate: today, StartTime: "09:00", EndTime: "10:30", Location: "A201", Status: "待签到", CreatedAt: now, UpdatedAt: now},
			{ID: 2, ClassID: 2, CourseID: 2, TeacherID: 2, ScheduleType: "常规课", ScheduleDate: today, StartTime: "14:00", EndTime: "15:30", Location: "B103", Status: "已完成", CreatedAt: now, UpdatedAt: now},
			{ID: 3, ClassID: 3, CourseID: 3, TeacherID: 3, ScheduleType: "常规课", ScheduleDate: tomorrow, StartTime: "10:00", EndTime: "11:30", Location: "Art-2", Status: "待上课", CreatedAt: now, UpdatedAt: now},
		}

		scheduleErr := tx.Create(&schedules).Error
		if scheduleErr != nil {
			return scheduleErr
		}

		attendanceRecords := []edumodel.AttendanceRecord{
			{ID: 1, ScheduleID: 2, StudentID: 3, Status: "请假", Remark: "家长上午已请假", CheckedAt: &now, UpdatedBy: "林老师", CreatedAt: now, UpdatedAt: now},
		}

		attendanceRecordErr := tx.Create(&attendanceRecords).Error
		if attendanceRecordErr != nil {
			return attendanceRecordErr
		}

		homeworks := []edumodel.Homework{
			{ID: 1, ScheduleID: 1, ClassID: 1, Title: "思维训练第 4 讲课后练习", Content: "完成练习册第 12-15 页，并整理两道易错题的解题步骤。", SubmissionNote: "下节课前带回纸质作业，家长协助检查书写完整度。", CreatedByUserID: 1, Status: "published", CreatedAt: now, UpdatedAt: now},
			{ID: 2, ScheduleID: 2, ClassID: 2, Title: "英语阅读分级复述任务", Content: "完成本周阅读卡第 3 篇并录一段 60 秒英文复述。", SubmissionNote: "家长可拍视频发到班级群，老师下次课统一点评。", CreatedByUserID: 1, Status: "published", CreatedAt: now, UpdatedAt: now},
		}

		homeworkErr := tx.Create(&homeworks).Error
		if homeworkErr != nil {
			return homeworkErr
		}

		classFeedbacks := []edumodel.ClassFeedback{
			{ID: 1, ScheduleID: 1, ClassID: 1, Summary: "课堂整体专注度不错，能跟上本节推理题节奏。", LearningStatus: "大部分同学能独立完成基础题，个别同学在多步骤表达上还需要提醒。", NextSuggestion: "下次课前建议再复习一次本周错题，带着自己的思路来讲解。", ParentNotice: "请家长关注孩子列式步骤是否完整，不只看最终答案。", CreatedByUserID: 1, CreatedAt: now, UpdatedAt: now},
			{ID: 2, ScheduleID: 2, ClassID: 2, Summary: "本节重点放在阅读理解和复述表达，课堂互动稳定。", LearningStatus: "学生对关键词抓取已经有进步，但整句复述还需要更多练习。", NextSuggestion: "课后多做跟读和复述训练，下节课继续检查语音语调。", ParentNotice: "如果孩子本周无法按时提交复述视频，请提前在群里说明。", CreatedByUserID: 1, CreatedAt: now, UpdatedAt: now},
		}

		classFeedbackErr := tx.Create(&classFeedbacks).Error
		if classFeedbackErr != nil {
			return classFeedbackErr
		}

		publishAtOne := now.Add(-6 * time.Hour)
		publishAtTwo := now.Add(-2 * time.Hour)
		publishAtThree := now.Add(2 * time.Hour)
		relatedClassID := uint64(3)

		notices := []edumodel.Notice{
			{ID: 1, Title: "端午节放假安排", Content: "本周放假安排请注意查看。", NoticeType: "校区通知", TargetScope: "全部学员家长", AuthorName: "运营老师", Status: "已发送", PublishAt: &publishAtOne, CreatedAt: now, UpdatedAt: now},
			{ID: 2, Title: "六月续费提醒名单确认", Content: "请班主任确认续费提醒名单。", NoticeType: "续费提醒", TargetScope: "待续费学员家长", AuthorName: "班主任", Status: "草稿", PublishAt: &publishAtTwo, CreatedAt: now, UpdatedAt: now},
			{ID: 3, Title: "周末美术课材料准备说明", Content: "请家长提前准备水彩笔与画纸。", NoticeType: "课程通知", TargetScope: "少儿创意美术班", AuthorName: "教务老师", RelatedClassID: &relatedClassID, Status: "待发送", PublishAt: &publishAtThree, CreatedAt: now, UpdatedAt: now},
		}

		noticeErr := tx.Create(&notices).Error
		if noticeErr != nil {
			return noticeErr
		}

		return nil
	})
}

func teacherItemsFromDemo() []TeacherItem {
	items := make([]TeacherItem, 0, len(demo.Teachers()))
	for _, item := range demo.Teachers() {
		items = append(items, TeacherItem{
			ID:             uint64(item.ID),
			Name:           item.Name,
			Mobile:         item.Mobile,
			Title:          item.Title,
			MainSubject:    item.MainSubject,
			EmploymentType: item.EmploymentType,
			WeeklyHours:    item.WeeklyHours,
			Campus:         item.Campus,
			Status:         item.Status,
			Remark:         item.Remark,
		})
	}
	return items
}

func teacherItemFromDemo(rawID string) (TeacherItem, bool, error) {
	item, found := demo.FindTeacher(rawID)
	if !found {
		return TeacherItem{}, false, nil
	}
	return TeacherItem{
		ID:             uint64(item.ID),
		Name:           item.Name,
		Mobile:         item.Mobile,
		Title:          item.Title,
		MainSubject:    item.MainSubject,
		EmploymentType: item.EmploymentType,
		WeeklyHours:    item.WeeklyHours,
		Campus:         item.Campus,
		Status:         item.Status,
		Remark:         item.Remark,
	}, true, nil
}

func courseItemsFromDemo(source []demo.Course) []CourseItem {
	items := make([]CourseItem, 0, len(source))
	for _, item := range source {
		items = append(items, CourseItem{
			ID:                    uint64(item.ID),
			Name:                  item.Name,
			Category:              item.Category,
			Description:           item.Description,
			AgeRange:              item.AgeRange,
			LessonDurationMinutes: item.LessonDurationMinutes,
			TotalLessons:          item.TotalLessons,
			DeliveryType:          item.DeliveryType,
			Status:                item.Status,
			ClassCount:            item.ClassCount,
		})
	}
	return items
}

func courseItemFromDemo(rawID string) (CourseItem, bool, error) {
	item, found := demo.FindCourse(rawID)
	if !found {
		return CourseItem{}, false, nil
	}

	return CourseItem{
		ID:                    uint64(item.ID),
		Name:                  item.Name,
		Category:              item.Category,
		Description:           item.Description,
		AgeRange:              item.AgeRange,
		LessonDurationMinutes: item.LessonDurationMinutes,
		TotalLessons:          item.TotalLessons,
		DeliveryType:          item.DeliveryType,
		Status:                item.Status,
		ClassCount:            item.ClassCount,
	}, true, nil
}

func courseItemFromPayload(id uint64, input CoursePayload) CourseItem {
	return CourseItem{
		ID:                    id,
		Name:                  input.Name,
		Category:              input.Category,
		Description:           input.Description,
		AgeRange:              input.AgeRange,
		LessonDurationMinutes: input.LessonDurationMinutes,
		TotalLessons:          input.TotalLessons,
		DeliveryType:          input.DeliveryType,
		Status:                input.Status,
		ClassCount:            0,
	}
}

func studentItemsFromDemo() []StudentItem {
	return studentItemsFromDemoWithSlice(demo.Students())
}

func studentItemsFromDemoWithSlice(source []demo.Student) []StudentItem {
	items := make([]StudentItem, 0, len(source))
	for _, item := range source {
		items = append(items, StudentItem{
			ID:             uint64(item.ID),
			Name:           item.Name,
			Grade:          item.Grade,
			ParentName:     item.ParentName,
			ParentMobile:   item.ParentMobile,
			Campus:         item.Campus,
			ClassID:        uint64(item.ClassID),
			ClassName:      item.ClassName,
			RemainingHours: item.RemainingHours,
			Status:         item.Status,
		})
	}
	return items
}

func studentItemFromDemo(rawID string) (StudentItem, bool, error) {
	item, found := demo.FindStudent(rawID)
	if !found {
		return StudentItem{}, false, nil
	}
	return StudentItem{
		ID:             uint64(item.ID),
		Name:           item.Name,
		Grade:          item.Grade,
		ParentName:     item.ParentName,
		ParentMobile:   item.ParentMobile,
		Campus:         item.Campus,
		ClassID:        uint64(item.ClassID),
		ClassName:      item.ClassName,
		RemainingHours: item.RemainingHours,
		Status:         item.Status,
	}, true, nil
}

func studentProfileFromStudentItem(
	studentItem StudentItem,
	schoolName string,
	gender string,
	remark string,
) StudentProfile {
	return StudentProfile{
		ID:             studentItem.ID,
		Name:           studentItem.Name,
		Gender:         gender,
		SchoolName:     schoolName,
		Grade:          studentItem.Grade,
		ParentName:     studentItem.ParentName,
		ParentMobile:   studentItem.ParentMobile,
		Campus:         studentItem.Campus,
		RemainingHours: studentItem.RemainingHours,
		Status:         studentItem.Status,
		Remark:         remark,
	}
}

func classItemsFromDemo(source []demo.Class) []ClassItem {
	items := make([]ClassItem, 0, len(source))
	for _, item := range source {
		items = append(items, ClassItem{
			ID:             uint64(item.ID),
			CourseID:       uint64(item.CourseID),
			Name:           item.Name,
			CourseName:     item.CourseName,
			TeacherID:      uint64(item.TeacherID),
			TeacherName:    item.TeacherName,
			Campus:         item.Campus,
			StudentCount:   item.StudentCount,
			Capacity:       item.Capacity,
			WeeklySchedule: item.WeeklySchedule,
			StartDate:      item.StartDate,
			EndDate:        item.EndDate,
			Status:         item.Status,
			Remark:         item.Remark,
		})
	}
	return items
}

func classItemFromDemo(rawID string) (ClassItem, bool, error) {
	item, found := demo.FindClass(rawID)
	if !found {
		return ClassItem{}, false, nil
	}
	return ClassItem{
		ID:             uint64(item.ID),
		CourseID:       uint64(item.CourseID),
		Name:           item.Name,
		CourseName:     item.CourseName,
		TeacherID:      uint64(item.TeacherID),
		TeacherName:    item.TeacherName,
		Campus:         item.Campus,
		StudentCount:   item.StudentCount,
		Capacity:       item.Capacity,
		WeeklySchedule: item.WeeklySchedule,
		StartDate:      item.StartDate,
		EndDate:        item.EndDate,
		Status:         item.Status,
		Remark:         item.Remark,
	}, true, nil
}

func scheduleItemsFromDemo(source []demo.Schedule) []ScheduleItem {
	items := make([]ScheduleItem, 0, len(source))
	for _, item := range source {
		items = append(items, ScheduleItem{
			ID:               uint64(item.ID),
			ClassID:          uint64(item.ClassID),
			SourceScheduleID: uint64(item.SourceScheduleID),
			ClassName:        item.ClassName,
			CourseName:       item.CourseName,
			TeacherID:        uint64(item.TeacherID),
			TeacherName:      item.TeacherName,
			Campus:           item.Campus,
			Classroom:        item.Classroom,
			ScheduleType:     "常规课",
			LessonDate:       item.LessonDate,
			LessonTime:       item.LessonTime,
			AttendanceStatus: item.AttendanceStatus,
		})
	}
	return items
}

func noticeItemsFromDemo(source []demo.Notice) []NoticeItem {
	items := make([]NoticeItem, 0, len(source))
	for _, item := range source {
		items = append(items, NoticeItem{
			ID:                uint64(item.ID),
			Title:             item.Title,
			Content:           item.Content,
			Category:          item.Category,
			TargetScope:       item.TargetScope,
			RelatedClassID:    uint64(item.RelatedClassID),
			RelatedScheduleID: uint64(item.RelatedScheduleID),
			Status:            item.Status,
			PublishAt:         item.PublishAt,
			Author:            item.Author,
		})
	}
	return items
}

func (s *Service) noticeItemsFromDemoWithFilter(filter NoticeFilter) []NoticeItem {
	items := noticeItemsFromDemo(demo.Notices())
	filteredItems := make([]NoticeItem, 0, len(items))
	filterDate := strings.TrimSpace(filter.Date)
	filterStatus := strings.TrimSpace(filter.Status)

	for _, item := range items {
		if filter.ClassID > 0 && item.RelatedClassID != filter.ClassID {
			continue
		}
		if filterStatus != "" && item.Status != filterStatus {
			continue
		}
		if filterDate != "" && !strings.HasPrefix(item.PublishAt, filterDate) {
			continue
		}

		filteredItems = append(filteredItems, item)
	}

	return filteredItems
}

func noticeTargetItemsFromDemo(source []demo.NoticeTarget) []NoticeTargetItem {
	items := make([]NoticeTargetItem, 0, len(source))
	for _, item := range source {
		items = append(items, NoticeTargetItem{
			Name:   item.Name,
			Type:   item.Type,
			Campus: item.Campus,
		})
	}
	return items
}

func attendanceItemsFromDemo(source []demo.AttendanceItem) []AttendanceItem {
	items := make([]AttendanceItem, 0, len(source))
	for _, item := range source {
		items = append(items, AttendanceItem{
			StudentID:    uint64(item.StudentID),
			StudentName:  item.StudentName,
			Grade:        item.Grade,
			ParentMobile: item.ParentMobile,
			Status:       item.Status,
			Remark:       item.Remark,
			UpdatedBy:    item.UpdatedBy,
			UpdatedAt:    item.UpdatedAt,
		})
	}
	return items
}

func homeworkItemsFromDemo(source []demo.Homework) []HomeworkItem {
	items := make([]HomeworkItem, 0, len(source))
	for _, item := range source {
		items = append(items, HomeworkItem{
			ID:             uint64(item.ID),
			ScheduleID:     uint64(item.ScheduleID),
			ClassID:        uint64(item.ClassID),
			ClassName:      item.ClassName,
			CourseName:     item.CourseName,
			TeacherName:    item.TeacherName,
			LessonDate:     item.LessonDate,
			Title:          item.Title,
			Content:        item.Content,
			SubmissionNote: item.SubmissionNote,
			Status:         item.Status,
		})
	}

	return items
}

func homeworkItemFromDemo(rawScheduleID string) (HomeworkItem, bool, error) {
	item, found := demo.HomeworkBySchedule(rawScheduleID)
	if !found {
		return HomeworkItem{}, false, nil
	}

	return HomeworkItem{
		ID:             uint64(item.ID),
		ScheduleID:     uint64(item.ScheduleID),
		ClassID:        uint64(item.ClassID),
		ClassName:      item.ClassName,
		CourseName:     item.CourseName,
		TeacherName:    item.TeacherName,
		LessonDate:     item.LessonDate,
		Title:          item.Title,
		Content:        item.Content,
		SubmissionNote: item.SubmissionNote,
		Status:         item.Status,
	}, true, nil
}

func homeworkItemFromDemoWithSchedule(scheduleItem ScheduleItem) (HomeworkItem, bool) {
	item, found := demo.HomeworkBySchedule(fmt.Sprintf("%d", scheduleItem.ID))
	if !found {
		return HomeworkItem{}, false
	}

	normalizedStatus := normalizeHomeworkStatus(item.Status)
	if normalizedStatus == "" {
		normalizedStatus = "published"
	}

	return HomeworkItem{
		ID:             uint64(item.ID),
		ScheduleID:     scheduleItem.ID,
		ClassID:        scheduleItem.ClassID,
		ClassName:      scheduleItem.ClassName,
		CourseName:     scheduleItem.CourseName,
		TeacherName:    scheduleItem.TeacherName,
		LessonDate:     scheduleItem.LessonDate,
		Title:          item.Title,
		Content:        item.Content,
		SubmissionNote: item.SubmissionNote,
		Status:         normalizedStatus,
	}, true
}

func feedbackItemsFromDemo(source []demo.Feedback) []FeedbackItem {
	items := make([]FeedbackItem, 0, len(source))
	for _, item := range source {
		items = append(items, FeedbackItem{
			ID:             uint64(item.ID),
			ScheduleID:     uint64(item.ScheduleID),
			ClassID:        uint64(item.ClassID),
			ClassName:      item.ClassName,
			CourseName:     item.CourseName,
			TeacherName:    item.TeacherName,
			LessonDate:     item.LessonDate,
			Summary:        item.Summary,
			LearningStatus: item.LearningStatus,
			NextSuggestion: item.NextSuggestion,
			ParentNotice:   item.ParentNotice,
		})
	}

	return items
}

func feedbackItemFromDemo(rawScheduleID string) (FeedbackItem, bool, error) {
	item, found := demo.FeedbackBySchedule(rawScheduleID)
	if !found {
		return FeedbackItem{}, false, nil
	}

	return FeedbackItem{
		ID:             uint64(item.ID),
		ScheduleID:     uint64(item.ScheduleID),
		ClassID:        uint64(item.ClassID),
		ClassName:      item.ClassName,
		CourseName:     item.CourseName,
		TeacherName:    item.TeacherName,
		LessonDate:     item.LessonDate,
		Summary:        item.Summary,
		LearningStatus: item.LearningStatus,
		NextSuggestion: item.NextSuggestion,
		ParentNotice:   item.ParentNotice,
	}, true, nil
}

func feedbackItemFromDemoWithSchedule(scheduleItem ScheduleItem) (FeedbackItem, bool) {
	item, found := demo.FeedbackBySchedule(fmt.Sprintf("%d", scheduleItem.ID))
	if !found {
		return FeedbackItem{}, false
	}

	return FeedbackItem{
		ID:             uint64(item.ID),
		ScheduleID:     scheduleItem.ID,
		ClassID:        scheduleItem.ClassID,
		ClassName:      scheduleItem.ClassName,
		CourseName:     scheduleItem.CourseName,
		TeacherName:    scheduleItem.TeacherName,
		LessonDate:     scheduleItem.LessonDate,
		Summary:        item.Summary,
		LearningStatus: item.LearningStatus,
		NextSuggestion: item.NextSuggestion,
		ParentNotice:   item.ParentNotice,
	}, true
}

func summarizeAttendanceItems(items []AttendanceItem) (int, int, int, int) {
	presentCount := 0
	leaveCount := 0
	absentCount := 0
	pendingCount := 0

	for _, item := range items {
		switch normalizeAttendanceItemStatus(item.Status) {
		case "已到", "补签":
			presentCount++
		case "请假":
			leaveCount++
		case "缺席":
			absentCount++
		default:
			pendingCount++
		}
	}

	return presentCount, leaveCount, absentCount, pendingCount
}

func teacherRecentSchedules(items []ScheduleItem, limit int) []ScheduleItem {
	if len(items) == 0 || limit <= 0 {
		return []ScheduleItem{}
	}

	today := startOfDay(time.Now())
	upcomingItems := make([]ScheduleItem, 0, len(items))
	pastItems := make([]ScheduleItem, 0, len(items))

	for _, item := range items {
		lessonDate, parseErr := time.ParseInLocation(dateLayout, item.LessonDate, time.Local)
		if parseErr != nil {
			upcomingItems = append(upcomingItems, item)
			continue
		}

		if startOfDay(lessonDate).Before(today) {
			pastItems = append(pastItems, item)
			continue
		}

		upcomingItems = append(upcomingItems, item)
	}

	recentItems := make([]ScheduleItem, 0, limit)
	for _, item := range upcomingItems {
		recentItems = append(recentItems, item)
		if len(recentItems) >= limit {
			return recentItems
		}
	}

	for index := len(pastItems) - 1; index >= 0; index-- {
		recentItems = append(recentItems, pastItems[index])
		if len(recentItems) >= limit {
			break
		}
	}

	return recentItems
}

func defaultAttendanceItemStatus(scheduleStatus string, index int, total int) string {
	switch scheduleStatus {
	case "已完成":
		return "已到"
	case "请假待批":
		if total > 0 && index == total-1 {
			return "请假"
		}
		return "已到"
	default:
		return "待确认"
	}
}

func normalizeAttendanceItemStatus(status string) string {
	switch strings.TrimSpace(status) {
	case "已到", "请假", "缺席", "补签", "待确认":
		return strings.TrimSpace(status)
	default:
		return ""
	}
}

func sessionAttendanceStatus(baseStatus string, pendingCount int) string {
	switch strings.TrimSpace(baseStatus) {
	case "待上课", "已调课", "已停课", "请假待批":
		return strings.TrimSpace(baseStatus)
	case "已完成":
		if pendingCount > 0 {
			return "待签到"
		}
		return "已完成"
	default:
		if pendingCount > 0 {
			return "待签到"
		}
		return "已完成"
	}
}

func nextSavedAttendanceStatus(pendingCount int) string {
	if pendingCount > 0 {
		return "待签到"
	}

	return "已完成"
}

func normalizeHomeworkStatus(status string) string {
	switch strings.TrimSpace(status) {
	case "draft", "published":
		return strings.TrimSpace(status)
	default:
		return ""
	}
}

func guardianRelationFromName(name string) string {
	switch {
	case strings.Contains(name, "先生"):
		return "父亲"
	case strings.Contains(name, "女士"):
		return "母亲"
	default:
		return "家长"
	}
}

func (s *Service) courseQuery(filter CourseFilter) *gorm.DB {
	courseQuery := s.db.Table("courses AS c").
		Select(`
c.id,
c.name,
c.category,
COALESCE(c.description, '') AS description,
c.age_range,
c.lesson_duration_minutes,
c.total_lessons,
c.delivery_type,
c.status,
COUNT(DISTINCT cl.id) AS class_count
`).
		Joins("LEFT JOIN classes AS cl ON cl.course_id = c.id")

	keyword := strings.TrimSpace(filter.Keyword)
	if keyword != "" {
		likeKeyword := "%" + keyword + "%"
		courseQuery = courseQuery.Where(
			"(c.name LIKE ? OR c.category LIKE ? OR c.description LIKE ?)",
			likeKeyword,
			likeKeyword,
			likeKeyword,
		)
	}

	category := strings.TrimSpace(filter.Category)
	if category != "" {
		courseQuery = courseQuery.Where("c.category = ?", category)
	}

	status := strings.TrimSpace(filter.Status)
	if status != "" {
		courseQuery = courseQuery.Where("c.status = ?", status)
	}

	return courseQuery.Group(`
c.id,
c.name,
c.category,
c.description,
c.age_range,
c.lesson_duration_minutes,
c.total_lessons,
c.delivery_type,
c.status
`)
}

func formatLessonTime(startTime string, endTime string) string {
	return shortClock(startTime) + "-" + shortClock(endTime)
}

func shortClock(value string) string {
	if len(value) >= 5 {
		return value[:5]
	}
	return value
}

func formatDateTime(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format(dateTimeLayout)
}

func parseOptionalDate(raw string) (*time.Time, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}

	parsedValue, parseErr := time.ParseInLocation(dateLayout, trimmed, time.Local)
	if parseErr != nil {
		return nil, parseErr
	}

	normalizedValue := startOfDay(parsedValue)
	return &normalizedValue, nil
}

func parseClassDates(startDate string, endDate string) (*time.Time, *time.Time, error) {
	parsedStartDate, startParseErr := parseOptionalDate(startDate)
	if startParseErr != nil {
		return nil, nil, startParseErr
	}

	parsedEndDate, endParseErr := parseOptionalDate(endDate)
	if endParseErr != nil {
		return nil, nil, endParseErr
	}

	return parsedStartDate, parsedEndDate, nil
}

func datePointer(value time.Time) *time.Time {
	return &value
}

func startOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}
