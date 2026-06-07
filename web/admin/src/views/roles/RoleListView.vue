<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import {
  createRole,
  fetchRoleDetail,
  fetchRoleList,
  saveRolePermissions,
  updateRole,
  type AccessRole,
  type AccessRolePayload,
} from "../../api/education";
import { useAuthStore } from "../../stores/auth";

type PermissionGroup = {
  key: string;
  label: string;
  description: string;
  permissions: Array<{ code: string; label: string; description: string }>;
};

const authStore = useAuthStore();

const loading = ref(false);
const saving = ref(false);
const roleDialogVisible = ref(false);
const permissionDialogVisible = ref(false);
const editingRoleId = ref<number | null>(null);
const roles = ref<AccessRole[]>([]);
const permissionGroups = ref<PermissionGroup[]>([]);
const selectedPermissions = ref<string[]>([]);
const activePermissionRole = ref<AccessRole | null>(null);
const formRef = ref<FormInstance>();

const filters = reactive({
  keyword: "",
  status: "",
});

const form = reactive<AccessRolePayload>(defaultForm());

const rules: FormRules<AccessRolePayload> = {
  name: [{ required: true, message: "请输入角色名称", trigger: "blur" }],
  code: [{ required: true, message: "请输入角色编码", trigger: "blur" }],
  status: [{ required: true, message: "请选择角色状态", trigger: "change" }],
};

const filteredRoles = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase();

  return roles.value.filter((role) => {
    const matchesKeyword =
      keyword.length === 0 ||
      [role.name, role.code, role.description].join(" ").toLowerCase().includes(keyword);
    const matchesStatus = filters.status.length === 0 || role.status === filters.status;
    return matchesKeyword && matchesStatus;
  });
});

const enabledCount = computed(() => {
  return roles.value.filter((role) => role.status === "启用").length;
});

const customCount = computed(() => {
  return roles.value.filter((role) => !["super_admin", "campus_owner", "front_desk", "teacher"].includes(role.code)).length;
});

const averagePermissionCount = computed(() => {
  if (roles.value.length === 0) {
    return 0;
  }

  return Math.round(
    roles.value.reduce((sum, role) => sum + role.permissionCount, 0) / roles.value.length,
  );
});

const dialogTitle = computed(() => {
  return editingRoleId.value ? "编辑角色" : "新建角色";
});

const canManageRoles = computed(() => authStore.hasPermission("roles:manage"));

function defaultForm(): AccessRolePayload {
  return {
    name: "",
    code: "",
    description: "",
    status: "启用",
  };
}

function resetForm() {
  Object.assign(form, defaultForm());
  editingRoleId.value = null;
  formRef.value?.clearValidate();
}

function openCreateDialog() {
  resetForm();
  roleDialogVisible.value = true;
}

function openEditDialog(role: AccessRole) {
  editingRoleId.value = role.id;
  Object.assign(form, {
    name: role.name,
    code: role.code,
    description: role.description,
    status: role.status,
  });
  roleDialogVisible.value = true;
  formRef.value?.clearValidate();
}

function closeRoleDialog() {
  roleDialogVisible.value = false;
  resetForm();
}

async function loadRoles() {
  loading.value = true;

  try {
    const result = await fetchRoleList();
    roles.value = result.list;
  } catch (error) {
    console.error(error);
    ElMessage.error("角色权限工作台加载失败");
  } finally {
    loading.value = false;
  }
}

async function submitRoleForm() {
  const formNode = formRef.value;
  if (!formNode) {
    return;
  }

  const valid = await formNode.validate().catch(() => false);
  if (!valid) {
    return;
  }

  saving.value = true;

  try {
    const payload = {
      name: form.name.trim(),
      code: form.code.trim(),
      description: form.description.trim(),
      status: form.status,
    };

    if (editingRoleId.value) {
      await updateRole(editingRoleId.value, payload);
      ElMessage.success("角色已更新");
    } else {
      await createRole(payload);
      ElMessage.success("角色已创建");
    }

    closeRoleDialog();
    await loadRoles();
  } catch (error: any) {
    console.error(error);
    const message = error?.response?.data?.message ?? "角色保存失败";
    ElMessage.error(message);
  } finally {
    saving.value = false;
  }
}

async function openPermissionDialog(role: AccessRole) {
  activePermissionRole.value = role;
  permissionDialogVisible.value = true;

  try {
    const result = await fetchRoleDetail(role.id);
    permissionGroups.value = result.permissionGroups;
    selectedPermissions.value = [...result.role.permissions];
  } catch (error) {
    console.error(error);
    ElMessage.error("权限配置加载失败");
  }
}

