<template>
  <div class="avatar-uploader">
    <div class="avatar-circle" @click="triggerSelect">
      <el-image class="avatar-img" :src="avatarUrl" fit="cover">
        <template #error>
          <div class="avatar-fallback">{{ fallbackText }}</div>
        </template>
      </el-image>
      <div class="avatar-mask">
        <el-icon v-if="!uploading" :size="18"><Camera/></el-icon>
        <el-icon v-else :size="18" class="is-loading"><Loading/></el-icon>
        <span>{{ uploading ? t('uploader.uploading') : t('uploader.changeAvatar') }}</span>
      </div>
    </div>
    <input
        ref="fileInput"
        type="file"
        accept="image/*"
        class="avatar-input"
        @change="onFileChange"
    />
    <div class="avatar-tip">{{ t('uploader.tip', {n: maxSize}) }}</div>
  </div>
</template>

<script setup lang="ts">
import {computed, ref} from "vue";
import {Camera, Loading} from "@element-plus/icons-vue";
import {ElMessage} from "element-plus";
import {uploadAvatar} from "@/api/user";
import {useUserStore} from "@/stores/user";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const emit = defineEmits<{ (e: 'changed'): void }>()

const userStore = useUserStore()
const fileInput = ref<HTMLInputElement>()
const uploading = ref(false)
const maxSize = 20

const avatarUrl = computed(() => userStore.state.userInfo.avatar || '/image/avatar.jpg')
const fallbackText = computed(() => (userStore.state.userInfo.username || 'U').slice(0, 1).toUpperCase())

const triggerSelect = () => {
  if (uploading.value) return
  fileInput.value?.click()
}

const onFileChange = (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return

  // 类型校验（与后端白名单保持一致）
  const allowed = [
    'image/jpeg',
    'image/png',
    'image/gif',
    'image/webp',
    'image/svg+xml',
    'image/tiff',
    'image/x-icon',
    'image/vnd.microsoft.icon',
  ]
  if (!allowed.includes(file.type)) {
    ElMessage.error(t('uploader.onlyImageError'))
    return
  }

  // 大小校验
  if (file.size / 1024 / 1024 >= maxSize) {
    ElMessage.error(t('uploader.sizeError', {n: maxSize}))
    return
  }

  uploadFile(file)
}

const uploadFile = async (file: File) => {
  uploading.value = true
  try {
    const res = await uploadAvatar(file)
    if (res.code === 0) {
      userStore.state.userInfo.avatar = res.data.url
      ElMessage.success(t('uploader.success'))
      emit('changed')
    }
  } finally {
    uploading.value = false
  }
}
</script>

<style scoped lang="scss">
.avatar-uploader {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: var(--sp-2);

  .avatar-circle {
    position: relative;
    width: 80px;
    height: 80px;
    border-radius: 50%;
    overflow: hidden;
    cursor: pointer;
    border: 1px solid var(--border);
    background-color: var(--bg-elevated);

    .avatar-img {
      width: 100%;
      height: 100%;
    }

    .avatar-fallback {
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: var(--fs-24);
      font-weight: 600;
      font-family: var(--font-serif);
      color: var(--text-primary);
    }

    .avatar-mask {
      position: absolute;
      inset: 0;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: var(--sp-1);
      background: color-mix(in srgb, var(--text-primary) 55%, transparent);
      color: var(--bg);
      font-size: var(--fs-12);
      opacity: 0;
      transition: opacity 0.2s ease-out;
    }

    &:hover .avatar-mask {
      opacity: 1;
    }
  }

  .avatar-input {
    display: none;
  }

  .avatar-tip {
    font-size: var(--fs-12);
    color: var(--text-muted);
  }
}
</style>
