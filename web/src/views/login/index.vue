<template>
  <div class="login">
    <!-- 添加一些内容以展示风格变化 -->
    <h1>欢迎登录</h1>
    <p>请稍候，正在登录中...</p>
  </div>
</template>

<script setup lang="ts">
import {useRoute, useRouter} from "vue-router";
import service, {type ApiResponse} from "@/utils/request";
import type {LoginResponse} from "@/api/user";
import {useUserStore} from "@/stores/user";

ElMessage.info('登录中...')

const route = useRoute();
const router = useRouter();

const flag = route.query.flag as string;
const code = route.query.code as string;

interface QQLoginRequest {
  flag: string;
  code: string;
}
const qqLogin = (data: QQLoginRequest): Promise<ApiResponse<LoginResponse>> => {
  return service({
    url: '/user/login',
    method: 'post',
    params: data
  });
}

const loginIn = async () => {
  const qqLoginRequest: QQLoginRequest = {
    flag: flag,
    code: code,
  }
  const res = await qqLogin(qqLoginRequest)
  if (res.code === 0) {
    const userStore = useUserStore()
    userStore.state.userInfo = res.data.user
    userStore.state.accessToken = res.data.access_token
    userStore.state.isUserLoggedInBefore = true
    router.push({name: 'index'}).then(() => {
      ElMessage.success(res.msg)
    })
  } else {
    router.push({name: 'index'}).then()
  }
}

loginIn()
</script>

<style scoped>
.login {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  font-family: 'Roboto', sans-serif;
}

h1 {
  font-size: 3em;
  margin-bottom: 0.5em;
}

p {
  font-size: 1.5em;
}
</style>
