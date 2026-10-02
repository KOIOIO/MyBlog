import type {Directive, DirectiveBinding} from 'vue'

/**
 * v-reveal —— 滚动进入视口块级 fade-up（Folio 交互改造规范 §4.2）
 *
 * 用法：
 *   <div v-reveal>...</div>                       块级 fade-up（16px / 250ms ease-out）
 *   <ul v-reveal="{ stagger: 60 }"><li/>...</ul>   子项级联：delay = 60ms × index，总级联 ≤ 240ms
 *
 * 兜底（不吞内容）：
 *   - prefers-reduced-motion / IO 不可用 / 挂载时元素已在视口 → 立即加 is-revealed，不播放动画
 *   - .reveal 类仅由本指令在挂载时添加；JS 完全不可用时元素保持原生可见
 */

export interface RevealOptions {
  /** 子项级联步长（ms），缺省 0 = 不级联 */
  stagger?: number
}

const REVEAL_DURATION = 250
const ROOT_MARGIN_BOTTOM = '-15%' // 触发线落在视口高度 85% 处
const MAX_STAGGER_DELAY = 240

function prefersReducedMotion(): boolean {
  if (typeof window === 'undefined' || !window.matchMedia) return false
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

function immediateShow(el: HTMLElement) {
  el.classList.add('is-revealed')
}

function setup(el: HTMLElement, binding: DirectiveBinding<RevealOptions | undefined>) {
  const opts = binding.value ?? {}
  const stagger = typeof opts.stagger === 'number' && opts.stagger > 0 ? opts.stagger : 0

  // 级联模式：观察容器，子项逐个 .reveal；否则自身就是目标
  const targets: HTMLElement[] = stagger
    ? Array.from(el.children).filter((c): c is HTMLElement => c instanceof HTMLElement)
    : [el]

  targets.forEach((t, i) => {
    t.classList.add('reveal')
    if (stagger) {
      const delay = Math.min(i * stagger, MAX_STAGGER_DELAY)
      if (delay > 0) t.style.transitionDelay = `${delay}ms`
    }
  })

  const showAll = () => targets.forEach(immediateShow)

  // 1) reduced-motion：立即显示
  if (prefersReducedMotion()) {
    showAll()
    return
  }
  // 2) IO 不可用：立即显示
  if (typeof IntersectionObserver === 'undefined') {
    showAll()
    return
  }
  // 3) 挂载时已在视口（越过 85% 触发线）：立即显示，不播动画
  const rect = el.getBoundingClientRect()
  if (rect.bottom > 0 && rect.top < window.innerHeight * 0.85) {
    showAll()
    return
  }

  const io = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (entry.isIntersecting) {
          showAll()
          io.unobserve(el)
          break
        }
      }
    },
    {rootMargin: `0px 0px ${ROOT_MARGIN_BOTTOM} 0px`, threshold: 0},
  )
  io.observe(el)
  ;(el as RevealElement).__revealIO = io
}

interface RevealElement extends HTMLElement {
  __revealIO?: IntersectionObserver
}

const revealDirective: Directive<HTMLElement, RevealOptions | undefined> = {
  mounted(el, binding) {
    try {
      setup(el, binding)
    } catch {
      // 任何异常都不能吞内容：去掉隐藏类、直接显示
      el.classList.remove('reveal')
      el.classList.add('is-revealed')
    }
  },
  unmounted(el) {
    const host = el as RevealElement
    if (host.__revealIO) {
      host.__revealIO.disconnect()
      delete host.__revealIO
    }
  },
}

export default revealDirective
