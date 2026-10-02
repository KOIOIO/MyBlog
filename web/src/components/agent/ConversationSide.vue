<template>
  <div class="conv-side">
    <div v-if="conversations.length === 0" class="no-history">{{ t('pages.agent.noHistory') }}</div>
    <div
        v-for="conv in conversations"
        :key="conv.id"
        class="conv-item"
        :class="{active: conv.id === currentId}"
        @click="emit('select', conv.id)"
    >
      <div class="conv-title">{{ conv.title || t('pages.agent.untitled') }}</div>
      <div class="conv-time">{{ formatTime(conv.updated_at) }}</div>
      <el-popconfirm
          :title="t('pages.agent.deleteConfirm')"
          width="220"
          @confirm="emit('remove', conv.id)"
      >
        <template #reference>
          <el-icon class="del-btn" @click.stop><Delete/></el-icon>
        </template>
      </el-popconfirm>
    </div>
  </div>
</template>

<script setup lang="ts">
import {Delete} from "@element-plus/icons-vue";
import dayjs from "dayjs";
import {useI18n} from "vue-i18n";
import type {AgentConversation} from "@/api/agent";

const props = defineProps<{
    conversations: AgentConversation[];
    currentId: number;
}>();

const emit = defineEmits<{
    (e: 'select', id: number): void;
    (e: 'remove', id: number): void;
}>();

const {t} = useI18n();

const formatTime = (s: string): string => dayjs(s).format('MM-DD HH:mm');
</script>

<style scoped lang="scss">
.conv-side {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 6px;

  .no-history {
    padding: 24px 8px;
    text-align: center;
    font-size: 13px;
    color: var(--text-secondary);
  }

  .conv-item {
    position: relative;
    padding: 10px 12px;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.2s;

    &:hover {
      background: var(--el-fill-color-light);

      .del-btn { opacity: 1; }
    }

    &.active {
      background: var(--el-color-primary-light-9);
      border: 1px solid var(--el-color-primary-light-7);
    }

    .conv-title {
      font-size: 14px;
      color: var(--text-body);
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      padding-right: 28px;
    }

    .conv-time {
      margin-top: 4px;
      font-size: 12px;
      color: var(--text-secondary);
    }

    .del-btn {
      position: absolute;
      top: 12px;
      right: 10px;
      opacity: 0;
      color: var(--el-color-danger);
      transition: opacity 0.2s;
    }
  }
}
</style>
