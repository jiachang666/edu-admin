<script setup lang="ts">
import { EditPen, Plus, View } from "@element-plus/icons-vue";
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import {
  createStudent,
  fetchStudentList,
  updateStudent,
  type Student,
  type StudentPayload,
} from "../../api/education";
import { useAuthStore } from "../../stores/auth";

const router = useRouter();
const authStore = useAuthStore();

const loading = ref(false);
const saving = ref(false);
const dialogVisible = ref(false);
const editingStudentId = ref<number | null>(null);
const students = ref<Student[]>([]);
const formRef = ref<FormInstance>();

const filters = reactive({
  keyword: "",
  status: "",
  campus: "",
  grade: "",
});

const form = reactive<StudentPayload>(defaultForm());

const statusOptions = ["在读", "待续费", "停课", "结课"];
const genderOptions = ["男", "女", "未填写"];

const rules: FormRules<StudentPayload> = {
  name: [{ required: true, message: "请输入学员姓名", trigger: "blur" }],
  gradeName: [{ required: true, message: "请输入学员年级", trigger: "blur" }],
  campus: [{ required: true, message: "请输入所属校区", trigger: "blur" }],
  status: [{ required: true, message: "请选择学员状态", trigger: "change" }],
  guardianName: [{ required: true, message: "请输入家长姓名", trigger: "blur" }],
  guardianMobile: [{ required: true, message: "请输入家长手机号", trigger: "blur" }],
};

const filteredStudents = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase();

  return students.value.filter((student) => {
    const matchesKeyword =
      keyword.length === 0 ||
      [
        student.name,
        student.grade,
        student.parentName,
        student.parentMobile,
        student.className,
        student.campus,
      ]
        .join(" ")
        .toLowerCase()
        .includes(keyword);
    const matchesStatus = filters.status.length === 0 || student.status === filters.status;
    const matchesCampus = filters.campus.length === 0 || student.campus === filters.campus;
    const matchesGrade = filters.grade.length === 0 || student.grade === filters.grade;

    return matchesKeyword && matchesStatus && matchesCampus && matchesGrade;
  });
});

const activeCount = computed(() => {
  return students.value.filter((student) => student.status === "在读").length;
});

const renewalCount = computed(() => {
  return students.value.filter((student) => student.status === "待续费").length;
});

const classCount = computed(() => {
  return new Set(
    students.value.filter((student) => student.classId > 0).map((student) => student.classId),
  ).size;
});

const remainingHoursTotal = computed(() => {
  return students.value.reduce((sum, student) => sum + student.remainingHours, 0);
});

const campusOptions = computed(() => {
  return Array.from(new Set(students.value.map((student) => student.campus).filter(Boolean)));
});

const gradeOptions = computed(() => {
  return Array.from(new Set(students.value.map((student) => student.grade).filter(Boolean)));
});

const dialogTitle = computed(() => {
  return editingStudentId.value ? "编辑学员" : "新增学员";
});

const canManageStudents = computed(() => authStore.hasPermission("students:manage"));

function defaultForm(): StudentPayload {
  return {
    name: "",
    gender: "未填写",
    schoolName: "",
    gradeName: "",
    campus: "",
    remainingHours: 0,
    status: "在读",
    remark: "",
    guardianName: "",
    guardianMobile: "",
    guardianRelation: "家长",
  };
}

function resetForm() {
  Object.assign(form, defaultForm());
  editingStudentId.value = null;
  formRef.value?.clearValidate();
}

function openCreateDialog() {
  resetForm();
  dialogVisible.value = true;
}

async function openEditDialog(student: Student) {
  editingStudentId.value = student.id;
  Object.assign(form, {
    name: student.name,
    gender: "未填写",
    schoolName: "",
    gradeName: student.grade,
    campus: student.campus,
    remainingHours: student.remainingHours,
    status: student.status,
    remark: "",
    guardianName: student.parentName,
    guardianMobile: student.parentMobile,
    guardianRelation: "家长",
  });

  dialogVisible.value = true;
  formRef.value?.clearValidate();
}

function closeDialog() {
  dialogVisible.value = false;
  resetForm();
}

