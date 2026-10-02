<template>
  <div :class="{'web-navbar': true, show: isShow}">
    <scroll-progress/>
    <div class="container">
      <logo/>
      <div class="web-menu">
        <el-menu mode="horizontal" :ellipsis="false" :router="true" :default-active="$route.path">
          <template v-for="item in menuList">
            <el-menu-item :index="item.name"><span>{{ t(item.titleKey) }}</span></el-menu-item>
          </template>
        </el-menu>
      </div>
      <div class="nav-actions">
        <button class="icon-btn" :aria-label="t('nav.search')" @click="$router.push('/search')">
          <el-icon :size="18"><Search /></el-icon>
        </button>
        <theme-toggle/>
        <el-dropdown class="lang-switcher" trigger="click" @command="onLangChange">
          <button class="lang-trigger" :aria-label="t('common.language')">
            <el-icon :size="16"><Switch /></el-icon>
            <span class="lang-current">{{ currentLangLabel }}</span>
            <el-icon :size="12" class="lang-caret"><CaretBottom /></el-icon>
          </button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item
                  v-for="opt in langOptions"
                  :key="opt.value"
                  :command="opt.value"
                  :class="{ 'is-current': opt.value === currentLocale }"
              >
                <span class="opt-dot" :class="{ 'show': opt.value === currentLocale }"/>
                <span>{{ t(opt.labelKey) }}</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <auth-popover/>
      </div>
    </div>
    <back-to-top/>
  </div>
</template>

<script setup lang="ts">
import AuthPopover from "@/components/common/AuthPopover.vue";
import Logo from "@/components/widgets/Logo.vue";
import ThemeToggle from "@/components/common/ThemeToggle.vue";
import ScrollProgress from "@/components/common/ScrollProgress.vue";
import BackToTop from "@/components/common/BackToTop.vue";
import {Search, Switch, CaretBottom} from "@element-plus/icons-vue";
import {ref, computed} from "vue";
import {onUnmounted} from "vue";
import {useI18n} from "vue-i18n";
import {setLocale, getLocale} from "@/i18n";

const {t} = useI18n()

const isShow = ref(true)

const props = defineProps<{
  noScroll?: boolean
}>()

if (!props.noScroll) {
  isShow.value = false
  window.addEventListener("scroll", scroll)
  scroll()
}

function scroll() {
  let top = document.documentElement.scrollTop
  isShow.value = top >= 8;
}

onUnmounted(() => {
  if (!props.noScroll) {
    window.removeEventListener("scroll", scroll)
  }
})

interface MenuItem {
  titleKey: string;
  name: string;
}

const menuList: MenuItem[] = [
  {titleKey: "nav.home", name: "/"},
  {titleKey: "nav.search", name: "/search"},
  {titleKey: "nav.news", name: "/news"},
  {titleKey: "nav.forum", name: "/forum"},
  {titleKey: "nav.agent", name: "/agent"},
  {titleKey: "nav.friendLink", name: "/friend-link"}
]

const langOptions = [
  {value: 'zh-CN', labelKey: 'common.langZhCN'},
  {value: 'en', labelKey: 'common.langEn'},
  {value: 'zh-TW', labelKey: 'common.langZhTW'},
] as const

const currentLocale = ref(getLocale())
const currentLangLabel = computed(() => {
  const opt = langOptions.find(o => o.value === currentLocale.value)
  return opt ? t(opt.labelKey) : currentLocale.value
})

const onLangChange = (value: string) => {
  setLocale(value)
  currentLocale.value = getLocale()
}

</script>


<style scoped lang="scss">
.web-navbar {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 100;
  display: flex;
  justify-content: center;
  background-color: transparent;
  backdrop-filter: blur(0px);
  -webkit-backdrop-filter: blur(0px);
  transition: background-color 200ms ease-out, border-color 200ms ease-out, backdrop-filter 200ms ease-out;
  border-bottom: 1px solid transparent;

  &.show {
    background-color: color-mix(in srgb, var(--bg-elevated) 82%, transparent);
    backdrop-filter: blur(12px);
    -webkit-backdrop-filter: blur(12px);
    border-bottom-color: var(--border);
  }

  .container {
    display: flex;
    align-items: center;
    max-width: var(--content-width);
    width: 100%;
    height: 56px;
    padding: 0 var(--sp-5);

    .web-menu {
      margin-left: var(--sp-6);
      flex: 1;

      .el-menu {
        background-color: transparent;
        border-bottom: none;
        height: 56px;

        :deep(.el-menu-item) {
          background-color: transparent;
          border-bottom: none;
          color: var(--text-body);
          font-size: var(--fs-14);
          height: 56px;
          line-height: 56px;
          position: relative;

          &::after {
            content: "";
            position: absolute;
            left: 0;
            bottom: 0;
            width: 100%;
            height: 2px;
            background-color: var(--accent);
            transform: scaleX(0);
            transform-origin: left;
            transition: transform 150ms var(--ease-out);
          }

          &:hover,
          &.is-active {
            color: var(--accent);
            background-color: transparent;

            &::after {
              transform: scaleX(1);
            }
          }
        }
      }
    }

    .nav-actions {
      display: flex;
      align-items: center;
      gap: var(--sp-2);
      margin-left: auto;

      .icon-btn {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        width: 36px;
        height: 36px;
        border: none;
        border-radius: var(--radius-sm);
        background: transparent;
        color: var(--text-body);
        cursor: pointer;
        transition: background-color 150ms ease-out, color 150ms ease-out;

        &:hover {
          background: var(--accent-weak);
          color: var(--accent);
        }

        &:focus-visible {
          outline: 2px solid var(--accent);
          outline-offset: 2px;
        }
      }

      .lang-switcher {
        display: inline-flex;
        align-items: center;

        .lang-trigger {
          display: inline-flex;
          align-items: center;
          gap: var(--sp-1);
          height: 36px;
          padding: 0 var(--sp-2);
          border: none;
          border-radius: var(--radius-sm);
          background: transparent;
          color: var(--text-body);
          font-size: var(--fs-14);
          cursor: pointer;
          transition: background-color 150ms ease-out, color 150ms ease-out;

          &:hover {
            background: var(--accent-weak);
            color: var(--accent);
          }

          &:focus-visible {
            outline: 2px solid var(--accent);
            outline-offset: 2px;
          }

          .lang-caret {
            color: var(--text-muted);
          }
        }
      }

      .auth-popover {
        margin-left: var(--sp-1);
      }
    }
  }
}
</style>

<style lang="scss">
.el-dropdown-menu__item.is-current {
  color: var(--accent);

  .opt-dot {
    background-color: var(--accent);
  }
}

.el-dropdown-menu__item .opt-dot {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  margin-right: var(--sp-2);
  background-color: transparent;
}
</style>
