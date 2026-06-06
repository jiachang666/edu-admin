import http from "./http";

type ApiEnvelope<T> = {
  code: number;
  message: string;
  data: T;
  requestId: string;
};

export type PagedResult<T> = {
  list: T[];
  total: number;
  page: number;
  pageSize: number;
};

export type Teacher = {
  id: number;
  name: string;
  mobile: string;
  title: string;
  mainSubject: string;
  employmentType: string;
  weeklyHours: number;
  campus: string;
  status: string;
  remark: string;
};

export type TeacherPayload = {
  name: string;
  mobile: string;
  title: string;
  mainSubject: string;
  employmentType: string;
  weeklyHours: number;
  campus: string;
  status: string;
  remark: string;
};

export type TeacherDetail = {
  teacher: Teacher;
  classes: SchoolClass[];
  recentSchedules: Schedule[];
};

export type Course = {
  id: number;
  name: string;
  category: string;
  description: string;
  ageRange: string;
  lessonDurationMinutes: number;
  totalLessons: number;
  deliveryType: string;
  status: string;
  classCount: number;
};

export type CoursePayload = {
  name: string;
  category: string;
  description: string;
  ageRange: string;
  lessonDurationMinutes: number;
  totalLessons: number;
  deliveryType: string;
  status: string;
};

export type Student = {
  id: number;
  name: string;
  grade: string;
  parentName: string;
  parentMobile: string;
  campus: string;
  classId: number;
  className: string;
  remainingHours: number;
  status: string;
};

export type StudentPayload = {
  name: string;
  gender: string;
  schoolName: string;
  gradeName: string;
  campus: string;
  remainingHours: number;
  status: string;
  remark: string;
  guardianName: string;
  guardianMobile: string;
  guardianRelation: string;
};

export type StudentProfile = {
  id: number;
  name: string;
  gender: string;
  schoolName: string;
  grade: string;
  parentName: string;
  parentMobile: string;
  campus: string;
  remainingHours: number;
  status: string;
  remark: string;
};

export type StudentGuardian = {
  id: number;
  name: string;
  relation: string;
  mobile: string;
  isPrimary: boolean;
};

export type StudentGuardianPayload = {
  name: string;
  relation: string;
  mobile: string;
  isPrimary: boolean;
};

export type StudentAttendanceRecord = {
  scheduleId: number;
  classId: number;
  className: string;
  courseName: string;
  teacherName: string;
  campus: string;
  classroom: string;
  lessonDate: string;
  lessonTime: string;
  status: string;
  remark: string;
};

export type StudentDetail = {
  student: StudentProfile;
  guardians: StudentGuardian[];
  classes: SchoolClass[];
  recentSchedules: Schedule[];
  recentAttendance: StudentAttendanceRecord[];
  recentHomeworks: Homework[];
  recentFeedbacks: Feedback[];
};

export type SchoolClass = {
  id: number;
  courseId: number;
  name: string;
  courseName: string;
  teacherId: number;
  teacherName: string;
  campus: string;
  studentCount: number;
  capacity: number;
  weeklySchedule: string;
  startDate: string;
  endDate: string;
  status: string;
  remark: string;
};

export type SchoolClassPayload = {
  name: string;
  courseId: number;
  teacherId: number;
  campus: string;
  capacity: number;
  weeklySchedule: string;
  startDate: string;
  endDate: string;
  status: string;
  remark: string;
};

export type ClassDetail = {
  class: SchoolClass;
  students: Student[];
  upcomingSchedules: Schedule[];
  recentAttendance: AttendanceSession[];
  recentHomeworks: Homework[];
  recentNotices: Notice[];
};

export type Schedule = {
  id: number;
  classId: number;
  className: string;
  courseName: string;
  teacherId: number;
  teacherName: string;
  campus: string;
  classroom: string;
  scheduleType: string;
  lessonDate: string;
  lessonTime: string;
  attendanceStatus: string;
  remark: string;
};

export type SchedulePayload = {
  classId: number;
  scheduleType: string;
  lessonDate: string;
  startTime: string;
  endTime: string;
  classroom: string;
  remark: string;
};

