<template>
    <div>
        <div class="flex items-center justify-between mb-6">
            <h1 class="text-2xl font-bold" :style="{ color: 'var(--text-primary)' }">私信</h1>
        </div>

        <n-tabs v-model:value="activeTab" type="line">
            <!-- Friends Tab -->
            <n-tab-pane name="friends" tab="好友列表">
                <n-spin :show="friendsLoading">
                    <n-empty v-if="!friendsLoading && friendList.length === 0" description="暂无好友" />

                    <div v-else class="space-y-3">
                        <div v-for="friend in friendList" :key="friend.id"
                            class="flex items-center gap-4 p-4 rounded-lg transition-all hover:shadow-md cursor-pointer"
                            :style="{ backgroundColor: 'var(--hover-bg)' }" @click="viewUserProfile(friend.id)">
                            <n-avatar :size="48" :src="friend.avatar || undefined">
                                <User v-if="!friend.avatar" />
                            </n-avatar>

                            <div class="flex-1 min-w-0">
                                <div class="flex items-center gap-2 mb-1">
                                    <h3 class="font-medium truncate" :style="{ color: 'var(--text-primary)' }">
                                        {{ friend.nickname || friend.username }}
                                    </h3>
                                    <span class="text-xs px-2 py-0.5 rounded" :style="{
                                        backgroundColor: 'rgba(24, 160, 88, 0.1)',
                                        color: '#18a058'
                                    }">
                                        Rating: {{ friend.rating }}
                                    </span>
                                </div>
                                <p class="text-sm truncate" :style="{ color: 'var(--text-secondary)' }">
                                    {{ friend.introduction || '这个人很懒，什么都没留下' }}
                                </p>
                            </div>

                            <n-button text size="small" @click.stop="startChat">
                                <MessageSquare :size="18" />
                            </n-button>
                        </div>
                    </div>
                </n-spin>
            </n-tab-pane>

            <!-- Friend Requests Tab -->
            <n-tab-pane name="requests" tab="好友请求">
                <n-spin :show="requestsLoading">
                    <n-empty v-if="!requestsLoading && requestList.length === 0" description="暂无好友请求" />

                    <div v-else class="space-y-3">
                        <div v-for="request in requestList" :key="request.id" class="p-4 rounded-lg border" :style="{
                            backgroundColor: 'var(--card-bg)',
                            borderColor: request.status === 0 ? 'var(--primary-color)' : 'var(--border-color)'
                        }">
                            <div class="flex items-center gap-4">
                                <n-avatar :size="48" :src="request.avatar || undefined">
                                    <User v-if="!request.avatar" />
                                </n-avatar>

                                <div class="flex-1 min-w-0">
                                    <div class="flex items-center gap-2 mb-1">
                                        <h3 class="font-medium" :style="{ color: 'var(--text-primary)' }">
                                            {{ request.nickname || request.username }}
                                        </h3>
                                        <span class="text-xs px-2 py-0.5 rounded" :style="{
                                            backgroundColor: 'rgba(24, 160, 88, 0.1)',
                                            color: '#18a058'
                                        }">
                                            Rating: {{ request.rating }}
                                        </span>
                                        <n-tag v-if="request.status === 0" type="warning" size="small">
                                            待处理
                                        </n-tag>
                                        <n-tag v-else-if="request.status === 1" type="success" size="small">
                                            已接受
                                        </n-tag>
                                        <n-tag v-else type="error" size="small">
                                            已拒绝
                                        </n-tag>
                                    </div>
                                    <p class="text-sm mb-2" :style="{ color: 'var(--text-secondary)' }">
                                        {{ request.verification || '无验证消息' }}
                                    </p>
                                    <p class="text-xs" :style="{ color: 'var(--text-tertiary)' }">
                                        {{ formatTime(request.created_at) }}
                                    </p>
                                </div>

                                <div v-if="request.status === 0" class="flex gap-2">
                                    <n-button type="success" size="small" @click="handleRequest(request.id, true)">
                                        接受
                                    </n-button>
                                    <n-button type="error" size="small" @click="handleRequest(request.id, false)">
                                        拒绝
                                    </n-button>
                                </div>
                            </div>
                        </div>
                    </div>
                </n-spin>
            </n-tab-pane>

            <!-- Search Results Tab -->
            <n-tab-pane name="search" tab="搜索用户">
                <div class="mb-4 gap-2 flex items-center">
                    <n-input v-model:value="userSearchKeyword" placeholder="输入用户名或昵称搜索" clearable
                        />
                    <n-button type="info" @click="handleUserSearch">
                        <template #icon>
                            <Search :size="16" />
                        </template>
                    </n-button>
                </div>

                <n-spin :show="searchLoading">
                    <n-empty v-if="userSearchResults.length === 0 " description="请输入关键词搜索用户" />

                    <div v-else class="space-y-3">
                        <div v-for="user in userSearchResults" :key="user.id"
                            class="flex items-center gap-4 p-4 rounded-lg border"
                            :style="{ backgroundColor: 'var(--card-bg)' }">
                            <n-avatar :size="48" :src="user.avatar || undefined">
                                <User v-if="!user.avatar" />
                            </n-avatar>

                            <div class="flex-1 min-w-0">
                                <div class="flex items-center gap-2 mb-1 cursor-pointer" @click="$router.push({name : 'UserHomePage', params: { id: user.id }})">
                                    <h3 class="font-medium" :style="{ color: 'var(--text-primary)' }">
                                        {{ user.nickname || user.username }}
                                    </h3>
                                    <span class="text-xs px-2 py-0.5 rounded" :style="{
                                        backgroundColor: 'rgba(24, 160, 88, 0.1)',
                                        color: '#18a058'
                                    }">
                                        Rating: {{ user.rating }}
                                    </span>
                                </div>
                                <p class="text-sm truncate" :style="{ color: 'var(--text-secondary)' }">
                                    {{ user.introduction || '这个人很懒，什么都没留下' }}
                                </p>
                            </div>

                            <n-button type="primary" size="small" @click="showAddFriendModal(user)">
                                添加好友
                            </n-button>
                        </div>
                    </div>
                </n-spin>
            </n-tab-pane>
        </n-tabs>

        <!-- Add Friend Modal -->
        <n-modal v-model:show="showAddFriend" preset="dialog" title="添加好友">
            <div v-if="selectedUser" class="mb-4">
                <div class="flex items-center gap-3 mb-4">
                    <n-avatar :size="60" :src="selectedUser.avatar || undefined">
                        <User v-if="!selectedUser.avatar" />
                    </n-avatar>
                    <div>
                        <h3 class="font-medium" :style="{ color: 'var(--text-primary)' }">
                            {{ selectedUser.nickname || selectedUser.username }}
                        </h3>
                        <p class="text-sm" :style="{ color: 'var(--text-secondary)' }">
                            {{ selectedUser.introduction || '这个人很懒，什么都没留下' }}
                        </p>
                    </div>
                </div>
                <n-input v-model:value="verificationMessage" type="textarea" placeholder="输入验证消息（可选）"
                    :autosize="{ minRows: 2, maxRows: 4 }" />
            </div>
            <template #action>
                <n-button @click="showAddFriend = false">取消</n-button>
                <n-button type="primary" :loading="sendingRequest" @click="sendFriendRequest">
                    发送请求
                </n-button>
            </template>
        </n-modal>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { User, Search, MessageSquare } from 'lucide-vue-next'
