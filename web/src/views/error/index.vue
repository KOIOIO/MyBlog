<template>

  <web-navbar :noScroll="true"/>
  <div class="error">
    <div class="error-inner">
      <div class="error-code">404</div>
      <div class="error-text">{{ t('system.error.notFound') }}</div>
      <router-link class="error-link" :to="{ name: 'index' }">{{ t('system.error.backHome') }}</router-link>
    </div>
  </div>

  <web-footer/>
</template>

<script setup lang="ts">

import WebNavbar from "@/components/layout/WebNavbar.vue";
import WebFooter from "@/components/layout/WebFooter.vue";
import {useI18n} from "vue-i18n";

const {t} = useI18n()
</script>

<style scoped lang="scss">
.error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  padding: var(--sp-9) var(--sp-5);
  background: var(--bg);

  .error-inner {
    text-align: center;
  }

  .error-code {
    font-family: var(--font-serif);
    font-size: var(--fs-48);
    font-weight: 700;
    line-height: 1;
    color: var(--text-muted);
    margin-bottom: var(--sp-4);
  }

  .error-text {
    font-size: var(--fs-16);
    color: var(--text-body);
    margin-bottom: var(--sp-6);
  }

  .error-link {
    font-size: var(--fs-14);
    color: var(--accent);
  }
}
</style>
