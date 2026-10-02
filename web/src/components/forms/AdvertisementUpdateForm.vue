<template>
  <div class="advertisement-update-form">
    <el-form
        :model="advertisementUpdateFormData"
        :validate-on-rule-change="false"
    >

      <el-form-item :label="t('forms.advertisementUpdate.image')" prop="ad_image">
        <el-image :src="props.advertisement.ad_image" alt=""/>
      </el-form-item>
      <el-form-item :label="t('forms.advertisementUpdate.link')" prop="link">
        <el-input
            v-model="advertisementUpdateFormData.link"
            size="large"
            :placeholder="t('forms.advertisementUpdate.linkPlaceholder')"
        />
      </el-form-item>
      <el-form-item :label="t('forms.advertisementUpdate.title')" prop="title">
        <el-input
            v-model="advertisementUpdateFormData.title"
            size="large"
            :placeholder="t('forms.advertisementUpdate.titlePlaceholder')"
        />
      </el-form-item>
      <el-form-item :label="t('forms.advertisementUpdate.content')" prop="content">
        <el-input
            v-model="advertisementUpdateFormData.content"
            size="large"
            :placeholder="t('forms.advertisementUpdate.contentPlaceholder')"
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
              @click="layoutStore.state.advertisementUpdateVisible = false"
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
import {type Advertisement, advertisementUpdate, type AdvertisementUpdateRequest} from "@/api/advertisement";
import {useLayoutStore} from "@/stores/layout";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const layoutStore = useLayoutStore()

const props = defineProps<{
  advertisement: Advertisement;
}>();

const advertisementUpdateFormData = reactive<AdvertisementUpdateRequest>({
  id: props.advertisement.id,
  link: props.advertisement.link,
  title: props.advertisement.title,
  content: props.advertisement.content,
})

const submitForm = async () => {
  const res = await advertisementUpdate(advertisementUpdateFormData)
  if (res.code === 0) {
    ElMessage.success(res.msg)
    layoutStore.state.shouldRefreshAdvertisementTable = true
    layoutStore.state.advertisementUpdateVisible = false
  }
}
</script>

<style scoped lang="scss">
.advertisement-update-form {
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
