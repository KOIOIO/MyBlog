<template>
  <div class="detail-page">
    <web-navbar :noScroll="true"/>
    <div class="page">
      <div v-if="post" class="article">
        <div class="meta-line">
          <span class="meta-category">{{ categoryLabel(post.category) }}</span>
          <span class="meta-sep">·</span>
          <span>{{ formatTime(post.created_at) }}</span>
          <span class="meta-sep">·</span>
          <span>👁 {{ post.view_count }}</span>
        </div>

        <h1 class="article-title">{{ post.title }}</h1>

        <div class="author-row">
          <el-avatar :size="32" :src="post.user.avatar" class="author-avatar">
            {{ (post.user.username || 'U').slice(0, 1).toUpperCase() }}
          </el-avatar>
          <span class="author-name">{{ post.user.username }}</span>
          <span class="author-time">{{ formatTime(post.created_at) }}</span>
        </div>

        <div class="article-body">{{ post.content }}</div>

        <div v-if="postImages.length" class="image-list">
          <img v-for="(img, idx) in postImages" :key="idx" :src="img" alt=""/>
        </div>

        <div v-if="postTags.length" class="tag-row">
          <span v-for="tag in postTags" :key="tag" class="tag-item">{{ tagLabel(tag) }}</span>
        </div>

        <div class="action-row">
          <button class="like-btn" :class="{ liked }" @click="onLike">
            <span class="like-icon">❤️</span>
            <span>{{ post.like_count }}</span>
          </button>
          <span class="comment-stat">💬 {{ post.comment_count }}</span>
        </div>

        <div class="comment-section">
          <div class="comment-input-box">
            <el-input
                v-model="commentText"
                type="textarea"
                :autosize="{ minRows: 3, maxRows: 6 }"
                :placeholder="userStore.isLoggedIn ? t('pages.forum.detail.commentPlaceholder') : t('pages.forum.detail.loginCommentPlaceholder')"
            />
            <div class="comment-input-actions">
              <el-button v-if="userStore.isLoggedIn" type="primary" size="small" :loading="commenting" @click="submitComment">{{ t('common.publish') }}</el-button>
              <el-button v-else size="small" type="primary" @click="openLogin">{{ t('auth.login') }}</el-button>
            </div>
          </div>

          <h3 class="comment-title">{{ t('pages.forum.detail.commentTitle', {count: comments.length}) }}</h3>

          <div v-if="comments.length" class="comment-list">
            <ForumCommentItem
                v-for="c in topLevelComments"
                :key="c.id"
                :comment="c"
                :post-id="post.id"
                @refresh="loadDetail"
            />
          </div>
          <div v-else class="comment-empty">{{ t('pages.forum.detail.commentEmpty') }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, ref} from "vue";
import {useRoute} from "vue-router";
import {ElMessage} from "element-plus";
import WebNavbar from "@/components/layout/WebNavbar.vue";
import ForumCommentItem from "@/components/forum/ForumCommentItem.vue";
import {forumComment, forumDetail, forumLike, type ForumComment, type ForumPost} from "@/api/forum";
import {useUserStore} from "@/stores/user";
import {useLayoutStore} from "@/stores/layout";
import {useI18n} from "vue-i18n";
import {categoryLabel, tagLabel} from "@/i18n/meta";

const {t} = useI18n();

const route = useRoute();
const userStore = useUserStore();
const layoutStore = useLayoutStore();

const post = ref<ForumPost | null>(null);
const comments = ref<ForumComment[]>([]);
const commentText = ref('');
const commenting = ref(false);
const liked = ref(false);

const postId = Number(route.params.id);

// 兼容后端返回 images/tags 为 JSON 字符串或数组两种格式
const postImages = computed<string[]>(() => {
    const val = post.value?.images;
    if (!val) return [];
    if (Array.isArray(val)) return val;
    try {
        const parsed = JSON.parse(val as unknown as string);
        return Array.isArray(parsed) ? parsed : [];
    } catch {
        return [];
    }
});

const postTags = computed<string[]>(() => {
    const val = post.value?.tags;
    if (!val) return [];
    if (Array.isArray(val)) return val;
    try {
        const parsed = JSON.parse(val as unknown as string);
        return Array.isArray(parsed) ? parsed : [];
    } catch {
        return [];
    }
});

const formatTime = (time: string) => {
    if (!time) return '';
    const d = new Date(time);
    if (isNaN(d.getTime())) return time;
    return d.toLocaleString();
};

const loadDetail = async () => {
    const res = await forumDetail({id: postId});
    if (res.code === 0) {
        post.value = res.data;
        comments.value = res.data.comments || [];
    }
};

onMounted(loadDetail);

const topLevelComments = computed(() => comments.value);

