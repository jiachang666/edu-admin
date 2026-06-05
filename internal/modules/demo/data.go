package demo

import (
	"strconv"
	"sync"
	"time"
)

const (
	dateLayout     = "2006-01-02"
	dateTimeLayout = "2006-01-02 15:04"
)

type Teacher struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Mobile         string `json:"mobile"`
	MainSubject    string `json:"mainSubject"`
	EmploymentType string `json:"employmentType"`
	WeeklyHours    int    `json:"weeklyHours"`
	Campus         string `json:"campus"`
	Status         string `json:"status"`
}

type Student struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Grade          string `json:"grade"`
	ParentName     string `json:"parentName"`
	ParentMobile   string `json:"parentMobile"`
	Campus         string `json:"campus"`
	ClassID        int    `json:"classId"`
	ClassName      string `json:"className"`
	RemainingHours int    `json:"remainingHours"`
	Status         string `json:"status"`
}

type Course struct {
	ID                    int    `json:"id"`
	Name                  string `json:"name"`
	Category              string `json:"category"`
	Description           string `json:"description"`
	AgeRange              string `json:"ageRange"`
	LessonDurationMinutes int    `json:"lessonDurationMinutes"`
	TotalLessons          int    `json:"totalLessons"`
	DeliveryType          string `json:"deliveryType"`
	Status                string `json:"status"`
	ClassCount            int    `json:"classCount"`
}

type Class struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	CourseName     string `json:"courseName"`
	TeacherID      int    `json:"teacherId"`
	TeacherName    string `json:"teacherName"`
	Campus         string `json:"campus"`
	StudentCount   int    `json:"studentCount"`
	Capacity       int    `json:"capacity"`
	WeeklySchedule string `json:"weeklySchedule"`
	Status         string `json:"status"`
}

type Schedule struct {
	ID               int    `json:"id"`
	ClassID          int    `json:"classId"`
	ClassName        string `json:"className"`
	CourseName       string `json:"courseName"`
	TeacherName      string `json:"teacherName"`
	Campus           string `json:"campus"`
	Classroom        string `json:"classroom"`
	LessonDate       string `json:"lessonDate"`
	LessonTime       string `json:"lessonTime"`
	AttendanceStatus string `json:"attendanceStatus"`
}

type Notice struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	TargetScope string `json:"targetScope"`
	Status      string `json:"status"`
	PublishAt   string `json:"publishAt"`
	Author      string `json:"author"`
}

type Option struct {
	Value int    `json:"value"`
	Label string `json:"label"`
}

type NoticeTarget struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Campus string `json:"campus"`
}

type AttendanceItem struct {
	StudentID    int    `json:"studentId"`
	StudentName  string `json:"studentName"`
	Grade        string `json:"grade"`
	ParentMobile string `json:"parentMobile"`
	Status       string `json:"status"`
	Remark       string `json:"remark"`
}

var (
	attendanceStoreMu           sync.RWMutex
	attendanceOverrideStore     = map[int][]AttendanceItem{}
	scheduleStatusOverrideStore = map[int]string{}
)

func Teachers() []Teacher {
	return []Teacher{
		{ID: 1, Name: "周老师", Mobile: "13800000001", MainSubject: "数学思维", EmploymentType: "全职", WeeklyHours: 18, Campus: "明发校区", Status: "在职"},
		{ID: 2, Name: "林老师", Mobile: "13800000002", MainSubject: "英语阅读", EmploymentType: "全职", WeeklyHours: 16, Campus: "百汇校区", Status: "在职"},
		{ID: 3, Name: "陈老师", Mobile: "13800000003", MainSubject: "创意美术", EmploymentType: "兼职", WeeklyHours: 12, Campus: "明发校区", Status: "排课中"},
	}
}

