<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import {
  createCourse,
  fetchCourseList,
  updateCourse,
  type Course,
  type CoursePayload,
} from "../../api/education";

const statusOptions = ["启用", "停用"];
const deliveryTypeOptions = ["线下", "线上", "线上线下结合"];

const loading = ref(false);
const saving = ref(false);
const dialogVisible = ref(false);
const editingCourseId = ref<number | null>(null);
const courses = ref<Course[]>([]);
const formRef = ref<FormInstance>();

const filters = reactive({
  keyword: "",
  category: "",
  status: "",
});

const form = reactive<CoursePayload>(defaultForm());

const rules: FormRules<CoursePayload> = {
  name: [{ required: true, message: "请输入课程名称", trigger: "blur" }],
  category: [{ required: true, message: "请输入课程分类", trigger: "blur" }],
  lessonDurationMinutes: [
    { required: true, message: "请输入单次时长", trigger: "blur" },
  ],
  totalLessons: [
    { required: true, message: "请输入建议总课时", trigger: "blur" },
  ],
  deliveryType: [
    { required: true, message: "请选择授课方式", trigger: "change" },
  ],
  status: [{ required: true, message: "请选择课程状态", trigger: "change" }],
};

const categoryOptions = computed(() => {
  const categories = new Set<string>();

  for (const course of courses.value) {
    if (course.category) {
      categories.add(course.category);
    }
  }

  return Array.from(categories);
});

const enabledCount = computed(() => {
  return courses.value.filter((course) => course.status === "启用").length;
});

const categoryCount = computed(() => {
  return categoryOptions.value.length;
});

const linkedClassCount = computed(() => {
  return courses.value.reduce((total, course) => total + course.classCount, 0);
});

const dialogTitle = computed(() => {
  return editingCourseId.value ? "编辑课程" : "新增课程";
});

function defaultForm(): CoursePayload {
  return {
    name: "",
    category: "",
    description: "",
    ageRange: "",
    lessonDurationMinutes: 90,
    totalLessons: 16,
    deliveryType: "线下",
    status: "启用",
  };
}

function resetForm() {
  Object.assign(form, defaultForm());
  editingCourseId.value = null;
  formRef.value?.clearValidate();
}

function openCreateDialog() {
  resetForm();
  dialogVisible.value = true;
}

function openEditDialog(course: Course) {
  editingCourseId.value = course.id;
  Object.assign(form, {
    name: course.name,
    category: course.category,
    description: course.description,
    ageRange: course.ageRange,
    lessonDurationMinutes: course.lessonDurationMinutes,
    totalLessons: course.totalLessons,
    deliveryType: course.deliveryType,
    status: course.status,
  });
  dialogVisible.value = true;
}

function closeDialog() {
  dialogVisible.value = false;
  resetForm();
}

function buildPayload(): CoursePayload {
  return {
    name: form.name.trim(),
    category: form.category.trim(),
    description: form.description.trim(),
    ageRange: form.ageRange.trim(),
    lessonDurationMinutes: form.lessonDurationMinutes,
    totalLessons: form.totalLessons,
    deliveryType: form.deliveryType,
    status: form.status,
  };
}

async function loadCourses() {
  loading.value = true;

  try {
    const result = await fetchCourseList({
      keyword: filters.keyword.trim() || undefined,
      category: filters.category || undefined,
      status: filters.status || undefined,
    });
    courses.value = result.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("课程列表加载失败");
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

    if (editingCourseId.value) {
      await updateCourse(editingCourseId.value, payload);
      ElMessage.success("课程已更新");
    } else {
      await createCourse(payload);
      ElMessage.success("课程已创建");
    }

    closeDialog();
    await loadCourses();
  } catch (error) {
    console.error(error);
    ElMessage.error("课程保存失败");
  } finally {
    saving.value = false;
  }
}

function handleSearch() {
  void loadCourses();
}

function handleReset() {
  filters.keyword = "";
  filters.category = "";
  filters.status = "";
  void loadCourses();
}

onMounted(() => {
  void loadCourses();
});
</script>

