import { createRouter, createWebHistory } from "vue-router";
import AdminLayout from "../layouts/AdminLayout.vue";
import { useAuthStore } from "../stores/auth";
import LoginView from "../views/auth/LoginView.vue";
import DashboardView from "../views/dashboard/DashboardView.vue";
import TeacherListView from "../views/teachers/TeacherListView.vue";
import StudentListView from "../views/students/StudentListView.vue";
import ClassListView from "../views/classes/ClassListView.vue";
import ScheduleListView from "../views/schedules/ScheduleListView.vue";
import NoticeListView from "../views/notices/NoticeListView.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/login",
      name: "login",
      component: LoginView,
    },
    {
      path: "/",
      component: AdminLayout,
      children: [
        { path: "", redirect: "/dashboard" },
        { path: "dashboard", name: "dashboard", component: DashboardView },
        { path: "teachers", name: "teachers", component: TeacherListView },
        { path: "students", name: "students", component: StudentListView },
        { path: "classes", name: "classes", component: ClassListView },
        { path: "schedules", name: "schedules", component: ScheduleListView },
        { path: "notices", name: "notices", component: NoticeListView },
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
