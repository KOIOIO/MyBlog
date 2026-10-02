<template>
  <div class="logo">
    <span class="brand-word">Folio<span class="brand-dot">.</span></span>
    <div class="brand-info">
      <div class="slogan">{{ websiteStore.state.websiteInfo.slogan }}</div>
      <div class="slogan-en">{{ websiteStore.state.websiteInfo.slogan_en }}</div>
    </div>
  </div>
</template>


<script setup lang="ts">
import {useWebsiteStore} from "@/stores/website";

const websiteStore = useWebsiteStore()
</script>

<style scoped lang="scss">
.logo {
  display: flex;
  align-items: center;

  .brand-word {
    font-size: var(--fs-20);
    font-weight: 700;
    color: var(--text-primary);
    letter-spacing: -0.01em;
    line-height: var(--lh-title);

    .brand-dot {
      color: var(--accent);
    }
  }

  .brand-info {
    display: flex;
    flex-direction: column;
    margin-left: var(--sp-3);

    .slogan {
      font-size: var(--fs-14);
      color: var(--text-primary);
      line-height: var(--lh-title);
    }

    .slogan-en {
      font-size: var(--fs-12);
      color: var(--text-muted);
      line-height: var(--lh-title);
    }
  }
}

/* 后台侧栏：Logo 区域 */
.dashboard .logo {
  padding: var(--sp-4) var(--sp-5);
  height: 56px;

  .brand-word {
    font-size: var(--fs-20);
  }

  .brand-info {
    display: flex;
  }
}

/* 后台折叠态：只留图标字 */
.dashboard .collapsed .logo {
  justify-content: center;
  padding: var(--sp-4) 0;

  .brand-info {
    display: none;
  }
}

/* 前台导航中的 Logo */
.web-navbar .logo {
  height: 56px;
  padding: 0;

  .brand-info {
    display: none;
  }
}
</style>
