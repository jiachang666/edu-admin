-- MySQL 8.0 draft schema for edu admin v0.1
-- Focus: first-release core workflow only

CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `display_name` varchar(64) NOT NULL,
  `mobile` varchar(32) DEFAULT NULL,
  `status` varchar(32) NOT NULL DEFAULT 'active',
  `last_login_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_username` (`username`),
  KEY `idx_users_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='System login accounts';

CREATE TABLE IF NOT EXISTS `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `code` varchar(64) NOT NULL,
  `status` varchar(32) NOT NULL DEFAULT 'active',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_roles_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Role definitions';

CREATE TABLE IF NOT EXISTS `user_roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `role_id` bigint unsigned NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_roles_user_role` (`user_id`, `role_id`),
  KEY `idx_user_roles_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='User role mapping';

CREATE TABLE IF NOT EXISTS `teachers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned DEFAULT NULL,
  `name` varchar(64) NOT NULL,
  `mobile` varchar(32) DEFAULT NULL,
  `title` varchar(64) DEFAULT NULL,
  `status` varchar(32) NOT NULL DEFAULT 'active',
  `remark` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_teachers_user_id` (`user_id`),
  KEY `idx_teachers_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Teacher profiles';

CREATE TABLE IF NOT EXISTS `students` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `gender` varchar(16) DEFAULT NULL,
  `birth_date` date DEFAULT NULL,
  `school_name` varchar(128) DEFAULT NULL,
  `grade_name` varchar(64) DEFAULT NULL,
  `status` varchar(32) NOT NULL DEFAULT 'active',
  `remark` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_students_name` (`name`),
  KEY `idx_students_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Student profiles';

CREATE TABLE IF NOT EXISTS `student_guardians` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `student_id` bigint unsigned NOT NULL,
  `name` varchar(64) NOT NULL,
  `relation` varchar(32) DEFAULT NULL,
  `mobile` varchar(32) NOT NULL,
  `is_primary` tinyint(1) NOT NULL DEFAULT 0,
  `remark` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_student_guardians_student_id` (`student_id`),
  KEY `idx_student_guardians_mobile` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Student guardians and contacts';

CREATE TABLE IF NOT EXISTS `courses` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,
  `category` varchar(64) DEFAULT NULL,
  `intro` text,
  `age_range` varchar(64) DEFAULT NULL,
  `lesson_duration_minutes` int unsigned DEFAULT NULL,
  `total_lessons` int unsigned DEFAULT NULL,
  `delivery_type` varchar(32) DEFAULT NULL,
  `status` varchar(32) NOT NULL DEFAULT 'active',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_courses_name` (`name`),
  KEY `idx_courses_category` (`category`),
  KEY `idx_courses_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Course templates';

CREATE TABLE IF NOT EXISTS `classes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,
  `course_id` bigint unsigned NOT NULL,
  `teacher_id` bigint unsigned NOT NULL,
  `assistant_teacher_name` varchar(64) DEFAULT NULL,
  `start_date` date DEFAULT NULL,
  `end_date` date DEFAULT NULL,
  `status` varchar(32) NOT NULL DEFAULT 'preparing',
  `remark` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_classes_course_id` (`course_id`),
  KEY `idx_classes_teacher_id` (`teacher_id`),
  KEY `idx_classes_status` (`status`),
  KEY `idx_classes_start_date` (`start_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Active teaching classes';

CREATE TABLE IF NOT EXISTS `class_students` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `class_id` bigint unsigned NOT NULL,
  `student_id` bigint unsigned NOT NULL,
  `join_date` date DEFAULT NULL,
  `leave_date` date DEFAULT NULL,
  `status` varchar(32) NOT NULL DEFAULT 'active',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_class_students_class_id` (`class_id`),
  KEY `idx_class_students_student_id` (`student_id`),
  KEY `idx_class_students_status` (`status`),
  KEY `idx_class_students_class_student` (`class_id`, `student_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Student-class membership';

