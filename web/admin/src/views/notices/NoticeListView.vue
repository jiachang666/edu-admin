<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import {
  createNotice,
  fetchClassList,
  fetchNotice,
  fetchNoticeList,
  fetchNoticeTargets,
  fetchScheduleList,
  fetchStudentList,
  sendNotice,
  updateNotice,
  type Notice,
  type NoticePayload,
  type NoticeTarget,
  type Schedule,
  type SchoolClass,
  type Student,
} from "../../api/education";

const loading = ref(false);
const saving = ref(false);
const sendingNoticeId = ref<number | null>(null);
const dialogVisible = ref(false);
const targetDrawerVisible = ref(false);
const editingNoticeId = ref<number | null>(null);
const currentTargetTitle = ref("");
const notices = ref<Notice[]>([]);
const classOptions = ref<SchoolClass[]>([]);
const scheduleOptions = ref<Schedule[]>([]);
const studentOptions = ref<Student[]>([]);
const noticeTargets = ref<NoticeTarget[]>([]);
const formRef = ref<FormInstance>();

const filters = reactive({
  keyword: "",
  status: "",
  category: "",
  classId: null as number | null,
  dateRange: [] as string[],
});

const form = reactive<NoticePayload>(defaultForm());

const statusOptions = ["草稿", "待发送", "已发送"];
const categoryOptions = ["校区通知", "课程通知", "续费提醒", "调课通知", "停课提醒"];

const rules: FormRules<NoticePayload> = {
  title: [{ required: true, message: "请输入通知标题", trigger: "blur" }],
  content: [{ required: true, message: "请输入通知内容", trigger: "blur" }],
  category: [{ required: true, message: "请选择通知分类", trigger: "change" }],
  targetScope: [{ required: true, message: "请输入通知范围", trigger: "blur" }],
  status: [{ required: true, message: "请选择通知状态", trigger: "change" }],
  author: [{ required: true, message: "请输入发起人", trigger: "blur" }],
};

const classNameMap = computed(() => {
  return new Map(classOptions.value.map((item) => [item.id, item.name]));
});

const scheduleNameMap = computed(() => {
  return new Map(
    scheduleOptions.value.map((item) => [
      item.id,
      `${item.lessonDate} ${item.className} ${item.lessonTime}`,
    ]),
  );
});

const studentNameMap = computed(() => {
  return new Map(studentOptions.value.map((item) => [item.id, item.name]));
});

const schedulePickerOptions = computed(() => {
  if (!form.relatedClassId) {
    return scheduleOptions.value;
  }

  return scheduleOptions.value.filter((item) => item.classId === form.relatedClassId);
});

const studentPickerOptions = computed(() => {
  if (!form.relatedClassId) {
    return studentOptions.value;
  }

  return studentOptions.value.filter((item) => item.classId === form.relatedClassId);
});

const filteredNotices = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase();

  return notices.value.filter((notice) => {
    const matchesKeyword =
      keyword.length === 0 ||
      [notice.title, notice.content, notice.targetScope, notice.author, describeNoticeRelation(notice)]
        .join(" ")
        .toLowerCase()
        .includes(keyword);
    const matchesStatus = filters.status.length === 0 || notice.status === filters.status;
    const matchesCategory =
      filters.category.length === 0 || notice.category === filters.category;
    const matchesClass =
      filters.classId === null || notice.relatedClassId === filters.classId;
    const matchesDate = matchesDateRange(extractDate(notice.publishAt), filters.dateRange);

    return matchesKeyword && matchesStatus && matchesCategory && matchesClass && matchesDate;
  });
});

const sentCount = computed(() => {
  return notices.value.filter((notice) => notice.status === "已发送").length;
});

const pendingCount = computed(() => {
  return notices.value.filter((notice) => notice.status === "待发送").length;
});

const draftCount = computed(() => {
  return notices.value.filter((notice) => notice.status === "草稿").length;
});

const relatedScheduleCount = computed(() => {
  return notices.value.filter((notice) => notice.relatedScheduleId > 0).length;
});

const dialogTitle = computed(() => {
  return editingNoticeId.value ? "编辑通知" : "新建通知";
});

function defaultForm(): NoticePayload {
  return {
    title: "",
    content: "",
    category: "校区通知",
    targetScope: "",
    relatedClassId: 0,
    relatedScheduleId: 0,
    studentIds: [],
    status: "草稿",
    author: "教务老师",
  };
}

function extractDate(value: string) {
  return value.slice(0, 10);
}

