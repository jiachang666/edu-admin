<script setup lang="ts">
import { ElMessage } from "element-plus";
import { onMounted, ref } from "vue";
import { fetchClassList, type SchoolClass } from "../../api/education";

const loading = ref(false);
const classes = ref<SchoolClass[]>([]);

async function loadClasses() {
  loading.value = true;

  try {
    const result = await fetchClassList();
    classes.value = result.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("班级列表加载失败");
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  void loadClasses();
});
</script>

<template>
  <div class="page-card">
    <div class="page-header">
      <div>
        <h2>班级管理</h2>
        <p class="soft-text">先把班级、课程、老师和容量关系串起来。</p>
      </div>
      <el-tag type="info">共 {{ classes.length }} 个班</el-tag>
    </div>

    <el-table v-loading="loading" :data="classes" stripe>
      <el-table-column label="班级名称" prop="name" min-width="180" />
      <el-table-column label="课程" prop="courseName" width="140" />
      <el-table-column label="老师" prop="teacherName" width="120" />
      <el-table-column label="校区" prop="campus" width="120" />
      <el-table-column label="人数" width="100">
        <template #default="{ row }">
          {{ row.studentCount }}/{{ row.capacity }}
        </template>
      </el-table-column>
      <el-table-column label="固定排课" prop="weeklySchedule" min-width="180" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === '开班中' ? 'success' : 'warning'">
            {{ row.status }}
          </el-tag>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>
