<template>
  <div class="article-create-form">
    <el-form
        :model="articleCreateFormData"
        :validate-on-rule-change="false"
    >
      <el-form-item :label="t('forms.articleCreate.cover')" prop="cover">
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

          <el-image v-if="articleCreateFormData.cover" :src="articleCreateFormData.cover" alt=""/>

          <div v-else class="upload-content">
            <div class="container">
              <component is="UploadFilled" class="upload-filled"></component>
              <div class="el-upload__text">
                {{ t('forms.articleCreate.dragPrefix') }}<em>{{ t('forms.articleCreate.dragClick') }}</em>
              </div>
            </div>
          </div>

          <template #tip>
            <div class="el-upload__tip">
              {{ t('forms.articleCreate.uploadTip') }}
            </div>
          </template>
        </el-upload>

        <el-input
            v-model="articleCreateFormData.cover"
            size="large"
            disabled
        />
      </el-form-item>
      <el-form-item :label="t('forms.articleCreate.title')" prop="title" class="form-item--title">
        <el-input
            v-model="articleCreateFormData.title"
            size="large"
            :placeholder="t('forms.articleCreate.titlePlaceholder')"
        />
      </el-form-item>
      <el-form-item :label="t('forms.articleCreate.category')" prop="category">
        <el-select
            v-model="articleCreateFormData.category"
            size="large"
            :placeholder="t('forms.articleCreate.categoryPlaceholder')"
            style="width: 100%"
            @change="handleCategoryChange"
        >
          <el-option :label="categoryLabel('技术')" value="技术"/>
          <el-option :label="categoryLabel('生活')" value="生活"/>
        </el-select>
      </el-form-item>
      <el-form-item :label="t('forms.articleCreate.tags')" prop="tags">
        <el-select
            v-model="articleCreateFormData.tags"
            multiple
            filterable
            collapse-tags
            collapse-tags-tooltip
            size="large"
            style="width: 100%"
            :placeholder="articleCreateFormData.category ? t('forms.articleCreate.tagPlaceholder') : t('forms.articleCreate.tagPlaceholderNoCategory')"
        >
          <el-option
              v-for="item in filteredTagOptions"
              :key="item.tag"
              :label="tagLabel(item.tag)"
              :value="item.tag"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('forms.articleCreate.abstract')" prop="abstract">
        <el-input
            v-model="articleCreateFormData.abstract"
            type="textarea"
            :placeholder="t('forms.articleCreate.abstractPlaceholder')"
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
              @click="layoutStore.state.articleCreateVisible = false"
          >{{ t('common.cancel') }}
          </el-button>
        </div>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import {computed, defineProps, onMounted, reactive, ref} from "vue";
import {ElMessage} from "element-plus";
import {articleCreate, articleTags, type ArticleCreateRequest} from "@/api/article";
import type {ApiResponse} from "@/utils/request";
import type {ImageUploadResponse} from "@/api/image";
import {useUserStore} from "@/stores/user";
import {useLayoutStore} from "@/stores/layout";
import {useI18n} from "vue-i18n";
import {categoryLabel, tagLabel} from "@/i18n/meta";

const {t} = useI18n()

interface TagOption {
  tag: string;
  group: string;
  number: number;
}

const props = defineProps<{
  title: string;
  content: string;
}>();

const userStore = useUserStore()
const layoutStore = useLayoutStore()

const path = ref(import.meta.env.VITE_BASE_API)

const articleCreateFormData = reactive<ArticleCreateRequest>({
  cover: '',
  title: props.title,
  category: '',
  tags: [],
  abstract: '',
  content: props.content,
})

const tagOptions = ref<TagOption[]>([])

const filteredTagOptions = computed<TagOption[]>(() => {
  if (!articleCreateFormData.category) {
    return tagOptions.value
  }
  const group = articleCreateFormData.category === '技术' ? 'tech' : 'life'
  return tagOptions.value.filter(item => item.group === group)
})

const handleCategoryChange = () => {
  // 切换分类后清空已选标签，避免跨分类残留
  articleCreateFormData.tags = []
}

onMounted(async () => {
  const res = await articleTags()
  if (res.code === 0) {
    tagOptions.value = res.data as TagOption[]
  }
})

const handleSuccess = (res: ApiResponse<ImageUploadResponse>) => {
  if (res.code === 0) {
    articleCreateFormData.cover = res.data.url
    ElMessage.success(res.msg)
  }
}


const submitForm = async () => {
  const res = await articleCreate(articleCreateFormData)
  if (res.code === 0) {
    ElMessage.success(res.msg)
    layoutStore.state.articleCreateVisible = false
  }
}
</script>

<style scoped lang="scss">
.article-create-form {
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

      /* 标题输入框加大字号 */
      .form-item--title .el-input__wrapper {
        font-size: var(--fs-18);
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
