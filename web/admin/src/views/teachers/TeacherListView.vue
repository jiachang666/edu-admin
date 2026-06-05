<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import {
  createTeacher,
  fetchTeacherList,
  updateTeacher,
  type Teacher,
  type TeacherPayload,
} from "../../api/education";
import { useAuthStore } from "../../stores/auth";

const authStore = useAuthStore();

const loading = ref(false);
const saving = ref(false);
const dialogVisible = ref(false);
const editingTeacherId = ref<number | null>(null);
const teachers = ref<Teacher[]>([]);
const formRef = ref<FormInstance>();

const filters = reactive({
  keyword: "",
  status: "",
  employmentType: "",
  campus: "",
});

const form = reactive<TeacherPayload>(defaultForm());

const statusOptions = ["在职", "排课中", "停用"];
const employmentTypeOptions = ["全职", "兼职", "合作老师"];

const rules: FormRules<TeacherPayload> = {
  name: [{ required: true, message: "请输入老师姓名", trigger: "blur" }],
  mainSubject: [{ required: true, message: "请输入主教科目", trigger: "blur" }],
  employmentType: [{ required: true, message: "请选择用工类型", trigger: "change" }],
  campus: [{ required: true, message: "请输入所属校区", trigger: "blur" }],
  status: [{ required: true, message: "请选择老师状态", trigger: "change" }],
};

const filteredTeachers = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase();

  return teachers.value.filter((teacher) => {
    const matchesKeyword =
      keyword.length === 0 ||
      [
        teacher.name,
        teacher.title,
        teacher.mainSubject,
        teacher.mobile,
        teacher.campus,
        teacher.remark,
      ]
        .join(" ")
        .toLowerCase()
        .includes(keyword);
    const matchesStatus = filters.status.length === 0 || teacher.status === filters.status;
    const matchesEmploymentType =
      filters.employmentType.length === 0 || teacher.employmentType === filters.employmentType;
    const matchesCampus = filters.campus.length === 0 || teacher.campus === filters.campus;

    return matchesKeyword && matchesStatus && matchesEmploymentType && matchesCampus;
  });
});

const activeCount = computed(() => {
  return teachers.value.filter((teacher) => teacher.status === "在职").length;
});

const fullTimeCount = computed(() => {
  return teachers.value.filter((teacher) => teacher.employmentType === "全职").length;
});

const campusCount = computed(() => {
  return new Set(teachers.value.map((teacher) => teacher.campus).filter(Boolean)).size;
});

const totalWeeklyHours = computed(() => {
  return teachers.value.reduce((sum, teacher) => sum + teacher.weeklyHours, 0);
});

const campusOptions = computed(() => {
  return Array.from(new Set(teachers.value.map((teacher) => teacher.campus).filter(Boolean)));
});

const dialogTitle = computed(() => {
  return editingTeacherId.value ? "编辑老师" : "新增老师";
});

const canManageTeachers = computed(() => authStore.hasPermission("teachers:manage"));

function defaultForm(): TeacherPayload {
  return {
    name: "",
    mobile: "",
    title: "",
    mainSubject: "",
    employmentType: "全职",
    weeklyHours: 12,
    campus: "",
    status: "在职",
    remark: "",
  };
}

function resetForm() {
  Object.assign(form, defaultForm());
  editingTeacherId.value = null;
  formRef.value?.clearValidate();
}

function openCreateDialog() {
  resetForm();
  dialogVisible.value = true;
}

function openEditDialog(teacher: Teacher) {
  editingTeacherId.value = teacher.id;
  Object.assign(form, {
    name: teacher.name,
    mobile: teacher.mobile,
    title: teacher.title,
    mainSubject: teacher.mainSubject,
    employmentType: teacher.employmentType,
    weeklyHours: teacher.weeklyHours,
    campus: teacher.campus,
    status: teacher.status,
    remark: teacher.remark,
  });
  dialogVisible.value = true;
  formRef.value?.clearValidate();
}

function closeDialog() {
  dialogVisible.value = false;
  resetForm();
}

function buildPayload(): TeacherPayload {
  return {
    name: form.name.trim(),
    mobile: form.mobile.trim(),
    title: form.title.trim(),
    mainSubject: form.mainSubject.trim(),
    employmentType: form.employmentType,
    weeklyHours: Math.max(0, Number(form.weeklyHours) || 0),
    campus: form.campus.trim(),
    status: form.status,
    remark: form.remark.trim(),
  };
}

async function loadTeachers() {
  loading.value = true;

  try {
    const result = await fetchTeacherList();
    teachers.value = result.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("老师列表加载失败");
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

    if (editingTeacherId.value) {
      await updateTeacher(editingTeacherId.value, payload);
      ElMessage.success("老师资料已更新");
    } else {
      await createTeacher(payload);
      ElMessage.success("老师已创建");
    }

    closeDialog();
    await loadTeachers();
  } catch (error: any) {
    console.error(error);
    const message = error?.response?.data?.message ?? "老师资料保存失败";
    ElMessage.error(message);
  } finally {
    saving.value = false;
  }
}

function handleReset() {
  filters.keyword = "";
  filters.status = "";
  filters.employmentType = "";
  filters.campus = "";
}

onMounted(() => {
  void loadTeachers();
});
</script>

