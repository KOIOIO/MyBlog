<template>
  <div class="gaode-config">
    <el-col :span="12">
      <div class="info">
        <div class="title">
          <el-row>高德配置</el-row>
        </div>
        <div class="content">
          <el-form
              :model="gaodeInfo"
              :validate-on-rule-change="false"
              label-width="auto"
              style="max-width: 400px"
          >
            <el-form-item label="是否开启">
              <el-switch v-model="gaodeInfo.enable" @change="updateGaodeInfo"/>
            </el-form-item>
            <el-form-item label="高德密钥">
              <el-input @change="updateGaodeInfo" v-model="gaodeInfo.key" type="password" show-password/>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </el-col>
  </div>
</template>

<script setup lang="ts">
import {ref, watch} from "vue";
import {type Gaode, getGaode, updateGaode} from "@/api/config";

const gaodeInfo = ref<Gaode>({
  enable: false,
  key: '',
})

const getGaodeInfo = async () => {
  const res = await getGaode()
  if (res.code === 0) {
    gaodeInfo.value = res.data
  }
}

getGaodeInfo()

const shouldRefreshInfo = ref(false)
watch(() => shouldRefreshInfo.value, (newVal) => {
  if (newVal) {
    getGaodeInfo()
    shouldRefreshInfo.value = false
  }
})

const updateGaodeInfo = async () => {
  const res = await updateGaode(gaodeInfo.value)
  console.log(gaodeInfo.value)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  } else {
    shouldRefreshInfo.value = true
  }
}
</script>

<style scoped lang="scss">
.gaode-config {
  display: flex;

  .info {
    .title {
      border-left: 5px solid #8A2BE2; /* 蓝紫色 */
      padding-left: 10px;
      color: #8A2BE2; /* 蓝紫色 */
    }

    .content {
      margin: 20px;
      background-color: #F0F8FF; /* 浅蓝色背景 */
      padding: 20px;
      border-radius: 10px;
      box-shadow: 0 0 10px rgba(138, 43, 178, 0.5); /* 蓝紫色阴影 */
    }
  }

  .el-form-item__label {
    color: #8A2BE2; /* 蓝紫色 */
  }

  .el-switch__core {
    background-color: #8A2BE2; /* 蓝紫色 */
  }

  .el-input__inner {
    border-color: #8A2BE2; /* 蓝紫色 */
  }

  .el-input__inner:focus {
    border-color: #8A2BE2; /* 蓝紫色 */
    box-shadow: 0 0 5px rgba(138, 43, 178, 0.5); /* 蓝紫色阴影 */
  }
}
</style>
