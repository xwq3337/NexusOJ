<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia'
const store = useUserStore()
const { id } = storeToRefs(store)
const form = ref<Problem>({
  id: 0,
  accept: 0,
  submission: 0,
  title: '',
  context: '',
  difficulty: 0,
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
  tips: '',
  user_id: id.value,
})
import { useCache } from '@/stores/Cache'
import { ElMessage } from 'element-plus'
import ProblemEdit from './components/ProblemEdit.vue'
import type { Problem } from '@nexusoj/type'
import { problemApi } from '@nexusoj/server'
import { pick } from 'lodash'

const cache = useCache()
const { addCache, getCache } = cache

watch(() => form.value, () => {
  addCache('problemCreateForm', form.value)
}, { deep: true })
onMounted(() => { form.value = getCache('problemCreateForm') || form.value })
const onSubmit = async () => {
  const obj = pick(form.value, ['title', 'context', 'judge_case', 'difficulty', 'judge_config', 'judge_sample', 'tags', 'input_description', 'output_description', 'tips', 'user_id'])
  await problemApi.createProblem(obj)
    .then((res) => {
      const { code } = res
      if (code == 200) {
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
