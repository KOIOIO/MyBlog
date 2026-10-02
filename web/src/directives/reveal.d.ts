import type {Directive} from 'vue'
import type {RevealOptions} from './reveal'

declare module 'vue' {
  export interface ComponentCustomDirectives {
    /** 滚动进入视口块级 fade-up；v-reveal="{ stagger: 60 }" 启用子项级联 */
    reveal: Directive<HTMLElement, RevealOptions | undefined>
  }
}

export {}
