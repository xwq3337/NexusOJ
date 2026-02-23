<template>
  <div class="animate-fade-in">
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
      <!-- Left Sidebar - User Profile -->
      <div class="lg:col-span-1">
        <UserSideBar :user-id="userId" />
      </div>

      <!-- Main Content -->
      <div class="lg:col-span-3">
        <!-- User Stats Cards -->
        <UserStatsCard :user-id="userId" />

        <!-- User Heatmap -->
        <div class="mb-6">
          <UserHeatmap :user-id="userId" />
        </div>

        <!-- Tabs for different sections -->
        <div
          class="rounded-xl"
          :style="{
            backgroundColor: 'var(--surface-primary)',
            borderColor: 'var(--border-color)',
            borderWidth: '1px',
            borderStyle: 'solid',
            paddingLeft : '10px'
          }"
        >
          <n-tabs
            v-model:value="activeTab"
            type="line"
            :tab-style="{
              color: 'var(--text-secondary)',
              fontWeight: '500'
            }"
          >
            <n-tab-pane name="submissions" tab="提交记录">
              <div class="p-6">
                <p class="text-center py-8" :style="{ color: 'var(--text-tertiary)' }">
                  提交记录组件待实现...
                </p>
              </div>
            </n-tab-pane>

            <n-tab-pane name="solutions" tab="题解">
              <div class="p-6">
                <p class="text-center py-8" :style="{ color: 'var(--text-tertiary)' }">
                  题解组件待实现...
                </p>
              </div>
            </n-tab-pane>

            <n-tab-pane name="blogs" tab="博客">
              <div class="p-6">
                <p class="text-center py-8" :style="{ color: 'var(--text-tertiary)' }">
                  博客组件待实现...
                </p>
              </div>
            </n-tab-pane>

            <n-tab-pane name="discussions" tab="讨论">
              <div class="p-6">
                <p class="text-center py-8" :style="{ color: 'var(--text-tertiary)' }">
                  讨论组件待实现...
                </p>
              </div>
            </n-tab-pane>

            <n-tab-pane name="gallery" tab="展示">
              <div class="p-6">
                <p class="text-center py-8" :style="{ color: 'var(--text-tertiary)' }">
                  展示组件待实现...
                </p>
              </div>
            </n-tab-pane>
          </n-tabs>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { NTabs, NTabPane } from 'naive-ui'
import UserSideBar from '@/pages/user/components/UserSideBar.vue'
import UserStatsCard from '@/pages/user/components/UserStatsCard.vue'
import UserHeatmap from '@/pages/user/components/UserHeatmap.vue'

const route = useRoute()
const activeTab = ref('submissions')

// 获取用户 ID
const userId = ref<string>('')
userId.value = route.params.id as string

</script>
