<template>
  <div class="publish-page">
    <web-navbar :noScroll="true"/>
    <div class="page">
      <div v-if="userStore.isLoggedIn" class="editor">
        <input
            v-model="title"
            class="title-input"
            :placeholder="t('pages.forum.publish.titlePlaceholder')"
            maxlength="80"
        />

        <div class="field">
          <label class="field-label">{{ t('pages.forum.publish.categoryLabel') }}</label>
          <el-radio-group v-model="category">
            <el-radio value="技术">{{ categoryLabel('技术') }}</el-radio>
            <el-radio value="生活">{{ categoryLabel('生活') }}</el-radio>
          </el-radio-group>
        </div>

        <div class="field">
          <label class="field-label">{{ t('pages.forum.publish.tagLabel') }}</label>
          <el-select
              v-model="tags"
              multiple
              :placeholder="t('pages.forum.publish.selectTagPlaceholder')"
              style="width: 100%"
            >
              <el-option
                  v-for="opt in tagOptions"
                  :key="opt.tag"
                  :label="tagLabel(opt.tag)"
                  :value="opt.tag"
              />
            </el-select>
        </div>

        <div class="field">
          <el-input
              v-model="content"
              type="textarea"
              :rows="10"
              :placeholder="t('pages.forum.publish.contentPlaceholder')"
          />
        </div>

        <div class="field">
          <label class="field-label">{{ t('pages.forum.publish.imageLabel') }}</label>
          <ForumImageUploader v-model="images"/>
        </div>

        <div class="actions">
          <el-button @click="onCancel">{{ t('common.cancel') }}</el-button>
          <el-button type="primary" :loading="submitting" @click="onPublish">{{ t('common.publish') }}</el-button>
        </div>
      </div>
      <div v-else class="need-login">
        <p class="need-login-text">{{ t('pages.forum.publish.needLoginText') }}</p>
        <el-button type="primary" @click="openLogin">{{ t('pages.forum.publish.goLogin') }}</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, ref} from "vue";
import {useRouter} from "vue-router";
import {ElMessage} from "element-plus";
import WebNavbar from "@/components/layout/WebNavbar.vue";
import ForumImageUploader from "@/components/forum/ForumImageUploader.vue";
import {forumPublish, forumTags, type BlogTag} from "@/api/forum";
import {useUserStore} from "@/stores/user";
import {useLayoutStore} from "@/stores/layout";
import {useI18n} from "vue-i18n";
import {categoryLabel, tagLabel} from "@/i18n/meta";

const {t} = useI18n();

const userStore = useUserStore();
const layoutStore = useLayoutStore();
const router = useRouter();

const title = ref('');
const category = ref('技术');
const tags = ref<string[]>([]);
const content = ref('');
const images = ref<string[]>([]);
const submitting = ref(false);

const allTags = ref<BlogTag[]>([]);

onMounted(async () => {
    const res = await forumTags();
    if (res.code === 0) {
        allTags.value = res.data || [];
    }
});

const tagOptions = computed(() => {
    const group = category.value === '技术' ? 'tech' : 'life';
    return allTags.value.filter(t => t.group === group);
});

const onCancel = () => {
    router.push({name: 'forum'});
};

const onPublish = async () => {
    if (!title.value.trim()) {
        ElMessage.warning(t('pages.forum.publish.titleRequired'));
        return;
    }
    if (!content.value.trim()) {
        ElMessage.warning(t('pages.forum.publish.contentRequired'));
        return;
    }
    submitting.value = true;
    try {
        const res = await forumPublish({
            title: title.value.trim(),
            content: content.value.trim(),
            category: category.value,
            tags: tags.value,
            images: images.value,
        });
        if (res.code === 0) {
            ElMessage.success(t('pages.forum.publish.publishSuccess'));
            router.push({name: 'forum-detail', params: {id: res.data.id}});
        }
    } finally {
        submitting.value = false;
    }
};

const openLogin = () => {
    layoutStore.state.popoverVisible = true;
    layoutStore.state.loginVisible = true;
};
</script>

<style scoped lang="scss">
.publish-page {
  background-color: var(--bg);
  min-height: 100vh;

  .page {
    max-width: var(--reading-width);
    margin: 0 auto;
    padding: calc(70px + var(--sp-6)) var(--sp-4) var(--sp-9);
  }

  .editor {
    display: flex;
    flex-direction: column;
    gap: var(--sp-5);
  }

  .title-input {
    width: 100%;
    border: none;
    outline: none;
    background: transparent;
    font-size: var(--fs-24);
    font-weight: 600;
    color: var(--text-primary);
    line-height: var(--lh-title);
    padding: var(--sp-2) 0;
    border-bottom: 1px solid var(--border);
    transition: border-color 150ms ease-out;

    &::placeholder {
      color: var(--text-muted);
    }

    &:focus {
      border-bottom-color: var(--accent);
    }
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: var(--sp-2);

    .field-label {
      font-size: var(--fs-12);
      color: var(--text-muted);
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }
  }

  .actions {
    display: flex;
    justify-content: flex-end;
    gap: var(--sp-2);
    margin-top: var(--sp-2);
  }

  .need-login {
    padding: var(--sp-9) 0;
    text-align: center;

    .need-login-text {
      color: var(--text-muted);
      font-size: var(--fs-14);
      margin-bottom: var(--sp-4);
    }
  }
}
</style>
