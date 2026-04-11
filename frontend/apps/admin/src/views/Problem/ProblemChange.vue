<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia'
import { message } from 'ant-design-vue'
const problem_id = Number(useRoute().query.id);
const store = useUserStore()
const { id } = storeToRefs(store)
const form = ref<Problem>({
  id: 0,
  user_id: id.value,
  accept: 0,
  submission: 0,
  title: '',
  context: '',
  difficulty: 1,
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
})
onMounted(async () => {
  await problemApi.getProblemDetail(String(problem_id)).then((res) => {
    const { code, info } = res
    if (code == 200 && info) {
      const { problem } = info
      form.value = problem
    } else {
      message.error('获取题目信息失败')
    }
  }).catch((e) => {
    message.error('获取题目信息失败')
    console.log(e)
  })
})
import { ElMessage } from 'element-plus'
import { useRoute } from 'vue-router'
import ProblemEdit from './components/ProblemEdit.vue'
import { problemApi } from '@nexusoj/server'
import type { Problem } from '@nexusoj/type'
import { pick } from 'lodash'
const onSubmit = async () => {
  const obj = pick(form.value, ['title', 'context', 'judge_case', 'judge_config', 'judge_sample', 'tags', 'input_description', 'output_description', 'tips'])
  await problemApi.updateProblem(problem_id, {
    ...obj,
    difficulty: +form.value.difficulty,
    user_id: id.value,
  })
    .then((res) => {
      if (res.code == 200) {
        ElMessage({ message: '题目更改成功', type: 'success' })
      } else {
        ElMessage({ message: '题目更改失败', type: 'error' })
      }
    })
    .catch((err) => {
      ElMessage({ message: '题目更改失败', type: 'error' })
      console.log(err)
    })
}


</script>

<template>
  <ProblemEdit :form="form" :onSubmit="onSubmit" label="更新题目" />
</template>

<style scoped></style>
