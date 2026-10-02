<template>
  <div class="password-reset-form">
    <el-form
        ref="passwordResetForm"
        :model="passwordResetFormData"
        :rules="rules"
        :validate-on-rule-change="false"
        hide-required-asterisk
    >
      <el-form-item :label="t('forms.passwordReset.oldPassword')" prop="password">
        <el-input
            v-model="passwordResetFormData.password"
            size="large"
            :placeholder="t('forms.passwordReset.oldPasswordPlaceholder')"
        />
      </el-form-item>
      <el-form-item :label="t('forms.passwordReset.newPassword')" prop="new_password">
        <el-input
            v-model="passwordResetFormData.new_password"
            size="large"
            :placeholder="t('forms.passwordReset.newPasswordPlaceholder')"
        />
      </el-form-item>
      <el-form-item>
        <div class="button-group">
          <el-button
              type="primary"
              size="large"
              @click="submitForm">{{ t('common.confirm') }}
          </el-button>
          <el-button
              size="large"
              @click="layoutStore.state.passwordResetVisible = false"
          >{{ t('common.cancel') }}
          </el-button>
        </div>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import {useLayoutStore} from "@/stores/layout";
import {reactive, ref} from "vue";
import {ElMessage, type FormInstance, type FormRules} from "element-plus";
import {userResetPassword, type UserResetPasswordRequest} from "@/api/user";
import router from "@/router";
import {useUserStore} from "@/stores/user";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const layoutStore = useLayoutStore()
const userStore = useUserStore()

const passwordResetForm = ref<FormInstance>()

const passwordResetFormData = reactive<UserResetPasswordRequest>({
  password: '',
  new_password: '',
})

const rules = reactive<FormRules<UserResetPasswordRequest>>({
  password: [{
    required: true,
    min: 8,
    max: 20,
    trigger: 'change',
    message: t('forms.passwordReset.passwordLength')
  }],
  new_password: [{
    required: true,
    min: 8,
    max: 20,
    trigger: 'change',
    message: t('forms.passwordReset.passwordLength')
  }]
})

const submitForm = async () => {
  const isValid: boolean = await new Promise((resolve) => {
    passwordResetForm.value?.validate((valid: boolean) => {
      resolve(valid)
    })
  })

  if (isValid) {
    const res = await userResetPassword(passwordResetFormData)
    if (res.code === 0) {
      ElMessage.success(res.msg)
      layoutStore.state.passwordResetVisible = false
      router.push({name: "index"}).then(
          userStore.reset
      )
    }
  }
}
</script>

<style scoped lang="scss">
.password-reset-form {
  .el-form {
    .el-form-item {
      .button-group {
        margin-left: auto;
        display: flex;
        gap: var(--sp-2);

        .el-button {
          height: 40px;
          min-width: 88px;
          font-weight: 500;
        }
      }
    }
  }
}
</style>