CREATE TABLE IF NOT EXISTS `class_schedules` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `class_id` bigint unsigned NOT NULL,
  `course_id` bigint unsigned NOT NULL,
  `teacher_id` bigint unsigned NOT NULL,
  `schedule_type` varchar(32) NOT NULL DEFAULT 'normal',
  `schedule_date` date NOT NULL,
  `start_time` time NOT NULL,
  `end_time` time NOT NULL,
  `location` varchar(128) DEFAULT NULL,
  `status` varchar(32) NOT NULL DEFAULT 'scheduled',
  `origin_schedule_id` bigint unsigned DEFAULT NULL,
  `remark` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_class_schedules_class_id` (`class_id`),
  KEY `idx_class_schedules_teacher_id` (`teacher_id`),
  KEY `idx_class_schedules_schedule_date` (`schedule_date`),
  KEY `idx_class_schedules_status` (`status`),
  KEY `idx_class_schedules_origin_schedule_id` (`origin_schedule_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Concrete class sessions';

CREATE TABLE IF NOT EXISTS `attendance_records` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `schedule_id` bigint unsigned NOT NULL,
  `class_id` bigint unsigned NOT NULL,
  `student_id` bigint unsigned NOT NULL,
  `status` varchar(32) NOT NULL DEFAULT 'present',
  `checkin_by_user_id` bigint unsigned NOT NULL,
  `checked_at` datetime DEFAULT NULL,
  `remark` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_attendance_schedule_student` (`schedule_id`, `student_id`),
  KEY `idx_attendance_class_id` (`class_id`),
  KEY `idx_attendance_student_id` (`student_id`),
  KEY `idx_attendance_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Attendance per student per session';

CREATE TABLE IF NOT EXISTS `homeworks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `schedule_id` bigint unsigned NOT NULL,
  `class_id` bigint unsigned NOT NULL,
  `title` varchar(128) NOT NULL,
  `content` text NOT NULL,
  `submission_note` text,
  `created_by_user_id` bigint unsigned NOT NULL,
  `status` varchar(32) NOT NULL DEFAULT 'draft',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_homeworks_schedule_id` (`schedule_id`),
  KEY `idx_homeworks_class_id` (`class_id`),
  KEY `idx_homeworks_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Homework for a class session';

CREATE TABLE IF NOT EXISTS `class_feedbacks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `schedule_id` bigint unsigned NOT NULL,
  `class_id` bigint unsigned NOT NULL,
  `summary` text,
  `learning_status` text,
  `next_suggestion` text,
  `parent_notice` text,
  `created_by_user_id` bigint unsigned NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_class_feedbacks_schedule_id` (`schedule_id`),
  KEY `idx_class_feedbacks_class_id` (`class_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Class-level feedback after a session';

CREATE TABLE IF NOT EXISTS `notices` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(128) NOT NULL,
  `content` text NOT NULL,
  `notice_type` varchar(32) DEFAULT 'general',
  `sender_user_id` bigint unsigned NOT NULL,
  `related_class_id` bigint unsigned DEFAULT NULL,
  `related_schedule_id` bigint unsigned DEFAULT NULL,
  `status` varchar(32) NOT NULL DEFAULT 'draft',
  `sent_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_notices_sender_user_id` (`sender_user_id`),
  KEY `idx_notices_related_class_id` (`related_class_id`),
  KEY `idx_notices_related_schedule_id` (`related_schedule_id`),
  KEY `idx_notices_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Notice master records';

CREATE TABLE IF NOT EXISTS `notice_targets` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `notice_id` bigint unsigned NOT NULL,
  `target_type` varchar(32) NOT NULL,
  `class_id` bigint unsigned DEFAULT NULL,
  `student_id` bigint unsigned DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_notice_targets_notice_id` (`notice_id`),
  KEY `idx_notice_targets_class_id` (`class_id`),
  KEY `idx_notice_targets_student_id` (`student_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Notice recipients or target scope';

CREATE TABLE IF NOT EXISTS `operation_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `module` varchar(64) NOT NULL,
  `action` varchar(64) NOT NULL,
  `target_type` varchar(64) DEFAULT NULL,
  `target_id` bigint unsigned DEFAULT NULL,
  `content` text,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_operation_logs_user_id` (`user_id`),
  KEY `idx_operation_logs_module` (`module`),
  KEY `idx_operation_logs_action` (`action`),
  KEY `idx_operation_logs_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Operation audit logs';

-- Suggested seed roles:
-- INSERT INTO roles (name, code, status) VALUES
-- ('Super Admin', 'super_admin', 'active'),
-- ('Principal', 'principal', 'active'),
-- ('Staff', 'staff', 'active'),
-- ('Teacher', 'teacher', 'active');
