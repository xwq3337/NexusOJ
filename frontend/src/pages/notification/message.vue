<template>
    <n-card>
        <n-split direction="horizontal" v-model:size="split" style="height: 200px" :max="0.5" :min="0.2">
            <template #1>
                <n-list>
                    <template #header>
                        <div class="flex justify-center items-center gap-2">
                            <span class="flex-1 text-xl font-semibold">好友列表</span>
                            <n-button circle class="border-none" @click="showModal = true">
                                <UserPlus2 :size="16" />
                            </n-button>
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
                                                borderColor: request.status === 0 ? 'var(--primary-color)' : 'var(--border-color)'
                                            }">
                                            <div class="flex items-center gap-4">
                                                <n-avatar :size="48" :src="request.avatar || undefined">
                                                    <User v-if="!request.avatar" />
                                                </n-avatar>

                                                <div class="flex-1 min-w-0">
                                                    <div class="flex items-center gap-2 mb-1">
                                                        <h3 class="font-medium"
                                                            :style="{ color: 'var(--text-primary)' }">
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
                                                        <n-tag v-else-if="request.status === 1" type="success"
                                                            size="small">
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
                                                        {{ formatRelativeTime(request.created_at) }}
                                                    </p>
                                                </div>

                                                <div v-if="request.status === 0" class="flex gap-2">
                                                    <n-button type="success" size="small"
                                                        @click="handleRequest(request.id, true)">
                                                        接受
                                                    </n-button>
                                                    <n-button type="error" size="small"
                                                        @click="handleRequest(request.id, false)">
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
                    <n-list-item v-for="(value, index) in friendList">
                        {{ value }}
                    </n-list-item>
                </n-list>
            </template>
            <template #2>
                <!-- 聊天页面 -->
                <div class="p-4">
                    <h2 class="text-2xl font-bold mb-4">聊天页面</h2>
                    <p class="text-gray-500">这是与好友的聊天记录。</p>
                </div>
            </template>
        </n-split>
    </n-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
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
    useMessage,
} from 'naive-ui'
const showModal = ref(false)
const split = ref(0.2)
import { userApi } from '@/services/user'
import type { User as UserType, FriendShip } from '@/types/user'
import { formatRelativeTime } from '@/utils/format'
import UserSearch from '@/components/UserSearch.vue'

const router = useRouter()
const message = useMessage()

// State
const activeTab = ref('friends')

const friendsLoading = ref(false)
const requestsLoading = ref(false)

const friendList = ref<UserType[]>([])
const requestList = ref<(FriendShip & {
    username: string
    nickname: string
    avatar: string
    rating: number
    verification?: string
})[]>([])

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


const viewUserProfile = (userId: string) => {
    router.push({ name: 'UserHomePage', params: { id: userId } })
}


const startChat = () => {
    // TODO: Implement chat functionality
    message.info('聊天功能开发中')
}

onMounted(() => {
    loadFriendList()
    loadFriendRequests()
})
</script>
