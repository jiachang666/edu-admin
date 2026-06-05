<script setup lang="ts">
import { Calendar, Clock, EditPen, RefreshRight } from "@element-plus/icons-vue";
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import {
  cancelLesson,
  createMakeupLesson,
  createSchedule,
  fetchClassList,
  fetchScheduleList,
  rescheduleLesson,
  updateSchedule,
  type Schedule,
  type ScheduleActionPayload,
  type SchedulePayload,
  type SchoolClass,
} from "../../api/education";

type DialogMode = "create" | "edit" | "reschedule" | "cancel" | "makeup";

const router = useRouter();

const loading = ref(false);
const saving = ref(false);
const dialogVisible = ref(false);
const schedules = ref<Schedule[]>([]);
const classes = ref<SchoolClass[]>([]);
const selectedSchedule = ref<Schedule | null>(null);
const dialogMode = ref<DialogMode>("create");
const formRef = ref<FormInstance>();

const filters = reactive({
  keyword: "",
  scheduleType: "",
  status: "",
  campus: "",
});

const scheduleForm = reactive<SchedulePayload>({
  classId: 0,
  scheduleType: "常规课",
  lessonDate: "",
  startTime: "",
  endTime: "",
  classroom: "",
  remark: "",
});

const actionForm = reactive<ScheduleActionPayload>({
  lessonDate: "",
  startTime: "",
  endTime: "",
  classroom: "",
  remark: "",
});

const scheduleTypeOptions = ["常规课", "调课", "补课", "试听课"];
const statusOptions = ["待上课", "待签到", "已完成", "已调课", "已停课", "请假待批"];

const rules: FormRules<SchedulePayload> = {
  classId: [{ required: true, message: "请选择班级", trigger: "change" }],
  lessonDate: [{ required: true, message: "请选择上课日期", trigger: "change" }],
  startTime: [{ required: true, message: "请选择开始时间", trigger: "change" }],
  endTime: [{ required: true, message: "请选择结束时间", trigger: "change" }],
  classroom: [{ required: true, message: "请输入教室或地点", trigger: "blur" }],
};

const actionRules = computed<FormRules<ScheduleActionPayload>>(() => {
  if (dialogMode.value === "cancel") {
    return {
      remark: [{ required: true, message: "请填写停课原因", trigger: "blur" }],
    };
  }

  return {
    lessonDate: [{ required: true, message: "请选择日期", trigger: "change" }],
    startTime: [{ required: true, message: "请选择开始时间", trigger: "change" }],
    endTime: [{ required: true, message: "请选择结束时间", trigger: "change" }],
    classroom: [{ required: true, message: "请输入教室或地点", trigger: "blur" }],
  };
});

const filteredSchedules = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase();

  return schedules.value.filter((item) => {
    const matchesKeyword =
      keyword.length === 0 ||
      [item.className, item.courseName, item.teacherName, item.classroom, item.remark]
        .join(" ")
        .toLowerCase()
        .includes(keyword);
    const matchesType =
      filters.scheduleType.length === 0 || item.scheduleType === filters.scheduleType;
    const matchesStatus = filters.status.length === 0 || item.attendanceStatus === filters.status;
    const matchesCampus = filters.campus.length === 0 || item.campus === filters.campus;

    return matchesKeyword && matchesType && matchesStatus && matchesCampus;
  });
});

const doneCount = computed(() => {
  return schedules.value.filter((item) => item.attendanceStatus === "已完成").length;
});

const pendingCount = computed(() => {
  return schedules.value.filter((item) => item.attendanceStatus === "待签到").length;
});

const campusCount = computed(() => {
  return new Set(schedules.value.map((item) => item.campus).filter(Boolean)).size;
});

const waitingCount = computed(() => {
  return schedules.value.filter((item) => item.attendanceStatus === "待上课").length;
});

const campusOptions = computed(() => {
  return Array.from(new Set(schedules.value.map((item) => item.campus).filter(Boolean)));
});

const dialogTitle = computed(() => {
  switch (dialogMode.value) {
    case "create":
      return "新建排课";
    case "edit":
      return "编辑排课";
    case "reschedule":
      return "调课";
    case "cancel":
      return "停课";
    case "makeup":
      return "补课";
    default:
      return "排课操作";
  }
});

