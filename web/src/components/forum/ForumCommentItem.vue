<template>
  <div class="comment-node" :class="{ child: depth > 0 }">
    <div class="node-row">
      <el-avatar :size="depth > 0 ? 24 : 32" :src="comment.user.avatar" class="node-avatar">
        {{ (comment.user.username || 'U').slice(0, 1).toUpperCase() }}
      </el-avatar>
      <div class="node-body">
        <div class="node-head">
          <span class="node-name">{{ comment.user.username }}</span>
          <span class="node-time">{{ formatTime(comment.created_at) }}</span>
        </div>
        <div class="node-content">{{ comment.content }}</div>
        <div class="node-actions">
          <button class="reply-btn" @click="toggleReply">{{ t('comment.reply') }}</button>
        </div>

        <div class="reply-collapse" :class="{ open: replying }">
          <div class="reply-collapse-inner">
            <el-input
                v-model="replyText"
                type="textarea"
                :autosize="{ minRows: 2, maxRows: 5 }"
                :placeholder="t('components.forumComment.replyPlaceholder')"
            />
            <div class="reply-actions">
              <el-button size="small" @click="cancelReply">{{ t('common.cancel') }}</el-button>
              <el-button size="small" type="primary" :loading="submitting" @click="submitReply">{{ t('common.publish') }}</el-button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="comment.children && comment.children.length" class="node-children">
      <ForumCommentItem
          v-for="child in comment.children"
          :key="child.id"
          :comment="child"
          :post-id="postId"
          :depth="depth + 1"
          @refresh="emit('refresh')"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref} from "vue";
import {ElMessage} from "element-plus";
import {forumComment, type ForumComment} from "@/api/forum";
import {useUserStore} from "@/stores/user";
import {useLayoutStore} from "@/stores/layout";
import {useI18n} from "vue-i18n";

const {t} = useI18n();

const props = withDefaults(defineProps<{
    comment: ForumComment;
    postId: number;
    depth?: number;
}>(), {depth: 0});

const emit = defineEmits<{ (e: 'refresh'): void }>();

const userStore = useUserStore();
const layoutStore = useLayoutStore();

const replying = ref(false);
const replyText = ref('');
const submitting = ref(false);

const formatTime = (time: string) => {
    if (!time) return '';
    const d = new Date(time);
    if (isNaN(d.getTime())) return time;
    return d.toLocaleString();
};

const toggleReply = () => {
    if (!userStore.isLoggedIn) {
        ElMessage.warning(t('components.forumComment.loginRequired'));
        layoutStore.state.popoverVisible = true;
        layoutStore.state.loginVisible = true;
        return;
    }
    replying.value = !replying.value;
};

const cancelReply = () => {
    replying.value = false;
    replyText.value = '';
};

const submitReply = async () => {
    const content = replyText.value.trim();
    if (!content) {
        ElMessage.warning(t('components.forumComment.replyRequired'));
        return;
    }
    submitting.value = true;
    try {
        const res = await forumComment({
            post_id: props.postId,
            parent_id: props.comment.id,
            content,
        });
        if (res.code === 0) {
            ElMessage.success(t('components.forumComment.replySuccess'));
            cancelReply();
            emit('refresh');
        }
    } finally {
        submitting.value = false;
    }
};
</script>

<style scoped lang="scss">
.comment-node {
  padding: var(--sp-4) 0;
  animation: kf-fade-up 200ms var(--ease-out, cubic-bezier(.16, 1, .3, 1)) both;

  &.child {
    margin-left: var(--sp-4);
    padding-left: var(--sp-4);
    border-left: 2px solid var(--border);
  }

  .node-row {
    display: flex;
    gap: var(--sp-3);

    .node-avatar {
      flex-shrink: 0;
      background-color: var(--accent-weak);
      color: var(--accent);
      font-size: var(--fs-12);
      font-weight: 600;
    }

    .node-body {
      flex: 1;
      min-width: 0;

      .node-head {
        display: flex;
        align-items: baseline;
        gap: var(--sp-2);

        .node-name {
          font-size: var(--fs-14);
          font-weight: 600;
          color: var(--text-primary);
        }

        .node-time {
          font-size: var(--fs-12);
          color: var(--text-muted);
        }
      }

      .node-content {
        margin-top: var(--sp-1);
        font-size: var(--fs-14);
        color: var(--text-body);
        line-height: var(--lh-body);
        white-space: pre-wrap;
        word-break: break-word;
      }

      .node-actions {
        margin-top: var(--sp-1);

        .reply-btn {
          border: none;
          background: transparent;
          padding: 0;
          font-size: var(--fs-12);
          color: var(--accent);
          cursor: pointer;
          transition: opacity 150ms ease-out;

          &:hover {
            opacity: 0.8;
          }
        }
      }

      .reply-collapse {
        display: grid;
        grid-template-rows: 0fr;
        margin-top: var(--sp-2);
        transition: grid-template-rows 220ms var(--ease-in-out, cubic-bezier(.4, 0, .2, 1));

        &.open {
          grid-template-rows: 1fr;
        }

        > .reply-collapse-inner {
          overflow: hidden;
          min-height: 0;
        }

        .reply-actions {
          display: flex;
          justify-content: flex-end;
          gap: var(--sp-2);
          margin-top: var(--sp-2);
        }
      }
    }
  }

  .node-children {
    margin-top: var(--sp-1);
  }
}
</style>
