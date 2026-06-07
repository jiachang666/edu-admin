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
    authStore.setSession(
      result.accessToken,
      result.user,
      result.roles,
      result.roleNames,
      result.permissions,
    );
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
    <div class="login-stage">
      <section class="login-poster">
        <div class="login-eyebrow">Education Admin</div>
        <h1>培训机构管理系统</h1>
        <p>统一处理学员、班级、排课、签到、作业和通知。</p>
      </section>

      <section class="login-card">
        <div class="login-card-header">
          <h2 class="login-card-title">登录</h2>
        </div>

        <div class="login-demo-strip">
          <div>
            <strong>默认账号</strong>
            <span>账号 `admin`，密码使用当前预填内容</span>
          </div>
          <el-tag type="info">可直接登录</el-tag>
        </div>

        <el-form label-position="top" @submit.prevent="submit">
          <el-form-item label="账号">
            <el-input v-model="form.username" placeholder="请输入登录账号" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input
              v-model="form.password"
              placeholder="请输入密码"
              show-password
              type="password"
            />
          </el-form-item>
          <el-button
            class="full-width"
            :loading="loading"
            type="primary"
            @click="submit"
          >
            登录并进入后台
          </el-button>
        </el-form>
      </section>
    </div>
  </div>
</template>
