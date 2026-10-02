<template>
  <div class="gaode-config">
    <el-col :span="12">
      <div class="info">
        <div class="page-title">{{ t('system.config.gaode.title') }}</div>
        <div class="content">
          <el-form
              :model="gaodeInfo"
              :validate-on-rule-change="false"
              label-width="auto"
              style="max-width: 400px"
          >
            <el-form-item :label="t('system.config.gaode.enable')">
              <el-switch v-model="gaodeInfo.enable" @change="updateGaodeInfo"/>
            </el-form-item>
            <el-form-item :label="t('system.config.gaode.key')">
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
import {useI18n} from "vue-i18n";

const {t} = useI18n()

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
  .page-title {
    font-size: var(--fs-20);
    font-weight: 600;
    color: var(--text-primary);
    line-height: var(--lh-title);
    margin-bottom: var(--sp-4);
  }

  .content {
    background: var(--bg-elevated);
    border: 1px solid var(--border);
    border-radius: var(--radius-md);
    padding: var(--sp-5);
  }
}
</style>
