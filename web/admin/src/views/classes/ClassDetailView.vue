<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  addStudentsToClass,
  fetchClassDetail,
  fetchStudentList,
  removeStudentFromClass,
  type AttendanceSession,
  type ClassDetail,
  type Homework,
  type Student,
} from "../../api/education";

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const saving = ref(false);
const dialogVisible = ref(false);
const classDetail = ref<ClassDetail | null>(null);
const allStudents = ref<Student[]>([]);
const formRef = ref<FormInstance>();

const form = reactive({
  studentIds: [] as number[],
});

const rules: FormRules<typeof form> = {
  studentIds: [{ required: true, message: "请至少选择一个学员", trigger: "change" }],
};

const classId = computed(() => {
  const parsedValue = Number(route.params.id);
  if (!Number.isInteger(parsedValue) || parsedValue <= 0) {
    return null;
  }

  return parsedValue;
});

const currentClass = computed(() => {
  return classDetail.value?.class ?? null;
});

const availableStudents = computed(() => {
  const currentStudentIds = new Set((classDetail.value?.students ?? []).map((item) => item.id));
  return allStudents.value.filter((item) => !currentStudentIds.has(item.id));
});

const attendanceSummary = computed(() => {
  return (classDetail.value?.recentAttendance ?? []).reduce(
    (summary, item) => {
      summary.present += item.presentCount;
      summary.leave += item.leaveCount;
      summary.absent += item.absentCount;
      summary.pending += item.pendingCount;
      return summary;
    },
    { present: 0, leave: 0, absent: 0, pending: 0 },
  );
});

function openAddStudentDialog() {
  form.studentIds = [];
  dialogVisible.value = true;
  formRef.value?.clearValidate();
}

function closeDialog() {
  dialogVisible.value = false;
  form.studentIds = [];
}

async function loadClassDetail() {
  const currentClassId = classId.value;
  if (currentClassId === null) {
    ElMessage.error("班级编号不正确");
    return;
  }

  loading.value = true;

  try {
    const [detailResult, studentResult] = await Promise.all([
      fetchClassDetail(currentClassId),
      fetchStudentList(),
    ]);
    classDetail.value = detailResult;
    allStudents.value = studentResult.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("班级详情加载失败");
  } finally {
    loading.value = false;
  }
}

async function handleAddStudents() {
  const formNode = formRef.value;
  if (!formNode || classId.value === null) {
    return;
  }

  const valid = await formNode.validate().catch(() => false);
  if (!valid) {
    return;
  }

  saving.value = true;

  try {
    await addStudentsToClass(classId.value, {
      studentIds: form.studentIds,
    });
    ElMessage.success("学员已加入班级");
    closeDialog();
    await loadClassDetail();
  } catch (error) {
    console.error(error);
    ElMessage.error("加入班级失败");
  } finally {
    saving.value = false;
  }
}

async function handleRemoveStudent(studentId: number) {
  if (classId.value === null) {
    return;
  }

  try {
    await removeStudentFromClass(classId.value, studentId);
    ElMessage.success("学员已移出班级");
    await loadClassDetail();
  } catch (error) {
    console.error(error);
    ElMessage.error("移出班级失败");
  }
}

function openAttendance(session: AttendanceSession) {
  void router.push(`/attendance?scheduleId=${session.id}`);
}

function openHomework(scheduleId: number) {
  void router.push(`/homeworks?scheduleId=${scheduleId}`);
}

function openNoticeCenter() {
  void router.push("/notices");
}

function openSchedules() {
  void router.push("/schedules");
}

function homeworkRoute(homework: Homework) {
  return `/homeworks?scheduleId=${homework.scheduleId}`;
}

onMounted(() => {
  void loadClassDetail();
});
</script>