function matchesDateRange(value: string, dateRange: string[]) {
  if (dateRange.length !== 2) {
    return true;
  }

  const [dateFrom, dateTo] = dateRange;
  if (!dateFrom || !dateTo || value.length === 0) {
    return true;
  }

  return value >= dateFrom && value <= dateTo;
}

function buildStudentScope(studentIds: number[]) {
  const studentNames = studentIds
    .map((studentId) => studentNameMap.value.get(studentId) ?? "")
    .filter(Boolean);

  if (studentNames.length === 0) {
    return "指定学员家长";
  }

  if (studentNames.length === 1) {
    return `${studentNames[0]}家长`;
  }

  if (studentNames.length === 2) {
    return `${studentNames[0]}、${studentNames[1]}家长`;
  }

  return `${studentNames[0]}等${studentNames.length}位学员家长`;
}

function describeNoticeRelation(notice: Notice) {
  const relationParts: string[] = [];

  if (notice.relatedClassId > 0) {
    relationParts.push(classNameMap.value.get(notice.relatedClassId) ?? `班级 #${notice.relatedClassId}`);
  }

  if (notice.relatedScheduleId > 0) {
    relationParts.push(scheduleNameMap.value.get(notice.relatedScheduleId) ?? `课程安排 #${notice.relatedScheduleId}`);
  }

  if (notice.studentIds.length > 0) {
    relationParts.push(`指定 ${notice.studentIds.length} 位学员`);
  }

  if (relationParts.length === 0) {
    return notice.targetScope || "未绑定具体对象";
  }

  return relationParts.join(" / ");
}

function resetForm() {
  Object.assign(form, defaultForm());
  editingNoticeId.value = null;
  formRef.value?.clearValidate();
}

function openCreateDialog() {
  resetForm();
  dialogVisible.value = true;
}

async function openEditDialog(notice: Notice) {
  try {
    const detail = await fetchNotice(notice.id);
    editingNoticeId.value = detail.id;
    Object.assign(form, {
      title: detail.title,
      content: detail.content,
      category: detail.category,
      targetScope: detail.targetScope,
      relatedClassId: detail.relatedClassId,
      relatedScheduleId: detail.relatedScheduleId,
      studentIds: [...detail.studentIds],
      status: detail.status,
      author: detail.author,
    });
    dialogVisible.value = true;
    formRef.value?.clearValidate();
  } catch (error) {
    console.error(error);
    ElMessage.error("通知详情加载失败");
  }
}

function closeDialog() {
  dialogVisible.value = false;
  resetForm();
}

function buildPayload(): NoticePayload {
  return {
    title: form.title.trim(),
    content: form.content.trim(),
    category: form.category,
    targetScope: form.targetScope.trim(),
    relatedClassId: Number(form.relatedClassId) || 0,
    relatedScheduleId: Number(form.relatedScheduleId) || 0,
    studentIds: Array.from(new Set(form.studentIds.map((item) => Number(item)).filter((item) => item > 0))),
    status: form.status,
    author: form.author.trim(),
  };
}

function syncTargetScopeFromClass(classId: number) {
  if (!classId) {
    return;
  }

  const relatedClass = classOptions.value.find((item) => item.id === classId);
  if (!relatedClass) {
    return;
  }

  if (!form.targetScope.trim()) {
    form.targetScope = `${relatedClass.name}家长群`;
  }
}

function handleClassChange(rawValue: number | string | undefined) {
  const classId = Number(rawValue || 0);
  form.relatedClassId = classId;

  if (classId > 0) {
    form.studentIds = form.studentIds.filter((studentId) => {
      return studentOptions.value.some((item) => item.id === studentId && item.classId === classId);
    });
  }

  if (form.relatedScheduleId > 0) {
    const relatedSchedule = scheduleOptions.value.find((item) => item.id === form.relatedScheduleId);
    if (relatedSchedule && relatedSchedule.classId !== classId) {
      form.relatedScheduleId = 0;
    }
  }

  syncTargetScopeFromClass(classId);
}

function handleScheduleChange(rawValue: number | string | undefined) {
  const scheduleId = Number(rawValue || 0);
  form.relatedScheduleId = scheduleId;
  if (!scheduleId) {
    return;
  }

  const relatedSchedule = scheduleOptions.value.find((item) => item.id === scheduleId);
  if (!relatedSchedule) {
    return;
  }

  if (relatedSchedule.classId > 0) {
    form.relatedClassId = relatedSchedule.classId;
  }

  if (!form.targetScope.trim()) {
    form.targetScope = `${relatedSchedule.className}家长群`;
  }
}

