<template>
  <div class="qiniu-config">
    <el-col :span="12">
      <div class="info">
        <div class="page-title">{{ t('system.config.qiniu.title') }}</div>
        <div class="content">
          <el-form
              :model="qiniuInfo"
              :validate-on-rule-change="false"
              label-width="auto"
              style="max-width: 400px"
          >
            <el-form-item :label="t('system.config.qiniu.zone')">
              <el-input @change="updateQiniuInfo" v-model="qiniuInfo.zone"/>
            </el-form-item>
            <el-form-item :label="t('system.config.qiniu.bucket')">
              <el-input @change="updateQiniuInfo" v-model.number="qiniuInfo.bucket"/>
            </el-form-item>
            <el-form-item :label="t('system.config.qiniu.accessKey')">
              <el-input @change="updateQiniuInfo" v-model="qiniuInfo.access_key" type="password" show-password/>
            </el-form-item>
            <el-form-item :label="t('system.config.qiniu.secretKey')">
              <el-input @change="updateQiniuInfo" v-model="qiniuInfo.secret_key" type="password" show-password/>
            </el-form-item>
            <el-form-item :label="t('system.config.qiniu.cdnDomain')">
              <el-input @change="updateQiniuInfo" v-model="qiniuInfo.img_path"/>
            </el-form-item>
            <el-form-item :label="t('system.config.qiniu.useCdn')">
              <el-switch v-model="qiniuInfo.use_cdn_domains" @change="updateQiniuInfo"/>
            </el-form-item>
            <el-form-item :label="t('system.config.qiniu.useHttps')">
              <el-switch v-model="qiniuInfo.use_https" @change="updateQiniuInfo"/>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </el-col>
  </div>
</template>

<script setup lang="ts">
import {ref, watch} from "vue";
import {type Qiniu, getQiniu, updateQiniu} from "@/api/config";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const qiniuInfo = ref<Qiniu>({
  zone: '',
  bucket: '',
  img_path: '',
  access_key: '',
  secret_key: '',
  use_https: false,
  use_cdn_domains: false,
})

const getQiniuInfo = async () => {
  const res = await getQiniu()
  if (res.code === 0) {
    qiniuInfo.value = res.data
  }
}

getQiniuInfo()

const shouldRefreshInfo = ref(false)
watch(() => shouldRefreshInfo.value, (newVal) => {
  if (newVal) {
    getQiniuInfo()
    shouldRefreshInfo.value = false
  }
})

const updateQiniuInfo = async () => {
  const res = await updateQiniu(qiniuInfo.value)
  console.log(qiniuInfo.value)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  } else {
    shouldRefreshInfo.value = true
  }
}
</script>

<style scoped lang="scss">
.qiniu-config {
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
