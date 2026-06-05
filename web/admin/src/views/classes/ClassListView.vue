<script setup lang="ts">
import { ElMessage } from "element-plus";
import { computed, onMounted, ref } from "vue";
import { fetchClassList, type SchoolClass } from "../../api/education";

const loading = ref(false);
const classes = ref<SchoolClass[]>([]);

const runningCount = computed(() => {
  return classes.value.filter((item) => item.status === "开班中").length;
});

const remainingSeats = computed(() => {
  return classes.value.reduce((total, item) => total + (item.capacity - item.studentCount), 0);
});

const campusCount = computed(() => {
  return new Set(classes.value.map((item) => item.campus)).size;
});

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
  <div class="page-stack">
    <section class="page-hero">
      <div class="page-hero__copy">
        <span class="section-kicker">Classroom Matrix</span>
        <h2>把班级、课程、老师和容量安排放进同一个稳定的班级矩阵里。</h2>
        <p>
          这页更适合教务看整体开班情况：哪些班已开、还有多少余量、目前覆盖了几个校区，后面再接新增和编辑流程会更顺。
        </p>
      </div>

      <div class="metric-strip">
        <article class="metric-tile">
          <span>班级总数</span>
          <strong>{{ classes.length }}</strong>
          <small>当前已经建立的全部班级</small>
        </article>
        <article class="metric-tile">
          <span>开班中</span>
          <strong>{{ runningCount }}</strong>
          <small>正在正常运转的班级数量</small>
        </article>
        <article class="metric-tile">
          <span>剩余名额</span>
          <strong>{{ remainingSeats }}</strong>
          <small>按容量减去当前人数得到</small>
        </article>
        <article class="metric-tile">
          <span>覆盖校区</span>
          <strong>{{ campusCount }}</strong>
          <small>方便看班级分布是否均衡</small>
        </article>
      </div>
    </section>

    <section class="page-card page-card--table">
      <div class="page-header">
        <div>
          <h2>班级列表</h2>
          <p class="soft-text">先把班级、课程、老师和容量关系串起来。</p>
        </div>
        <div class="section-note">班级视图</div>
      </div>

      <div class="data-table-shell">
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
    </section>
  </div>
</template>
