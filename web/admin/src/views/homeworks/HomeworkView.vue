<script setup lang="ts">
import { ElMessage } from "element-plus";
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  fetchClassList,
  fetchFeedbackList,
  fetchHomeworkList,
  fetchScheduleFeedback,
  fetchScheduleHomework,
  fetchScheduleList,
  fetchTeacherOptions,
  saveScheduleFeedback,
  saveScheduleHomework,
  type Feedback,
  type Homework,
  type Schedule,
  type SchoolClass,
  type SelectOption,
} from "../../api/education";

const route = useRoute();
const router = useRouter();

const listLoading = ref(false);
const detailLoading = ref(false);
const savingHomework = ref(false);
const savingFeedback = ref(false);
const homeworkItems = ref<Homework[]>([]);
const feedbackItems = ref<Feedback[]>([]);
const schedules = ref<Schedule[]>([]);
const classOptions = ref<SchoolClass[]>([]);
const teacherOptions = ref<SelectOption[]>([]);
const selectedSchedule = ref<Schedule | null>(null);

const filters = reactive({
  keyword: "",
  status: "",
  classId: null as number | null,
  teacherId: null as number | null,
  dateRange: [] as string[],
});

const homeworkForm = reactive({
  title: "",
  content: "",
  submissionNote: "",
  status: "published",
});

const feedbackForm = reactive({
  summary: "",
  learningStatus: "",
  nextSuggestion: "",
  parentNotice: "",
});

const homeworkStatusOptions = [
  { label: "已发布", value: "published" },
  { label: "草稿", value: "draft" },
];

const teacherNameMap = computed(() => {
  return new Map(teacherOptions.value.map((item) => [item.value, item.label]));
});

const feedbackCount = computed(() => {
  return feedbackItems.value.length;
});

const filteredHomeworkItems = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase();

  return homeworkItems.value.filter((item) => {
    const matchesKeyword =
      keyword.length === 0 ||
      [item.title, item.className, item.courseName, item.teacherName]
        .join(" ")
        .toLowerCase()
        .includes(keyword);
    const matchesStatus = filters.status.length === 0 || item.status === filters.status;
    const matchesClass = filters.classId === null || item.classId === filters.classId;
    const matchesTeacher =
      filters.teacherId === null ||
      schedules.value.some((scheduleItem) => {
        return (
          scheduleItem.id === item.scheduleId &&
          scheduleItem.teacherId === filters.teacherId
        );
      });
    const matchesDate = matchesDateRange(item.lessonDate, filters.dateRange);

    return matchesKeyword && matchesStatus && matchesClass && matchesTeacher && matchesDate;
  });
});

const filteredFeedbackItems = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase();

  return feedbackItems.value.filter((item) => {
    const matchesKeyword =
      keyword.length === 0 ||
      [item.className, item.courseName, item.teacherName, item.summary, item.learningStatus]
        .join(" ")
        .toLowerCase()
        .includes(keyword);
    const matchesClass = filters.classId === null || item.classId === filters.classId;
    const matchesTeacher =
      filters.teacherId === null ||
      schedules.value.some((scheduleItem) => {
        return (
          scheduleItem.id === item.scheduleId &&
          scheduleItem.teacherId === filters.teacherId
        );
      });
    const matchesDate = matchesDateRange(item.lessonDate, filters.dateRange);

    return matchesKeyword && matchesClass && matchesTeacher && matchesDate;
  });
});

const publishedCount = computed(() => {
  return homeworkItems.value.filter((item) => item.status === "published").length;
});

const draftCount = computed(() => {
  return homeworkItems.value.filter((item) => item.status === "draft").length;
});

const todayScheduleOptions = computed(() => {
  return schedules.value.filter((item) => item.attendanceStatus !== "待上课");
});

const schedulePickerOptions = computed(() => {
  const selectedTeacherName =
    filters.teacherId === null ? "" : (teacherNameMap.value.get(filters.teacherId) ?? "");

  return todayScheduleOptions.value.filter((item) => {
    const matchesClass = filters.classId === null || item.classId === filters.classId;
    const matchesTeacher =
      selectedTeacherName.length === 0 || item.teacherName === selectedTeacherName;

    return matchesClass && matchesTeacher;
  });
});

