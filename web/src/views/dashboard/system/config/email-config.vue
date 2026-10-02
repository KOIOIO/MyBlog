<template>
  <div class="email-config">
    <el-col :span="12">
      <div class="info">
        <div class="page-title">{{ t('system.config.email.title') }}</div>
        <div class="content">
          <el-form
              :model="emailInfo"
              :validate-on-rule-change="false"
              label-width="auto"
              style="max-width: 400px"
          >
            <el-form-item :label="t('system.config.email.host')">
              <el-input @change="updateEmailInfo" v-model="emailInfo.host"/>
            </el-form-item>
            <el-form-item :label="t('system.config.email.port')">
              <el-input @change="updateEmailInfo" v-model.number="emailInfo.port"/>
            </el-form-item>
            <el-form-item :label="t('system.config.email.from')">
              <el-input @change="updateEmailInfo" v-model="emailInfo.from"/>
            </el-form-item>
            <el-form-item :label="t('system.config.email.nickname')">
              <el-input @change="updateEmailInfo" v-model="emailInfo.nickname"/>
            </el-form-item>
            <el-form-item :label="t('system.config.email.secret')">
              <el-input @change="updateEmailInfo" v-model="emailInfo.secret" type="password" show-password/>
            </el-form-item>
            <el-form-item :label="t('system.config.email.ssl')">
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
import {useI18n} from "vue-i18n";

const {t} = useI18n()

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
