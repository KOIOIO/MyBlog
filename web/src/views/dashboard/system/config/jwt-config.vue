<template>
  <div class="jwt-config">
    <el-col :span="12">
      <div class="info">
        <div class="page-title">{{ t('system.config.jwt.title') }}</div>
        <div class="content">
          <el-form
              :model="jwtInfo"
              :validate-on-rule-change="false"
              label-width="auto"
              style="max-width: 400px"
          >
            <el-form-item :label="t('system.config.jwt.accessTokenSecret')">
              <el-input @change="updateJwtInfo" v-model="jwtInfo.access_token_secret" type="password" show-password/>
            </el-form-item>
            <el-form-item :label="t('system.config.jwt.accessTokenExpiry')">
              <el-input @change="updateJwtInfo" v-model="jwtInfo.access_token_expiry_time"/>
            </el-form-item>
            <el-form-item :label="t('system.config.jwt.refreshTokenSecret')">
              <el-input @change="updateJwtInfo" v-model.number="jwtInfo.refresh_token_secret" type="password"
                        show-password/>
            </el-form-item>
            <el-form-item :label="t('system.config.jwt.refreshTokenExpiry')">
              <el-input @change="updateJwtInfo" v-model="jwtInfo.refresh_token_expiry_time"/>
            </el-form-item>
            <el-form-item :label="t('system.config.jwt.issuer')">
              <el-input @change="updateJwtInfo" v-model="jwtInfo.issuer"/>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </el-col>
  </div>
</template>

<script setup lang="ts">
import {ref, watch} from "vue";
import {type Jwt, getJwt, updateJwt} from "@/api/config";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const jwtInfo = ref<Jwt>({
  access_token_secret: '',
  refresh_token_secret: '',
  access_token_expiry_time: '',
  refresh_token_expiry_time: '',
  issuer: '',
})

const getJwtInfo = async () => {
  const res = await getJwt()
  if (res.code === 0) {
    jwtInfo.value = res.data
  }
}

getJwtInfo()

const shouldRefreshInfo = ref(false)
watch(() => shouldRefreshInfo.value, (newVal) => {
  if (newVal) {
    getJwtInfo()
    shouldRefreshInfo.value = false
  }
})

const updateJwtInfo = async () => {
  const res = await updateJwt(jwtInfo.value)
  console.log(jwtInfo.value)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  } else {
    shouldRefreshInfo.value = true
  }
}
</script>

<style scoped lang="scss">
.jwt-config {
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
