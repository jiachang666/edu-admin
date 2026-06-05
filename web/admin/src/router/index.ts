import { createRouter, createWebHistory } from "vue-router";
import AdminLayout from "../layouts/AdminLayout.vue";
import { useAuthStore } from "../stores/auth";
import LoginView from "../views/auth/LoginView.vue";
import DashboardView from "../views/dashboard/DashboardView.vue";
import TeacherListView from "../views/teachers/TeacherListView.vue";
import StudentListView from "../views/students/StudentListView.vue";
import CourseListView from "../views/courses/CourseListView.vue";
import ClassListView from "../views/classes/ClassListView.vue";
import ScheduleListView from "../views/schedules/ScheduleListView.vue";
import AttendanceView from "../views/attendance/AttendanceView.vue";
import NoticeListView from "../views/notices/NoticeListView.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/login",
      name: "login",
      component: LoginView,
      meta: {
        title: "登录",
        eyebrow: "Campus Access",
        description: "进入培训机构的日常运营工作台。",
      },
    },
    {
      path: "/",
      component: AdminLayout,
      children: [
        { path: "", redirect: "/dashboard" },
        {
          path: "dashboard",
          name: "dashboard",
          component: DashboardView,
          meta: {
            title: "首页总览",
            eyebrow: "Operations Deck",
            description: "把今天的课程、待办、通知和班级节奏放在同一块总览板上。",
          },
        },
        {
          path: "teachers",
          name: "teachers",
          component: TeacherListView,
          meta: {
            title: "老师管理",
            eyebrow: "Faculty Ledger",
            description: "统一查看老师的科目、校区分布和当前授课负载。",
          },
        },
        {
          path: "students",
          name: "students",
          component: StudentListView,
          meta: {
            title: "学员管理",
            eyebrow: "Student Ledger",
            description: "把学员、家长、班级归属和续费状态收拢到同一份台账里。",
          },
        },
        {
          path: "courses",
          name: "courses",
          component: CourseListView,
          meta: {
            title: "课程管理",
            eyebrow: "Curriculum Board",
            description: "先把课程模板定清楚，后面的建班、排课和通知才会顺。",
          },
        },
        {
          path: "classes",
          name: "classes",
          component: ClassListView,
          meta: {
            title: "班级管理",
            eyebrow: "Classroom Matrix",
            description: "关注班级组合、容量安排和固定上课节奏。",
          },
        },
        {
          path: "schedules",
          name: "schedules",
          component: ScheduleListView,
          meta: {
            title: "排课管理",
            eyebrow: "Schedule Rail",
            description: "看清楚每天谁来上课、什么时候上、在哪里上。",
          },
        },
        {
          path: "attendance",
          name: "attendance",
          component: AttendanceView,
          meta: {
            title: "签到管理",
            eyebrow: "Attendance Desk",
            description: "把待签到班级、点名结果和历史记录收进一个顺手的签到工作台。",
          },
        },
        {
          path: "notices",
          name: "notices",
          component: NoticeListView,
          meta: {
            title: "通知管理",
            eyebrow: "Notice Dispatch",
            description: "追踪消息范围、发送状态和还没处理完的提醒。",
          },
        },
      ],
    },
  ],
});

router.beforeEach((to) => {
  const authStore = useAuthStore();
  const isLoginPage = to.path === "/login";

  if (!isLoginPage && !authStore.isLoggedIn) {
    return "/login";
  }

  if (isLoginPage && authStore.isLoggedIn) {
    return "/dashboard";
  }

  return true;
});

export default router;
