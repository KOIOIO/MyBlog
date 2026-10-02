<template>
  <div class="news-page">
    <web-navbar :noScroll="true"/>
    <div class="page">
      <el-tabs v-model="activeTab" @tab-click="handleNewsTabClick">
        <el-tab-pane v-for="item in newsTypeList" :key="item.name" :name="item.name">
          <template #label>
            <el-image :src="item.src" alt="" class="tab-icon"></el-image>
            <span>{{ item.label }}</span>
          </template>
          <div class="news-list" :key="activeTab">
            <div v-for="(row, idx) in newsTableData" :key="idx" class="news-item"
                 :style="{ animationDelay: Math.min(idx, 7) * 60 + 'ms' }"
                 @click="handleNewsTableClick(row)">
              <span class="item-index">{{ idx + 1 }}</span>
              <img v-if="row.image" class="item-cover" :src="row.image" referrerpolicy="no-referrer" alt=""/>
              <div class="item-body">
                <h3 class="item-title">{{ row.title }}</h3>
                <p class="item-desc">{{ row.description }}</p>
              </div>
              <span v-if="row.popularity" class="item-popularity">{{ row.popularity }}</span>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import WebNavbar from "@/components/layout/WebNavbar.vue";
import {computed, ref} from "vue";
import {type HotItem, websiteNews, type WebsiteNewsRequest} from "@/api/website";
import type {TabsPaneContext} from "element-plus";
import {useI18n} from "vue-i18n";

const {t} = useI18n();

const newsTableData = ref<HotItem[]>()
const activeTab = ref("baidu")

interface newsTypeItem {
  name: string;
  label: string;
  src: string;
}

const newsTypeList = computed<newsTypeItem[]>(() => [
  {
    name: "baidu",
    label: t('pages.news.baidu'),
    src: "/image/baidu.png"
  },
  {
    name: "bilibili",
    label: t('pages.news.bilibili'),
    src: "/image/bilibili.png"
  },
  {
    name: "kuaishou",
    label: t('pages.news.kuaishou'),
    src: "/image/kuaishou.png"
  },
  {
    name: "toutiao",
    label: t('pages.news.toutiao'),
    src: "/image/toutiao.png"
  }
])

const handleNewsTabClick = (tab: TabsPaneContext, _: Event) => {
  getNewsTableData(tab.paneName as string)
}

let newsMap = new Map<string, HotItem[]>()

const getNewsTableData = async (source: string) => {
  if (!newsMap.has(source)) {
    const newsRequest: WebsiteNewsRequest = {
      source: source
    }
    const res = await websiteNews(newsRequest)
    if (res.code === 0) {
      newsMap.set(source, res.data.hot_list)
    }
  }
  newsTableData.value = newsMap.get(source)
}

getNewsTableData("baidu")

const handleNewsTableClick = (item: HotItem) => {
  window.open(item.url)
}
</script>

<style scoped lang="scss">
.news-page {
  background-color: var(--bg);
  min-height: 100vh;

  .page {
    max-width: var(--reading-width);
    margin: 0 auto;
    padding: calc(70px + var(--sp-6)) var(--sp-4) var(--sp-9);
  }

  :deep(.el-tabs__item) {
    .tab-icon {
      width: 18px;
      height: 18px;
      margin-right: var(--sp-1);
      vertical-align: middle;
    }
  }

  .news-list {
    border-top: 1px solid var(--border);
  }

  .news-item {
    display: flex;
    align-items: flex-start;
    gap: var(--sp-4);
    padding: var(--sp-4) var(--sp-2);
    border-bottom: 1px solid var(--border);
    cursor: pointer;
    transition: background-color 150ms ease-out, transform 80ms ease-out;
    animation: kf-fade-up var(--dur-mid) var(--ease-out) backwards;

    &:hover {
      background-color: var(--bg-elevated);
    }

    &:active {
      transform: translateY(1px);
    }

    .item-index {
      flex-shrink: 0;
      width: 24px;
      font-size: var(--fs-14);
      font-weight: 600;
      color: var(--text-muted);
      font-family: var(--font-mono);
      padding-top: 2px;
    }

    .item-cover {
      flex-shrink: 0;
      width: 120px;
      height: 80px;
      object-fit: cover;
      border-radius: var(--radius-sm);
    }

    .item-body {
      flex: 1;
      min-width: 0;

      .item-title {
        font-size: var(--fs-16);
        font-weight: 600;
        color: var(--text-primary);
        line-height: var(--lh-title);
        margin: 0 0 var(--sp-1);
        transition: color 150ms ease-out;
      }

      &:hover .item-title {
        color: var(--accent);
      }

      .item-desc {
        font-size: var(--fs-14);
        color: var(--text-muted);
        line-height: var(--lh-body);
        margin: 0;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }
    }

    .item-popularity {
      flex-shrink: 0;
      font-size: var(--fs-12);
      color: var(--text-muted);
      font-family: var(--font-mono);
    }
  }
}

</style>