function handleStudentChange(studentIds: number[]) {
  form.studentIds = Array.from(new Set(studentIds.map((item) => Number(item)).filter((item) => item > 0)));
  if (form.studentIds.length === 0) {
    return;
  }

  if (!form.targetScope.trim()) {
    form.targetScope = buildStudentScope(form.studentIds);
  }
}

function buildNoticeQuery() {
  const [dateFrom, dateTo] = filters.dateRange;

  return {
    classId: filters.classId ?? undefined,
    status: filters.status || undefined,
    noticeType: filters.category || undefined,
    dateFrom: dateFrom || undefined,
    dateTo: dateTo || undefined,
  };
}

async function loadNoticeOptions() {
  try {
    const [classResult, scheduleResult, studentResult] = await Promise.all([
      fetchClassList(),
      fetchScheduleList(),
      fetchStudentList(),
    ]);
    classOptions.value = classResult.list;
    scheduleOptions.value = scheduleResult.list;
    studentOptions.value = studentResult.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("通知关联数据加载失败");
  }
}

async function loadNotices() {
  loading.value = true;

  try {
    const noticeResult = await fetchNoticeList(buildNoticeQuery());
    notices.value = noticeResult.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("通知列表加载失败");
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

    if (editingNoticeId.value) {
      await updateNotice(editingNoticeId.value, payload);
      ElMessage.success("通知已更新");
    } else {
      await createNotice(payload);
      ElMessage.success("通知已创建");
    }

    closeDialog();
    await loadNotices();
  } catch (error) {
    console.error(error);
    ElMessage.error("通知保存失败");
  } finally {
    saving.value = false;
  }
}

async function handleSend(notice: Notice) {
  sendingNoticeId.value = notice.id;

  try {
    await sendNotice(notice.id);
    ElMessage.success("通知已发送");
    await loadNotices();
  } catch (error) {
    console.error(error);
    ElMessage.error("通知发送失败");
  } finally {
    sendingNoticeId.value = null;
  }
}

async function openTargets(notice: Notice) {
  currentTargetTitle.value = notice.title;
  targetDrawerVisible.value = true;

  try {
    noticeTargets.value = await fetchNoticeTargets(notice.id);
  } catch (error) {
    console.error(error);
    noticeTargets.value = [];
    ElMessage.error("通知范围加载失败");
  }
}

async function handleSearch() {
  await loadNotices();
}

async function handleResetFilters() {
  filters.keyword = "";
  filters.status = "";
  filters.category = "";
  filters.classId = null;
  filters.dateRange = [];
  await loadNotices();
}

onMounted(async () => {
  await Promise.all([loadNoticeOptions(), loadNotices()]);
});
</script>

