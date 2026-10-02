<template>
  <div class="breadcrumb">
    <el-breadcrumb separator="/">
      <template v-for="item in route.matched">
        <el-breadcrumb-item v-if="item.components" :to="{ path: item.path }"> {{ crumbTitle(item) }}</el-breadcrumb-item>
        <el-breadcrumb-item v-else>{{ crumbTitle(item) }}</el-breadcrumb-item>
      </template>
    </el-breadcrumb>
  </div>
</template>

<script setup lang="ts">
import {useRoute, type RouteRecordNameGeneric} from "vue-router";
import {useI18n} from "vue-i18n";

const route = useRoute()
const {t} = useI18n()

// 路由 meta.title 由路由文件（分片B）定义为硬编码中文，此处按路由 name 映射到 i18n key；
// 未收录的路由（如系统配置子页）回退显示原始 meta.title，保证其它分片页面不空白。
const breadcrumbKeyMap: Record<string, string> = {
  dashboard: 'breadcrumb.root',
  home: 'menu.home',
  'user-center': 'menu.userCenter.title',
  'user-info': 'menu.userCenter.info',
  'user-star': 'menu.userCenter.star',
  'user-comment': 'menu.userCenter.comment',
  'user-feedback': 'menu.userCenter.feedback',
  users: 'menu.users.title',
  'user-list': 'menu.users.list',
  articles: 'menu.articles.title',
  'article-publish': 'menu.articles.publish',
  'comment-list': 'menu.articles.commentList',
  'article-list': 'menu.articles.list',
  'dashboard-forum': 'menu.forum.title',
  'dashboard-forum-comments': 'menu.forum.comments',
  images: 'menu.images.title',
  'image-list': 'menu.images.list',
  system: 'menu.system.title',
  'feedback-list': 'menu.system.feedback',
  'advertisement-list': 'menu.system.advertisement',
  'friend-link-list': 'menu.system.friendLink',
  'login-logs': 'menu.system.loginLogs',
  'app-config': 'menu.system.appConfig',
  // 应用配置子页（系统配置 / 网站配置 / 邮箱 / 七牛 / JWT / 高德）
  'site-config': 'menu.config.site',
  'system-config': 'menu.config.system',
  'email-config': 'menu.config.email',
  'qiniu-config': 'menu.config.qiniu',
  'jwt-config': 'menu.config.jwt',
  'gaode-config': 'menu.config.gaode',
}

// RouteLocationMatched.name 为 RouteRecordNameGeneric（string | symbol | null），
// 此处放宽参数类型以兼容；仅在 name 存在且命中映射表时取 i18n 文案，否则回退 meta.title。
const crumbTitle = (item: { name?: RouteRecordNameGeneric | null; meta?: { title?: string } }): string => {
  const key = item.name != null ? breadcrumbKeyMap[String(item.name)] : undefined
  return key ? t(key) : (item.meta?.title ?? '')
}
</script>

<style scoped lang="scss">
.breadcrumb {
  display: flex;
  align-items: center;

  :deep(.el-breadcrumb) {
    font-size: var(--fs-14);
    padding: 0;
    line-height: 1;

    .el-breadcrumb__inner {
      color: var(--text-muted);
      font-weight: 400;
      transition: color 150ms var(--ease-out);
    }

    .el-breadcrumb__item .el-breadcrumb__inner a,
    .el-breadcrumb__item.is-link .el-breadcrumb__inner:hover {
      transition: color 150ms var(--ease-out);
    }

    .el-breadcrumb__inner a:focus-visible {
      outline: 2px solid var(--accent);
      outline-offset: 2px;
    }

    .el-breadcrumb__item:last-child .el-breadcrumb__inner {
      color: var(--text-primary);
      font-weight: 500;
    }

    .el-breadcrumb__separator {
      color: var(--border);
      margin: 0 var(--sp-2);
    }
  }
}
</style>
