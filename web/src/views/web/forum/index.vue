<template>
  <div class="forum-page">
    <web-navbar :noScroll="true"/>
    <div class="page">
      <div class="container">
        <main class="main">
          <header class="page-head">
            <div class="head-text">
              <h1 class="page-title">{{ t('pages.forum.title') }}</h1>
              <p class="page-desc">{{ t('pages.forum.desc') }}</p>
            </div>
            <el-button v-if="userStore.isLoggedIn" type="primary" @click="goPublish">{{ t('pages.forum.publishPost') }}</el-button>
            <el-button v-else type="primary" @click="requireLogin">{{ t('pages.forum.publishPost') }}</el-button>
          </header>

          <div class="filter-tabs">
            <span
                v-for="tab in tabs"
                :key="tab.value"
                class="tab-item"
                :class="{ active: category === tab.value }"
                @click="onTab(tab.value)"
            >{{ tab.label }}</span>
          </div>

          <div v-if="list.length" class="post-list">
            <PostCard v-for="post in list" :key="post.id" :post="post"/>
          </div>
          <div v-else class="empty">{{ loading ? t('common.loading') : t('pages.forum.emptyPost') }}</div>

          <div class="pager">
            <el-pagination
                v-if="total > pageSize"
                background
                layout="prev, pager, next"
                :total="total"
                :page-size="pageSize"
                :current-page="page"
                @current-change="onPageChange"
            />
          </div>
        </main>

        <TagCloudSidebar :selected-tag="selectedTag" @select="onTagSelect"/>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, ref} from "vue";
import {ElMessage} from "element-plus";
import WebNavbar from "@/components/layout/WebNavbar.vue";
import PostCard from "@/components/forum/PostCard.vue";
import TagCloudSidebar from "@/components/forum/TagCloudSidebar.vue";
import {forumList, type ForumPost} from "@/api/forum";
import {useUserStore} from "@/stores/user";
import {useLayoutStore} from "@/stores/layout";
import {useRouter} from "vue-router";
import {useI18n} from "vue-i18n";
import {categoryLabel} from "@/i18n/meta";

const {t} = useI18n();

const userStore = useUserStore();
const layoutStore = useLayoutStore();
const router = useRouter();

const tabs = computed(() => [
    {label: t('common.all'), value: ''},
    {label: categoryLabel('技术'), value: '技术'},
    {label: categoryLabel('生活'), value: '生活'},
]);

const list = ref<ForumPost[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = ref(10);
const category = ref('');
const selectedTag = ref('');
const loading = ref(false);

const loadList = async () => {
    loading.value = true;
    try {
        const res = await forumList({
            page: page.value,
            page_size: pageSize.value,
            category: category.value || undefined,
            tag: selectedTag.value || undefined,
        });
        if (res.code === 0) {
            list.value = res.data.list || [];
            total.value = res.data.total || 0;
        }
    } finally {
        loading.value = false;
    }
};

onMounted(loadList);

const onTab = (val: string) => {
    category.value = val;
    page.value = 1;
    loadList();
};

const onTagSelect = (tag: string) => {
    selectedTag.value = selectedTag.value === tag ? '' : tag;
    page.value = 1;
    loadList();
};

const onPageChange = (p: number) => {
    page.value = p;
    loadList();
    window.scrollTo({top: 0, behavior: 'smooth'});
};

const requireLogin = () => {
    ElMessage.warning(t('pages.forum.loginRequired'));
    layoutStore.state.popoverVisible = true;
    layoutStore.state.loginVisible = true;
};

const goPublish = () => {
    router.push({name: 'forum-publish'});
};
</script>

<style scoped lang="scss">
.forum-page {
  background-color: var(--bg);
  min-height: 100vh;

  .page {
    padding: calc(70px + var(--sp-6)) var(--sp-4) var(--sp-9);
  }

  .container {
    display: flex;
    justify-content: center;
    gap: var(--sp-7);
    max-width: var(--content-width);
    margin: 0 auto;
  }

  .main {
    flex: 1;
    min-width: 0;
  }

  .page-head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    margin-bottom: var(--sp-5);

    .page-title {
      font-size: var(--fs-24);
      font-weight: 600;
      color: var(--text-primary);
      margin: 0 0 var(--sp-1);
    }

    .page-desc {
      font-size: var(--fs-14);
      color: var(--text-muted);
      margin: 0;
    }
  }

  .filter-tabs {
    display: flex;
    gap: var(--sp-5);
    border-bottom: 1px solid var(--border);
    margin-bottom: var(--sp-2);

    .tab-item {
      position: relative;
      padding: var(--sp-2) 0;
      font-size: var(--fs-14);
      color: var(--text-body);
      cursor: pointer;
      transition: color 150ms ease-out;

      &::after {
        content: "";
        position: absolute;
        left: 0;
        bottom: -1px;
        width: 100%;
        height: 2px;
        background-color: var(--accent);
        transform: scaleX(0);
        transition: transform 180ms ease-out;
      }

      &:hover,
      &.active {
        color: var(--accent);

        &::after {
          transform: scaleX(1);
        }
      }
    }
  }

  .empty {
    padding: var(--sp-9) 0;
    text-align: center;
    font-size: var(--fs-14);
    color: var(--text-muted);
  }

  .pager {
    display: flex;
    justify-content: center;
    margin-top: var(--sp-6);
  }
}
</style>