async function savePermissions() {
  const role = activePermissionRole.value;
  if (!role) {
    return;
  }

  saving.value = true;

  try {
    await saveRolePermissions(role.id, {
      permissions: selectedPermissions.value,
    });
    ElMessage.success("角色权限已保存");
    permissionDialogVisible.value = false;
    await loadRoles();
  } catch (error: any) {
    console.error(error);
    const message = error?.response?.data?.message ?? "权限保存失败";
    ElMessage.error(message);
  } finally {
    saving.value = false;
  }
}

function resetFilters() {
  filters.keyword = "";
  filters.status = "";
}

onMounted(() => {
  void loadRoles();
});
</script>

<template>
  <div class="page-stack">
    <section class="page-card page-card--table list-card">
      <div class="page-header">
        <div class="list-card__heading">
          <h2>角色列表</h2>
          <span class="list-card__count">共 {{ filteredRoles.length }} 条</span>
        </div>
        <div class="page-actions">
          <el-button v-if="canManageRoles" type="primary" @click="openCreateDialog">
            新建角色
          </el-button>
        </div>
      </div>

      <div class="metric-strip metric-strip--compact list-card__metrics">
        <article class="metric-tile">
          <span>角色总数</span>
          <strong>{{ roles.length }}</strong>
        </article>
        <article class="metric-tile">
          <span>启用角色</span>
          <strong>{{ enabledCount }}</strong>
        </article>
        <article class="metric-tile">
          <span>自定义角色</span>
          <strong>{{ customCount }}</strong>
        </article>
        <article class="metric-tile">
          <span>平均权限数</span>
          <strong>{{ averagePermissionCount }}</strong>
        </article>
      </div>

      <div class="filter-bar list-card__filters">
        <div class="toolbar-filters">
          <el-input v-model="filters.keyword" class="toolbar-field" clearable placeholder="搜索角色名称、编码或说明" />
          <el-select v-model="filters.status" class="toolbar-field" clearable placeholder="全部状态">
            <el-option label="启用" value="启用" />
            <el-option label="停用" value="停用" />
          </el-select>
        </div>

        <div class="toolbar-actions">
          <el-button plain @click="resetFilters">重置筛选</el-button>
        </div>
      </div>

      <div class="data-table-shell">
        <el-table v-loading="loading" :data="filteredRoles" stripe>
          <el-table-column label="角色" min-width="220">
            <template #default="{ row }">
              <div class="table-primary">
                <strong>{{ row.name }}</strong>
                <small>{{ row.code }}</small>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="说明" min-width="260" prop="description" />
          <el-table-column label="账号数" width="100" prop="userCount" />
          <el-table-column label="权限数" width="100" prop="permissionCount" />
          <el-table-column label="状态" width="110">
            <template #default="{ row }">
              <el-tag :type="row.status === '启用' ? 'success' : 'info'">
                {{ row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <div class="table-link-group">
                <el-button v-if="canManageRoles" link type="primary" @click="openEditDialog(row)">
                  编辑
                </el-button>
                <el-button link type="primary" @click="openPermissionDialog(row)">
                  权限配置
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </section>

    <el-dialog v-model="roleDialogVisible" :title="dialogTitle" width="680px" @closed="closeRoleDialog">
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <div class="dialog-grid">
          <el-form-item label="角色名称" prop="name">
            <el-input v-model="form.name" placeholder="例如 校区主管" />
          </el-form-item>
          <el-form-item label="角色编码" prop="code">
            <el-input v-model="form.code" placeholder="例如 campus_manager" />
          </el-form-item>
          <el-form-item class="full-span" label="角色说明">
            <el-input
              v-model="form.description"
              :rows="3"
              placeholder="写清这个角色负责什么，后面更容易维护"
              type="textarea"
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
          <el-button @click="closeRoleDialog">取消</el-button>
          <el-button :loading="saving" type="primary" @click="submitRoleForm">保存角色</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog
      v-model="permissionDialogVisible"
      :title="activePermissionRole ? `${activePermissionRole.name} · 权限配置` : '权限配置'"
      width="860px"
    >
      <div class="stack-list stack-list--spacious">
        <article
          v-for="group in permissionGroups"
          :key="group.key"
          class="detail-note"
        >
          <strong>{{ group.label }}</strong>
          <p>{{ group.description }}</p>
          <el-checkbox-group v-model="selectedPermissions" class="permission-grid">
            <el-checkbox
              v-for="permission in group.permissions"
              :key="permission.code"
              :label="permission.code"
            >
              <span class="permission-label">{{ permission.label }}</span>
              <small>{{ permission.description }}</small>
            </el-checkbox>
          </el-checkbox-group>
        </article>
      </div>

      <template #footer>
        <div class="dialog-actions">
          <el-button @click="permissionDialogVisible = false">取消</el-button>
          <el-button :loading="saving" type="primary" @click="savePermissions">保存权限</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>
