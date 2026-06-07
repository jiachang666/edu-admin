<script setup lang="ts">
import { ElMessage } from "element-plus";
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { fetchTeacherDetail, type Schedule, type TeacherDetail } from "../../api/education";

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const teacherDetail = ref<TeacherDetail | null>(null);

const teacherId = computed(() => {
  const parsedValue = Number(route.params.id);
  if (!Number.isInteger(parsedValue) || parsedValue <= 0) {
    return null;
  }

  return parsedValue;
});

const currentTeacher = computed(() => {
  return teacherDetail.value?.teacher ?? null;
});

const currentClasses = computed(() => {
  return teacherDetail.value?.classes ?? [];
});

const currentSchedules = computed(() => {
  return teacherDetail.value?.recentSchedules ?? [];
});

const runningClassCount = computed(() => {
  return currentClasses.value.filter((item) => item.status === "开班中").length;
});

const campusCoverage = computed(() => {
  return new Set(currentClasses.value.map((item) => item.campus).filter(Boolean)).size;
});

const upcomingScheduleCount = computed(() => {
  return currentSchedules.value.filter((item) => scheduleIsTodayOrLater(item)).length;
});

async function loadTeacherDetail() {
  const currentTeacherId = teacherId.value;
  if (currentTeacherId === null) {
    ElMessage.error("老师编号不正确");
    return;
  }

  loading.value = true;

  try {
    teacherDetail.value = await fetchTeacherDetail(currentTeacherId);
  } catch (error) {
    console.error(error);
    ElMessage.error("老师详情加载失败");
  } finally {
    loading.value = false;
  }
}

function scheduleIsTodayOrLater(item: Schedule) {
  const dateText = item.lessonDate.trim();
  if (!dateText) {
    return false;
  }

  const lessonDate = new Date(`${dateText}T00:00:00`);
  if (Number.isNaN(lessonDate.getTime())) {
    return false;
  }

  const today = new Date();
  today.setHours(0, 0, 0, 0);

  return lessonDate.getTime() >= today.getTime();
}

function scheduleStatusTone(status: string) {
  switch (status) {
    case "已完成":
      return "success";
    case "待签到":
      return "warning";
    case "待上课":
      return "primary";
    case "已停课":
      return "danger";
    default:
      return "info";
  }
}

function openClassDetail(classId: number) {
  void router.push(`/classes/${classId}`);
}

function openScheduleDetail(scheduleId: number) {
  void router.push(`/schedules/${scheduleId}`);
}

function openScheduleList() {
  void router.push("/schedules");
}

onMounted(() => {
  void loadTeacherDetail();
});
</script>