<template>
  <div class="page-stack">
    <section class="page-hero">
      <div class="page-hero__copy">
        <span class="section-kicker">Faculty Ledger</span>
        <h2>把老师的科目、职级、校区和授课负载放进一套能直接维护的师资台账里。</h2>
        <p>
          现在这页已经不只是盘点视图了。教务和负责人可以直接新增老师、补齐资料、更新状态，把排课前最常用的师资底册先理顺。
        </p>
      </div>

      <div class="metric-strip">
        <article class="metric-tile">
          <span>老师总数</span>
          <strong>{{ teachers.length }}</strong>
          <small>当前系统里已经登记的全部师资</small>
        </article>
        <article class="metric-tile">
          <span>在职人数</span>
          <strong>{{ activeCount }}</strong>
          <small>当前可正常安排课程的老师</small>
        </article>
        <article class="metric-tile">
          <span>全职老师</span>
          <strong>{{ fullTimeCount }}</strong>
          <small>便于快速识别稳定排班主力</small>
        </article>
        <article class="metric-tile">
          <span>覆盖校区</span>
          <strong>{{ campusCount }}</strong>
          <small>看清老师分布是否均衡</small>
        </article>
        <article class="metric-tile">
          <span>周课时合计</span>
          <strong>{{ totalWeeklyHours }}</strong>
          <small>帮助判断当前整体授课负载</small>
        </article>
      </div>
    </section>

    <section class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>老师列表</h2>
          <p class="soft-text">支持筛选、新增和编辑，方便你把师资台账真正维护起来。</p>
        </div>
        <div class="page-actions">
          <div class="section-note">师资工作台</div>
          <el-button v-if="canManageTeachers" type="primary" @click="openCreateDialog">
            新增老师
          </el-button>
        </div>
      </div>

      <div class="page-toolbar">
        <div class="toolbar-filters">
          <el-input
            v-model="filters.keyword"
            class="toolbar-field"
            clearable
            placeholder="搜索姓名、科目、职级、手机、校区"
          />
          <el-select
            v-model="filters.employmentType"
            class="toolbar-field"
            clearable
            placeholder="用工类型"
          >
            <el-option
              v-for="option in employmentTypeOptions"
              :key="option"
              :label="option"
              :value="option"
            />
          </el-select>
          <el-select
            v-model="filters.status"
            class="toolbar-field"
            clearable
            placeholder="老师状态"
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
        <el-table v-loading="loading" :data="filteredTeachers" stripe>
          <el-table-column label="老师" min-width="170">
            <template #default="{ row }">
              <div class="table-primary">
                <strong>{{ row.name }}</strong>
                <small>{{ row.title || "未填写职级" }}</small>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="主教科目" prop="mainSubject" width="140" />
          <el-table-column label="校区" prop="campus" width="120" />
          <el-table-column label="类型" prop="employmentType" width="110" />
          <el-table-column label="周课时" prop="weeklyHours" width="100" />
          <el-table-column label="手机号" prop="mobile" width="140" />
          <el-table-column label="备注" min-width="220">
            <template #default="{ row }">
              <span class="muted-cell">{{ row.remark || "暂无备注" }}</span>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === '在职' ? 'success' : row.status === '排课中' ? 'warning' : 'info'">
                {{ row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column v-if="canManageTeachers" label="操作" width="100" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="openEditDialog(row)">编辑</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </section>

    <el-dialog
      :model-value="dialogVisible"
      :title="dialogTitle"
      width="760px"
      destroy-on-close
      @close="closeDialog"
      @update:model-value="dialogVisible = $event"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <div class="dialog-grid">
          <el-form-item label="老师姓名" prop="name">
            <el-input v-model="form.name" placeholder="例如：周老师" />
          </el-form-item>
          <el-form-item label="职级/标签">
            <el-input v-model="form.title" placeholder="例如：资深数学老师 / 教研负责人" />
          </el-form-item>
          <el-form-item label="主教科目" prop="mainSubject">
            <el-input v-model="form.mainSubject" placeholder="例如：数学思维" />
          </el-form-item>
          <el-form-item label="用工类型" prop="employmentType">
            <el-select v-model="form.employmentType" placeholder="请选择用工类型">
              <el-option
                v-for="option in employmentTypeOptions"
                :key="option"
                :label="option"
                :value="option"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="所属校区" prop="campus">
            <el-input v-model="form.campus" placeholder="例如：明发校区" />
          </el-form-item>
          <el-form-item label="老师状态" prop="status">
            <el-select v-model="form.status" placeholder="请选择老师状态">
              <el-option
                v-for="option in statusOptions"
                :key="option"
                :label="option"
                :value="option"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="手机号">
            <el-input v-model="form.mobile" placeholder="用于联系或登录关联" />
          </el-form-item>
          <el-form-item label="周课时">
            <el-input-number v-model="form.weeklyHours" class="full-width" :min="0" :max="60" />
          </el-form-item>
        </div>

        <el-form-item label="备注">
          <el-input
            v-model="form.remark"
            type="textarea"
            :rows="4"
            placeholder="可以记录擅长方向、带班偏好或近期排班说明"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-actions">
          <el-button @click="closeDialog">取消</el-button>
          <el-button type="primary" :loading="saving" @click="submitForm">保存老师资料</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>
