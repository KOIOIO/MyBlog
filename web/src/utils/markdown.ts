import MarkdownIt from "markdown-it";

/**
 * Agent 回复的 Markdown 渲染。
 * - html: false 禁止原始 HTML 注入，防止 XSS
 * - linkify 自动识别链接；breaks 支持换行
 * - ```mermaid 代码块渲染为 <div class="mermaid">，由 mermaid.run 负责成图
 */
const md: MarkdownIt = new MarkdownIt({
    html: false,
    linkify: true,
    breaks: true,
});

const defaultFence = md.renderer.rules.fence;

md.renderer.rules.fence = (tokens, idx, options, env, self) => {
    const token = tokens[idx];
    const info = token.info.trim();
    if (info === "mermaid") {
        return `<div class="mermaid">${token.content}</div>`;
    }
    return defaultFence ? defaultFence(tokens, idx, options, env, self) : "";
};

export function renderMarkdown(src: string): string {
    return md.render(src || "");
}
