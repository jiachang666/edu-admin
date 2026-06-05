<script setup lang="ts">
import { ElMessage } from "element-plus";
import { onMounted, ref } from "vue";
import { fetchNoticeList, type Notice } from "../../api/education";

const loading = ref(false);
const notices = ref<Notice[]>([]);

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
  <div class="page-card">
    <div class="page-header">
      <div>
        <h2>通知管理</h2>
        <p class="soft-text">先把通知分类、发送范围和当前状态跑通。</p>
      </div>
      <el-tag type="info">共 {{ notices.length }} 条</el-tag>
    </div>

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
</template>
