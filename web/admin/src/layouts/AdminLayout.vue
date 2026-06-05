<script setup lang="ts">
import { ElMessage } from "element-plus";
import { useRouter } from "vue-router";
import { logout } from "../api/auth";
import { useAuthStore } from "../stores/auth";

const menus = [
  { path: "/dashboard", label: "首页" },
  { path: "/teachers", label: "老师" },
  { path: "/students", label: "学员" },
  { path: "/courses", label: "课程" },
  { path: "/classes", label: "班级" },
  { path: "/schedules", label: "排课" },
  { path: "/notices", label: "通知" },
];

const router = useRouter();
const authStore = useAuthStore();

async function handleLogout() {
  try {
    await logout();
  } catch (error) {
    console.warn("logout request failed", error);
  }

  authStore.clearSession();
  await router.push("/login");
  ElMessage.success("已退出登录");
}
</script>

<template>
  <div class="admin-layout">
    <aside class="sidebar">
      <div class="brand">Edu Admin</div>
      <div class="nav-list">
        <RouterLink
          v-for="menu in menus"
          :key="menu.path"
          class="nav-link"
          :to="menu.path"
        >
          {{ menu.label }}
        </RouterLink>
      </div>
      <div class="sidebar-footer">
        <div class="sidebar-user">
          <div class="sidebar-user-label">当前账号</div>
          <div class="sidebar-user-name">
            {{ authStore.userName || "System Admin" }}
          </div>
        </div>
        <el-button plain @click="handleLogout">退出登录</el-button>
      </div>
    </aside>
    <main class="content">
      <RouterView />
    </main>
  </div>
</template>
