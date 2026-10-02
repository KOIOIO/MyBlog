<template>
  <aside class="tag-sidebar">
    <div v-for="group in groups" :key="group.key" class="tag-group">
      <div class="group-title">{{ categoryLabel(group.title) }}</div>
      <div class="tag-list">
        <span
            v-for="item in group.items"
            :key="item.tag"
            class="tag-item chip"
            :class="{ active: item.tag === selectedTag }"
            @click="onSelect(item.tag)"
        >
          <span class="tag-name">{{ tagLabel(item.tag) }}</span>
          <span class="tag-num">{{ item.number }}</span>
        </span>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import {computed, onMounted, ref} from "vue";
import {forumTags, type BlogTag} from "@/api/forum";
import {useI18n} from "vue-i18n";
import {categoryLabel, tagLabel} from "@/i18n/meta";

const {t} = useI18n();

defineProps<{ selectedTag?: string }>();

const emit = defineEmits<{ (e: 'select', tag: string): void }>();

const tags = ref<BlogTag[]>([]);

const groups = computed(() => {
    const tech = tags.value.filter(t => t.group === 'tech').slice(0, 15);
    const life = tags.value.filter(t => t.group === 'life').slice(0, 15);
    return [
        {key: 'tech', title: '技术', items: tech},
        {key: 'life', title: '生活', items: life},
    ];
});

const loadTags = async () => {
    const res = await forumTags();
    if (res.code === 0) {
        tags.value = res.data || [];
    }
};

onMounted(loadTags);

const onSelect = (tag: string) => {
    emit('select', tag);
};
</script>

<style scoped lang="scss">
.tag-sidebar {
  width: 280px;
  flex-shrink: 0;
  position: sticky;
  top: 80px;
  align-self: flex-start;
  display: flex;
  flex-direction: column;
  gap: var(--sp-5);

  .tag-group {
    .group-title {
      font-size: var(--fs-12);
      font-weight: 500;
      color: var(--text-muted);
      text-transform: uppercase;
      letter-spacing: 0.05em;
      margin-bottom: var(--sp-3);
    }

    .tag-list {
      display: flex;
      flex-wrap: wrap;
      gap: var(--sp-2);
    }

    .tag-item {
      display: inline-flex;
      align-items: center;
      gap: var(--sp-1);
      font-size: var(--fs-12);
      color: var(--text-body);
      background: var(--accent-weak);
      padding: 3px var(--sp-2);
      border-radius: var(--radius-sm);
      cursor: pointer;
      transition: color 150ms ease-out, background-color 150ms ease-out, transform 80ms ease-out;

      .tag-num {
        color: var(--text-muted);
        font-size: var(--fs-12);
      }

      &:hover {
        color: var(--accent);
      }

      &:active {
        transform: translateY(1px);
      }

      &.active {
        background: var(--accent);
        color: var(--bg-elevated);

        .tag-num {
          color: var(--bg-elevated);
          opacity: 0.8;
        }
      }
    }
  }
}
</style>
