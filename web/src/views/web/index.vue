<template>
  <div class="web">
    <el-container>
      <router-view/>
      <el-footer>
        <web-footer/>
      </el-footer>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import {useI18n} from "vue-i18n";
import WebFooter from "@/components/layout/WebFooter.vue";

const {t} = useI18n();
</script>

<style scoped lang="scss">
.web {
  min-height: 100vh;
  background-color: var(--bg);

  .el-container {
    flex-direction: column;
    background-color: transparent;
  }

  .el-footer {
    height: auto;
    padding: 0;
    background-color: transparent;
    color: var(--text-body);
  }
}
</style>
