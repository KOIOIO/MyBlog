<template>
  <div class="article-publish">
    <div class="title">
      <div class="left">
        <el-form :inline="true">
          <el-form-item label="文章标题">
            <el-input v-model="title" placeholder="请输入文章标题" clearable/>
          </el-form-item>
        </el-form>
      </div>
      <div class="right">
        <el-text>自动保存</el-text>
        <el-switch v-model="isAutoSaveEnabled"/>
        <el-button type="danger" icon="Delete" @click="title='';text=''">清空文章</el-button>
        <el-button type="success" icon="Plus" @click="layoutStore.state.articleCreateVisible=true">发布文章</el-button>


        <!--  这里必须销毁，不然不会重新加载props-->
        <el-dialog
            v-model="articleCreateVisible"
            width="500"
            align-center
            destroy-on-close
            :before-close="articleCreateVisibleSynchronization"
        >
          <template #header>
            发布文章
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
import {MdEditor} from 'md-editor-v3';
import 'md-editor-v3/lib/style.css';
import axios from "axios";
import {useLayoutStore} from "@/stores/layout";
import ArticleCreateForm from "@/components/forms/ArticleCreateForm.vue";
import type {AxiosResponse} from "axios";
import type {ApiResponse} from "@/utils/request";
import type {ImageUploadResponse} from "@/api/image";

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
  background-color: #1e1e2f; /* 深蓝紫色背景 */

  .title {
    display: flex;
    color: #ffffff; /* 白色文字 */

    .left {
      .el-input {
        min-width: 400px;
        background-color: #2e2e4f; /* 深蓝紫色输入框背景 */
        color: #ffffff; /* 白色文字 */
      }
    }

    .right {
      margin-left: auto;

      .el-switch {
        margin-left: 20px;
        margin-right: 20px;
      }

      .el-button {
        background-color: #3e3e6f; /* 深蓝紫色按钮背景 */
        color: #ffffff; /* 白色文字 */
        border-color: #3e3e6f; /* 深蓝紫色边框 */
      }

      .el-button:hover {
        background-color: #5e5e8f; /* 浅蓝紫色按钮背景 */
        border-color: #5e5e8f; /* 浅蓝紫色边框 */
      }
    }
  }

  .md-editor {
    height: 95%;
    background-color: #2e2e4f; /* 深蓝紫色编辑器背景 */
    color: #ffffff; /* 白色文字 */
  }
}
</style>

<style lang="scss">
.article-publish .md-editor .md-editor-toolbar-wrapper .md-editor-toolbar svg.md-editor-icon {
  height: 24px;
  width: 24px;
  fill: #ffffff; /* 白色图标 */
}
</style>
