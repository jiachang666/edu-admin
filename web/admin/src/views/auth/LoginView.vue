<script setup lang="ts">
import { Calendar, Message, Reading, UserFilled } from "@element-plus/icons-vue";
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

const featureCards = [
  {
    title: "班级与学员",
    description: "把谁在班里、家长怎么联系、当前状态如何放在同一处。",
    icon: UserFilled,
  },
  {
    title: "排课与签到",
    description: "每天几点上、谁来上、还差哪些签到，一眼就能看到。",
    icon: Calendar,
  },
  {
    title: "课程与通知",
    description: "课程模板和消息发布都能接着走，流程不会断在半路。",
    icon: Message,
  },
];

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
        <div class="login-eyebrow">Training Operations Console</div>
        <h1>把培训机构的日常节奏，收进一张安静又清楚的工作台。</h1>
        <p>
          从老师、学员、课程、班级到排课与通知，这个后台先把最常用的运营路径整理顺，让教务和老师都能快速接手。
        </p>

        <div class="login-feature-grid">
          <article
            v-for="feature in featureCards"
            :key="feature.title"
            class="login-feature-card"
          >
            <el-icon size="18">
              <component :is="feature.icon" />
            </el-icon>
            <strong>{{ feature.title }}</strong>
            <span>{{ feature.description }}</span>
          </article>
        </div>

        <div class="login-quote-card">
          <strong>Preview Story</strong>
          <p>
            这不是一个泛后台壳子，而是一套更贴近真实教培场景的起步工作台，适合先跑顺“学员
            - 班级 - 排课 - 通知”的主流程。
          </p>
        </div>
      </section>

      <section class="login-card">
        <div class="login-card-header">
          <div class="section-kicker">
            <el-icon><Reading /></el-icon>
            <span>进入工作台</span>
          </div>
          <h2 class="login-card-title">欢迎回来</h2>
          <p class="login-card-copy">
            当前是演示环境，适合先确认页面结构、视觉风格和业务流程。
          </p>
        </div>

        <div class="login-demo-strip">
          <div>
            <strong>默认体验账号</strong>
            <span>账号 admin，密码可直接使用当前预填内容</span>
          </div>
          <el-tag type="success">Demo Ready</el-tag>
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

        <div class="login-form-note">
          如果登录失败，通常是后端服务还没启动；启动成功后，这里会直接进入首页总览。
        </div>
      </section>
    </div>
  </div>
</template>
