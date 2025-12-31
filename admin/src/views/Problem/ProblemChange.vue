<script setup lang="ts">
import { onMounted, ref } from 'vue'
import Request from '@/api/axios'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia'
import { message } from 'ant-design-vue'
import { HttpStatusCode } from 'axios'
const problem_id = Number(useRoute().query.id);
const store = useUserStore()
const { id } = storeToRefs(store)
const form = ref({
  title: '',
  context: '',
  difficulty: "1",
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
})
onMounted(async () => {
  await Request.get({ url: `/problem/${problem_id}` }).then((res) => {
    if (res.code == HttpStatusCode.Ok) {
      form.value = res.info
    } else {
      message.error('获取题目信息失败')
    }
  }).catch((err) => {
    message.error('获取题目信息失败')
    console.log(err)
  })
})
import { ElMessage } from 'element-plus'
import { useRoute } from 'vue-router'
import ProblemEdit from './components/ProblemEdit.vue'

const onSubmit = async () => {
  await Request.post({
    url: `/problem/change`,
    data: {
      id: problem_id,
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