func Students() []Student {
	return []Student{
		{ID: 1, Name: "李一诺", Grade: "三年级", ParentName: "李女士", ParentMobile: "13900000001", Campus: "明发校区", ClassID: 1, ClassName: "周末奥数提高班", RemainingHours: 18, Status: "在读"},
		{ID: 2, Name: "王梓涵", Grade: "四年级", ParentName: "王先生", ParentMobile: "13900000002", Campus: "明发校区", ClassID: 1, ClassName: "周末奥数提高班", RemainingHours: 12, Status: "在读"},
		{ID: 3, Name: "陈可欣", Grade: "五年级", ParentName: "陈女士", ParentMobile: "13900000003", Campus: "百汇校区", ClassID: 2, ClassName: "英语阅读进阶班", RemainingHours: 24, Status: "在读"},
		{ID: 4, Name: "张沐阳", Grade: "二年级", ParentName: "张女士", ParentMobile: "13900000004", Campus: "明发校区", ClassID: 3, ClassName: "少儿创意美术班", RemainingHours: 10, Status: "待续费"},
	}
}

func Courses() []Course {
	return []Course{
		{ID: 1, Name: "数学思维", Category: "数学", Description: "围绕数感、推理和应用题训练，适合周末持续提升。", AgeRange: "8-10岁", LessonDurationMinutes: 90, TotalLessons: 24, DeliveryType: "线下", Status: "启用", ClassCount: 1},
		{ID: 2, Name: "英语阅读", Category: "英语", Description: "通过分级阅读和表达训练，帮助孩子建立英文阅读习惯。", AgeRange: "9-11岁", LessonDurationMinutes: 90, TotalLessons: 24, DeliveryType: "线下", Status: "启用", ClassCount: 1},
		{ID: 3, Name: "创意美术", Category: "美术", Description: "以主题创作为主，兼顾色彩感受和动手表达。", AgeRange: "6-8岁", LessonDurationMinutes: 90, TotalLessons: 16, DeliveryType: "线下", Status: "启用", ClassCount: 1},
	}
}

func Classes() []Class {
	return []Class{
		{ID: 1, Name: "周末奥数提高班", CourseName: "数学思维", TeacherID: 1, TeacherName: "周老师", Campus: "明发校区", StudentCount: 2, Capacity: 16, WeeklySchedule: "周六 09:00-10:30", Status: "开班中"},
		{ID: 2, Name: "英语阅读进阶班", CourseName: "英语阅读", TeacherID: 2, TeacherName: "林老师", Campus: "百汇校区", StudentCount: 1, Capacity: 12, WeeklySchedule: "周六 14:00-15:30", Status: "开班中"},
		{ID: 3, Name: "少儿创意美术班", CourseName: "创意美术", TeacherID: 3, TeacherName: "陈老师", Campus: "明发校区", StudentCount: 1, Capacity: 10, WeeklySchedule: "周日 10:00-11:30", Status: "待满班"},
	}
}

func Schedules() []Schedule {
	now := time.Now()
	today := now.Format(dateLayout)
	tomorrow := now.AddDate(0, 0, 1).Format(dateLayout)

	schedules := []Schedule{
		{ID: 1, ClassID: 1, ClassName: "周末奥数提高班", CourseName: "数学思维", TeacherName: "周老师", Campus: "明发校区", Classroom: "A201", LessonDate: today, LessonTime: "09:00-10:30", AttendanceStatus: "待签到"},
		{ID: 2, ClassID: 2, ClassName: "英语阅读进阶班", CourseName: "英语阅读", TeacherName: "林老师", Campus: "百汇校区", Classroom: "B103", LessonDate: today, LessonTime: "14:00-15:30", AttendanceStatus: "已完成"},
		{ID: 3, ClassID: 3, ClassName: "少儿创意美术班", CourseName: "创意美术", TeacherName: "陈老师", Campus: "明发校区", Classroom: "Art-2", LessonDate: tomorrow, LessonTime: "10:00-11:30", AttendanceStatus: "待上课"},
	}

	attendanceStoreMu.RLock()
	defer attendanceStoreMu.RUnlock()

	for index := range schedules {
		if status, found := scheduleStatusOverrideStore[schedules[index].ID]; found {
			schedules[index].AttendanceStatus = status
		}
	}

	return schedules
}

