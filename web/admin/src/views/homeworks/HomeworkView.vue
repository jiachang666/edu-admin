<script setup lang="ts">
import { ElMessage } from "element-plus";
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  fetchHomeworkList,
  fetchScheduleFeedback,
  fetchScheduleHomework,
  fetchScheduleList,
  saveScheduleFeedback,
  saveScheduleHomework,
  type Feedback,
  type Homework,
  type Schedule,
} from "../../api/education";

const route = useRoute();
const router = useRouter();

const listLoading = ref(false);
const detailLoading = ref(false);
const savingHomework = ref(false);
const savingFeedback = ref(false);
const homeworkItems = ref<Homework[]>([]);
const schedules = ref<Schedule[]>([]);
const selectedSchedule = ref<Schedule | null>(null);

const filters = reactive({
  keyword: "",
  status: "",
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

    return matchesKeyword && matchesStatus;
  });
});

const publishedCount = computed(() => {
  return homeworkItems.value.filter((item) => item.status === "published").length;
});

const draftCount = computed(() => {
  return homeworkItems.value.filter((item) => item.status === "draft").length;
});

const coveredClassCount = computed(() => {
  return new Set(homeworkItems.value.map((item) => item.classId)).size;
});

const todayScheduleOptions = computed(() => {
  return schedules.value.filter((item) => item.attendanceStatus !== "待上课");
});

function parseScheduleId(rawValue: unknown) {
  const parsedValue = Number(rawValue);
  if (!Number.isInteger(parsedValue) || parsedValue <= 0) {
    return null;
  }

  return parsedValue;
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
    const [homeworkResult, scheduleResult] = await Promise.all([
      fetchHomeworkList(),
      fetchScheduleList(),
    ]);
    homeworkItems.value = homeworkResult.list;
    schedules.value = scheduleResult.list;
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
    schedules.value.find((item) => item.id === preferredId) ??
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

function handleResetFilters() {
  filters.keyword = "";
  filters.status = "";
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
    <section class="page-hero">
      <div class="page-hero__copy">
        <span class="section-kicker">Homework Studio</span>
        <h2>把每次课后的作业布置和班级反馈放到同一条课后处理线上，不再散在群消息里。</h2>
        <p>
          这页先围绕老师和教务最常做的两件事来做：一是按上课安排快速发布作业，二是补上班级级别的课后反馈，让家长和机构都能回看。
        </p>
      </div>

      <div class="metric-strip">
        <article class="metric-tile">
          <span>作业总数</span>
          <strong>{{ homeworkItems.length }}</strong>
          <small>当前系统里可回看的全部课后作业</small>
        </article>
        <article class="metric-tile">
          <span>已发布</span>
          <strong>{{ publishedCount }}</strong>
          <small>已经可直接通知家长和学员的作业</small>
        </article>
        <article class="metric-tile">
          <span>草稿</span>
          <strong>{{ draftCount }}</strong>
          <small>还在整理中的课后任务内容</small>
        </article>
        <article class="metric-tile">
          <span>覆盖班级</span>
          <strong>{{ coveredClassCount }}</strong>
          <small>目前已经做过课后跟进的班级数量</small>
        </article>
      </div>
    </section>

    <section class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>作业列表</h2>
          <p class="soft-text">先看哪些班已经布置过作业，再进入单次课程安排继续编辑。</p>
        </div>
        <div class="section-note">作业视图</div>
      </div>

      <div class="page-toolbar">
        <div class="toolbar-filters">
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

    <section class="page-card">
      <div class="page-header">
        <div>
          <h2>按课程安排处理课后内容</h2>
          <p class="soft-text">一次上课对应一份作业和一份班级反馈，这样后面回看才不会乱。</p>
        </div>
        <div class="section-note">单次上课</div>
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
              v-for="item in todayScheduleOptions"
              :key="item.id"
              :label="`${item.lessonDate} ${item.className} ${item.lessonTime}`"
              :value="item.id"
            />
          </el-select>
        </div>

        <div class="toolbar-actions">
          <div class="soft-text" v-if="selectedSchedule">
            {{ selectedSchedule.className }} · {{ selectedSchedule.courseName }} ·
            {{ selectedSchedule.teacherName }}
          </div>
        </div>
      </div>

      <div v-loading="detailLoading" class="homework-editor-grid">
        <section class="page-card homework-editor-card">
          <div class="page-header">
            <div>
              <h3>作业内容</h3>
              <p class="soft-text">优先把标题、内容和提交说明整理清楚。</p>
            </div>
            <div class="section-note">老师布置</div>
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
            <div>
              <h3>课后反馈</h3>
              <p class="soft-text">先按班级维度记录课堂表现和下次建议。</p>
            </div>
            <div class="section-note">班级反馈</div>
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
