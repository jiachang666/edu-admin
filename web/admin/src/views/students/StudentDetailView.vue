<script setup lang="ts">
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  addStudentsToClass,
  createStudentGuardian,
  deleteStudentGuardian,
  fetchClassList,
  fetchStudentDetail,
  removeStudentFromClass,
  updateStudentGuardian,
  type Homework,
  type SchoolClass,
  type StudentDetail,
  type StudentGuardian,
  type StudentGuardianPayload,
} from "../../api/education";
import { useAuthStore } from "../../stores/auth";

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();

const loading = ref(false);
const classSaving = ref(false);
const guardianSaving = ref(false);
const classDialogVisible = ref(false);
const guardianDialogVisible = ref(false);
const editingGuardianId = ref<number | null>(null);
const studentDetail = ref<StudentDetail | null>(null);
const allClasses = ref<SchoolClass[]>([]);
const classFormRef = ref<FormInstance>();
const guardianFormRef = ref<FormInstance>();

const classForm = reactive({
  classId: null as number | null,
});

const guardianForm = reactive<StudentGuardianPayload>(defaultGuardianForm());

const classRules: FormRules<typeof classForm> = {
  classId: [{ required: true, message: "请选择要加入的班级", trigger: "change" }],
};

const guardianRules: FormRules<StudentGuardianPayload> = {
  name: [{ required: true, message: "请输入家长姓名", trigger: "blur" }],
  mobile: [{ required: true, message: "请输入家长手机号", trigger: "blur" }],
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

const canManageStudents = computed(() => authStore.hasPermission("students:manage"));
const canManageClasses = computed(() => authStore.hasPermission("classes:manage"));

const guardianDialogTitle = computed(() => {
  return editingGuardianId.value ? "编辑家长" : "新增家长";
});

function defaultGuardianForm(): StudentGuardianPayload {
  return {
    name: "",
    relation: "",
    mobile: "",
    isPrimary: false,
  };
}

function resetGuardianForm() {
  Object.assign(guardianForm, defaultGuardianForm());
  editingGuardianId.value = null;
  guardianFormRef.value?.clearValidate();
}

function openJoinClassDialog() {
  classForm.classId = null;
  classDialogVisible.value = true;
  classFormRef.value?.clearValidate();
}

function closeClassDialog() {
  classDialogVisible.value = false;
  classForm.classId = null;
}

function openGuardianCreateDialog() {
  resetGuardianForm();
  guardianForm.isPrimary = (studentDetail.value?.guardians.length ?? 0) === 0;
  guardianDialogVisible.value = true;
}

function openGuardianEditDialog(guardian: StudentGuardian) {
  editingGuardianId.value = guardian.id;
  Object.assign(guardianForm, {
    name: guardian.name,
    relation: guardian.relation,
    mobile: guardian.mobile,
    isPrimary: guardian.isPrimary,
  });
  guardianDialogVisible.value = true;
  guardianFormRef.value?.clearValidate();
}

function closeGuardianDialog() {
  guardianDialogVisible.value = false;
  resetGuardianForm();
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
  const formNode = classFormRef.value;
  const currentStudentId = studentId.value;
  if (!formNode || currentStudentId === null || classForm.classId === null) {
    return;
  }

  const valid = await formNode.validate().catch(() => false);
  if (!valid) {
    return;
  }

  classSaving.value = true;

  try {
    await addStudentsToClass(classForm.classId, {
      studentIds: [currentStudentId],
    });
    ElMessage.success("学员已加入班级");
    closeClassDialog();
    await loadStudentDetail();
  } catch (error: any) {
    console.error(error);
    const message = error?.response?.data?.message ?? "加入班级失败";
    ElMessage.error(message);
  } finally {
    classSaving.value = false;
  }
}

async function handleRemoveFromClass(classId: number) {
  const currentStudentId = studentId.value;
  if (currentStudentId === null) {
    return;
  }

  try {
    await ElMessageBox.confirm("移出后这个学员将不再属于该班级，确定继续吗？", "移出班级", {
      type: "warning",
      confirmButtonText: "确认移出",
      cancelButtonText: "取消",
    });
  } catch {
    return;
  }

  try {
    await removeStudentFromClass(classId, currentStudentId);
    ElMessage.success("学员已移出班级");
    await loadStudentDetail();
  } catch (error: any) {
    console.error(error);
    const message = error?.response?.data?.message ?? "移出班级失败";
    ElMessage.error(message);
  }
}

function buildGuardianPayload(): StudentGuardianPayload {
  return {
    name: guardianForm.name.trim(),
    relation: guardianForm.relation.trim(),
    mobile: guardianForm.mobile.trim(),
    isPrimary: guardianForm.isPrimary,
  };
}

async function submitGuardianForm() {
  const formNode = guardianFormRef.value;
  const currentStudentId = studentId.value;
  if (!formNode || currentStudentId === null) {
    return;
  }

  const valid = await formNode.validate().catch(() => false);
  if (!valid) {
    return;
  }

  guardianSaving.value = true;

  try {
    const payload = buildGuardianPayload();

    if (editingGuardianId.value) {
      await updateStudentGuardian(currentStudentId, editingGuardianId.value, payload);
      ElMessage.success("家长资料已更新");
    } else {
      await createStudentGuardian(currentStudentId, payload);
      ElMessage.success("家长已新增");
    }

    closeGuardianDialog();
    await loadStudentDetail();
  } catch (error: any) {
    console.error(error);
    const message = error?.response?.data?.message ?? "家长资料保存失败";
    ElMessage.error(message);
  } finally {
    guardianSaving.value = false;
  }
}

async function handleDeleteGuardian(guardian: StudentGuardian) {
  const currentStudentId = studentId.value;
  if (currentStudentId === null) {
    return;
  }

  try {
    await ElMessageBox.confirm(
      `删除后将无法直接联系 ${guardian.name}，确定继续吗？`,
      "删除家长",
      {
        type: "warning",
        confirmButtonText: "确认删除",
        cancelButtonText: "取消",
      },
    );
  } catch {
    return;
  }

  try {
    await deleteStudentGuardian(currentStudentId, guardian.id);
    ElMessage.success("家长已删除");
    await loadStudentDetail();
  } catch (error: any) {
    console.error(error);
    const message = error?.response?.data?.message ?? "删除家长失败";
    ElMessage.error(message);
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
    <div class="student-detail-grid">
      <section class="page-card">
        <div class="page-header">
          <div>
            <h2>{{ currentStudent?.name || "学员详情" }}</h2>
            <p class="soft-text">
              {{ currentStudent?.grade || "未填写年级" }} ·
              {{ currentStudent?.campus || "未填写校区" }} ·
              {{ currentStudent?.status || "未填写状态" }}
            </p>
          </div>
        </div>

        <div class="metric-strip metric-strip--compact">
          <article class="metric-tile">
            <span>当前班级</span>
            <strong>{{ currentClasses.length }}</strong>
          </article>
          <article class="metric-tile">
            <span>最近课程</span>
            <strong>{{ studentDetail?.recentSchedules.length ?? 0 }}</strong>
          </article>
          <article class="metric-tile">
            <span>最近作业</span>
            <strong>{{ studentDetail?.recentHomeworks.length ?? 0 }}</strong>
          </article>
          <article class="metric-tile">
            <span>剩余课时</span>
            <strong>{{ currentStudent?.remainingHours ?? 0 }}</strong>
          </article>
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
            <p class="soft-text">主要联系人和备用联系人会集中显示在这里。</p>
          </div>
          <div class="page-actions">
            <el-button v-if="canManageStudents" type="primary" @click="openGuardianCreateDialog">
              新增家长
            </el-button>
          </div>
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
            <div v-if="canManageStudents" class="table-link-group">
              <el-button link type="primary" @click="openGuardianEditDialog(guardian)">编辑</el-button>
              <el-button link type="danger" @click="handleDeleteGuardian(guardian)">删除</el-button>
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
          <p class="soft-text">可以直接把学员加入班级，或者从班级移出。</p>
        </div>
        <div class="page-actions">
          <el-button v-if="canManageClasses" type="primary" @click="openJoinClassDialog">
            加入班级
          </el-button>
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
          <el-table-column label="操作" :width="canManageClasses ? 160 : 100" fixed="right">
            <template #default="{ row }">
              <div class="table-link-group">
                <el-button link type="primary" @click="openClassDetail(row.id)">查看班级</el-button>
                <el-button
                  v-if="canManageClasses"
                  link
                  type="danger"
                  @click="handleRemoveFromClass(row.id)"
                >
                  移出班级
                </el-button>
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
            <p class="soft-text">最近到课、请假和缺席会集中显示在这里。</p>
          </div>
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
            <p class="soft-text">最近几次上课安排会集中显示在这里。</p>
          </div>
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
            <p class="soft-text">最近布置过的作业可以直接回看。</p>
          </div>
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
            <p class="soft-text">老师最近的课堂反馈和家长提示会集中显示在这里。</p>
          </div>
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
      :model-value="classDialogVisible"
      title="把学员加入班级"
      width="560px"
      destroy-on-close
      @close="closeClassDialog"
      @update:model-value="classDialogVisible = $event"
    >
      <el-form ref="classFormRef" :model="classForm" :rules="classRules" label-position="top">
        <el-form-item label="选择班级" prop="classId">
          <el-select v-model="classForm.classId" class="full-width" placeholder="请选择要加入的班级">
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
          <el-button @click="closeClassDialog">取消</el-button>
          <el-button type="primary" :loading="classSaving" @click="handleJoinClass">确认加入</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog
      :model-value="guardianDialogVisible"
      :title="guardianDialogTitle"
      width="620px"
      destroy-on-close
      @close="closeGuardianDialog"
      @update:model-value="guardianDialogVisible = $event"
    >
      <el-form ref="guardianFormRef" :model="guardianForm" :rules="guardianRules" label-position="top">
        <div class="dialog-grid">
          <el-form-item label="家长姓名" prop="name">
            <el-input v-model="guardianForm.name" placeholder="例如：李女士" />
          </el-form-item>
          <el-form-item label="关系">
            <el-input v-model="guardianForm.relation" placeholder="例如：母亲 / 父亲 / 监护人" />
          </el-form-item>
          <el-form-item label="手机号" prop="mobile">
            <el-input v-model="guardianForm.mobile" placeholder="用于联系和通知" />
          </el-form-item>
          <el-form-item label="联系优先级">
            <el-select v-model="guardianForm.isPrimary" placeholder="请选择">
              <el-option :value="true" label="设为主要联系人" />
              <el-option :value="false" label="设为备用联系人" />
            </el-select>
          </el-form-item>
        </div>
      </el-form>

      <template #footer>
        <div class="dialog-actions">
          <el-button @click="closeGuardianDialog">取消</el-button>
          <el-button type="primary" :loading="guardianSaving" @click="submitGuardianForm">
            保存家长资料
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>