function matchesDateRange(value: string, dateRange: string[]) {
  if (dateRange.length !== 2) {
    return true;
  }

  const [dateFrom, dateTo] = dateRange;
  if (!dateFrom || !dateTo) {
    return true;
  }

  return value >= dateFrom && value <= dateTo;
}

function parseScheduleId(rawValue: unknown) {
  const parsedValue = Number(rawValue);
  if (!Number.isInteger(parsedValue) || parsedValue <= 0) {
    return null;
  }

  return parsedValue;
}

function buildListQuery() {
  const [dateFrom, dateTo] = filters.dateRange;

  return {
    classId: filters.classId ?? undefined,
    teacherId: filters.teacherId ?? undefined,
    dateFrom: dateFrom || undefined,
    dateTo: dateTo || undefined,
  };
}

function resetForms() {
  homeworkForm.title = "";
  homeworkForm.content = "";
  homeworkForm.submissionNote = "";
  homeworkForm.status = "published";
  feedbackForm.summary = "";
  feedbackForm.learningStatus = "";
  feedbackForm.nextSuggestion = "";
  feedbackForm.parentNotice = "";
}

async function loadListData() {
  listLoading.value = true;

  try {
    const query = buildListQuery();
    const [homeworkResult, feedbackResult, scheduleResult, classResult, teacherResult] =
      await Promise.all([
        fetchHomeworkList(query),
        fetchFeedbackList(query),
        fetchScheduleList(),
        fetchClassList(),
        fetchTeacherOptions(),
      ]);
    homeworkItems.value = homeworkResult.list;
    feedbackItems.value = feedbackResult.list;
    schedules.value = scheduleResult.list;
    classOptions.value = classResult.list;
    teacherOptions.value = teacherResult;
  } catch (error) {
    console.error(error);
    ElMessage.error("作业与反馈列表加载失败");
  } finally {
    listLoading.value = false;
  }
}

async function loadDetail(scheduleId: number) {
  const scheduleItem = schedules.value.find((item) => item.id === scheduleId) ?? null;
  selectedSchedule.value = scheduleItem;
  if (!scheduleItem) {
    resetForms();
    return;
  }

  detailLoading.value = true;

  try {
    const [homeworkResult, feedbackResult] = await Promise.all([
      fetchScheduleHomework(scheduleId),
      fetchScheduleFeedback(scheduleId),
    ]);

    homeworkForm.title = String(homeworkResult.title ?? "");
    homeworkForm.content = String(homeworkResult.content ?? "");
    homeworkForm.submissionNote = String(homeworkResult.submissionNote ?? "");
    homeworkForm.status = String(homeworkResult.status ?? "published");

    feedbackForm.summary = String(feedbackResult.summary ?? "");
    feedbackForm.learningStatus = String(feedbackResult.learningStatus ?? "");
    feedbackForm.nextSuggestion = String(feedbackResult.nextSuggestion ?? "");
    feedbackForm.parentNotice = String(feedbackResult.parentNotice ?? "");
  } catch (error) {
    console.error(error);
    ElMessage.error("作业或反馈明细加载失败");
  } finally {
    detailLoading.value = false;
  }
}

async function ensureSelectedSchedule(preferredId: number | null) {
  const nextSchedule =
    schedulePickerOptions.value.find((item) => item.id === preferredId) ??
    schedulePickerOptions.value[0] ??
    todayScheduleOptions.value.find((item) => item.id === preferredId) ??
    todayScheduleOptions.value[0] ??
    schedules.value[0] ??
    null;

  selectedSchedule.value = nextSchedule;
  if (!nextSchedule) {
    resetForms();
    return;
  }

  if (parseScheduleId(route.query.scheduleId) !== nextSchedule.id) {
    await router.replace({
      query: {
        ...route.query,
        scheduleId: String(nextSchedule.id),
      },
    });
  }

  await loadDetail(nextSchedule.id);
}

async function handleSelectSchedule(scheduleId: number) {
  await ensureSelectedSchedule(scheduleId);
}

async function handleSearch() {
  await loadListData();
  await ensureSelectedSchedule(selectedSchedule.value?.id ?? parseScheduleId(route.query.scheduleId));
}

async function handleResetFilters() {
  filters.keyword = "";
  filters.status = "";
  filters.classId = null;
  filters.teacherId = null;
  filters.dateRange = [];
  await loadListData();
  await ensureSelectedSchedule(parseScheduleId(route.query.scheduleId));
}

