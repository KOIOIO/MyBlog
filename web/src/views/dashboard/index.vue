<template>
  <div class="dashboard">
    <el-container>
      <el-aside :class="{collapsed: isCollapse}">
        <logo/>
        <dashboard-menu/>
      </el-aside>
      <el-container>
        <el-header>
          <div class="header-top">
            <breadcrumb/>
            <div class="header-top-right">
              <auth-popover/>
            </div>
          </div>
          <dashboard-tag/>
        </el-header>
        <el-main>
          <router-view/>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import Logo from "@/components/widgets/Logo.vue";
import Breadcrumb from "@/components/layout/Breadcrumb.vue";
import DashboardMenu from "@/components/layout/DashboardMenu.vue";
import {computed} from "vue";
import DashboardTag from "@/components/layout/DashboardTag.vue";
import AuthPopover from "@/components/common/AuthPopover.vue";
import {useLayoutStore} from "@/stores/layout";

const store = useLayoutStore()
const isCollapse = computed(() => store.state.isCollapse)
</script>

<style scoped lang="scss">
.dashboard {
  display: flex;
  min-height: 100vh;
  background-color: var(--bg);

  .el-container {
    background-color: transparent;
  }

  .el-aside {
    width: 240px;
    height: 100vh;
    background-color: var(--bg-elevated);
    border-right: 1px solid var(--border);
    transition: width 200ms ease;
    overflow-x: hidden;
    &::-webkit-scrollbar{
      display: none;
    }
  }

  .el-aside.collapsed {
    width: 64px;
  }

  .el-header {
    height: auto;
    padding: 0;
    background-color: var(--bg-elevated);
    border-bottom: 1px solid var(--border);

    .header-top {
      display: flex;
      align-items: center;
      height: 56px;
      padding: 0 var(--sp-5);

      .header-top-right {
        margin-left: auto;
        margin-top: auto;
        margin-bottom: auto;
      }
    }
  }

  .el-main{
    height: calc(100vh - 56px);
    padding: var(--sp-5);
    background-color: var(--bg);
    &::-webkit-scrollbar{
      display: none;
    }
  }
}
</style>
