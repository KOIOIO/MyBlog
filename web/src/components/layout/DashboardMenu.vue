<template>
  <div class="dashboard-menu">
    <div class="collapse-button">
      <el-button @click="handleCollapse">
        <el-icon>
          <component :is="icon"></component>
        </el-icon>
      </el-button>
    </div>
    <el-menu :collapse="isCollapse" :collapse-transition="false" :router="true" :default-active="$route.path">
      <template v-for="item in menuList">
        <el-menu-item v-if="!item.subItems" :index="generatePathForSingleItem(item)">
          <el-icon>
            <component :is="item.icon"></component>
          </el-icon>
          <span>{{ item.title }}</span>
        </el-menu-item>
        <el-sub-menu v-else-if="!item.admin_role||userStore.isAdmin" :index="item.name">
          <template #title>
            <el-icon>
              <component :is="item.icon"></component>
            </el-icon>
            <span>{{ item.title }}</span>
          </template>
          <el-menu-item v-for="subItem in item.subItems" :index="generatePathForSubItem(item,subItem)"
                        @click="handleClick(subItem)">
            <el-icon>
              <component :is="subItem.icon"></component>
            </el-icon>
            <span>{{ subItem.title }}</span>
          </el-menu-item>
        </el-sub-menu>
      </template>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import type {Tag} from "@/stores/tag";
import {computed} from "vue";
import {useI18n} from "vue-i18n";
import {useLayoutStore} from "@/stores/layout";
import {useTagStore} from "@/stores/tag";
import {useUserStore} from "@/stores/user";

const {t} = useI18n()
const userStore = useUserStore()

interface MenuItem {
  title: string;
  name: string;
  icon: string;
  subItems?: MenuItem[];
  admin_role?:boolean;
}

const menuList = computed<MenuItem[]>(() => [
  {
    title: t("menu.home"),
    name: "",
    icon: "House"
  },
  {
    title: t("menu.userCenter.title"),
    name: "user-center",
    icon: "Monitor",
    subItems: [
      {
        title: t("menu.userCenter.info"),
        name: "user-info",
        icon: "Postcard"
      },
      {
        title: t("menu.userCenter.star"),
        name: "user-star",
        icon: "Star"
      },
      {
        title: t("menu.userCenter.comment"),
        name: "user-comment",
        icon: "ChatDotRound"
      },
      {
        title: t("menu.userCenter.feedback"),
        name: "user-feedback",
        icon: "Message"
      }
    ]
  },
  {
    title: t("menu.users.title"),
    name: "users",
    icon: "User",
    admin_role:true,
    subItems: [
      {
        title: t("menu.users.list"),
        name: "user-list",
        icon: "SetUp"
      }
    ]
  },
  {
    title: t("menu.articles.title"),
    name: "articles",
    icon: "Document",
    admin_role:true,
    subItems: [
      {
        title: t("menu.articles.publish"),
        name: "article-publish",
        icon: "Collection"
      },
      {
        title: t("menu.articles.commentList"),
        name: "comment-list",
        icon: "ChatLineRound"
      },
      {
        title: t("menu.articles.list"),
        name: "article-list",
        icon: "DocumentCopy"
      }
    ]
  },
  {
    title: t("menu.forum.title"),
    name: "forum",
    icon: "ChatDotRound",
    subItems: [
      {
        title: t("menu.forum.list"),
        name: "",
        icon: "Tickets"
      },
      {
        title: t("menu.forum.comments"),
        name: "comments",
        icon: "ChatLineRound"
      }
    ]
  },
  {
    title: t("menu.images.title"),
    name: "images",
    icon: "Picture",
    admin_role:true,
    subItems: [
      {
        title: t("menu.images.list"),
        name: "image-list",
        icon: "PictureRounded"
      }
    ]
  },
  {
    title: t("menu.system.title"),
    name: "system",
    icon: "Coin",
    admin_role:true,
    subItems: [
      {
        title: t("menu.system.feedback"),
        name: "feedback-list",
        icon: "Position"
      },
      {
        title: t("menu.system.advertisement"),
        name: "advertisement-list",
        icon: "Connection"
      },
      {
        title: t("menu.system.friendLink"),
        name: "friend-link-list",
        icon: "Link"
      },
      {
        title: t("menu.system.loginLogs"),
        name: "login-logs",
        icon: "memo"
      },
      {
        title: t("menu.system.appConfig"),
        name: "app-config",
        icon: "setting"
      }
    ]
  }
])

const layoutStore = useLayoutStore()
const tagStore = useTagStore()
const isCollapse = computed(() => layoutStore.state.isCollapse)
const icon = computed(() => layoutStore.state.isCollapse ? "Expand" : "Fold")

const handleCollapse = () => {
  layoutStore.state.isCollapse = !layoutStore.state.isCollapse
}

function generatePathForSingleItem(item: MenuItem): string {
  return '/dashboard/' + item.name;
}

function generatePathForSubItem(parentItem: MenuItem, subItem: MenuItem): string {
  if (!subItem.name) {
    return '/dashboard/' + parentItem.name;
  }
  return '/dashboard/' + parentItem.name + '/' + subItem.name;
}

function handleClick(subItem: MenuItem) {
  const newTag: Tag = {
    title: subItem.title,
    name: subItem.name
  }
  const exists = tagStore.state.tags.some(tag => tag.name === newTag.name);
  if (exists) {
    return;
  }
  tagStore.state.tags.push(newTag);
}
</script>

<style scoped lang="scss">
.dashboard-menu {
  display: flex;
  flex-direction: column;
  height: 100%;

  .collapse-button {
    display: flex;
    padding: var(--sp-2) var(--sp-3);

    .el-button {
      margin-left: auto;
      border: none;
      background: transparent;
      color: var(--text-muted);

      &:hover {
        background: var(--accent-weak);
        color: var(--accent);
      }
    }
  }

  .el-menu {
    flex: 1;
    border-right: none;
    background-color: transparent;
    --el-menu-item-height: 44px;
    --el-menu-sub-item-height: 40px;
    --el-menu-text-color: var(--text-body);
    --el-menu-hover-text-color: var(--accent);
    --el-menu-active-color: var(--accent);

    :deep(.el-menu-item),
    :deep(.el-sub-menu__title) {
      position: relative;
      height: 44px;
      line-height: 44px;
      color: var(--text-body);
      transition: background-color 150ms ease-out, color 150ms ease-out, transform 80ms ease-out;

      &:hover {
        background-color: var(--accent-weak);
        color: var(--accent);
      }

      &:active {
        transform: translateY(1px);
      }

      &:focus-visible {
        outline: 2px solid var(--accent);
        outline-offset: -2px;
      }
    }

    :deep(.el-menu-item.is-active) {
      background-color: var(--accent-weak);
      color: var(--accent);

      &::before {
        content: "";
        position: absolute;
        left: 0;
        top: 0;
        bottom: 0;
        width: 3px;
        background-color: var(--accent);
      }
    }
  }

  .el-popper {
    .el-menu-item {
      height: 44px;
      width: 140px;
    }
  }
}

.collapsed .dashboard-menu .collapse-button .el-button {
  margin-right: auto;
  margin-left: auto;
}
</style>

<style lang="scss">
.el-menu--popup {
  min-width: 140px;
}
</style>
