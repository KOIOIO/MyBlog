<template>
  <div id="app">
    <el-config-provider :locale="epLocale">
      <!-- §4.4 路由过渡：进入 opacity 150ms ease-out / 离开 80ms 淡出（CSS 在 base.css .page-*） -->
      <router-view v-slot="{Component}">
        <transition name="page" mode="out-in">
          <component :is="Component"/>
        </transition>
      </router-view>
    </el-config-provider>
  </div>
</template>

<script setup lang="ts">
import {computed} from 'vue'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'
import zhTw from 'element-plus/es/locale/lang/zh-tw'
import {i18n} from '@/i18n'

const epLocale = computed(() => {
  const locale = (i18n.global.locale as {value: string}).value
  switch (locale) {
    case 'en':
      return en
    case 'zh-TW':
      return zhTw
    default:
      return zhCn
  }
})
</script>

<style scoped>
</style>