export type ScheduleActionPayload = {
  lessonDate: string;
  startTime: string;
  endTime: string;
  classroom: string;
  remark: string;
};

export type ScheduleDetail = {
  schedule: Schedule;
  class: SchoolClass;
  students: Student[];
  attendance: AttendanceSession;
  homework?: Homework;
  feedback?: Feedback;
  relatedNotices: Notice[];
};

export type Notice = {
  id: number;
  title: string;
  content: string;
  category: string;
  targetScope: string;
  relatedClassId: number;
  status: string;
  publishAt: string;
  author: string;
};

export type NoticeTarget = {
  name: string;
  type: string;
  campus: string;
};

export type NoticePayload = {
  title: string;
  content: string;
  category: string;
  targetScope: string;
  relatedClassId: number;
  status: string;
  author: string;
};

export type AttendanceEntry = {
  studentId: number;
  studentName: string;
  grade: string;
  parentMobile: string;
  status: string;
  remark: string;
  updatedBy: string;
  updatedAt: string;
};

export type AttendanceSession = {
  id: number;
  classId: number;
  className: string;
  courseName: string;
  teacherName: string;
  campus: string;
  classroom: string;
  lessonDate: string;
  lessonTime: string;
  attendanceStatus: string;
  studentCount: number;
  presentCount: number;
  leaveCount: number;
  absentCount: number;
  pendingCount: number;
};

export type AttendanceRecord = {
  scheduleId: number;
  classId: number;
  className: string;
  studentId: number;
  studentName: string;
  teacherName: string;
  lessonDate: string;
  lessonTime: string;
  status: string;
  remark: string;
  updatedBy: string;
  updatedAt: string;
  parentMobile: string;
};

export type AttendanceDetail = {
  schedule: Schedule;
  items: AttendanceEntry[];
};

export type Homework = {
  id: number;
  scheduleId: number;
  classId: number;
  className: string;
  courseName: string;
  teacherName: string;
  lessonDate: string;
  title: string;
  content: string;
  submissionNote: string;
  status: string;
};

export type Feedback = {
  id: number;
  scheduleId: number;
  classId: number;
  className: string;
  courseName: string;
  teacherName: string;
  lessonDate: string;
  summary: string;
  learningStatus: string;
  nextSuggestion: string;
  parentNotice: string;
};

export type DashboardOverview = {
  todayCourses: number;
  todayPendingCheck: number;
  todayLeaveCount: number;
  todayAbsentCount: number;
  studentCount: number;
  classCount: number;
  pendingActionCount: number;
  upcomingLessons: Schedule[];
  latestNotices: Notice[];
};

export type UserRole = {
  id: number;
  code: string;
  name: string;
};

export type AccessUser = {
  id: number;
  username: string;
  displayName: string;
  mobile: string;
  status: string;
  lastLoginAt: string;
  roles: UserRole[];
  primaryRoleCode: string;
  primaryRoleName: string;
};

export type AccessUserPayload = {
  username: string;
  password: string;
  displayName: string;
  mobile: string;
  roleCode: string;
  status: string;
};

export type AccessRole = {
  id: number;
  name: string;
  code: string;
  description: string;
  status: string;
  userCount: number;
  permissionCount: number;
  permissions: string[];
};

export type AccessRolePayload = {
  name: string;
  code: string;
  description: string;
  status: string;
};

export type OperationLog = {
  id: number;
  userId: number;
  userName: string;
  module: string;
  action: string;
  targetType: string;
  targetId: number;
  content: string;
  createdAt: string;
};

export type SelectOption = {
  value: number;
  label: string;
};

async function unwrap<T>(request: Promise<{ data: ApiEnvelope<T> }>) {
  const response = await request;
  return response.data.data;
}

export function fetchDashboardOverview() {
  return unwrap<DashboardOverview>(http.get("/dashboard/overview"));
}

export function fetchUserList() {
  return unwrap<PagedResult<AccessUser>>(http.get("/users"));
}

export function createUser(payload: AccessUserPayload) {
  return unwrap<AccessUser>(http.post("/users", payload));
}

