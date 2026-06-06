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

var ErrInvalidAttendanceStatus = errors.New("invalid attendance status")

const attendanceRecordSyntheticFactor uint64 = 1000000

type Service struct {
	db *gorm.DB
}

type Scope struct {
	UserID         uint64
	PrimaryRole    string
	TeacherID      uint64
	RestrictToSelf bool
}

type Option struct {
	Value uint64 `json:"value"`
	Label string `json:"label"`
}

type TeacherItem struct {
	ID             uint64 `json:"id"`
	UserID         uint64 `json:"userId" gorm:"column:user_id"`
	UserName       string `json:"userName" gorm:"column:user_name"`
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

type TeacherFilter struct {
	Keyword        string
	Status         string
	EmploymentType string
	Campus         string
}

type CourseFilter struct {
	Keyword  string
	Category string
	Status   string
}

type StudentFilter struct {
	Keyword string
	Status  string
	ClassID uint64
}

type ClassFilter struct {
	Keyword   string
	Status    string
	CourseID  uint64
	TeacherID uint64
	Scope     Scope
}

type ScheduleFilter struct {
	ClassID   uint64
	TeacherID uint64
	DateFrom  string
	DateTo    string
	Status    string
	Scope     Scope
}

type TeacherPayload struct {
	UserID         uint64
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
	RecentFeedbacks   []FeedbackItem          `json:"recentFeedbacks"`
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
	ID                uint64   `json:"id"`
	Title             string   `json:"title"`
	Content           string   `json:"content"`
	Category          string   `json:"category"`
	TargetScope       string   `json:"targetScope" gorm:"column:target_scope"`
	RelatedClassID    uint64   `json:"relatedClassId" gorm:"column:related_class_id"`
	RelatedScheduleID uint64   `json:"relatedScheduleId" gorm:"column:related_schedule_id"`
	StudentIDs        []uint64 `json:"studentIds"`
	Status            string   `json:"status"`
	PublishAt         string   `json:"publishAt"`
	Author            string   `json:"author"`
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
	StudentIDs        []uint64 `json:"studentIds"`
	Status            string   `json:"status"`
	Author            string   `json:"author"`
}

type NoticeFilter struct {
	ClassID    uint64
	Status     string
	NoticeType string
	Date       string
	DateFrom   string
	DateTo     string
	Scope      Scope
}

type AttendanceItem struct {
	RecordID     uint64 `json:"recordId"`
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
	DateFrom  string
	DateTo    string
	Status    string
	Scope     Scope
}

type AttendanceRecordItem struct {
	ID           uint64 `json:"id"`
	ScheduleID   uint64 `json:"scheduleId"`
	ClassID      uint64 `json:"classId"`
	ClassName    string `json:"className"`
	StudentID    uint64 `json:"studentId"`
	StudentName  string `json:"studentName"`
	TeacherID    uint64 `json:"teacherId"`
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

type AttendanceRecordUpdatePayload struct {
	Status string `json:"status"`
	Remark string `json:"remark"`
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

type HomeworkFilter struct {
	ClassID   uint64
	TeacherID uint64
	DateFrom  string
	DateTo    string
	Scope     Scope
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

type FeedbackFilter struct {
	ClassID   uint64
	TeacherID uint64
	DateFrom  string
	DateTo    string
	Scope     Scope
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
		&edumodel.NoticeTarget{},
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
	return s.OverviewWithScope(Scope{})
}

func (s *Service) OverviewWithScope(scope Scope) (map[string]any, error) {
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

	attendanceSessions, sessionErr := s.AttendanceSessionsWithScope(scope)
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

	var draftHomeworkCount int64
	draftHomeworkErr := s.db.Model(&edumodel.Homework{}).
		Where("status = ?", "draft").
		Count(&draftHomeworkCount).Error
	if draftHomeworkErr != nil {
		return nil, draftHomeworkErr
	}

	upcomingLessons, scheduleErr := s.SchedulesWithScope(scope)
	if scheduleErr != nil {
		return nil, scheduleErr
	}

	latestNotices, noticeErr := s.NoticesWithFilter(NoticeFilter{Scope: scope})
	if noticeErr != nil {
		return nil, noticeErr
	}

	if len(upcomingLessons) > 5 {
		upcomingLessons = upcomingLessons[:5]
	}

	if len(latestNotices) > 5 {
		latestNotices = latestNotices[:5]
	}

	pendingItems := make([]map[string]any, 0, 3)
	if todayPendingCheck > 0 {
		pendingItems = append(pendingItems, map[string]any{
			"key":   "attendance",
			"label": "待签到班级",
			"count": todayPendingCheck,
		})
	}
	if pendingActionCount > 0 {
		pendingItems = append(pendingItems, map[string]any{
			"key":   "notice",
			"label": "待发送通知",
			"count": pendingActionCount,
		})
	}
	if draftHomeworkCount > 0 {
		pendingItems = append(pendingItems, map[string]any{
			"key":   "homework",
			"label": "待补作业反馈",
			"count": draftHomeworkCount,
		})
	}

	return map[string]any{
		"todayCourses":         todayCourses,
		"todayPendingCheck":    todayPendingCheck,
		"todayLeaveCount":      todayLeaveCount,
		"todayAbsentCount":     todayAbsentCount,
		"studentCount":         studentCount,
		"classCount":           classCount,
		"pendingActionCount":   pendingActionCount,
		"pendingHomeworkCount": draftHomeworkCount,
		"pendingItems":         pendingItems,
		"upcomingLessons":      upcomingLessons,
		"latestNotices":        latestNotices,
	}, nil
}

func (s *Service) ScopeForUser(userID uint64, primaryRole string) (Scope, error) {
	scope := Scope{
		UserID:      userID,
		PrimaryRole: strings.TrimSpace(primaryRole),
	}

	if scope.UserID == 0 || scope.PrimaryRole != "teacher" {
		return scope, nil
	}

	teacherID, found, teacherErr := s.TeacherIDByUserID(scope.UserID)
	if teacherErr != nil {
		return Scope{}, teacherErr
	}
	if !found {
		return scope, nil
	}

	scope.TeacherID = teacherID
	scope.RestrictToSelf = true
	return scope, nil
}

func (s *Service) ScopeAllowsTeacher(scope Scope, teacherID uint64) bool {
	if !scope.RestrictToSelf {
		return true
	}

	return scope.TeacherID > 0 && scope.TeacherID == teacherID
}

func (s *Service) TeacherIDByUserID(userID uint64) (uint64, bool, error) {
	if userID == 0 {
		return 0, false, nil
	}

	if s.db == nil {
		if userID == 4 {
			return 1, true, nil
		}

		return 0, false, nil
	}

	type teacherUserRow struct {
		ID uint64 `gorm:"column:id"`
	}

	var row teacherUserRow
	findErr := s.db.Model(&edumodel.Teacher{}).
		Select("id").
		Where("user_id = ?", userID).
		Limit(1).
		Scan(&row).Error
	if findErr != nil {
		return 0, false, findErr
	}
	if row.ID == 0 {
		return 0, false, nil
	}

	return row.ID, true, nil
}

func (s *Service) Teachers() ([]TeacherItem, error) {
	return s.TeachersWithFilter(TeacherFilter{})
}

func (s *Service) TeachersWithFilter(filter TeacherFilter) ([]TeacherItem, error) {
	if s.db == nil {
		return s.teacherItemsFromDemoWithFilter(filter), nil
	}

	var items []TeacherItem
	listErr := s.teacherQuery(filter).
		Order("t.id ASC").
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
	findErr := s.teacherQuery(TeacherFilter{}).
		Where("t.id = ?", rawID).
		Limit(1).
		Scan(&item).Error
	if findErr != nil {
		return TeacherItem{}, false, findErr
	}

	if item.ID == 0 {
		return TeacherItem{}, false, nil
	}

	return item, true, nil
}

func (s *Service) CreateTeacher(input TeacherPayload, operator Operator) (TeacherItem, error) {
	if s.db == nil {
		return TeacherItem{
			ID:             1,
			UserID:         input.UserID,
			UserName:       demoTeacherUserName(input.UserID),
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

	bindErr := s.ensureTeacherUserBinding(input.UserID, 0)
	if bindErr != nil {
		return TeacherItem{}, bindErr
	}

	var createdTeacherID uint64
	createErr := s.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		teacher := edumodel.Teacher{
			UserID:         optionalUint64Pointer(input.UserID),
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

		createTeacherErr := tx.Create(&teacher).Error
		if createTeacherErr != nil {
			return createTeacherErr
		}

		createdTeacherID = teacher.ID
		return s.recordBusinessOperationTx(
			tx,
			operator,
			"teachers",
			"create",
			"teacher",
			teacher.ID,
			fmt.Sprintf("创建老师 %s。", teacher.Name),
		)
	})
	if createErr != nil {
		return TeacherItem{}, createErr
	}

	createdItem, found, teacherErr := s.Teacher(fmt.Sprintf("%d", createdTeacherID))
	if teacherErr != nil {
		return TeacherItem{}, teacherErr
	}
	if !found {
		return TeacherItem{ID: createdTeacherID}, nil
	}

	return createdItem, nil
}

func (s *Service) UpdateTeacher(rawID string, input TeacherPayload, operator Operator) (TeacherItem, bool, error) {
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
		item.UserID = input.UserID
		item.UserName = demoTeacherUserName(input.UserID)

		return item, true, nil
	}

	teacherID, parseErr := strconv.ParseUint(rawID, 10, 64)
	if parseErr != nil {
		return TeacherItem{}, false, nil
	}

	bindErr := s.ensureTeacherUserBinding(input.UserID, teacherID)
	if bindErr != nil {
		return TeacherItem{}, false, bindErr
	}

	updateErr := s.db.Transaction(func(tx *gorm.DB) error {
		var teacher edumodel.Teacher
		findTeacherErr := tx.Where("id = ?", teacherID).Take(&teacher).Error
		if findTeacherErr != nil {
			return findTeacherErr
		}

		updateResult := tx.Model(&edumodel.Teacher{}).
			Where("id = ?", teacherID).
			Updates(map[string]any{
				"user_id":         optionalUint64Pointer(input.UserID),
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
			return updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return s.recordBusinessOperationTx(
			tx,
			operator,
			"teachers",
			"update",
			"teacher",
			teacherID,
			fmt.Sprintf("更新老师 %s 的资料。", input.Name),
		)
	})
	if errors.Is(updateErr, gorm.ErrRecordNotFound) {
		return TeacherItem{}, false, nil
	}
	if updateErr != nil {
		return TeacherItem{}, false, updateErr
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

func (s *Service) CreateCourse(input CoursePayload, operator Operator) (CourseItem, error) {
	if s.db == nil {
		return courseItemFromPayload(1, input), nil
	}

	var createdCourseID uint64
	createErr := s.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		course := edumodel.Course{
			Name:                  input.Name,
			Category:              input.Category,
			Description:           input.Description,
			AgeRange:              input.AgeRange,
			LessonDurationMinutes: input.LessonDurationMinutes,
			TotalLessons:          input.TotalLessons,
			DeliveryType:          input.DeliveryType,
			Status:                input.Status,
			CreatedAt:             now,
			UpdatedAt:             now,
		}

		createCourseErr := tx.Create(&course).Error
		if createCourseErr != nil {
			return createCourseErr
		}

		createdCourseID = course.ID
		return s.recordBusinessOperationTx(
			tx,
			operator,
			"courses",
			"create",
			"course",
			course.ID,
			fmt.Sprintf("创建课程 %s。", course.Name),
		)
	})
	if createErr != nil {
		return CourseItem{}, createErr
	}

	createdItem, found, detailErr := s.Course(fmt.Sprintf("%d", createdCourseID))
	if detailErr != nil {
		return CourseItem{}, detailErr
	}
	if !found {
		return courseItemFromPayload(createdCourseID, input), nil
	}

	return createdItem, nil
}

func (s *Service) UpdateCourse(rawID string, input CoursePayload, operator Operator) (CourseItem, bool, error) {
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

	courseID, parseErr := strconv.ParseUint(rawID, 10, 64)
	if parseErr != nil {
		return CourseItem{}, false, nil
	}

	updateErr := s.db.Transaction(func(tx *gorm.DB) error {
		var course edumodel.Course
		findCourseErr := tx.Where("id = ?", courseID).Take(&course).Error
		if findCourseErr != nil {
			return findCourseErr
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
			"updated_at":              time.Now(),
		}

		updateResult := tx.Model(&edumodel.Course{}).
			Where("id = ?", courseID).
			Updates(updateValues)
		if updateResult.Error != nil {
			return updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return s.recordBusinessOperationTx(
			tx,
			operator,
			"courses",
			"update",
			"course",
			courseID,
			fmt.Sprintf("更新课程 %s。", input.Name),
		)
	})
	if errors.Is(updateErr, gorm.ErrRecordNotFound) {
		return CourseItem{}, false, nil
	}
	if updateErr != nil {
		return CourseItem{}, false, updateErr
	}

	updatedItem, found, detailErr := s.Course(rawID)
	if detailErr != nil {
		return CourseItem{}, false, detailErr
	}

	return updatedItem, found, nil
}

func (s *Service) Students() ([]StudentItem, error) {
	return s.StudentsWithFilter(StudentFilter{})
}

func (s *Service) StudentsWithFilter(filter StudentFilter) ([]StudentItem, error) {
	if s.db == nil {
		return s.studentItemsFromDemoWithFilter(filter), nil
	}

	var items []StudentItem
	listErr := s.studentQuery(filter).
		Order("s.id ASC").
		Scan(&items).Error
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

func (s *Service) CreateStudent(input StudentPayload, operator Operator) (StudentDetail, error) {
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
		return s.recordBusinessOperationTx(
			tx,
			operator,
			"students",
			"create",
			"student",
			student.ID,
			fmt.Sprintf("新增学员 %s。", student.Name),
		)
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

func (s *Service) UpdateStudent(rawID string, input StudentPayload, operator Operator) (StudentDetail, bool, error) {
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
			createGuardianErr := tx.Create(&guardian).Error
			if createGuardianErr != nil {
				return createGuardianErr
			}

			return s.recordBusinessOperationTx(
				tx,
				operator,
				"students",
				"update",
				"student",
				studentID,
				fmt.Sprintf("更新学员 %s。", input.Name),
			)
		}
		if findGuardianErr != nil {
			return findGuardianErr
		}

		updateGuardianErr := tx.Model(&edumodel.StudentGuardian{}).
			Where("id = ?", primaryGuardian.ID).
			Updates(map[string]any{
				"name":       input.GuardianName,
				"relation":   input.GuardianRelation,
				"mobile":     input.GuardianMobile,
				"is_primary": true,
				"updated_at": time.Now(),
			}).Error
		if updateGuardianErr != nil {
			return updateGuardianErr
		}

		return s.recordBusinessOperationTx(
			tx,
			operator,
			"students",
			"update",
			"student",
			studentID,
			fmt.Sprintf("更新学员 %s。", input.Name),
		)
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

func (s *Service) CreateStudentGuardian(rawStudentID string, input StudentGuardianPayload, operator Operator) (StudentGuardianItem, bool, error) {
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

		createGuardianErr := tx.Create(&createdGuardian).Error
		if createGuardianErr != nil {
			return createGuardianErr
		}

		return s.recordBusinessOperationTx(
			tx,
			operator,
			"students",
			"create_guardian",
			"guardian",
			createdGuardian.ID,
			fmt.Sprintf("为学员 %s 新增家长联系人 %s。", student.Name, createdGuardian.Name),
		)
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

func (s *Service) UpdateStudentGuardian(rawStudentID string, rawGuardianID string, input StudentGuardianPayload, operator Operator) (StudentGuardianItem, bool, error) {
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

		updateGuardianErr := tx.Model(&edumodel.StudentGuardian{}).
			Where("id = ?", guardianID).
			Updates(map[string]any{
				"name":       input.Name,
				"relation":   input.Relation,
				"mobile":     input.Mobile,
				"is_primary": input.IsPrimary,
				"updated_at": time.Now(),
			}).Error
		if updateGuardianErr != nil {
			return updateGuardianErr
		}

		studentItem, found, studentErr := s.Student(rawStudentID)
		if studentErr != nil {
			return studentErr
		}
		if !found {
			return gorm.ErrRecordNotFound
		}

		return s.recordBusinessOperationTx(
			tx,
			operator,
			"students",
			"update_guardian",
			"guardian",
			guardianID,
			fmt.Sprintf("更新学员 %s 的家长联系人 %s。", studentItem.Name, input.Name),
		)
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

func (s *Service) DeleteStudentGuardian(rawStudentID string, rawGuardianID string, operator Operator) (bool, error) {
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
		var removedGuardianName string
		for _, guardian := range guardians {
			if guardian.ID == guardianID {
				targetFound = true
				removedPrimary = guardian.IsPrimary
				removedGuardianName = guardian.Name
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

		studentItem, found, studentErr := s.Student(rawStudentID)
		if studentErr != nil {
			return studentErr
		}
		if !found {
			return gorm.ErrRecordNotFound
		}

		recordErr := s.recordBusinessOperationTx(
			tx,
			operator,
			"students",
			"delete_guardian",
			"guardian",
			guardianID,
			fmt.Sprintf("删除学员 %s 的家长联系人 %s。", studentItem.Name, removedGuardianName),
		)
		if recordErr != nil {
			return recordErr
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
	return s.ClassesWithFilter(ClassFilter{})
}

func (s *Service) ClassesWithFilter(filter ClassFilter) ([]ClassItem, error) {
	if s.db == nil {
		return s.classItemsFromDemoWithFilter(filter), nil
	}

	var items []ClassItem
	listErr := s.classQuery(filter).
		Order("c.id ASC").
		Scan(&items).Error
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

func (s *Service) CreateClass(input ClassPayload, operator Operator) (ClassItem, error) {
	if s.db == nil {
		return ClassItem{}, nil
	}

	startDate, endDate, parseErr := parseClassDates(input.StartDate, input.EndDate)
	if parseErr != nil {
		return ClassItem{}, parseErr
	}

	var createdClassID uint64
	createErr := s.db.Transaction(func(tx *gorm.DB) error {
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

		createClassErr := tx.Create(&classItem).Error
		if createClassErr != nil {
			return createClassErr
		}

		createdClassID = classItem.ID
		return s.recordBusinessOperationTx(
			tx,
			operator,
			"classes",
			"create",
			"class",
			classItem.ID,
			fmt.Sprintf("新建班级 %s。", classItem.Name),
		)
	})
	if createErr != nil {
		return ClassItem{}, createErr
	}

	createdItem, found, classErr := s.Class(fmt.Sprintf("%d", createdClassID))
	if classErr != nil {
		return ClassItem{}, classErr
	}
	if !found {
		return ClassItem{}, nil
	}

	return createdItem, nil
}

func (s *Service) UpdateClass(rawID string, input ClassPayload, operator Operator) (ClassItem, bool, error) {
	if s.db == nil {
		return ClassItem{}, false, nil
	}

	startDate, endDate, parseErr := parseClassDates(input.StartDate, input.EndDate)
	if parseErr != nil {
		return ClassItem{}, false, parseErr
	}

	classID, parseErr := strconv.ParseUint(rawID, 10, 64)
	if parseErr != nil {
		return ClassItem{}, false, nil
	}

	updateErr := s.db.Transaction(func(tx *gorm.DB) error {
		var classItem edumodel.Class
		findClassErr := tx.Where("id = ?", classID).Take(&classItem).Error
		if findClassErr != nil {
			return findClassErr
		}

		updateResult := tx.Model(&edumodel.Class{}).
			Where("id = ?", classID).
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
			return updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return s.recordBusinessOperationTx(
			tx,
			operator,
			"classes",
			"update",
			"class",
			classID,
			fmt.Sprintf("更新班级 %s。", input.Name),
		)
	})
	if errors.Is(updateErr, gorm.ErrRecordNotFound) {
		return ClassItem{}, false, nil
	}
	if updateErr != nil {
		return ClassItem{}, false, updateErr
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
	const recentItemLimit = 4

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

	classSchedules, scheduleErr := s.classSchedules(rawClassID)
	if scheduleErr != nil {
		return ClassDetail{}, false, scheduleErr
	}

	upcomingSchedules, upcomingErr := s.UpcomingSchedules(rawClassID)
	if upcomingErr != nil {
		return ClassDetail{}, false, upcomingErr
	}

	recentAttendance, attendanceErr := s.recentAttendanceSessions(classSchedules, recentItemLimit)
	if attendanceErr != nil {
		return ClassDetail{}, false, attendanceErr
	}

	recentHomeworks, homeworkErr := s.HomeworksWithFilter(HomeworkFilter{ClassID: classItem.ID})
	if homeworkErr != nil {
		return ClassDetail{}, false, homeworkErr
	}
	if len(recentHomeworks) > recentItemLimit {
		recentHomeworks = recentHomeworks[:recentItemLimit]
	}

	recentFeedbacks, feedbackErr := s.FeedbacksWithFilter(FeedbackFilter{ClassID: classItem.ID})
	if feedbackErr != nil {
		return ClassDetail{}, false, feedbackErr
	}
	if len(recentFeedbacks) > recentItemLimit {
		recentFeedbacks = recentFeedbacks[:recentItemLimit]
	}

	recentNotices, noticeErr := s.NoticesWithFilter(NoticeFilter{ClassID: classItem.ID})
	if noticeErr != nil {
		return ClassDetail{}, false, noticeErr
	}
	if len(recentNotices) > recentItemLimit {
		recentNotices = recentNotices[:recentItemLimit]
	}

	return ClassDetail{
		Class:             classItem,
		Students:          students,
		UpcomingSchedules: upcomingSchedules,
		RecentAttendance:  recentAttendance,
		RecentHomeworks:   recentHomeworks,
		RecentFeedbacks:   recentFeedbacks,
		RecentNotices:     recentNotices,
	}, true, nil
}

func (s *Service) AddStudentsToClass(rawClassID string, studentIDs []uint64, operator Operator) (bool, error) {
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

	saveErr := s.db.Transaction(func(tx *gorm.DB) error {
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

		createRelationErr := tx.Clauses(clause.OnConflict{
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
		if createRelationErr != nil {
			return createRelationErr
		}

		return s.recordBusinessOperationTx(
			tx,
			operator,
			"classes",
			"add_students",
			"class",
			classItem.ID,
			fmt.Sprintf("向班级 %s 添加 %d 名学员。", classItem.Name, len(studentIDs)),
		)
	})
	if saveErr != nil {
		return false, saveErr
	}

	return true, nil
}

func (s *Service) RemoveStudentFromClass(rawClassID string, rawStudentID string, operator Operator) (bool, error) {
	if s.db == nil {
		return true, nil
	}

	classItem, classFound, classErr := s.Class(rawClassID)
	if classErr != nil {
		return false, classErr
	}
	if !classFound {
		return false, nil
	}

	studentItem, studentFound, studentErr := s.Student(rawStudentID)
	if studentErr != nil {
		return false, studentErr
	}
	if !studentFound {
		return false, nil
	}

	updateErr := s.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		updateResult := tx.Model(&edumodel.ClassStudent{}).
			Where("class_id = ? AND student_id = ?", rawClassID, rawStudentID).
			Updates(map[string]any{
				"status":     "已移出",
				"leave_date": datePointer(startOfDay(now)),
				"updated_at": now,
			})
		if updateResult.Error != nil {
			return updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return s.recordBusinessOperationTx(
			tx,
			operator,
			"classes",
			"remove_student",
			"class",
			classItem.ID,
			fmt.Sprintf("将学员 %s 移出班级 %s。", studentItem.Name, classItem.Name),
		)
	})
	if errors.Is(updateErr, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if updateErr != nil {
		return false, updateErr
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
	return s.SchedulesWithFilter(ScheduleFilter{})
}

func (s *Service) SchedulesWithScope(scope Scope) ([]ScheduleItem, error) {
	return s.SchedulesWithFilter(ScheduleFilter{Scope: scope})
}

func (s *Service) SchedulesWithFilter(filter ScheduleFilter) ([]ScheduleItem, error) {
	if s.db == nil {
		return scheduleItemsFromDemoWithFilter(filter), nil
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
WHERE 1 = 1
`

	args := make([]any, 0, 5)
	if filter.ClassID > 0 {
		query += " AND s.class_id = ?"
		args = append(args, filter.ClassID)
	}
	if filter.TeacherID > 0 {
		query += " AND s.teacher_id = ?"
		args = append(args, filter.TeacherID)
	}
	if filter.Scope.RestrictToSelf && filter.Scope.TeacherID > 0 {
		query += " AND s.teacher_id = ?"
		args = append(args, filter.Scope.TeacherID)
	}
	filterDateFrom := strings.TrimSpace(filter.DateFrom)
	if filterDateFrom != "" {
		query += " AND s.schedule_date >= ?"
		args = append(args, filterDateFrom)
	}
	filterDateTo := strings.TrimSpace(filter.DateTo)
	if filterDateTo != "" {
		query += " AND s.schedule_date <= ?"
		args = append(args, filterDateTo)
	}
	filterStatus := strings.TrimSpace(filter.Status)
	if filterStatus != "" {
		query += " AND s.status = ?"
		args = append(args, filterStatus)
	}

	query += `
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
	listErr := s.db.Raw(query, args...).Scan(&rows).Error
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

func (s *Service) ClassAccessible(rawID string, scope Scope) (ClassItem, bool, error) {
	classItem, found, classErr := s.Class(rawID)
	if classErr != nil {
		return ClassItem{}, false, classErr
	}
	if !found || !s.ScopeAllowsTeacher(scope, classItem.TeacherID) {
		return ClassItem{}, false, nil
	}

	return classItem, true, nil
}

func (s *Service) ScheduleAccessible(rawID string, scope Scope) (ScheduleItem, bool, error) {
	scheduleItem, found, scheduleErr := s.Schedule(rawID)
	if scheduleErr != nil {
		return ScheduleItem{}, false, scheduleErr
	}
	if !found || !s.ScopeAllowsTeacher(scope, scheduleItem.TeacherID) {
		return ScheduleItem{}, false, nil
	}

	return scheduleItem, true, nil
}

func (s *Service) CreateSchedule(payload SchedulePayload, operator Operator) (ScheduleItem, error) {
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

	var createdScheduleID uint64
	createErr := s.db.Transaction(func(tx *gorm.DB) error {
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

		createScheduleErr := tx.Create(&record).Error
		if createScheduleErr != nil {
			return createScheduleErr
		}

		createdScheduleID = record.ID
		return s.recordBusinessOperationTx(
			tx,
			operator,
			"schedules",
			"create",
			"schedule",
			record.ID,
			fmt.Sprintf(
				"为班级 %s 新建%s排课，时间 %s %s。",
				classItem.Name,
				scheduleType,
				lessonDate.Format(dateLayout),
				formatLessonTime(payload.StartTime, payload.EndTime),
			),
		)
	})
	if createErr != nil {
		return ScheduleItem{}, createErr
	}

	createdItem, itemFound, detailErr := s.Schedule(fmt.Sprintf("%d", createdScheduleID))
	if detailErr != nil {
		return ScheduleItem{}, detailErr
	}
	if !itemFound {
		return ScheduleItem{}, nil
	}

	return createdItem, nil
}

func (s *Service) UpdateSchedule(rawID string, payload SchedulePayload, operator Operator) (ScheduleItem, bool, error) {
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

	scheduleID, scheduleParseErr := strconv.ParseUint(rawID, 10, 64)
	if scheduleParseErr != nil {
		return ScheduleItem{}, false, nil
	}

	updateErr := s.db.Transaction(func(tx *gorm.DB) error {
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

		updateResult := tx.Model(&edumodel.ClassSchedule{}).
			Where("id = ?", scheduleID).
			Updates(updateValues)
		if updateResult.Error != nil {
			return updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return s.recordBusinessOperationTx(
			tx,
			operator,
			"schedules",
			"update",
			"schedule",
			scheduleID,
			fmt.Sprintf(
				"更新班级 %s 的排课，时间调整为 %s %s。",
				classItem.Name,
				lessonDate.Format(dateLayout),
				formatLessonTime(payload.StartTime, payload.EndTime),
			),
		)
	})
	if errors.Is(updateErr, gorm.ErrRecordNotFound) {
		return ScheduleItem{}, false, nil
	}
	if updateErr != nil {
		return ScheduleItem{}, false, updateErr
	}

	updatedItem, itemFound, detailErr := s.Schedule(rawID)
	if detailErr != nil {
		return ScheduleItem{}, false, detailErr
	}

	return updatedItem, itemFound, nil
}

func (s *Service) Reschedule(rawID string, payload ScheduleActionPayload, operator Operator) (ScheduleItem, bool, error) {
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

	replacementItem, createErr := s.CreateSchedule(replacementPayload, Operator{})
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

		return s.recordBusinessOperationTx(
			tx,
			operator,
			"schedules",
			"reschedule",
			"schedule",
			scheduleItem.ID,
			fmt.Sprintf(
				"将班级 %s 的课程从 %s %s 调整到 %s %s。",
				scheduleItem.ClassName,
				scheduleItem.LessonDate,
				scheduleItem.LessonTime,
				replacementItem.LessonDate,
				replacementItem.LessonTime,
			),
		)
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

func (s *Service) CancelSchedule(rawID string, payload ScheduleActionPayload, operator Operator) (ScheduleItem, bool, error) {
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

	scheduleID, scheduleParseErr := strconv.ParseUint(rawID, 10, 64)
	if scheduleParseErr != nil {
		return ScheduleItem{}, false, nil
	}

	updateErr := s.db.Transaction(func(tx *gorm.DB) error {
		updateResult := tx.Model(&edumodel.ClassSchedule{}).
			Where("id = ?", scheduleID).
			Updates(map[string]any{
				"status":     "已停课",
				"remark":     strings.TrimSpace(payload.Remark),
				"updated_at": time.Now(),
			})
		if updateResult.Error != nil {
			return updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return s.recordBusinessOperationTx(
			tx,
			operator,
			"schedules",
			"cancel",
			"schedule",
			scheduleID,
			fmt.Sprintf("将班级 %s 的课程 %s %s 标记为停课。", scheduleItem.ClassName, scheduleItem.LessonDate, scheduleItem.LessonTime),
		)
	})
	if errors.Is(updateErr, gorm.ErrRecordNotFound) {
		return ScheduleItem{}, false, nil
	}
	if updateErr != nil {
		return ScheduleItem{}, false, updateErr
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

func (s *Service) CreateMakeupSchedule(rawID string, payload ScheduleActionPayload, operator Operator) (ScheduleItem, bool, error) {
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

	createdItem, createErr := s.CreateSchedule(makeupPayload, Operator{})
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
	sourceErr := s.db.Transaction(func(tx *gorm.DB) error {
		updateSourceErr := tx.Model(&edumodel.ClassSchedule{}).
			Where("id = ?", createdItem.ID).
			Updates(map[string]any{
				"source_schedule_id": sourceScheduleID,
				"status":             "待上课",
				"updated_at":         time.Now(),
			}).Error
		if updateSourceErr != nil {
			return updateSourceErr
		}

		return s.recordBusinessOperationTx(
			tx,
			operator,
			"schedules",
			"makeup",
			"schedule",
			createdItem.ID,
			fmt.Sprintf("为班级 %s 创建补课安排，时间 %s %s。", originalItem.ClassName, createdItem.LessonDate, createdItem.LessonTime),
		)
	})
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
	items, listErr := s.classSchedules(rawClassID)
	if listErr != nil {
		return nil, listErr
	}

	today := startOfDay(time.Now())
	filteredItems := make([]ScheduleItem, 0, len(items))
	for _, item := range items {
		lessonDate, parseErr := time.ParseInLocation(dateLayout, item.LessonDate, time.Local)
		if parseErr == nil && startOfDay(lessonDate).Before(today) {
			continue
		}

		filteredItems = append(filteredItems, item)
	}

	return filteredItems, nil
}

func (s *Service) classSchedules(rawClassID string) ([]ScheduleItem, error) {
	items, listErr := s.Schedules()
	if listErr != nil {
		return nil, listErr
	}

	filteredItems := make([]ScheduleItem, 0, len(items))
	for _, item := range items {
		if fmt.Sprintf("%d", item.ClassID) != rawClassID {
			continue
		}

		filteredItems = append(filteredItems, item)
	}

	return filteredItems, nil
}

func (s *Service) recentAttendanceSessions(items []ScheduleItem, limit int) ([]AttendanceSessionItem, error) {
	if len(items) == 0 || limit <= 0 {
		return []AttendanceSessionItem{}, nil
	}

	today := startOfDay(time.Now())
	recentItems := make([]AttendanceSessionItem, 0, limit)
	for index := len(items) - 1; index >= 0; index-- {
		scheduleItem := items[index]
		lessonDate, parseErr := time.ParseInLocation(dateLayout, scheduleItem.LessonDate, time.Local)
		if parseErr == nil && startOfDay(lessonDate).After(today) {
			continue
		}

		attendanceItems, attendanceErr := s.attendanceFromSchedule(scheduleItem)
		if attendanceErr != nil {
			return nil, attendanceErr
		}

		presentCount, leaveCount, absentCount, pendingCount := summarizeAttendanceItems(attendanceItems)
		recentItems = append(recentItems, AttendanceSessionItem{
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
		if len(recentItems) >= limit {
			break
		}
	}

	return recentItems, nil
}

func (s *Service) Attendance(rawScheduleID string) ([]AttendanceItem, error) {
	if s.db == nil {
		return attendanceItemsFromDemo(rawScheduleID, demo.Attendance(rawScheduleID)), nil
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
	return s.AttendanceSessionsWithScope(Scope{})
}

func (s *Service) AttendanceSessionsWithScope(scope Scope) ([]AttendanceSessionItem, error) {
	schedules, scheduleErr := s.SchedulesWithScope(scope)
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
	schedules, scheduleErr := s.SchedulesWithScope(filter.Scope)
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
		if !matchesDateRange(scheduleItem.LessonDate, filter.DateFrom, filter.DateTo) {
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

			items = append(items, attendanceRecordItemFromAttendanceItem(scheduleItem, attendanceItem))
		}
	}

	return items, nil
}

func (s *Service) AttendanceRecord(rawRecordID string) (AttendanceRecordItem, bool, error) {
	if s.db == nil {
		recordID, parseErr := strconv.ParseUint(strings.TrimSpace(rawRecordID), 10, 64)
		if parseErr != nil || recordID == 0 {
			return AttendanceRecordItem{}, false, nil
		}

		scheduleID, studentID, parsed := attendanceRecordPartsFromID(recordID)
		if !parsed {
			return AttendanceRecordItem{}, false, nil
		}

		return s.attendanceRecordByScheduleStudent(scheduleID, studentID)
	}

	query := `
SELECT
  ar.id,
  ar.schedule_id,
  s.class_id,
  COALESCE(c.name, '') AS class_name,
  ar.student_id,
  COALESCE(st.name, '') AS student_name,
  COALESCE(s.teacher_id, 0) AS teacher_id,
  COALESCE(t.name, '') AS teacher_name,
  COALESCE(DATE_FORMAT(s.schedule_date, '%Y-%m-%d'), '') AS lesson_date,
  COALESCE(s.start_time, '') AS start_time,
  COALESCE(s.end_time, '') AS end_time,
  COALESCE(ar.status, '') AS status,
  COALESCE(ar.remark, '') AS remark,
  COALESCE(ar.updated_by, '') AS updated_by,
  ar.checked_at AS checked_at,
  COALESCE(g.mobile, '') AS parent_mobile
FROM attendance_records AS ar
LEFT JOIN class_schedules AS s
  ON s.id = ar.schedule_id
LEFT JOIN classes AS c
  ON c.id = s.class_id
LEFT JOIN teachers AS t
  ON t.id = s.teacher_id
LEFT JOIN students AS st
  ON st.id = ar.student_id
LEFT JOIN student_guardians AS g
  ON g.student_id = st.id AND g.is_primary = 1
WHERE ar.id = ?
LIMIT 1
`

	type attendanceRecordRow struct {
		ID           uint64     `gorm:"column:id"`
		ScheduleID   uint64     `gorm:"column:schedule_id"`
		ClassID      uint64     `gorm:"column:class_id"`
		ClassName    string     `gorm:"column:class_name"`
		StudentID    uint64     `gorm:"column:student_id"`
		StudentName  string     `gorm:"column:student_name"`
		TeacherID    uint64     `gorm:"column:teacher_id"`
		TeacherName  string     `gorm:"column:teacher_name"`
		LessonDate   string     `gorm:"column:lesson_date"`
		StartTime    string     `gorm:"column:start_time"`
		EndTime      string     `gorm:"column:end_time"`
		Status       string     `gorm:"column:status"`
		Remark       string     `gorm:"column:remark"`
		UpdatedBy    string     `gorm:"column:updated_by"`
		CheckedAt    *time.Time `gorm:"column:checked_at"`
		ParentMobile string     `gorm:"column:parent_mobile"`
	}

	var row attendanceRecordRow
	findErr := s.db.Raw(query, rawRecordID).Scan(&row).Error
	if findErr != nil {
		return AttendanceRecordItem{}, false, findErr
	}
	if row.ID > 0 {
		return AttendanceRecordItem{
			ID:           row.ID,
			ScheduleID:   row.ScheduleID,
			ClassID:      row.ClassID,
			ClassName:    row.ClassName,
			StudentID:    row.StudentID,
			StudentName:  row.StudentName,
			TeacherID:    row.TeacherID,
			TeacherName:  row.TeacherName,
			LessonDate:   row.LessonDate,
			LessonTime:   formatLessonTime(row.StartTime, row.EndTime),
			Status:       normalizeAttendanceItemStatus(row.Status),
			Remark:       row.Remark,
			UpdatedBy:    row.UpdatedBy,
			UpdatedAt:    formatDateTime(row.CheckedAt),
			ParentMobile: row.ParentMobile,
		}, true, nil
	}

	recordID, parseErr := strconv.ParseUint(strings.TrimSpace(rawRecordID), 10, 64)
	if parseErr != nil || recordID == 0 {
		return AttendanceRecordItem{}, false, nil
	}

	scheduleID, studentID, parsed := attendanceRecordPartsFromID(recordID)
	if !parsed {
		return AttendanceRecordItem{}, false, nil
	}

	return s.attendanceRecordByScheduleStudent(scheduleID, studentID)
}

func (s *Service) UpdateAttendanceRecord(rawRecordID string, payload AttendanceRecordUpdatePayload, operator Operator) (AttendanceRecordItem, bool, error) {
	nextStatus := normalizeAttendanceItemStatus(payload.Status)
	if nextStatus == "" {
		return AttendanceRecordItem{}, false, ErrInvalidAttendanceStatus
	}

	recordItem, found, recordErr := s.AttendanceRecord(rawRecordID)
	if recordErr != nil {
		return AttendanceRecordItem{}, false, recordErr
	}
	if !found {
		return AttendanceRecordItem{}, false, nil
	}

	trimmedRemark := strings.TrimSpace(payload.Remark)
	now := time.Now()
	editorName := operationDisplayName(operator, "")

	if s.db == nil {
		scheduleItem, scheduleFound, scheduleErr := s.Schedule(fmt.Sprintf("%d", recordItem.ScheduleID))
		if scheduleErr != nil {
			return AttendanceRecordItem{}, false, scheduleErr
		}
		if !scheduleFound {
			return AttendanceRecordItem{}, false, nil
		}

		currentItems, currentErr := s.attendanceFromSchedule(scheduleItem)
		if currentErr != nil {
			return AttendanceRecordItem{}, false, currentErr
		}

		demoItems := make([]demo.AttendanceItem, 0, len(currentItems))
		itemFound := false
		for _, item := range currentItems {
			nextItem := item
			if item.RecordID == recordItem.ID {
				nextItem.Status = nextStatus
				nextItem.Remark = trimmedRemark
				nextItem.UpdatedBy = editorName
				nextItem.UpdatedAt = now.Format(dateTimeLayout)
				itemFound = true
			}

			demoItems = append(demoItems, demo.AttendanceItem{
				StudentID:    int(nextItem.StudentID),
				StudentName:  nextItem.StudentName,
				Grade:        nextItem.Grade,
				ParentMobile: nextItem.ParentMobile,
				Status:       nextItem.Status,
				Remark:       nextItem.Remark,
				UpdatedBy:    nextItem.UpdatedBy,
				UpdatedAt:    nextItem.UpdatedAt,
			})
		}
		if !itemFound {
			return AttendanceRecordItem{}, false, nil
		}
		if !demo.SaveAttendance(fmt.Sprintf("%d", recordItem.ScheduleID), demoItems) {
			return AttendanceRecordItem{}, false, nil
		}

		return s.attendanceRecordByScheduleStudent(recordItem.ScheduleID, recordItem.StudentID)
	}

	scheduleItem, scheduleFound, scheduleErr := s.Schedule(fmt.Sprintf("%d", recordItem.ScheduleID))
	if scheduleErr != nil {
		return AttendanceRecordItem{}, false, scheduleErr
	}
	if !scheduleFound {
		return AttendanceRecordItem{}, false, nil
	}

	currentItems, currentErr := s.attendanceFromSchedule(scheduleItem)
	if currentErr != nil {
		return AttendanceRecordItem{}, false, currentErr
	}

	itemFound := false
	for index := range currentItems {
		if currentItems[index].RecordID != recordItem.ID {
			continue
		}

		currentItems[index].Status = nextStatus
		currentItems[index].Remark = trimmedRemark
		currentItems[index].UpdatedBy = editorName
		currentItems[index].UpdatedAt = now.Format(dateTimeLayout)
		itemFound = true
		break
	}
	if !itemFound {
		return AttendanceRecordItem{}, false, nil
	}

	_, _, _, pendingCount := summarizeAttendanceItems(currentItems)
	sessionStatus := nextSavedAttendanceStatus(pendingCount)
	targetRecordID := recordItem.ID

	transactionErr := s.db.Transaction(func(tx *gorm.DB) error {
		record := edumodel.AttendanceRecord{
			ScheduleID: recordItem.ScheduleID,
			StudentID:  recordItem.StudentID,
			Status:     nextStatus,
			Remark:     trimmedRemark,
			CheckedAt:  &now,
			UpdatedBy:  editorName,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		saveErr := tx.Clauses(clause.OnConflict{
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
		}).Create(&record).Error
		if saveErr != nil {
			return saveErr
		}
		if record.ID > 0 {
			targetRecordID = record.ID
		}

		sessionUpdateErr := tx.Model(&edumodel.ClassSchedule{}).
			Where("id = ?", recordItem.ScheduleID).
			Updates(map[string]any{
				"status":     sessionStatus,
				"updated_at": now,
			}).Error
		if sessionUpdateErr != nil {
			return sessionUpdateErr
		}

		content := fmt.Sprintf(
			"更正班级 %s 在 %s %s 的学员 %s 签到结果为 %s。",
			recordItem.ClassName,
			recordItem.LessonDate,
			recordItem.LessonTime,
			recordItem.StudentName,
			nextStatus,
		)

		return s.recordBusinessOperationTx(tx, operator, "attendance", "update_record", "attendance_record", targetRecordID, content)
	})
	if errors.Is(transactionErr, gorm.ErrRecordNotFound) {
		return AttendanceRecordItem{}, false, nil
	}
	if transactionErr != nil {
		return AttendanceRecordItem{}, false, transactionErr
	}

	return s.attendanceRecordByScheduleStudent(recordItem.ScheduleID, recordItem.StudentID)
}

func (s *Service) attendanceRecordByScheduleStudent(scheduleID uint64, studentID uint64) (AttendanceRecordItem, bool, error) {
	if scheduleID == 0 || studentID == 0 {
		return AttendanceRecordItem{}, false, nil
	}

	scheduleItem, found, scheduleErr := s.Schedule(fmt.Sprintf("%d", scheduleID))
	if scheduleErr != nil {
		return AttendanceRecordItem{}, false, scheduleErr
	}
	if !found {
		return AttendanceRecordItem{}, false, nil
	}

	attendanceItems, attendanceErr := s.attendanceFromSchedule(scheduleItem)
	if attendanceErr != nil {
		return AttendanceRecordItem{}, false, attendanceErr
	}

	for _, attendanceItem := range attendanceItems {
		if attendanceItem.StudentID != studentID {
			continue
		}

		return attendanceRecordItemFromAttendanceItem(scheduleItem, attendanceItem), true, nil
	}

	return AttendanceRecordItem{}, false, nil
}

func (s *Service) SaveAttendance(rawScheduleID string, payload AttendanceSavePayload, updatedBy string, operator Operator) (bool, error) {
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
	editorName := operationDisplayName(operator, updatedBy)
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
		var existingRecordCount int64
		countErr := tx.Model(&edumodel.AttendanceRecord{}).
			Where("schedule_id = ?", scheduleID).
			Count(&existingRecordCount).Error
		if countErr != nil {
			return countErr
		}

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

		action := "create"
		content := fmt.Sprintf("保存班级 %s 在 %s %s 的签到结果。", scheduleItem.ClassName, scheduleItem.LessonDate, scheduleItem.LessonTime)
		if existingRecordCount > 0 {
			action = "update"
			content = fmt.Sprintf("更新班级 %s 在 %s %s 的签到结果。", scheduleItem.ClassName, scheduleItem.LessonDate, scheduleItem.LessonTime)
		}

		return s.recordBusinessOperationTx(tx, operator, "attendance", action, "schedule", scheduleID, content)
	})
	if transactionErr != nil {
		return false, transactionErr
	}

	return true, nil
}

func (s *Service) attendanceFromSchedule(scheduleItem ScheduleItem) ([]AttendanceItem, error) {
	if s.db == nil {
		return attendanceItemsFromDemo(fmt.Sprintf("%d", scheduleItem.ID), demo.Attendance(fmt.Sprintf("%d", scheduleItem.ID))), nil
	}

	query := `
SELECT
  COALESCE(ar.id, 0) AS attendance_record_id,
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
		AttendanceRecordID uint64     `gorm:"column:attendance_record_id"`
		StudentID          uint64     `gorm:"column:student_id"`
		StudentName        string     `gorm:"column:student_name"`
		Grade              string     `gorm:"column:grade"`
		ParentMobile       string     `gorm:"column:parent_mobile"`
		AttendanceStatus   string     `gorm:"column:attendance_status"`
		AttendanceRemark   string     `gorm:"column:attendance_remark"`
		UpdatedBy          string     `gorm:"column:updated_by"`
		CheckedAt          *time.Time `gorm:"column:checked_at"`
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

		recordID := row.AttendanceRecordID
		if recordID == 0 {
			recordID = attendanceRecordIDFromValues(scheduleItem.ID, row.StudentID)
		}

		items = append(items, AttendanceItem{
			RecordID:     recordID,
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
	return s.HomeworksWithFilter(HomeworkFilter{})
}

func (s *Service) HomeworksWithFilter(filter HomeworkFilter) ([]HomeworkItem, error) {
	if s.db == nil {
		return s.homeworkItemsFromDemoWithFilter(filter), nil
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
WHERE 1 = 1
`

	args := make([]any, 0, 4)
	if filter.ClassID > 0 {
		query += " AND h.class_id = ?"
		args = append(args, filter.ClassID)
	}
	if filter.TeacherID > 0 {
		query += " AND s.teacher_id = ?"
		args = append(args, filter.TeacherID)
	}
	if filter.Scope.RestrictToSelf && filter.Scope.TeacherID > 0 {
		query += " AND s.teacher_id = ?"
		args = append(args, filter.Scope.TeacherID)
	}
	if strings.TrimSpace(filter.DateFrom) != "" {
		query += " AND s.schedule_date >= ?"
		args = append(args, strings.TrimSpace(filter.DateFrom))
	}
	if strings.TrimSpace(filter.DateTo) != "" {
		query += " AND s.schedule_date <= ?"
		args = append(args, strings.TrimSpace(filter.DateTo))
	}

	query += `
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
	listErr := s.db.Raw(query, args...).Scan(&rows).Error
	if listErr != nil {
		return nil, listErr
	}

	items := make([]HomeworkItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, HomeworkItem{
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
		})
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

func (s *Service) SaveHomework(rawScheduleID string, payload HomeworkPayload, operator Operator) (HomeworkItem, bool, error) {
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

	saveErr := s.db.Transaction(func(tx *gorm.DB) error {
		var existingHomework edumodel.Homework
		findHomeworkErr := tx.Model(&edumodel.Homework{}).
			Select("id").
			Where("schedule_id = ?", scheduleID).
			Limit(1).
			Scan(&existingHomework).Error
		if findHomeworkErr != nil {
			return findHomeworkErr
		}

		now := time.Now()
		record := edumodel.Homework{
			ScheduleID:      scheduleID,
			ClassID:         scheduleItem.ClassID,
			Title:           trimmedTitle,
			Content:         strings.TrimSpace(payload.Content),
			SubmissionNote:  strings.TrimSpace(payload.SubmissionNote),
			CreatedByUserID: operationUserID(operator),
			Status:          normalizedStatus,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		createHomeworkErr := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "schedule_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"title",
				"content",
				"submission_note",
				"status",
				"updated_at",
			}),
		}).Create(&record).Error
		if createHomeworkErr != nil {
			return createHomeworkErr
		}

		action := "create"
		targetID := record.ID
		content := fmt.Sprintf("为班级 %s 保存作业《%s》。", scheduleItem.ClassName, trimmedTitle)
		if existingHomework.ID > 0 {
			action = "update"
			targetID = existingHomework.ID
			content = fmt.Sprintf("更新班级 %s 的作业《%s》。", scheduleItem.ClassName, trimmedTitle)
		}

		return s.recordBusinessOperationTx(tx, operator, "homeworks", action, "homework", targetID, content)
	})
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
	return s.FeedbacksWithFilter(FeedbackFilter{})
}

func (s *Service) FeedbacksWithFilter(filter FeedbackFilter) ([]FeedbackItem, error) {
	if s.db == nil {
		return s.feedbackItemsFromDemoWithFilter(filter), nil
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
WHERE 1 = 1
`

	args := make([]any, 0, 4)
	if filter.ClassID > 0 {
		query += " AND f.class_id = ?"
		args = append(args, filter.ClassID)
	}
	if filter.TeacherID > 0 {
		query += " AND s.teacher_id = ?"
		args = append(args, filter.TeacherID)
	}
	if filter.Scope.RestrictToSelf && filter.Scope.TeacherID > 0 {
		query += " AND s.teacher_id = ?"
		args = append(args, filter.Scope.TeacherID)
	}
	if strings.TrimSpace(filter.DateFrom) != "" {
		query += " AND s.schedule_date >= ?"
		args = append(args, strings.TrimSpace(filter.DateFrom))
	}
	if strings.TrimSpace(filter.DateTo) != "" {
		query += " AND s.schedule_date <= ?"
		args = append(args, strings.TrimSpace(filter.DateTo))
	}

	query += `
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
	listErr := s.db.Raw(query, args...).Scan(&rows).Error
	if listErr != nil {
		return nil, listErr
	}

	items := make([]FeedbackItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, FeedbackItem{
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
		})
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

func (s *Service) SaveFeedback(rawScheduleID string, payload FeedbackPayload, operator Operator) (FeedbackItem, bool, error) {
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

	saveErr := s.db.Transaction(func(tx *gorm.DB) error {
		var existingFeedback edumodel.ClassFeedback
		findFeedbackErr := tx.Model(&edumodel.ClassFeedback{}).
			Select("id").
			Where("schedule_id = ?", scheduleID).
			Limit(1).
			Scan(&existingFeedback).Error
		if findFeedbackErr != nil {
			return findFeedbackErr
		}

		now := time.Now()
		record := edumodel.ClassFeedback{
			ScheduleID:      scheduleID,
			ClassID:         scheduleItem.ClassID,
			Summary:         strings.TrimSpace(payload.Summary),
			LearningStatus:  strings.TrimSpace(payload.LearningStatus),
			NextSuggestion:  strings.TrimSpace(payload.NextSuggestion),
			ParentNotice:    strings.TrimSpace(payload.ParentNotice),
			CreatedByUserID: operationUserID(operator),
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		createFeedbackErr := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "schedule_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"summary",
				"learning_status",
				"next_suggestion",
				"parent_notice",
				"updated_at",
			}),
		}).Create(&record).Error
		if createFeedbackErr != nil {
			return createFeedbackErr
		}

		action := "create"
		targetID := record.ID
		content := fmt.Sprintf("为班级 %s 保存课后反馈。", scheduleItem.ClassName)
		if existingFeedback.ID > 0 {
			action = "update"
			targetID = existingFeedback.ID
			content = fmt.Sprintf("更新班级 %s 的课后反馈。", scheduleItem.ClassName)
		}

		return s.recordBusinessOperationTx(tx, operator, "feedbacks", action, "feedback", targetID, content)
	})
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
WHERE 1 = 1
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
	queryArgs := make([]any, 0, 6)
	if filter.ClassID > 0 {
		query += " AND COALESCE(related_class_id, 0) = ?"
		queryArgs = append(queryArgs, filter.ClassID)
	}
	query, queryArgs = appendNoticeScopeFilter(query, queryArgs, filter.Scope)

	filterDate := strings.TrimSpace(filter.Date)
	filterStatus := strings.TrimSpace(filter.Status)
	filterNoticeType := strings.TrimSpace(filter.NoticeType)
	filterDateFrom := strings.TrimSpace(filter.DateFrom)
	filterDateTo := strings.TrimSpace(filter.DateTo)
	if filterStatus != "" {
		query += " AND status = ?"
		queryArgs = append(queryArgs, filterStatus)
	}
	if filterNoticeType != "" {
		query += " AND notice_type = ?"
		queryArgs = append(queryArgs, filterNoticeType)
	}
	if filterDate != "" {
		query += " AND DATE(COALESCE(publish_at, created_at)) = ?"
		queryArgs = append(queryArgs, filterDate)
	}
	if filterDateFrom != "" {
		query += " AND DATE(COALESCE(publish_at, created_at)) >= ?"
		queryArgs = append(queryArgs, filterDateFrom)
	}
	if filterDateTo != "" {
		query += " AND DATE(COALESCE(publish_at, created_at)) <= ?"
		queryArgs = append(queryArgs, filterDateTo)
	}
	query += " ORDER BY COALESCE(publish_at, created_at) DESC, id DESC"

	listErr := s.db.Raw(
		query,
		queryArgs...,
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
			StudentIDs:        []uint64{},
			Status:            row.Status,
			PublishAt:         publishAt,
			Author:            row.Author,
		})
	}

	return items, nil
}

func (s *Service) Notice(rawID string) (NoticeItem, bool, error) {
	if s.db == nil {
		rawItem, found := demo.FindNotice(rawID)
		if !found {
			return NoticeItem{}, false, nil
		}

		return NoticeItem{
			ID:                uint64(rawItem.ID),
			Title:             rawItem.Title,
			Content:           rawItem.Content,
			Category:          rawItem.Category,
			TargetScope:       rawItem.TargetScope,
			RelatedClassID:    uint64(rawItem.RelatedClassID),
			RelatedScheduleID: uint64(rawItem.RelatedScheduleID),
			StudentIDs:        noticeStudentIDsFromDemo(rawItem.StudentIDs),
			Status:            rawItem.Status,
			PublishAt:         rawItem.PublishAt,
			Author:            rawItem.Author,
		}, true, nil
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
WHERE id = ?
LIMIT 1
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

	var row noticeRow
	findErr := s.db.Raw(query, rawID).Scan(&row).Error
	if findErr != nil {
		return NoticeItem{}, false, findErr
	}
	if row.ID == 0 {
		return NoticeItem{}, false, nil
	}

	studentIDs, studentErr := s.noticeStudentIDs(rawID)
	if studentErr != nil {
		return NoticeItem{}, false, studentErr
	}

	publishAt := formatDateTime(row.PublishAt)
	if publishAt == "" {
		publishAt = row.CreatedAt.Format(dateTimeLayout)
	}

	return NoticeItem{
		ID:                row.ID,
		Title:             row.Title,
		Content:           row.Content,
		Category:          row.Category,
		TargetScope:       row.TargetScope,
		RelatedClassID:    row.RelatedClassID,
		RelatedScheduleID: row.RelatedScheduleID,
		StudentIDs:        studentIDs,
		Status:            row.Status,
		PublishAt:         publishAt,
		Author:            row.Author,
	}, true, nil
}

func (s *Service) NoticeAccessible(rawNoticeID string, scope Scope) (bool, error) {
	if !scope.RestrictToSelf || scope.TeacherID == 0 {
		return true, nil
	}

	if s.db == nil {
		noticeItem, found, noticeErr := s.Notice(rawNoticeID)
		if noticeErr != nil {
			return false, noticeErr
		}
		if !found {
			return false, nil
		}

		return s.noticeMatchesScope(scope, noticeItem)
	}

	query := `
SELECT id
FROM notices
WHERE id = ?
`

	queryArgs := []any{rawNoticeID}
	query, queryArgs = appendNoticeScopeFilter(query, queryArgs, scope)
	query += " LIMIT 1"

	type noticeAccessRow struct {
		ID uint64 `gorm:"column:id"`
	}

	var row noticeAccessRow
	findErr := s.db.Raw(query, queryArgs...).Scan(&row).Error
	if findErr != nil {
		return false, findErr
	}

	return row.ID > 0, nil
}

func (s *Service) NoticeTargets(rawNoticeID string) ([]NoticeTargetItem, error) {
	if s.db == nil {
		return noticeTargetItemsFromDemo(demo.NoticeTargets(rawNoticeID)), nil
	}

	targetItems, targetErr := s.noticeTargetsFromRelations(rawNoticeID)
	if targetErr != nil {
		return nil, targetErr
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

	if row.ScopeName == "" && row.ClassName == "" && row.ScheduleID == 0 && len(targetItems) == 0 {
		return []NoticeTargetItem{}, nil
	}

	items := make([]NoticeTargetItem, 0, len(targetItems)+2)
	items = append(items, targetItems...)

	hasClassTarget := false
	hasScheduleTarget := false
	for _, item := range targetItems {
		if item.Type == "关联班级" {
			hasClassTarget = true
		}
		if item.Type == "关联课程安排" {
			hasScheduleTarget = true
		}
	}

	if !hasClassTarget && row.ClassName != "" {
		items = append(items, NoticeTargetItem{
			Name:   row.ClassName,
			Type:   "关联班级",
			Campus: row.Campus,
		})
	}
	if len(targetItems) == 0 && row.ScopeName != "" {
		items = append(items, NoticeTargetItem{
			Name:   row.ScopeName,
			Type:   "通知范围",
			Campus: row.Campus,
		})
	}
	if !hasScheduleTarget && row.ScheduleID > 0 && row.ScheduleDate != nil {
		items = append(items, NoticeTargetItem{
			Name:   fmt.Sprintf("%s %s", row.ScheduleDate.Format(dateLayout), formatLessonTime(row.StartTime, row.EndTime)),
			Type:   "关联课程安排",
			Campus: row.Campus,
		})
	}

	return items, nil
}

func (s *Service) noticeTargetsFromRelations(rawNoticeID string) ([]NoticeTargetItem, error) {
	query := `
SELECT
  nt.target_type,
  COALESCE(c.name, '') AS class_name,
  COALESCE(c.campus, '') AS class_campus,
  COALESCE(st.name, '') AS student_name,
  COALESCE(st.campus, '') AS student_campus
FROM notice_targets AS nt
LEFT JOIN classes AS c
  ON c.id = nt.class_id
LEFT JOIN students AS st
  ON st.id = nt.student_id
WHERE nt.notice_id = ?
ORDER BY nt.id ASC
`

	type noticeRelationRow struct {
		TargetType    string `gorm:"column:target_type"`
		ClassName     string `gorm:"column:class_name"`
		ClassCampus   string `gorm:"column:class_campus"`
		StudentName   string `gorm:"column:student_name"`
		StudentCampus string `gorm:"column:student_campus"`
	}

	var rows []noticeRelationRow
	listErr := s.db.Raw(query, rawNoticeID).Scan(&rows).Error
	if listErr != nil {
		if strings.Contains(strings.ToLower(listErr.Error()), "doesn't exist") {
			return []NoticeTargetItem{}, nil
		}
		return nil, listErr
	}

	items := make([]NoticeTargetItem, 0, len(rows))
	for _, row := range rows {
		switch row.TargetType {
		case "class":
			if row.ClassName == "" {
				continue
			}
			items = append(items, NoticeTargetItem{
				Name:   row.ClassName,
				Type:   "关联班级",
				Campus: row.ClassCampus,
			})
		case "student":
			if row.StudentName == "" {
				continue
			}
			items = append(items, NoticeTargetItem{
				Name:   row.StudentName,
				Type:   "指定学员",
				Campus: row.StudentCampus,
			})
		}
	}

	return items, nil
}

func (s *Service) CreateNotice(input NoticePayload, operator Operator) (NoticeItem, error) {
	if len(input.StudentIDs) > 0 {
		studentNames, studentErr := s.noticeStudentNames(input.StudentIDs)
		if studentErr != nil {
			return NoticeItem{}, studentErr
		}
		if strings.TrimSpace(input.TargetScope) == "" || strings.TrimSpace(input.TargetScope) == "指定学员" {
			input.TargetScope = buildStudentTargetScope(studentNames)
		}
	}

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
			StudentIDs:        cloneNoticeStudentIDs(input.StudentIDs),
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
			StudentIDs:        demoNoticeStudentIDs(input.StudentIDs),
			Status:            item.Status,
			PublishAt:         item.PublishAt,
			Author:            item.Author,
		})

		return item, nil
	}

	var createdNoticeID uint64
	createErr := s.db.Transaction(func(tx *gorm.DB) error {
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

		createNoticeErr := tx.Create(&record).Error
		if createNoticeErr != nil {
			return createNoticeErr
		}

		targetSaveErr := s.replaceNoticeTargetsTx(tx, record.ID, input)
		if targetSaveErr != nil {
			return targetSaveErr
		}

		createdNoticeID = record.ID
		return s.recordBusinessOperationTx(
			tx,
			operator,
			"notices",
			"create",
			"notice",
			record.ID,
			fmt.Sprintf("创建通知《%s》。", input.Title),
		)
	})
	if createErr != nil {
		return NoticeItem{}, createErr
	}

	createdItem, found, detailErr := s.Notice(fmt.Sprintf("%d", createdNoticeID))
	if detailErr != nil {
		return NoticeItem{}, detailErr
	}
	if !found {
		return NoticeItem{
			ID:                createdNoticeID,
			Title:             input.Title,
			Content:           input.Content,
			Category:          input.Category,
			TargetScope:       input.TargetScope,
			RelatedClassID:    input.RelatedClassID,
			RelatedScheduleID: input.RelatedScheduleID,
			StudentIDs:        cloneNoticeStudentIDs(input.StudentIDs),
			Status:            input.Status,
			PublishAt:         "",
			Author:            input.Author,
		}, nil
	}

	return createdItem, nil
}

func (s *Service) UpdateNotice(rawID string, input NoticePayload, operator Operator) (NoticeItem, bool, error) {
	if len(input.StudentIDs) > 0 {
		studentNames, studentErr := s.noticeStudentNames(input.StudentIDs)
		if studentErr != nil {
			return NoticeItem{}, false, studentErr
		}
		if strings.TrimSpace(input.TargetScope) == "" || strings.TrimSpace(input.TargetScope) == "指定学员" {
			input.TargetScope = buildStudentTargetScope(studentNames)
		}
	}

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
			StudentIDs:        cloneNoticeStudentIDs(input.StudentIDs),
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
			StudentIDs:        demoNoticeStudentIDs(input.StudentIDs),
			Status:            item.Status,
			PublishAt:         item.PublishAt,
			Author:            item.Author,
		})

		return item, true, nil
	}

	noticeID, parseErr := strconv.ParseUint(rawID, 10, 64)
	if parseErr != nil {
		return NoticeItem{}, false, nil
	}

	updateErr := s.db.Transaction(func(tx *gorm.DB) error {
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

		updateResult := tx.Model(&edumodel.Notice{}).
			Where("id = ?", noticeID).
			Updates(updateValues)
		if updateResult.Error != nil {
			return updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		targetSaveErr := s.replaceNoticeTargetsTx(tx, noticeID, input)
		if targetSaveErr != nil {
			return targetSaveErr
		}

		return s.recordBusinessOperationTx(
			tx,
			operator,
			"notices",
			"update",
			"notice",
			noticeID,
			fmt.Sprintf("更新通知《%s》。", input.Title),
		)
	})
	if errors.Is(updateErr, gorm.ErrRecordNotFound) {
		return NoticeItem{}, false, nil
	}
	if updateErr != nil {
		return NoticeItem{}, false, updateErr
	}

	updatedItem, found, detailErr := s.Notice(rawID)
	if detailErr != nil {
		return NoticeItem{}, false, detailErr
	}

	return updatedItem, found, nil
}

func (s *Service) SendNotice(rawID string, operator Operator) (NoticeItem, bool, error) {
	item, found, itemErr := s.Notice(rawID)
	if itemErr != nil || !found {
		return item, found, itemErr
	}

	if s.db == nil {
		input := NoticePayload{
			Title:             item.Title,
			Content:           item.Content,
			Category:          item.Category,
			TargetScope:       item.TargetScope,
			RelatedClassID:    item.RelatedClassID,
			RelatedScheduleID: item.RelatedScheduleID,
			StudentIDs:        cloneNoticeStudentIDs(item.StudentIDs),
			Status:            "已发送",
			Author:            item.Author,
		}

		return s.UpdateNotice(rawID, input, Operator{})
	}

	noticeID, parseErr := strconv.ParseUint(rawID, 10, 64)
	if parseErr != nil {
		return NoticeItem{}, false, nil
	}

	sendErr := s.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		updateResult := tx.Model(&edumodel.Notice{}).
			Where("id = ?", noticeID).
			Updates(map[string]any{
				"status":     "已发送",
				"publish_at": now,
				"updated_at": now,
			})
		if updateResult.Error != nil {
			return updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return s.recordBusinessOperationTx(
			tx,
			operator,
			"notices",
			"send",
			"notice",
			noticeID,
			fmt.Sprintf("发送通知《%s》。", item.Title),
		)
	})
	if errors.Is(sendErr, gorm.ErrRecordNotFound) {
		return NoticeItem{}, false, nil
	}
	if sendErr != nil {
		return NoticeItem{}, false, sendErr
	}

	return s.Notice(rawID)
}

func (s *Service) createScheduleNotice(rawScheduleID string, input NoticePayload) error {
	noticeItem, noticeFound, noticeErr := s.findNoticeByScheduleAndCategory(rawScheduleID, input.Category)
	if noticeErr != nil {
		return noticeErr
	}
	if noticeFound {
		_, _, updateErr := s.UpdateNotice(fmt.Sprintf("%d", noticeItem.ID), input, Operator{})
		return updateErr
	}

	_, createErr := s.CreateNotice(input, Operator{})
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

func buildStudentTargetScope(studentNames []string) string {
	if len(studentNames) == 0 {
		return "指定学员家长"
	}

	if len(studentNames) == 1 {
		return fmt.Sprintf("%s家长", studentNames[0])
	}

	if len(studentNames) == 2 {
		return fmt.Sprintf("%s、%s家长", studentNames[0], studentNames[1])
	}

	return fmt.Sprintf("%s等%d位学员家长", studentNames[0], len(studentNames))
}

func (s *Service) noticeStudentNames(studentIDs []uint64) ([]string, error) {
	if len(studentIDs) == 0 {
		return []string{}, nil
	}

	if s.db == nil {
		names := make([]string, 0, len(studentIDs))
		for _, studentID := range studentIDs {
			item, found := demo.FindStudent(fmt.Sprintf("%d", studentID))
			if found {
				names = append(names, item.Name)
			}
		}
		return names, nil
	}

	var students []edumodel.Student
	listErr := s.db.
		Where("id IN ?", studentIDs).
		Order("id ASC").
		Find(&students).Error
	if listErr != nil {
		return nil, listErr
	}

	nameMap := make(map[uint64]string, len(students))
	for _, student := range students {
		nameMap[student.ID] = student.Name
	}

	names := make([]string, 0, len(studentIDs))
	for _, studentID := range studentIDs {
		if nameMap[studentID] == "" {
			continue
		}
		names = append(names, nameMap[studentID])
	}

	return names, nil
}

func (s *Service) noticeStudentIDs(rawNoticeID string) ([]uint64, error) {
	if s.db == nil {
		rawItem, found := demo.FindNotice(rawNoticeID)
		if !found {
			return []uint64{}, nil
		}

		return noticeStudentIDsFromDemo(rawItem.StudentIDs), nil
	}

	query := `
SELECT
  COALESCE(student_id, 0) AS student_id
FROM notice_targets
WHERE notice_id = ?
  AND target_type = 'student'
  AND student_id IS NOT NULL
ORDER BY id ASC
`

	type noticeStudentRow struct {
		StudentID uint64 `gorm:"column:student_id"`
	}

	var rows []noticeStudentRow
	listErr := s.db.Raw(query, rawNoticeID).Scan(&rows).Error
	if listErr != nil {
		if strings.Contains(strings.ToLower(listErr.Error()), "doesn't exist") {
			return []uint64{}, nil
		}
		return nil, listErr
	}

	studentIDs := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row.StudentID == 0 {
			continue
		}
		studentIDs = append(studentIDs, row.StudentID)
	}

	return studentIDs, nil
}

func (s *Service) replaceNoticeTargets(noticeID uint64, input NoticePayload) error {
	if s.db == nil {
		return nil
	}

	return s.replaceNoticeTargetsTx(s.db, noticeID, input)
}

func (s *Service) replaceNoticeTargetsTx(tx *gorm.DB, noticeID uint64, input NoticePayload) error {
	if tx == nil {
		return nil
	}

	deleteErr := tx.Where("notice_id = ?", noticeID).Delete(&edumodel.NoticeTarget{}).Error
	if deleteErr != nil {
		return deleteErr
	}

	targets := make([]edumodel.NoticeTarget, 0, len(input.StudentIDs)+1)
	now := time.Now()

	if input.RelatedClassID > 0 {
		classID := input.RelatedClassID
		targets = append(targets, edumodel.NoticeTarget{
			NoticeID:   noticeID,
			TargetType: "class",
			ClassID:    &classID,
			CreatedAt:  now,
		})
	}

	for _, studentID := range input.StudentIDs {
		currentStudentID := studentID
		targets = append(targets, edumodel.NoticeTarget{
			NoticeID:   noticeID,
			TargetType: "student",
			StudentID:  &currentStudentID,
			CreatedAt:  now,
		})
	}

	if len(targets) == 0 {
		return nil
	}

	return tx.Create(&targets).Error
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
			{ID: 1, UserID: uint64Pointer(4), Name: "周老师", Mobile: "13800000001", MainSubject: "数学思维", EmploymentType: "全职", WeeklyHours: 18, Campus: "明发校区", Status: "在职", CreatedAt: now, UpdatedAt: now},
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

		noticeTargets := []edumodel.NoticeTarget{
			{ID: 1, NoticeID: 3, TargetType: "class", ClassID: &relatedClassID, CreatedAt: now},
		}

		noticeTargetErr := tx.Create(&noticeTargets).Error
		if noticeTargetErr != nil {
			return noticeTargetErr
		}

		return nil
	})
}

func teacherItemsFromDemo() []TeacherItem {
	items := make([]TeacherItem, 0, len(demo.Teachers()))
	for _, item := range demo.Teachers() {
		items = append(items, TeacherItem{
			ID:             uint64(item.ID),
			UserID:         demoTeacherUserID(item.ID),
			UserName:       demoTeacherUserName(demoTeacherUserID(item.ID)),
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

func (s *Service) teacherItemsFromDemoWithFilter(filter TeacherFilter) []TeacherItem {
	items := teacherItemsFromDemo()
	filteredItems := make([]TeacherItem, 0, len(items))
	keyword := strings.TrimSpace(filter.Keyword)
	status := strings.TrimSpace(filter.Status)
	employmentType := strings.TrimSpace(filter.EmploymentType)
	campus := strings.TrimSpace(filter.Campus)

	for _, item := range items {
		if keyword != "" && !matchesTeacherKeyword(item, keyword) {
			continue
		}
		if status != "" && item.Status != status {
			continue
		}
		if employmentType != "" && item.EmploymentType != employmentType {
			continue
		}
		if campus != "" && item.Campus != campus {
			continue
		}

		filteredItems = append(filteredItems, item)
	}

	return filteredItems
}

func teacherItemFromDemo(rawID string) (TeacherItem, bool, error) {
	item, found := demo.FindTeacher(rawID)
	if !found {
		return TeacherItem{}, false, nil
	}
	return TeacherItem{
		ID:             uint64(item.ID),
		UserID:         demoTeacherUserID(item.ID),
		UserName:       demoTeacherUserName(demoTeacherUserID(item.ID)),
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

func (s *Service) ensureTeacherUserBinding(userID uint64, excludeTeacherID uint64) error {
	if s.db == nil || userID == 0 {
		return nil
	}

	var user edumodel.User
	findUserErr := s.db.Where("id = ?", userID).Take(&user).Error
	if errors.Is(findUserErr, gorm.ErrRecordNotFound) {
		return ErrTeacherUserNotFound
	}
	if findUserErr != nil {
		return findUserErr
	}

	roleCodes, roleErr := s.userRoleCodes(userID)
	if roleErr != nil {
		return roleErr
	}
	if !containsString(roleCodes, "teacher") {
		return ErrTeacherUserRoleInvalid
	}

	alreadyBound, boundErr := s.teacherUserBound(userID, excludeTeacherID)
	if boundErr != nil {
		return boundErr
	}
	if alreadyBound {
		return ErrTeacherUserAlreadyBound
	}

	return nil
}

func (s *Service) userRoleCodes(userID uint64) ([]string, error) {
	var rows []userRoleRecord
	listErr := s.db.Table("user_roles AS ur").
		Select("ur.user_id, ur.role_id, r.code AS role_code, r.name AS role_name").
		Joins("JOIN roles AS r ON r.id = ur.role_id").
		Where("ur.user_id = ?", userID).
		Scan(&rows).Error
	if listErr != nil {
		return nil, listErr
	}

	roleCodes := make([]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.RoleCode) == "" {
			continue
		}
		roleCodes = append(roleCodes, row.RoleCode)
	}

	return uniqueSortedStrings(roleCodes), nil
}

func (s *Service) teacherUserBound(userID uint64, excludeTeacherID uint64) (bool, error) {
	query := s.db.Model(&edumodel.Teacher{}).Where("user_id = ?", userID)
	if excludeTeacherID > 0 {
		query = query.Where("id <> ?", excludeTeacherID)
	}

	var count int64
	countErr := query.Count(&count).Error
	if countErr != nil {
		return false, countErr
	}

	return count > 0, nil
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

func (s *Service) studentItemsFromDemoWithFilter(filter StudentFilter) []StudentItem {
	items := studentItemsFromDemo()
	filteredItems := make([]StudentItem, 0, len(items))
	keyword := strings.TrimSpace(filter.Keyword)
	status := strings.TrimSpace(filter.Status)

	for _, item := range items {
		if keyword != "" && !matchesStudentKeyword(item, keyword) {
			continue
		}
		if status != "" && item.Status != status {
			continue
		}
		if filter.ClassID > 0 && item.ClassID != filter.ClassID {
			continue
		}

		filteredItems = append(filteredItems, item)
	}

	return filteredItems
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

func (s *Service) classItemsFromDemoWithFilter(filter ClassFilter) []ClassItem {
	items := classItemsFromDemo(demo.Classes())
	filteredItems := make([]ClassItem, 0, len(items))
	keyword := strings.TrimSpace(filter.Keyword)
	status := strings.TrimSpace(filter.Status)

	for _, item := range items {
		if keyword != "" && !matchesClassKeyword(item, keyword) {
			continue
		}
		if status != "" && item.Status != status {
			continue
		}
		if filter.CourseID > 0 && item.CourseID != filter.CourseID {
			continue
		}
		if filter.TeacherID > 0 && item.TeacherID != filter.TeacherID {
			continue
		}

		filteredItems = append(filteredItems, item)
	}

	return filteredItems
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

func scheduleItemsFromDemoWithFilter(filter ScheduleFilter) []ScheduleItem {
	items := scheduleItemsFromDemo(demo.Schedules())
	filteredItems := make([]ScheduleItem, 0, len(items))
	filterStatus := strings.TrimSpace(filter.Status)

	for _, item := range items {
		if filter.ClassID > 0 && item.ClassID != filter.ClassID {
			continue
		}
		if filter.TeacherID > 0 && item.TeacherID != filter.TeacherID {
			continue
		}
		if filter.Scope.RestrictToSelf && filter.Scope.TeacherID > 0 && item.TeacherID != filter.Scope.TeacherID {
			continue
		}
		if filterStatus != "" && item.AttendanceStatus != filterStatus {
			continue
		}
		if !matchesDateRange(item.LessonDate, filter.DateFrom, filter.DateTo) {
			continue
		}

		filteredItems = append(filteredItems, item)
	}

	return filteredItems
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
			StudentIDs:        noticeStudentIDsFromDemo(item.StudentIDs),
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
	filterNoticeType := strings.TrimSpace(filter.NoticeType)

	for _, item := range items {
		if filter.ClassID > 0 && item.RelatedClassID != filter.ClassID {
			continue
		}
		if filterStatus != "" && item.Status != filterStatus {
			continue
		}
		if filterNoticeType != "" && item.Category != filterNoticeType {
			continue
		}

		noticeDate := noticeDateValue(item.PublishAt)
		if filterDate != "" && noticeDate != filterDate {
			continue
		}
		if !matchesDateRange(noticeDate, filter.DateFrom, filter.DateTo) {
			continue
		}
		scopeMatched, scopeErr := s.noticeMatchesScope(filter.Scope, item)
		if scopeErr != nil || !scopeMatched {
			continue
		}

		filteredItems = append(filteredItems, item)
	}

	return filteredItems
}

func appendNoticeScopeFilter(query string, queryArgs []any, scope Scope) (string, []any) {
	if !scope.RestrictToSelf || scope.TeacherID == 0 {
		return query, queryArgs
	}

	query += ` AND (
COALESCE(related_schedule_id, 0) IN (
  SELECT id FROM class_schedules WHERE teacher_id = ?
)
OR COALESCE(related_class_id, 0) IN (
  SELECT id FROM classes WHERE teacher_id = ?
)
OR id IN (
  SELECT nt.notice_id
  FROM notice_targets AS nt
  JOIN classes AS c
    ON c.id = nt.class_id
  WHERE nt.target_type = 'class'
    AND c.teacher_id = ?
)
OR id IN (
  SELECT nt.notice_id
  FROM notice_targets AS nt
  JOIN class_students AS cs
    ON cs.student_id = nt.student_id
   AND cs.status = ?
  JOIN classes AS c
    ON c.id = cs.class_id
  WHERE nt.target_type = 'student'
    AND c.teacher_id = ?
)
)`

	queryArgs = append(
		queryArgs,
		scope.TeacherID,
		scope.TeacherID,
		scope.TeacherID,
		activeClassStudentStatus,
		scope.TeacherID,
	)

	return query, queryArgs
}

func (s *Service) noticeMatchesScope(scope Scope, noticeItem NoticeItem) (bool, error) {
	if !scope.RestrictToSelf {
		return true, nil
	}
	if scope.TeacherID == 0 {
		return false, nil
	}

	if noticeItem.RelatedScheduleID > 0 {
		scheduleItem, scheduleFound, scheduleErr := s.Schedule(fmt.Sprintf("%d", noticeItem.RelatedScheduleID))
		if scheduleErr != nil {
			return false, scheduleErr
		}
		if scheduleFound && scheduleItem.TeacherID == scope.TeacherID {
			return true, nil
		}
	}

	if noticeItem.RelatedClassID > 0 {
		classItem, classFound, classErr := s.Class(fmt.Sprintf("%d", noticeItem.RelatedClassID))
		if classErr != nil {
			return false, classErr
		}
		if classFound && classItem.TeacherID == scope.TeacherID {
			return true, nil
		}
	}

	return s.anyStudentBelongsToTeacher(noticeItem.StudentIDs, scope.TeacherID)
}

func (s *Service) NoticePayloadAccessible(input NoticePayload, scope Scope) (bool, error) {
	if !scope.RestrictToSelf {
		return true, nil
	}
	if scope.TeacherID == 0 {
		return false, nil
	}

	hasScopedTarget := false
	if input.RelatedScheduleID > 0 {
		scheduleItem, found, scheduleErr := s.Schedule(fmt.Sprintf("%d", input.RelatedScheduleID))
		if scheduleErr != nil {
			return false, scheduleErr
		}
		if !found || scheduleItem.TeacherID != scope.TeacherID {
			return false, nil
		}
		hasScopedTarget = true
	}

	if input.RelatedClassID > 0 {
		classItem, found, classErr := s.Class(fmt.Sprintf("%d", input.RelatedClassID))
		if classErr != nil {
			return false, classErr
		}
		if !found || classItem.TeacherID != scope.TeacherID {
			return false, nil
		}
		hasScopedTarget = true
	}

	if len(input.StudentIDs) > 0 {
		allMatched, matchErr := s.allStudentsBelongToTeacher(input.StudentIDs, scope.TeacherID)
		if matchErr != nil {
			return false, matchErr
		}
		if !allMatched {
			return false, nil
		}
		hasScopedTarget = true
	}

	return hasScopedTarget, nil
}

func (s *Service) anyStudentBelongsToTeacher(studentIDs []uint64, teacherID uint64) (bool, error) {
	if teacherID == 0 || len(studentIDs) == 0 {
		return false, nil
	}

	if s.db == nil {
		for _, studentID := range studentIDs {
			classes := demo.StudentClasses(fmt.Sprintf("%d", studentID))
			for _, classItem := range classes {
				if uint64(classItem.TeacherID) == teacherID {
					return true, nil
				}
			}
		}

		return false, nil
	}

	query := `
SELECT cs.student_id
FROM class_students AS cs
JOIN classes AS c
  ON c.id = cs.class_id
WHERE cs.student_id IN ?
  AND cs.status = ?
  AND c.teacher_id = ?
LIMIT 1
`

	type studentScopeRow struct {
		StudentID uint64 `gorm:"column:student_id"`
	}

	var row studentScopeRow
	findErr := s.db.Raw(query, studentIDs, activeClassStudentStatus, teacherID).Scan(&row).Error
	if findErr != nil {
		return false, findErr
	}

	return row.StudentID > 0, nil
}

func (s *Service) allStudentsBelongToTeacher(studentIDs []uint64, teacherID uint64) (bool, error) {
	uniqueStudentIDs := uniqueUint64s(studentIDs)
	if teacherID == 0 || len(uniqueStudentIDs) == 0 {
		return false, nil
	}

	if s.db == nil {
		for _, studentID := range uniqueStudentIDs {
			classes := demo.StudentClasses(fmt.Sprintf("%d", studentID))
			matched := false
			for _, classItem := range classes {
				if uint64(classItem.TeacherID) != teacherID {
					continue
				}
				matched = true
				break
			}
			if !matched {
				return false, nil
			}
		}

		return true, nil
	}

	type studentMatchRow struct {
		StudentID uint64 `gorm:"column:student_id"`
	}

	var rows []studentMatchRow
	query := `
SELECT DISTINCT cs.student_id
FROM class_students AS cs
JOIN classes AS c
  ON c.id = cs.class_id
WHERE cs.student_id IN ?
  AND cs.status = ?
  AND c.teacher_id = ?
`

	listErr := s.db.Raw(query, uniqueStudentIDs, activeClassStudentStatus, teacherID).Scan(&rows).Error
	if listErr != nil {
		return false, listErr
	}

	matchedStudentIDs := make(map[uint64]struct{}, len(rows))
	for _, row := range rows {
		matchedStudentIDs[row.StudentID] = struct{}{}
	}

	for _, studentID := range uniqueStudentIDs {
		if _, exists := matchedStudentIDs[studentID]; !exists {
			return false, nil
		}
	}

	return true, nil
}

func matchesDateRange(targetDate string, dateFrom string, dateTo string) bool {
	targetDate = strings.TrimSpace(targetDate)
	dateFrom = strings.TrimSpace(dateFrom)
	dateTo = strings.TrimSpace(dateTo)

	if targetDate == "" {
		return dateFrom == "" && dateTo == ""
	}

	if dateFrom != "" && targetDate < dateFrom {
		return false
	}

	if dateTo != "" && targetDate > dateTo {
		return false
	}

	return true
}

func noticeDateValue(rawValue string) string {
	trimmedValue := strings.TrimSpace(rawValue)
	if len(trimmedValue) >= len(dateLayout) {
		return trimmedValue[:len(dateLayout)]
	}

	return trimmedValue
}

func cloneNoticeStudentIDs(source []uint64) []uint64 {
	if len(source) == 0 {
		return []uint64{}
	}

	items := make([]uint64, len(source))
	copy(items, source)
	return items
}

func noticeStudentIDsFromDemo(source []int) []uint64 {
	if len(source) == 0 {
		return []uint64{}
	}

	items := make([]uint64, 0, len(source))
	for _, item := range source {
		if item <= 0 {
			continue
		}
		items = append(items, uint64(item))
	}

	return items
}

func demoNoticeStudentIDs(source []uint64) []int {
	if len(source) == 0 {
		return []int{}
	}

	items := make([]int, 0, len(source))
	for _, item := range source {
		if item == 0 {
			continue
		}
		items = append(items, int(item))
	}

	return items
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

func attendanceItemsFromDemo(rawScheduleID string, source []demo.AttendanceItem) []AttendanceItem {
	items := make([]AttendanceItem, 0, len(source))
	for _, item := range source {
		recordID := attendanceRecordIDFromParts(rawScheduleID, uint64(item.StudentID))
		items = append(items, AttendanceItem{
			RecordID:     recordID,
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

func attendanceRecordIDFromParts(rawScheduleID string, studentID uint64) uint64 {
	scheduleID, parseErr := strconv.ParseUint(strings.TrimSpace(rawScheduleID), 10, 64)
	if parseErr != nil || scheduleID == 0 || studentID == 0 {
		return 0
	}

	return attendanceRecordIDFromValues(scheduleID, studentID)
}

func attendanceRecordIDFromValues(scheduleID uint64, studentID uint64) uint64 {
	if scheduleID == 0 || studentID == 0 {
		return 0
	}

	return scheduleID*attendanceRecordSyntheticFactor + studentID
}

func attendanceRecordPartsFromID(recordID uint64) (uint64, uint64, bool) {
	if recordID < attendanceRecordSyntheticFactor {
		return 0, 0, false
	}

	scheduleID := recordID / attendanceRecordSyntheticFactor
	studentID := recordID % attendanceRecordSyntheticFactor
	if scheduleID == 0 || studentID == 0 {
		return 0, 0, false
	}

	return scheduleID, studentID, true
}

func attendanceRecordItemFromAttendanceItem(scheduleItem ScheduleItem, attendanceItem AttendanceItem) AttendanceRecordItem {
	return AttendanceRecordItem{
		ID:           attendanceItem.RecordID,
		ScheduleID:   scheduleItem.ID,
		ClassID:      scheduleItem.ClassID,
		ClassName:    scheduleItem.ClassName,
		StudentID:    attendanceItem.StudentID,
		StudentName:  attendanceItem.StudentName,
		TeacherID:    scheduleItem.TeacherID,
		TeacherName:  scheduleItem.TeacherName,
		LessonDate:   scheduleItem.LessonDate,
		LessonTime:   scheduleItem.LessonTime,
		Status:       attendanceItem.Status,
		Remark:       attendanceItem.Remark,
		UpdatedBy:    attendanceItem.UpdatedBy,
		UpdatedAt:    attendanceItem.UpdatedAt,
		ParentMobile: attendanceItem.ParentMobile,
	}
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

func (s *Service) homeworkItemsFromDemoWithFilter(filter HomeworkFilter) []HomeworkItem {
	items := homeworkItemsFromDemo(demo.Homeworks())
	if len(items) == 0 {
		return items
	}

	filteredItems := make([]HomeworkItem, 0, len(items))
	for _, item := range items {
		if filter.ClassID > 0 && item.ClassID != filter.ClassID {
			continue
		}

		scheduleItem, scheduleFound, scheduleErr := s.Schedule(fmt.Sprintf("%d", item.ScheduleID))
		if scheduleErr != nil {
			continue
		}

		if filter.TeacherID > 0 {
			if !scheduleFound || scheduleItem.TeacherID != filter.TeacherID {
				continue
			}
		}

		if !matchesDateRange(item.LessonDate, filter.DateFrom, filter.DateTo) {
			continue
		}

		filteredItems = append(filteredItems, item)
	}

	return filteredItems
}

func (s *Service) feedbackItemsFromDemoWithFilter(filter FeedbackFilter) []FeedbackItem {
	items := feedbackItemsFromDemo(demo.Feedbacks())
	if len(items) == 0 {
		return items
	}

	filteredItems := make([]FeedbackItem, 0, len(items))
	for _, item := range items {
		if filter.ClassID > 0 && item.ClassID != filter.ClassID {
			continue
		}

		scheduleItem, scheduleFound, scheduleErr := s.Schedule(fmt.Sprintf("%d", item.ScheduleID))
		if scheduleErr != nil {
			continue
		}

		if filter.TeacherID > 0 {
			if !scheduleFound || scheduleItem.TeacherID != filter.TeacherID {
				continue
			}
		}

		if !matchesDateRange(item.LessonDate, filter.DateFrom, filter.DateTo) {
			continue
		}

		filteredItems = append(filteredItems, item)
	}

	return filteredItems
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

func (s *Service) teacherQuery(filter TeacherFilter) *gorm.DB {
	teacherQuery := s.db.Table("teachers AS t").
		Select(`
t.id,
COALESCE(t.user_id, 0) AS user_id,
COALESCE(u.display_name, '') AS user_name,
t.name,
COALESCE(t.mobile, '') AS mobile,
COALESCE(t.title, '') AS title,
COALESCE(t.main_subject, '') AS main_subject,
COALESCE(t.employment_type, '') AS employment_type,
COALESCE(t.weekly_hours, 0) AS weekly_hours,
COALESCE(t.campus, '') AS campus,
COALESCE(t.status, '') AS status,
COALESCE(t.remark, '') AS remark
`).
		Joins("LEFT JOIN users AS u ON u.id = t.user_id")

	keyword := strings.TrimSpace(filter.Keyword)
	if keyword != "" {
		likeKeyword := "%" + keyword + "%"
		teacherQuery = teacherQuery.Where(
			`(t.name LIKE ? OR COALESCE(t.title, '') LIKE ? OR COALESCE(t.main_subject, '') LIKE ? OR COALESCE(t.mobile, '') LIKE ? OR COALESCE(t.campus, '') LIKE ? OR COALESCE(t.remark, '') LIKE ? OR COALESCE(u.display_name, '') LIKE ? OR COALESCE(u.username, '') LIKE ?)`,
			likeKeyword,
			likeKeyword,
			likeKeyword,
			likeKeyword,
			likeKeyword,
			likeKeyword,
			likeKeyword,
			likeKeyword,
		)
	}

	status := strings.TrimSpace(filter.Status)
	if status != "" {
		teacherQuery = teacherQuery.Where("t.status = ?", status)
	}

	employmentType := strings.TrimSpace(filter.EmploymentType)
	if employmentType != "" {
		teacherQuery = teacherQuery.Where("t.employment_type = ?", employmentType)
	}

	campus := strings.TrimSpace(filter.Campus)
	if campus != "" {
		teacherQuery = teacherQuery.Where("t.campus = ?", campus)
	}

	return teacherQuery
}

func (s *Service) studentQuery(filter StudentFilter) *gorm.DB {
	studentQuery := s.db.Table("students AS s").
		Select(`
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
`).
		Joins("LEFT JOIN student_guardians AS g ON g.student_id = s.id AND g.is_primary = 1").
		Joins(`
LEFT JOIN (
  SELECT student_id, MIN(class_id) AS class_id
  FROM class_students
  WHERE status = ?
  GROUP BY student_id
) AS cs ON cs.student_id = s.id
`, activeClassStudentStatus).
		Joins("LEFT JOIN classes AS c ON c.id = cs.class_id")

	keyword := strings.TrimSpace(filter.Keyword)
	if keyword != "" {
		likeKeyword := "%" + keyword + "%"
		studentQuery = studentQuery.Where(
			"(s.name LIKE ? OR s.grade_name LIKE ? OR COALESCE(g.name, '') LIKE ? OR COALESCE(g.mobile, '') LIKE ?)",
			likeKeyword,
			likeKeyword,
			likeKeyword,
			likeKeyword,
		)
	}

	status := strings.TrimSpace(filter.Status)
	if status != "" {
		studentQuery = studentQuery.Where("s.status = ?", status)
	}

	if filter.ClassID > 0 {
		studentQuery = studentQuery.Where(
			`EXISTS (
  SELECT 1
  FROM class_students AS filter_cs
  WHERE filter_cs.student_id = s.id
    AND filter_cs.class_id = ?
    AND filter_cs.status = ?
)`,
			filter.ClassID,
			activeClassStudentStatus,
		)
	}

	return studentQuery
}

func (s *Service) classQuery(filter ClassFilter) *gorm.DB {
	classQuery := s.db.Table("classes AS c").
		Select(`
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
c.status,
COALESCE(c.remark, '') AS remark
`).
		Joins("LEFT JOIN courses AS co ON co.id = c.course_id").
		Joins("LEFT JOIN teachers AS t ON t.id = c.teacher_id").
		Joins("LEFT JOIN class_students AS cs ON cs.class_id = c.id AND cs.status = ?", activeClassStudentStatus)

	keyword := strings.TrimSpace(filter.Keyword)
	if keyword != "" {
		likeKeyword := "%" + keyword + "%"
		classQuery = classQuery.Where(
			"(c.name LIKE ? OR COALESCE(co.name, '') LIKE ? OR COALESCE(t.name, '') LIKE ? OR COALESCE(c.campus, '') LIKE ?)",
			likeKeyword,
			likeKeyword,
			likeKeyword,
			likeKeyword,
		)
	}

	status := strings.TrimSpace(filter.Status)
	if status != "" {
		classQuery = classQuery.Where("c.status = ?", status)
	}

	if filter.CourseID > 0 {
		classQuery = classQuery.Where("c.course_id = ?", filter.CourseID)
	}

	if filter.TeacherID > 0 {
		classQuery = classQuery.Where("c.teacher_id = ?", filter.TeacherID)
	}

	if filter.Scope.RestrictToSelf && filter.Scope.TeacherID > 0 {
		classQuery = classQuery.Where("c.teacher_id = ?", filter.Scope.TeacherID)
	}

	return classQuery.Group(`
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
c.status,
c.remark
`)
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

func matchesStudentKeyword(item StudentItem, keyword string) bool {
	return strings.Contains(item.Name, keyword) ||
		strings.Contains(item.Grade, keyword) ||
		strings.Contains(item.ParentName, keyword) ||
		strings.Contains(item.ParentMobile, keyword)
}

func matchesClassKeyword(item ClassItem, keyword string) bool {
	return strings.Contains(item.Name, keyword) ||
		strings.Contains(item.CourseName, keyword) ||
		strings.Contains(item.TeacherName, keyword) ||
		strings.Contains(item.Campus, keyword)
}

func matchesTeacherKeyword(item TeacherItem, keyword string) bool {
	return strings.Contains(item.Name, keyword) ||
		strings.Contains(item.Title, keyword) ||
		strings.Contains(item.MainSubject, keyword) ||
		strings.Contains(item.Mobile, keyword) ||
		strings.Contains(item.Campus, keyword) ||
		strings.Contains(item.Remark, keyword) ||
		strings.Contains(item.UserName, keyword)
}

func optionalUint64Pointer(value uint64) *uint64 {
	if value == 0 {
		return nil
	}

	pointerValue := value
	return &pointerValue
}

func uint64Pointer(value uint64) *uint64 {
	pointerValue := value
	return &pointerValue
}

func demoTeacherUserID(teacherID int) uint64 {
	switch teacherID {
	case 1:
		return 4
	default:
		return 0
	}
}

func demoTeacherUserName(userID uint64) string {
	switch userID {
	case 4:
		return "周老师"
	default:
		return ""
	}
}

func containsString(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}

	return false
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

func hasOperator(operator Operator) bool {
	return operator.UserID > 0 || strings.TrimSpace(operator.DisplayName) != ""
}

func operationUserID(operator Operator) uint64 {
	if operator.UserID > 0 {
		return operator.UserID
	}

	return 1
}

func operationDisplayName(operator Operator, fallback string) string {
	displayName := strings.TrimSpace(operator.DisplayName)
	if displayName != "" {
		return displayName
	}

	displayName = strings.TrimSpace(fallback)
	if displayName != "" {
		return displayName
	}

	return "系统管理员"
}

func (s *Service) recordBusinessOperationTx(
	tx *gorm.DB,
	operator Operator,
	module string,
	action string,
	targetType string,
	targetID uint64,
	content string,
) error {
	if !hasOperator(operator) {
		return nil
	}

	return s.recordOperationTx(tx, operator, module, action, targetType, targetID, content)
}
