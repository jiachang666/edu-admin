package service

import (
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