<template>
  <div class="page-stack">
    <section class="page-card page-card--table list-card">
      <div class="page-header">
        <div class="list-card__heading">
          <h2>通知列表</h2>
          <span class="list-card__count">共 {{ filteredNotices.length }} 条</span>
        </div>
        <div class="page-actions">
          <el-button type="primary" @click="openCreateDialog">新建通知</el-button>
        </div>
      </div>

      <div class="metric-strip metric-strip--compact list-card__metrics">
        <article class="metric-tile">
          <span>通知总数</span>
          <strong>{{ notices.length }}</strong>
        </article>
        <article class="metric-tile">
          <span>已发送</span>
          <strong>{{ sentCount }}</strong>
        </article>
        <article class="metric-tile">
          <span>待发送</span>
          <strong>{{ pendingCount }}</strong>
        </article>
        <article class="metric-tile">
          <span>草稿中</span>
          <strong>{{ draftCount }}</strong>
        </article>
        <article class="metric-tile">
          <span>关联课程</span>
          <strong>{{ relatedScheduleCount }}</strong>
        </article>
      </div>

      <div class="filter-bar list-card__filters">
        <div class="toolbar-filters">
          <el-select
            v-model="filters.classId"
            class="toolbar-field"
            clearable
            filterable
            placeholder="按班级查看"
          >
            <el-option
              v-for="item in classOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
          <el-date-picker
            v-model="filters.dateRange"
            class="toolbar-field"
            clearable
            type="daterange"
            value-format="YYYY-MM-DD"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
          />
          <el-input
            v-model="filters.keyword"
            class="toolbar-field"
            clearable
            placeholder="搜索标题、内容、范围或发起人"
          />
          <el-select
            v-model="filters.category"
            class="toolbar-field"
            clearable
            placeholder="通知分类"
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
            placeholder="通知状态"
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
          <el-button type="primary" @click="handleSearch">筛选结果</el-button>
          <el-button @click="handleResetFilters">重置筛选</el-button>
        </div>
      </div>

      <div class="data-table-shell">
        <el-table v-loading="loading" :data="filteredNotices" stripe>
          <el-table-column label="标题" min-width="220">
            <template #default="{ row }">
              <div class="table-primary">
                <strong>{{ row.title }}</strong>
                <small>{{ row.content }}</small>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="分类" prop="category" width="140" />
          <el-table-column label="关联对象" min-width="220">
            <template #default="{ row }">
              <div class="table-primary">
                <strong>{{ describeNoticeRelation(row) }}</strong>
                <small>{{ row.targetScope }}</small>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="row.status === '已发送' ? 'success' : row.status === '草稿' ? 'info' : 'warning'">
                {{ row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="时间" prop="publishAt" width="180" />
          <el-table-column label="发起人" prop="author" width="120" />
          <el-table-column label="操作" min-width="220" fixed="right">
            <template #default="{ row }">
              <div class="table-link-group">
                <el-button link type="primary" @click="openEditDialog(row)">编辑</el-button>
                <el-button link type="primary" @click="openTargets(row)">看范围</el-button>
                <el-button
                  link
                  type="primary"
                  :disabled="row.status === '已发送'"
                  :loading="sendingNoticeId === row.id"
                  @click="handleSend(row)"
                >
                  {{ row.status === "已发送" ? "已发送" : "立即发送" }}
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
          <el-form-item label="通知标题" prop="title">
            <el-input v-model="form.title" placeholder="例如：周末奥数班调课提醒" />
          </el-form-item>
          <el-form-item label="通知分类" prop="category">
            <el-select v-model="form.category" placeholder="请选择通知分类">
              <el-option
                v-for="category in categoryOptions"
                :key="category"
                :label="category"
                :value="category"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="关联班级">
            <el-select
              v-model="form.relatedClassId"
              clearable
              filterable
              placeholder="不关联具体班级也可以"
              @change="handleClassChange"
            >
              <el-option
                v-for="schoolClass in classOptions"
                :key="schoolClass.id"
                :label="schoolClass.name"
                :value="schoolClass.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="通知状态" prop="status">
            <el-select v-model="form.status" placeholder="请选择通知状态">
              <el-option
                v-for="status in statusOptions"
                :key="status"
                :label="status"
                :value="status"
              />
            </el-select>
          </el-form-item>
          <el-form-item class="full-span" label="关联课程安排">
            <el-select
              v-model="form.relatedScheduleId"
              clearable
              filterable
              placeholder="把通知和某次具体上课安排绑定起来"
              @change="handleScheduleChange"
            >
              <el-option
                v-for="item in schedulePickerOptions"
                :key="item.id"
                :label="`${item.lessonDate} ${item.className} ${item.lessonTime}`"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item class="full-span" label="指定学员">
            <el-select
              v-model="form.studentIds"
              multiple
              filterable
              collapse-tags
              collapse-tags-tooltip
              placeholder="需要单独通知到某几位学员时可以直接选"
              @change="handleStudentChange"
            >
              <el-option
                v-for="item in studentPickerOptions"
                :key="item.id"
                :label="`${item.name} · ${item.className || item.grade}`"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
        </div>

        <el-form-item label="通知范围" prop="targetScope">
          <el-input
            v-model="form.targetScope"
            placeholder="例如：周末奥数提高班家长群 / 全部学员家长 / 指定学员家长"
          />
        </el-form-item>

        <el-form-item label="通知内容" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="5"
            placeholder="把调课说明、停课原因、材料准备或日常提醒写清楚"
          />
        </el-form-item>

        <el-form-item label="发起人" prop="author">
          <el-input v-model="form.author" placeholder="例如：教务老师" />
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-actions">
          <el-button @click="closeDialog">取消</el-button>
          <el-button type="primary" :loading="saving" @click="submitForm">保存通知</el-button>
        </div>
      </template>
    </el-dialog>

    <el-drawer
      v-model="targetDrawerVisible"
      title="通知影响范围"
      size="460px"
      destroy-on-close
    >
      <div class="page-stack">
        <section class="page-card notice-target-card">
          <div class="page-header notice-target-card__header">
            <h3>{{ currentTargetTitle || "通知范围" }}</h3>
          </div>

          <div v-if="noticeTargets.length === 0" class="soft-empty">
            当前没有可展示的目标范围。
          </div>

          <div v-else class="notice-target-list">
            <article
              v-for="target in noticeTargets"
              :key="`${target.name}-${target.type}-${target.campus}`"
              class="notice-target-item"
            >
              <strong>{{ target.name }}</strong>
              <small>{{ target.type }} · {{ target.campus || "未区分校区" }}</small>
            </article>
          </div>
        </section>
      </div>
    </el-drawer>
  </div>
</template>
