<script setup lang="ts">
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

onMounted(() => {
  void loadStudents();
});
</script>

<template>
  <div class="page-stack">
    <section class="page-hero">
      <div class="page-hero__copy">
        <span class="section-kicker">Student Ledger</span>
        <h2>把学员、家长、校区和课时余额放在同一套能直接维护的学员台账里。</h2>
        <p>
          这页现在既能快速盘点，也能直接新增和编辑。先把学员底册和首位联系人整理清楚，后面的分班、续费和跟进会顺很多。
        </p>
      </div>

      <div class="metric-strip">
        <article class="metric-tile">
          <span>学员总数</span>
          <strong>{{ students.length }}</strong>
          <small>当前系统里已经登记的全部学员</small>
        </article>
        <article class="metric-tile">
          <span>在读学员</span>
          <strong>{{ activeCount }}</strong>
          <small>当前重点跟进的稳定在读群体</small>
        </article>
        <article class="metric-tile">
          <span>待续费</span>
          <strong>{{ renewalCount }}</strong>
          <small>方便班主任和教务优先跟进</small>
        </article>
        <article class="metric-tile">
          <span>已分班级</span>
          <strong>{{ classCount }}</strong>
          <small>当前学员已进入的班级数量</small>
        </article>
        <article class="metric-tile">
          <span>剩余课时合计</span>
          <strong>{{ remainingHoursTotal }}</strong>
          <small>帮助判断整体课时消耗节奏</small>
        </article>
      </div>
    </section>

    <section class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>学员列表</h2>
          <p class="soft-text">支持直接新增和编辑，也能继续进入详情页维护家长和班级关系。</p>
        </div>
        <div class="page-actions">
          <div class="section-note">学员工作台</div>
          <el-button v-if="canManageStudents" type="primary" @click="openCreateDialog">
            新增学员
          </el-button>
        </div>
      </div>

      <div class="page-toolbar">
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
        <el-table v-loading="loading" :data="filteredStudents" stripe>
          <el-table-column label="姓名" min-width="150">
            <template #default="{ row }">
              <div class="table-primary">
                <el-button link type="primary" @click="openStudentDetail(row.id)">
                  {{ row.name }}
                </el-button>
                <small>{{ row.grade }}</small>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="所属班级" prop="className" min-width="180" />
          <el-table-column label="家长" prop="parentName" width="120" />
          <el-table-column label="联系电话" prop="parentMobile" width="140" />
          <el-table-column label="剩余课时" prop="remainingHours" width="100" />
          <el-table-column label="校区" prop="campus" width="120" />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === '在读' ? 'success' : row.status === '待续费' ? 'warning' : 'info'">
                {{ row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" :width="canManageStudents ? 160 : 120" fixed="right">
            <template #default="{ row }">
              <div class="table-link-group">
                <el-button link type="primary" @click="openStudentDetail(row.id)">进入详情</el-button>
                <el-button
                  v-if="canManageStudents"
                  link
                  type="primary"
                  @click="openEditDialog(row)"
                >
                  编辑资料
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
