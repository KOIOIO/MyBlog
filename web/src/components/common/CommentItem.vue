<template>
  <div class="comment-item">
    <div v-for="item in comments" :key="item.id" class="comment-row">
      <div class="item-card">
        <div class="title">
          <el-popover width="280">
            <template #reference>
              <el-avatar :src="item.user.avatar" class="avatar"/>
            </template>
            <template #default>
              <user-card :uuid="''"
                         :user-card-info="{uuid:item.user.uuid,username:item.user.username,avatar:item.user.avatar,address:item.user.address,signature:item.user.signature}"/>
            </template>
          </el-popover>
          <div class="name">
            {{ item.user.username }}
          </div>
        </div>
        <MdPreview class="content" :modelValue="item.content"/>
        <div class="meta">
          <span class="time">{{ getTime(item.created_at) }}</span>
          <div class="actions">
            <el-button v-if="replyFlag===item.id" type="primary" link size="small" @click="submitReply(item);content=''">{{ t('common.confirm') }}</el-button>
            <el-button v-if="replyFlag===item.id" link size="small" @click="content='';replyFlag=0">{{ t('common.cancel') }}</el-button>
            <el-button v-if="!(replyFlag===item.id)" type="primary" link size="small" @click="replyFlag=item.id">{{ t('comment.reply') }}</el-button>
            <el-button v-if="(item.user_uuid===userStore.state.userInfo.uuid||userStore.isAdmin)&&!(replyFlag===item.id)" type="danger" link size="small"
                       @click="handleDelete(item.id)">
              {{ t('common.delete') }}
            </el-button>
          </div>
        </div>
        <div class="reply-collapse" :class="{ open: replyFlag===item.id }">
          <div class="reply-collapse-inner">
            <el-input class="comment-input" v-model="content" :autosize="{ minRows: 2, maxRows: 6 }" type="textarea"
                      :placeholder="t('comment.replyPlaceholder')" maxlength="320"/>
            <div class="comment-tool">
              <el-popover width="448" trigger="click">
                <template #reference>
                  <el-avatar class="emoji-trigger">😊</el-avatar>
                </template>
                <template #default>
                  <div class="emoji-panel">
                    <span
                        v-for="emoji in emojis"
                        :key="emoji"
                        class="emoji-item"
                        @click="content=content+emoji"
                    >{{ emoji }}</span>
                  </div>
                </template>
              </el-popover>
            </div>
          </div>
        </div>
      </div>
      <div v-if="item.children && item.children.length">
        <div v-if="!isExpanded(item.id)" class="children-toggle" @click="toggleChildren(item.id)">
          {{ t('comment.repliesToggle', {n: item.children.length}) }}
        </div>
        <div v-else class="item-children">
          <comment-item :comments="item.children"/>
          <div class="children-toggle" @click="toggleChildren(item.id)">{{ t('comment.collapse') }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  type Comment,
  commentCreate,
  type CommentCreateRequest,
  commentDelete,
  type CommentDeleteRequest
} from "@/api/comment";
import {userCard} from "@/api/user";
import UserCard from "@/components/widgets/UserCard.vue";
import {useUserStore} from "@/stores/user";
import {MdPreview} from "md-editor-v3";
import {ref} from "vue";
import {useLayoutStore} from "@/stores/layout";
import {useI18n} from "vue-i18n";

defineProps<{
  comments: Comment[];
}>();

const userStore = useUserStore()
const {t} = useI18n()
const getTime = (date: Date): string => {
  const past = new Date(date).getTime()
  const diff = Math.floor((Date.now() - past) / 1000)
  if (diff < 60) return t('comment.justNow')
  const mins = Math.floor(diff / 60)
  if (mins < 60) return t('comment.minutesAgo', {n: mins})
  const hours = Math.floor(mins / 60)
  if (hours < 24) return t('comment.hoursAgo', {n: hours})
  const days = Math.floor(hours / 24)
  if (days < 7) return t('comment.daysAgo', {n: days})
  return new Date(date).toLocaleDateString()
}

const replyFlag = ref(0)
const content = ref('')
const expanded = ref<Set<number>>(new Set())
const isExpanded = (id: number): boolean => expanded.value.has(id)
const toggleChildren = (id: number): void => {
  const s = new Set(expanded.value)
  if (s.has(id)) s.delete(id)
  else s.add(id)
  expanded.value = s
}
const emojis = ['😀','😁','😂','🤣','😊','😇','🙂','😉','😍','🥰','😘','😋','😛','🤪','🤗','🤭','🤔','😏','🙄','😬','😮','😲','🥺','😢','😭','😤','😡','🤯','😱','😰','🥳','😎','🤓','👍','👎','👏','🙏','💪','🔥','❤️'];

const layoutStore = useLayoutStore()

