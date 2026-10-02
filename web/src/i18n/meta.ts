import { i18n } from './index'

// 展示层元数据（分类/标签/排序）三语言映射函数。
// 后端数据与 DB 原值（中文 / 英文标签 / 排序枚举）一律不动，仅在展示时翻译 label。
// 模板中直接调用这些函数即可随语言切换更新（i18n.global.t 读取 locale ref 具响应性）。

const t = (key: string): string => (i18n.global.t as (...args: any[]) => string)(key)

/**
 * 分类展示：cat 为 DB 原值 '技术' | '生活'（或其他未知值），返回当前语言 label。
 * 未命中映射时回退原值，保证不空白。
 */
export function categoryLabel(cat: string): string {
  const key = cat === '技术' ? 'tech' : cat === '生活' ? 'life' : ''
  if (!key) return cat
  const label = t(`meta.categories.${key}`)
  return label.startsWith('meta.categories.') ? cat : label
}

/**
 * 标签展示：tag 为 DB 原值（中文或英文），返回当前语言 label。
 * 英文标签（Go/Redis/Docker 等）未在语言包定义，t() 返回原 key，回退原值。
 */
export function tagLabel(tag: string): string {
  const label = t(`meta.tags.${tag}`)
  return label.startsWith('meta.tags.') ? tag : label
}

/**
 * 排序展示：sort 为枚举值 default|time|comment|view|like（或对应的中文 label），返回当前语言 label。
 * 未命中映射时回退原值。
 */
export function sortLabel(sort: string): string {
  const label = t(`meta.sorts.${sort}`)
  return label.startsWith('meta.sorts.') ? sort : label
}
