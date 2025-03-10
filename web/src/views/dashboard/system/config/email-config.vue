<template>
  <div class="email-config">
    <el-col :span="12">
      <div class="info">
        <div class="title">
          <el-row>邮箱配置</el-row>
        </div>
        <div class="content">
          <el-form
              :model="emailInfo"
              :validate-on-rule-change="false"
              label-width="auto"
              style="max-width: 400px"
          >
            <el-form-item label="服务器地址">
              <el-input @change="updateEmailInfo" v-model="emailInfo.host"/>
            </el-form-item>
            <el-form-item label="服务器端口">
              <el-input @change="updateEmailInfo" v-model.number="emailInfo.port"/>
            </el-form-item>
            <el-form-item label="发件人邮箱">
              <el-input @change="updateEmailInfo" v-model="emailInfo.from"/>
            </el-form-item>
            <el-form-item label="发件人昵称">
              <el-input @change="updateEmailInfo" v-model="emailInfo.nickname"/>
            </el-form-item>
            <el-form-item label="邮箱密钥">
              <el-input @change="updateEmailInfo" v-model="emailInfo.secret" type="password" show-password/>
            </el-form-item>
            <el-form-item label="使用SSL">
              <el-switch v-model="emailInfo.is_ssl" @change="updateEmailInfo"/>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </el-col>
  </div>
</template>

<script setup lang="ts">
import {ref, watch} from "vue";
import {type Email, getEmail, updateEmail} from "@/api/config";

const emailInfo = ref<Email>({
  host: '',
  port: 0,
  from: '',
  nickname: '',
  secret: '',
  is_ssl: false,
})

const getEmailInfo = async () => {
  const res = await getEmail()
  if (res.code === 0) {
    emailInfo.value = res.data
  }
}

getEmailInfo()

const shouldRefreshInfo = ref(false)
watch(() => shouldRefreshInfo.value, (newVal) => {
  if (newVal) {
    getEmailInfo()
    shouldRefreshInfo.value = false
  }
})

const updateEmailInfo = async () => {
  const res = await updateEmail(emailInfo.value)
  console.log(emailInfo.value)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  } else {
    shouldRefreshInfo.value = true
  }
}
</script>

<style scoped lang="scss">
.email-config {
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
      box-shadow: 0 0 10px rgba(138, 43, 178, 0.5); /* 添加阴影效果 */
    }
  }

  .el-form-item__label {
    color: #8A2BE2; /* 更改标签颜色为蓝紫色 */
  }

  .el-input__inner, .el-switch__core {
    border-color: #8A2BE2; /* 更改输入框和开关的边框颜色 */
  }

  .el-input__inner:focus, .el-switch__core:focus {
    border-color: #8A2BE2; /* 更改输入框和开关的聚焦边框颜色 */
  }

  .el-button--primary {
    background-color: #8A2BE2; /* 更改按钮颜色为蓝紫色 */
    border-color: #8A2BE2; /* 更改按钮边框颜色为蓝紫色 */
  }
}
</style>
