<script setup lang="ts">
import { ElMessage } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { fetchClassList, type SchoolClass } from "../../api/education";

const router = useRouter();

const loading = ref(false);
const classes = ref<SchoolClass[]>([]);

const filters = reactive({
  keyword: "",
  status: "",
  teacher: "",
});

const filteredClasses = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase();

  return classes.value.filter((item) => {
    const matchesKeyword =
      keyword.length === 0 ||
      [item.name, item.courseName, item.teacherName, item.campus]
        .join(" ")
        .toLowerCase()
        .includes(keyword);
    const matchesStatus = filters.status.length === 0 || item.status === filters.status;
    const matchesTeacher = filters.teacher.length === 0 || item.teacherName === filters.teacher;

    return matchesKeyword && matchesStatus && matchesTeacher;
  });
});

const runningCount = computed(() => {
  return classes.value.filter((item) => item.status === "开班中").length;
});

const remainingSeats = computed(() => {
  return classes.value.reduce((total, item) => total + (item.capacity - item.studentCount), 0);
});

const campusCount = computed(() => {
  return new Set(classes.value.map((item) => item.campus)).size;
});

const teacherOptions = computed(() => {
  return Array.from(new Set(classes.value.map((item) => item.teacherName).filter(Boolean)));
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

function handleResetFilters() {
  filters.keyword = "";
  filters.status = "";
  filters.teacher = "";
}

function openDetail(classId: number) {
  void router.push(`/classes/${classId}`);
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
        <h2>先从班级列表看整体，再进入某一个班，把学员、排课和课后动作连到一起。</h2>
        <p>
          这页现在不仅能看班级矩阵，也能直接进入班级详情中心页。后面老师和教务的大多数动作，都可以从那个入口继续往下走。
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
          <p class="soft-text">先找到目标班级，再进入详情页处理学员、排课、签到和通知。</p>
        </div>
        <div class="section-note">班级视图</div>
      </div>

      <div class="page-toolbar">
        <div class="toolbar-filters">
          <el-input
            v-model="filters.keyword"
            class="toolbar-field"
            clearable
            placeholder="搜索班级、课程、老师或校区"
          />
          <el-select
            v-model="filters.teacher"
            class="toolbar-field"
            clearable
            placeholder="主讲老师"
          >
            <el-option
              v-for="teacher in teacherOptions"
              :key="teacher"
              :label="teacher"
              :value="teacher"
            />
          </el-select>
          <el-select
            v-model="filters.status"
            class="toolbar-field"
            clearable
            placeholder="班级状态"
          >
            <el-option label="开班中" value="开班中" />
            <el-option label="待满班" value="待满班" />
          </el-select>
        </div>
        <el-button @click="handleResetFilters">重置筛选</el-button>
      </div>

      <div class="data-table-shell">
        <el-table v-loading="loading" :data="filteredClasses" stripe>
          <el-table-column label="班级名称" min-width="200">
            <template #default="{ row }">
              <el-button link type="primary" @click="openDetail(row.id)">{{ row.name }}</el-button>
            </template>
          </el-table-column>
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
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="openDetail(row.id)">进入详情</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </section>
  </div>
</template>
