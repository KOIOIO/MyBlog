<template>
  <div class="friend-link-marquee">
    <div class="marquee-track">
      <!-- 内容复制为两组，配合 translateX(-50%) 实现无缝循环 -->
      <div v-for="(group, g) in [0, 1]" :key="g" class="marquee-group">
        <div
          v-for="(item, i) in items"
          :key="`${item.name}-${g}-${i}`"
          class="marquee-item"
          :title="item.name"
          @click="open(item.link)"
        >
          <img class="marquee-logo" :src="item.logo" :alt="item.name" loading="lazy"/>
          <span class="marquee-name">{{ item.name }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { FriendLink } from '@/api/friend-link'

defineProps<{ items: FriendLink[] }>()

const open = (link: string) => {
  window.open(link)
}
</script>

<style scoped lang="scss">
.friend-link-marquee {
  overflow: hidden;
  margin-bottom: var(--sp-6);
  padding: var(--sp-2) 0;
  -webkit-mask-image: linear-gradient(90deg, transparent, #000 10%, #000 90%, transparent);
  mask-image: linear-gradient(90deg, transparent, #000 10%, #000 90%, transparent);

  &:hover .marquee-track {
    animation-play-state: paused;
  }

  .marquee-track {
    display: flex;
    width: max-content;
    animation: marquee-scroll 32s linear infinite;
  }

  .marquee-group {
    display: flex;
  }

  .marquee-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--sp-1);
    margin-right: var(--sp-8);
    cursor: pointer;

    .marquee-logo {
      height: 44px;
      max-width: 140px;
      object-fit: contain;
      filter: grayscale(1);
      opacity: 0.75;
      transition: filter 200ms ease-out, opacity 200ms ease-out;
    }

    .marquee-name {
      max-width: 140px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      font-size: var(--fs-12);
      color: var(--text-muted);
    }

    &:hover .marquee-logo {
      filter: grayscale(0);
      opacity: 1;
    }
  }

  @media (max-width: 768px) {
    .marquee-logo {
      height: 36px;
    }
    .marquee-item {
      margin-right: var(--sp-6);
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .marquee-track {
      animation: none;
      flex-wrap: wrap;
      justify-content: center;
    }
    // 隐藏复制出来的第二组，只展示一轮
    .marquee-group:nth-child(2) {
      display: none;
    }
  }
}

@keyframes marquee-scroll {
  from {
    transform: translateX(0);
  }
  to {
    transform: translateX(-50%);
  }
}
</style>
