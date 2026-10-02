<template>
  <div class="dashboard-tag">
    <el-tag
        v-for="tag in tags"
        :closable="tag.name!=='home'"
        :key="tag.name"
        :effect="route.name === tag.name?'light':'plain'"
        type="info"
        @close="handleClose(tag.name)"
        @click="handleTag(tag.name)"
    >
      {{ tag.title }}
    </el-tag>
    <el-button class="close-button" size="small" @click="closeAllTags">
      {{ t('tag.closeAll') }}
    </el-button>
  </div>
</template>

<script lang="ts" setup>
import {useTagStore} from "@/stores/tag";
import {useRoute, useRouter} from "vue-router";
import {computed} from "vue";
import {useI18n} from "vue-i18n";

const {t} = useI18n()
const route = useRoute()
const store = useTagStore()
const router = useRouter()

const tags = computed(() => store.state.tags)
const handleClose = (tagName: string) => {
  const tags = store.state.tags;

  // 如果要删除的 tag 不是当前路由，则直接删除
  if (tagName !== route.name) {
    store.state.tags = tags.filter(tag => tag.name !== tag.name);
    return;
  }

  // 如果要删除的 tag 是当前路由，先找到要删除的 tag 的索引
  const index = tags.findIndex(tag => tag.name === tag.name);

  // 如果找到了该 tag
  if (index !== -1) {
    // 先删除该 tag
    store.state.tags = tags.filter(tag => tag.name !== tag.name);

    // 计算要跳转的上一个 tag 的名称
    const previousTag = index > 0 ? tags[index - 1].name : null;

    // 跳转到上一个 tag 或默认路由
    if (previousTag) {
      router.push({name: previousTag});
    } else {
      router.push({name: 'home'});
    }
  }
}

const handleTag = (tagName: string) => {
  router.push({name: tagName})
}


const closeAllTags = () => {
  store.state.tags = [
    {
      title: t("menu.home"),
      name: "home"
    }
  ];
  router.push({name: 'home'});
}
</script>

<style scoped lang="scss">
.dashboard-tag {
  display: flex;
  align-items: center;
  gap: var(--sp-1);
  padding: 0 var(--sp-5);
  border-bottom: 1px solid var(--border);
  overflow-x: auto;

  .el-tag {
    font-size: var(--fs-14);
    margin-right: 0;
    height: 36px;
    line-height: 34px;
    padding: 0 var(--sp-3);
    border-radius: var(--radius-sm) var(--radius-sm) 0 0;
    border: none;
    border-bottom: 2px solid transparent;
    background: transparent;
    color: var(--text-muted);
    cursor: pointer;
    transition: color 150ms ease-out, border-color 150ms ease-out, background-color 150ms ease-out;

    &:hover {
      color: var(--text-primary);
      background: var(--accent-weak);
    }

    &.is-light {
      color: var(--accent);
      border-bottom-color: var(--accent);
      background: transparent;
    }

    :deep(.el-tag__close) {
      color: var(--text-muted);

      &:hover {
        background: var(--accent-weak);
        color: var(--accent);
      }
    }
  }

  .close-button {
    margin-left: auto;
    flex-shrink: 0;
  }
}
</style>
