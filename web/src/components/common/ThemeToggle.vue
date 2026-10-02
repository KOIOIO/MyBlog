<template>
  <button class="theme-toggle" :aria-label="isDark ? t('auth.themeToggleToLight') : t('auth.themeToggleToDark')" @click="toggle">
    <span class="theme-icon" :class="{ switching: animating }">
      <el-icon :size="18">
        <Sunny v-if="!isDark"/>
        <Moon v-else/>
      </el-icon>
    </span>
  </button>
</template>

<script setup lang="ts">
import {ref} from "vue";
import {Sunny, Moon} from "@element-plus/icons-vue";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const isDark = ref(document.documentElement.getAttribute("data-theme") === "dark")
const animating = ref(false)
let removeTimer: ReturnType<typeof setTimeout> | undefined

const toggle = () => {
  // §4.5：setAttribute 前挂 theme-switching class（全局色值过渡 250ms），300ms 后摘
  const root = document.documentElement
  root.classList.add("theme-switching")
  animating.value = true

  isDark.value = !isDark.value
  root.setAttribute("data-theme", isDark.value ? "dark" : "light")
  localStorage.setItem("myblog-theme", isDark.value ? "dark" : "light")

  if (removeTimer) clearTimeout(removeTimer)
  removeTimer = setTimeout(() => {
    root.classList.remove("theme-switching")
    animating.value = false
  }, 300)
}
</script>

<style scoped lang="scss">
.theme-toggle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: var(--radius-sm);
  background: transparent;
  color: var(--text-body);
  cursor: pointer;
  transition: background-color 150ms ease-out, color 150ms ease-out;

  &:hover {
    background: var(--accent-weak);
    color: var(--accent);
  }

  &:focus-visible {
    outline: 2px solid var(--accent);
    outline-offset: 2px;
  }

  .theme-icon {
    display: inline-flex;
    transition: transform 200ms ease-out, opacity 200ms ease-out;
  }

  /* §4.5：切换时轻微旋转/缩放换场 */
  .theme-icon.switching {
    animation: theme-icon-swap 200ms ease-out;
  }
}

@keyframes theme-icon-swap {
  0% { transform: rotate(-90deg) scale(0.6); opacity: 0; }
  100% { transform: rotate(0deg) scale(1); opacity: 1; }
}
</style>