const selectedClassSummary = computed(() => {
  return classes.value.find((item) => item.id === scheduleForm.classId) ?? null;
});

function resetScheduleForm() {
  scheduleForm.classId = 0;
  scheduleForm.scheduleType = "常规课";
  scheduleForm.lessonDate = "";
  scheduleForm.startTime = "";
  scheduleForm.endTime = "";
  scheduleForm.classroom = "";
  scheduleForm.remark = "";
}

function resetActionForm() {
  actionForm.lessonDate = "";
  actionForm.startTime = "";
  actionForm.endTime = "";
  actionForm.classroom = "";
  actionForm.remark = "";
}

function closeDialog() {
  dialogVisible.value = false;
  selectedSchedule.value = null;
  resetScheduleForm();
  resetActionForm();
  formRef.value?.clearValidate();
}

function syncScheduleFormWithItem(item: Schedule) {
  const [startTime = "", endTime = ""] = item.lessonTime.split("-");

  scheduleForm.classId = item.classId;
  scheduleForm.scheduleType = item.scheduleType || "常规课";
  scheduleForm.lessonDate = item.lessonDate;
  scheduleForm.startTime = startTime;
  scheduleForm.endTime = endTime;
  scheduleForm.classroom = item.classroom;
  scheduleForm.remark = item.remark || "";
}

function syncActionFormWithItem(item: Schedule) {
  const [startTime = "", endTime = ""] = item.lessonTime.split("-");

  actionForm.lessonDate = item.lessonDate;
  actionForm.startTime = startTime;
  actionForm.endTime = endTime;
  actionForm.classroom = item.classroom;
  actionForm.remark = item.remark || "";
}

function openCreateDialog() {
  dialogMode.value = "create";
  selectedSchedule.value = null;
  resetScheduleForm();
  dialogVisible.value = true;
  formRef.value?.clearValidate();
}

function openEditDialog(item: Schedule) {
  dialogMode.value = "edit";
  selectedSchedule.value = item;
  syncScheduleFormWithItem(item);
  dialogVisible.value = true;
  formRef.value?.clearValidate();
}

function openRescheduleDialog(item: Schedule) {
  dialogMode.value = "reschedule";
  selectedSchedule.value = item;
  syncActionFormWithItem(item);
  dialogVisible.value = true;
  formRef.value?.clearValidate();
}

function openCancelDialog(item: Schedule) {
  dialogMode.value = "cancel";
  selectedSchedule.value = item;
  resetActionForm();
  actionForm.remark = item.remark || "";
  dialogVisible.value = true;
  formRef.value?.clearValidate();
}

function openMakeupDialog(item: Schedule) {
  dialogMode.value = "makeup";
  selectedSchedule.value = item;
  syncActionFormWithItem(item);
  actionForm.lessonDate = "";
  actionForm.remark = item.remark ? `补课说明：${item.remark}` : "";
  dialogVisible.value = true;
  formRef.value?.clearValidate();
}

function handleResetFilters() {
  filters.keyword = "";
  filters.scheduleType = "";
  filters.status = "";
  filters.campus = "";
}

async function loadPageData() {
  loading.value = true;

  try {
    const [scheduleResult, classResult] = await Promise.all([fetchScheduleList(), fetchClassList()]);
    schedules.value = scheduleResult.list;
    classes.value = classResult.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("排课工作台加载失败");
  } finally {
    loading.value = false;
  }
}

