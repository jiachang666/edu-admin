<script setup lang="ts">
import {
  Avatar,
  Bell,
  Calendar,
  Checked,
  CollectionTag,
  House,
  Memo,
  Notebook,
  Reading,
  School,
  SwitchButton,
  Tickets,
  UserFilled,
} from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { logout } from "../api/auth";
import { useAuthStore } from "../stores/auth";

const menus = [
  {
    path: "/dashboard",
    label: "首页总览",
    hint: "今天的班、课和待办",
    icon: House,
    permission: "dashboard:view",
  },
  {
    path: "/users",
    label: "账号管理",
    hint: "登录账号与启停状态",
    icon: Avatar,
    permission: "users:view",
  },
  {
    path: "/roles",
    label: "角色权限",
    hint: "岗位能看什么、能做什么",
    icon: Tickets,
    permission: "roles:view",
  },
  {
    path: "/operation-logs",
    label: "操作记录",
    hint: "回看关键操作轨迹",
    icon: CollectionTag,
    permission: "operation_logs:view",
  },
  {
    path: "/teachers",
    label: "老师",
    hint: "师资分布与授课负载",
    icon: UserFilled,
    permission: "teachers:view",
  },
  {
    path: "/students",
    label: "学员",
    hint: "学员与家长台账",
    icon: Notebook,
    permission: "students:view",
  },
  {
    path: "/courses",
    label: "课程",
    hint: "课程模板与适用阶段",
    icon: Reading,
    permission: "courses:view",
  },
  {
    path: "/classes",
    label: "班级",
    hint: "班级组合与容量安排",
    icon: School,
    permission: "classes:view",
  },
  {
    path: "/schedules",
    label: "排课",
    hint: "每天的上课节奏",
    icon: Calendar,
    permission: "schedules:view",
  },
  {
    path: "/attendance",
    label: "签到",
    hint: "点名与历史回看",
    icon: Checked,
    permission: "attendance:view",
  },
  {
    path: "/homeworks",
    label: "作业反馈",
    hint: "课后作业与班级反馈",
    icon: Memo,
    permission: "homeworks:view",
  },
  {
    path: "/notices",
    label: "通知",
    hint: "消息发送与回看",
    icon: Bell,
    permission: "notices:view",
  },
];

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const visibleMenus = computed(() => {
  return menus.filter((menu) => authStore.hasPermission(menu.permission));
});

const activeMenu = computed(() => {
  return visibleMenus.value.find((menu) => route.path.startsWith(menu.path)) ?? visibleMenus.value[0] ?? menus[0];
});

const pageTitle = computed(() => {
  return String(route.meta.title ?? activeMenu.value.label);
});

const pageEyebrow = computed(() => {
  return String(route.meta.eyebrow ?? "Campus Console");
});

const pageDescription = computed(() => {
  return String(route.meta.description ?? activeMenu.value.hint);
});

const formattedDate = computed(() => {
  return new Intl.DateTimeFormat("zh-CN", {
    month: "long",
    day: "numeric",
    weekday: "short",
  }).format(new Date());
});

const userInitial = computed(() => {
  const name = authStore.userName.trim();
  if (name.length === 0) {
    return "A";
  }

  return name.slice(0, 1).toUpperCase();
});

const roleLabel = computed(() => {
  return authStore.roleNames[0] || "演示账号";
});

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
      <div class="sidebar-glow sidebar-glow--top"></div>
      <div class="sidebar-glow sidebar-glow--bottom"></div>

      <div class="brand-panel">
        <div class="brand-kicker">Edu Admin</div>
        <div class="brand-title">Schoolhouse Ledger</div>
        <p class="brand-copy">
          为培训机构整理每天的人、班、课与通知，让教务工作看起来更从容。
        </p>
        <div class="brand-stamp">
          <el-icon><CollectionTag /></el-icon>
          <span>开源演示版</span>
        </div>
      </div>

      <div class="sidebar-section-label">工作区导航</div>

      <div class="nav-list">
        <RouterLink
          v-for="menu in visibleMenus"
          :key="menu.path"
          class="nav-link"
          :to="menu.path"
        >
          <span class="nav-icon">
            <el-icon>
              <component :is="menu.icon" />
            </el-icon>
          </span>
          <span class="nav-copy">
            <strong>{{ menu.label }}</strong>
            <small>{{ menu.hint }}</small>
          </span>
        </RouterLink>
      </div>

      <div class="sidebar-footer">
        <div class="sidebar-user-card">
          <div class="sidebar-user-avatar">
            {{ userInitial }}
          </div>
          <div class="sidebar-user">
            <div class="sidebar-user-label">当前账号</div>
            <div class="sidebar-user-name">
              {{ authStore.userName || "System Admin" }}
            </div>
            <div class="sidebar-user-caption">
              {{ roleLabel }} · 当前菜单会按权限自动收起
            </div>
          </div>
        </div>

        <el-button class="sidebar-ghost-button" plain @click="handleLogout">
          <el-icon><SwitchButton /></el-icon>
          <span>退出登录</span>
        </el-button>
      </div>
    </aside>

    <div class="app-stage">
      <header class="app-topbar">
        <div class="topbar-copy">
          <span class="topbar-eyebrow">{{ pageEyebrow }}</span>
          <h1>{{ pageTitle }}</h1>
          <p>{{ pageDescription }}</p>
        </div>

        <div class="topbar-cards">
          <article class="topbar-card">
            <span class="topbar-card-label">今日节奏</span>
            <strong>{{ formattedDate }}</strong>
            <small>适合先处理排课、签到和通知节奏</small>
          </article>

          <article class="topbar-card topbar-card--accent">
            <span class="topbar-card-label">当前焦点</span>
            <strong>{{ activeMenu.label }}</strong>
            <small>{{ activeMenu.hint }}</small>
          </article>
        </div>
      </header>

      <main class="content">
        <RouterView />
      </main>
    </div>
  </div>
</template>
