<template>
    <n-card>
        <n-split direction="horizontal" v-model:size="split" style="height: calc(100vh - 14rem)" :max="0.45"
            :min="0.25">
            <template #1>
                <n-list>
                    <template #header>
                        <div class="flex justify-center items-center gap-2">
                            <span class="flex-1 text-xl font-semibold">好友列表</span>
                            <div class="mr-4">
                                <n-button text @click="showModal = true">
                                    <UserPlus2 :size="16" />
                                </n-button>
                            </div>
                        </div>
                        <n-modal v-model:show="showModal">
                            <n-card style="width: 37.5rem" title="搜索用户" :bordered="false" size="huge" role="dialog"
                                aria-modal="true">
                                <template #header-extra> </template>
                                <UserSearch />
                                <n-divider />
                                <n-spin :show="requestsLoading">
                                    <n-empty v-if="!requestsLoading && requestList.length === 0" description="暂无好友请求" />

                                    <div v-else class="space-y-3">
                                        <div v-for="request in requestList" :key="request.id"
                                            class="p-4 rounded-lg border" :style="{
                                                backgroundColor: 'var(--card-bg)',
                                                borderColor: request.status === 'rejected' ? 'var(--primary-color)' : 'var(--border-color)'
                                            }">
                                            <div class="flex items-center gap-4">
                                                <n-avatar :size="48" :src="request.friend_avatar || undefined">
                                                    <User v-if="!request.friend_avatar" />
                                                </n-avatar>

                                                <div class="flex-1 min-w-0">
                                                    <div class="flex items-center gap-2 mb-1">
                                                        <h3 class="font-medium"
                                                            :style="{ color: 'var(--text-primary)' }">
                                                            {{ request.friend_nickname || request.friend_username }}
                                                        </h3>
                                                        <span class="text-xs px-2 py-0.5 rounded" :style="{
                                                            backgroundColor: 'rgba(24, 160, 88, 0.1)',
                                                            color: '#18a058'
                                                        }">
                                                            Rating: {{ 1000 }}
                                                        </span>
                                                        <n-tag v-if="request.status === 'pending'" type="warning"
                                                            size="small">
                                                            待处理
                                                        </n-tag>
                                                        <n-tag v-else-if="request.status === 'accepted'" type="success"
                                                            size="small">
                                                            已接受
                                                        </n-tag>
                                                        <n-tag v-else type="error" size="small">
                                                            已拒绝
                                                        </n-tag>
                                                    </div>
                                                    <p class="text-sm mb-2" :style="{ color: 'var(--text-secondary)' }">
                                                        {{ request.message || '无验证消息' }}
                                                    </p>
                                                    <p class="text-xs" :style="{ color: 'var(--text-tertiary)' }">
                                                        {{ formatRelativeTime(request.created_at) }}
                                                    </p>
                                                </div>

                                                <div v-if="request.status === 'pending'" class="flex gap-2">
                                                    <n-button type="success" size="small"
                                                        @click="handleRequest(request.id, 'accepted')">
                                                        接受
                                                    </n-button>
                                                    <n-button type="error" size="small"
                                                        @click="handleRequest(request.id, 'rejected')">
                                                        拒绝
                                                    </n-button>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </n-spin>
                            </n-card>
                        </n-modal>
                    </template>
                    <n-list-item v-for="(value, index) in friendList" :key="value.id">
                        <div class="flex items-center gap-4 p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors cursor-pointer"
                            @click="startChat(value.friend_id)">
                            <n-avatar :size="48" :src="value.friend_avatar || undefined"
                                @click.stop="viewUserProfile(value.friend_id)">
                                <User v-if="!value.friend_avatar" />
                            </n-avatar>

                            <div class="flex-1 min-w-0">
                                <div class="flex items-center gap-2 mb-1">
                                    <h3 class="font-medium truncate" :style="{ color: 'var(--text-primary)' }">
                                        {{ value.remark || value.friend_nickname || value.friend_username }}
                                    </h3>
                                    <span v-if="value.remark" class="text-xs"
                                        :style="{ color: 'var(--text-tertiary)' }">
                                        ({{ value.friend_username }})
                                    </span>
                                </div>
                                <p class="text-sm truncate" :style="{ color: 'var(--text-secondary)' }">
                                    {{ value.friend_nickname && value.remark ? value.friend_nickname : '添加于 ' +
                                        formatRelativeTime(value.created_at) }}
                                </p>
                            </div>
                        </div>
                    </n-list-item>
                </n-list>
            </template>
            <template #2>
                <!-- 聊天页面 -->
                <div class="flex flex-col h-full">
                    <!-- 聊天头部 -->
                    <div class="flex items-center justify-between p-4 border-b"
                        :style="{ borderColor: 'var(--border-color)' }">
                        <div v-if="currentChatFriend" class="flex items-center gap-3">
                            <n-avatar :size="40" :src="currentChatFriend.friend_avatar || undefined">
                                <User v-if="!currentChatFriend.friend_avatar" />
                            </n-avatar>
                            <div>
                                <h2 class="text-lg font-semibold" :style="{ color: 'var(--text-primary)' }">
                                    {{ currentChatFriend.remark || currentChatFriend.friend_nickname || currentChatFriend.friend_username }}
                                </h2>
                                <p v-if="currentChatFriend.remark" class="text-xs" :style="{ color: 'var(--text-secondary)' }">
                                    {{ currentChatFriend.friend_nickname || currentChatFriend.friend_username }}
                                </p>
                            </div>
                        </div>
                        <h2 v-else class="text-lg font-semibold" :style="{ color: 'var(--text-primary)' }">聊天页面</h2>
                    </div>

                    <!-- 聊天消息区域 -->
                    <div ref="messagesContainer" class="flex-1 overflow-y-auto p-4 space-y-4">
                        <!-- 加载更多触发器 -->
                        <div ref="loadMoreTrigger" class="flex justify-center py-2">
                            <n-spin v-if="isLoadingMore" size="small" />
                            <p v-else-if="!hasMoreMessages && chatMessages.length > 0" class="text-xs" :style="{ color: 'var(--text-tertiary)' }">
                                没有更多消息了
                            </p>
                        </div>
                        
                        <n-spin :show="messagesLoading">
                            <div v-if="currentChatFriendId && chatMessages.length > 0" class="space-y-4">
                                <div v-for="msg in chatMessages" :key="msg.id"
                                    :class="['flex', isCurrentUser(msg) ? 'justify-end' : 'justify-start']">
                                    <div :class="['flex max-w-75%', isCurrentUser(msg) ? 'flex-row-reverse' : 'flex-row']">
                                        <!-- 头像 -->
                                        <n-avatar :size="40"
                                            :src="isCurrentUser(msg) ? userStore.avatar : currentChatFriend?.friend_avatar || undefined"
                                            class="shrink-0">
                                            <User v-if="!isCurrentUser(msg) && !currentChatFriend?.friend_avatar" />
                                        </n-avatar>

                                        <!-- 消息内容 -->
                                        <div :class="['mx-3', isCurrentUser(msg) ? 'text-right' : 'text-left']">
                                            <div :class="['inline-block px-4 py-2 rounded-lg', isCurrentUser(msg)
                                                ? 'bg-blue-500 text-white'
                                                : 'bg-gray-200 dark:bg-gray-700']"
                                                :style="!isCurrentUser(msg) ? { backgroundColor: 'var(--card-bg)' } : {}">
                                                <p class="text-sm whitespace-pre-wrap wrap-break-word">{{ msg.content }}</p>
                                            </div>
                                            <p class="text-xs mt-1"
                                                :style="{ color: 'var(--text-tertiary)' }">
                                                {{ formatRelativeTime(msg.created_at) }}
                                            </p>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div v-else-if="currentChatFriendId && chatMessages.length === 0"
                                class="flex justify-center">
                                <p class="text-sm" :style="{ color: 'var(--text-tertiary)' }">暂无聊天记录</p>
                            </div>
                            <div v-else class="flex justify-center">
                                <p class="text-sm" :style="{ color: 'var(--text-tertiary)' }">选择一个好友开始聊天</p>
                            </div>
                        </n-spin>
                    </div>

                    <!-- 消息输入区域 -->
                    <div class="p-4 border-t" :style="{ borderColor: 'var(--border-color)' }">
                        <div class="flex gap-3">
                            <n-input type="textarea" placeholder="输入消息..." :autosize="{ minRows: 1, maxRows: 4 }"
                                disabled />
                            <n-button type="primary" disabled>发送</n-button>
                        </div>
                    </div>
                </div>
            </template>
        </n-split>
    </n-card>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import { User, MessageSquare, UserPlus2 } from 'lucide-vue-next'
