<template>
    <div class="mb-4 gap-2 flex items-center">
        <n-input v-model:value="userSearchKeyword" placeholder="输入用户名或昵称搜索" clearable
            @keydown.enter="handleUserSearch" />
        <n-button type="info" @click="handleUserSearch">
            <template #icon>
                <Search :size="16" />
            </template>
        </n-button>
    </div>

    <n-spin :show="searchLoading">
        <n-empty v-if="userSearchResults.length === 0" description="请输入关键词搜索用户" />

        <div v-else class="space-y-3">
            <div v-for="user in userSearchResults" :key="user.id" class="flex items-center gap-4 p-4 rounded-lg border"
                :style="{ backgroundColor: 'var(--card-bg)' }">
                <n-avatar :size="48" :src="user.avatar || undefined">
                    <User v-if="!user.avatar" />
                </n-avatar>

                <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2 mb-1 cursor-pointer"
                        @click="$router.push({ name: 'UserHomePage', params: { id: user.id } })">
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
    </n-spin>

</template>


<script setup lang="ts">
import { ref } from 'vue'
import { Search, User } from 'lucide-vue-next'
import { NAvatar,NInput,NModal, NButton, NEmpty, NSpin } from 'naive-ui'
import { useMessage } from 'naive-ui'
import { useUserStore } from '@/stores/useUserStore'
import { userApi } from '@/services/user'
import type { User as UserType } from '@/types/user'
const message = useMessage()
const searchKeyword = ref('')
const userSearchKeyword = ref('')
const searchLoading = ref(false)
const userSearchResults = ref<UserType[]>([])

const verificationMessage = ref('')
const showAddFriend = ref(false)
const sendingRequest = ref(false)
const selectedUser = ref<UserType | null>(null)

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

</script>