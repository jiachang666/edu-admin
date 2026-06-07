<script setup lang="ts">
import { EditPen, Plus, View } from "@element-plus/icons-vue";
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import {
  createTeacher,
  fetchTeacherList,
  updateTeacher,
  type Teacher,
  type TeacherPayload,
} from "../../api/education";
import { useAuthStore } from "../../stores/auth";

const router = useRouter();
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

function openTeacherDetail(teacherId: number) {
  void router.push(`/teachers/${teacherId}`);
}

function teacherInitial(name: string) {
  const trimmedName = name.trim();
  if (trimmedName.length === 0) {
    return "师";
  }

  return trimmedName.slice(0, 1);
}

function teacherStatusTagType(status: string) {
  if (status === "在职") {
    return "success";
  }
  if (status === "排课中") {
    return "warning";
  }

  return "info";
}

onMounted(() => {
  void loadTeachers();
});
</script>

<template>
  <div class="page-stack">
    <section class="page-card page-card--table list-card">
      <div class="page-header">
        <div class="list-card__heading">
          <h2>老师列表</h2>
          <span class="list-card__count">共 {{ filteredTeachers.length }} 条</span>
        </div>
        <div class="page-actions">
          <el-button
            v-if="canManageTeachers"
            class="teacher-create-button"
            :icon="Plus"
            type="primary"
            @click="openCreateDialog"
          >
            新增老师
          </el-button>
        </div>
      </div>

      <div class="metric-strip metric-strip--compact list-card__metrics">
        <article class="metric-tile">
          <span>老师总数</span>
          <strong>{{ teachers.length }}</strong>
        </article>
        <article class="metric-tile">
          <span>在职人数</span>
          <strong>{{ activeCount }}</strong>
        </article>
        <article class="metric-tile">
          <span>全职老师</span>
          <strong>{{ fullTimeCount }}</strong>
        </article>
        <article class="metric-tile">
          <span>覆盖校区</span>
          <strong>{{ campusCount }}</strong>
        </article>
        <article class="metric-tile">
          <span>周课时合计</span>
          <strong>{{ totalWeeklyHours }}</strong>
        </article>
      </div>

      <div class="filter-bar list-card__filters">
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
        <el-table v-loading="loading" class="teacher-table" :data="filteredTeachers" stripe>
          <el-table-column label="老师" min-width="240">
            <template #default="{ row }">
              <div class="teacher-name-cell">
                <button
                  class="teacher-avatar-button"
                  type="button"
                  :aria-label="`查看${row.name}详情`"
                  @click="openTeacherDetail(row.id)"
                >
                  {{ teacherInitial(row.name) }}
                </button>
                <div class="teacher-name-copy">
                  <button class="teacher-name-link" type="button" @click="openTeacherDetail(row.id)">
                    {{ row.name }}
                  </button>
                  <span>{{ row.title || "未填写职级" }}</span>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="主教科目" width="140">
            <template #default="{ row }">
              <el-tag class="teacher-subject-tag" type="primary">
                {{ row.mainSubject || "未填写" }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="校区" prop="campus" width="128" />
          <el-table-column label="类型" prop="employmentType" width="112" />
          <el-table-column label="周课时" width="104">
            <template #default="{ row }">
              <span class="teacher-hours-pill">
                <strong>{{ row.weeklyHours }}</strong>
                <span>课时</span>
              </span>
            </template>
          </el-table-column>
          <el-table-column label="手机号" prop="mobile" width="140" />
          <el-table-column label="备注" min-width="220">
            <template #default="{ row }">
              <span class="muted-cell teacher-note">{{ row.remark || "暂无备注" }}</span>
            </template>
          </el-table-column>
          <el-table-column align="center" label="状态" width="104">
            <template #default="{ row }">
              <el-tag :type="teacherStatusTagType(row.status)">
                {{ row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column align="right" label="操作" :width="canManageTeachers ? 178 : 98" fixed="right">
            <template #default="{ row }">
              <div class="teacher-action-group">
                <el-button
                  class="table-action-button"
                  :icon="View"
                  plain
                  size="small"
                  @click="openTeacherDetail(row.id)"
                >
                  详情
                </el-button>
                <el-button
                  v-if="canManageTeachers"
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
