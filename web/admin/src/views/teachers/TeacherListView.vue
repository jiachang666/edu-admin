<script setup lang="ts">
import { ElMessage } from "element-plus";
import { onMounted, ref } from "vue";
import { fetchTeacherList, type Teacher } from "../../api/education";

const loading = ref(false);
const teachers = ref<Teacher[]>([]);

async function loadTeachers() {
  loading.value = true;

  try {
    const result = await fetchTeacherList();
    teachers.value = result.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("老师列表加载失败");
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  void loadTeachers();
});
</script>

<template>
  <div class="page-card">
    <div class="page-header">
      <div>
        <h2>老师管理</h2>
        <p class="soft-text">先展示老师台账，后面再继续补详情和编辑流程。</p>
      </div>
      <el-tag type="info">共 {{ teachers.length }} 位</el-tag>
    </div>

    <el-table v-loading="loading" :data="teachers" stripe>
      <el-table-column label="姓名" prop="name" width="120" />
      <el-table-column label="主教科目" prop="mainSubject" width="140" />
      <el-table-column label="校区" prop="campus" width="120" />
      <el-table-column label="类型" prop="employmentType" width="100" />
      <el-table-column label="周课时" prop="weeklyHours" width="100" />
      <el-table-column label="手机号" prop="mobile" width="140" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === '在职' ? 'success' : 'warning'">
            {{ row.status }}
          </el-tag>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>
