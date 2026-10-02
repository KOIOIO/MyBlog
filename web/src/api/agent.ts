import type {ApiResponse} from "@/utils/request";
import service from "@/utils/request";
import {useUserStore} from "@/stores/user";

export interface AgentConversation {
    id: number;
    uuid: string;
    title: string;
    created_at: string;
    updated_at: string;
}

export interface AgentMessage {
    id: number;
    conversation_id: number;
    role: 'user' | 'assistant';
    content: string;
    article_ids: number[];
    created_at: string;
}

export const createConversation = (): Promise<ApiResponse<AgentConversation>> => {
    return service({
        url: '/agent/conversations',
        method: 'post',
    });
}

export const listConversations = (): Promise<ApiResponse<AgentConversation[]>> => {
    return service({
        url: '/agent/conversations',
        method: 'get',
    });
}

export const listMessages = (conversationId: number): Promise<ApiResponse<AgentMessage[]>> => {
    return service({
        url: `/agent/conversations/${conversationId}/messages`,
        method: 'get',
    });
}

export const deleteConversation = (conversationId: number): Promise<ApiResponse<null>> => {
    return service({
        url: `/agent/conversations/${conversationId}`,
        method: 'delete',
    });
}

export interface ChatStreamPayload {
    conversationId?: number;
    message: string;
    articleIds: number[];
}

/**
 * 流式对话：原生 fetch + ReadableStream 解析 SSE（不走 axios 拦截器）。
 * 服务端逐块返回 data: {"delta": "..."}，结束 data: {"done": true}，错误 data: {"error": "..."}。
 */
export const chatStream = async (
    payload: ChatStreamPayload,
    onDelta: (delta: string) => void,
    signal?: AbortSignal,
): Promise<void> => {
    const userStore = useUserStore();
    const baseURL = import.meta.env.VITE_BASE_API || '';
    const res = await fetch(`${baseURL}/agent/chat`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'x-access-token': userStore.state.accessToken,
        },
        body: JSON.stringify({
            conversation_id: payload.conversationId ?? 0,
            message: payload.message,
            article_ids: payload.articleIds,
        }),
        signal,
    });
    if (!res.ok || !res.body) {
        throw new Error(`chat request failed: ${res.status}`);
    }
    const reader = res.body.getReader();
    const decoder = new TextDecoder();
    let buf = '';
    for (;;) {
        const {done, value} = await reader.read();
        if (done) {
            break;
        }
        buf += decoder.decode(value, {stream: true});
        const lines = buf.split('\n');
        buf = lines.pop() ?? '';
        for (const line of lines) {
            const s = line.trim();
            if (!s.startsWith('data:')) {
                continue;
            }
            const data = s.slice(5).trim();
            if (!data) {
                continue;
            }
            let json: {delta?: string; done?: boolean; error?: string};
            try {
                json = JSON.parse(data);
            } catch {
                continue;
            }
            if (json.error) {
                throw new Error(json.error);
            }
            if (json.delta) {
                onDelta(json.delta);
            }
            if (json.done) {
                return;
            }
        }
    }
}