async function submitDialog() {
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
    if (dialogMode.value === "create") {
      await createSchedule({
        classId: scheduleForm.classId,
        scheduleType: scheduleForm.scheduleType,
        lessonDate: scheduleForm.lessonDate,
        startTime: scheduleForm.startTime,
        endTime: scheduleForm.endTime,
        classroom: scheduleForm.classroom.trim(),
        remark: scheduleForm.remark.trim(),
      });
      ElMessage.success("新排课已创建");
    } else if (dialogMode.value === "edit" && selectedSchedule.value) {
      await updateSchedule(selectedSchedule.value.id, {
        classId: scheduleForm.classId,
        scheduleType: scheduleForm.scheduleType,
        lessonDate: scheduleForm.lessonDate,
        startTime: scheduleForm.startTime,
        endTime: scheduleForm.endTime,
        classroom: scheduleForm.classroom.trim(),
        remark: scheduleForm.remark.trim(),
      });
      ElMessage.success("排课已更新");
    } else if (dialogMode.value === "reschedule" && selectedSchedule.value) {
      await rescheduleLesson(selectedSchedule.value.id, {
        lessonDate: actionForm.lessonDate,
        startTime: actionForm.startTime,
        endTime: actionForm.endTime,
        classroom: actionForm.classroom.trim(),
        remark: actionForm.remark.trim(),
      });
      ElMessage.success("调课已完成");
    } else if (dialogMode.value === "cancel" && selectedSchedule.value) {
      await cancelLesson(selectedSchedule.value.id, {
        remark: actionForm.remark.trim(),
      });
      ElMessage.success("停课已记录");
    } else if (dialogMode.value === "makeup" && selectedSchedule.value) {
      await createMakeupLesson(selectedSchedule.value.id, {
        lessonDate: actionForm.lessonDate,
        startTime: actionForm.startTime,
        endTime: actionForm.endTime,
        classroom: actionForm.classroom.trim(),
        remark: actionForm.remark.trim(),
      });
      ElMessage.success("补课安排已创建");
    }

    closeDialog();
    await loadPageData();
  } catch (error) {
    console.error(error);
    ElMessage.error("排课操作失败");
  } finally {
    saving.value = false;
  }
}

function openDetail(scheduleId: number) {
  void router.push(`/schedules/${scheduleId}`);
}

function openAttendance(scheduleId: number) {
  void router.push({
    path: "/attendance",
    query: { scheduleId: String(scheduleId) },
  });
}

function openHomework(scheduleId: number) {
  void router.push({
    path: "/homeworks",
    query: { scheduleId: String(scheduleId) },
  });
}

function attendanceTagType(status: string) {
  switch (status) {
    case "已完成":
      return "success";
    case "待签到":
      return "warning";
    case "已调课":
      return "primary";
    case "已停课":
      return "danger";
    case "请假待批":
      return "warning";
    default:
      return "info";
  }
}

