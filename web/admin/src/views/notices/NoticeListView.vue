<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import {
  createNotice,
  fetchClassList,
  fetchNoticeList,
  fetchNoticeTargets,
  sendNotice,
  updateNotice,
  type Notice,
  type NoticePayload,
  type NoticeTarget,
  type SchoolClass,
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
const noticeTargets = ref<NoticeTarget[]>([]);
const formRef = ref<FormInstance>();

const filters = reactive({
  keyword: "",
  status: "",
  category: "",
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

const filteredNotices = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase();

  return notices.value.filter((notice) => {
    const matchesKeyword =
      keyword.length === 0 ||
      [notice.title, notice.content, notice.targetScope, notice.author]
        .join(" ")
        .toLowerCase()
        .includes(keyword);
    const matchesStatus = filters.status.length === 0 || notice.status === filters.status;
    const matchesCategory =
      filters.category.length === 0 || notice.category === filters.category;

    return matchesKeyword && matchesStatus && matchesCategory;
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

const relatedClassCount = computed(() => {
  return notices.value.filter((notice) => notice.relatedClassId > 0).length;
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
    status: "草稿",
    author: "教务老师",
  };
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

function openEditDialog(notice: Notice) {
  editingNoticeId.value = notice.id;
  Object.assign(form, {
    title: notice.title,
    content: notice.content,
    category: notice.category,
    targetScope: notice.targetScope,
    relatedClassId: notice.relatedClassId,
    status: notice.status,
    author: notice.author,
  });
  dialogVisible.value = true;
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
    relatedClassId: form.relatedClassId,
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

async function loadNotices() {
  loading.value = true;

  try {
    const [noticeResult, classResult] = await Promise.all([
      fetchNoticeList(),
      fetchClassList(),
    ]);
    notices.value = noticeResult.list;
    classOptions.value = classResult.list;
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

function handleResetFilters() {
  filters.keyword = "";
  filters.status = "";
  filters.category = "";
}

onMounted(() => {
  void loadNotices();
});
</script>

<template>
  <div class="page-stack">
    <section class="page-hero">
      <div class="page-hero__copy">
        <span class="section-kicker">Notice Dispatch</span>
        <h2>把新建、编辑、发送和影响范围放在同一块通知工作台里，消息才能真正可追踪。</h2>
        <p>
          这页现在不只是看列表了。教务和负责人可以直接整理通知内容、确认发送范围，再把调课、停课和日常提醒真正发出去。
        </p>
      </div>

      <div class="metric-strip">
        <article class="metric-tile">
          <span>通知总数</span>
          <strong>{{ notices.length }}</strong>
          <small>当前可回看的全部消息记录</small>
        </article>
        <article class="metric-tile">
          <span>已发送</span>
          <strong>{{ sentCount }}</strong>
          <small>已经推送到目标范围的通知</small>
        </article>
        <article class="metric-tile">
          <span>待发送</span>
          <strong>{{ pendingCount }}</strong>
          <small>接下来最该优先处理的消息</small>
        </article>
        <article class="metric-tile">
          <span>关联班级</span>
          <strong>{{ relatedClassCount }}</strong>
          <small>已经和具体班级绑定的通知数量</small>
        </article>
      </div>
    </section>

    <section class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>通知列表</h2>
          <p class="soft-text">先确认哪些通知还是草稿，哪些要发送，哪些已经发出并可回看影响范围。</p>
        </div>
        <div class="page-actions">
          <div class="section-note">通知视图</div>
          <el-button type="primary" @click="openCreateDialog">新建通知</el-button>
        </div>
      </div>

      <div class="page-toolbar">
        <div class="toolbar-filters">
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
        <el-button @click="handleResetFilters">重置筛选</el-button>
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
          <el-table-column label="范围" prop="targetScope" min-width="180" />
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
      width="680px"
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
              placeholder="不关联具体班级也可以"
              @change="syncTargetScopeFromClass(Number($event || 0))"
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
        </div>

        <el-form-item label="通知范围" prop="targetScope">
          <el-input
            v-model="form.targetScope"
            placeholder="例如：周末奥数提高班家长群 / 全部学员家长 / 待续费学员家长"
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
            <div>
              <h3>{{ currentTargetTitle || "通知范围" }}</h3>
              <p class="soft-text">发错通知时，至少能先知道影响到了谁。</p>
            </div>
            <div class="section-note">目标范围</div>
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
