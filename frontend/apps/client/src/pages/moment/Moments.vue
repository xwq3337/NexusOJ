<script setup lang="ts">
import { ref, computed, onBeforeUnmount } from 'vue'
import { useClipboard } from '@vueuse/core'
import {
  NInput, NButton, NModal, NTab, NTabs, NTag, NAvatar, NTooltip,
  useMessage,
} from 'naive-ui'
import {
  Heart, MessageCircle, Share2, Bookmark, Plus, Image, Video, Music,
  Send,
} from 'lucide-vue-next'
// @ts-expect-error plyr uses CJS-style export
import Plyr from 'plyr'
import 'plyr/dist/plyr.css'
import { MOCK_MOMENTS, type Moment, type MomentMedia } from './mock'

const message = useMessage()

// ============ 数据 ============
const moments = ref<Moment[]>([...MOCK_MOMENTS])
const activeTab = ref('recommend')
const showPublish = ref(false)
const publishContent = ref('')
const publishMediaType = ref<'none' | 'image' | 'video' | 'audio'>('none')
const expandedComments = ref<Set<number>>(new Set())
const showShareModal = ref(false)
const shareMomentId = ref<number>(0)

const filteredMoments = computed(() => {
  if (activeTab.value === 'all') return moments.value
  return moments.value.filter(m => m.type === activeTab.value || m.type === 'recommend')
})

// ============ Plyr 指令 ============
const plyrInstances = new Set<Plyr>()

const vPlyr = {
  mounted(el: HTMLVideoElement | HTMLAudioElement) {
    const instance = new Plyr(el, {
      controls: ['play', 'progress', 'current-time', 'mute', 'volume'],
      hideControls: false,
      clickToPlay: true,
    })
    plyrInstances.add(instance)
  },
  beforeUnmount(_el: HTMLVideoElement | HTMLAudioElement) {
    // plyr will be destroyed in onBeforeUnmount
  },
}

onBeforeUnmount(() => {
  plyrInstances.forEach(p => p.destroy())
  plyrInstances.clear()
})

// ============ 交互 ============
const toggleLike = (moment: Moment) => {
  moment.isLiked = !moment.isLiked
  moment.likes += moment.isLiked ? 1 : -1
}

const toggleBookmark = (moment: Moment) => {
  moment.isBookmarked = !moment.isBookmarked
  moment.bookmarks += moment.isBookmarked ? 1 : -1
}

const toggleComments = (id: number) => {
  if (expandedComments.value.has(id)) {
    expandedComments.value.delete(id)
  } else {
    expandedComments.value.add(id)
  }
}

const handleShare = (id: number) => {
  shareMomentId.value = id
  showShareModal.value = true
}

const { copy } = useClipboard()
const copyShareLink = () => {
  copy(`${window.location.origin}/moments/${shareMomentId.value}`)
  message.success('链接已复制到剪贴板')
  showShareModal.value = false
}

// ============ 发布 ============
const handlePublish = () => {
  if (!publishContent.value.trim()) {
    message.warning('请输入内容')
    return
  }

  const mediaList: MomentMedia[] = []
  if (publishMediaType.value === 'video') {
    mediaList.push({ type: 'video', url: 'https://www.w3schools.com/html/mov_bbb.mp4', duration: '0:10' })
  } else if (publishMediaType.value === 'audio') {
    mediaList.push({ type: 'audio', url: 'https://www.w3schools.com/html/horse.mp3', duration: '0:05' })
  } else if (publishMediaType.value === 'image') {
    mediaList.push({ type: 'image', url: `https://picsum.photos/seed/new${Date.now()}/800/600` })
  }

  moments.value.unshift({
    id: Date.now(),
    user: { id: 0, username: 'me', nickname: '我', avatar: 'https://picsum.photos/seed/me/100/100' },
    content: publishContent.value,
    tags: [],
    media: mediaList,
    likes: 0,
    comments: 0,
    shares: 0,
    bookmarks: 0,
    isLiked: false,
    isBookmarked: false,
    createdAt: '刚刚',
    commentList: [],
    type: activeTab.value as Moment['type'],
  })

  publishContent.value = ''
  publishMediaType.value = 'none'
  showPublish.value = false
  message.success('发布成功！')
}

// ============ 工具函数 ============
const formatCount = (n: number): string => {
  if (n >= 10000) return (n / 10000).toFixed(1) + 'w'
  if (n >= 1000) return (n / 1000).toFixed(1) + 'k'
  return String(n)
}

