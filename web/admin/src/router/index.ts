import { createRouter, createWebHistory } from "vue-router";
import AdminLayout from "../layouts/AdminLayout.vue";
import { useAuthStore } from "../stores/auth";
import LoginView from "../views/auth/LoginView.vue";
import DashboardView from "../views/dashboard/DashboardView.vue";
import TeacherListView from "../views/teachers/TeacherListView.vue";
import TeacherDetailView from "../views/teachers/TeacherDetailView.vue";
import StudentListView from "../views/students/StudentListView.vue";
import StudentDetailView from "../views/students/StudentDetailView.vue";
import CourseListView from "../views/courses/CourseListView.vue";
import ClassListView from "../views/classes/ClassListView.vue";
import ClassDetailView from "../views/classes/ClassDetailView.vue";
import ScheduleListView from "../views/schedules/ScheduleListView.vue";
import ScheduleDetailView from "../views/schedules/ScheduleDetailView.vue";
import AttendanceView from "../views/attendance/AttendanceView.vue";
import HomeworkView from "../views/homeworks/HomeworkView.vue";
import NoticeListView from "../views/notices/NoticeListView.vue";
import UserListView from "../views/users/UserListView.vue";
import RoleListView from "../views/roles/RoleListView.vue";
import OperationLogView from "../views/audit/OperationLogView.vue";

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
          path: "users",
          name: "users",
          component: UserListView,
          meta: {
            title: "账号管理",
            eyebrow: "Access Desk",
            description: "统一管理后台账号、使用身份和启停状态，确保谁能进系统一目了然。",
          },
        },
        {
          path: "roles",
          name: "roles",
          component: RoleListView,
          meta: {
            title: "角色权限",
            eyebrow: "Role Matrix",
            description: "把不同岗位能看到什么、能操作什么收拢到同一张权限工作台里。",
          },
        },
        {
          path: "operation-logs",
          name: "operation-logs",
          component: OperationLogView,
          meta: {
            title: "操作记录",
            eyebrow: "Audit Trail",
            description: "回看关键操作是谁做的、做了什么、是什么时候发生的。",
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
          path: "teachers/:id",
          name: "teacher-detail",
          component: TeacherDetailView,
          meta: {
            title: "老师详情",
            eyebrow: "Teacher Hub",
            description: "把老师的基本资料、负责班级和近期课程安排集中在同一页里，方便教务快速判断谁在带什么班。",
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
          path: "students/:id",
          name: "student-detail",
          component: StudentDetailView,
          meta: {
            title: "学员详情",
            eyebrow: "Student Hub",
            description: "把学员、家长、班级、上课和课后跟进收进同一页，方便老师和教务快速接手。",
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
          path: "classes/:id",
          name: "class-detail",
          component: ClassDetailView,
          meta: {
            title: "班级详情",
            eyebrow: "Class Hub",
            description: "把班级里的学员、排课、签到、作业和通知收进同一个业务中心页。",
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
          path: "schedules/:id",
          name: "schedule-detail",
          component: ScheduleDetailView,
          meta: {
            title: "课程安排详情",
            eyebrow: "Lesson Hub",
            description: "把一次具体上课的签到、作业、反馈和相关通知入口集中起来。",
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
          path: "homeworks",
          name: "homeworks",
          component: HomeworkView,
          meta: {
            title: "作业与反馈",
            eyebrow: "Homework Studio",
            description: "围绕每次上课安排快速布置作业，并补上班级级别的课后反馈。",
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

function normalizePathSegment(path: string) {
  const segments = path.split("/").filter(Boolean);
  if (segments.length === 0) {
    return "/dashboard";
  }

  return `/${segments[0]}`;
}

const routePermissionMap: Record<string, string> = {
  "/dashboard": "dashboard:view",
  "/users": "users:view",
  "/roles": "roles:view",
  "/operation-logs": "operation_logs:view",
  "/teachers": "teachers:view",
  "/students": "students:view",
  "/courses": "courses:view",
  "/classes": "classes:view",
  "/schedules": "schedules:view",
  "/attendance": "attendance:view",
  "/homeworks": "homeworks:view",
  "/notices": "notices:view",
};

function firstAllowedPath(authStore: ReturnType<typeof useAuthStore>) {
  const orderedPaths = [
    "/dashboard",
    "/users",
    "/roles",
    "/operation-logs",
    "/teachers",
    "/students",
    "/courses",
    "/classes",
    "/schedules",
    "/attendance",
    "/homeworks",
    "/notices",
  ];

  const matchedPath = orderedPaths.find((path) => {
    const permission = routePermissionMap[path];
    return !permission || authStore.hasPermission(permission);
  });

  return matchedPath ?? "/dashboard";
}

router.beforeEach((to) => {
  const authStore = useAuthStore();
  const isLoginPage = to.path === "/login";

  if (!isLoginPage && !authStore.isLoggedIn) {
    return "/login";
  }

  if (isLoginPage && authStore.isLoggedIn) {
    return "/dashboard";
  }

  if (!isLoginPage && authStore.token && authStore.permissions.length === 0) {
    return authStore
      .hydrateSession()
      .then(() => true)
      .catch(() => {
        authStore.clearSession();
        return "/login";
      });
  }

  if (!isLoginPage) {
    const basePath = normalizePathSegment(to.path);
    const requiredPermission = routePermissionMap[basePath];
    if (requiredPermission && !authStore.hasPermission(requiredPermission)) {
      return firstAllowedPath(authStore);
    }
  }

  return true;
});

export default router;
