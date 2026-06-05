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
  mainSubject: string;
  employmentType: string;
  weeklyHours: number;
  campus: string;
  status: string;
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
  name: string;
  courseName: string;
  teacherId: number;
  teacherName: string;
  campus: string;
  studentCount: number;
  capacity: number;
  weeklySchedule: string;
  status: string;
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

async function unwrap<T>(request: Promise<{ data: ApiEnvelope<T> }>) {
  const response = await request;
  return response.data.data;
}

export function fetchDashboardOverview() {
  return unwrap<DashboardOverview>(http.get("/dashboard/overview"));
}

export function fetchTeacherList() {
  return unwrap<PagedResult<Teacher>>(http.get("/teachers"));
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

export function fetchStudentList() {
  return unwrap<PagedResult<Student>>(http.get("/students"));
}

export function fetchStudentDetail(studentId: number) {
  return unwrap<StudentDetail>(http.get(`/students/${studentId}`));
}

export function fetchClassList() {
  return unwrap<PagedResult<SchoolClass>>(http.get("/classes"));
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
