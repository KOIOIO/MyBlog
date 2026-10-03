<template>
  <div class="friend-link-page">
    <web-navbar :noScroll="true"/>
    <div class="page">
      <h1 class="page-title">{{ t('pages.friendLink.title') }}</h1>
      <friend-link-marquee v-if="friendLinkList.length" :items="friendLinkList"/>
      <div class="list">
        <div v-for="(item, i) in friendLinkList" :key="item.name" class="link-card"
             :style="{ animationDelay: Math.min(i, 5) * 60 + 'ms' }"
             @click="handleFriendLinkJumps(item.link)">
          <img class="card-logo" :src="item.logo" :alt="item.name" loading="lazy"/>
          <div class="card-body">
            <h3 class="card-name">{{ item.name }}</h3>
            <p class="card-desc">{{ item.description }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import WebNavbar from "@/components/layout/WebNavbar.vue";
import FriendLinkMarquee from "@/components/pages/FriendLinkMarquee.vue";
import {type FriendLink, friendLinkInfo} from "@/api/friend-link";
import {ref} from "vue";
import {useI18n} from "vue-i18n";

const {t} = useI18n();

const friendLinkList = ref<FriendLink[]>([])

const getFriendLinkInfo = async () => {
  const res = await friendLinkInfo()
  if (res.code === 0) {
    friendLinkList.value = res.data.list
  }
}

getFriendLinkInfo()

const handleFriendLinkJumps = (link: string) => {
  window.open(link)
}
</script>

<style scoped lang="scss">
.friend-link-page {
  background-color: var(--bg);
  min-height: 100vh;

  .page {
    max-width: var(--content-width);
    margin: 0 auto;
    padding: calc(70px + var(--sp-6)) var(--sp-4) var(--sp-9);
  }

  .page-title {
    font-size: var(--fs-30);
    font-weight: 700;
    color: var(--text-primary);
    margin-bottom: var(--sp-6);
  }

  .list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: var(--sp-4);
  }

  .link-card {
    display: flex;
    align-items: flex-start;
    gap: var(--sp-3);
    padding: var(--sp-4);
    background: var(--bg-elevated);
    border: 1px solid var(--border);
    border-radius: var(--radius-md);
    cursor: pointer;
    transition: border-color 150ms cubic-bezier(.16,1,.3,1),
                box-shadow 150ms cubic-bezier(.16,1,.3,1),
                transform 150ms cubic-bezier(.16,1,.3,1);
    animation: kf-fade-up var(--dur-mid) var(--ease-out) backwards;

    &:hover {
      border-color: var(--el-border-color-dark);
      box-shadow: var(--shadow-md);
      transform: translateY(-2px);
    }

    &:active {
      transform: translateY(1px);
      box-shadow: var(--shadow-sm);
      transition-duration: 80ms;
    }

    .card-logo {
      height: 48px;
      max-width: 96px;
      object-fit: contain;
      border-radius: var(--radius-sm);
      flex-shrink: 0;
    }

    .card-body {
      min-width: 0;

      .card-name {
        font-size: var(--fs-16);
        font-weight: 600;
        color: var(--text-primary);
        margin: 0 0 var(--sp-1);
        transition: color 150ms ease-out;
      }

      &:hover .card-name {
        color: var(--accent);
      }

      .card-desc {
        font-size: var(--fs-14);
        color: var(--text-muted);
        line-height: var(--lh-body);
        margin: 0;
      }
    }
  }
}

</style>
