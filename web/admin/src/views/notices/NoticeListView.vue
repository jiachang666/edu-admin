<script setup lang="ts">
import { ElMessage } from "element-plus";
import { computed, onMounted, ref } from "vue";
import { fetchNoticeList, type Notice } from "../../api/education";

const loading = ref(false);
const notices = ref<Notice[]>([]);

const sentCount = computed(() => {
  return notices.value.filter((notice) => notice.status === "已发送").length;
});

const pendingCount = computed(() => {
  return notices.value.filter((notice) => notice.status === "待发送").length;
});

const draftCount = computed(() => {
  return notices.value.filter((notice) => notice.status === "草稿").length;
});

async function loadNotices() {
  loading.value = true;

  try {
    const result = await fetchNoticeList();
    notices.value = result.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("通知列表加载失败");
  } finally {
    loading.value = false;
  }
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
        <h2>把通知的分类、范围和发送进度统一放进一张更好回看的消息派发板。</h2>
        <p>
          当前页面适合运营和班主任先确认哪些通知已经发出、哪些仍在草稿或待发送状态，避免临时消息遗落。
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
          <span>草稿</span>
          <strong>{{ draftCount }}</strong>
          <small>还在编辑中的消息内容</small>
        </article>
      </div>
    </section>

    <section class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>通知列表</h2>
          <p class="soft-text">先把通知分类、发送范围和当前状态跑通。</p>
        </div>
        <div class="section-note">通知视图</div>
      </div>

      <div class="data-table-shell">
        <el-table v-loading="loading" :data="notices" stripe>
          <el-table-column label="标题" prop="title" min-width="220" />
          <el-table-column label="分类" prop="category" width="140" />
          <el-table-column label="范围" prop="targetScope" min-width="180" />
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="row.status === '已发送' ? 'success' : 'warning'">
                {{ row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="时间" prop="publishAt" width="180" />
          <el-table-column label="发起人" prop="author" width="120" />
        </el-table>
      </div>
    </section>
  </div>
</template>