function buildPayload(): StudentPayload {
  return {
    name: form.name.trim(),
    gender: form.gender.trim(),
    schoolName: form.schoolName.trim(),
    gradeName: form.gradeName.trim(),
    campus: form.campus.trim(),
    remainingHours: Math.max(0, Number(form.remainingHours) || 0),
    status: form.status,
    remark: form.remark.trim(),
    guardianName: form.guardianName.trim(),
    guardianMobile: form.guardianMobile.trim(),
    guardianRelation: form.guardianRelation.trim(),
  };
}

async function loadStudents() {
  loading.value = true;

  try {
    const result = await fetchStudentList();
    students.value = result.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("学员列表加载失败");
  } finally {
    loading.value = false;
  }
}

async function submitForm() {
  const formNode = formRef.value;
  if (!formNode) {
    return;
  }

  const valid = await formNode.validate().catch(() => false);
  if (!valid) {
    return;
  }

  saving.value = true;

  try {
    const payload = buildPayload();

    if (editingStudentId.value) {
      await updateStudent(editingStudentId.value, payload);
      ElMessage.success("学员资料已更新");
    } else {
      await createStudent(payload);
      ElMessage.success("学员已创建");
    }

    closeDialog();
    await loadStudents();
  } catch (error: any) {
    console.error(error);
    const message = error?.response?.data?.message ?? "学员资料保存失败";
    ElMessage.error(message);
  } finally {
    saving.value = false;
  }
}

function openStudentDetail(studentId: number) {
  void router.push(`/students/${studentId}`);
}

function handleReset() {
  filters.keyword = "";
  filters.status = "";
  filters.campus = "";
  filters.grade = "";
}

function studentInitial(name: string) {
  const trimmedName = name.trim();
  if (trimmedName.length === 0) {
    return "生";
  }

  return trimmedName.slice(0, 1);
}

function studentStatusTagType(status: string) {
  if (status === "在读") {
    return "success";
  }
  if (status === "待续费") {
    return "warning";
  }
  if (status === "停课") {
    return "danger";
  }

  return "info";
}

onMounted(() => {
  void loadStudents();
});
</script>

