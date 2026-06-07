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
  CaretBottom,
} from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { computed, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { logout } from "../api/auth";
import { useAuthStore } from "../stores/auth";

type MenuLeaf = {
  path: string;
  label: string;
  hint: string;
  icon: unknown;
  permission: string;
};

type MenuSection = {
  key: string;
  label: string;
  hint: string;
  children: MenuLeaf[];
};

type MenuGroup = {
  key: string;
  label: string;
  hint: string;
  children: MenuSection[];
};

const menuLeaves = {
  dashboard: {
    path: "/dashboard",
    label: "首页总览",
    hint: "今天的班、课和待办",
    icon: House,
    permission: "dashboard:view",
  },
  users: {
    path: "/users",
    label: "账号管理",
    hint: "登录账号与启停状态",
    icon: Avatar,
    permission: "users:view",
  },
  roles: {
    path: "/roles",
    label: "角色权限",
    hint: "岗位能看什么、能做什么",
    icon: Tickets,
    permission: "roles:view",
  },
  operationLogs: {
    path: "/operation-logs",
    label: "操作记录",
    hint: "回看关键操作轨迹",
    icon: CollectionTag,
    permission: "operation_logs:view",
  },
  teachers: {
    path: "/teachers",
    label: "老师",
    hint: "师资分布与授课负载",
    icon: UserFilled,
    permission: "teachers:view",
  },
  students: {
    path: "/students",
    label: "学员",
    hint: "学员与家长台账",
    icon: Notebook,
    permission: "students:view",
  },
  courses: {
    path: "/courses",
    label: "课程",
    hint: "课程模板与适用阶段",
    icon: Reading,
    permission: "courses:view",
  },
  classes: {
    path: "/classes",
    label: "班级",
    hint: "班级组合与容量安排",
    icon: School,
    permission: "classes:view",
  },
  schedules: {
    path: "/schedules",
    label: "排课",
    hint: "每天的上课节奏",
    icon: Calendar,
    permission: "schedules:view",
  },
  attendance: {
    path: "/attendance",
    label: "签到",
    hint: "点名与历史回看",
    icon: Checked,
    permission: "attendance:view",
  },
  homeworks: {
    path: "/homeworks",
    label: "作业反馈",
    hint: "课后作业与班级反馈",
    icon: Memo,
    permission: "homeworks:view",
  },
  notices: {
    path: "/notices",
    label: "通知",
    hint: "消息发送与回看",
    icon: Bell,
    permission: "notices:view",
  },
} satisfies Record<string, MenuLeaf>;

const menuGroups: MenuGroup[] = [
  {
    key: "overview",
    label: "工作总览",
    hint: "今天的重点安排",
    children: [
      {
        key: "overview-console",
        label: "控制台",
        hint: "总览与待办",
        children: [menuLeaves.dashboard],
      },
    ],
  },
  {
    key: "education",
    label: "教务中心",
    hint: "人、班、课日常处理",
    children: [
      {
        key: "education-people",
        label: "人员班级",
        hint: "学员、老师、班级",
        children: [menuLeaves.students, menuLeaves.teachers, menuLeaves.classes],
      },
      {
        key: "education-curriculum",
        label: "课程排课",
        hint: "课程与课表",
        children: [menuLeaves.courses, menuLeaves.schedules],
      },
      {
        key: "education-classroom",
        label: "课堂课后",
        hint: "签到与课后",
        children: [menuLeaves.attendance, menuLeaves.homeworks],
      },
      {
        key: "education-notice",
        label: "通知协同",
        hint: "通知与沟通",
        children: [menuLeaves.notices],
      },
    ],
  },
  {
    key: "system",
    label: "平台设置",
    hint: "账号、权限、记录",
    children: [
      {
        key: "system-access",
        label: "权限控制",
        hint: "账号与权限",
        children: [menuLeaves.users, menuLeaves.roles],
      },
      {
        key: "system-audit",
        label: "系统回看",
        hint: "关键操作记录",
        children: [menuLeaves.operationLogs],
      },
    ],
  },
];

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

function filterMenuGroups(groups: MenuGroup[]) {
  return groups
    .map((group) => {
      const visibleSections = group.children
        .map((section) => {
          const visibleChildren = section.children.filter((leaf) => authStore.hasPermission(leaf.permission));
          if (visibleChildren.length === 0) {
            return null;
          }

          return {
            ...section,
            children: visibleChildren,
          };
        })
        .filter((section): section is MenuSection => section !== null);

      if (visibleSections.length === 0) {
        return null;
      }

      return {
        ...group,
        children: visibleSections,
      };
    })
    .filter((group): group is MenuGroup => group !== null);
}

function flattenMenuLeaves(groups: MenuGroup[]) {
  return groups.flatMap((group) => {
    return group.children.flatMap((section) => section.children);
  });
}

const visibleMenuGroups = computed(() => {
  return filterMenuGroups(menuGroups);
});

const visibleMenus = computed(() => {
  return flattenMenuLeaves(visibleMenuGroups.value);
});

const activeMenu = computed(() => {
  return visibleMenus.value.find((menu) => route.path.startsWith(menu.path)) ?? visibleMenus.value[0] ?? menuLeaves.dashboard;
});

const pageTitle = computed(() => {
  return String(route.meta.title ?? activeMenu.value.label);
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

const expandedGroupKey = ref<string | null>(null);
const expandedSectionKey = ref<string | null>(null);

function sectionIsActive(section: MenuSection) {
  return section.children.some((leaf) => route.path.startsWith(leaf.path));
}

function groupIsActive(group: MenuGroup) {
  return group.children.some((section) => sectionIsActive(section));
}

const activeGroupKey = computed(() => {
  return visibleMenuGroups.value.find((group) => groupIsActive(group))?.key ?? visibleMenuGroups.value[0]?.key ?? null;
});

const activeSectionKey = computed(() => {
  const currentGroup = visibleMenuGroups.value.find((group) => group.key === activeGroupKey.value);
  return currentGroup?.children.find((section) => sectionIsActive(section))?.key ?? currentGroup?.children[0]?.key ?? null;
});

function groupIsExpanded(groupKey: string) {
  return expandedGroupKey.value === groupKey;
}

function sectionIsExpanded(sectionKey: string) {
  return expandedSectionKey.value === sectionKey;
}

function preferredSectionKey(group: MenuGroup) {
  return group.children.find((section) => sectionIsActive(section))?.key ?? group.children[0]?.key ?? null;
}

function toggleGroup(group: MenuGroup) {
  if (expandedGroupKey.value === group.key) {
    expandedGroupKey.value = null;
    expandedSectionKey.value = null;
    return;
  }

  expandedGroupKey.value = group.key;
  expandedSectionKey.value = preferredSectionKey(group);
}

function toggleSection(section: MenuSection) {
  if (expandedSectionKey.value === section.key) {
    expandedSectionKey.value = null;
    return;
  }

  expandedSectionKey.value = section.key;
}

watch([visibleMenuGroups, activeGroupKey, activeSectionKey], ([groups, groupKey, sectionKey]) => {
  if (groups.length === 0) {
    expandedGroupKey.value = null;
    expandedSectionKey.value = null;
    return;
  }

  const groupVisible = groups.some((group) => group.key === expandedGroupKey.value);
  if (!groupVisible || expandedGroupKey.value === null) {
    expandedGroupKey.value = groupKey;
  }

  const currentGroup = groups.find((group) => group.key === expandedGroupKey.value);
  const sectionVisible = currentGroup?.children.some((section) => section.key === expandedSectionKey.value) ?? false;

  if (!sectionVisible || expandedSectionKey.value === null) {
    expandedSectionKey.value = sectionKey;
  }
}, { immediate: true });

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
        <div class="brand-kicker">Education Admin</div>
        <div class="brand-title">Sunrise Board</div>
      </div>

      <div class="sidebar-section-label">导航目录</div>

      <div class="nav-list">
        <section
          v-for="group in visibleMenuGroups"
          :key="group.key"
          class="nav-group"
        >
          <button
            class="nav-group-toggle"
            :class="{
              'nav-group-toggle--active': groupIsActive(group),
              'nav-group-toggle--expanded': groupIsExpanded(group.key),
            }"
            type="button"
            @click="toggleGroup(group)"
          >
            <div class="nav-group-toggle-copy">
              <strong>{{ group.label }}</strong>
            </div>

            <div class="nav-group-toggle-meta">
              <el-icon class="nav-group-toggle-icon">
                <CaretBottom />
              </el-icon>
            </div>
          </button>

          <div
            v-show="groupIsExpanded(group.key)"
            class="nav-group-body"
          >
            <article
              v-for="section in group.children"
              :key="section.key"
              class="nav-section"
              :class="{
                'nav-section--active': sectionIsActive(section),
                'nav-section--expanded': sectionIsExpanded(section.key),
              }"
            >
              <button
                class="nav-section-header nav-section-toggle"
                type="button"
                @click="toggleSection(section)"
              >
                <div class="nav-section-copy">
                  <strong>{{ section.label }}</strong>
                </div>
                <div class="nav-section-meta">
                  <el-icon class="nav-section-toggle-icon">
                    <CaretBottom />
                  </el-icon>
                </div>
              </button>

              <div
                v-show="sectionIsExpanded(section.key)"
                class="nav-sublinks"
              >
                <RouterLink
                  v-for="menu in section.children"
                  :key="menu.path"
                  class="nav-sublink"
                  :to="menu.path"
                >
                  <span class="nav-sublink-icon">
                    <el-icon>
                      <component :is="menu.icon" />
                    </el-icon>
                  </span>
                  <span class="nav-sublink-copy">
                    <strong>{{ menu.label }}</strong>
                  </span>
                </RouterLink>
              </div>
            </article>
          </div>
        </section>
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
              {{ roleLabel }}
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
          <h1>{{ pageTitle }}</h1>
          <div class="topbar-meta">
            <span class="topbar-chip topbar-chip--primary">{{ activeMenu.label }}</span>
            <span class="topbar-chip">{{ roleLabel }}</span>
          </div>
        </div>
      </header>

      <main class="content">
        <RouterView />
      </main>
    </div>
  </div>
</template>