<template>
  <div v-loading="loading" class="page-stack">
    <section class="page-card">
      <div class="page-header">
        <div>
          <h2>{{ currentClass?.name || "班级详情" }}</h2>
          <p class="soft-text">
            {{ currentClass?.courseName || "未填写课程" }} ·
            {{ currentClass?.teacherName || "未填写老师" }} ·
            {{ currentClass?.campus || "未填写校区" }}
          </p>
        </div>
        <div class="page-actions">
          <el-button plain @click="openSchedules">查看排课</el-button>
          <el-button plain @click="openNoticeCenter">查看通知</el-button>
        </div>
      </div>

      <div class="metric-strip metric-strip--compact">
        <article class="metric-tile">
          <span>当前人数</span>
          <strong>{{ currentClass?.studentCount ?? 0 }}</strong>
        </article>
        <article class="metric-tile">
          <span>近期课程</span>
          <strong>{{ classDetail?.upcomingSchedules.length ?? 0 }}</strong>
        </article>
        <article class="metric-tile">
          <span>最近作业</span>
          <strong>{{ classDetail?.recentHomeworks.length ?? 0 }}</strong>
        </article>
        <article class="metric-tile">
          <span>关联通知</span>
          <strong>{{ classDetail?.recentNotices.length ?? 0 }}</strong>
        </article>
      </div>

      <div class="stats-grid">
        <article class="stat-card">
          <span class="stat-label">课程</span>
          <strong class="stat-value">{{ currentClass?.courseName || "-" }}</strong>
        </article>
        <article class="stat-card">
          <span class="stat-label">主讲老师</span>
          <strong class="stat-value">{{ currentClass?.teacherName || "-" }}</strong>
        </article>
        <article class="stat-card">
          <span class="stat-label">校区</span>
          <strong class="stat-value">{{ currentClass?.campus || "-" }}</strong>
        </article>
        <article class="stat-card">
          <span class="stat-label">固定排课</span>
          <strong class="stat-value">{{ currentClass?.weeklySchedule || "-" }}</strong>
        </article>
      </div>
    </section>

    <section class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>学员名单</h2>
          <p class="soft-text">可以直接给班级加人，或者把不再上课的学员移出。</p>
        </div>
        <div class="page-actions">
          <el-button type="primary" @click="openAddStudentDialog">添加学员</el-button>
        </div>
      </div>

      <div class="data-table-shell">
        <el-table :data="classDetail?.students ?? []" stripe>
          <el-table-column label="学员" prop="name" width="120" />
          <el-table-column label="年级" prop="grade" width="100" />
          <el-table-column label="家长" prop="parentName" width="120" />
          <el-table-column label="家长手机" prop="parentMobile" width="150" />
          <el-table-column label="剩余课时" prop="remainingHours" width="100" />
          <el-table-column label="状态" prop="status" width="100" />
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="{ row }">
              <el-button link type="danger" @click="handleRemoveStudent(row.id)">移出班级</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </section>

    <section class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>近期安排</h2>
          <p class="soft-text">接下来最需要处理的课程安排会集中显示在这里。</p>
        </div>
      </div>

      <div class="data-table-shell">
        <el-table :data="classDetail?.upcomingSchedules ?? []" stripe>
          <el-table-column label="日期" prop="lessonDate" width="120" />
          <el-table-column label="时间" prop="lessonTime" width="120" />
          <el-table-column label="课程" prop="courseName" width="140" />
          <el-table-column label="老师" prop="teacherName" width="120" />
          <el-table-column label="教室" prop="classroom" width="120" />
          <el-table-column label="状态" prop="attendanceStatus" width="120" />
          <el-table-column label="操作" min-width="180" fixed="right">
            <template #default="{ row }">
              <div class="table-link-group">
                <el-button link type="primary" @click="openAttendance(row)">
                  {{ row.attendanceStatus === "待签到" ? "去签到" : "看签到" }}
                </el-button>
                <el-button link type="primary" @click="openHomework(row.id)">作业反馈</el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </section>

    <div class="class-detail-grid">
      <section class="page-card">
        <div class="page-header">
          <div>
            <h3>最近签到</h3>
            <p class="soft-text">近期到课、请假和缺席的整体情况会集中显示在这里。</p>
          </div>
        </div>

        <div class="stats-grid stats-grid--compact">
          <article class="stat-card">
            <span class="stat-label">已到</span>
            <strong class="stat-value">{{ attendanceSummary.present }}</strong>
          </article>
          <article class="stat-card">
            <span class="stat-label">请假</span>
            <strong class="stat-value">{{ attendanceSummary.leave }}</strong>
          </article>
          <article class="stat-card">
            <span class="stat-label">缺席</span>
            <strong class="stat-value">{{ attendanceSummary.absent }}</strong>
          </article>
          <article class="stat-card">
            <span class="stat-label">待确认</span>
            <strong class="stat-value">{{ attendanceSummary.pending }}</strong>
          </article>
        </div>

        <div class="stack-list">
          <article
            v-for="item in classDetail?.recentAttendance ?? []"
            :key="item.id"
            class="stack-item"
          >
            <div>
              <strong>{{ item.lessonDate }} {{ item.lessonTime }}</strong>
              <small>{{ item.courseName }} · {{ item.attendanceStatus }}</small>
            </div>
            <el-button link type="primary" @click="openAttendance(item)">去处理</el-button>
          </article>
        </div>
      </section>

      <section class="page-card">
        <div class="page-header">
          <div>
            <h3>最近作业反馈</h3>
            <p class="soft-text">老师上完课后留下的作业和反馈可以在这里快速回看。</p>
          </div>
        </div>

        <div class="stack-list">
          <article
            v-for="item in classDetail?.recentHomeworks ?? []"
            :key="item.id"
            class="stack-item stack-item--stretch"
          >
            <div>
              <strong>{{ item.title }}</strong>
              <small>{{ item.lessonDate }} · {{ item.teacherName }} · {{ item.status === "published" ? "已发布" : "草稿" }}</small>
            </div>
            <el-button link type="primary" @click="router.push(homeworkRoute(item))">查看作业</el-button>
          </article>

          <div v-if="(classDetail?.recentHomeworks.length ?? 0) === 0" class="soft-empty">
            这个班目前还没有作业记录。
          </div>
        </div>
      </section>
    </div>

    <section class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>关联通知</h2>
          <p class="soft-text">和这个班相关的通知会集中显示在这里。</p>
        </div>
        <div class="page-actions">
          <el-button plain @click="openNoticeCenter">进入通知中心</el-button>
        </div>
      </div>

      <div class="stack-list stack-list--spacious">
        <article
          v-for="notice in classDetail?.recentNotices ?? []"
          :key="notice.id"
          class="stack-item stack-item--stretch"
        >
          <div>
            <strong>{{ notice.title }}</strong>
            <small>{{ notice.category }} · {{ notice.status }} · {{ notice.publishAt || "未发送" }}</small>
          </div>
          <el-button link type="primary" @click="openNoticeCenter">去通知页</el-button>
        </article>

        <div v-if="(classDetail?.recentNotices.length ?? 0) === 0" class="soft-empty">
          这个班目前还没有关联通知。
        </div>
      </div>
    </section>

    <el-dialog
      :model-value="dialogVisible"
      title="添加学员到班级"
      width="560px"
      destroy-on-close
      @close="closeDialog"
      @update:model-value="dialogVisible = $event"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="选择学员" prop="studentIds">
          <el-select
            v-model="form.studentIds"
            multiple
            filterable
            placeholder="请选择要加入班级的学员"
          >
            <el-option
              v-for="student in availableStudents"
              :key="student.id"
              :label="`${student.name} · ${student.grade} · ${student.parentName}`"
              :value="student.id"
            />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-actions">
          <el-button @click="closeDialog">取消</el-button>
          <el-button type="primary" :loading="saving" @click="handleAddStudents">确认添加</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>
