<script setup lang="ts">
import { ElMessage } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import { fetchOperationLogList, type OperationLog } from "../../api/education";

const loading = ref(false);
const logs = ref<OperationLog[]>([]);

const filters = reactive({
  keyword: "",
  module: "",
  action: "",
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
    const matchesModule = filters.module.length === 0 || log.module === filters.module;
    const matchesAction = filters.action.length === 0 || log.action === filters.action;

    return matchesKeyword && matchesModule && matchesAction;
  });
});

const operatorCount = computed(() => {
  return new Set(logs.value.map((log) => log.userId)).size;
});

const moduleCount = computed(() => {
  return new Set(logs.value.map((log) => log.module)).size;
});

const loginCount = computed(() => {
  return logs.value.filter((log) => log.action === "login").length;
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

async function loadLogs() {
  loading.value = true;

  try {
    const result = await fetchOperationLogList();
    logs.value = result.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("操作记录加载失败");
  } finally {
    loading.value = false;
  }
}

function resetFilters() {
  filters.keyword = "";
  filters.module = "";
  filters.action = "";
}

onMounted(() => {
  void loadLogs();
});
</script>

<template>
  <div class="page-stack">
    <section class="page-hero">
      <div class="page-hero__copy">
        <span class="section-kicker">Audit Trail</span>
        <h2>把关键动作留下来，后面追问题时才能清楚知道是谁在什么时候做了什么。</h2>
        <p>
          首版先覆盖账号、角色和登录这几类关键动作，先把“有迹可查”这条底线建立起来。后面你还可以继续把排课、通知这些业务操作补进来。
        </p>
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
        <div>
          <h2>操作记录列表</h2>
          <p class="soft-text">可以按模块、动作和关键词筛选，适合定位“谁做了什么”的问题。</p>
        </div>
        <div class="section-note">关键留痕</div>
      </div>

      <div class="page-toolbar">
        <div class="toolbar-filters">
          <el-input v-model="filters.keyword" class="toolbar-field" clearable placeholder="搜索操作人、说明或对象编号" />
          <el-select v-model="filters.module" class="toolbar-field" clearable placeholder="全部模块">
            <el-option
              v-for="module in moduleOptions"
              :key="module"
              :label="module"
              :value="module"
            />
          </el-select>
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

      <div class="detail-note">
        <strong>当前已接入的留痕重点</strong>
        <p>
          登录、账号新建编辑启停、角色新建编辑、角色权限保存，这几类动作现在都会自动进入记录里。后续可以继续把排课、停课、补课和通知发送也接进来。
        </p>
      </div>
    </section>
  </div>
</template>
