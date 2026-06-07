package service

import (
	"fmt"
	"testing"
	"time"

	"edu-admin/internal/modules/demo"
)

func TestClassDetailIncludesRecentFeedbacks(t *testing.T) {
	svc := New(nil)

	detail, found, detailErr := svc.ClassDetail("1")
	if detailErr != nil {
		t.Fatalf("ClassDetail returned error: %v", detailErr)
	}
	if !found {
		t.Fatal("expected class detail to be found")
	}
	if len(detail.RecentFeedbacks) == 0 {
		t.Fatal("expected class detail to include recent feedbacks")
	}
	if len(detail.UpcomingSchedules) == 0 {
		t.Fatal("expected class detail to include upcoming schedules")
	}
	if len(detail.RecentAttendance) == 0 {
		t.Fatal("expected class detail to include recent attendance")
	}
}

func TestClassDetailSkipsFutureAttendanceSessions(t *testing.T) {
	svc := New(nil)

	detail, found, detailErr := svc.ClassDetail("3")
	if detailErr != nil {
		t.Fatalf("ClassDetail returned error: %v", detailErr)
	}
	if !found {
		t.Fatal("expected class detail to be found")
	}
	if len(detail.UpcomingSchedules) == 0 {
		t.Fatal("expected class detail to keep future schedules")
	}
	if len(detail.RecentAttendance) != 0 {
		t.Fatalf("expected future-only class to have no recent attendance, got %d", len(detail.RecentAttendance))
	}
}

func TestTeacherNoticeScopeMatchesStudentTargets(t *testing.T) {
	svc := New(nil)

	demo.SaveNotice(demo.Notice{
		ID:          9101,
		Title:       "teacher-own-student",
		Content:     "own student target",
		Category:    "作业提醒",
		TargetScope: "李一诺家长",
		StudentIDs:  []int{1},
		Status:      "草稿",
		PublishAt:   "2026-06-06 14:40",
		Author:      "系统管理员",
	})
	demo.SaveNotice(demo.Notice{
		ID:          9102,
		Title:       "teacher-other-student",
		Content:     "other student target",
		Category:    "作业提醒",
		TargetScope: "陈可欣家长",
		StudentIDs:  []int{3},
		Status:      "草稿",
		PublishAt:   "2026-06-06 14:41",
		Author:      "系统管理员",
	})

	scope := Scope{UserID: 4, PrimaryRole: "teacher", TeacherID: 1, RestrictToSelf: true}
	notices, noticeErr := svc.NoticesWithFilter(NoticeFilter{Scope: scope})
	if noticeErr != nil {
		t.Fatalf("NoticesWithFilter returned error: %v", noticeErr)
	}

	hasOwnStudentNotice := false
	hasOtherStudentNotice := false
	for _, item := range notices {
		if item.ID == 9101 {
			hasOwnStudentNotice = true
		}
		if item.ID == 9102 {
			hasOtherStudentNotice = true
		}
	}

	if !hasOwnStudentNotice {
		t.Fatal("expected teacher to see notice for own student")
	}
	if hasOtherStudentNotice {
		t.Fatal("expected teacher to not see notice for another teacher's student")
	}

	ownAccessible, ownErr := svc.NoticeAccessible("9101", scope)
	if ownErr != nil {
		t.Fatalf("NoticeAccessible returned error for own notice: %v", ownErr)
	}
	if !ownAccessible {
		t.Fatal("expected own student notice to be accessible")
	}

	otherAccessible, otherErr := svc.NoticeAccessible("9102", scope)
	if otherErr != nil {
		t.Fatalf("NoticeAccessible returned error for other notice: %v", otherErr)
	}
	if otherAccessible {
		t.Fatal("expected other teacher student notice to be inaccessible")
	}
}

func TestSchedulesWithFilterAppliesClassTeacherAndStatus(t *testing.T) {
	svc := New(nil)

	items, listErr := svc.SchedulesWithFilter(ScheduleFilter{
		ClassID:   1,
		TeacherID: 1,
		Status:    "待签到",
	})
	if listErr != nil {
		t.Fatalf("SchedulesWithFilter returned error: %v", listErr)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 schedule, got %d", len(items))
	}
	if items[0].ID != 1 {
		t.Fatalf("expected schedule 1, got %d", items[0].ID)
	}
}

func TestSchedulesWithFilterAppliesDateRangeAndTeacherScope(t *testing.T) {
	svc := New(nil)

	tomorrow := time.Now().AddDate(0, 0, 1).Format(dateLayout)
	items, listErr := svc.SchedulesWithFilter(ScheduleFilter{
		DateFrom: tomorrow,
		DateTo:   tomorrow,
		Scope: Scope{
			UserID:         3,
			PrimaryRole:    "teacher",
			TeacherID:      3,
			RestrictToSelf: true,
		},
	})
	if listErr != nil {
		t.Fatalf("SchedulesWithFilter returned error: %v", listErr)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 schedule in scope, got %d", len(items))
	}
	if items[0].TeacherID != 3 {
		t.Fatalf("expected teacher 3 schedule, got teacher %d", items[0].TeacherID)
	}
	if items[0].LessonDate != tomorrow {
		t.Fatalf("expected lesson date %s, got %s", tomorrow, items[0].LessonDate)
	}
}