const onLike = async () => {
    if (!userStore.isLoggedIn) {
        ElMessage.warning(t('pages.forum.detail.likeLoginRequired'));
        layoutStore.state.popoverVisible = true;
        layoutStore.state.loginVisible = true;
        return;
    }
    const res = await forumLike({post_id: postId});
    if (res.code === 0 && post.value) {
        liked.value = res.data.liked;
        post.value.like_count = res.data.like_count;
    }
};

const submitComment = async () => {
    const content = commentText.value.trim();
    if (!content) {
        ElMessage.warning(t('pages.forum.detail.commentRequired'));
        return;
    }
    commenting.value = true;
    try {
        const res = await forumComment({
            post_id: postId,
            parent_id: 0,
            content,
        });
        if (res.code === 0) {
            ElMessage.success(t('pages.forum.detail.commentSuccess'));
            commentText.value = '';
            loadDetail();
        }
    } finally {
        commenting.value = false;
    }
};

const openLogin = () => {
    layoutStore.state.popoverVisible = true;
    layoutStore.state.loginVisible = true;
};
</script>

<style scoped lang="scss">
.detail-page {
  background-color: var(--bg);
  min-height: 100vh;

  .page {
    max-width: var(--reading-width);
    margin: 0 auto;
    padding: calc(70px + var(--sp-6)) var(--sp-4) var(--sp-9);
  }

  .meta-line {
    font-size: var(--fs-14);
    color: var(--text-muted);
    margin-bottom: var(--sp-3);
    display: flex;
    align-items: center;
    gap: var(--sp-2);

    .meta-category {
      color: var(--accent);
    }

    .meta-sep {
      color: var(--border);
    }
  }

  .article-title {
    font-family: var(--font-serif);
    font-size: var(--fs-36);
    font-weight: 700;
    line-height: var(--lh-title);
    color: var(--text-primary);
    margin: 0 0 var(--sp-4);
    letter-spacing: -0.01em;
  }

  .author-row {
    display: flex;
    align-items: center;
    gap: var(--sp-2);
    margin-bottom: var(--sp-6);

    .author-avatar {
      background-color: var(--accent-weak);
      color: var(--accent);
      font-size: var(--fs-14);
      font-weight: 600;
    }

    .author-name {
      font-size: var(--fs-14);
      font-weight: 600;
      color: var(--text-primary);
    }

    .author-time {
      font-size: var(--fs-12);
      color: var(--text-muted);
    }
  }

  .article-body {
    font-size: var(--fs-16);
    line-height: var(--lh-body);
    color: var(--text-body);
    white-space: pre-wrap;
    word-break: break-word;
    margin-bottom: var(--sp-5);
  }

  .image-list {
    display: flex;
    flex-direction: column;
    gap: var(--sp-3);
    margin-bottom: var(--sp-5);

    img {
      width: 100%;
      border-radius: var(--radius-md);
      display: block;
    }
  }

  .tag-row {
    display: flex;
    flex-wrap: wrap;
    gap: var(--sp-2);
    margin-bottom: var(--sp-5);

    .tag-item {
      font-size: var(--fs-12);
      color: var(--text-body);
      background: var(--accent-weak);
      padding: 3px var(--sp-2);
      border-radius: var(--radius-sm);
    }
  }

  .action-row {
    display: flex;
    align-items: center;
    gap: var(--sp-5);
    padding-top: var(--sp-4);
    border-top: 1px solid var(--border);
    margin-bottom: var(--sp-7);

    .like-btn {
      display: inline-flex;
      align-items: center;
      gap: var(--sp-1);
      border: none;
      background: transparent;
      padding: var(--sp-1) var(--sp-2);
      border-radius: var(--radius-sm);
      font-size: var(--fs-14);
      color: var(--text-muted);
      cursor: pointer;
      transition: color 150ms ease-out, background-color 150ms ease-out;

      .like-icon {
        filter: grayscale(1);
        transition: filter 150ms ease-out;
      }

      &:hover {
        color: var(--accent);
      }

      &.liked {
        color: var(--accent);

        .like-icon {
          filter: none;
          animation: kf-pop 150ms var(--ease-pop, cubic-bezier(.34, 1.56, .64, 1));
        }
      }
    }

    .comment-stat {
      font-size: var(--fs-14);
      color: var(--text-muted);
    }
  }

  .comment-section {
    padding-top: var(--sp-5);
    border-top: 1px solid var(--border);

    .comment-input-box {
      margin-bottom: var(--sp-6);

      .comment-input-actions {
        display: flex;
        justify-content: flex-end;
        margin-top: var(--sp-2);
      }
    }

    .comment-title {
      font-size: var(--fs-18);
      font-weight: 600;
      color: var(--text-primary);
      margin-bottom: var(--sp-3);
    }

    .comment-empty {
      padding: var(--sp-5) 0;
      text-align: center;
      font-size: var(--fs-14);
      color: var(--text-muted);
    }
  }
}
</style>
