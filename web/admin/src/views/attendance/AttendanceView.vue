<script setup lang="ts">
import { Check, RefreshRight } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  fetchAttendanceRecordList,
  fetchAttendanceSessionList,
  fetchClassList,
  fetchScheduleAttendance,
  fetchStudentList,
  saveScheduleAttendance,
  type AttendanceRecord,
  type AttendanceEntry,
  type AttendanceSession,
  type SchoolClass,
  type Student,
} from "../../api/education";

const route = useRoute();
const router = useRouter();

const sessionLoading = ref(false);
const detailLoading = ref(false);
const saving = ref(false);
const sessions = ref<AttendanceSession[]>([]);
const historyLoading = ref(false);
const historyRecords = ref<AttendanceRecord[]>([]);
const attendanceItems = ref<AttendanceEntry[]>([]);
const selectedSessionId = ref<number | null>(null);
const classOptions = ref<SchoolClass[]>([]);
const studentOptions = ref<Student[]>([]);

const filters = reactive({
  keyword: "",
  status: "",
});

const historyFilters = reactive({
  classId: null as number | null,
  studentId: null as number | null,
  date: "",
  status: "",
});

const sessionStatusOptions = ["待签到", "已完成", "待上课"];
const entryStatusOptions = ["待确认", "已到", "请假", "缺席", "补签"];
const historyStatusOptions = ["已到", "请假", "缺席", "补签", "待确认"];

const filteredSessions = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase();

  return sessions.value.filter((item) => {
    const matchesKeyword =
      keyword.length === 0 ||
      [item.className, item.courseName, item.teacherName, item.campus]
        .join(" ")
        .toLowerCase()
        .includes(keyword);
    const matchesStatus = filters.status.length === 0 || item.attendanceStatus === filters.status;

    return matchesKeyword && matchesStatus;
  });
});

const selectedSession = computed(() => {
  return sessions.value.find((item) => item.id === selectedSessionId.value) ?? null;
});

const selectedSessionLocked = computed(() => {
  return selectedSession.value?.attendanceStatus === "待上课";
});

const pendingSessionCount = computed(() => {
  return sessions.value.filter((item) => item.attendanceStatus === "待签到").length;
});

const completedSessionCount = computed(() => {
  return sessions.value.filter((item) => item.attendanceStatus === "已完成").length;
});

const leaveCount = computed(() => {
  return sessions.value.reduce((total, item) => total + item.leaveCount, 0);
});

const absentCount = computed(() => {
  return sessions.value.reduce((total, item) => total + item.absentCount, 0);
});

const detailSummary = computed(() => {
  return summarizeAttendance(attendanceItems.value);
});

const historySessions = computed(() => {
  return sessions.value.filter((item) => item.attendanceStatus !== "待上课");
});

function summarizeAttendance(items: AttendanceEntry[]) {
  return items.reduce(
    (summary, item) => {
      switch (item.status) {
        case "已到":
        case "补签":
          summary.presentCount += 1;
          break;
        case "请假":
          summary.leaveCount += 1;
          break;
        case "缺席":
          summary.absentCount += 1;
          break;
        default:
          summary.pendingCount += 1;
          break;
      }

      return summary;
    },
    {
      presentCount: 0,
      leaveCount: 0,
      absentCount: 0,
      pendingCount: 0,
    },
  );
}

function parseScheduleId(rawValue: unknown) {
  const parsedValue = Number(rawValue);
  if (!Number.isInteger(parsedValue) || parsedValue <= 0) {
    return null;
  }

  return parsedValue;
}

function resolveSessionId(preferredId: number | null) {
  if (preferredId !== null && sessions.value.some((item) => item.id === preferredId)) {
    return preferredId;
  }

  const pendingItem = sessions.value.find((item) => item.attendanceStatus === "待签到");
  if (pendingItem) {
    return pendingItem.id;
  }

  return sessions.value[0]?.id ?? null;
}

async function loadAttendanceDetail(sessionId: number) {
  detailLoading.value = true;

  try {
    const result = await fetchScheduleAttendance(sessionId);
    attendanceItems.value = result.items;
  } catch (error) {
    console.error(error);
    attendanceItems.value = [];
    ElMessage.error("签到明细加载失败");
  } finally {
    detailLoading.value = false;
  }
}

