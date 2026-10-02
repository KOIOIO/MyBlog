<template>
  <div class="user-info">
    <el-col :span="12">
      <div class="info">
        <div class="page-title">{{ t('system.userCenter.info.title') }}</div>
        <div class="content">
          <el-form
              ref="userChangeInfoForm"
              :model="userChangeInfoFormData"
              :rules="rules"
              :validate-on-rule-change="false"
              hide-required-asterisk
              label-width="auto"
              style="max-width: 400px"
          >
            <el-form-item :label="t('system.userCenter.info.avatar')">
              <avatar-uploader @changed="onAvatarChanged"/>
            </el-form-item>
            <el-form-item label="uuid">
              {{ userInfo.uuid }}
            </el-form-item>
            <el-form-item :label="t('system.userCenter.info.username')" prop="username">
              <el-input @change="updateUserInfo" v-model="userChangeInfoFormData.username"/>
            </el-form-item>
            <el-form-item :label="t('system.userCenter.info.address')" prop="address">
              <el-input @change="updateUserInfo" v-model="userChangeInfoFormData.address"/>
            </el-form-item>
            <el-form-item :label="t('system.userCenter.info.signature')" prop="signature">
              <el-input @change="updateUserInfo" v-model="userChangeInfoFormData.signature" type="textarea" :rows="2"/>
            </el-form-item>
            <el-form-item :label="t('system.userCenter.info.email')">
              {{ userInfo.email }}
            </el-form-item>
            <el-form-item :label="t('system.userCenter.info.role')">
              {{ userInfo.role_id === 1 ? t('system.userCenter.info.normalUser') : t('system.userCenter.info.admin') }}
            </el-form-item>
            <el-form-item :label="t('system.userCenter.info.registerSource')">
              {{ userInfo.register }}
            </el-form-item>
          </el-form>
        </div>
      </div>
      <div class="operation" v-if="userStore.state.userInfo.register==='邮箱'">
        <div class="page-title">{{ t('system.userCenter.info.operation') }}</div>
        <div class="content">
          <el-button @click="layoutStore.state.passwordResetVisible = true">{{ t('system.userCenter.info.changePassword') }}</el-button>
        </div>
        <el-dialog
            v-model="passwordResetVisible"
            width="500"
            align-center
            destroy-on-close
            :before-close="passwordResetVisibleSynchronization"
        >
          <template #header>
            {{ t('system.userCenter.info.changePassword') }}
          </template>
          <password-reset-form/>
          <template #footer>
          </template>
        </el-dialog>
      </div>
    </el-col>
    <el-col :span="12">
      <div class="card">
        <div class="page-title">{{ t('system.userCenter.info.userCard') }}</div>
      </div>
      <div class="content">
        <user-card :key="cardKey" :uuid="userInfo.uuid" :user-card-info="null"/>
      </div>
    </el-col>
  </div>
</template>

<script setup lang="ts">
import {useUserStore} from "@/stores/user";
import {reactive, ref, watch} from "vue";
import {userChangeInfo, type UserChangeInfoRequest} from "@/api/user";
import type {FormInstance, FormRules} from "element-plus";
import UserCard from "@/components/widgets/UserCard.vue";
import {useLayoutStore} from "@/stores/layout";
import PasswordResetForm from "@/components/forms/PasswordResetForm.vue";
import AvatarUploader from "@/components/common/AvatarUploader.vue";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const userStore = useUserStore()
const layoutStore = useLayoutStore()

const userInfo = ref(userStore.state.userInfo)

const userChangeInfoForm = ref<FormInstance>()

const userChangeInfoFormData = reactive<UserChangeInfoRequest>({
  username: userInfo.value.username,
  address: userInfo.value.address,
  signature: userInfo.value.signature,
})

const rules = reactive<FormRules<UserChangeInfoRequest>>({
  username: [{
    required:true,
    max:20,
    trigger:'blur',
    message:t('system.userCenter.info.usernameMax')
  }],
  address: [{
    max: 200,
    trigger: 'blur',
    message:t('system.userCenter.info.addressMax')
  }],
  signature: [{
    max: 320,
    trigger: 'blur',
    message:t('system.userCenter.info.signatureMax')
  }]
})

const cardKey = ref(0)

const onAvatarChanged = () => {
  cardKey.value += 1
}

const updateUserInfo = async () => {
  const isValid: boolean = await new Promise((resolve) => {
    userChangeInfoForm.value?.validate((valid: boolean) => {
      resolve(valid)
    })
  })

  if (isValid) {
    const res = await userChangeInfo(userChangeInfoFormData)
    if (res.code === 0) {
      cardKey.value += 1
    }
  }
}

const passwordResetVisible = ref(layoutStore.state.passwordResetVisible)
watch(
    () => layoutStore.state.passwordResetVisible,
    (newValue) => {
      passwordResetVisible.value = newValue
    }
)

const passwordResetVisibleSynchronization = () => {
  layoutStore.state.passwordResetVisible = false
}
</script>

<style scoped lang="scss">
.user-info {
  display: flex;
  gap: var(--sp-5);
  align-items: flex-start;

  .page-title {
    font-size: var(--fs-20);
    font-weight: 600;
    color: var(--text-primary);
    line-height: var(--lh-title);
    margin-bottom: var(--sp-4);
  }

  .info,
  .operation,
  .card {
    background: var(--bg-elevated);
    border: 1px solid var(--border);
    border-radius: var(--radius-md);
    padding: var(--sp-5);
    margin-bottom: var(--sp-5);
  }

  .info .content,
  .operation .content {
    .avatar-uploader {
      margin-bottom: var(--sp-2);
    }
  }

  .card {
    margin-bottom: var(--sp-3);

    .page-title {
      margin-bottom: 0;
    }
  }
}
</style>
