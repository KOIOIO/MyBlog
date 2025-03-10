<template>
  <div class="qq-config">
    <el-col :span="12">
      <div class="info">
        <div class="title">
          <el-row>QQ登录配置</el-row>
        </div>
        <div class="content">
          <el-form
              :model="qqInfo"
              :validate-on-rule-change="false"
              label-width="auto"
              style="max-width: 400px"
          >

            <el-form-item label="启用QQ登录">
              <el-switch v-model="qqInfo.enable" @change="updateQQInfo"/>
            </el-form-item>
            <el-form-item label="应用ID">
              <el-input @change="updateQQInfo" v-model="qqInfo.app_id"/>
            </el-form-item>
            <el-form-item label="应用密钥">
              <el-input @change="updateQQInfo" v-model.number="qqInfo.app_key" type="password" show-password/>
            </el-form-item>
            <el-form-item label="网站回调域">
              <el-input @change="updateQQInfo" v-model="qqInfo.redirect_uri"/>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </el-col>
  </div>
</template>

<script setup lang="ts">
import {ref, watch} from "vue";
import {getQQ, type QQ, updateQQ} from "@/api/config";

const qqInfo = ref<QQ>({
  enable: false,
  app_id: '',
  app_key:'',
  redirect_uri: '',
})

const getQQInfo = async () => {
  const res = await getQQ()
  if (res.code === 0) {
    qqInfo.value = res.data
  }
}

getQQInfo()

const shouldRefreshInfo = ref(false)
watch(() => shouldRefreshInfo.value, (newVal) => {
  if (newVal) {
    getQQInfo()
    shouldRefreshInfo.value = false
  }
})

const updateQQInfo = async () => {
  const res = await updateQQ(qqInfo.value)
  console.log(qqInfo.value)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  } else {
    shouldRefreshInfo.value = true
  }
}
</script>

<style scoped lang="scss">
.qq-config {
  display: flex;

  .info {
    .title {
      border-left: 5px solid #8A2BE2; /* 更改为蓝紫色 */
      padding-left: 10px;
      color: #8A2BE2; /* 更改为蓝紫色 */
    }

    .content {
      margin: 20px;
      background-color: #F0F8FF; /* 添加背景色 */
      border-radius: 10px; /* 添加圆角 */
      padding: 20px; /* 添加内边距 */
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); /* 添加阴影 */
    }
  }

  .el-form-item__label {
    color: #8A2BE2; /* 更改标签颜色为蓝紫色 */
  }

  .el-input__inner, .el-switch__core {
    border-color: #8A2BE2; /* 更改输入框和开关的边框颜色 */
  }

  .el-input__inner:focus, .el-switch__core:focus {
    border-color: #8A2BE2; /* 更改输入框和开关聚焦时的边框颜色 */
  }
}
</style>
