<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import {
  createUser,
  disableUser,
  enableUser,
  fetchRoleList,
  fetchUserList,
  updateUser,
  type AccessRole,
  type AccessUser,
  type AccessUserPayload,
} from "../../api/education";
import { useAuthStore } from "../../stores/auth";

const authStore = useAuthStore();

const loading = ref(false);
const saving = ref(false);
const dialogVisible = ref(false);
const editingUserId = ref<number | null>(null);
const users = ref<AccessUser[]>([]);
const roles = ref<AccessRole[]>([]);
const formRef = ref<FormInstance>();

const filters = reactive({
  keyword: "",
  status: "",
  roleCode: "",
});

const form = reactive<AccessUserPayload>(defaultForm());

const rules: FormRules<AccessUserPayload> = {
  username: [{ required: true, message: "请输入登录账号", trigger: "blur" }],
  displayName: [{ required: true, message: "请输入姓名", trigger: "blur" }],
  roleCode: [{ required: true, message: "请选择角色", trigger: "change" }],
  status: [{ required: true, message: "请选择状态", trigger: "change" }],
};

const filteredUsers = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase();

  return users.value.filter((user) => {
    const matchesKeyword =
      keyword.length === 0 ||
      [user.username, user.displayName, user.mobile, user.primaryRoleName]
        .join(" ")
        .toLowerCase()
        .includes(keyword);
    const matchesStatus = filters.status.length === 0 || user.status === filters.status;
    const matchesRole = filters.roleCode.length === 0 || user.primaryRoleCode === filters.roleCode;

    return matchesKeyword && matchesStatus && matchesRole;
  });
});

const enabledCount = computed(() => {
  return users.value.filter((user) => user.status === "启用").length;
});

const disabledCount = computed(() => {
  return users.value.filter((user) => user.status === "停用").length;
});

const recentLoginCount = computed(() => {
  return users.value.filter((user) => user.lastLoginAt.length > 0).length;
});

const dialogTitle = computed(() => {
  return editingUserId.value ? "编辑账号" : "新建账号";
});

const canManageUsers = computed(() => authStore.hasPermission("users:manage"));

function defaultForm(): AccessUserPayload {
  return {
    username: "",
    password: "",
    displayName: "",
    mobile: "",
    roleCode: "",
    status: "启用",
  };
}

function resetForm() {
  Object.assign(form, defaultForm());
  editingUserId.value = null;
  formRef.value?.clearValidate();
}

function openCreateDialog() {
  resetForm();
  dialogVisible.value = true;
}

function openEditDialog(user: AccessUser) {
  editingUserId.value = user.id;
  Object.assign(form, {
    username: user.username,
    password: "",
    displayName: user.displayName,
    mobile: user.mobile,
    roleCode: user.primaryRoleCode,
    status: user.status,
  });
  dialogVisible.value = true;
  formRef.value?.clearValidate();
}

function closeDialog() {
  dialogVisible.value = false;
  resetForm();
}

function buildPayload() {
  return {
    username: form.username.trim(),
    password: form.password.trim(),
    displayName: form.displayName.trim(),
    mobile: form.mobile.trim(),
    roleCode: form.roleCode,
    status: form.status,
  };
}

async function loadPageData() {
  loading.value = true;

  try {
    const [userResult, roleResult] = await Promise.all([fetchUserList(), fetchRoleList()]);
    users.value = userResult.list;
    roles.value = roleResult.list.filter((role) => role.status === "启用");
  } catch (error) {
    console.error(error);
    ElMessage.error("账号工作台加载失败");
  } finally {
    loading.value = false;
  }
}

async function submitForm() {
  const formNode = formRef.value;
  if (!formNode) {
    return;
  }

  const valid = await formNode.validate().catch(() => false);
  if (!valid) {
    return;
  }

  if (!editingUserId.value && form.password.trim().length === 0) {
    ElMessage.warning("新建账号时需要填写密码");
    return;
  }

  saving.value = true;

  try {
    const payload = buildPayload();

    if (editingUserId.value) {
      await updateUser(editingUserId.value, payload);
      ElMessage.success("账号已更新");
    } else {
      await createUser(payload);
      ElMessage.success("账号已创建");
    }

    closeDialog();
    await loadPageData();
  } catch (error: any) {
    console.error(error);
    const message = error?.response?.data?.message ?? "账号保存失败";
    ElMessage.error(message);
  } finally {
    saving.value = false;
  }
}

async function toggleUserStatus(user: AccessUser) {
  try {
    if (user.status === "启用") {
      await disableUser(user.id);
      ElMessage.success("账号已停用");
    } else {
      await enableUser(user.id);
      ElMessage.success("账号已启用");
    }

    await loadPageData();
  } catch (error: any) {
    console.error(error);
    const message = error?.response?.data?.message ?? "账号状态更新失败";
    ElMessage.error(message);
  }
}

