import type {ApiResponse} from "@/utils/request";
import service from "@/utils/request";

export interface ForumUser {
    uuid: string;
    username: string;
    avatar: string;
}

export interface ForumPost {
    id: number;
    user_id: number;
    title: string;
    content: string;
    category: string;
    tags: string[];
    images: string[];
    like_count: number;
    comment_count: number;
    view_count: number;
    created_at: string;
    user: ForumUser;
}

export interface ForumComment {
    id: number;
    post_id: number;
    parent_id: number;
    user_id: number;
    content: string;
    created_at: string;
    user: ForumUser;
    children?: ForumComment[];
}

export interface BlogTag {
    tag: string;
    group: 'tech' | 'life';
    number: number;
}

export interface ForumListParams {
    page?: number;
    page_size?: number;
    category?: string;
    tag?: string;
}

export interface ForumListData {
    list: ForumPost[];
    total: number;
}

export const forumList = (params: ForumListParams): Promise<ApiResponse<ForumListData>> => {
    return service({
        url: '/forum/list',
        method: 'get',
        params: params,
    });
}

// 后端 forum/detail 返回扁平结构：帖子字段展开在 data 顶层，comments 同级
export interface ForumDetailData extends ForumPost {
    comments: ForumComment[];
}

export const forumDetail = (params: { id: number }): Promise<ApiResponse<ForumDetailData>> => {
    return service({
        url: '/forum/detail',
        method: 'get',
        params: params,
    });
}

export const forumTags = (): Promise<ApiResponse<BlogTag[]>> => {
    return service({
        url: '/forum/tags',
        method: 'get',
    });
}

export interface ForumPublishData {
    title: string;
    content: string;
    category: string;
    tags: string[];
    images: string[];
}

export const forumPublish = (data: ForumPublishData): Promise<ApiResponse<{ id: number }>> => {
    return service({
        url: '/forum/publish',
        method: 'post',
        data: data,
    });
}

export const forumUpload = (file: File): Promise<ApiResponse<{ url: string }>> => {
    const formData = new FormData();
    formData.append('image', file);
    return service({
        url: '/forum/upload',
        method: 'post',
        data: formData,
        headers: {
            'Content-Type': 'multipart/form-data',
        },
    });
}

export const forumLike = (data: { post_id: number }): Promise<ApiResponse<{ liked: boolean; like_count: number }>> => {
    return service({
        url: '/forum/like',
        method: 'post',
        data: data,
    });
}

export interface ForumCommentData {
    post_id: number;
    parent_id: number;
    content: string;
}

export const forumComment = (data: ForumCommentData): Promise<ApiResponse<{ id: number }>> => {
    return service({
        url: '/forum/comment',
        method: 'post',
        data: data,
    });
}
