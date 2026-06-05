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

export type Schedule = {
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
};

export type Notice = {
  id: number;
  title: string;
  category: string;
  targetScope: string;
  status: string;
  publishAt: string;
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

export function fetchClassList() {
  return unwrap<PagedResult<SchoolClass>>(http.get("/classes"));
}

export function fetchScheduleList() {
  return unwrap<PagedResult<Schedule>>(http.get("/schedules"));
}

export function fetchNoticeList() {
  return unwrap<PagedResult<Notice>>(http.get("/notices"));
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
