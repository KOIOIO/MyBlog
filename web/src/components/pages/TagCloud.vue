<template>
  <el-card class="tag-cloud">
    <el-row class="title">{{ t('components.tagCloud.title') }}</el-row>
    <div class="tag-list">
      <span v-for="(item, i) in tagCloudArray" :key="item.tag" class="tag-item"
            :style="{ animationDelay: Math.min(i, 5) * 60 + 'ms' }"
            @click="handleSearchJumps(item.tag)">
        {{ tagLabel(item.tag) }}
      </span>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import {ref} from "vue";
import {type ArticleTag, articleTags} from "@/api/article";
import {useI18n} from "vue-i18n";
import {tagLabel} from "@/i18n/meta";

const {t} = useI18n();

const tagTypes = ["primary", "success", "info", "warning", "danger"]

interface TagCloudItem {
  tag: string;
  number: number;
  type: string;
}

const tagCloudArray = ref<TagCloudItem[]>([])

const getTagCloudArray = async () => {
  let tagsArray: ArticleTag[]
  const res = await articleTags()
  if (res.code === 0) {
    tagsArray = res.data
    for (let i = 0; i < tagsArray.length; i++) {
      const item = tagsArray[i];
      const tagCloud: TagCloudItem = {
        tag: item.tag,
        number: item.number,
        type: tagTypes[i % tagTypes.length]
      }
      tagCloudArray.value.push(tagCloud);
    }
  }
}

getTagCloudArray()

const handleSearchJumps = (tag: string) => {
  window.open("/search?tag=" + tag)
}
</script>

<style scoped lang="scss">
.tag-cloud {
  margin-bottom: var(--sp-4);
  background: var(--bg-elevated);
  border: 1px solid var(--border);

  .title {
    font-size: var(--fs-16);
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: var(--sp-3);
  }

  .tag-list {
    display: flex;
    flex-wrap: wrap;
    gap: var(--sp-2);
  }

  .tag-item {
    display: inline-block;
    font-size: var(--fs-14);
    color: var(--text-body);
    background: var(--accent-weak);
    padding: 2px var(--sp-2);
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: color 150ms cubic-bezier(.16,1,.3,1), transform 150ms cubic-bezier(.16,1,.3,1);
    animation: kf-fade-up var(--dur-mid) var(--ease-out) backwards;

    &:hover {
      color: var(--accent);
      transform: scale(1.05);
    }

    &:active {
      transform: scale(1.02) translateY(1px);
      transition-duration: 80ms;
    }
  }
}

</style>