import {
    NButton,
    NEmpty,
    NSplit,
    NSpin,
    NTabs,
    NCard,
    NDivider,
    NTabPane,
    NAvatar,
    NTag,
    NModal,
    NList,
    NListItem,
    NInput,
    useMessage,
} from 'naive-ui'
const showModal = ref(false)
const split = ref(0.2)
import { userApi } from '@/services/user'
import type { FriendShip, FriendShipRequest } from '@/types/user'
import { formatRelativeTime } from '@/utils/format'
import UserSearch from '@/components/UserSearch.vue'

const router = useRouter()
import { useIntersectionObserver } from '@vueuse/core'
const message = useMessage()

const friendsLoading = ref(false)
const requestsLoading = ref(false)

const friendList = ref<FriendShip[]>([])
const requestList = ref<FriendShipRequest[]>([])

// Methods
const loadFriendList = async () => {
    friendsLoading.value = true
    try {
        const res = await userApi.getFriendList()
        if (res.code === 200 && res.info) {
            friendList.value = res.info
        }
    } catch (error) {
        console.error('Failed to load friend list:', error)
    } finally {
        friendsLoading.value = false
    }
}

const loadFriendRequests = async () => {
    requestsLoading.value = true
    try {
        const res = await userApi.getFriendRequestList()
        if (res.code === 200 && res.info) {
            requestList.value = res.info
        }
    } catch (error) {
        console.error('Failed to load friend requests:', error)
    } finally {
        requestsLoading.value = false
    }
}

