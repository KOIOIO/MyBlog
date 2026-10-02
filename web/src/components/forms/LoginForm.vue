<template>
  <div class="login-form">
    <el-form
        ref="loginForm"
        :model="loginFormData"
        :rules="rules"
        :validate-on-rule-change="false"
        hide-required-asterisk
        @keyup.enter="submitForm"
    >
      <el-form-item :label="t('forms.login.email')" prop="email">
        <el-input
            v-model="loginFormData.email"
            size="large"
            :placeholder="t('forms.login.emailPlaceholder')"
            suffix-icon="user"
        />
      </el-form-item>
      <el-form-item :label="t('forms.login.password')" prop="password">
        <el-input
            v-model="loginFormData.password"
            show-password
            size="large"
            type="password"
            :placeholder="t('forms.login.passwordPlaceholder')"
        />
      </el-form-item>
      <el-form-item :label="t('forms.login.captcha')" prop="captcha">
        <div class="captcha">
          <el-input
              v-model="loginFormData.captcha"
              :placeholder="t('forms.login.captchaPlaceholder')"
              size="large"
          />
          <el-image :src="picPath" alt="" @click="loginVerify"/>
        </div>
      </el-form-item>
      <el-form-item>
        <el-button
            type="primary"
            size="large"
            @click="submitForm"
        >{{ t('forms.login.submit') }}
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref,} from "vue";
import type {LoginRequest} from "@/api/user";
import {useUserStore} from "@/stores/user";
import {captcha} from "@/api/base";
import type {FormInstance, FormRules} from 'element-plus';
import {useLayoutStore} from "@/stores/layout";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const loginForm = ref<FormInstance>()

const loginFormData = reactive<LoginRequest>({
  email: "",
  password: "",
  captcha: "",
  captcha_id: "",
})

const rules = reactive<FormRules<LoginRequest>>({
  email: [{
    required: true,
    type: 'email',
    trigger: 'blur',
    message: t('forms.login.emailFormat')
  }],
  password: [{
    required: true,
    min: 8,
    max: 20,
    trigger: 'change',
    message: t('forms.login.passwordLength')
  }],
  captcha: [{
    required: true,
    len: 6,
    trigger: 'blur',
    message: t('forms.login.captchaLength')
  }],
})

const userStore = useUserStore()
const layoutStore = useLayoutStore()

const picPath = ref('')

const loginVerify = () => {
  captcha().then(async (res) => {
    picPath.value = res.data.pic_path
    loginFormData.captcha_id = res.data.captcha_id
  })
}

loginVerify()

const submitForm = async () => {
  const isValid: boolean = await new Promise((resolve) => {
    loginForm.value?.validate((valid: boolean) => {
      resolve(valid)
    })
  })

  if (isValid) {
    const flag = await userStore.loginIn(loginFormData)

    if (!flag) {
      loginVerify()
    } else {
      layoutStore.state.loginVisible = false
      layoutStore.state.popoverVisible = false
    }
  }
  loginVerify()
  return false
}
</script>


<style scoped lang="scss">
.login-form {
  .el-form {

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
.login-form .el-form .captcha .el-input__wrapper {
  height: 40px;
}

/* 密码可见性切换图标：色渡 150ms + 按压反馈（§8） */
.login-form .el-input__password {
  cursor: pointer;
  transition: color 150ms ease-out, transform 80ms ease-out;

  &:active {
    transform: scale(0.92);
  }
}
</style>
