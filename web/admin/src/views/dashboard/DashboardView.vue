<script setup lang="ts">
import { ElMessage } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { fetchDashboardOverview, type DashboardOverview } from "../../api/education";

const loading = ref(false);
const router = useRouter();
const overview = reactive<DashboardOverview>({
  todayCourses: 0,
  todayPendingCheck: 0,
  todayLeaveCount: 0,
  todayAbsentCount: 0,
  studentCount: 0,
  classCount: 0,
  pendingActionCount: 0,
  upcomingLessons: [],
  latestNotices: [],
});

const summaryCards = computed(() => [
  { label: "今日课程", value: overview.todayCourses, tone: "blue" },
  { label: "待签到", value: overview.todayPendingCheck, tone: "orange" },
  { label: "请假人数", value: overview.todayLeaveCount, tone: "teal" },
  { label: "在读学员", value: overview.studentCount, tone: "green" },
  { label: "开班数量", value: overview.classCount, tone: "indigo" },
  { label: "待处理事项", value: overview.pendingActionCount, tone: "red" },
]);

const signalCards = computed(() => [
  {
    label: "今日节奏",
    value: `${overview.todayCourses} 节课`,
    description: "先盯紧今天需要签到和即将开始的班级。",
  },
  {
    label: "运营焦点",
    value: `${overview.pendingActionCount} 条待处理`,
    description: "通知草稿、请假处理和临时调整建议今天清掉。",
  },
  {
    label: "学员规模",
    value: `${overview.studentCount} 位在库学员`,
    description: "当前总览更适合负责人和教务快速盘点全局。",
  },
]);

const quickActions = [
  {
    label: "进入签到",
    description: "先处理今天还没完成点名的班级。",
    path: "/attendance",
  },
  {
    label: "查看排课",
    description: "快速确认接下来几节课的时间和教室。",
    path: "/schedules",
  },
  {
    label: "发布通知",
    description: "把调课、放假或提醒消息尽快发出去。",
    path: "/notices",
  },
];

function openPath(path: string) {
  void router.push(path);
}

function noticeStatusType(status: string) {
  if (status === "已发送") {
    return "success";
  }

  if (status === "草稿") {
    return "info";
  }

  return "warning";
}

async function openAttendance(scheduleId: number) {
  await router.push({
    path: "/attendance",
    query: { scheduleId: String(scheduleId) },
  });
}

async function openHomework(scheduleId: number) {
  await router.push({
    path: "/homeworks",
    query: { scheduleId: String(scheduleId) },
  });
}

async function loadOverview() {
  loading.value = true;

  try {
    const result = await fetchDashboardOverview();
    Object.assign(overview, result);
  } catch (error) {
    console.error(error);
    ElMessage.error("首页数据加载失败");
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  void loadOverview();
});
</script>

<template>
  <div class="page-stack">
    <section class="hero-board">
      <div class="hero-board__main">
        <h2>
          今天共有 {{ overview.todayCourses }} 节课，其中还有
          {{ overview.todayPendingCheck }} 节待签到。
        </h2>

        <div class="hero-actions">
          <button
            v-for="action in quickActions"
            :key="action.path"
            class="hero-action-card"
            type="button"
            @click="openPath(action.path)"
          >
            <strong>{{ action.label }}</strong>
            <span>{{ action.description }}</span>
          </button>
        </div>
      </div>

      <div class="hero-board__side">
        <article
          v-for="signal in signalCards"
          :key="signal.label"
          class="signal-card"
        >
          <span>{{ signal.label }}</span>
          <strong>{{ signal.value }}</strong>
          <p>{{ signal.description }}</p>
        </article>
      </div>
    </section>

    <section class="page-card">
      <div class="page-header">
        <h2>今日概况</h2>
      </div>

      <div v-loading="loading" class="stats-grid">
        <article
          v-for="item in summaryCards"
          :key="item.label"
          class="stat-card"
          :data-tone="item.tone"
        >
          <div class="stat-label">{{ item.label }}</div>
          <div class="stat-value">{{ item.value }}</div>
        </article>
      </div>
    </section>

    <div class="dashboard-grid">
      <section class="page-card dashboard-panel">
        <div class="page-header">
          <h3>近期课表</h3>
        </div>

        <div v-loading="loading" class="agenda-list">
          <div v-if="overview.upcomingLessons.length === 0" class="soft-empty">
            还没有可展示的排课安排。
          </div>

          <article
            v-for="lesson in overview.upcomingLessons"
            v-else
            :key="lesson.id"
            class="agenda-item"
          >
            <div class="agenda-time">
              <strong>{{ lesson.lessonTime }}</strong>
              <span>{{ lesson.lessonDate }}</span>
            </div>

            <div class="agenda-copy">
              <strong>{{ lesson.className }}</strong>
              <p>{{ lesson.courseName }} · {{ lesson.teacherName }} · {{ lesson.classroom }}</p>
            </div>

            <div class="agenda-meta">
              <el-tag :type="lesson.attendanceStatus === '待签到' ? 'warning' : 'success'">
                {{ lesson.attendanceStatus }}
              </el-tag>
              <div class="agenda-actions">
                <el-button link type="primary" @click="openAttendance(lesson.id)">
                  {{ lesson.attendanceStatus === "待签到" ? "去签到" : "看记录" }}
                </el-button>
                <el-button link type="primary" @click="openHomework(lesson.id)">
                  作业反馈
                </el-button>
              </div>
            </div>
          </article>
        </div>
      </section>

      <section class="page-card dashboard-panel">
        <div class="page-header">
          <h3>最近通知</h3>
        </div>

        <div v-loading="loading" class="notice-feed">
          <div v-if="overview.latestNotices.length === 0" class="soft-empty">
            还没有可展示的通知记录。
          </div>

          <article
            v-for="notice in overview.latestNotices"
            v-else
            :key="`${notice.title}-${notice.publishAt}`"
            class="notice-feed__item"
          >
            <div class="notice-feed__main">
              <strong>{{ notice.title }}</strong>
              <p>{{ notice.targetScope }} · {{ notice.category }}</p>
            </div>

            <div class="notice-feed__meta">
              <el-tag :type="noticeStatusType(notice.status)">
                {{ notice.status }}
              </el-tag>
              <small>{{ notice.publishAt }}</small>
            </div>
          </article>
        </div>
      </section>
    </div>
  </div>
</template>
