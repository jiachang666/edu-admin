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

async function openAttendance(scheduleId: number) {
  await router.push({
    path: "/attendance",
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
        <span class="section-kicker">Daily Pulse</span>
        <h2>
          今天共有 {{ overview.todayCourses }} 节课，其中还有
          {{ overview.todayPendingCheck }} 节待签到。
        </h2>
        <p>
          这块首页先帮你把“今天最值得先处理什么”直接摆出来，避免教务和负责人一进后台还要自己翻来翻去找重点。
        </p>
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
        <div>
          <h2>运营概况</h2>
          <p class="soft-text">先用演示数据把核心运营面板跑通。</p>
        </div>
        <div class="section-note">演示环境</div>
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

    <div class="panel-grid panel-grid--dashboard">
      <section class="page-card page-card--table">
        <div class="page-header">
          <div>
            <h3>近期课表</h3>
            <p class="soft-text">帮助教务快速确认今天和明天的安排。</p>
          </div>
          <div class="section-note">排课视图</div>
        </div>

        <div class="data-table-shell">
          <el-table v-loading="loading" :data="overview.upcomingLessons" stripe>
            <el-table-column label="日期" prop="lessonDate" width="120" />
            <el-table-column label="时间" prop="lessonTime" width="140" />
            <el-table-column label="班级" prop="className" min-width="180" />
            <el-table-column label="课程" prop="courseName" width="140" />
            <el-table-column label="老师" prop="teacherName" width="120" />
            <el-table-column label="教室" prop="classroom" width="120" />
            <el-table-column label="状态" width="120">
              <template #default="{ row }">
                <el-tag>{{ row.attendanceStatus }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="处理" width="120" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openAttendance(row.id)">
                  {{ row.attendanceStatus === "待签到" ? "去签到" : "看记录" }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </section>

      <section class="page-card page-card--table">
        <div class="page-header">
          <div>
            <h3>最近通知</h3>
            <p class="soft-text">方便运营和班主任查看还没处理完的消息。</p>
          </div>
          <div class="section-note">消息视图</div>
        </div>

        <div class="data-table-shell">
          <el-table v-loading="loading" :data="overview.latestNotices" stripe>
            <el-table-column label="标题" prop="title" min-width="220" />
            <el-table-column label="分类" prop="category" width="140" />
            <el-table-column label="发送范围" prop="targetScope" min-width="180" />
            <el-table-column label="状态" width="120">
              <template #default="{ row }">
                <el-tag :type="row.status === '已发送' ? 'success' : 'warning'">
                  {{ row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="时间" prop="publishAt" width="180" />
          </el-table>
        </div>
      </section>
    </div>
  </div>
</template>