const handleRequest = async (requestId: string, status: "accepted" | "rejected") => {
    try {
        const res = await userApi.HandleFriendRequest(requestId, status)
        if (res.code === 200) {
            message.success(status === "accepted" ? '已接受好友请求' : '已拒绝好友请求')
            await loadFriendRequests()
            if (status === "accepted") {
                await loadFriendList()
            }
        } else {
            message.error(res.msg || '操作失败')
        }
    } catch (error) {
        console.error('Failed to handle friend request:', error)
        message.error('操作失败')
    }
}


const viewUserProfile = (userId: string) => {
    router.push({ name: 'UserHomePage', params: { id: userId } })
}

import type { ChatMessage } from '@/types/chat'
import { useUserStore } from '@/stores/useUserStore'

const userStore = useUserStore()
const page = ref(1)
const currentChatFriendId = ref<string | null>(null)
const chatMessages = ref<ChatMessage[]>([])
const currentChatFriend = ref<FriendShip | null>(null)
const messagesLoading = ref(false)
const messagesContainer = ref<HTMLElement | null>(null)
const loadMoreTrigger = ref<HTMLElement | null>(null)
const isLoadingMore = ref(false)
const hasMoreMessages = ref(true)
const allPagesLoaded = ref<Set<string>>(new Set())

const scrollToBottom = () => {
    nextTick(() => {
        if (messagesContainer.value) {
            messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
        }
    })
}
const loadMoreMessages = async () => {
    if (!currentChatFriendId.value || isLoadingMore.value || !hasMoreMessages.value) return
    
    // 防止重复加载同一页
    const pageKey = `${currentChatFriendId.value}_page_${page.value + 1}`
    if (allPagesLoaded.value.has(pageKey)) {
        hasMoreMessages.value = false
        return
    }
    
    isLoadingMore.value = true
    const oldScrollHeight = messagesContainer.value?.scrollHeight || 0
    const oldScrollTop = messagesContainer.value?.scrollTop || 0
    
    try {
        const nextPage = page.value + 1
        const res = await userApi.getChatRecordList(currentChatFriendId.value, nextPage)
        
        if (res.code === 200 && res.info) {
            if (res.info.length === 0) {
                hasMoreMessages.value = false
                allPagesLoaded.value.add(pageKey)
            } else {
                // 将新消息反转后添加到现有消息前面
                const newMessages = [...res.info].reverse()
                chatMessages.value = [...newMessages, ...chatMessages.value]
                page.value = nextPage
                allPagesLoaded.value.add(pageKey)
                
                // 恢复滚动位置
                await nextTick()
                if (messagesContainer.value) {
                    const newScrollHeight = messagesContainer.value.scrollHeight
                    messagesContainer.value.scrollTop = oldScrollTop + (newScrollHeight - oldScrollHeight)
                }
            }
        }
    } catch (error) {
        console.error('Failed to load more messages:', error)
    } finally {
        isLoadingMore.value = false
    }
}

// 设置 IntersectionObserver
useIntersectionObserver(
    loadMoreTrigger,
    ([{ isIntersecting }]) => {
        if (isIntersecting) {
            loadMoreMessages()
        }
    },
    {
        threshold: 0.1
    }
)

// 监听 currentChatFriendId 变化，重置分页状态
watch(currentChatFriendId, (newId, oldId) => {
    if (newId !== oldId) {
        page.value = 1
        hasMoreMessages.value = true
        allPagesLoaded.value.clear()
    }
})


const startChat = async (friend_id: string) => {
    currentChatFriendId.value = friend_id
    messagesLoading.value = true

    try {
        const res = await userApi.getChatRecordList(friend_id, page.value)
        if (res.code === 200 && res.info) {
            // 消息返回是倒序的（最新的在前），需要反转成正序显示
            chatMessages.value = [...res.info].reverse()
            // Find friend info from friendList
            currentChatFriend.value = friendList.value.find(f => f.friend_id === friend_id) || null
            // 滚动到底部显示最新消息
            scrollToBottom()
        }
    } catch (error) {
        console.error('Failed to load chat messages:', error)
        message.error('加载聊天记录失败')
    } finally {
        messagesLoading.value = false
    }
}

const isCurrentUser = (message: ChatMessage) => {
    return message.sender_id === userStore.id
}

onMounted(() => {
    loadFriendList()
    loadFriendRequests()
})
</script>
