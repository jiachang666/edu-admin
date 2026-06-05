<script setup lang="ts">
import { ElMessage } from "element-plus";
import { onMounted, ref } from "vue";
import { fetchStudentList, type Student } from "../../api/education";

const loading = ref(false);
const students = ref<Student[]>([]);

async function loadStudents() {
  loading.value = true;

  try {
    const result = await fetchStudentList();
    students.value = result.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("学员列表加载失败");
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  void loadStudents();
});
</script>

<template>
  <div class="page-card">
    <div class="page-header">
      <div>
        <h2>学员管理</h2>
        <p class="soft-text">覆盖学员、家长、班级和剩余课时这几个最核心字段。</p>
      </div>
      <el-tag type="info">共 {{ students.length }} 人</el-tag>
    </div>

    <el-table v-loading="loading" :data="students" stripe>
      <el-table-column label="姓名" prop="name" width="120" />
      <el-table-column label="年级" prop="grade" width="100" />
      <el-table-column label="所属班级" prop="className" min-width="180" />
      <el-table-column label="家长" prop="parentName" width="120" />
      <el-table-column label="联系电话" prop="parentMobile" width="140" />
      <el-table-column label="剩余课时" prop="remainingHours" width="100" />
      <el-table-column label="校区" prop="campus" width="120" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === '在读' ? 'success' : 'warning'">
            {{ row.status }}
          </el-tag>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>
