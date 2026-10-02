<template>
  <div class="advertisement-create-form">
    <el-form
        :model="advertisementCreateFormData"
        :validate-on-rule-change="false"
    >
      <el-form-item :label="t('forms.advertisementCreate.image')" prop="ad_image">
        <el-upload
            :action="`${path}/image/upload`"
            drag
            with-credentials
            :headers="{'x-access-token':userStore.state.accessToken}"
            :show-file-list="false"
            :on-success="handleSuccess"
            :on-error="handleSuccess"
            name="image"
        >

          <el-image v-if="advertisementCreateFormData.ad_image" :src="advertisementCreateFormData.ad_image" alt=""/>

          <div v-else class="upload-content">
            <div class="container">
              <component is="UploadFilled" class="upload-filled"></component>
              <div class="el-upload__text">
                {{ t('forms.advertisementCreate.dragPrefix') }}<em>{{ t('forms.advertisementCreate.dragClick') }}</em>
              </div>
            </div>
          </div>

          <template #tip>
            <div class="el-upload__tip">
              {{ t('forms.advertisementCreate.uploadTip') }}
            </div>
          </template>
        </el-upload>

        <el-input
            v-model="advertisementCreateFormData.ad_image"
            size="large"
            disabled
        />
      </el-form-item>
      <el-form-item :label="t('forms.advertisementCreate.link')" prop="link">
        <el-input
            v-model="advertisementCreateFormData.link"
            size="large"
            :placeholder="t('forms.advertisementCreate.linkPlaceholder')"
        />
      </el-form-item>
      <el-form-item :label="t('forms.advertisementCreate.title')" prop="title">
        <el-input
            v-model="advertisementCreateFormData.title"
            size="large"
            :placeholder="t('forms.advertisementCreate.titlePlaceholder')"
        />
      </el-form-item>
      <el-form-item :label="t('forms.advertisementCreate.content')" prop="content">
        <el-input
            v-model="advertisementCreateFormData.content"
            size="large"
            :placeholder="t('forms.advertisementCreate.contentPlaceholder')"
        />
      </el-form-item>
      <el-form-item>
        <div class="button-group">
          <el-button
              type="primary"
              size="large"
              @click="submitForm"
          >{{ t('common.confirm') }}
          </el-button>
          <el-button
              size="large"
              @click="layoutStore.state.advertisementCreateVisible = false"
          >{{ t('common.cancel') }}
          </el-button>
        </div>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref} from "vue";
import {ElMessage} from "element-plus";
import {advertisementCreate} from "@/api/advertisement";
import type {AdvertisementCreateRequest} from "@/api/advertisement";
import type {ApiResponse} from "@/utils/request";
import type {ImageUploadResponse} from "@/api/image";
import {useUserStore} from "@/stores/user";
import {useLayoutStore} from "@/stores/layout";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const userStore = useUserStore()
const layoutStore = useLayoutStore()

const path = ref(import.meta.env.VITE_BASE_API)

const advertisementCreateFormData = reactive<AdvertisementCreateRequest>({
  ad_image: "",
  link: "",
  title: "",
  content: "",
})
const handleSuccess = (res: ApiResponse<ImageUploadResponse>) => {
  if (res.code === 0) {
    advertisementCreateFormData.ad_image = res.data.url
    ElMessage.success(res.msg)
  }
}


const submitForm = async () => {
  const res = await advertisementCreate(advertisementCreateFormData)
  if (res.code === 0) {
    ElMessage.success(res.msg)
    layoutStore.state.shouldRefreshAdvertisementTable = true
    layoutStore.state.advertisementCreateVisible = false
  }
}
</script>

<style scoped lang="scss">
.advertisement-create-form {
  .el-form {
    .el-form-item {
      .el-image {
        height: 120px;
        width: 100%;
        object-fit: cover;
        border-radius: var(--radius-sm);
      }

      .upload-content {
        display: flex;
        height: 140px;
        width: 100%;
        border: 1px dashed var(--accent);
        border-radius: var(--radius-sm);
        background-color: transparent;
        transition: background-color 0.15s ease-out;

        &:hover {
          background-color: var(--accent-weak);
        }

        .container {
          margin: auto;
          text-align: center;
          color: var(--text-muted);

          .upload-filled {
            height: 32px;
            width: 32px;
            color: var(--accent);
          }

          .el-upload__text {
            font-size: var(--fs-14);
            margin-top: var(--sp-2);

            em {
              color: var(--accent);
              font-style: normal;
            }
          }
        }
      }

      .el-upload__tip {
        color: var(--text-muted);
        font-size: var(--fs-12);
      }

      .button-group {
        margin-left: auto;
        display: flex;
        gap: var(--sp-2);

        .el-button {
          height: 40px;
          min-width: 88px;
          font-weight: 500;
        }
      }
    }
  }
}
</style>

<style lang="scss">
.el-upload {
  --el-upload-dragger-padding-horizontal: 0px;
  --el-upload-dragger-padding-vertical: 0px;
  line-height: 0;
}
</style>
