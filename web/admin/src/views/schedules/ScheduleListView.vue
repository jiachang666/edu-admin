<script setup lang="ts">
import { ElMessage } from "element-plus";
import { onMounted, ref } from "vue";
import { fetchScheduleList, type Schedule } from "../../api/education";

const loading = ref(false);
const schedules = ref<Schedule[]>([]);

async function loadSchedules() {
  loading.value = true;

  try {
    const result = await fetchScheduleList();
    schedules.value = result.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("排课列表加载失败");
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  void loadSchedules();
});
</script>

<template>
  <div class="page-card">
    <div class="page-header">
      <div>
        <h2>排课管理</h2>
        <p class="soft-text">先把上课时间、教室和签到状态展示出来。</p>
      </div>
      <el-tag type="info">共 {{ schedules.length }} 条</el-tag>
    </div>

    <el-table v-loading="loading" :data="schedules" stripe>
      <el-table-column label="日期" prop="lessonDate" width="120" />
      <el-table-column label="时间" prop="lessonTime" width="140" />
      <el-table-column label="班级" prop="className" min-width="180" />
      <el-table-column label="课程" prop="courseName" width="140" />
      <el-table-column label="老师" prop="teacherName" width="120" />
      <el-table-column label="校区" prop="campus" width="120" />
      <el-table-column label="教室" prop="classroom" width="120" />
      <el-table-column label="签到状态" width="120">
        <template #default="{ row }">
          <el-tag
            :type="row.attendanceStatus === '已完成' ? 'success' : 'warning'"
          >
            {{ row.attendanceStatus }}
          </el-tag>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>
