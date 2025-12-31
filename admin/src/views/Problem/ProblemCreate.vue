<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import Request from '@/api/axios'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia'
import { HttpStatusCode } from 'axios'
const store = useUserStore()
const { id } = storeToRefs(store)
const form = ref<Problem>({
  title: '',
  context: '',
  difficulty: '0',
  judge_case: [
    {
      input: '',
      expected: '',
    },
  ],
  judge_config: {
    time_limit: 64,
    memory_limit: 1,
  },
  judge_sample: [
    {
      input: '',
      expected: '',
    },
  ],
  tags: [],
  input_description: '',
  output_description: '',
  tip: '',
  user_id: id.value,
})
import { useCache } from '@/stores/Cache'
import { ElMessage } from 'element-plus'
import ProblemEdit from './components/ProblemEdit.vue'
import type { Problem } from './type'
const cache = useCache()
const { addCache, getCache } = cache

watch(() => form.value, () => {
  addCache('problemCreateForm', form.value)
}, { deep: true })
onMounted(() => { form.value = getCache('problemCreateForm') || form.value })

const onSubmit = async () => {
  await Request.post({
    url: `/problem/create`,
    data: {
      title: form.value.title,
      context: form.value.context,
      difficulty: +form.value.difficulty,
      judge_case: form.value.judge_case,
      judge_config: form.value.judge_config,
      judge_sample: form.value.judge_sample,
      tags: form.value.tags,
      input_description: form.value.input_description,
      output_description: form.value.output_description,
      user_id: id.value,
      tip: form.value.tip,
    },
  })
    .then((response) => {
      if (response.code == HttpStatusCode.Ok) {
        ElMessage({ message: '题目创建成功', type: 'success' })
      } else {
        ElMessage({ message: '题目创建失败', type: 'error' })
      }
    })
    .catch((err) => {
      ElMessage({ message: '题目创建失败', type: 'error' })
      console.log(err)
    })
}

</script>

<template>
  <ProblemEdit label="创建题目" :form="form" @submit="onSubmit" />
</template>

<style scoped></style>
