<template>
  <div class="friend-link-update-form">
    <el-form
        :model="friendLinkUpdateFormData"
        :validate-on-rule-change="false"
    >

      <el-form-item :label="t('forms.friendLinkUpdate.logo')" prop="logo">
        <el-image :src="props.friendLink.logo" alt=""/>
      </el-form-item>
      <el-form-item :label="t('forms.friendLinkUpdate.link')" prop="link">
        <el-input
            v-model="friendLinkUpdateFormData.link"
            size="large"
            :placeholder="t('forms.friendLinkUpdate.linkPlaceholder')"
        />
      </el-form-item>
      <el-form-item :label="t('forms.friendLinkUpdate.name')" prop="name">
        <el-input
            v-model="friendLinkUpdateFormData.name"
            size="large"
            :placeholder="t('forms.friendLinkUpdate.namePlaceholder')"
        />
      </el-form-item>
      <el-form-item :label="t('forms.friendLinkUpdate.description')" prop="description">
        <el-input
            v-model="friendLinkUpdateFormData.description"
            size="large"
            :placeholder="t('forms.friendLinkUpdate.descriptionPlaceholder')"
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
              @click="layoutStore.state.friendLinkUpdateVisible = false"
          >{{ t('common.cancel') }}
          </el-button>
        </div>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import {defineProps, reactive} from 'vue';
import {ElMessage} from "element-plus";
import {type FriendLink, friendLinkUpdate, type FriendLinkUpdateRequest} from "@/api/friend-link";
import {useLayoutStore} from "@/stores/layout";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const layoutStore = useLayoutStore()

const props = defineProps<{
  friendLink: FriendLink;
}>();

const friendLinkUpdateFormData = reactive<FriendLinkUpdateRequest>({
  id: props.friendLink.id,
  link: props.friendLink.link,
  name: props.friendLink.name,
  description: props.friendLink.description,
})

const submitForm = async () => {
  const res = await friendLinkUpdate(friendLinkUpdateFormData)
  if (res.code === 0) {
    ElMessage.success(res.msg)
    layoutStore.state.shouldRefreshFriendLinkTable = true
    layoutStore.state.friendLinkUpdateVisible = false
  }
}
</script>

<style scoped lang="scss">
.friend-link-update-form {
  .el-form {
    .el-form-item {
      .el-image {
        height: 120px;
        width: 100%;
        object-fit: cover;
        border-radius: var(--radius-sm);
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
