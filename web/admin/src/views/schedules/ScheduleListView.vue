<script setup lang="ts">
import { ElMessage } from "element-plus";
import { computed, onMounted, ref } from "vue";
import { fetchScheduleList, type Schedule } from "../../api/education";

const loading = ref(false);
const schedules = ref<Schedule[]>([]);

const doneCount = computed(() => {
  return schedules.value.filter((item) => item.attendanceStatus === "已完成").length;
});

const pendingCount = computed(() => {
  return schedules.value.filter((item) => item.attendanceStatus === "待签到").length;
});

const campusCount = computed(() => {
  return new Set(schedules.value.map((item) => item.campus)).size;
});

async function loadSchedules() {
  loading.value = true;

  try {
    const result = await fetchScheduleList();
    schedules.value = result.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("排课列表加载失败");
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  void loadSchedules();
});
</script>

<template>
  <div class="page-stack">
    <section class="page-hero">
      <div class="page-hero__copy">
        <span class="section-kicker">Schedule Rail</span>
        <h2>把每天什么时候上、谁来上、在哪上，用一条更清楚的排课视图串起来。</h2>
        <p>
          目前先展示时间、教室和签到进度，后面继续补调课和补课操作时，这块会成为整套流程的时间轴中心。
        </p>
      </div>

      <div class="metric-strip">
        <article class="metric-tile">
          <span>排课总数</span>
          <strong>{{ schedules.length }}</strong>
          <small>当前演示环境可见的全部课程安排</small>
        </article>
        <article class="metric-tile">
          <span>已完成</span>
          <strong>{{ doneCount }}</strong>
          <small>已经完成签到或课程闭环的安排</small>
        </article>
        <article class="metric-tile">
          <span>待签到</span>
          <strong>{{ pendingCount }}</strong>
          <small>优先处理正在等老师确认的课程</small>
        </article>
        <article class="metric-tile">
          <span>涉及校区</span>
          <strong>{{ campusCount }}</strong>
          <small>当前排课已经覆盖的校区数量</small>
        </article>
      </div>
    </section>

    <section class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>排课列表</h2>
          <p class="soft-text">先把上课时间、教室和签到状态展示出来。</p>
        </div>
        <div class="section-note">排课视图</div>
      </div>

      <div class="data-table-shell">
        <el-table v-loading="loading" :data="schedules" stripe>
          <el-table-column label="日期" prop="lessonDate" width="120" />
          <el-table-column label="时间" prop="lessonTime" width="140" />
          <el-table-column label="班级" prop="className" min-width="180" />
          <el-table-column label="课程" prop="courseName" width="140" />
          <el-table-column label="老师" prop="teacherName" width="120" />
          <el-table-column label="校区" prop="campus" width="120" />
          <el-table-column label="教室" prop="classroom" width="120" />
          <el-table-column label="签到状态" width="120">
            <template #default="{ row }">
              <el-tag
                :type="row.attendanceStatus === '已完成' ? 'success' : 'warning'"
              >
                {{ row.attendanceStatus }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </section>
  </div>
</template>
