<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  addStudentsToClass,
  fetchClassList,
  fetchStudentDetail,
  removeStudentFromClass,
  type Homework,
  type SchoolClass,
  type StudentDetail,
  type StudentGuardian,
} from "../../api/education";

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const saving = ref(false);
const dialogVisible = ref(false);
const studentDetail = ref<StudentDetail | null>(null);
const allClasses = ref<SchoolClass[]>([]);
const formRef = ref<FormInstance>();

const form = reactive({
  classId: null as number | null,
});

const rules: FormRules<typeof form> = {
  classId: [{ required: true, message: "请选择要加入的班级", trigger: "change" }],
};

const studentId = computed(() => {
  const parsedValue = Number(route.params.id);
  if (!Number.isInteger(parsedValue) || parsedValue <= 0) {
    return null;
  }

  return parsedValue;
});

const currentStudent = computed(() => {
  return studentDetail.value?.student ?? null;
});

const currentClasses = computed(() => {
  return studentDetail.value?.classes ?? [];
});

const availableClasses = computed(() => {
  const currentClassIds = new Set(currentClasses.value.map((item) => item.id));
  return allClasses.value.filter((item) => !currentClassIds.has(item.id));
});

const primaryGuardian = computed(() => {
  return (
    studentDetail.value?.guardians.find((item) => item.isPrimary) ??
    studentDetail.value?.guardians[0] ??
    null
  );
});

const attendanceSummary = computed(() => {
  return (studentDetail.value?.recentAttendance ?? []).reduce(
    (summary, item) => {
      switch (item.status) {
        case "已到":
        case "补签":
          summary.present += 1;
          break;
        case "请假":
          summary.leave += 1;
          break;
        case "缺席":
          summary.absent += 1;
          break;
        default:
          summary.pending += 1;
          break;
      }
      return summary;
    },
    { present: 0, leave: 0, absent: 0, pending: 0 },
  );
});

function openJoinClassDialog() {
  form.classId = null;
  dialogVisible.value = true;
  formRef.value?.clearValidate();
}

function closeDialog() {
  dialogVisible.value = false;
  form.classId = null;
}

async function loadStudentDetail() {
  const currentStudentId = studentId.value;
  if (currentStudentId === null) {
    ElMessage.error("学员编号不正确");
    return;
  }

  loading.value = true;

  try {
    const [detailResult, classResult] = await Promise.all([
      fetchStudentDetail(currentStudentId),
      fetchClassList(),
    ]);
    studentDetail.value = detailResult;
    allClasses.value = classResult.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("学员详情加载失败");
  } finally {
    loading.value = false;
  }
}

async function handleJoinClass() {
  const formNode = formRef.value;
  const currentStudentId = studentId.value;
  if (!formNode || currentStudentId === null || form.classId === null) {
    return;
  }

  const valid = await formNode.validate().catch(() => false);
  if (!valid) {
    return;
  }

  saving.value = true;

  try {
    await addStudentsToClass(form.classId, {
      studentIds: [currentStudentId],
    });
    ElMessage.success("学员已加入班级");
    closeDialog();
    await loadStudentDetail();
  } catch (error) {
    console.error(error);
    ElMessage.error("加入班级失败");
  } finally {
    saving.value = false;
  }
}

async function handleRemoveFromClass(classId: number) {
  const currentStudentId = studentId.value;
  if (currentStudentId === null) {
    return;
  }

  try {
    await removeStudentFromClass(classId, currentStudentId);
    ElMessage.success("学员已移出班级");
    await loadStudentDetail();
  } catch (error) {
    console.error(error);
    ElMessage.error("移出班级失败");
  }
}

function openClassDetail(classId: number) {
  void router.push(`/classes/${classId}`);
}

function openAttendance(scheduleId: number) {
  void router.push(`/attendance?scheduleId=${scheduleId}`);
}

function openHomework(scheduleId: number) {
  void router.push(`/homeworks?scheduleId=${scheduleId}`);
}

function homeworkRoute(item: Homework) {
  return `/homeworks?scheduleId=${item.scheduleId}`;
}

function guardianLabel(guardian: StudentGuardian) {
  if (guardian.relation.length > 0) {
    return `${guardian.name} · ${guardian.relation}`;
  }

  return guardian.name;
}

onMounted(() => {
  void loadStudentDetail();
});
</script>