func Notices() []Notice {
	now := time.Now()

	return []Notice{
		{ID: 1, Title: "端午节放假安排", Category: "校区通知", TargetScope: "全部学员家长", Status: "已发送", PublishAt: now.Add(-6 * time.Hour).Format(dateTimeLayout), Author: "运营老师"},
		{ID: 2, Title: "六月续费提醒名单确认", Category: "续费提醒", TargetScope: "待续费学员家长", Status: "草稿", PublishAt: now.Add(-2 * time.Hour).Format(dateTimeLayout), Author: "班主任"},
		{ID: 3, Title: "周末美术课材料准备说明", Category: "课程通知", TargetScope: "少儿创意美术班", Status: "待发送", PublishAt: now.Add(2 * time.Hour).Format(dateTimeLayout), Author: "教务老师"},
	}
}

func TeacherOptions() []Option {
	teachers := Teachers()
	options := make([]Option, 0, len(teachers))
	for _, teacher := range teachers {
		options = append(options, Option{Value: teacher.ID, Label: teacher.Name})
	}
	return options
}

func CourseOptions() []Option {
	courses := Courses()
	options := make([]Option, 0, len(courses))
	for _, course := range courses {
		options = append(options, Option{Value: course.ID, Label: course.Name})
	}
	return options
}

func FindTeacher(rawID string) (Teacher, bool) {
	for _, teacher := range Teachers() {
		if matchID(teacher.ID, rawID) {
			return teacher, true
		}
	}
	return Teacher{}, false
}

func FindCourse(rawID string) (Course, bool) {
	for _, course := range Courses() {
		if matchID(course.ID, rawID) {
			return course, true
		}
	}
	return Course{}, false
}

func FindStudent(rawID string) (Student, bool) {
	for _, student := range Students() {
		if matchID(student.ID, rawID) {
			return student, true
		}
	}
	return Student{}, false
}

func FindClass(rawID string) (Class, bool) {
	for _, classItem := range Classes() {
		if matchID(classItem.ID, rawID) {
			return classItem, true
		}
	}
	return Class{}, false
}

func FindSchedule(rawID string) (Schedule, bool) {
	for _, schedule := range Schedules() {
		if matchID(schedule.ID, rawID) {
			return schedule, true
		}
	}
	return Schedule{}, false
}

func FindNotice(rawID string) (Notice, bool) {
	for _, notice := range Notices() {
		if matchID(notice.ID, rawID) {
			return notice, true
		}
	}
	return Notice{}, false
}

func StudentClasses(rawStudentID string) []Class {
	student, found := FindStudent(rawStudentID)
	if !found {
		return []Class{}
	}

	classItem, classFound := FindClass(strconv.Itoa(student.ClassID))
	if !classFound {
		return []Class{}
	}

	return []Class{classItem}
}

func ClassStudents(rawClassID string) []Student {
	students := Students()
	items := make([]Student, 0, len(students))

	for _, student := range students {
		if matchID(student.ClassID, rawClassID) {
			items = append(items, student)
		}
	}

	return items
}

func UpcomingSchedules(rawClassID string) []Schedule {
	schedules := Schedules()
	items := make([]Schedule, 0, len(schedules))

	for _, schedule := range schedules {
		if matchID(schedule.ClassID, rawClassID) {
			items = append(items, schedule)
		}
	}

	return items
}

func NoticeTargets(rawNoticeID string) []NoticeTarget {
	notice, found := FindNotice(rawNoticeID)
	if !found {
		return []NoticeTarget{}
	}

	switch notice.TargetScope {
	case "全部学员家长":
		return []NoticeTarget{
			{Name: "明发校区学员家长", Type: "家长群", Campus: "明发校区"},
			{Name: "百汇校区学员家长", Type: "家长群", Campus: "百汇校区"},
		}
	case "待续费学员家长":
		return []NoticeTarget{
			{Name: "张沐阳家长", Type: "个人", Campus: "明发校区"},
		}
	default:
		return []NoticeTarget{
			{Name: "少儿创意美术班家长群", Type: "班级群", Campus: "明发校区"},
		}
	}
}

