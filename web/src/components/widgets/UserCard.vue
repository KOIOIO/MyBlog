<template>
  <div class="user-card">
    <div class="user-info">
      <el-avatar :size="50" :src="userCardInfo.avatar"/>
      <el-row>{{ userCardInfo.username }}</el-row>
      <el-row>{{ userCardInfo.address }}</el-row>
    </div>
    <el-row class="uuid">{{ t('auth.uuid') }}：{{ userCardInfo.uuid }}</el-row>
    <div class="container">
      <el-row>{{ t('auth.signature') }}：{{ userCardInfo.signature }}</el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import {defineProps, ref} from "vue";
import {userCard, type UserCardRequest, type UserCardResponse} from "@/api/user";
import {useI18n} from "vue-i18n";

const {t} = useI18n();

const props = defineProps<{
  uuid: string;
  userCardInfo: UserCardResponse | null;
}>();

const userCardInfo = ref<UserCardResponse>({
  uuid: "",
  username: "",
  avatar: "",
  address: "",
  signature: "",
})

const getUserCard = async () => {
  if (props.userCardInfo) {
    userCardInfo.value = props.userCardInfo
  } else {
    const userCardRequest: UserCardRequest = {
      uuid: props.uuid,
    }

    const res = await userCard(userCardRequest)
    if (res.code === 0) {
      userCardInfo.value = res.data
      sendDataToParent(res.data)
    }
  }
}

getUserCard()

// 定义 emit 事件及其类型
const emit = defineEmits<{
  (event: 'userCardInfo', message: UserCardResponse): void;
}>();

// 定义触发事件的函数
const sendDataToParent = (userCardInfo:UserCardResponse) => {
  emit('userCardInfo', userCardInfo);
};
</script>

<style scoped lang="scss">
.user-card {
  max-width: 320px;
  padding: var(--sp-5);
  background-color: var(--bg-elevated);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  transition: border-color 150ms ease-out, box-shadow 150ms ease-out;

  &:hover {
    border-color: var(--accent);
    box-shadow: var(--shadow-md);
  }

  .user-info {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;

    :deep(.el-avatar) {
      flex-shrink: 0;
      transition: opacity 150ms ease-out;
    }

    &:hover :deep(.el-avatar) {
      opacity: 0.85;
    }

    .el-row:first-of-type {
      margin-top: var(--sp-3);
      font-size: var(--fs-16);
      font-weight: 600;
      color: var(--text-primary);
      line-height: var(--lh-title);
    }

    .el-row:last-of-type {
      margin-top: var(--sp-1);
      font-size: var(--fs-12);
      color: var(--text-muted);
    }
  }

  .uuid {
    margin-top: var(--sp-4);
    padding-top: var(--sp-3);
    border-top: 1px solid var(--border);
    font-size: var(--fs-12);
    color: var(--text-muted);
    font-family: var(--font-mono);
    word-break: break-all;
  }

  .container {
    margin-top: var(--sp-3);
    background-color: transparent;

    .el-row {
      padding: 0;
      font-size: var(--fs-14);
      color: var(--text-body);
      line-height: var(--lh-body);
    }
  }
}
</style>