export function updateUser(id: number, payload: AccessUserPayload) {
  return unwrap<AccessUser>(http.patch(`/users/${id}`, payload));
}

export function enableUser(id: number) {
  return unwrap<AccessUser>(http.post(`/users/${id}/enable`));
}

export function disableUser(id: number) {
  return unwrap<AccessUser>(http.post(`/users/${id}/disable`));
}

export function fetchRoleList() {
  return unwrap<PagedResult<AccessRole>>(http.get("/roles"));
}

export function fetchRoleDetail(id: number) {
  return unwrap<{
    role: AccessRole;
    permissionGroups: Array<{
      key: string;
      label: string;
      description: string;
      permissions: Array<{ code: string; label: string; description: string }>;
    }>;
  }>(http.get(`/roles/${id}`));
}

export function createRole(payload: AccessRolePayload) {
  return unwrap<AccessRole>(http.post("/roles", payload));
}

export function updateRole(id: number, payload: AccessRolePayload) {
  return unwrap<AccessRole>(http.patch(`/roles/${id}`, payload));
}

export function saveRolePermissions(id: number, payload: { permissions: string[] }) {
  return unwrap<AccessRole>(http.put(`/roles/${id}/permissions`, payload));
}

export function fetchOperationLogList() {
  return unwrap<PagedResult<OperationLog>>(http.get("/operation-logs"));
}

export function fetchTeacherList() {
  return unwrap<PagedResult<Teacher>>(http.get("/teachers"));
}

export function createTeacher(payload: TeacherPayload) {
  return unwrap<Teacher>(http.post("/teachers", payload));
}

export function updateTeacher(id: number, payload: TeacherPayload) {
  return unwrap<Teacher>(http.patch(`/teachers/${id}`, payload));
}

export function fetchTeacherDetail(teacherId: number) {
  return unwrap<TeacherDetail>(http.get(`/teachers/${teacherId}`));
}

export function fetchTeacherOptions() {
  return unwrap<SelectOption[]>(http.get("/teachers/options"));
}

export function fetchCourseList(params?: {
  keyword?: string;
  category?: string;
  status?: string;
}) {
  return unwrap<PagedResult<Course>>(http.get("/courses", { params }));
}

export function createCourse(payload: CoursePayload) {
  return unwrap<Course>(http.post("/courses", payload));
}

export function updateCourse(id: number, payload: CoursePayload) {
  return unwrap<Course>(http.patch(`/courses/${id}`, payload));
}

export function fetchCourseOptions() {
  return unwrap<SelectOption[]>(http.get("/courses/options"));
}

export function fetchStudentList() {
  return unwrap<PagedResult<Student>>(http.get("/students"));
}

export function createStudent(payload: StudentPayload) {
  return unwrap<StudentDetail>(http.post("/students", payload));
}

export function updateStudent(id: number, payload: StudentPayload) {
  return unwrap<StudentDetail>(http.patch(`/students/${id}`, payload));
}

export function fetchStudentDetail(studentId: number) {
  return unwrap<StudentDetail>(http.get(`/students/${studentId}`));
}

export function createStudentGuardian(studentId: number, payload: StudentGuardianPayload) {
  return unwrap<StudentGuardian>(http.post(`/students/${studentId}/guardians`, payload));
}

export function updateStudentGuardian(
  studentId: number,
  guardianId: number,
  payload: StudentGuardianPayload,
) {
  return unwrap<StudentGuardian>(http.patch(`/students/${studentId}/guardians/${guardianId}`, payload));
}

export function deleteStudentGuardian(studentId: number, guardianId: number) {
  return unwrap<{ deleted: boolean }>(http.delete(`/students/${studentId}/guardians/${guardianId}`));
}

export function fetchClassList() {
  return unwrap<PagedResult<SchoolClass>>(http.get("/classes"));
}

export function createClass(payload: SchoolClassPayload) {
  return unwrap<SchoolClass>(http.post("/classes", payload));
}

export function updateClass(id: number, payload: SchoolClassPayload) {
  return unwrap<SchoolClass>(http.patch(`/classes/${id}`, payload));
}

export function fetchClassDetail(classId: number) {
  return unwrap<ClassDetail>(http.get(`/classes/${classId}`));
}

