<script setup lang="ts">
import { userApi } from '@/services/user';
import { User } from '@/types/user';
import { NDropdown, NAvatar, NText, useMessage, NButton } from 'naive-ui'
import { h, onMounted, ref } from 'vue'
const props = defineProps<{
    user_id: string
}>()
const user = ref<User>()
onMounted(async () => {
    await userApi.getInfoById(props.user_id).then(response => {
        user.value = response.info
    }).catch(error => {
        console.error('获取用户信息失败:', error)
    })
})
function renderHeader() {
    return h(
        'div',
        {
            style: 'display: flex; align-items: center; padding: 8px 12px;'
        },
        [
            h(NAvatar, {
                round: true,
                style: 'margin-right: 12px;',
                src: user.value?.avatar || 'https://07akioni.oss-cn-beijing.aliyuncs.com/demo1.JPG'
            }),
            h('div', null, [
                h('div', null, [h(NText, { depth: 2 }, { default: () => user.value?.username || '未知用户' })]),
                h('div', { style: 'font-size: 12px;' }, [
                    h(
                        NText,
                        { depth: 3 },
                        { default: () => '毫无疑问，你是最亮的星' }
                    )
                ])
            ]),

        ]
    )
}

const message = useMessage()
const options = [
    {
        key: 'header',
        type: 'render',
        render: renderHeader
    },
    {
        key: 'header-divider',
        type: 'divider'
    },
    {
        key: 'body',
        type: "render",
        render: () => h('div', { style: 'padding: 8px 12px;' }, h(
            'div', { style: 'padding: 8px 12px;display: flex; align-items: center; justify-content: space-between;' }, [
            h(NButton, { depth: 2, type: "info" }, { default: () => '关注' }),
            h(NButton, { depth: 2, type: "info", ghost: true }, { default: () => '私信' })
        ]
        ))
    },
]

function handleSelect(key: string | number) {
    message.info(String(key))
}
</script>
<template>
    <n-dropdown trigger="hover" :options="options" @select="handleSelect">
        <slot></slot>
    </n-dropdown>
</template>