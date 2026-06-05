<script setup lang="ts">
import { ElMessage } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import { fetchDashboardOverview, type DashboardOverview } from "../../api/education";

const loading = ref(false);
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
    <section class="page-card">
      <div class="page-header">
        <div>
          <h2>首页概况</h2>
          <p class="soft-text">先用演示数据把核心运营面板跑通。</p>
        </div>
        <el-tag type="primary">演示版</el-tag>
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

    <section class="page-card">
      <div class="page-header">
        <div>
          <h3>近期课表</h3>
          <p class="soft-text">帮助教务快速确认今天和明天的安排。</p>
        </div>
      </div>

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
      </el-table>
    </section>

    <section class="page-card">
      <div class="page-header">
        <div>
          <h3>最近通知</h3>
          <p class="soft-text">方便运营和班主任查看还没处理完的消息。</p>
        </div>
      </div>

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
    </section>
  </div>
</template>