func Attendance(rawScheduleID string) []AttendanceItem {
	schedule, found := FindSchedule(rawScheduleID)
	if !found {
		return []AttendanceItem{}
	}

	scheduleID, parseErr := strconv.Atoi(rawScheduleID)
	if parseErr == nil {
		attendanceStoreMu.RLock()
		overrideItems, hasOverride := attendanceOverrideStore[scheduleID]
		attendanceStoreMu.RUnlock()
		if hasOverride {
			return cloneAttendanceItems(overrideItems)
		}
	}

	students := ClassStudents(strconv.Itoa(schedule.ClassID))
	items := make([]AttendanceItem, 0, len(students))

	for index, student := range students {
		status := "待确认"
		remark := ""

		if schedule.ID == 1 && index > 0 {
			status = "已到"
		}

		if schedule.ID == 2 {
			status = "请假"
			remark = "家长上午已请假"
		}

		if schedule.AttendanceStatus == "已完成" && schedule.ID != 2 {
			status = "已到"
		}

		if schedule.AttendanceStatus == "已完成" && schedule.ID != 2 && index == len(students)-1 {
			remark = "课堂表现积极"
		}

		items = append(items, AttendanceItem{
			StudentID:    student.ID,
			StudentName:  student.Name,
			Grade:        student.Grade,
			ParentMobile: student.ParentMobile,
			Status:       status,
			Remark:       remark,
		})
	}

	return items
}

func Overview() map[string]any {
	students := Students()
	classes := Classes()
	schedules := Schedules()
	notices := Notices()
	today := time.Now().Format(dateLayout)

	todayCourses := 0
	todayPendingCheck := 0
	todayLeaveCount := 0
	todayAbsentCount := 0
	pendingActionCount := 0

	for _, schedule := range schedules {
		if schedule.LessonDate == today {
			todayCourses++
		}
		if schedule.AttendanceStatus == "待签到" {
			todayPendingCheck++
		}

		for _, item := range Attendance(strconv.Itoa(schedule.ID)) {
			switch item.Status {
			case "请假":
				todayLeaveCount++
			case "缺席":
				todayAbsentCount++
			}
		}
	}

	for _, notice := range notices {
		if notice.Status == "草稿" || notice.Status == "待发送" {
			pendingActionCount++
		}
	}

	return map[string]any{
		"todayCourses":       todayCourses,
		"todayPendingCheck":  todayPendingCheck,
		"todayLeaveCount":    todayLeaveCount,
		"todayAbsentCount":   todayAbsentCount,
		"studentCount":       len(students),
		"classCount":         len(classes),
		"pendingActionCount": pendingActionCount,
		"upcomingLessons":    Schedules(),
		"latestNotices":      notices,
	}
}

func matchID(id int, rawID string) bool {
	return strconv.Itoa(id) == rawID
}

func SaveAttendance(rawScheduleID string, items []AttendanceItem) bool {
	scheduleID, parseErr := strconv.Atoi(rawScheduleID)
	if parseErr != nil {
		return false
	}

	attendanceStoreMu.Lock()
	defer attendanceStoreMu.Unlock()

	attendanceOverrideStore[scheduleID] = cloneAttendanceItems(items)
	scheduleStatusOverrideStore[scheduleID] = deriveAttendanceSessionStatus(items)

	return true
}

func cloneAttendanceItems(source []AttendanceItem) []AttendanceItem {
	items := make([]AttendanceItem, 0, len(source))
	for _, item := range source {
		items = append(items, item)
	}
	return items
}

func deriveAttendanceSessionStatus(items []AttendanceItem) string {
	for _, item := range items {
		if item.Status == "待确认" {
			return "待签到"
		}
	}

	return "已完成"
}
