<template>
  <div class="friend-link-create-form">
    <el-form
        :model="friendLinkCreateFormData"
        :validate-on-rule-change="false"
    >
      <el-form-item :label="t('forms.friendLinkCreate.logo')" prop="logo">
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

          <el-image v-if="friendLinkCreateFormData.logo" :src="friendLinkCreateFormData.logo" alt=""/>

          <div v-else class="upload-content">
            <div class="container">
              <component is="UploadFilled" class="upload-filled"></component>
              <div class="el-upload__text">
                {{ t('forms.friendLinkCreate.dragPrefix') }}<em>{{ t('forms.friendLinkCreate.dragClick') }}</em>
              </div>
            </div>
          </div>

          <template #tip>
            <div class="el-upload__tip">
              {{ t('forms.friendLinkCreate.uploadTip') }}
            </div>
          </template>
        </el-upload>

        <el-input
            v-model="friendLinkCreateFormData.logo"
            size="large"
            disabled
        />
      </el-form-item>
      <el-form-item :label="t('forms.friendLinkCreate.link')" prop="link">
        <el-input
            v-model="friendLinkCreateFormData.link"
            size="large"
            :placeholder="t('forms.friendLinkCreate.linkPlaceholder')"
        />
      </el-form-item>
      <el-form-item :label="t('forms.friendLinkCreate.name')" prop="name">
        <el-input
            v-model="friendLinkCreateFormData.name"
            size="large"
            :placeholder="t('forms.friendLinkCreate.namePlaceholder')"
        />
      </el-form-item>
      <el-form-item :label="t('forms.friendLinkCreate.description')" prop="description">
        <el-input
            v-model="friendLinkCreateFormData.description"
            size="large"
            :placeholder="t('forms.friendLinkCreate.descriptionPlaceholder')"
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
              @click="layoutStore.state.friendLinkCreateVisible = false"
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
import {friendLinkCreate} from "@/api/friend-link";
import type {FriendLinkCreateRequest} from "@/api/friend-link";
import type {ApiResponse} from "@/utils/request";
import type {ImageUploadResponse} from "@/api/image";
import {useUserStore} from "@/stores/user";
import {useLayoutStore} from "@/stores/layout";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const userStore = useUserStore()
const layoutStore = useLayoutStore()

const path = ref(import.meta.env.VITE_BASE_API)

const friendLinkCreateFormData = reactive<FriendLinkCreateRequest>({
  logo: "",
  link: "",
  name: "",
  description: "",
})
const handleSuccess = (res: ApiResponse<ImageUploadResponse>) => {
  if (res.code === 0) {
    friendLinkCreateFormData.logo = res.data.url
    ElMessage.success(res.msg)
  }
}


const submitForm = async () => {
  const res = await friendLinkCreate(friendLinkCreateFormData)
  if (res.code === 0) {
    ElMessage.success(res.msg)
    layoutStore.state.shouldRefreshFriendLinkTable = true
    layoutStore.state.friendLinkCreateVisible = false
  }
}
</script>

<style scoped lang="scss">
.friend-link-create-form {
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