const submitReply = async (item: Comment) => {
  const commentCreateRequest: CommentCreateRequest = {
    article_id: item.article_id,
    p_id: item.id,
    content: content.value,
  }
  const res = await commentCreate(commentCreateRequest)
  if (res.code === 0) {
    ElMessage.success(res.msg)
    layoutStore.state.shouldRefreshCommentList = true
    replyFlag.value = 0
  }
}

const handleDelete = async (id: number) => {
  let ids: number[] = []
  ids.push(id)
  const commentDeleteRequest: CommentDeleteRequest = {
    ids: ids
  }
  const res = await commentDelete(commentDeleteRequest)
  if (res.code === 0) {
    ElMessage.success(res.msg)
    layoutStore.state.shouldRefreshCommentList = true
  }
}
</script>

<style scoped lang="scss">
.comment-item {
  .comment-row {
    animation: kf-fade-up 200ms var(--ease-out, cubic-bezier(.16, 1, .3, 1)) both;
  }

  .item-card {
    padding: var(--sp-2) 0;
    border-bottom: 1px solid var(--border);

    .title {
      display: flex;
      align-items: center;
      gap: var(--sp-2);

      .avatar {
        width: 32px;
        height: 32px;
        cursor: pointer;
        flex-shrink: 0;
        border-radius: 6px;
      }

      .name {
        font-size: var(--fs-14);
        font-weight: 600;
        line-height: 20px;
        color: var(--text-primary);
      }
    }

    .content {
      margin-top: 4px;
      font-size: var(--fs-14);
      color: var(--text-body);
      line-height: 1.5;

      /* 评论正文：MdPreview 全链路压平（padding/字号/行高），紧凑如 B 站 */
      :deep(.md-editor),
      :deep(.md-editor-content),
      :deep(.md-editor-preview-wrapper) {
        padding: 0 !important;
      }

      :deep(.md-editor-preview) {
        padding: 0 !important;
        font-size: var(--fs-14) !important;
        line-height: 1.5 !important;
      }

      :deep(.md-editor-preview p) {
        margin: 0 !important;
        font-size: var(--fs-14) !important;
        line-height: 1.5 !important;
      }
    }

    /* B站式操作行：时间 / 回复 / 删除 同行，紧贴内容 */
    .meta {
      display: flex;
      align-items: center;
      gap: var(--sp-3);
      margin-top: 4px;
      min-height: 16px;

      .time {
        font-size: var(--fs-12);
        line-height: 16px;
        color: var(--text-muted);
      }

      .actions {
        display: flex;
        align-items: center;
        gap: var(--sp-2);

        :deep(.el-button) {
          font-size: var(--fs-12);
          line-height: 16px;
          padding: 0;
          height: auto;
          min-height: 0;
        }
      }
    }

    .reply-collapse {
      display: grid;
      grid-template-rows: 0fr;
      transition: grid-template-rows 220ms var(--ease-in-out, cubic-bezier(.4, 0, .2, 1));

      &.open {
        grid-template-rows: 1fr;
      }

      > .reply-collapse-inner {
        overflow: hidden;
        min-height: 0;
      }

      .comment-input {
        margin-top: var(--sp-2);
      }

      .comment-tool {
        margin-right: auto;
        margin-top: var(--sp-2);

        .emoji-trigger {
          background-color: unset;
          font-size: var(--fs-20);
          cursor: pointer;
        }

        .emoji-panel {
          display: flex;
          flex-wrap: wrap;
          gap: var(--sp-1);
          padding: var(--sp-2);

          .emoji-item {
            width: 34px;
            height: 34px;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            font-size: var(--fs-20);
            border-radius: 6px;
            cursor: pointer;
            transition: background-color 120ms ease;

            &:hover {
              background-color: var(--bg-elevated);
            }
          }
        }
      }
    }
  }


  .item-children {
    margin-left: var(--sp-4);
    margin-top: var(--sp-1);
    padding: var(--sp-1) 0 var(--sp-1) var(--sp-3);
    border-left: 2px solid var(--border);
    border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
    background: color-mix(in srgb, var(--bg-elevated) 45%, transparent);

    :deep(.el-avatar.avatar) {
      width: 24px;
      height: 24px;
      border-radius: 6px;
    }

    .item-card {
      padding: var(--sp-1) 0;
      border-bottom: none;
    }
  }

  /* B站式折叠入口：共 N 条回复，点击查看 */
  .children-toggle {
    margin-top: var(--sp-2);
    font-size: var(--fs-12);
    color: var(--text-muted);
    cursor: pointer;
    user-select: none;
    transition: color 120ms ease;

    &:hover {
      color: var(--accent);
    }
  }
}

.el-popover.el-popper {
  .el-image {
    height: 50px;
    width: 50px;
  }
}
</style>
