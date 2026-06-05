<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import {
  createClass,
  fetchClassList,
  fetchCourseOptions,
  fetchTeacherOptions,
  updateClass,
  type SchoolClass,
  type SchoolClassPayload,
  type SelectOption,
} from "../../api/education";
import { useAuthStore } from "../../stores/auth";

const router = useRouter();
const authStore = useAuthStore();

const loading = ref(false);
const saving = ref(false);
const dialogVisible = ref(false);
const editingClassId = ref<number | null>(null);
const classes = ref<SchoolClass[]>([]);
const teacherOptions = ref<SelectOption[]>([]);
const courseOptions = ref<SelectOption[]>([]);
const formRef = ref<FormInstance>();

const filters = reactive({
  keyword: "",
  status: "",
  teacher: "",
  campus: "",
});

const form = reactive<SchoolClassPayload>(defaultForm());

const statusOptions = ["开班中", "待满班", "已结课", "已停班"];

const rules: FormRules<SchoolClassPayload> = {
  name: [{ required: true, message: "请输入班级名称", trigger: "blur" }],
  courseId: [{ required: true, message: "请选择课程", trigger: "change" }],
  teacherId: [{ required: true, message: "请选择主讲老师", trigger: "change" }],
  campus: [{ required: true, message: "请输入所属校区", trigger: "blur" }],
  capacity: [{ required: true, message: "请输入班级容量", trigger: "blur" }],
  weeklySchedule: [{ required: true, message: "请输入固定排课", trigger: "blur" }],
  status: [{ required: true, message: "请选择班级状态", trigger: "change" }],
};

const filteredClasses = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase();

  return classes.value.filter((item) => {
    const matchesKeyword =
      keyword.length === 0 ||
      [
        item.name,
        item.courseName,
        item.teacherName,
        item.campus,
        item.weeklySchedule,
        item.remark,
      ]
        .join(" ")
        .toLowerCase()
        .includes(keyword);
    const matchesStatus = filters.status.length === 0 || item.status === filters.status;
    const matchesTeacher = filters.teacher.length === 0 || item.teacherName === filters.teacher;
    const matchesCampus = filters.campus.length === 0 || item.campus === filters.campus;

    return matchesKeyword && matchesStatus && matchesTeacher && matchesCampus;
  });
});

const runningCount = computed(() => {
  return classes.value.filter((item) => item.status === "开班中").length;
});

const remainingSeats = computed(() => {
  return classes.value.reduce((total, item) => total + Math.max(item.capacity - item.studentCount, 0), 0);
});

const campusCount = computed(() => {
  return new Set(classes.value.map((item) => item.campus).filter(Boolean)).size;
});

const teacherNameOptions = computed(() => {
  return Array.from(new Set(classes.value.map((item) => item.teacherName).filter(Boolean)));
});

const campusOptions = computed(() => {
  return Array.from(new Set(classes.value.map((item) => item.campus).filter(Boolean)));
});

const dialogTitle = computed(() => {
  return editingClassId.value ? "编辑班级" : "新增班级";
});

const canManageClasses = computed(() => authStore.hasPermission("classes:manage"));

function defaultForm(): SchoolClassPayload {
  return {
    name: "",
    courseId: 0,
    teacherId: 0,
    campus: "",
    capacity: 16,
    weeklySchedule: "",
    startDate: "",
    endDate: "",
    status: "待满班",
    remark: "",
  };
}

function resetForm() {
  Object.assign(form, defaultForm());
  editingClassId.value = null;
  formRef.value?.clearValidate();
}

function openCreateDialog() {
  resetForm();
  dialogVisible.value = true;
}

function openEditDialog(item: SchoolClass) {
  editingClassId.value = item.id;
  Object.assign(form, {
    name: item.name,
    courseId: item.courseId,
    teacherId: item.teacherId,
    campus: item.campus,
    capacity: item.capacity,
    weeklySchedule: item.weeklySchedule,
    startDate: item.startDate,
    endDate: item.endDate,
    status: item.status,
    remark: item.remark,
  });
  dialogVisible.value = true;
  formRef.value?.clearValidate();
}

