<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia'
import { message } from 'ant-design-vue'
const problem_id = Number(useRoute().query.id);
const store = useUserStore()
const { id } = storeToRefs(store)
const form = ref<Problem>({
  id : 0,
  user_id: id.value,
  accept : 0,
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
    const {code ,info } = res
    if (code == 200 && info){
      form.value = info
    }else {
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

const onSubmit = async () => {
  await problemApi.updateProblem(String(problem_id), {
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
    tips: form.value.tips,
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
