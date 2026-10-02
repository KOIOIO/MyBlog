<template>
  <div class="scroll-progress" :class="{ hidden: !visible }" :style="{ transform: `scaleX(${progress})` }" aria-hidden="true"/>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";

/**
 * 阅读进度条（§5）：
 * - fixed top 2px 高，accent 色，transform scaleX(transform-origin:left) 跟手，无 transition；
 * - 滚动监听 rAF 节流；
 * - 进度 <2% 或不可滚动时 opacity 0（150ms 过渡）。
 */
const progress = ref(0);
const visible = ref(false);
let rafId = 0;
let pending = false;

const update = () => {
  pending = false;
  const doc = document.documentElement;
  const scrollTop = window.scrollY || doc.scrollTop || 0;
  const max = (doc.scrollHeight || 0) - (doc.clientHeight || 0);
  const p = max > 0 ? Math.min(1, Math.max(0, scrollTop / max)) : 0;
  progress.value = p;
  // 进度 <2% 或页面不可滚动时隐藏
  visible.value = max > 0 && p >= 0.02;
};

const onScroll = () => {
  if (pending) return;
  pending = true;
  rafId = requestAnimationFrame(update);
};

onMounted(() => {
  update();
  window.addEventListener("scroll", onScroll, { passive: true });
  window.addEventListener("resize", onScroll, { passive: true });
});

onUnmounted(() => {
  window.removeEventListener("scroll", onScroll);
  window.removeEventListener("resize", onScroll);
  if (rafId) cancelAnimationFrame(rafId);
});
</script>

<style scoped lang="scss">
.scroll-progress {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background-color: var(--accent);
  transform-origin: left center;
  /* 跟手：transform 无 transition；仅 opacity 150ms 淡入淡出 */
  transition: opacity 150ms var(--ease-out);
  opacity: 1;
  z-index: 200;
  pointer-events: none;

  &.hidden {
    opacity: 0;
  }
}
</style>