async function ensureSelectedSession(preferredId: number | null) {
  const nextId = resolveSessionId(preferredId);
  if (nextId === null) {
    selectedSessionId.value = null;
    attendanceItems.value = [];
    return;
  }

  selectedSessionId.value = nextId;

  const currentQueryId = parseScheduleId(route.query.scheduleId);
  if (currentQueryId !== nextId) {
    await router.replace({
      query: {
        ...route.query,
        scheduleId: String(nextId),
      },
    });
  }

  await loadAttendanceDetail(nextId);
}

async function loadSessions() {
  sessionLoading.value = true;

  try {
    const result = await fetchAttendanceSessionList();
    sessions.value = result.list;
    await ensureSelectedSession(parseScheduleId(route.query.scheduleId));
  } catch (error) {
    console.error(error);
    ElMessage.error("签到台加载失败");
  } finally {
    sessionLoading.value = false;
  }
}

async function loadHistoryRecords() {
  historyLoading.value = true;

  try {
    const result = await fetchAttendanceRecordList({
      mode: "records",
      classId: historyFilters.classId ?? undefined,
      studentId: historyFilters.studentId ?? undefined,
      date: historyFilters.date || undefined,
      status: historyFilters.status || undefined,
    });
    historyRecords.value = result.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("签到记录加载失败");
  } finally {
    historyLoading.value = false;
  }
}

async function loadFilterOptions() {
  try {
    const [classResult, studentResult] = await Promise.all([fetchClassList(), fetchStudentList()]);
    classOptions.value = classResult.list;
    studentOptions.value = studentResult.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("签到筛选项加载失败");
  }
}

async function handleSelectSession(sessionId: number) {
  await ensureSelectedSession(sessionId);
}

function handleResetFilters() {
  filters.keyword = "";
  filters.status = "";
}

function handleResetHistoryFilters() {
  historyFilters.classId = null;
  historyFilters.studentId = null;
  historyFilters.date = "";
  historyFilters.status = "";
  void loadHistoryRecords();
}

function handleMarkAllPresent() {
  attendanceItems.value = attendanceItems.value.map((item) => ({
    ...item,
    status: "已到",
  }));
}

function handleResetPending() {
  attendanceItems.value = attendanceItems.value.map((item) => ({
    ...item,
    status: "待确认",
    remark: "",
  }));
}

async function handleSaveAttendance() {
  const sessionId = selectedSessionId.value;
  if (sessionId === null) {
    return;
  }

  saving.value = true;

  try {
    await saveScheduleAttendance(sessionId, {
      items: attendanceItems.value.map((item) => ({
        studentId: item.studentId,
        status: item.status,
        remark: item.remark,
      })),
    });
    ElMessage.success("签到结果已保存");
    await loadSessions();
    await loadHistoryRecords();
  } catch (error) {
    console.error(error);
    ElMessage.error("签到结果保存失败");
  } finally {
    saving.value = false;
  }
}

function sessionProgressText(item: AttendanceSession) {
  return `已到 ${item.presentCount} · 请假 ${item.leaveCount} · 缺席 ${item.absentCount} · 待确认 ${item.pendingCount}`;
}

function sessionActionLabel(item: AttendanceSession) {
  if (item.attendanceStatus === "待签到") {
    return "开始签到";
  }
  if (item.attendanceStatus === "待上课") {
    return "查看名单";
  }

  return "回看记录";
}

function openHistoryRecord(record: AttendanceRecord) {
  void ensureSelectedSession(record.scheduleId);
}

watch(
  () => route.query.scheduleId,
  (value) => {
    const scheduleId = parseScheduleId(value);
    if (scheduleId === null || scheduleId === selectedSessionId.value) {
      return;
    }

    void ensureSelectedSession(scheduleId);
  },
);

onMounted(() => {
  void loadSessions();
  void loadHistoryRecords();
  void loadFilterOptions();
});
</script>