<template>
  <div v-loading="loading" class="page-stack">
    <section class="page-card">
      <div class="page-header">
        <div>
          <h2>{{ currentTeacher?.name || "老师详情" }}</h2>
          <p class="soft-text">
            {{ currentTeacher?.mainSubject || "未填写科目" }} ·
            {{ currentTeacher?.campus || "未填写校区" }} ·
            {{ currentTeacher?.status || "未填写状态" }}
          </p>
        </div>
        <div class="page-actions">
          <el-button plain @click="openScheduleList">查看排课总表</el-button>
        </div>
      </div>

      <div class="metric-strip metric-strip--compact">
        <article class="metric-tile">
          <span>负责班级</span>
          <strong>{{ currentClasses.length }}</strong>
        </article>
        <article class="metric-tile">
          <span>开班中</span>
          <strong>{{ runningClassCount }}</strong>
        </article>
        <article class="metric-tile">
          <span>近期课程</span>
          <strong>{{ currentSchedules.length }}</strong>
        </article>
        <article class="metric-tile">
          <span>待上课</span>
          <strong>{{ upcomingScheduleCount }}</strong>
        </article>
      </div>

      <div class="detail-info-grid">
        <article class="detail-info-card">
          <span>姓名</span>
          <strong>{{ currentTeacher?.name || "-" }}</strong>
        </article>
        <article class="detail-info-card">
          <span>职级 / 标签</span>
          <strong>{{ currentTeacher?.title || "未填写" }}</strong>
        </article>
        <article class="detail-info-card">
          <span>主教科目</span>
          <strong>{{ currentTeacher?.mainSubject || "-" }}</strong>
        </article>
        <article class="detail-info-card">
          <span>用工类型</span>
          <strong>{{ currentTeacher?.employmentType || "-" }}</strong>
        </article>
        <article class="detail-info-card">
          <span>所属校区</span>
          <strong>{{ currentTeacher?.campus || "-" }}</strong>
        </article>
        <article class="detail-info-card">
          <span>老师状态</span>
          <strong>{{ currentTeacher?.status || "-" }}</strong>
        </article>
      </div>

      <div class="stats-grid stats-grid--compact">
        <article class="stat-card" data-tone="blue">
          <span class="stat-label">周课时</span>
          <strong class="stat-value">{{ currentTeacher?.weeklyHours ?? 0 }}</strong>
        </article>
        <article class="stat-card" data-tone="teal">
          <span class="stat-label">覆盖校区</span>
          <strong class="stat-value">{{ campusCoverage }}</strong>
        </article>
        <article class="stat-card" data-tone="orange">
          <span class="stat-label">已排课程</span>
          <strong class="stat-value">{{ currentSchedules.length }}</strong>
        </article>
        <article class="stat-card" data-tone="green">
          <span class="stat-label">联系电话</span>
          <strong class="stat-value teacher-detail-phone">{{ currentTeacher?.mobile || "-" }}</strong>
        </article>
      </div>

      <div class="detail-note" v-if="currentTeacher?.remark">
        <strong>备注</strong>
        <p>{{ currentTeacher.remark }}</p>
      </div>
    </section>

    <div class="teacher-detail-grid">
      <section class="page-card page-card--table">
        <div class="page-header">
          <div>
            <h2>负责班级</h2>
            <p class="soft-text">当前班级和固定排课会集中显示在这里。</p>
          </div>
        </div>

        <div class="data-table-shell">
          <el-table :data="currentClasses" stripe>
            <el-table-column label="班级" min-width="200">
              <template #default="{ row }">
                <div class="table-primary">
                  <el-button link type="primary" @click="openClassDetail(row.id)">
                    {{ row.name }}
                  </el-button>
                  <small>{{ row.courseName }} · {{ row.weeklySchedule || "未设置固定排课" }}</small>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="校区" prop="campus" width="120" />
            <el-table-column label="人数" width="100">
              <template #default="{ row }">
                {{ row.studentCount }}/{{ row.capacity }}
              </template>
            </el-table-column>
            <el-table-column label="起止时间" min-width="180">
              <template #default="{ row }">
                {{ row.startDate || "未填" }} 至 {{ row.endDate || "未填" }}
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === '开班中' ? 'success' : row.status === '待满班' ? 'warning' : 'info'">
                  {{ row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openClassDetail(row.id)">查看班级</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div v-if="currentClasses.length === 0" class="soft-empty teacher-detail-empty">
          这个老师目前还没有关联班级。
        </div>
      </section>

      <section class="page-card">
        <div class="page-header">
          <div>
            <h2>近期排课</h2>
            <p class="soft-text">最近要上的课会集中显示在这里。</p>
          </div>
        </div>

        <div class="stack-list">
          <article
            v-for="schedule in currentSchedules"
            :key="schedule.id"
            class="stack-item stack-item--stretch"
          >
            <div>
              <strong>{{ schedule.className }}</strong>
              <small>
                {{ schedule.courseName }} · {{ schedule.lessonDate }} {{ schedule.lessonTime }}
              </small>
              <small>
                {{ schedule.campus || "未填写校区" }} / {{ schedule.classroom || "未填写教室" }}
              </small>
            </div>
            <div class="teacher-schedule-actions">
              <el-tag :type="scheduleStatusTone(schedule.attendanceStatus)">
                {{ schedule.attendanceStatus }}
              </el-tag>
              <el-button link type="primary" @click="openScheduleDetail(schedule.id)">
                查看课程
              </el-button>
            </div>
          </article>

          <div v-if="currentSchedules.length === 0" class="soft-empty">
            这个老师目前还没有最近课程安排。
          </div>
        </div>
      </section>
    </div>
  </div>
</template>
