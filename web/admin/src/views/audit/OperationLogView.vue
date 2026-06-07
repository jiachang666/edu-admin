<script setup lang="ts">
import { ElMessage } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import {
  fetchOperationLogList,
  fetchUserList,
  type AccessUser,
  type OperationLog,
} from "../../api/education";

const loading = ref(false);
const logs = ref<OperationLog[]>([]);
const operators = ref<AccessUser[]>([]);

const filters = reactive({
  keyword: "",
  userId: null as number | null,
  module: "",
  action: "",
  dateRange: [] as string[],
});

const filteredLogs = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase();

  return logs.value.filter((log) => {
    const matchesKeyword =
      keyword.length === 0 ||
      [log.userName, log.content, log.targetType, String(log.targetId)]
        .join(" ")
        .toLowerCase()
        .includes(keyword);
    const matchesUser = filters.userId === null || log.userId === filters.userId;
    const matchesModule = filters.module.length === 0 || log.module === filters.module;
    const matchesAction = filters.action.length === 0 || log.action === filters.action;

    return matchesKeyword && matchesUser && matchesModule && matchesAction;
  });
});

const operatorCount = computed(() => {
  return new Set(logs.value.map((log) => log.userId)).size;
});

const moduleCount = computed(() => {
  return new Set(logs.value.map((log) => log.module)).size;
});

const latestLogTime = computed(() => {
  return logs.value[0]?.createdAt ?? "暂无记录";
});

const moduleOptions = computed(() => {
  return Array.from(new Set(logs.value.map((log) => log.module)));
});

const actionOptions = computed(() => {
  return Array.from(new Set(logs.value.map((log) => log.action)));
});

const operatorOptions = computed(() => {
  if (operators.value.length > 0) {
    return operators.value.map((item) => ({
      label: item.displayName || item.username,
      value: item.id,
    }));
  }

  return Array.from(
    new Map(logs.value.map((log) => [log.userId, { label: log.userName, value: log.userId }])).values(),
  );
});

function buildLogQuery() {
  const [dateFrom, dateTo] = filters.dateRange;

  return {
    userId: filters.userId ?? undefined,
    module: filters.module || undefined,
    dateFrom: dateFrom || undefined,
    dateTo: dateTo || undefined,
  };
}

async function loadLogs() {
  loading.value = true;

  try {
    const result = await fetchOperationLogList(buildLogQuery());
    logs.value = result.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("操作记录加载失败");
  } finally {
    loading.value = false;
  }
}

async function loadOperators() {
  try {
    const result = await fetchUserList();
    operators.value = result.list;
  } catch (error) {
    console.error(error);
    operators.value = [];
  }
}

function handleSearch() {
  void loadLogs();
}

function resetFilters() {
  filters.keyword = "";
  filters.userId = null;
  filters.module = "";
  filters.action = "";
  filters.dateRange = [];
  void loadLogs();
}

onMounted(() => {
  void Promise.all([loadLogs(), loadOperators()]);
});
</script>

<template>
  <div class="page-stack">
    <section class="page-hero">
      <div class="page-hero__copy">
        <h2>操作记录</h2>
      </div>

      <div class="metric-strip">
        <article class="metric-tile">
          <span>记录总数</span>
          <strong>{{ logs.length }}</strong>
          <small>当前已经被系统留痕的关键动作</small>
        </article>
        <article class="metric-tile">
          <span>涉及账号</span>
          <strong>{{ operatorCount }}</strong>
          <small>至少有过一次关键操作的账号数量</small>
        </article>
        <article class="metric-tile">
          <span>模块种类</span>
          <strong>{{ moduleCount }}</strong>
          <small>目前已经接入留痕的业务范围</small>
        </article>
        <article class="metric-tile">
          <span>最新记录</span>
          <strong>{{ latestLogTime }}</strong>
          <small>方便快速确认系统最近有没有被操作</small>
        </article>
      </div>
    </section>

    <section class="page-card page-card--table">
      <div class="page-header">
        <h2>操作记录列表</h2>
      </div>

      <div class="page-toolbar">
        <div class="toolbar-filters">
          <el-select
            v-model="filters.userId"
            class="toolbar-field"
            clearable
            filterable
            placeholder="全部操作人"
          >
            <el-option
              v-for="item in operatorOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
          <el-select v-model="filters.module" class="toolbar-field" clearable placeholder="全部模块">
            <el-option
              v-for="module in moduleOptions"
              :key="module"
              :label="module"
              :value="module"
            />
          </el-select>
          <el-date-picker
            v-model="filters.dateRange"
            class="toolbar-field"
            clearable
            type="daterange"
            value-format="YYYY-MM-DD"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
          />
          <el-input
            v-model="filters.keyword"
            class="toolbar-field"
            clearable
            placeholder="搜索操作人、说明或对象编号"
          />
          <el-select v-model="filters.action" class="toolbar-field" clearable placeholder="全部动作">
            <el-option
              v-for="action in actionOptions"
              :key="action"
              :label="action"
              :value="action"
            />
          </el-select>
        </div>

        <div class="toolbar-actions">
          <el-button type="primary" @click="handleSearch">筛选记录</el-button>
          <el-button plain @click="resetFilters">重置筛选</el-button>
        </div>
      </div>

      <div class="data-table-shell">
        <el-table v-loading="loading" :data="filteredLogs" stripe>
          <el-table-column label="时间" prop="createdAt" width="190" />
          <el-table-column label="操作人" width="140">
            <template #default="{ row }">
              <div class="table-primary">
                <strong>{{ row.userName }}</strong>
                <small>ID {{ row.userId }}</small>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="模块" prop="module" width="120" />
          <el-table-column label="动作" prop="action" width="150" />
          <el-table-column label="对象" width="150">
            <template #default="{ row }">
              <span>{{ row.targetType }} #{{ row.targetId }}</span>
            </template>
          </el-table-column>
          <el-table-column label="说明" min-width="360" prop="content" />
        </el-table>
      </div>
    </section>
  </div>
</template>