<template>
  <div class="page-stack">
    <section class="page-hero">
      <div class="page-hero__copy">
        <span class="section-kicker">Attendance Desk</span>
        <h2>把今天要点的名、点完后的结果和历史回看，收进同一个更顺手的签到工作台。</h2>
        <p>
          这块先围绕“快”来做：先看哪些班还没点名，再在单页里改状态、写备注、保存结果，尽量不让老师和教务在几个页面之间来回跳。
        </p>
      </div>

      <div class="metric-strip">
        <article class="metric-tile">
          <span>签到场次</span>
          <strong>{{ sessions.length }}</strong>
          <small>当前演示环境里可回看的全部签到安排</small>
        </article>
        <article class="metric-tile">
          <span>待签到</span>
          <strong>{{ pendingSessionCount }}</strong>
          <small>优先处理今天还没确认完成的班级</small>
        </article>
        <article class="metric-tile">
          <span>已完成</span>
          <strong>{{ completedSessionCount }}</strong>
          <small>已经保存过签到结果的课程场次</small>
        </article>
        <article class="metric-tile">
          <span>异常人数</span>
          <strong>{{ leaveCount + absentCount }}</strong>
          <small>把请假和缺席一起先盯住，方便后续跟进</small>
        </article>
      </div>
    </section>

    <section class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>签到工作台</h2>
          <p class="soft-text">先找出今天该处理的班级，再直接进入点名明细。</p>
        </div>
        <div class="section-note">单页处理</div>
      </div>

      <div class="page-toolbar">
        <div class="toolbar-filters">
          <el-input
            v-model="filters.keyword"
            class="toolbar-field"
            clearable
            placeholder="搜索班级、课程、老师或校区"
          />
          <el-select
            v-model="filters.status"
            class="toolbar-field"
            clearable
            placeholder="签到状态"
          >
            <el-option
              v-for="item in sessionStatusOptions"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
        </div>

        <div class="toolbar-actions">
          <el-button plain @click="handleResetFilters">重置筛选</el-button>
          <el-button plain @click="loadSessions">
            <el-icon><RefreshRight /></el-icon>
            <span>刷新列表</span>
          </el-button>
        </div>
      </div>

      <div class="data-table-shell">
        <el-table v-loading="sessionLoading" :data="filteredSessions" stripe>
          <el-table-column label="日期" prop="lessonDate" width="120" />
          <el-table-column label="时间" prop="lessonTime" width="140" />
          <el-table-column label="班级" prop="className" min-width="180" />
          <el-table-column label="课程" prop="courseName" width="140" />
          <el-table-column label="老师" prop="teacherName" width="120" />
          <el-table-column label="进度" min-width="260">
            <template #default="{ row }">
              {{ sessionProgressText(row) }}
            </template>
          </el-table-column>
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="row.attendanceStatus === '已完成' ? 'success' : 'warning'">
                {{ row.attendanceStatus }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="handleSelectSession(row.id)">
                {{ sessionActionLabel(row) }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </section>

    <section v-if="selectedSession" class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>{{ selectedSession.className }} · 签到明细</h2>
          <p class="soft-text">
            {{ selectedSession.lessonDate }} {{ selectedSession.lessonTime }} ·
            {{ selectedSession.courseName }} · {{ selectedSession.teacherName }} ·
            {{ selectedSession.classroom }}
          </p>
        </div>
        <div class="section-note">
          {{ selectedSessionLocked ? "未开始" : "可保存" }}
        </div>
      </div>

      <div class="metric-strip">
        <article class="metric-tile">
          <span>已到</span>
          <strong>{{ detailSummary.presentCount }}</strong>
          <small>包含正常到课和补签学员</small>
        </article>
        <article class="metric-tile">
          <span>请假</span>
          <strong>{{ detailSummary.leaveCount }}</strong>
          <small>可在备注里补充原因</small>
        </article>
        <article class="metric-tile">
          <span>缺席</span>
          <strong>{{ detailSummary.absentCount }}</strong>
          <small>后续适合班主任继续跟进</small>
        </article>
        <article class="metric-tile">
          <span>待确认</span>
          <strong>{{ detailSummary.pendingCount }}</strong>
          <small>还没最终落定的签到结果</small>
        </article>
      </div>

      <div class="page-toolbar">
        <div class="toolbar-filters">
          <span class="soft-text">
            默认更适合先把少数异常学员改出来，再一次性保存整场签到结果。
          </span>
        </div>

        <div class="toolbar-actions">
          <el-button plain :disabled="selectedSessionLocked || detailLoading" @click="handleResetPending">
            恢复待确认
          </el-button>
          <el-button plain :disabled="selectedSessionLocked || detailLoading" @click="handleMarkAllPresent">
            <el-icon><Check /></el-icon>
            <span>全员到课</span>
          </el-button>
          <el-button
            type="primary"
            :disabled="selectedSessionLocked || detailLoading"
            :loading="saving"
            @click="handleSaveAttendance"
          >
            保存签到
          </el-button>
        </div>
      </div>

      <div v-if="selectedSessionLocked" class="attendance-lock-note">
        这节课还没开始，当前先提供名单预览；等进入上课时段后再保存正式签到结果会更合适。
      </div>

      <div class="data-table-shell">
        <el-table v-loading="detailLoading" :data="attendanceItems" stripe>
          <el-table-column label="学员" prop="studentName" width="140" />
          <el-table-column label="年级" prop="grade" width="110" />
          <el-table-column label="家长电话" prop="parentMobile" width="140" />
          <el-table-column label="签到状态" width="160">
            <template #default="{ row }">
              <el-select
                v-model="row.status"
                class="attendance-status-field"
                :disabled="selectedSessionLocked"
              >
                <el-option
                  v-for="item in entryStatusOptions"
                  :key="item"
                  :label="item"
                  :value="item"
                />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="备注" min-width="240">
            <template #default="{ row }">
              <el-input
                v-model="row.remark"
                class="attendance-remark-field"
                :disabled="selectedSessionLocked"
                placeholder="补充请假原因、课堂表现或补签说明"
              />
            </template>
          </el-table-column>
        </el-table>
      </div>
    </section>

    <section class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>签到记录回看</h2>
          <p class="soft-text">按班级、学员和日期回看历史签到，也能快速跳回对应场次继续更正。</p>
        </div>
        <div class="section-note">历史记录</div>
      </div>

      <div class="page-toolbar">
        <div class="toolbar-filters">
          <el-select
            v-model="historyFilters.classId"
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
            v-model="historyFilters.studentId"
            class="toolbar-field"
            clearable
            filterable
            placeholder="按学员查看"
          >
            <el-option
              v-for="item in studentOptions"
              :key="item.id"
              :label="`${item.name} · ${item.grade}`"
              :value="item.id"
            />
          </el-select>
          <el-date-picker
            v-model="historyFilters.date"
            class="toolbar-field"
            clearable
            type="date"
            value-format="YYYY-MM-DD"
            placeholder="按日期筛选"
          />
          <el-select
            v-model="historyFilters.status"
            class="toolbar-field"
            clearable
            placeholder="异常状态"
          >
            <el-option
              v-for="item in historyStatusOptions"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
        </div>

        <div class="toolbar-actions">
          <el-button plain @click="handleResetHistoryFilters">重置筛选</el-button>
          <el-button plain @click="loadHistoryRecords">
            <el-icon><RefreshRight /></el-icon>
            <span>刷新记录</span>
          </el-button>
        </div>
      </div>

      <div class="data-table-shell">
        <el-table v-loading="historyLoading" :data="historyRecords" stripe>
          <el-table-column label="日期" prop="lessonDate" width="120" />
          <el-table-column label="时间" prop="lessonTime" width="120" />
          <el-table-column label="班级" prop="className" min-width="180" />
          <el-table-column label="学员" prop="studentName" width="120" />
          <el-table-column label="老师" prop="teacherName" width="120" />
          <el-table-column label="家长电话" prop="parentMobile" width="140" />
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag
                :type="
                  row.status === '已到' || row.status === '补签'
                    ? 'success'
                    : row.status === '请假'
                      ? 'warning'
                      : row.status === '缺席'
                        ? 'danger'
                        : 'info'
                "
              >
                {{ row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="备注" min-width="180">
            <template #default="{ row }">
              <span class="muted-cell">{{ row.remark || "暂无备注" }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作人" prop="updatedBy" width="120" />
          <el-table-column label="最后更新时间" prop="updatedAt" width="180" />
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="openHistoryRecord(row)">
                去更正
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div v-if="historyRecords.length === 0" class="soft-empty teacher-detail-empty">
        当前筛选条件下还没有签到记录。
      </div>
    </section>
  </div>
</template>