function scheduleTone(type: string) {
  switch (type) {
    case "调课":
      return "blue";
    case "补课":
      return "teal";
    case "试听课":
      return "orange";
    default:
      return "green";
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
        <span class="section-kicker">
          <el-icon><Calendar /></el-icon>
          Schedule Rail
        </span>
        <h2>把新增、改时间、停课和补课收进同一个排课工作台，教务每天就有一条清楚主线。</h2>
        <p>
          这页不再只是看列表。我们把排课信息、调整动作和后续签到作业入口都放到一起，方便先处理当天安排，再进入单次课程详情继续跟进。
        </p>
      </div>

      <div class="metric-strip">
        <article class="metric-tile">
          <span>排课总数</span>
          <strong>{{ schedules.length }}</strong>
          <small>当前环境里全部可处理的课程安排</small>
        </article>
        <article class="metric-tile">
          <span>待上课</span>
          <strong>{{ waitingCount }}</strong>
          <small>还没到签到环节的课程安排</small>
        </article>
        <article class="metric-tile">
          <span>待签到</span>
          <strong>{{ pendingCount }}</strong>
          <small>老师上完课后优先要跟进的场次</small>
        </article>
        <article class="metric-tile">
          <span>已完成</span>
          <strong>{{ doneCount }}</strong>
          <small>签到结果已经闭环的课程安排</small>
        </article>
        <article class="metric-tile">
          <span>涉及校区</span>
          <strong>{{ campusCount }}</strong>
          <small>方便看本周安排覆盖了多少教学点</small>
        </article>
      </div>
    </section>

    <section class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>排课工作台</h2>
          <p class="soft-text">先筛出目标课程，再直接新建、编辑、调课、停课或补课。</p>
        </div>
        <div class="page-actions">
          <div class="section-note">单页处理</div>
          <el-button type="primary" @click="openCreateDialog">新建排课</el-button>
        </div>
      </div>

      <div class="page-toolbar">
        <div class="toolbar-filters">
          <el-input
            v-model="filters.keyword"
            class="toolbar-field"
            clearable
            placeholder="搜索班级、课程、老师、教室或备注"
          />
          <el-select
            v-model="filters.scheduleType"
            class="toolbar-field"
            clearable
            placeholder="排课类型"
          >
            <el-option
              v-for="item in scheduleTypeOptions"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
          <el-select
            v-model="filters.status"
            class="toolbar-field"
            clearable
            placeholder="当前状态"
          >
            <el-option
              v-for="item in statusOptions"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
          <el-select
            v-model="filters.campus"
            class="toolbar-field"
            clearable
            placeholder="校区"
          >
            <el-option
              v-for="item in campusOptions"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
        </div>

        <div class="toolbar-actions">
          <el-button plain @click="handleResetFilters">重置筛选</el-button>
        </div>
      </div>

      <div class="data-table-shell">
        <el-table v-loading="loading" :data="filteredSchedules" stripe>
          <el-table-column label="日期" prop="lessonDate" width="120" />
          <el-table-column label="时间" prop="lessonTime" width="120" />
          <el-table-column label="班级" min-width="220">
            <template #default="{ row }">
              <div class="table-primary">
                <strong>{{ row.className }}</strong>
                <small>{{ row.courseName }} · {{ row.teacherName }}</small>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="校区/教室" min-width="160">
            <template #default="{ row }">
              <div class="table-primary">
                <strong>{{ row.classroom }}</strong>
                <small>{{ row.campus }}</small>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="类型" width="110">
            <template #default="{ row }">
              <el-tag :type="scheduleTone(row.scheduleType)">{{ row.scheduleType || "常规课" }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="110">
            <template #default="{ row }">
              <el-tag :type="attendanceTagType(row.attendanceStatus)">
                {{ row.attendanceStatus }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="备注" min-width="180">
            <template #default="{ row }">
              <span class="muted-cell">{{ row.remark || "暂无备注" }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" min-width="260" fixed="right">
            <template #default="{ row }">
              <div class="table-link-group">
                <el-button link type="primary" @click="openDetail(row.id)">进入详情</el-button>
                <el-button link type="primary" @click="openEditDialog(row)">编辑</el-button>
                <el-button link type="primary" @click="openRescheduleDialog(row)">调课</el-button>
                <el-button link type="danger" @click="openCancelDialog(row)">停课</el-button>
                <el-button link type="primary" @click="openMakeupDialog(row)">补课</el-button>
                <el-button link type="primary" @click="openAttendance(row.id)">
                  {{ row.attendanceStatus === "待签到" ? "去签到" : "看签到" }}
                </el-button>
                <el-button link type="primary" @click="openHomework(row.id)">作业反馈</el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </section>

    <section class="page-card">
      <div class="page-header">
        <div>
          <h2>排课动作说明</h2>
          <p class="soft-text">把常见动作的处理口径写清楚，减少教务和老师之间来回确认。</p>
        </div>
        <div class="section-note">处理规则</div>
      </div>

      <div class="quick-actions-grid">
        <button class="quick-action-card" type="button" @click="openCreateDialog">
          <strong>新增一节课</strong>
          <span>直接补上一条新的上课安排，适合常规加课、试听课或临时增加课程。</span>
        </button>
        <button class="quick-action-card" type="button" @click="router.push('/attendance')">
          <strong>去签到台</strong>
          <span>先把今天该点名的班级处理掉，再回来看哪些课需要继续跟进。</span>
        </button>
        <button class="quick-action-card" type="button" @click="router.push('/homeworks')">
          <strong>去课后处理</strong>
          <span>签到完成后，继续补作业和反馈，让这次课真正形成闭环。</span>
        </button>
        <button class="quick-action-card" type="button" @click="router.push('/notices')">
          <strong>去通知中心</strong>
          <span>调课、停课后的家长通知，可以直接去消息工作台统一整理和发送。</span>
        </button>
      </div>
    </section>

    <el-dialog
      :model-value="dialogVisible"
      :title="dialogTitle"
      width="640px"
      destroy-on-close
      @close="closeDialog"
      @update:model-value="dialogVisible = $event"
    >
      <template v-if="dialogMode === 'create' || dialogMode === 'edit'">
        <el-form ref="formRef" :model="scheduleForm" :rules="rules" label-position="top">
          <div class="dialog-grid">
            <el-form-item label="班级" prop="classId">
              <el-select v-model="scheduleForm.classId" class="full-width" filterable placeholder="请选择班级">
                <el-option
                  v-for="item in classes"
                  :key="item.id"
                  :label="`${item.name} · ${item.courseName} · ${item.teacherName}`"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="排课类型" prop="scheduleType">
              <el-select v-model="scheduleForm.scheduleType" class="full-width" placeholder="请选择排课类型">
                <el-option
                  v-for="item in scheduleTypeOptions"
                  :key="item"
                  :label="item"
                  :value="item"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="上课日期" prop="lessonDate">
              <el-date-picker
                v-model="scheduleForm.lessonDate"
                class="full-width"
                type="date"
                value-format="YYYY-MM-DD"
                placeholder="请选择日期"
              />
            </el-form-item>

            <el-form-item label="教室或地点" prop="classroom">
              <el-input v-model="scheduleForm.classroom" placeholder="例如 A201 或 线上会议室" />
            </el-form-item>

            <el-form-item label="开始时间" prop="startTime">
              <el-time-picker
                v-model="scheduleForm.startTime"
                class="full-width"
                format="HH:mm"
                value-format="HH:mm"
                placeholder="请选择开始时间"
              />
            </el-form-item>

            <el-form-item label="结束时间" prop="endTime">
              <el-time-picker
                v-model="scheduleForm.endTime"
                class="full-width"
                format="HH:mm"
                value-format="HH:mm"
                placeholder="请选择结束时间"
              />
            </el-form-item>

            <el-form-item class="full-span" label="备注">
              <el-input
                v-model="scheduleForm.remark"
                :rows="3"
                maxlength="255"
                show-word-limit
                type="textarea"
                placeholder="可填写调班说明、家长备注或教室安排"
              />
            </el-form-item>
          </div>
        </el-form>

        <div v-if="selectedClassSummary" class="detail-note">
          <strong>当前班级</strong>
          <p>
            {{ selectedClassSummary.name }} · {{ selectedClassSummary.courseName }} ·
            {{ selectedClassSummary.teacherName }} · {{ selectedClassSummary.campus }}
          </p>
        </div>
      </template>

      <template v-else>
        <el-form ref="formRef" :model="actionForm" :rules="actionRules" label-position="top">
          <div v-if="selectedSchedule" class="detail-note">
            <strong>原始安排</strong>
            <p>
              {{ selectedSchedule.lessonDate }} {{ selectedSchedule.lessonTime }} ·
              {{ selectedSchedule.className }} · {{ selectedSchedule.classroom }}
            </p>
          </div>

          <div v-if="dialogMode === 'cancel'" class="dialog-grid">
            <el-form-item class="full-span" label="停课原因" prop="remark">
              <el-input
                v-model="actionForm.remark"
                :rows="4"
                maxlength="255"
                show-word-limit
                type="textarea"
                placeholder="例如老师请假、节假日停课、家长统一请假"
              />
            </el-form-item>
          </div>

          <div v-else class="dialog-grid">
            <el-form-item label="新日期" prop="lessonDate">
              <el-date-picker
                v-model="actionForm.lessonDate"
                class="full-width"
                type="date"
                value-format="YYYY-MM-DD"
                placeholder="请选择日期"
              />
            </el-form-item>

            <el-form-item label="新教室或地点" prop="classroom">
              <el-input v-model="actionForm.classroom" placeholder="请输入新教室或地点" />
            </el-form-item>

            <el-form-item label="开始时间" prop="startTime">
              <el-time-picker
                v-model="actionForm.startTime"
                class="full-width"
                format="HH:mm"
                value-format="HH:mm"
                placeholder="请选择开始时间"
              />
            </el-form-item>

            <el-form-item label="结束时间" prop="endTime">
              <el-time-picker
                v-model="actionForm.endTime"
                class="full-width"
                format="HH:mm"
                value-format="HH:mm"
                placeholder="请选择结束时间"
              />
            </el-form-item>

            <el-form-item class="full-span" label="说明">
              <el-input
                v-model="actionForm.remark"
                :rows="3"
                maxlength="255"
                show-word-limit
                type="textarea"
                :placeholder="dialogMode === 'makeup' ? '可补充补课原因或家长确认信息' : '可补充调课原因或特别提醒'"
              />
            </el-form-item>
          </div>
        </el-form>
      </template>

      <template #footer>
        <div class="dialog-actions">
          <el-button @click="closeDialog">取消</el-button>
          <el-button type="primary" :icon="dialogMode === 'reschedule' ? RefreshRight : dialogMode === 'edit' ? EditPen : dialogMode === 'makeup' ? Clock : Calendar" :loading="saving" @click="submitDialog">
            确认
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>
