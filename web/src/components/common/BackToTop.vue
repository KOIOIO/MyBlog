<template>
  <button
    v-if="visible"
    class="back-to-top"
    :class="{ show: shown }"
    aria-label="Back to top"
    @click="toTop"
  >
    <el-icon :size="18"><Top /></el-icon>
  </button>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { Top } from "@element-plus/icons-vue";

/**
 * 回顶按钮（§5）：
 * - fixed bottom-right 24px，40px 圆形（bg-elevated + 1px border + shadow-md）；
 * - 滚动 >600px 出现（opacity + translateY(4px) 150ms var(--ease-out)）；
 * - 点击 window.scrollTo({top:0, behavior:'smooth'})；
 * - rAF 节流滚动监听。
 */
const visible = ref(false); // 是否渲染（>600px）
const shown = ref(false);    // 用于入场过渡
let rafId = 0;
let pending = false;

const update = () => {
  pending = false;
  const y = window.scrollY || document.documentElement.scrollTop || 0;
  if (y > 600) {
    if (!visible.value) visible.value = true;
    // 下一帧再置 shown 以触发过渡
    requestAnimationFrame(() => { shown.value = true; });
  } else {
    shown.value = false;
    visible.value = false;
  }
};

const onScroll = () => {
  if (pending) return;
  pending = true;
  rafId = requestAnimationFrame(update);
};

const toTop = () => {
  window.scrollTo({ top: 0, behavior: "smooth" });
};

onMounted(() => {
  update();
  window.addEventListener("scroll", onScroll, { passive: true });
});

onUnmounted(() => {
  window.removeEventListener("scroll", onScroll);
  if (rafId) cancelAnimationFrame(rafId);
});
</script>

<style scoped lang="scss">
.back-to-top {
  position: fixed;
  right: 24px;
  bottom: 24px;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  box-shadow: var(--shadow-md);
  color: var(--text-body);
  cursor: pointer;
  z-index: 150;
  opacity: 0;
  transform: translateY(4px);
  transition: opacity 150ms var(--ease-out),
    transform 150ms var(--ease-out),
    color 150ms var(--ease-out);

  &.show {
    opacity: 1;
    transform: translateY(0);
  }

  &:hover {
    color: var(--accent);
  }

  &:active {
    transform: translateY(0) scale(0.96);
  }

  &:focus-visible {
    outline: 2px solid var(--accent);
    outline-offset: 2px;
  }
}
</style>
