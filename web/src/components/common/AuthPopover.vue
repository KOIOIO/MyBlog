<template>
  <div class="auth-popover">
    <el-popover v-if="isLoggedIn" width="280" trigger="click">
      <template #reference>
        <el-avatar :size="36" :src="userStore.state.userInfo.avatar"/>
      </template>
      <template #default>
        <div class="title">
          <el-row>{{ t('auth.welcomeBack', {name: userStore.state.userInfo.username}) }}</el-row>
        </div>
        <div class="user-info">
          <el-avatar :size="50" :src="userStore.state.userInfo.avatar"/>
          <el-row>{{ userStore.state.userInfo.username }}</el-row>
          <el-row>{{ userStore.state.userInfo.address }}</el-row>
        </div>
        <el-row class="uuid">{{ t('auth.uuid') }}：{{ userStore.state.userInfo.uuid }}</el-row>
        <div class="container">
          <el-row>{{ t('auth.signature') }}：{{ userStore.state.userInfo.signature }}</el-row>
          <div class="action-button">
            <el-button @click="userStore.loginOut()">{{ t('auth.logout') }}</el-button>
            <el-button @click="goIndexOrToDashboard">{{ label }}</el-button>
          </div>
        </div>
      </template>
    </el-popover>

    <el-popover v-else :visible="popoverVisible" width="280">
      <template #reference>
        <button class="login-btn" @click="layoutStore.state.popoverVisible = !layoutStore.state.popoverVisible">{{ t('auth.login') }}</button>
      </template>
      <template #default>
        <div class="default" v-click-outside="onClickOutside">
          <div class="title">
            <el-row>{{ t('auth.loginDesc') }}</el-row>
          </div>
          <div class="auth-button">
            <div class="button-item">
              <el-button class="login" icon="User" @click="layoutStore.state.loginVisible = true"/>
              {{ t('auth.login') }}
              <el-dialog
                  v-model="loginVisible"
                  width="500"
                  destroy-on-close
                  align-center
                  :before-close="loginVisibleSynchronization"
              >
                <template #header>
                  {{ t('auth.login') }}
                </template>
                <login-form/>
                <template #footer>

                </template>
              </el-dialog>
            </div>

            <div class="button-item">
              <el-button class="register" icon="Edit" @click="layoutStore.state.registerVisible = true"/>
              {{ t('auth.register') }}
              <el-dialog
                  v-model="registerVisible"
                  width="640"
                  destroy-on-close
                  align-center
                  :before-close="registerVisibleSynchronization"
              >
                <template #header>
                  {{ t('auth.register') }}
                </template>
                <register-form/>
                <template #footer>

                </template>
              </el-dialog>
            </div>

            <div class="button-item">
              <el-button class="forgot-password" icon="Unlock" @click="layoutStore.state.forgotPasswordVisible = true"/>
              {{ t('auth.forgotPassword') }}
              <el-dialog
                  v-model="forgotPasswordVisible"
                  width="570"
                  destroy-on-close
                  align-center
                  :before-close="forgotPasswordVisibleSynchronization"
              >
                <template #header>
                  {{ t('auth.forgotPassword') }}
                </template>
                <forgot-password-form/>
                <template #footer>

                </template>
              </el-dialog>
            </div>
          </div>
        </div>
      </template>
    </el-popover>
  </div>
</template>

<script setup lang="ts">
import {computed, ref, watch} from 'vue'
import LoginForm from "@/components/forms/LoginForm.vue";
import {useLayoutStore} from "@/stores/layout";
import {ClickOutside as vClickOutside} from 'element-plus'
import {useUserStore} from "@/stores/user";
import router from "@/router";
import {useRoute} from "vue-router";
import RegisterForm from "@/components/forms/RegisterForm.vue";
import ForgotPasswordForm from "@/components/forms/ForgotPasswordForm.vue";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const userStore = useUserStore()
const route = useRoute()
const isLoggedIn = computed(() => userStore.isLoggedIn)
const isInDashboard = route.matched.some(record => record.name === 'dashboard')
const label = computed(() => isInDashboard ? t('auth.backHome') : t('auth.toDashboard'))
const goIndexOrToDashboard = (() => {
  if (isInDashboard) {
    router.push({name: 'index'})
  } else {
    router.push({name: 'home'})
  }
})

const layoutStore = useLayoutStore()

const popoverVisible = ref(layoutStore.state.popoverVisible)
watch(
    () => layoutStore.state.popoverVisible,
    (newValue) => {
      popoverVisible.value = newValue;
    }
)