async function handleSaveHomework() {
  if (!selectedSchedule.value) {
    return;
  }

  savingHomework.value = true;

  try {
    await saveScheduleHomework(selectedSchedule.value.id, {
      title: homeworkForm.title,
      content: homeworkForm.content,
      submissionNote: homeworkForm.submissionNote,
      status: homeworkForm.status,
    });
    ElMessage.success("作业内容已保存");
    await loadListData();
    await loadDetail(selectedSchedule.value.id);
  } catch (error) {
    console.error(error);
    ElMessage.error("作业保存失败");
  } finally {
    savingHomework.value = false;
  }
}

async function handleSaveFeedback() {
  if (!selectedSchedule.value) {
    return;
  }

  savingFeedback.value = true;

  try {
    await saveScheduleFeedback(selectedSchedule.value.id, {
      summary: feedbackForm.summary,
      learningStatus: feedbackForm.learningStatus,
      nextSuggestion: feedbackForm.nextSuggestion,
      parentNotice: feedbackForm.parentNotice,
    });
    ElMessage.success("课后反馈已保存");
    await loadListData();
    await loadDetail(selectedSchedule.value.id);
  } catch (error) {
    console.error(error);
    ElMessage.error("课后反馈保存失败");
  } finally {
    savingFeedback.value = false;
  }
}

watch(
  () => route.query.scheduleId,
  (value) => {
    const scheduleId = parseScheduleId(value);
    if (scheduleId === null || scheduleId === selectedSchedule.value?.id) {
      return;
    }

    void ensureSelectedSchedule(scheduleId);
  },
);

onMounted(async () => {
  await loadListData();
  await ensureSelectedSchedule(parseScheduleId(route.query.scheduleId));
});
</script>

