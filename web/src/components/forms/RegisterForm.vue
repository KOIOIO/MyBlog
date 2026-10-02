<template>
  <div class="register-form">
    <el-form
        ref="registerForm"
        :model="registerFormData"
        :rules="rules"
        :validate-on-rule-change="false"
        hide-required-asterisk
        @keyup.enter="submitForm"
    >

      <el-form-item :label="t('forms.register.username')" prop="username">
        <el-input
            v-model="registerFormData.username"
            size="large"
            :placeholder="t('forms.register.usernamePlaceholder')"
        />
      </el-form-item>
      <el-form-item :label="t('forms.register.password')" prop="password">
        <el-input
            v-model="registerFormData.password"
            show-password
            size="large"
            type="password"
            :placeholder="t('forms.register.passwordPlaceholder')"
        />
      </el-form-item>
      <el-form-item :label="t('forms.register.confirmPassword')">
        <el-input
            v-model="repeatPassword"
            show-password
            size="large"
            type="password"
            :placeholder="t('forms.register.confirmPasswordPlaceholder')"
        />
        <el-text v-if="registerFormData.password!==repeatPassword">{{ t('forms.register.passwordMismatch') }}</el-text>
      </el-form-item>
      <el-form-item :label="t('forms.register.email')" prop="email">
        <el-input
            v-model="registerFormData.email"
            size="large"
            :placeholder="t('forms.register.emailPlaceholder')"
        />
      </el-form-item>
      <el-form-item prop="captcha">
        <div class="captcha">
          <el-input
              v-model="emailRequest.captcha"
              :placeholder="t('forms.register.imageCaptchaPlaceholder')"
              size="large"
              maxlength="6"
              minlength="6"
          />
          <el-image :src="picPath" alt="" @click="emailVerify"/>
        </div>
      </el-form-item>
      <el-form-item :label="t('forms.register.emailCaptcha')" prop="verification_code">
        <div class="email-code">
          <el-input
              v-model="registerFormData.verification_code"
              size="large"
              :placeholder="t('forms.register.emailCaptchaPlaceholder')"
          />
          <el-button @click="sendCode">{{ t('forms.register.sendCode') }}</el-button>
        </div>
      </el-form-item>
      <el-form-item>
        <el-button
            type="primary"
            size="large"
            @click="submitForm"
        >{{ t('forms.register.submit') }}
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref,} from "vue";
import type {RegisterRequest} from "@/api/user";
import {useUserStore} from "@/stores/user";
import {captcha, type EmailRequest, sendEmailVerificationCode} from "@/api/base";
import type {FormInstance, FormRules} from 'element-plus';
import {useLayoutStore} from "@/stores/layout";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const registerForm = ref<FormInstance>()

const registerFormData = reactive<RegisterRequest>({
  username: "",
  password: "",
  email: "",
  verification_code: "",
})

const repeatPassword = ref('')

const emailRequest = reactive<EmailRequest>({
  email: "",
  captcha: "",
  captcha_id: "",
})


const rules = reactive<FormRules<RegisterRequest>>({
  username:[{
    required:true,
    max:20,
    trigger:'blur',
    message:t('forms.register.usernameMax')
  }],
  password:[{
    required:true,
    min:8,
    max:20,
    trigger:'change',
    message:t('forms.register.passwordLength')
  }],
  email: [{
    required: true,
    type:'email',
    trigger: 'blur',
    message: t('forms.register.emailFormat')
  }],
  verification_code: [{
    required: true,
    len: 6,
    trigger: 'blur',
    message: t('forms.register.codeLength')
  }],
})

const userStore = useUserStore()
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
  emailRequest.email = registerFormData.email
  const res = await sendEmailVerificationCode(emailRequest)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  } else {
    emailVerify()
  }
}

const submitForm = async () => {
  const isValid: boolean = await new Promise((resolve) => {
    registerForm.value?.validate((valid: boolean) => {
      resolve(valid)
    })
  })

  if (isValid) {
    const flag = await userStore.registerIn(registerFormData)

    if (flag) {
      layoutStore.state.registerVisible = false
      layoutStore.state.popoverVisible = false
    }
  }
  return false
}
</script>


<style scoped lang="scss">
.register-form {
  .el-form {

    .captcha {
      display: flex;
      gap: var(--sp-2);
      align-items: center;
      width: 100%;

      .el-input {
        flex: 1;
        min-width: 0;
      }

      .el-image {
        height: 40px;
        border-radius: var(--radius-sm);
        overflow: hidden;
        cursor: pointer;
        flex-shrink: 0;
      }
    }

    .email-code {
      display: flex;
      gap: var(--sp-2);
      align-items: center;
      width: 100%;

      .el-input {
        flex: 1;
        min-width: 0;
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
.register-form .el-form .captcha .el-input__wrapper {
  height: 40px;
}

/* 密码可见性切换图标：色渡 150ms + 按压反馈（§8） */
.register-form .el-input__password {
  cursor: pointer;
  transition: color 150ms ease-out, transform 80ms ease-out;

  &:active {
    transform: scale(0.92);
  }
}
</style>
