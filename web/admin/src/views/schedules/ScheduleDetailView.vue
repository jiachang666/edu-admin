<script setup lang="ts">
import { Bell, ChatDotRound, Finished, Reading } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { fetchScheduleDetail, type ScheduleDetail } from "../../api/education";

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const detail = ref<ScheduleDetail | null>(null);

const scheduleId = computed(() => {
  const parsedValue = Number(route.params.id);
  if (!Number.isInteger(parsedValue) || parsedValue <= 0) {
    return null;
  }

  return parsedValue;
});

const currentSchedule = computed(() => {
  return detail.value?.schedule ?? null;
});

const currentClass = computed(() => {
  return detail.value?.class ?? null;
});

const attendanceSummary = computed(() => {
  return detail.value?.attendance ?? null;
});

function statusTone(status: string) {
  switch (status) {
    case "已完成":
      return "success";
    case "待签到":
      return "warning";
    case "已停课":
      return "danger";
    case "已调课":
      return "primary";
    default:
      return "info";
  }
}

function homeworkStatusText() {
  if (!detail.value?.homework) {
    return "未布置";
  }

  return detail.value.homework.status === "published" ? "已发布" : "草稿";
}

function feedbackSummaryText() {
  if (!detail.value?.feedback) {
    return "老师还没写课后反馈";
  }
  if (detail.value.feedback.summary.trim()) {
    return detail.value.feedback.summary;
  }

  return "老师还没写课堂总结";
}

async function loadDetail() {
  const currentScheduleId = scheduleId.value;
  if (currentScheduleId === null) {
    ElMessage.error("课程安排编号不正确");
    return;
  }

  loading.value = true;

  try {
    detail.value = await fetchScheduleDetail(currentScheduleId);
  } catch (error) {
    console.error(error);
    ElMessage.error("课程安排详情加载失败");
  } finally {
    loading.value = false;
  }
}

function openAttendance() {
  if (!currentSchedule.value) {
    return;
  }

  void router.push({
    path: "/attendance",
    query: { scheduleId: String(currentSchedule.value.id) },
  });
}

function openHomework() {
  if (!currentSchedule.value) {
    return;
  }

  void router.push({
    path: "/homeworks",
    query: { scheduleId: String(currentSchedule.value.id) },
  });
}

function openClassDetail() {
  if (!currentClass.value) {
    return;
  }

  void router.push(`/classes/${currentClass.value.id}`);
}

function openSchedules() {
  void router.push("/schedules");
}

function openNotices() {
  void router.push("/notices");
}

onMounted(() => {
  void loadDetail();
});
</script>