func TestAttendanceRecordsApplyDateRange(t *testing.T) {
	svc := New(nil)

	tomorrow := time.Now().AddDate(0, 0, 1).Format(dateLayout)
	items, listErr := svc.AttendanceRecords(AttendanceRecordFilter{
		DateFrom: tomorrow,
		DateTo:   tomorrow,
	})
	if listErr != nil {
		t.Fatalf("AttendanceRecords returned error: %v", listErr)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 attendance record in date range, got %d", len(items))
	}
	if items[0].LessonDate != tomorrow {
		t.Fatalf("expected lesson date %s, got %s", tomorrow, items[0].LessonDate)
	}
	if items[0].ID == 0 {
		t.Fatal("expected attendance record id to be populated")
	}
}

func TestUpdateAttendanceRecordUpdatesSingleDemoRecord(t *testing.T) {
	svc := New(nil)

	originalItems, listErr := svc.Attendance("1")
	if listErr != nil {
		t.Fatalf("Attendance returned error: %v", listErr)
	}
	if len(originalItems) < 2 {
		t.Fatalf("expected at least 2 attendance items, got %d", len(originalItems))
	}
	defer restoreDemoAttendance(t, "1", originalItems)

	targetItem := originalItems[0]
	otherItem := originalItems[1]
	updatedItem, found, updateErr := svc.UpdateAttendanceRecord(
		fmt.Sprintf("%d", targetItem.RecordID),
		AttendanceRecordUpdatePayload{
			Status: "请假",
			Remark: "单条更正测试",
		},
		Operator{
			UserID:      1,
			DisplayName: "系统管理员",
		},
	)
	if updateErr != nil {
		t.Fatalf("UpdateAttendanceRecord returned error: %v", updateErr)
	}
	if !found {
		t.Fatal("expected attendance record to be found")
	}
	if updatedItem.ID != targetItem.RecordID {
		t.Fatalf("expected updated record id %d, got %d", targetItem.RecordID, updatedItem.ID)
	}
	if updatedItem.Status != "请假" {
		t.Fatalf("expected updated status 请假, got %s", updatedItem.Status)
	}
	if updatedItem.Remark != "单条更正测试" {
		t.Fatalf("expected updated remark to be saved, got %s", updatedItem.Remark)
	}

	afterItems, afterErr := svc.Attendance("1")
	if afterErr != nil {
		t.Fatalf("Attendance returned error after update: %v", afterErr)
	}

	var changedItem AttendanceItem
	var untouchedItem AttendanceItem
	for _, item := range afterItems {
		if item.RecordID == targetItem.RecordID {
			changedItem = item
		}
		if item.RecordID == otherItem.RecordID {
			untouchedItem = item
		}
	}

	if changedItem.RecordID == 0 {
		t.Fatal("expected changed attendance item to still exist")
	}
	if changedItem.Status != "请假" || changedItem.Remark != "单条更正测试" {
		t.Fatalf("expected changed item to be updated, got status=%s remark=%s", changedItem.Status, changedItem.Remark)
	}
	if untouchedItem.RecordID == 0 {
		t.Fatal("expected untouched attendance item to still exist")
	}
	if untouchedItem.Status != otherItem.Status || untouchedItem.Remark != otherItem.Remark {
		t.Fatalf("expected other item unchanged, got status=%s remark=%s", untouchedItem.Status, untouchedItem.Remark)
	}

	recordItem, recordFound, recordErr := svc.AttendanceRecord(fmt.Sprintf("%d", targetItem.RecordID))
	if recordErr != nil {
		t.Fatalf("AttendanceRecord returned error: %v", recordErr)
	}
	if !recordFound {
		t.Fatal("expected updated record to be queryable")
	}
	if recordItem.Status != "请假" || recordItem.Remark != "单条更正测试" {
		t.Fatalf("expected updated record values, got status=%s remark=%s", recordItem.Status, recordItem.Remark)
	}

	sessions, sessionErr := svc.AttendanceSessions()
	if sessionErr != nil {
		t.Fatalf("AttendanceSessions returned error: %v", sessionErr)
	}

	var scheduleOne AttendanceSessionItem
	for _, item := range sessions {
		if item.ID == 1 {
			scheduleOne = item
			break
		}
	}

	if scheduleOne.ID == 0 {
		t.Fatal("expected schedule 1 attendance session to exist")
	}
	if scheduleOne.AttendanceStatus != "已完成" {
		t.Fatalf("expected session status 已完成 after single update, got %s", scheduleOne.AttendanceStatus)
	}
	if scheduleOne.PendingCount != 0 {
		t.Fatalf("expected no pending attendance items, got %d", scheduleOne.PendingCount)
	}
}

func restoreDemoAttendance(t *testing.T, rawScheduleID string, items []AttendanceItem) {
	t.Helper()

	demoItems := make([]demo.AttendanceItem, 0, len(items))
	for _, item := range items {
		demoItems = append(demoItems, demo.AttendanceItem{
			StudentID:    int(item.StudentID),
			StudentName:  item.StudentName,
			Grade:        item.Grade,
			ParentMobile: item.ParentMobile,
			Status:       item.Status,
			Remark:       item.Remark,
			UpdatedBy:    item.UpdatedBy,
			UpdatedAt:    item.UpdatedAt,
		})
	}

	if !demo.SaveAttendance(rawScheduleID, demoItems) {
		t.Fatalf("failed to restore demo attendance for schedule %s", rawScheduleID)
	}
}
