import {defineStore} from 'pinia';
import {ref} from 'vue';
import {
    createConversation as apiCreate,
    listConversations as apiList,
    listMessages as apiMessages,
    deleteConversation as apiDelete,
    chatStream,
} from '@/api/agent';
import type {AgentConversation, AgentMessage} from '@/api/agent';

export const useAgentStore = defineStore('agent', () => {
    const conversations = ref<AgentConversation[]>([]);
    const currentConversationId = ref<number>(0);
    const messages = ref<AgentMessage[]>([]);
    const streaming = ref<boolean>(false);
    const abortController = ref<AbortController | null>(null);

    const loadConversations = async (): Promise<void> => {
        const res = await apiList();
        conversations.value = res.data;
    };

    const createConversation = async (): Promise<AgentConversation> => {
        const res = await apiCreate();
        await loadConversations();
        return res.data;
    };

    const openConversation = async (id: number): Promise<void> => {
        currentConversationId.value = id;
        const res = await apiMessages(id);
        messages.value = res.data;
    };

    const removeConversation = async (id: number): Promise<void> => {
        await apiDelete(id);
        if (currentConversationId.value === id) {
            currentConversationId.value = 0;
            messages.value = [];
        }
        await loadConversations();
    };

    const stopStreaming = (): void => {
        abortController.value?.abort();
        streaming.value = false;
    };

    /**
     * 发送消息：本地即时追加 user 消息与 assistant 占位，流式更新占位内容。
     * 发送结束若为新建会话（id 为 0），刷新列表并回填会话 ID。
     */
    const sendMessage = async (content: string, articleIds: number[]): Promise<void> => {
        streaming.value = true;
        const now = new Date().toISOString();
        const userMsg: AgentMessage = {
            id: Date.now(),
            conversation_id: currentConversationId.value,
            role: 'user',
            content,
            article_ids: articleIds,
            created_at: now,
        };
        const asstId = Date.now() + 1;
        const asstMsg: AgentMessage = {
            id: asstId,
            conversation_id: currentConversationId.value,
            role: 'assistant',
            content: '',
            article_ids: [],
            created_at: now,
        };
        messages.value.push(userMsg, asstMsg);

        const controller = new AbortController();
        abortController.value = controller;
        let full = '';
        try {
            await chatStream(
                {conversationId: currentConversationId.value || undefined, message: content, articleIds},
                (delta) => {
                    full += delta;
                    const last = messages.value[messages.value.length - 1];
                    if (last && last.role === 'assistant' && last.id === asstId) {
                        last.content = full;
                    }
                },
                controller.signal,
            );
            if (currentConversationId.value === 0) {
                await loadConversations();
                if (conversations.value.length > 0) {
                    currentConversationId.value = conversations.value[0].id;
                    userMsg.conversation_id = currentConversationId.value;
                    const asst = messages.value.find((m) => m.id === asstId);
                    if (asst) {
                        asst.conversation_id = currentConversationId.value;
                    }
                }
            }
        } catch (e) {
            if (!(e instanceof Error && e.name === 'AbortError')) {
                // 失败时移除空的 assistant 占位，保留用户消息可重发
                messages.value = messages.value.filter((m) => !(m.id === asstId && m.content === ''));
                throw e;
            }
        } finally {
            streaming.value = false;
            abortController.value = null;
        }
    };

    return {
        conversations,
        currentConversationId,
        messages,
        streaming,
        loadConversations,
        createConversation,
        openConversation,
        removeConversation,
        sendMessage,
        stopStreaming,
    };
});