function closeDialog() {
  dialogVisible.value = false;
  resetForm();
}

function buildPayload(): SchoolClassPayload {
  return {
    name: form.name.trim(),
    courseId: Number(form.courseId) || 0,
    teacherId: Number(form.teacherId) || 0,
    campus: form.campus.trim(),
    capacity: Math.max(1, Number(form.capacity) || 1),
    weeklySchedule: form.weeklySchedule.trim(),
    startDate: form.startDate.trim(),
    endDate: form.endDate.trim(),
    status: form.status,
    remark: form.remark.trim(),
  };
}

async function loadPageData() {
  loading.value = true;

  try {
    const [classResult, teacherResult, courseResult] = await Promise.all([
      fetchClassList(),
      fetchTeacherOptions(),
      fetchCourseOptions(),
    ]);
    classes.value = classResult.list;
    teacherOptions.value = teacherResult;
    courseOptions.value = courseResult;
  } catch (error) {
    console.error(error);
    ElMessage.error("班级工作台加载失败");
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

    if (editingClassId.value) {
      await updateClass(editingClassId.value, payload);
      ElMessage.success("班级资料已更新");
    } else {
      await createClass(payload);
      ElMessage.success("班级已创建");
    }

    closeDialog();
    await loadPageData();
  } catch (error: any) {
    console.error(error);
    const message = error?.response?.data?.message ?? "班级资料保存失败";
    ElMessage.error(message);
  } finally {
    saving.value = false;
  }
}

function handleResetFilters() {
  filters.keyword = "";
  filters.status = "";
  filters.teacher = "";
  filters.campus = "";
}

function openDetail(classId: number) {
  void router.push(`/classes/${classId}`);
}

function classStatusTagType(status: string) {
  switch (status) {
    case "开班中":
      return "success";
    case "待满班":
      return "warning";
    case "已结课":
      return "info";
    case "已停班":
      return "danger";
    default:
      return "info";
  }
}

onMounted(() => {
  void loadPageData();
});
</script>