<template>
  <div v-loading="loading" class="page-stack">
    <section class="page-hero">
      <div class="page-hero__copy">
        <span class="section-kicker">
          <el-icon><Reading /></el-icon>
          Lesson Hub
        </span>
        <h2>{{ currentSchedule?.className || "课程安排详情" }}</h2>
        <p>
          把这一次具体上课的签到、作业、反馈和相关通知都放到一起，方便老师和教务沿着同一条线接力处理。
        </p>
      </div>

      <div class="metric-strip">
        <article class="metric-tile">
          <span>当前状态</span>
          <strong>{{ currentSchedule?.attendanceStatus || "-" }}</strong>
          <small>先判断这节课现在卡在哪个环节</small>
        </article>
        <article class="metric-tile">
          <span>到课人数</span>
          <strong>{{ attendanceSummary?.presentCount ?? 0 }}</strong>
          <small>已经确认到课或补签的人数</small>
        </article>
        <article class="metric-tile">
          <span>异常人数</span>
          <strong>{{ (attendanceSummary?.leaveCount ?? 0) + (attendanceSummary?.absentCount ?? 0) }}</strong>
          <small>请假和缺席一起先盯住</small>
        </article>
        <article class="metric-tile">
          <span>关联通知</span>
          <strong>{{ detail?.relatedNotices.length ?? 0 }}</strong>
          <small>和这个班最近有关的提醒与消息</small>
        </article>
      </div>
    </section>

    <section class="page-card">
      <div class="page-header">
        <div>
          <h2>本次课程信息</h2>
          <p class="soft-text">先把这节课什么时候上、谁来上、在哪上看清楚。</p>
        </div>
        <div class="section-note">单次课程</div>
      </div>

      <div class="detail-info-grid">
        <article class="detail-info-card">
          <span>班级</span>
          <strong>{{ currentSchedule?.className || "-" }}</strong>
        </article>
        <article class="detail-info-card">
          <span>课程</span>
          <strong>{{ currentSchedule?.courseName || "-" }}</strong>
        </article>
        <article class="detail-info-card">
          <span>老师</span>
          <strong>{{ currentSchedule?.teacherName || "-" }}</strong>
        </article>
        <article class="detail-info-card">
          <span>日期</span>
          <strong>{{ currentSchedule?.lessonDate || "-" }}</strong>
        </article>
        <article class="detail-info-card">
          <span>时间</span>
          <strong>{{ currentSchedule?.lessonTime || "-" }}</strong>
        </article>
        <article class="detail-info-card">
          <span>校区 / 教室</span>
          <strong>{{ currentSchedule?.campus || "-" }} / {{ currentSchedule?.classroom || "-" }}</strong>
        </article>
      </div>

      <div class="detail-callout" v-if="currentSchedule?.remark">
        <strong>备注</strong>
        <p>{{ currentSchedule.remark }}</p>
      </div>
    </section>

    <div class="class-detail-grid">
      <section class="page-card">
        <div class="page-header">
          <div>
            <h3>签到概况</h3>
            <p class="soft-text">先看这节课还有多少人没确认，是否已经完成签到。</p>
          </div>
          <div class="section-note">Attendance</div>
        </div>

        <div class="stats-grid stats-grid--compact">
          <article class="stat-card" data-tone="green">
            <span class="stat-label">已到</span>
            <strong class="stat-value">{{ attendanceSummary?.presentCount ?? 0 }}</strong>
          </article>
          <article class="stat-card" data-tone="orange">
            <span class="stat-label">请假</span>
            <strong class="stat-value">{{ attendanceSummary?.leaveCount ?? 0 }}</strong>
          </article>
          <article class="stat-card" data-tone="red">
            <span class="stat-label">缺席</span>
            <strong class="stat-value">{{ attendanceSummary?.absentCount ?? 0 }}</strong>
          </article>
          <article class="stat-card" data-tone="indigo">
            <span class="stat-label">待确认</span>
            <strong class="stat-value">{{ attendanceSummary?.pendingCount ?? 0 }}</strong>
          </article>
        </div>

        <div class="stack-list">
          <article class="stack-item stack-item--stretch">
            <div>
              <strong>{{ currentSchedule?.lessonDate }} {{ currentSchedule?.lessonTime }}</strong>
              <small>
                共 {{ attendanceSummary?.studentCount ?? 0 }} 人 ·
                <el-tag :type="statusTone(attendanceSummary?.attendanceStatus || '')">
                  {{ attendanceSummary?.attendanceStatus || "待处理" }}
                </el-tag>
              </small>
            </div>
            <el-button type="primary" @click="openAttendance">
              {{ attendanceSummary?.attendanceStatus === "待签到" ? "去签到" : "查看签到" }}
            </el-button>
          </article>
        </div>
      </section>

      <section class="page-card">
        <div class="page-header">
          <div>
            <h3>课后内容</h3>
            <p class="soft-text">作业和反馈都围绕这一次课来处理，不再分散在别处。</p>
          </div>
          <div class="section-note">Homework</div>
        </div>

        <div class="stack-list">
          <article class="stack-item stack-item--stretch">
            <div>
              <strong>{{ detail?.homework?.title || "这节课还没有作业标题" }}</strong>
              <small>作业状态：{{ homeworkStatusText() }}</small>
            </div>
            <el-button link type="primary" @click="openHomework">查看作业</el-button>
          </article>

          <article class="stack-item stack-item--stretch">
            <div>
              <strong>课堂反馈</strong>
              <small>{{ feedbackSummaryText() }}</small>
            </div>
            <el-button link type="primary" @click="openHomework">查看反馈</el-button>
          </article>
        </div>
      </section>
    </div>

    <div class="class-detail-grid">
      <section class="page-card">
        <div class="page-header">
          <div>
            <h3>到课学员</h3>
            <p class="soft-text">先确认这节课涉及哪些学员，方便后面跟签到和通知。</p>
          </div>
          <div class="section-note">学生名单</div>
        </div>

        <div class="stack-list">
          <article
            v-for="student in detail?.students ?? []"
            :key="student.id"
            class="stack-item stack-item--stretch"
          >
            <div>
              <strong>{{ student.name }}</strong>
              <small>{{ student.grade }} · {{ student.parentName }} · {{ student.parentMobile }}</small>
            </div>
          </article>

          <div v-if="(detail?.students.length ?? 0) === 0" class="soft-empty">
            这节课暂时没有关联学员。
          </div>
        </div>
      </section>

      <section class="page-card">
        <div class="page-header">
          <div>
            <h3>相关通知</h3>
            <p class="soft-text">这里放和这个班近期有关的通知，方便继续发消息或回看历史。</p>
          </div>
          <div class="section-note">Notice</div>
        </div>

        <div class="stack-list">
          <article
            v-for="notice in detail?.relatedNotices ?? []"
            :key="notice.id"
            class="stack-item stack-item--stretch"
          >
            <div>
              <strong>{{ notice.title }}</strong>
              <small>{{ notice.category }} · {{ notice.status }} · {{ notice.publishAt || "未发送" }}</small>
            </div>
            <el-button link type="primary" @click="openNotices">去通知页</el-button>
          </article>

          <div v-if="(detail?.relatedNotices.length ?? 0) === 0" class="soft-empty">
            这节课目前还没有关联通知。
          </div>
        </div>
      </section>
    </div>

    <section class="page-card">
      <div class="page-header">
        <div>
          <h2>快捷入口</h2>
          <p class="soft-text">把这节课最常继续走的几个方向放在一起，减少来回切页。</p>
        </div>
        <div class="section-note">继续处理</div>
      </div>

      <div class="quick-actions-grid">
        <button class="quick-action-card" type="button" @click="openAttendance">
          <el-icon><Finished /></el-icon>
          <strong>处理签到</strong>
          <span>继续去点名、补签或回看这节课的到课情况。</span>
        </button>
        <button class="quick-action-card" type="button" @click="openHomework">
          <el-icon><ChatDotRound /></el-icon>
          <strong>处理作业反馈</strong>
          <span>直接继续补作业内容、课堂总结和家长提醒。</span>
        </button>
        <button class="quick-action-card" type="button" @click="openClassDetail">
          <el-icon><Reading /></el-icon>
          <strong>返回班级中心</strong>
          <span>回到班级页统一看这个班后续课程、学员和通知。</span>
        </button>
        <button class="quick-action-card" type="button" @click="openNotices">
          <el-icon><Bell /></el-icon>
          <strong>继续发通知</strong>
          <span>需要通知家长时，直接去消息工作台整理和发送。</span>
        </button>
      </div>

      <div class="stack-list stack-list--spacious">
        <article class="stack-item stack-item--stretch">
          <div>
            <strong>返回排课总览</strong>
            <small>如果还要继续改别的课程时间、停课或补课，可以回到排课工作台。</small>
          </div>
          <el-button link type="primary" @click="openSchedules">回排课页</el-button>
        </article>
      </div>
    </section>
  </div>
</template>