import {
    NButton,
    NEmpty,
    NInput,
    NSpin,
    NTabs,
    NTabPane,
    NAvatar,
    NTag,
    NModal,
    useMessage
} from 'naive-ui'
import { userApi } from '@/services/user'
import type { User as UserType, FriendShip } from '@/types/user'

const router = useRouter()
const message = useMessage()

// State
const activeTab = ref('friends')
const searchKeyword = ref('')
const userSearchKeyword = ref('')
const verificationMessage = ref('')
const showAddFriend = ref(false)
const sendingRequest = ref(false)
const selectedUser = ref<UserType | null>(null)

const friendsLoading = ref(false)
const requestsLoading = ref(false)
const searchLoading = ref(false)

const friendList = ref<UserType[]>([])
const requestList = ref<(FriendShip & {
    username: string
    nickname: string
    avatar: string
    rating: number
    verification?: string
})[]>([])
const userSearchResults = ref<UserType[]>([])

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

const handleRequest = async (friendshipId: string, accept: boolean) => {
    try {
        const res = await userApi.HandleFriendRequest(friendshipId, accept)
        if (res.code === 200) {
            message.success(accept ? '已接受好友请求' : '已拒绝好友请求')
            await loadFriendRequests()
            if (accept) {
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

const handleUserSearch = async () => {
    if (!userSearchKeyword.value?.trim()) {
        console.warn('Search keyword is empty')
        return
    }

    searchLoading.value = true
    try {
        const res = await userApi.searchUser(userSearchKeyword.value)
        if (res.code === 200) {
            userSearchResults.value = res.info ?? []
        }
    } catch (error) {
        console.error('Failed to search users:', error)
        message.error('搜索失败')
    } finally {
        searchLoading.value = false
    }
}

const showAddFriendModal = (user: UserType) => {
    selectedUser.value = user
    verificationMessage.value = ''
    showAddFriend.value = true
}

const sendFriendRequest = async () => {
    if (!selectedUser.value) return

    sendingRequest.value = true
    try {
        const res = await userApi.FirendRequest(
            selectedUser.value.id,
            verificationMessage.value || undefined
        )
        if (res.code === 200) {
            message.success('好友请求已发送')
            showAddFriend.value = false
        } else {
            message.error(res.msg || '发送失败')
        }
    } catch (error) {
        console.error('Failed to send friend request:', error)
        message.error('发送失败')
    } finally {
        sendingRequest.value = false
    }
}

const viewUserProfile = (userId: string) => {
    router.push({ name: 'UserHomePage', params: { id: userId } })
}

const startChat = () => {
    // TODO: Implement chat functionality
    message.info('聊天功能开发中')
}

const formatTime = (timeStr: string) => {
    const date = new Date(timeStr)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const minutes = Math.floor(diff / 60000)
    const hours = Math.floor(diff / 3600000)
    const days = Math.floor(diff / 86400000)

    if (minutes < 1) return '刚刚'
    if (minutes < 60) return `${minutes}分钟前`
    if (hours < 24) return `${hours}小时前`
    if (days < 7) return `${days}天前`
    return date.toLocaleDateString()
}

onMounted(() => {
    loadFriendList()
    loadFriendRequests()
})
</script>
