<template>
  <div class="system-config">
    <el-container>
      <el-aside width="160px">
        <el-menu :collapse-transition="false" :router="true" :default-active="$route.path">
          <template v-for="item in menuList">
            <el-menu-item :index="generatePathForSingleItem(item)">
              <span>{{ item.title }}</span>
            </el-menu-item>
          </template>
        </el-menu>
      </el-aside>
      <el-main>
        <router-view/>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import {computed} from "vue";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

interface MenuItem {
  title: string;
  name: string;
}

const menuList = computed<MenuItem[]>(() => [
  {
    title: t('system.appConfig.site'),
    name: "site-config",
  },
  {
    title: t('system.appConfig.system'),
    name: "system-config",
  },
  {
    title: t('system.appConfig.email'),
    name: "email-config",
  },
  {
    title: t('system.appConfig.qiniu'),
    name: "qiniu-config",
  },
  {
    title: t('system.appConfig.jwt'),
    name: "jwt-config",
  },
  {
    title: t('system.appConfig.gaode'),
    name: "gaode-config",
  }
])

function generatePathForSingleItem(item: MenuItem): string {
  return '/dashboard/system/app-config/' + item.name;
}
</script>

<style scoped lang="scss">
.system-config {
  color: var(--text-body);
}

.el-container {
  background: transparent;
}

.el-aside {
  border-right: 1px solid var(--border);
}

.el-menu {
  height: 100%;
  background: transparent;
  border-right: none;
}

.el-menu-item {
  color: var(--text-body);
  border-radius: var(--radius-sm);
  margin: 2px var(--sp-2);
  transition: background-color 150ms ease-out, color 150ms ease-out;

  &:hover {
    background-color: var(--accent-weak);
    color: var(--accent);
  }
}

.el-menu-item.is-active {
  background-color: var(--accent-weak);
  color: var(--accent);
}

.el-main {
  background: transparent;
  padding: var(--sp-5);
}
</style>
