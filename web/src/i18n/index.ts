import {createI18n} from 'vue-i18n'

const STORAGE_KEY = 'myblog-locale'
const SUPPORTED_LOCALES = ['zh-CN', 'en', 'zh-TW'] as const
type Locale = (typeof SUPPORTED_LOCALES)[number]

const DEFAULT_LOCALE: Locale = 'zh-CN'

// 合并 ./locales/<lang>/<module>.ts 下的全部模块
const modules = import.meta.glob('./locales/*/*.ts', {eager: true}) as Record<
    string,
    {default: Record<string, any>}
>

const messages: Record<string, any> = {}
for (const path of Object.keys(modules)) {
    // path 形如 ./locales/zh-CN/base.ts
    const segments = path.split('/')
    const lang = segments[2]
    if (!lang) continue
    if (!messages[lang]) messages[lang] = {}
    // 模块间不重复定义同 key；后加载者覆盖先加载者
    messages[lang] = {...messages[lang], ...modules[path].default}
}

function detectInitialLocale(): Locale {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored && (SUPPORTED_LOCALES as readonly string[]).includes(stored)) {
        return stored as Locale
    }
    return DEFAULT_LOCALE
}

const initialLocale = detectInitialLocale()

export const i18n = createI18n({
    legacy: false,
    locale: initialLocale,
    fallbackLocale: DEFAULT_LOCALE,
    messages: messages as Record<string, any>,
})

export function setLocale(locale: string) {
    if (!(SUPPORTED_LOCALES as readonly string[]).includes(locale)) {
        locale = DEFAULT_LOCALE
    }
    localStorage.setItem(STORAGE_KEY, locale)
    ;(i18n.global.locale as {value: string}).value = locale
    document.documentElement.lang = locale
}

export function getLocale(): string {
    return (i18n.global.locale as {value: string}).value
}

// 初始化 <html lang>
document.documentElement.lang = initialLocale