export function addStudentsToClass(classId: number, payload: { studentIds: number[] }) {
  return unwrap<{ added: boolean }>(http.post(`/classes/${classId}/students`, payload));
}

export function removeStudentFromClass(classId: number, studentId: number) {
  return unwrap<{ removed: boolean }>(http.delete(`/classes/${classId}/students/${studentId}`));
}

export function fetchScheduleList() {
  return unwrap<PagedResult<Schedule>>(http.get("/schedules"));
}

export function createSchedule(payload: SchedulePayload) {
  return unwrap<Schedule>(http.post("/schedules", payload));
}

export function fetchScheduleDetail(scheduleId: number) {
  return unwrap<ScheduleDetail>(http.get(`/schedules/${scheduleId}`));
}

export function updateSchedule(scheduleId: number, payload: SchedulePayload) {
  return unwrap<Schedule>(http.patch(`/schedules/${scheduleId}`, payload));
}

export function rescheduleLesson(scheduleId: number, payload: ScheduleActionPayload) {
  return unwrap<Schedule>(http.post(`/schedules/${scheduleId}/reschedule`, payload));
}

export function cancelLesson(scheduleId: number, payload: { remark: string }) {
  return unwrap<Schedule>(http.post(`/schedules/${scheduleId}/cancel`, payload));
}

export function createMakeupLesson(scheduleId: number, payload: ScheduleActionPayload) {
  return unwrap<Schedule>(http.post(`/schedules/${scheduleId}/makeup`, payload));
}

export function fetchNoticeList() {
  return unwrap<PagedResult<Notice>>(http.get("/notices"));
}

export function fetchNotice(id: number) {
  return unwrap<Notice>(http.get(`/notices/${id}`));
}

export function createNotice(payload: NoticePayload) {
  return unwrap<Notice>(http.post("/notices", payload));
}

export function updateNotice(id: number, payload: NoticePayload) {
  return unwrap<Notice>(http.patch(`/notices/${id}`, payload));
}

export function sendNotice(id: number) {
  return unwrap<Notice>(http.post(`/notices/${id}/send`));
}

export function fetchNoticeTargets(id: number) {
  return unwrap<NoticeTarget[]>(http.get(`/notices/${id}/targets`));
}

export function fetchAttendanceSessionList() {
  return unwrap<PagedResult<AttendanceSession>>(http.get("/attendance"));
}

export function fetchAttendanceRecordList(params?: {
  mode?: string;
  classId?: number;
  studentId?: number;
  date?: string;
  status?: string;
}) {
  return unwrap<PagedResult<AttendanceRecord>>(http.get("/attendance", { params }));
}

export function fetchScheduleAttendance(scheduleId: number) {
  return unwrap<AttendanceDetail>(http.get(`/schedules/${scheduleId}/attendance`));
}

export function saveScheduleAttendance(
  scheduleId: number,
  payload: { items: Array<{ studentId: number; status: string; remark: string }> },
) {
  return unwrap<{ saved: boolean }>(http.put(`/schedules/${scheduleId}/attendance`, payload));
}

export function fetchHomeworkList() {
  return unwrap<PagedResult<Homework>>(http.get("/homeworks"));
}

export function fetchScheduleHomework(scheduleId: number) {
  return unwrap<Partial<Homework>>(http.get(`/schedules/${scheduleId}/homework`));
}

export function saveScheduleHomework(
  scheduleId: number,
  payload: {
    title: string;
    content: string;
    submissionNote: string;
    status: string;
  },
) {
  return unwrap<Homework>(http.put(`/schedules/${scheduleId}/homework`, payload));
}

export function fetchScheduleFeedback(scheduleId: number) {
  return unwrap<Partial<Feedback>>(http.get(`/schedules/${scheduleId}/feedback`));
}

export function saveScheduleFeedback(
  scheduleId: number,
  payload: {
    summary: string;
    learningStatus: string;
    nextSuggestion: string;
    parentNotice: string;
  },
) {
  return unwrap<Feedback>(http.put(`/schedules/${scheduleId}/feedback`, payload));
}
