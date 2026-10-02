<template>
  <div class="system-config">
    <el-col :span="12">
      <div class="info">
        <div class="page-title">{{ t('system.config.system.title') }}</div>
        <div class="content">
          <el-form
              :model="systemInfo"
              :validate-on-rule-change="false"
              label-width="auto"
              style="max-width: 400px"
          >
            <el-form-item :label="t('system.config.system.multipoint')">
              <el-switch v-model="systemInfo.use_multipoint" @change="updateSystemInfo"/>
            </el-form-item>
            <el-form-item :label="t('system.config.system.sessionSecret')">
              <el-input @change="updateSystemInfo" v-model="systemInfo.sessions_secret" type="password" show-password/>
            </el-form-item>
            <el-form-item :label="t('system.config.system.ossType')">
              <el-select
                  @change="updateSystemInfo"
                  v-model="systemInfo.oss_type"
                  :placeholder="t('system.config.system.selectPlaceholder')"
                  style="width: 200px"
              >
                <el-option
                    v-for="item in ossTypeOptions"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                />
              </el-select>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </el-col>
  </div>
</template>

<script setup lang="ts">
import {computed, ref, watch} from "vue";
import {type System, getSystem, updateSystem} from "@/api/config";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const systemInfo = ref<System>({
  use_multipoint: false,
  sessions_secret: '',
  oss_type: '',
})

const ossTypeOptions = computed(() => [
  {
    value: 'local',
    label: t('system.config.system.ossLocal'),
  },
  {
    value: 'qiniu',
    label: t('system.config.system.ossQiniu'),
  },
])

const getSystemInfo = async () => {
  const res = await getSystem()
  if (res.code === 0) {
    systemInfo.value = res.data
  }
}

getSystemInfo()

const shouldRefreshInfo = ref(false)
watch(() => shouldRefreshInfo.value, (newVal) => {
  if (newVal) {
    getSystemInfo()
    shouldRefreshInfo.value = false
  }
})

const updateSystemInfo = async () => {
  const res = await updateSystem(systemInfo.value)
  console.log(systemInfo.value)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  } else {
    shouldRefreshInfo.value = true
  }
}
</script>

<style scoped lang="scss">
.system-config {
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
