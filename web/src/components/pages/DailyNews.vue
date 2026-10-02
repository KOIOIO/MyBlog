<template>
  <el-card class="daily-news">
    <el-row class="title">{{ t('components.dailyNews.title') }}</el-row>
    <el-tabs model-value="baidu" @tab-click="handleNewsTabClick">
      <el-tab-pane v-for="item in newsTypeList" :name="item.name">
        <template #label>
          <el-image :src="item.src" alt="" class="tab-icon"></el-image>
          <span>{{ item.label }}</span>
        </template>
        <el-table
            :data="newsTableData"
            :row-style="{height: '64px'}"
        >
          <el-table-column prop="index" :label="t('components.dailyNews.indexLabel')" width="60"/>
          <el-table-column :label="t('components.dailyNews.titleLabel')">
            <template #default="scope:{ row: any, column: any, $index: number }">
              <el-text class="news-title" @click="handleNewsTableClick(scope.row)">
                {{ scope.row.title }}
              </el-text>
            </template>
          </el-table-column>
          <el-table-column prop="popularity" :label="t('components.dailyNews.popularityLabel')" width="120"/>
        </el-table>
      </el-tab-pane>
    </el-tabs>
    <el-text class="tip">{{ t('components.dailyNews.tipText') }}
      <el-link href="/news">{{ t('components.dailyNews.tipLink') }}</el-link>
    </el-text>
  </el-card>
</template>

<script setup lang="ts">
import {computed, ref} from "vue";
import {type HotItem, websiteNews, type WebsiteNewsRequest} from "@/api/website";
import type {TabsPaneContext} from "element-plus";
import {useI18n} from "vue-i18n";

const {t} = useI18n();

const newsTableData = ref<HotItem[]>()

interface newsTypeItem {
  name: string;
  label: string;
  src: string;
}

const newsTypeList = computed<newsTypeItem[]>(() => [
  {
    name: "baidu",
    label: t('components.dailyNews.baidu'),
    src: "/image/baidu.png"
  },
  {
    name: "zhihu",
    label: t('components.dailyNews.zhihu'),
    src: "/image/zhihu.png"
  },
  {
    name: "kuaishou",
    label: t('components.dailyNews.kuaishou'),
    src: "/image/kuaishou.png"
  },
  {
    name: "toutiao",
    label: t('components.dailyNews.toutiao'),
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
      newsMap.set(source, res.data.hot_list.slice(0, 7))
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
.daily-news {
  margin-bottom: var(--sp-4);
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  animation: kf-fade-up var(--dur-mid) var(--ease-out) backwards;

  .title {
    font-size: var(--fs-16);
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: var(--sp-3);
  }

  :deep(.el-tabs__item) {
    .tab-icon {
      width: 16px;
      height: 16px;
      margin-right: var(--sp-1);
      vertical-align: middle;
    }
  }

  .news-title {
    cursor: pointer;
    transition: color 150ms ease-out, transform 80ms ease-out;

    &:hover {
      color: var(--accent);
    }

    &:active {
      transform: translateY(1px);
    }
  }

  .tip {
    display: flex;
    margin-top: var(--sp-3);
    font-size: var(--fs-12);
    color: var(--text-muted);
  }
}

</style>