const onClickOutside = () => {
  layoutStore.state.popoverVisible = false
}

const loginVisible = ref(layoutStore.state.loginVisible)
watch(
    () => layoutStore.state.loginVisible,
    (newValue) => {
      loginVisible.value = newValue
    }
)

const loginVisibleSynchronization = () => {
  layoutStore.state.loginVisible = false
}

const registerVisible = ref(layoutStore.state.registerVisible)
watch(
    () => layoutStore.state.registerVisible,
    (newValue) => {
      registerVisible.value = newValue
    }
)

const registerVisibleSynchronization = () => {
  layoutStore.state.registerVisible = false
}

const forgotPasswordVisible = ref(layoutStore.state.forgotPasswordVisible)
watch(
    () => layoutStore.state.forgotPasswordVisible,
    (newValue) => {
      forgotPasswordVisible.value = newValue
    }
)

const forgotPasswordVisibleSynchronization = () => {
  layoutStore.state.forgotPasswordVisible = false
}
</script>

<style scoped lang="scss">
.auth-popover {
  --el-text-color-disabled: transparent;

  .login-btn {
    height: 34px;
    padding: 0 var(--sp-4);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    background-color: transparent;
    color: var(--text-body);
    font-size: var(--fs-14);
    cursor: pointer;
    transition: color 150ms ease-out, border-color 150ms ease-out, background-color 150ms ease-out, transform 80ms ease-out;

    &:hover {
      color: var(--accent);
      border-color: var(--accent);
      background-color: var(--accent-weak);
    }

    &:active {
      transform: translateY(1px);
    }

    &:focus-visible {
      outline: 2px solid var(--accent);
      outline-offset: 2px;
    }
  }

  .el-avatar {
    cursor: pointer;
    transition: opacity 150ms ease-out;

    &:hover {
      opacity: 0.85;
    }
  }

  .el-avatar--icon {
    font-size: 24px;
  }
}
</style>

<style lang="scss">
.el-popover.el-popper {
  background-color: var(--bg-elevated);
  border: 1px solid var(--border);
  box-shadow: var(--shadow-sm);
  padding: 0;
  color: var(--text-body);

  .el-popper__arrow:before {
    background-color: var(--bg-elevated);
    border-color: var(--border);
  }

  .title {
    display: flex;
    height: 48px;
    border-bottom: 1px solid var(--border);

    .el-row {
      margin: auto auto;
      font-size: var(--fs-14);
      color: var(--text-primary);
    }
  }

  .user-info {
    text-align: center;
    padding: var(--sp-4) 0 var(--sp-2);

    .el-row {
      display: flex;
      justify-content: center;
      margin: var(--sp-1);
      font-size: var(--fs-14);
      color: var(--text-body);

      &:first-of-type {
        font-weight: 600;
        color: var(--text-primary);
      }
    }
  }

  .uuid {
    margin: var(--sp-2) var(--sp-5);
    font-size: var(--fs-12);
    color: var(--text-muted);
  }

  .container {
    background-color: transparent;

    > .el-row {
      padding: 0 var(--sp-5);
      font-size: var(--fs-14);
      color: var(--text-body);
    }

    .action-button {
      display: flex;
      justify-content: center;
      gap: var(--sp-2);
      padding: var(--sp-3) 0;

      .el-button {
        margin-bottom: 0;
      }
    }

    .el-divider {
      --el-bg-color: var(--border);
      font-size: var(--fs-12);
      color: var(--text-muted);
    }

  }

  .default {
    .title {
      display: flex;
      height: 48px;
      border-bottom: 1px solid var(--border);

      .el-row {
        margin: auto auto;
        font-size: var(--fs-14);
        color: var(--text-primary);
      }
    }

    .auth-button {
      display: flex;
      justify-content: center;
      padding: var(--sp-4) 0 var(--sp-2);

      .button-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 80px;
        gap: var(--sp-1);

        .el-button {
          border: none;
          --el-font-size-base: 24px;
          height: 48px;
          width: 48px;
          padding: 0;
          border-radius: var(--radius-sm);
          background-color: transparent;
          color: var(--text-body);
          transition: background-color 150ms ease-out, color 150ms ease-out, transform 80ms ease-out;

          &:hover {
            background-color: var(--accent-weak);
            color: var(--accent);
          }

          &:active {
            transform: translateY(1px);
          }

          &:focus-visible {
            outline: 2px solid var(--accent);
            outline-offset: 2px;
          }
        }

        font-size: var(--fs-12);
        color: var(--text-muted);
      }
    }
  }
}
</style>
