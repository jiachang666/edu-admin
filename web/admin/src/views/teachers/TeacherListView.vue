<script setup lang="ts">
import { ElMessage } from "element-plus";
import { computed, onMounted, ref } from "vue";
import { fetchTeacherList, type Teacher } from "../../api/education";

const loading = ref(false);
const teachers = ref<Teacher[]>([]);

const activeCount = computed(() => {
  return teachers.value.filter((teacher) => teacher.status === "在职").length;
});

const fullTimeCount = computed(() => {
  return teachers.value.filter((teacher) => teacher.employmentType === "全职").length;
});

const campusCount = computed(() => {
  return new Set(teachers.value.map((teacher) => teacher.campus)).size;
});

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
  <div class="page-stack">
    <section class="page-hero">
      <div class="page-hero__copy">
        <span class="section-kicker">Faculty Ledger</span>
        <h2>把老师的科目、校区和授课负载放进同一份清楚的师资台账里。</h2>
        <p>
          当前这页更偏向教务和负责人快速盘点：谁在职、谁是全职、目前覆盖了几个校区，先把基础视图稳定住。
        </p>
      </div>

      <div class="metric-strip">
        <article class="metric-tile">
          <span>老师总数</span>
          <strong>{{ teachers.length }}</strong>
          <small>当前收录在系统里的全部师资</small>
        </article>
        <article class="metric-tile">
          <span>在职人数</span>
          <strong>{{ activeCount }}</strong>
          <small>默认更关注当前可排课老师</small>
        </article>
        <article class="metric-tile">
          <span>全职老师</span>
          <strong>{{ fullTimeCount }}</strong>
          <small>便于快速区分排班主力</small>
        </article>
        <article class="metric-tile">
          <span>覆盖校区</span>
          <strong>{{ campusCount }}</strong>
          <small>看清老师分布是否均衡</small>
        </article>
      </div>
    </section>

    <section class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>老师列表</h2>
          <p class="soft-text">先展示老师台账，后面再继续补详情和编辑流程。</p>
        </div>
        <div class="section-note">师资视图</div>
      </div>

      <div class="data-table-shell">
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
    </section>
  </div>
</template>
