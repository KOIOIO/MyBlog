<template>
  <el-dialog v-model="visible" :title="t('pages.agent.selectArticle')" width="580px" append-to-body>
    <el-select
        v-model="selected"
        multiple
        filterable
        remote
        :remote-method="search"
        :loading="searching"
        :placeholder="t('pages.agent.searchArticlePlaceholder')"
        style="width: 100%"
    >
      <el-option v-for="item in options" :key="item.id" :label="item.title" :value="item.id"/>
    </el-select>
    <div class="tip">
      {{ t('pages.agent.selectedCount', {n: selected.length}) }}
    </div>
    <template #footer>
      <el-button @click="visible = false">{{ t('common.cancel') }}</el-button>
      <el-button type="primary" :disabled="selected.length === 0" @click="confirm">{{ t('common.confirm') }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import {ref} from "vue";
import {useI18n} from "vue-i18n";
import {articleSearch, type Article} from "@/api/article";

const {t} = useI18n();

const visible = ref(false);
const selected = ref<number[]>([]);
const options = ref<Array<{id: number; title: string}>>([]);
const searching = ref(false);

const emit = defineEmits<{ (e: 'confirm', ids: number[]): void }>();

const open = (initIds: number[] = []): void => {
    selected.value = [...initIds];
    options.value = [];
    visible.value = true;
};

const search = async (keyword: string): Promise<void> => {
    if (!keyword) {
        options.value = [];
        return;
    }
    searching.value = true;
    try {
        const res = await articleSearch({
            page: 1,
            page_size: 20,
            query: keyword,
            category: '',
            tag: '',
            sort: '',
            order: '',
        });
        options.value = res.data.list.map((hit) => ({
            id: Number(hit._id),
            title: hit._source.title,
        }));
    } finally {
        searching.value = false;
    }
};

const confirm = (): void => {
    emit('confirm', selected.value.slice(0, 5));
    visible.value = false;
};

defineExpose({open});
</script>

<style scoped lang="scss">
.tip {
  margin-top: 8px;
  font-size: 12px;
  color: var(--text-secondary);
}
</style>
