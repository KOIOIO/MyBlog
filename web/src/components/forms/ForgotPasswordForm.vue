<template>
  <div class="forgot-password-form">
    <el-form
        ref="forgotPasswordForm"
        :model="forgotPasswordFormData"
        :rules="rules"
        :validate-on-rule-change="false"
        hide-required-asterisk
        @keyup.enter="submitForm"
    >

      <el-form-item :label="t('forms.forgotPassword.email')" prop="email">
        <el-input
            v-model="forgotPasswordFormData.email"
            size="large"
            :placeholder="t('forms.forgotPassword.emailPlaceholder')"
        />
      </el-form-item>
      <el-form-item prop="captcha">
        <div class="captcha">
          <el-input
              v-model="emailRequest.captcha"
              :placeholder="t('forms.forgotPassword.imageCaptchaPlaceholder')"
              size="large"
              maxlength="6"
              minlength="6"
          />
          <el-image :src="picPath" alt="" @click="emailVerify"/>
          <el-button @click="sendCode">{{ t('forms.forgotPassword.sendCode') }}</el-button>
        </div>
      </el-form-item>
      <el-form-item :label="t('forms.forgotPassword.emailCaptcha')" prop="verification_code">
        <el-input
            v-model="forgotPasswordFormData.verification_code"
            size="large"
            :placeholder="t('forms.forgotPassword.emailCaptchaPlaceholder')"
        />
      </el-form-item>
      <el-form-item :label="t('forms.forgotPassword.password')" prop="new_password">
        <el-input
            v-model="forgotPasswordFormData.new_password"
            show-password
            size="large"
            type="password"
            :placeholder="t('forms.forgotPassword.newPasswordPlaceholder')"
        />
      </el-form-item>
      <el-form-item :label="t('forms.forgotPassword.confirmPassword')">
        <el-input
            v-model="repeatPassword"
            show-password
            size="large"
            type="password"
            :placeholder="t('forms.forgotPassword.confirmNewPasswordPlaceholder')"
        />
        <el-text v-if="forgotPasswordFormData.new_password!==repeatPassword">{{ t('forms.forgotPassword.passwordMismatch') }}</el-text>
      </el-form-item>
      <el-form-item>
        <el-button
            type="primary"
            size="large"
            @click="submitForm"
        >{{ t('forms.forgotPassword.submit') }}
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref,} from "vue";
import {forgotPassword, type ForgotPasswordRequest} from "@/api/user";
import {captcha, type EmailRequest, sendEmailVerificationCode} from "@/api/base";
import type {FormInstance, FormRules} from 'element-plus';
import {useLayoutStore} from "@/stores/layout";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const forgotPasswordForm = ref<FormInstance>()

const forgotPasswordFormData = reactive<ForgotPasswordRequest>({
  email: "",
  verification_code: "",
  new_password:"",
})

const repeatPassword = ref('')

const emailRequest = reactive<EmailRequest>({
  email: "",
  captcha: "",
  captcha_id: "",
})

const rules = reactive<FormRules<ForgotPasswordRequest>>({
  email: [{
    required: true,
    type:'email',
    trigger: 'blur',
    message: t('forms.forgotPassword.emailFormat')
  }],
  verification_code: [{
    required: true,
    len: 6,
    trigger: 'blur',
    message: t('forms.forgotPassword.codeLength')
  }],
  new_password:[{
    required:true,
    min:8,
    max:20,
    trigger:'change',
    message:t('forms.forgotPassword.passwordLength')
  }]
})

const layoutStore = useLayoutStore()

const picPath = ref('')

const emailVerify = () => {
  captcha().then(async (res) => {
    picPath.value = res.data.pic_path
    emailRequest.captcha_id = res.data.captcha_id
  })
}

emailVerify()

const sendCode = async () => {
  emailRequest.email = forgotPasswordFormData.email
  const res = await sendEmailVerificationCode(emailRequest)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  } else {
    emailVerify()
  }
}

const submitForm = async () => {
  const isValid: boolean = await new Promise((resolve) => {
    forgotPasswordForm.value?.validate((valid: boolean) => {
      resolve(valid)
    })
  })

  if (isValid) {
    const res = await forgotPassword(forgotPasswordFormData)

    if (res.code===0) {
      ElMessage.success(res.msg)
      layoutStore.state.forgotPasswordVisible = false
      layoutStore.state.popoverVisible = false
    }
  }
  return false
}
</script>


<style scoped lang="scss">
.forgot-password-form {
  .el-form {
    flex: 1;
    min-width: 0;

    .captcha {
      display: flex;
      gap: var(--sp-2);
      align-items: center;
      width: 100%;

      .el-input {
        flex: 1;
      }

      .el-image {
        height: 40px;
        border-radius: var(--radius-sm);
        overflow: hidden;
        cursor: pointer;
        flex-shrink: 0;
      }

      .el-button {
        height: 40px;
        flex-shrink: 0;
      }
    }

    .el-text {
      color: var(--el-color-danger);
      font-size: var(--fs-12);
    }

    .el-form-item:last-child {
      margin-bottom: 0;

      .el-button {
        width: 100%;
        height: 40px;
        font-weight: 500;
      }
    }
  }
}
</style>

<style lang="scss">
.forgot-password-form .el-form .captcha .el-input__wrapper {
  height: 40px;
}

/* 密码可见性切换图标：色渡 150ms + 按压反馈（§8） */
.forgot-password-form .el-input__password {
  cursor: pointer;
  transition: color 150ms ease-out, transform 80ms ease-out;

  &:active {
    transform: scale(0.92);
  }
}
</style>