<template>
  <div v-loading="loading" class="page-stack">
    <section class="page-hero">
      <div class="page-hero__copy">
        <span class="section-kicker">Student Hub</span>
        <h2>{{ currentStudent?.name || "学员详情" }}</h2>
        <p>
          把一个学员的家长信息、班级归属、最近上课、签到和课后反馈集中在一起，方便老师和教务接力跟进。
        </p>
      </div>

      <div class="metric-strip">
        <article class="metric-tile">
          <span>当前班级</span>
          <strong>{{ currentClasses.length }}</strong>
          <small>这个学员当前仍在读的班级数量</small>
        </article>
        <article class="metric-tile">
          <span>最近课程</span>
          <strong>{{ studentDetail?.recentSchedules.length ?? 0 }}</strong>
          <small>最近几次和这个学员相关的上课安排</small>
        </article>
        <article class="metric-tile">
          <span>最近作业</span>
          <strong>{{ studentDetail?.recentHomeworks.length ?? 0 }}</strong>
          <small>这个学员所在班级的近期课后作业</small>
        </article>
        <article class="metric-tile">
          <span>剩余课时</span>
          <strong>{{ currentStudent?.remainingHours ?? 0 }}</strong>
          <small>方便教务和班主任随手判断续费节奏</small>
        </article>
      </div>
    </section>

    <div class="student-detail-grid">
      <section class="page-card">
        <div class="page-header">
          <div>
            <h2>基本信息</h2>
            <p class="soft-text">先确认学员是谁、来自哪个校区、当前学习状态怎样。</p>
          </div>
          <div class="section-note">学员档案</div>
        </div>

        <div class="detail-info-grid">
          <article class="detail-info-card">
            <span>姓名</span>
            <strong>{{ currentStudent?.name || "-" }}</strong>
          </article>
          <article class="detail-info-card">
            <span>年级</span>
            <strong>{{ currentStudent?.grade || "-" }}</strong>
          </article>
          <article class="detail-info-card">
            <span>所在学校</span>
            <strong>{{ currentStudent?.schoolName || "-" }}</strong>
          </article>
          <article class="detail-info-card">
            <span>状态</span>
            <strong>{{ currentStudent?.status || "-" }}</strong>
          </article>
          <article class="detail-info-card">
            <span>校区</span>
            <strong>{{ currentStudent?.campus || "-" }}</strong>
          </article>
          <article class="detail-info-card">
            <span>性别</span>
            <strong>{{ currentStudent?.gender || "未填写" }}</strong>
          </article>
        </div>

        <div class="detail-note" v-if="currentStudent?.remark">
          <strong>备注</strong>
          <p>{{ currentStudent.remark }}</p>
        </div>
      </section>

      <section class="page-card">
        <div class="page-header">
          <div>
            <h2>家长信息</h2>
            <p class="soft-text">把主要联系人和备用联系人放在一起，方便随时联系。</p>
          </div>
          <div class="section-note">联系人</div>
        </div>

        <div class="stack-list">
          <article
            v-for="guardian in studentDetail?.guardians ?? []"
            :key="guardian.id"
            class="stack-item stack-item--stretch"
          >
            <div>
              <strong>{{ guardianLabel(guardian) }}</strong>
              <small>{{ guardian.mobile }} · {{ guardian.isPrimary ? "主要联系人" : "备用联系人" }}</small>
            </div>
          </article>

          <div v-if="(studentDetail?.guardians.length ?? 0) === 0" class="soft-empty">
            这个学员目前还没有联系人资料。
          </div>
        </div>

        <div class="detail-callout" v-if="primaryGuardian">
          <strong>默认联系对象</strong>
          <p>{{ primaryGuardian.name }}，电话 {{ primaryGuardian.mobile }}</p>
        </div>
      </section>
    </div>

    <section class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>所属班级</h2>
          <p class="soft-text">可以直接在这里把学员加入班级，或者从某个班级移出。</p>
        </div>
        <div class="page-actions">
          <div class="section-note">班级关系</div>
          <el-button type="primary" @click="openJoinClassDialog">加入班级</el-button>
        </div>
      </div>

      <div class="data-table-shell">
        <el-table :data="currentClasses" stripe>
          <el-table-column label="班级" min-width="200">
            <template #default="{ row }">
              <div class="table-primary">
                <strong>{{ row.name }}</strong>
                <small>{{ row.courseName }} · {{ row.teacherName }}</small>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="校区" prop="campus" width="120" />
          <el-table-column label="人数" width="100">
            <template #default="{ row }">
              {{ row.studentCount }}/{{ row.capacity }}
            </template>
          </el-table-column>
          <el-table-column label="固定排课" prop="weeklySchedule" min-width="180" />
          <el-table-column label="状态" prop="status" width="100" />
          <el-table-column label="操作" width="160" fixed="right">
            <template #default="{ row }">
              <div class="table-link-group">
                <el-button link type="primary" @click="openClassDetail(row.id)">查看班级</el-button>
                <el-button link type="danger" @click="handleRemoveFromClass(row.id)">移出班级</el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </section>

    <div class="student-detail-grid">
      <section class="page-card">
        <div class="page-header">
          <div>
            <h3>最近签到</h3>
            <p class="soft-text">先看这个学员最近到课、请假和缺席的情况。</p>
          </div>
          <div class="section-note">到课记录</div>
        </div>

        <div class="stats-grid stats-grid--compact">
          <article class="stat-card" data-tone="green">
            <span class="stat-label">已到</span>
            <strong class="stat-value">{{ attendanceSummary.present }}</strong>
          </article>
          <article class="stat-card" data-tone="orange">
            <span class="stat-label">请假</span>
            <strong class="stat-value">{{ attendanceSummary.leave }}</strong>
          </article>
          <article class="stat-card" data-tone="red">
            <span class="stat-label">缺席</span>
            <strong class="stat-value">{{ attendanceSummary.absent }}</strong>
          </article>
          <article class="stat-card" data-tone="indigo">
            <span class="stat-label">待确认</span>
            <strong class="stat-value">{{ attendanceSummary.pending }}</strong>
          </article>
        </div>

        <div class="stack-list">
          <article
            v-for="item in studentDetail?.recentAttendance ?? []"
            :key="`${item.scheduleId}-${item.classId}`"
            class="stack-item stack-item--stretch"
          >
            <div>
              <strong>{{ item.lessonDate }} {{ item.lessonTime }} · {{ item.status }}</strong>
              <small>
                {{ item.className }} · {{ item.courseName }} · {{ item.teacherName }}
                <template v-if="item.remark"> · {{ item.remark }}</template>
              </small>
            </div>
            <el-button link type="primary" @click="openAttendance(item.scheduleId)">去签到页</el-button>
          </article>

          <div v-if="(studentDetail?.recentAttendance.length ?? 0) === 0" class="soft-empty">
            这个学员目前还没有签到记录。
          </div>
        </div>
      </section>

      <section class="page-card">
        <div class="page-header">
          <div>
            <h3>最近课程</h3>
            <p class="soft-text">把最近几次上课安排放在一起，方便快速回看。</p>
          </div>
          <div class="section-note">上课安排</div>
        </div>

        <div class="stack-list">
          <article
            v-for="item in studentDetail?.recentSchedules ?? []"
            :key="item.id"
            class="stack-item stack-item--stretch"
          >
            <div>
              <strong>{{ item.lessonDate }} {{ item.lessonTime }}</strong>
              <small>{{ item.className }} · {{ item.courseName }} · {{ item.teacherName }} · {{ item.attendanceStatus }}</small>
            </div>
            <div class="table-link-group">
              <el-button link type="primary" @click="openAttendance(item.id)">看签到</el-button>
              <el-button link type="primary" @click="openHomework(item.id)">看作业</el-button>
            </div>
          </article>

          <div v-if="(studentDetail?.recentSchedules.length ?? 0) === 0" class="soft-empty">
            这个学员目前还没有最近课程安排。
          </div>
        </div>
      </section>
    </div>

    <div class="student-detail-grid">
      <section class="page-card">
        <div class="page-header">
          <div>
            <h3>最近作业</h3>
            <p class="soft-text">最近布置过的作业可以在这里快速回看。</p>
          </div>
          <div class="section-note">课后作业</div>
        </div>

        <div class="stack-list">
          <article
            v-for="item in studentDetail?.recentHomeworks ?? []"
            :key="item.id"
            class="stack-item stack-item--stretch"
          >
            <div>
              <strong>{{ item.title }}</strong>
              <small>{{ item.lessonDate }} · {{ item.className }} · {{ item.teacherName }}</small>
            </div>
            <el-button link type="primary" @click="router.push(homeworkRoute(item))">查看作业</el-button>
          </article>

          <div v-if="(studentDetail?.recentHomeworks.length ?? 0) === 0" class="soft-empty">
            这个学员目前还没有作业记录。
          </div>
        </div>
      </section>

      <section class="page-card">
        <div class="page-header">
          <div>
            <h3>最近反馈</h3>
            <p class="soft-text">老师最近的课堂反馈和家长提示集中放在这里。</p>
          </div>
          <div class="section-note">课后反馈</div>
        </div>

        <div class="stack-list">
          <article
            v-for="item in studentDetail?.recentFeedbacks ?? []"
            :key="item.id"
            class="stack-item stack-item--stretch"
          >
            <div>
              <strong>{{ item.className }} · {{ item.lessonDate }}</strong>
              <small>{{ item.summary || "老师暂未填写课堂总结" }}</small>
            </div>
            <el-button link type="primary" @click="openHomework(item.scheduleId)">查看反馈</el-button>
          </article>

          <div v-if="(studentDetail?.recentFeedbacks.length ?? 0) === 0" class="soft-empty">
            这个学员目前还没有反馈记录。
          </div>
        </div>
      </section>
    </div>

    <el-dialog
      :model-value="dialogVisible"
      title="把学员加入班级"
      width="560px"
      destroy-on-close
      @close="closeDialog"
      @update:model-value="dialogVisible = $event"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="选择班级" prop="classId">
          <el-select v-model="form.classId" class="full-width" placeholder="请选择要加入的班级">
            <el-option
              v-for="item in availableClasses"
              :key="item.id"
              :label="`${item.name} · ${item.courseName} · ${item.teacherName}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-actions">
          <el-button @click="closeDialog">取消</el-button>
          <el-button type="primary" :loading="saving" @click="handleJoinClass">确认加入</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>