<template>
  <div class="page-stack">
    <section class="page-hero">
      <div class="page-hero__copy">
        <span class="section-kicker">Curriculum Board</span>
        <h2>把课程模板的内容、阶段和授课方式先定清楚，后面的建班和排课才会稳。</h2>
        <p>
          这页现在已经不只是展示列表，也能直接新增和编辑课程。先把课程台账做清楚，首版闭环就更像一个真正可用的行业后台。
        </p>
      </div>

      <div class="metric-strip">
        <article class="metric-tile">
          <span>课程总数</span>
          <strong>{{ courses.length }}</strong>
          <small>当前系统里的全部课程模板</small>
        </article>
        <article class="metric-tile">
          <span>启用中</span>
          <strong>{{ enabledCount }}</strong>
          <small>仍然对外开班使用的课程</small>
        </article>
        <article class="metric-tile">
          <span>课程分类</span>
          <strong>{{ categoryCount }}</strong>
          <small>当前已覆盖的分类数量</small>
        </article>
        <article class="metric-tile">
          <span>关联班级</span>
          <strong>{{ linkedClassCount }}</strong>
          <small>课程模板目前被班级复用的总次数</small>
        </article>
      </div>
    </section>

    <section class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>课程列表</h2>
          <p class="soft-text">先把课程模板信息补齐，后面的建班和排课才有清晰基础。</p>
        </div>
        <div class="section-note">课程视图</div>
      </div>

      <div class="page-toolbar">
        <div class="toolbar-filters">
          <el-input
            v-model="filters.keyword"
            class="toolbar-field"
            clearable
            placeholder="搜索课程名称、分类、简介"
            @keyup.enter="handleSearch"
          />
          <el-select
            v-model="filters.category"
            class="toolbar-field"
            clearable
            placeholder="课程分类"
          >
            <el-option
              v-for="category in categoryOptions"
              :key="category"
              :label="category"
              :value="category"
            />
          </el-select>
          <el-select
            v-model="filters.status"
            class="toolbar-field"
            clearable
            placeholder="课程状态"
          >
            <el-option
              v-for="status in statusOptions"
              :key="status"
              :label="status"
              :value="status"
            />
          </el-select>
        </div>

        <div class="toolbar-actions">
          <el-button @click="handleReset">重置</el-button>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button type="success" @click="openCreateDialog">新增课程</el-button>
        </div>
      </div>

      <div class="data-table-shell">
        <el-table v-loading="loading" :data="courses" stripe>
          <el-table-column label="课程名称" prop="name" min-width="140" />
          <el-table-column label="课程分类" prop="category" width="120" />
          <el-table-column label="授课方式" prop="deliveryType" width="140" />
          <el-table-column label="适用年龄/阶段" prop="ageRange" min-width="140" />
          <el-table-column label="课时设置" min-width="160">
            <template #default="{ row }">
              {{ row.lessonDurationMinutes }} 分钟 / {{ row.totalLessons }} 课时
            </template>
          </el-table-column>
          <el-table-column label="关联班级" prop="classCount" width="100" />
          <el-table-column label="课程简介" min-width="240">
            <template #default="{ row }">
              <span class="muted-cell">
                {{ row.description || "暂无简介" }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === '启用' ? 'success' : 'info'">
                {{ row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="openEditDialog(row)">
                编辑
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <el-dialog
        :model-value="dialogVisible"
        :title="dialogTitle"
        width="720px"
        @close="closeDialog"
        @update:model-value="dialogVisible = $event"
      >
        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-position="top"
        >
          <div class="form-grid">
            <el-form-item label="课程名称" prop="name">
              <el-input v-model="form.name" placeholder="例如：数学思维" />
            </el-form-item>

            <el-form-item label="课程分类" prop="category">
              <el-input v-model="form.category" placeholder="例如：数学" />
            </el-form-item>

            <el-form-item label="适用年龄/阶段" prop="ageRange">
              <el-input v-model="form.ageRange" placeholder="例如：8-10岁" />
            </el-form-item>

            <el-form-item label="授课方式" prop="deliveryType">
              <el-select v-model="form.deliveryType" placeholder="选择授课方式">
                <el-option
                  v-for="deliveryType in deliveryTypeOptions"
                  :key="deliveryType"
                  :label="deliveryType"
                  :value="deliveryType"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="单次时长（分钟）" prop="lessonDurationMinutes">
              <el-input-number
                v-model="form.lessonDurationMinutes"
              :min="30"
              :step="15"
              class="full-width"
              />
            </el-form-item>

            <el-form-item label="建议总课时" prop="totalLessons">
              <el-input-number
                v-model="form.totalLessons"
              :min="1"
              :step="1"
              class="full-width"
              />
            </el-form-item>

            <el-form-item class="full-span" label="课程简介" prop="description">
              <el-input
                v-model="form.description"
                :rows="4"
                maxlength="300"
                placeholder="写一下课程适合谁、主要教什么、上课方式是什么"
                show-word-limit
                type="textarea"
              />
            </el-form-item>

            <el-form-item label="课程状态" prop="status">
              <el-select v-model="form.status" placeholder="选择课程状态">
                <el-option
                  v-for="status in statusOptions"
                  :key="status"
                  :label="status"
                  :value="status"
                />
              </el-select>
            </el-form-item>
          </div>
        </el-form>

        <template #footer>
          <el-button @click="closeDialog">取消</el-button>
          <el-button :loading="saving" type="primary" @click="submitForm">
            保存
          </el-button>
        </template>
      </el-dialog>
    </section>
  </div>
</template>
