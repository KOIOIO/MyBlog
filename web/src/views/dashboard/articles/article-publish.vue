<template>
  <div class="article-publish">
    <div class="title">
      <div class="left">
        <el-form :inline="true">
          <el-form-item :label="t('dashboard.articles.publish.title')">
            <el-input v-model="title" :placeholder="t('dashboard.articles.publish.titlePh')" clearable/>
          </el-form-item>
        </el-form>
      </div>
      <div class="right">
        <el-text>{{ t('dashboard.articles.publish.autoSave') }}</el-text>
        <el-switch v-model="isAutoSaveEnabled"/>
        <el-button type="danger" icon="Delete" @click="title='';text=''">{{ t('dashboard.articles.publish.clear') }}</el-button>
        <el-button type="success" icon="Plus" @click="layoutStore.state.articleCreateVisible=true">{{ t('dashboard.articles.publish.publish') }}</el-button>


        <!--  这里必须销毁，不然不会重新加载props-->
        <el-dialog
            v-model="articleCreateVisible"
            width="500"
            align-center
            destroy-on-close
            :before-close="articleCreateVisibleSynchronization"
        >
          <template #header>
            {{ t('dashboard.articles.publish.publish') }}
          </template>
          <article-create-form :title=title :content="text"/>
          <template #footer>
          </template>
        </el-dialog>
      </div>
    </div>
    <MdEditor v-model="text" @onUploadImg="onUploadImg" @onSave="onSave" @onChange="onChange"/>
  </div>
</template>

<script setup lang="ts">
import {ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {MdEditor} from 'md-editor-v3';
import 'md-editor-v3/lib/style.css';
import axios from "axios";
import {useLayoutStore} from "@/stores/layout";
import ArticleCreateForm from "@/components/forms/ArticleCreateForm.vue";
import type {AxiosResponse} from "axios";
import type {ApiResponse} from "@/utils/request";
import type {ImageUploadResponse} from "@/api/image";

const {t} = useI18n()
const layoutStore = useLayoutStore()

const savedIsAutoSaveEnabled = localStorage.getItem('isAutoSaveEnabled');
const isAutoSaveEnabled = savedIsAutoSaveEnabled ? ref(savedIsAutoSaveEnabled === 'true') : ref(true)
watch(() => isAutoSaveEnabled.value, (newIsAutoSaveEnabled) => {
  localStorage.setItem('isAutoSaveEnabled', String(newIsAutoSaveEnabled))
})

const title = ref('')
const savedArticle = localStorage.getItem('article')
const text = savedArticle ? ref(savedArticle) : ref('')

const onUploadImg = async (files: File[], callback: (urls: string[]) => void): Promise<void> => {
  const res = await Promise.all(
      files.map((file) => {
        return new Promise<AxiosResponse<ApiResponse<ImageUploadResponse>>>((resolve, reject) => {
          const form = new FormData();
          form.append('image', file);

          axios
              .post('/api/image/upload', form, {
                headers: {
                  'Content-Type': 'multipart/form-data',
                },
                withCredentials: true,
              })
              .then((response) => resolve(response))
              .catch((error) => reject(error));
        });
      })
  );

  callback(res.map((item) => item.data.data.url));
};

const onSave = (v: string, _: Promise<string>):void => {
  localStorage.setItem('article', v)
};

const onChange = (v: string):void => {
  if (isAutoSaveEnabled.value) {
    onSave(v,Promise.resolve(''))
  }
}

const articleCreateVisible = ref(layoutStore.state.articleCreateVisible);
watch(
    () => layoutStore.state.articleCreateVisible,
    (newValue) => {
      articleCreateVisible.value = newValue
    }
)

const articleCreateVisibleSynchronization = () => {
  layoutStore.state.articleCreateVisible = false
}
</script>

<style lang="scss">
.article-publish {
  height: 100%;
  display: flex;
  flex-direction: column;

  .title {
    display: flex;
    align-items: center;
    gap: var(--sp-4);
    margin-bottom: var(--sp-4);

    .left {
      flex: 1;
      min-width: 0;

      .el-form {
        display: flex;

        .el-form-item {
          margin-bottom: 0;
          flex: 1;
        }
      }

      /* 顶部标题输入框：无边框、大号 */
      .el-input__wrapper {
        box-shadow: none !important;
        background: transparent;
        padding: var(--sp-2) 0;
        border-radius: 0;

        &.is-focus {
          box-shadow: none !important;
        }
      }

      .el-input__inner {
        font-size: 28px;
        font-weight: 600;
        color: var(--text-primary);

        &::placeholder {
          color: var(--text-muted);
        }
      }
    }

    .right {
      display: flex;
      align-items: center;

      .el-text {
        color: var(--text-muted);
        font-size: var(--fs-14);
      }

      .el-switch {
        margin: 0 var(--sp-3);
      }
    }
  }

  /* 编辑器容器：elevated 卡片、细线、圆角 */
  .md-editor {
    flex: 1;
    border: 1px solid var(--border);
    border-radius: var(--radius-md);
    overflow: hidden;
    background-color: var(--bg-elevated);
    color: var(--text-body);
  }
}
</style>

<style lang="scss">
.article-publish .md-editor .md-editor-toolbar-wrapper .md-editor-toolbar svg.md-editor-icon {
  height: 20px;
  width: 20px;
  fill: currentColor;
}
</style>