const getImageGridClass = (count: number): string => {
  if (count === 1) return 'grid-cols-1'
  if (count === 2) return 'grid-cols-2'
  if (count === 4) return 'grid-cols-2'
  return 'grid-cols-3'
}

const getImageAspectClass = (count: number, _index: number): string => {
  if (count === 1) return 'aspect-video'
  if (count === 2) return 'aspect-[4/3]'
  return 'aspect-square'
}
</script>

<template>
  <div class="max-w-2xl mx-auto animate-fade-in">
    <!-- 顶部 Tab -->
    <div class="sticky top-0 z-10 py-4 mb-4" :style="{ backgroundColor: 'var(--bg-primary)' }">
      <n-tabs v-model:value="activeTab" type="line" animated>
        <n-tab name="recommend">推荐</n-tab>
        <n-tab name="following">关注</n-tab>
        <n-tab name="latest">最新</n-tab>
      </n-tabs>
    </div>

    <!-- 动态流 -->
    <div class="space-y-4 pb-20">
      <div
        v-for="moment in filteredMoments"
        :key="moment.id"
        class="rounded-xl p-5 transition-all duration-300 hover:shadow-lg"
        :style="{
          backgroundColor: 'var(--bg-card)',
          border: '1px solid var(--border-color)',
        }"
      >
        <!-- 用户信息 -->
        <div class="flex items-center gap-3 mb-3">
          <n-avatar :src="moment.user.avatar" :size="42" round />
          <div class="flex-1 min-w-0">
            <div class="font-medium text-sm" :style="{ color: 'var(--text-primary)' }">
              {{ moment.user.nickname }}
            </div>
            <div class="text-xs" :style="{ color: 'var(--text-tertiary)' }">
              @{{ moment.user.username }} · {{ moment.createdAt }}
            </div>
          </div>
        </div>

        <!-- 内容 -->
        <p class="text-sm leading-relaxed mb-3 whitespace-pre-wrap" :style="{ color: 'var(--text-primary)' }">
          {{ moment.content }}
        </p>

        <!-- 标签 -->
        <div v-if="moment.tags.length" class="flex flex-wrap gap-1.5 mb-3">
          <n-tag v-for="tag in moment.tags" :key="tag" size="small" :bordered="false" type="info">
            #{{ tag }}
          </n-tag>
        </div>

        <!-- 视频媒体 -->
        <div v-if="moment.media.length === 1 && moment.media[0].type === 'video'" class="mb-3 rounded-lg overflow-hidden">
          <video
            v-plyr
            :poster="moment.media[0].thumbnail"
            playsinline
            class="w-full"
          >
            <source :src="moment.media[0].url" type="video/mp4" />
          </video>
        </div>

        <!-- 音频媒体 -->
        <div v-if="moment.media.length === 1 && moment.media[0].type === 'audio'" class="mb-3">
          <div
            class="flex items-center gap-3 rounded-xl p-4"
            :style="{ backgroundColor: 'var(--surface-primary)' }"
          >
            <div
              class="w-12 h-12 rounded-full flex items-center justify-center shrink-0"
              :style="{ backgroundColor: 'var(--primary)', color: '#fff' }"
            >
              <Music :size="20" />
            </div>
            <div class="flex-1 min-w-0">
              <div class="text-xs mb-1" :style="{ color: 'var(--text-secondary)' }">
                🎵 音频 · {{ moment.media[0].duration }}
              </div>
              <audio v-plyr class="w-full">
                <source :src="moment.media[0].url" type="audio/mpeg" />
              </audio>
            </div>
          </div>
        </div>

        <!-- 图片媒体 -->
        <div
          v-if="moment.media.length > 0 && moment.media[0].type === 'image'"
          class="mb-3 grid gap-1.5 rounded-lg overflow-hidden"
          :class="getImageGridClass(moment.media.length)"
        >
          <div
            v-for="(media, idx) in moment.media"
            :key="idx"
            :class="getImageAspectClass(moment.media.length, idx)"
            class="overflow-hidden"
          >
            <img
              :src="media.url"
              :alt="`图片 ${idx + 1}`"
              class="w-full h-full object-cover transition-transform duration-300 hover:scale-105 cursor-pointer"
              loading="lazy"
            />
          </div>
        </div>

        <!-- 互动栏 -->
        <div class="flex items-center gap-1 pt-1">
          <n-tooltip trigger="hover">
            <template #trigger>
              <button
                class="flex items-center gap-1.5 px-3 py-1.5 rounded-full text-sm transition-all"
                :class="moment.isLiked ? 'text-red-500' : ''"
                :style="!moment.isLiked ? { color: 'var(--text-secondary)' } : {}"
                @click="toggleLike(moment)"
              >
                <Heart :size="18" :fill="moment.isLiked ? 'currentColor' : 'none'" />
                <span>{{ formatCount(moment.likes) }}</span>
              </button>
            </template>
            {{ moment.isLiked ? '取消点赞' : '点赞' }}
          </n-tooltip>

          <n-tooltip trigger="hover">
            <template #trigger>
              <button
                class="flex items-center gap-1.5 px-3 py-1.5 rounded-full text-sm transition-all"
                :style="{ color: 'var(--text-secondary)' }"
                @click="toggleComments(moment.id)"
              >
                <MessageCircle :size="18" />
                <span>{{ formatCount(moment.comments) }}</span>
              </button>
            </template>
            评论
          </n-tooltip>

          <n-tooltip trigger="hover">
            <template #trigger>
              <button
                class="flex items-center gap-1.5 px-3 py-1.5 rounded-full text-sm transition-all"
                :style="{ color: 'var(--text-secondary)' }"
                @click="handleShare(moment.id)"
              >
                <Share2 :size="18" />
                <span>{{ formatCount(moment.shares) }}</span>
              </button>
            </template>
            分享
          </n-tooltip>

          <div class="flex-1" />

          <n-tooltip trigger="hover">
            <template #trigger>
              <button
                class="flex items-center gap-1.5 px-3 py-1.5 rounded-full text-sm transition-all"
                :class="moment.isBookmarked ? 'text-amber-500' : ''"
                :style="!moment.isBookmarked ? { color: 'var(--text-secondary)' } : {}"
                @click="toggleBookmark(moment)"
              >
                <Bookmark :size="18" :fill="moment.isBookmarked ? 'currentColor' : 'none'" />
              </button>
            </template>
            {{ moment.isBookmarked ? '取消收藏' : '收藏' }}
          </n-tooltip>
        </div>

        <!-- 评论区 -->
        <div v-if="expandedComments.has(moment.id)" class="mt-4 pt-4" :style="{ borderTop: '1px solid var(--border-color)' }">
          <div v-if="moment.commentList.length" class="space-y-3">
            <div
              v-for="comment in moment.commentList"
              :key="comment.id"
              class="flex gap-2.5"
            >
              <n-avatar :src="comment.user.avatar" :size="28" round />
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-0.5">
                  <span class="text-xs font-medium" :style="{ color: 'var(--text-primary)' }">
                    {{ comment.user.nickname }}
                  </span>
                  <span class="text-xs" :style="{ color: 'var(--text-tertiary)' }">
                    {{ comment.createdAt }}
                  </span>
                </div>
                <p class="text-sm" :style="{ color: 'var(--text-secondary)' }">
                  {{ comment.content }}
                </p>
                <div class="flex items-center gap-3 mt-1">
                  <button class="flex items-center gap-1 text-xs" :style="{ color: 'var(--text-tertiary)' }">
                    <Heart :size="12" /> {{ comment.likes }}
                  </button>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="text-center py-4 text-sm" :style="{ color: 'var(--text-tertiary)' }">
            暂无评论，来说点什么吧～
          </div>

          <!-- 评论输入 -->
          <div class="flex items-center gap-2 mt-3">
            <n-input placeholder="写下你的评论..." size="small" round />
            <n-button size="small" circle type="primary">
              <Send :size="14" />
            </n-button>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div
        v-if="filteredMoments.length === 0"
        class="text-center py-20"
        :style="{ color: 'var(--text-tertiary)' }"
      >
        <p class="text-lg mb-2">暂无动态</p>
        <p class="text-sm">换个分类看看，或者发布第一条动态吧～</p>
      </div>
    </div>

    <!-- 发布按钮 (FAB) -->
    <button
      class="fixed bottom-8 right-8 w-14 h-14 rounded-full flex items-center justify-center shadow-2xl transition-all duration-300 hover:scale-110 z-20"
      :style="{ backgroundColor: 'var(--primary)', color: '#fff' }"
      @click="showPublish = true"
    >
      <Plus :size="24" />
    </button>

    <!-- 发布弹窗 -->
    <n-modal v-model:show="showPublish" preset="card" title="发布动态" style="max-width: 520px" :bordered="false">
      <n-input
        v-model:value="publishContent"
        type="textarea"
        placeholder="分享你的想法、解题思路、学习心得..."
        :autosize="{ minRows: 4, maxRows: 8 }"
      />

      <div class="flex items-center gap-2 mt-4">
        <span class="text-sm" :style="{ color: 'var(--text-secondary)' }">添加媒体：</span>
        <n-button
          size="small"
          :type="publishMediaType === 'image' ? 'primary' : 'default'"
          @click="publishMediaType = publishMediaType === 'image' ? 'none' : 'image'"
        >
          <Image :size="14" class="mr-1" /> 图片
        </n-button>
        <n-button
          size="small"
          :type="publishMediaType === 'video' ? 'primary' : 'default'"
          @click="publishMediaType = publishMediaType === 'video' ? 'none' : 'video'"
        >
          <Video :size="14" class="mr-1" /> 视频
        </n-button>
        <n-button
          size="small"
          :type="publishMediaType === 'audio' ? 'primary' : 'default'"
          @click="publishMediaType = publishMediaType === 'audio' ? 'none' : 'audio'"
        >
          <Music :size="14" class="mr-1" /> 音频
        </n-button>
      </div>

      <!-- 媒体预览 -->
      <div v-if="publishMediaType !== 'none'" class="mt-3 p-3 rounded-lg text-center" :style="{ backgroundColor: 'var(--surface-primary)' }">
        <Image v-if="publishMediaType === 'image'" :size="32" class="mx-auto mb-1" :style="{ color: 'var(--text-tertiary)' }" />
        <Video v-else-if="publishMediaType === 'video'" :size="32" class="mx-auto mb-1" :style="{ color: 'var(--text-tertiary)' }" />
        <Music v-else :size="32" class="mx-auto mb-1" :style="{ color: 'var(--text-tertiary)' }" />
        <p class="text-xs" :style="{ color: 'var(--text-tertiary)' }">
          已选择{{ publishMediaType === 'image' ? '图片' : publishMediaType === 'video' ? '视频' : '音频' }}（Mock 模式，自动填充示例媒体）
        </p>
      </div>

      <template #footer>
        <div class="flex justify-end gap-2">
          <n-button @click="showPublish = false">取消</n-button>
          <n-button type="primary" @click="handlePublish">
            <Send :size="14" class="mr-1" /> 发布
          </n-button>
        </div>
      </template>
    </n-modal>

    <!-- 分享弹窗 -->
    <n-modal v-model:show="showShareModal" preset="card" title="分享" style="max-width: 400px" :bordered="false">
      <div class="grid grid-cols-4 gap-4 py-4">
        <button class="flex flex-col items-center gap-2">
          <div class="w-12 h-12 rounded-full bg-green-500 text-white flex items-center justify-center">
            <Share2 :size="20" />
          </div>
          <span class="text-xs" :style="{ color: 'var(--text-secondary)' }">微信</span>
        </button>
        <button class="flex flex-col items-center gap-2">
          <div class="w-12 h-12 rounded-full bg-blue-500 text-white flex items-center justify-center">
            <Share2 :size="20" />
          </div>
          <span class="text-xs" :style="{ color: 'var(--text-secondary)' }">QQ</span>
        </button>
        <button class="flex flex-col items-center gap-2">
          <div class="w-12 h-12 rounded-full bg-red-500 text-white flex items-center justify-center">
            <Share2 :size="20" />
          </div>
          <span class="text-xs" :style="{ color: 'var(--text-secondary)' }">微博</span>
        </button>
        <button class="flex flex-col items-center gap-2" @click="copyShareLink">
          <div class="w-12 h-12 rounded-full flex items-center justify-center" :style="{ backgroundColor: 'var(--surface-secondary)' }">
            <Share2 :size="20" />
          </div>
          <span class="text-xs" :style="{ color: 'var(--text-secondary)' }">复制链接</span>
        </button>
      </div>
    </n-modal>
  </div>
</template>

<style scoped>
:deep(.plyr) {
  border-radius: 0.5rem;
}

:deep(.plyr--video) {
  aspect-ratio: 16 / 9;
}

:deep(.plyr--audio) {
  border-radius: 0.75rem;
}

:deep(.plyr__control--overlaid) {
  background: var(--primary);
}

:deep(.plyr--full-ui input[type="range"]) {
  color: var(--primary);
}
</style>