<template>
  <div class="page-stack">
    <section class="page-card page-card--table list-card">
      <div class="page-header">
        <div class="list-card__heading">
          <h2>学员列表</h2>
          <span class="list-card__count">共 {{ filteredStudents.length }} 条</span>
        </div>
        <div class="page-actions">
          <el-button
            v-if="canManageStudents"
            class="student-create-button"
            :icon="Plus"
            type="primary"
            @click="openCreateDialog"
          >
            新增学员
          </el-button>
        </div>
      </div>

      <div class="metric-strip metric-strip--compact list-card__metrics">
        <article class="metric-tile">
          <span>学员总数</span>
          <strong>{{ students.length }}</strong>
        </article>
        <article class="metric-tile">
          <span>在读学员</span>
          <strong>{{ activeCount }}</strong>
        </article>
        <article class="metric-tile">
          <span>待续费</span>
          <strong>{{ renewalCount }}</strong>
        </article>
        <article class="metric-tile">
          <span>已分班级</span>
          <strong>{{ classCount }}</strong>
        </article>
        <article class="metric-tile">
          <span>剩余课时</span>
          <strong>{{ remainingHoursTotal }}</strong>
        </article>
        <article class="metric-tile">
          <span>校区数量</span>
          <strong>{{ campusOptions.length }}</strong>
        </article>
      </div>

      <div class="filter-bar list-card__filters">
        <div class="toolbar-filters">
          <el-input
            v-model="filters.keyword"
            class="toolbar-field"
            clearable
            placeholder="搜索姓名、家长、电话、班级、校区"
          />
          <el-select
            v-model="filters.grade"
            class="toolbar-field"
            clearable
            placeholder="学员年级"
          >
            <el-option
              v-for="option in gradeOptions"
              :key="option"
              :label="option"
              :value="option"
            />
          </el-select>
          <el-select
            v-model="filters.status"
            class="toolbar-field"
            clearable
            placeholder="学员状态"
          >
            <el-option
              v-for="option in statusOptions"
              :key="option"
              :label="option"
              :value="option"
            />
          </el-select>
          <el-select
            v-model="filters.campus"
            class="toolbar-field"
            clearable
            placeholder="所属校区"
          >
            <el-option
              v-for="option in campusOptions"
              :key="option"
              :label="option"
              :value="option"
            />
          </el-select>
        </div>

        <div class="toolbar-actions">
          <el-button plain @click="handleReset">重置</el-button>
        </div>
      </div>

      <div class="data-table-shell">
        <el-table v-loading="loading" class="student-table" :data="filteredStudents" stripe>
          <el-table-column label="学员" min-width="240">
            <template #default="{ row }">
              <div class="student-name-cell">
                <button
                  class="student-avatar-button"
                  type="button"
                  :aria-label="`查看${row.name}详情`"
                  @click="openStudentDetail(row.id)"
                >
                  {{ studentInitial(row.name) }}
                </button>
                <div class="student-name-copy">
                  <button class="student-name-link" type="button" @click="openStudentDetail(row.id)">
                    {{ row.name }}
                  </button>
                  <span>{{ row.grade || "未填写年级" }} · {{ row.className || "暂未分班" }}</span>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="所属班级" min-width="160">
            <template #default="{ row }">
              <el-tag class="student-class-tag" :type="row.className ? 'primary' : 'info'">
                {{ row.className || "暂未分班" }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="家长" prop="parentName" width="120" />
          <el-table-column label="联系电话" width="148">
            <template #default="{ row }">
              <span class="table-nowrap">{{ row.parentMobile }}</span>
            </template>
          </el-table-column>
          <el-table-column label="剩余课时" width="104">
            <template #default="{ row }">
              <span class="student-hours-pill">
                <strong>{{ row.remainingHours }}</strong>
                <span>课时</span>
              </span>
            </template>
          </el-table-column>
          <el-table-column label="校区" prop="campus" width="120" />
          <el-table-column align="center" label="状态" width="104">
            <template #default="{ row }">
              <el-tag :type="studentStatusTagType(row.status)">
                {{ row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column align="right" label="操作" :width="canManageStudents ? 178 : 98" fixed="right">
            <template #default="{ row }">
              <div class="student-action-group">
                <el-button
                  class="table-action-button"
                  :icon="View"
                  plain
                  size="small"
                  @click="openStudentDetail(row.id)"
                >
                  详情
                </el-button>
                <el-button
                  v-if="canManageStudents"
                  class="table-action-button"
                  :icon="EditPen"
                  plain
                  size="small"
                  @click="openEditDialog(row)"
                >
                  编辑
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </section>

    <el-dialog
      :model-value="dialogVisible"
      :title="dialogTitle"
      width="820px"
      destroy-on-close
      @close="closeDialog"
      @update:model-value="dialogVisible = $event"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <div class="dialog-grid">
          <el-form-item label="学员姓名" prop="name">
            <el-input v-model="form.name" placeholder="例如：李一诺" />
          </el-form-item>
          <el-form-item label="性别">
            <el-select v-model="form.gender" placeholder="请选择性别">
              <el-option
                v-for="option in genderOptions"
                :key="option"
                :label="option"
                :value="option"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="所在学校">
            <el-input v-model="form.schoolName" placeholder="例如：实验小学" />
          </el-form-item>
          <el-form-item label="学员年级" prop="gradeName">
            <el-input v-model="form.gradeName" placeholder="例如：三年级" />
          </el-form-item>
          <el-form-item label="所属校区" prop="campus">
            <el-input v-model="form.campus" placeholder="例如：明发校区" />
          </el-form-item>
          <el-form-item label="学员状态" prop="status">
            <el-select v-model="form.status" placeholder="请选择学员状态">
              <el-option
                v-for="option in statusOptions"
                :key="option"
                :label="option"
                :value="option"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="剩余课时">
            <el-input-number v-model="form.remainingHours" class="full-width" :min="0" :max="300" />
          </el-form-item>
          <el-form-item label="家长关系">
            <el-input v-model="form.guardianRelation" placeholder="例如：母亲 / 父亲 / 监护人" />
          </el-form-item>
          <el-form-item label="家长姓名" prop="guardianName">
            <el-input v-model="form.guardianName" placeholder="例如：李女士" />
          </el-form-item>
          <el-form-item label="家长手机号" prop="guardianMobile">
            <el-input v-model="form.guardianMobile" placeholder="用于联系通知" />
          </el-form-item>
        </div>

        <el-form-item label="备注">
          <el-input
            v-model="form.remark"
            type="textarea"
            :rows="4"
            placeholder="可以记录学习特点、跟进提醒或特殊情况"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-actions">
          <el-button @click="closeDialog">取消</el-button>
          <el-button type="primary" :loading="saving" @click="submitForm">保存学员资料</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>
