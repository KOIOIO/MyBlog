<template>

  <web-navbar :noScroll="true"/>
  <div class="error">
  <el-image src="/image/404.jpg" alt=""/>
  </div>

  <web-footer/>
</template>

<script setup lang="ts">

import WebNavbar from "@/components/layout/WebNavbar.vue";
import WebFooter from "@/components/layout/WebFooter.vue";
</script>

<style scoped lang="scss">
.error{
  display: flex;
  background-color: #1e1e2f; /* 更改背景颜色为深蓝色 */
  .el-image {
    margin-top: 80px;
    margin-left: auto;
    margin-right: auto;
    height: 500px;
    width: auto;
    border: 2px solid #6a0dad; /* 添加紫色边框 */
    box-shadow: 0 0 20px #6a0dad; /* 添加紫色阴影 */
  }
}

</style>
