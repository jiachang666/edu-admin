<script setup lang="ts">
import { ElMessage } from "element-plus";
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { login } from "../../api/auth";
import { useAuthStore } from "../../stores/auth";

const router = useRouter();
const authStore = useAuthStore();
const loading = ref(false);

const form = reactive({
  username: "admin",
  password: "123456",
});

async function submit() {
  loading.value = true;

  try {
    const result = await login(form);
    authStore.setSession(result.accessToken, result.user.displayName);
    await router.push("/dashboard");
    ElMessage.success("登录成功");
  } catch (error) {
    console.error(error);
    ElMessage.error("登录失败，请确认后端服务已经启动");
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="login-page">
    <div class="login-card">
      <h1>Edu Admin</h1>
      <p>培训机构后台起步骨架</p>
      <el-form label-position="top">
        <el-form-item label="账号">
          <el-input v-model="form.username" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" show-password type="password" />
        </el-form-item>
        <el-button
          class="full-width"
          :loading="loading"
          type="primary"
          @click="submit"
        >
          登录
        </el-button>
      </el-form>
    </div>
  </div>
</template>
