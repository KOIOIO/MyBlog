<template>
  <div class="post-card" @click="goDetail">
    <div class="post-main">
      <div class="post-head">
        <span class="post-category">{{ categoryLabel(post.category) }}</span>
        <h3 class="post-title">{{ post.title }}</h3>
      </div>
      <p class="post-summary">{{ summary }}</p>
      <div class="post-meta">
        <el-avatar :size="20" :src="post.user.avatar" class="meta-avatar">
          {{ (post.user.username || 'U').slice(0, 1).toUpperCase() }}
        </el-avatar>
        <span class="meta-author">{{ post.user.username }}</span>
        <span class="meta-dot">·</span>
        <span class="meta-time">{{ formatTime(post.created_at) }}</span>
        <span class="meta-dot">·</span>
        <span class="meta-stat">👁 {{ post.view_count }}</span>
        <span class="meta-stat">💬 {{ post.comment_count }}</span>
        <span class="meta-stat">❤️ {{ post.like_count }}</span>
      </div>
    </div>
    <div v-if="images.length" class="post-thumb">
      <img :src="images[0]" alt="" />
    </div>
  </div>
</template>

<script setup lang="ts">
import {computed} from "vue";
import {useRouter} from "vue-router";
import type {ForumPost} from "@/api/forum";
import {useI18n} from "vue-i18n";
import {categoryLabel} from "@/i18n/meta";

const {t} = useI18n();

const props = defineProps<{ post: ForumPost }>();

const router = useRouter();

// 兼容后端返回 images 为 JSON 字符串或数组两种格式
const images = computed<string[]>(() => {
    const val = props.post.images;
    if (!val) return [];
    if (Array.isArray(val)) return val;
    try {
        const parsed = JSON.parse(val as unknown as string);
        return Array.isArray(parsed) ? parsed : [];
    } catch {
        return [];
    }
});

const summary = computed(() => {
    const text = (props.post.content || '').replace(/\s+/g, ' ').trim();
    return text.length > 100 ? text.slice(0, 100) + '…' : text;
});

const formatTime = (time: string) => {
    if (!time) return '';
    const d = new Date(time);
    if (isNaN(d.getTime())) return time;
    return d.toLocaleDateString();
};

const goDetail = () => {
    router.push({name: 'forum-detail', params: {id: props.post.id}});
};
</script>

<style scoped lang="scss">
.post-card {
  display: flex;
  align-items: flex-start;
  gap: var(--sp-4);
  padding: var(--sp-5) var(--sp-3);
  border-bottom: 1px solid var(--border);
  cursor: pointer;
  transition: background-color 150ms ease-out, transform 80ms ease-out;

  &:hover {
    background-color: var(--bg-elevated);
  }

  &:active {
    transform: translateY(1px);
  }

  .post-main {
    flex: 1;
    min-width: 0;
  }

  .post-head {
    display: flex;
    align-items: center;
    gap: var(--sp-2);
    margin-bottom: var(--sp-2);

    .post-category {
      flex-shrink: 0;
      font-size: var(--fs-12);
      color: var(--accent);
      background: var(--accent-weak);
      padding: 1px var(--sp-2);
      border-radius: var(--radius-sm);
    }

    .post-title {
      font-size: var(--fs-18);
      font-weight: 600;
      color: var(--text-primary);
      line-height: var(--lh-title);
      margin: 0;
      transition: color 150ms ease-out;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }

  &:hover .post-title {
    color: var(--accent);
  }

  .post-summary {
    font-size: var(--fs-14);
    color: var(--text-body);
    line-height: var(--lh-body);
    margin: 0 0 var(--sp-3);
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .post-meta {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: var(--sp-1);
    font-size: var(--fs-12);
    color: var(--text-muted);

    .meta-avatar {
      background-color: var(--accent-weak);
      color: var(--accent);
      font-size: var(--fs-12);
      font-weight: 600;
    }

    .meta-author {
      color: var(--text-muted);
    }

    .meta-dot {
      color: var(--border);
    }

    .meta-stat {
      display: inline-flex;
      align-items: center;
    }
  }

  .post-thumb {
    flex-shrink: 0;
    width: 80px;
    height: 80px;
    border-radius: var(--radius-sm);
    overflow: hidden;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }
}
</style>