<template>
  <div class="page-stack">
    <section class="page-hero">
      <div class="page-hero__copy">
        <span class="section-kicker">Classroom Matrix</span>
        <h2>先把班级、课程、老师和开班节奏维护清楚，再进入班级详情继续处理日常动作。</h2>
        <p>
          这页现在已经是班级工作台了。可以直接新建班级、调整老师和课程搭配，也能继续进入班级详情处理加人、排课、签到和通知。
        </p>
      </div>

      <div class="metric-strip">
        <article class="metric-tile">
          <span>班级总数</span>
          <strong>{{ classes.length }}</strong>
          <small>当前已经建立的全部班级</small>
        </article>
        <article class="metric-tile">
          <span>开班中</span>
          <strong>{{ runningCount }}</strong>
          <small>当前正在正常运转的班级</small>
        </article>
        <article class="metric-tile">
          <span>剩余名额</span>
          <strong>{{ remainingSeats }}</strong>
          <small>按容量减去当前人数得到</small>
        </article>
        <article class="metric-tile">
          <span>覆盖校区</span>
          <strong>{{ campusCount }}</strong>
          <small>帮助看清班级分布情况</small>
        </article>
      </div>
    </section>

    <section class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>班级列表</h2>
          <p class="soft-text">支持直接维护班级底册，也能继续进入详情页处理学员和排课。</p>
        </div>
        <div class="page-actions">
          <div class="section-note">班级工作台</div>
          <el-button v-if="canManageClasses" type="primary" @click="openCreateDialog">
            新增班级
          </el-button>
        </div>
      </div>

      <div class="page-toolbar">
        <div class="toolbar-filters">
          <el-input
            v-model="filters.keyword"
            class="toolbar-field"
            clearable
            placeholder="搜索班级、课程、老师、校区或备注"
          />
          <el-select
            v-model="filters.teacher"
            class="toolbar-field"
            clearable
            placeholder="主讲老师"
          >
            <el-option
              v-for="teacher in teacherNameOptions"
              :key="teacher"
              :label="teacher"
              :value="teacher"
            />
          </el-select>
          <el-select
            v-model="filters.status"
            class="toolbar-field"
            clearable
            placeholder="班级状态"
          >
            <el-option
              v-for="status in statusOptions"
              :key="status"
              :label="status"
              :value="status"
            />
          </el-select>
          <el-select
            v-model="filters.campus"
            class="toolbar-field"
            clearable
            placeholder="所属校区"
          >
            <el-option
              v-for="campus in campusOptions"
              :key="campus"
              :label="campus"
              :value="campus"
            />
          </el-select>
        </div>
        <div class="toolbar-actions">
          <el-button plain @click="handleResetFilters">重置筛选</el-button>
        </div>
      </div>

      <div class="data-table-shell">
        <el-table v-loading="loading" :data="filteredClasses" stripe>
          <el-table-column label="班级名称" min-width="220">
            <template #default="{ row }">
              <div class="table-primary">
                <el-button link type="primary" @click="openDetail(row.id)">{{ row.name }}</el-button>
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
          <el-table-column label="起止时间" min-width="180">
            <template #default="{ row }">
              {{ row.startDate || "未填" }} 至 {{ row.endDate || "未填" }}
            </template>
          </el-table-column>
          <el-table-column label="备注" min-width="220">
            <template #default="{ row }">
              <span class="muted-cell">{{ row.remark || "暂无备注" }}</span>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="classStatusTagType(row.status)">{{ row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" :width="canManageClasses ? 160 : 120" fixed="right">
            <template #default="{ row }">
              <div class="table-link-group">
                <el-button link type="primary" @click="openDetail(row.id)">进入详情</el-button>
                <el-button
                  v-if="canManageClasses"
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
      width="840px"
      destroy-on-close
      @close="closeDialog"
      @update:model-value="dialogVisible = $event"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <div class="dialog-grid">
          <el-form-item label="班级名称" prop="name">
            <el-input v-model="form.name" placeholder="例如：周末奥数提高班" />
          </el-form-item>
          <el-form-item label="课程" prop="courseId">
            <el-select v-model="form.courseId" placeholder="请选择课程">
              <el-option
                v-for="option in courseOptions"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="主讲老师" prop="teacherId">
            <el-select v-model="form.teacherId" placeholder="请选择主讲老师">
              <el-option
                v-for="option in teacherOptions"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="所属校区" prop="campus">
            <el-input v-model="form.campus" placeholder="例如：明发校区" />
          </el-form-item>
          <el-form-item label="班级容量" prop="capacity">
            <el-input-number v-model="form.capacity" class="full-width" :min="1" :max="80" />
          </el-form-item>
          <el-form-item label="班级状态" prop="status">
            <el-select v-model="form.status" placeholder="请选择班级状态">
              <el-option
                v-for="status in statusOptions"
                :key="status"
                :label="status"
                :value="status"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="固定排课" prop="weeklySchedule">
            <el-input v-model="form.weeklySchedule" placeholder="例如：周六 09:00-10:30" />
          </el-form-item>
          <el-form-item label="开始日期">
            <el-input v-model="form.startDate" placeholder="格式：2026-06-01" />
          </el-form-item>
          <el-form-item label="结束日期">
            <el-input v-model="form.endDate" placeholder="格式：2026-09-01" />
          </el-form-item>
        </div>

        <el-form-item label="备注">
          <el-input
            v-model="form.remark"
            type="textarea"
            :rows="4"
            placeholder="可以记录开班目标、适用人群或当前补位说明"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-actions">
          <el-button @click="closeDialog">取消</el-button>
          <el-button type="primary" :loading="saving" @click="submitForm">保存班级资料</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>
