<template>
  <div class="image-uploader">
    <div class="thumb-grid">
      <div v-for="(url, idx) in modelValue" :key="url + idx" class="thumb-item">
        <img :src="url" alt="" />
        <button class="thumb-del" @click="removeImage(idx)" :aria-label="t('common.delete')">×</button>
      </div>
      <div class="upload-trigger" @click="triggerSelect">
        <el-icon v-if="!uploading" :size="20"><Plus /></el-icon>
        <el-icon v-else :size="20" class="is-loading"><Loading /></el-icon>
        <span>{{ uploading ? t('uploader.uploading') : t('components.forumImageUploader.uploadImage') }}</span>
      </div>
    </div>
    <input
        ref="fileInput"
        type="file"
        accept="image/*"
        class="file-input"
        @change="onFileChange"
    />
  </div>
</template>

<script setup lang="ts">
import {ref} from "vue";
import {Plus, Loading} from "@element-plus/icons-vue";
import {ElMessage} from "element-plus";
import {forumUpload} from "@/api/forum";
import {useI18n} from "vue-i18n";

const {t} = useI18n();

const modelValue = defineModel<string[]>({required: true});

const fileInput = ref<HTMLInputElement>();
const uploading = ref(false);
const maxSize = 20;

const triggerSelect = () => {
    if (uploading.value) return;
    fileInput.value?.click();
};

const onFileChange = (event: Event) => {
    const input = event.target as HTMLInputElement;
    const files = Array.from(input.files || []);
    input.value = '';
    if (!files.length) return;

    const allowed = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
    for (const file of files) {
        if (!allowed.includes(file.type)) {
            ElMessage.error(t('uploader.onlyImageError'));
            continue;
        }
        if (file.size / 1024 / 1024 >= maxSize) {
            ElMessage.error(t('uploader.sizeError', {n: maxSize}));
            continue;
        }
        uploadFile(file);
    }
};

const uploadFile = async (file: File) => {
    uploading.value = true;
    try {
        const res = await forumUpload(file);
        if (res.code === 0) {
            modelValue.value = [...modelValue.value, res.data.url];
        }
    } finally {
        uploading.value = false;
    }
};

const removeImage = (idx: number) => {
    const next = [...modelValue.value];
    next.splice(idx, 1);
    modelValue.value = next;
};
</script>

<style scoped lang="scss">
.image-uploader {
  .thumb-grid {
    display: flex;
    flex-wrap: wrap;
    gap: var(--sp-2);
  }

  .thumb-item {
    position: relative;
    width: 80px;
    height: 80px;
    border-radius: var(--radius-sm);
    overflow: hidden;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .thumb-del {
      position: absolute;
      top: 2px;
      right: 2px;
      width: 18px;
      height: 18px;
      border: none;
      border-radius: 50%;
      background: color-mix(in srgb, var(--text-primary) 60%, transparent);
      color: var(--bg);
      font-size: 12px;
      line-height: 1;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
      transition: opacity 150ms ease-out, transform 80ms ease-out;

      &:hover {
        opacity: 0.85;
      }

      &:active {
        transform: scale(0.9);
      }
    }
  }

  .upload-trigger {
    width: 80px;
    height: 80px;
    border: 1px dashed var(--accent);
    border-radius: var(--radius-sm);
    background: transparent;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: var(--sp-1);
    color: var(--text-muted);
    font-size: var(--fs-12);
    cursor: pointer;
    transition: background-color 150ms ease-out, color 150ms ease-out, border-color 150ms ease-out, transform 80ms ease-out;

    &:hover {
      background: var(--accent-weak);
      color: var(--accent);
      border-color: var(--accent);
    }

    &:active {
      transform: translateY(1px);
    }
  }

  .file-input {
    display: none;
  }
}
</style>
