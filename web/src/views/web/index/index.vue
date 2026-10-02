<template>
  <div class="index">
    <web-navbar/>
    <div class="page">
      <article-list/>
    </div>
  </div>
</template>

<script setup lang="ts">
import {useI18n} from "vue-i18n";
import WebNavbar from "@/components/layout/WebNavbar.vue";
import ArticleList from "@/components/pages/ArticleList.vue";

const {t} = useI18n();
</script>

<style scoped lang="scss">
.index {
  background-color: var(--bg);
  min-height: 100vh;

  .page {
    max-width: var(--content-width);
    margin: 0 auto;
    padding: calc(70px + var(--sp-6)) var(--sp-4) var(--sp-9);
    animation: kf-fade-up var(--dur-mid) var(--ease-out) backwards;
  }
}

</style>
