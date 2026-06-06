package service

import (
	"strings"
	"testing"

	edumodel "edu-admin/internal/modules/edu/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestOperationLogQueryAppliesActionFilter(t *testing.T) {
	db, openErr := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:123456@tcp(127.0.0.1:3307)/edu-admin?charset=utf8mb4&parseTime=True&loc=Local",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun: true,
	})
	if openErr != nil {
		t.Fatalf("gorm.Open returned error: %v", openErr)
	}

	svc := New(nil)
	var logs []edumodel.OperationLog
	statement := svc.operationLogQuery(db, OperationLogFilter{
		UserID:   9,
		Module:   "attendance",
		Action:   "update_record",
		DateFrom: "2026-06-01",
		DateTo:   "2026-06-30",
	}).Find(&logs).Statement

	sqlText := statement.SQL.String()
	for _, expected := range []string{
		"user_id = ?",
		"module = ?",
		"action = ?",
		"DATE(created_at) >= ?",
		"DATE(created_at) <= ?",
	} {
		if !strings.Contains(sqlText, expected) {
			t.Fatalf("expected SQL to contain %q, got %s", expected, sqlText)
		}
	}

	if len(statement.Vars) != 5 {
		t.Fatalf("expected 5 SQL vars, got %d", len(statement.Vars))
	}
	if statement.Vars[0] != uint64(9) {
		t.Fatalf("expected first SQL var to be user id, got %#v", statement.Vars[0])
	}
	if statement.Vars[1] != "attendance" {
		t.Fatalf("expected second SQL var to be module, got %#v", statement.Vars[1])
	}
	if statement.Vars[2] != "update_record" {
		t.Fatalf("expected third SQL var to be action, got %#v", statement.Vars[2])
	}
	if statement.Vars[3] != "2026-06-01" {
		t.Fatalf("expected fourth SQL var to be dateFrom, got %#v", statement.Vars[3])
	}
	if statement.Vars[4] != "2026-06-30" {
		t.Fatalf("expected fifth SQL var to be dateTo, got %#v", statement.Vars[4])
	}
}