<template>
  <div class="page-stack">
    <section class="page-card page-card--table list-card">
      <div class="page-header">
        <div class="list-card__heading">
          <h2>作业列表</h2>
          <span class="list-card__count">共 {{ filteredHomeworkItems.length }} 条</span>
        </div>
      </div>

      <div class="metric-strip metric-strip--compact list-card__metrics">
        <article class="metric-tile">
          <span>作业总数</span>
          <strong>{{ homeworkItems.length }}</strong>
        </article>
        <article class="metric-tile">
          <span>已发布</span>
          <strong>{{ publishedCount }}</strong>
        </article>
        <article class="metric-tile">
          <span>草稿</span>
          <strong>{{ draftCount }}</strong>
        </article>
        <article class="metric-tile">
          <span>课后反馈</span>
          <strong>{{ feedbackCount }}</strong>
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
          <el-select
            v-model="filters.teacherId"
            class="toolbar-field"
            clearable
            filterable
            placeholder="按老师查看"
          >
            <el-option
              v-for="item in teacherOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
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
            placeholder="搜索作业标题、班级、课程或老师"
          />
          <el-select
            v-model="filters.status"
            class="toolbar-field"
            clearable
            placeholder="作业状态"
          >
            <el-option
              v-for="item in homeworkStatusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </div>

        <div class="toolbar-actions">
          <el-button type="primary" @click="handleSearch">筛选结果</el-button>
          <el-button plain @click="handleResetFilters">重置筛选</el-button>
        </div>
      </div>

      <div class="data-table-shell">
        <el-table v-loading="listLoading" :data="filteredHomeworkItems" stripe>
          <el-table-column label="日期" prop="lessonDate" width="120" />
          <el-table-column label="班级" prop="className" min-width="180" />
          <el-table-column label="课程" prop="courseName" width="140" />
          <el-table-column label="老师" prop="teacherName" width="120" />
          <el-table-column label="作业标题" prop="title" min-width="220" />
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="row.status === 'published' ? 'success' : 'warning'">
                {{ row.status === "published" ? "已发布" : "草稿" }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="handleSelectSchedule(row.scheduleId)">
                编辑内容
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </section>

    <section class="page-card page-card--table list-card">
      <div class="page-header">
        <div class="list-card__heading">
          <h2>课后反馈</h2>
          <span class="list-card__count">共 {{ filteredFeedbackItems.length }} 条</span>
        </div>
      </div>

      <div class="data-table-shell">
        <el-table v-loading="listLoading" :data="filteredFeedbackItems" stripe>
          <el-table-column label="日期" prop="lessonDate" width="120" />
          <el-table-column label="班级" prop="className" min-width="180" />
          <el-table-column label="课程" prop="courseName" width="140" />
          <el-table-column label="老师" prop="teacherName" width="120" />
          <el-table-column label="课堂反馈" min-width="260">
            <template #default="{ row }">
              <div class="table-primary">
                <strong>{{ row.summary || "暂无反馈摘要" }}</strong>
                <small>{{ row.learningStatus || row.parentNotice || "暂时还没有补充说明" }}</small>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="下次建议" min-width="220">
            <template #default="{ row }">
              <span class="muted-cell">{{ row.nextSuggestion || "暂无建议" }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="handleSelectSchedule(row.scheduleId)">
                查看并编辑
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div v-if="filteredFeedbackItems.length === 0" class="soft-empty">
        当前筛选条件下还没有课后反馈记录。
      </div>
    </section>

    <section class="page-card">
      <div class="page-header">
        <div>
          <h2>按课程安排处理</h2>
          <p v-if="selectedSchedule" class="soft-text">
            {{ selectedSchedule.className }} · {{ selectedSchedule.courseName }} ·
            {{ selectedSchedule.teacherName }}
          </p>
        </div>
      </div>

      <div class="page-toolbar">
        <div class="toolbar-filters">
          <el-select
            :model-value="selectedSchedule?.id"
            class="toolbar-field"
            placeholder="选择课程安排"
            @update:model-value="handleSelectSchedule"
          >
            <el-option
              v-for="item in schedulePickerOptions"
              :key="item.id"
              :label="`${item.lessonDate} ${item.className} ${item.lessonTime}`"
              :value="item.id"
            />
          </el-select>
        </div>
      </div>

      <div v-loading="detailLoading" class="homework-editor-grid">
        <section class="page-card homework-editor-card">
          <div class="page-header">
            <h3>作业内容</h3>
          </div>

          <div class="form-grid">
            <el-form-item class="full-span" label="作业标题">
              <el-input
                v-model="homeworkForm.title"
                placeholder="例如：思维训练第 4 讲课后练习"
              />
            </el-form-item>
            <el-form-item class="full-span" label="作业内容">
              <el-input
                v-model="homeworkForm.content"
                type="textarea"
                placeholder="填写本次课后需要完成的练习或复习内容"
              />
            </el-form-item>
            <el-form-item class="full-span" label="提交说明">
              <el-input
                v-model="homeworkForm.submissionNote"
                type="textarea"
                placeholder="例如：下节课前带回纸质作业，或发到班级群"
              />
            </el-form-item>
            <el-form-item label="作业状态">
              <el-select v-model="homeworkForm.status" class="full-width">
                <el-option
                  v-for="item in homeworkStatusOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </el-form-item>
          </div>

          <div class="homework-editor-actions">
            <el-button type="primary" :loading="savingHomework" @click="handleSaveHomework">
              保存作业
            </el-button>
          </div>
        </section>

        <section class="page-card homework-editor-card">
          <div class="page-header">
            <h3>课后反馈</h3>
          </div>

          <div class="form-grid">
            <el-form-item class="full-span" label="课堂表现">
              <el-input
                v-model="feedbackForm.summary"
                type="textarea"
                placeholder="概括本节课整体课堂状态和参与情况"
              />
            </el-form-item>
            <el-form-item class="full-span" label="学习情况">
              <el-input
                v-model="feedbackForm.learningStatus"
                type="textarea"
                placeholder="记录这次课掌握得好的点和仍需加强的点"
              />
            </el-form-item>
            <el-form-item class="full-span" label="下次建议">
              <el-input
                v-model="feedbackForm.nextSuggestion"
                type="textarea"
                placeholder="写给老师、教务和家长都能理解的下一步建议"
              />
            </el-form-item>
            <el-form-item class="full-span" label="家长关注事项">
              <el-input
                v-model="feedbackForm.parentNotice"
                type="textarea"
                placeholder="补充需要家长配合留意的事项"
              />
            </el-form-item>
          </div>

          <div class="homework-editor-actions">
            <el-button type="primary" :loading="savingFeedback" @click="handleSaveFeedback">
              保存反馈
            </el-button>
          </div>
        </section>
      </div>
    </section>
  </div>
</template>
