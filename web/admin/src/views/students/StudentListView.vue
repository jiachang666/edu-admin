<script setup lang="ts">
import { ElMessage } from "element-plus";
import { computed, onMounted, ref } from "vue";
import { fetchStudentList, type Student } from "../../api/education";

const loading = ref(false);
const students = ref<Student[]>([]);

const activeCount = computed(() => {
  return students.value.filter((student) => student.status === "在读").length;
});

const renewalCount = computed(() => {
  return students.value.filter((student) => student.status === "待续费").length;
});

const classCount = computed(() => {
  return new Set(
    students.value.filter((student) => student.classId > 0).map((student) => student.classId),
  ).size;
});

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
  <div class="page-stack">
    <section class="page-hero">
      <div class="page-hero__copy">
        <span class="section-kicker">Student Ledger</span>
        <h2>把学员、家长联系方式和班级归属放在一张更适合日常追踪的台账里。</h2>
        <p>
          当前页面先覆盖最常用字段，方便教务在不切页面的情况下快速确认谁在读、谁待续费、谁已经分到班里。
        </p>
      </div>

      <div class="metric-strip">
        <article class="metric-tile">
          <span>学员总数</span>
          <strong>{{ students.length }}</strong>
          <small>包含当前演示库中的全部学员</small>
        </article>
        <article class="metric-tile">
          <span>在读学员</span>
          <strong>{{ activeCount }}</strong>
          <small>首版默认重点关注在读状态</small>
        </article>
        <article class="metric-tile">
          <span>待续费</span>
          <strong>{{ renewalCount }}</strong>
          <small>适合班主任和教务跟进提醒</small>
        </article>
        <article class="metric-tile">
          <span>已分班级</span>
          <strong>{{ classCount }}</strong>
          <small>当前学员已经分布到的班级数</small>
        </article>
      </div>
    </section>

    <section class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>学员列表</h2>
          <p class="soft-text">覆盖学员、家长、班级和剩余课时这几个最核心字段。</p>
        </div>
        <div class="section-note">学员视图</div>
      </div>

      <div class="data-table-shell">
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
    </section>
  </div>
</template>