function resetFilters() {
  filters.keyword = "";
  filters.status = "";
  filters.roleCode = "";
}

onMounted(() => {
  void loadPageData();
});
</script>

<template>
  <div class="page-stack">
    <section class="page-card page-card--table list-card">
      <div class="page-header">
        <div class="list-card__heading">
          <h2>账号列表</h2>
          <span class="list-card__count">共 {{ filteredUsers.length }} 条</span>
        </div>
        <div class="page-actions">
          <el-button v-if="canManageUsers" type="primary" @click="openCreateDialog">
            新建账号
          </el-button>
        </div>
      </div>

      <div class="metric-strip metric-strip--compact list-card__metrics">
        <article class="metric-tile">
          <span>账号总数</span>
          <strong>{{ users.length }}</strong>
        </article>
        <article class="metric-tile">
          <span>启用账号</span>
          <strong>{{ enabledCount }}</strong>
        </article>
        <article class="metric-tile">
          <span>停用账号</span>
          <strong>{{ disabledCount }}</strong>
        </article>
        <article class="metric-tile">
          <span>有登录记录</span>
          <strong>{{ recentLoginCount }}</strong>
        </article>
      </div>

      <div class="filter-bar list-card__filters">
        <div class="toolbar-filters">
          <el-input v-model="filters.keyword" class="toolbar-field" clearable placeholder="搜索账号、姓名或手机号" />
          <el-select v-model="filters.status" class="toolbar-field" clearable placeholder="全部状态">
            <el-option label="启用" value="启用" />
            <el-option label="停用" value="停用" />
          </el-select>
          <el-select v-model="filters.roleCode" class="toolbar-field" clearable placeholder="全部角色">
            <el-option
              v-for="role in roles"
              :key="role.id"
              :label="role.name"
              :value="role.code"
            />
          </el-select>
        </div>

        <div class="toolbar-actions">
          <el-button plain @click="resetFilters">重置筛选</el-button>
        </div>
      </div>

      <div class="data-table-shell">
        <el-table v-loading="loading" :data="filteredUsers" stripe>
          <el-table-column label="账号" min-width="170">
            <template #default="{ row }">
              <div class="table-primary">
                <strong>{{ row.username }}</strong>
                <small>{{ row.displayName }}</small>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="角色" width="160">
            <template #default="{ row }">
              <el-tag>{{ row.primaryRoleName || "未分配" }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="手机号" prop="mobile" width="150" />
          <el-table-column label="最近登录" prop="lastLoginAt" width="190">
            <template #default="{ row }">
              <span class="muted-cell">{{ row.lastLoginAt || "还没有记录" }}</span>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="110">
            <template #default="{ row }">
              <el-tag :type="row.status === '启用' ? 'success' : 'info'">
                {{ row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column v-if="canManageUsers" label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <div class="table-link-group">
                <el-button link type="primary" @click="openEditDialog(row)">编辑</el-button>
                <el-button link type="primary" @click="toggleUserStatus(row)">
                  {{ row.status === "启用" ? "停用" : "启用" }}
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </section>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="680px" @closed="closeDialog">
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <div class="dialog-grid">
          <el-form-item label="登录账号" prop="username">
            <el-input v-model="form.username" placeholder="例如 staff01" />
          </el-form-item>
          <el-form-item label="姓名" prop="displayName">
            <el-input v-model="form.displayName" placeholder="例如 教务小林" />
          </el-form-item>
          <el-form-item label="手机号">
            <el-input v-model="form.mobile" placeholder="选填" />
          </el-form-item>
          <el-form-item label="角色" prop="roleCode">
            <el-select v-model="form.roleCode" placeholder="请选择角色">
              <el-option
                v-for="role in roles"
                :key="role.id"
                :label="role.name"
                :value="role.code"
              />
            </el-select>
          </el-form-item>
          <el-form-item class="full-span" :label="editingUserId ? '重置密码（留空则不修改）' : '登录密码'">
            <el-input
              v-model="form.password"
              :placeholder="editingUserId ? '如需修改，再输入新密码' : '请输入初始密码'"
              show-password
              type="password"
            />
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-radio-group v-model="form.status">
              <el-radio label="启用">启用</el-radio>
              <el-radio label="停用">停用</el-radio>
            </el-radio-group>
          </el-form-item>
        </div>
      </el-form>

      <template #footer>
        <div class="dialog-actions">
          <el-button @click="closeDialog">取消</el-button>
          <el-button :loading="saving" type="primary" @click="submitForm">保存账号</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>
