<script setup lang="ts">
import { computed, ref } from 'vue'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia'
import MarkdownPreview from '@/components/MarkdownPreview.vue'
const store = useUserStore()
const { id ,username} = storeToRefs(store)
const emit = defineEmits(['submit'])
const form = defineModel<Problem>('form',{required: true})
import DashBoardJudgeCase from './DashBoardJudgeCase.vue'
import DashBoardJudgeSample from './DashBoardJudgeSample.vue'
import type { Problem } from '../type'
const { label } = defineProps(['label'])

const check = () => {
  if (
    form.value.title === '' ||
    form.value.context === '' ||
    form.value.input_description === '' ||
    form.value.output_description === ''
  )
    return false
  if (form.value.judge_case.length === 0) return false
  return true
}

const display = computed(() => {
  let res = ''
  form.value.judge_sample.forEach((item, index: number) => {
    res += `### 样例输入${index + 1}\n\n\`\`\`\n${item.input} \n\`\`\`\n`
    res += `### 样例输出${index + 1}\n\n\`\`\`\n${item.expected}\n \`\`\`\n`
    if (index != form.value.judge_sample.length - 1) {
      res += `--- \n\n`
    }
  });
  return `# ${form.value.title}\n\n` +
    `时空限制: ${form.value.judge_config.time_limit}s  ${form.value.judge_config.memory_limit}MB\n\n` +
    `## 题目描述\n${form.value.context} \n\n` +
    `## 输入格式\n${form.value.input_description}\n\n` +
    `## 输出格式\n${form.value.output_description}\n\n ` +
    `## 样例\n${res}\n\n` +
    `## 提示\n ${form.value.tip} \n\n` +
    `### 造题人\n${username.value} `
})
const JudgeCaseVisiable = ref(false)
const JudgeSampleVisiable = ref(false)
const resetForm = () => {
  form.value.title = ''
  form.value.context = ''
  form.value.input_description = ''
  form.value.output_description = ''
  form.value.tip = ''
  form.value.judge_case = [
    {
      input: '',
      expected: '',
    },
  ]
  form.value.judge_config = {
    time_limit: 64,
    memory_limit: 1,
  }
  form.value.judge_sample = [
    {
      input: '',
      expected: '',
    },
  ]
  form.value.tags = []
  form.value.user_id = id.value
  form.value.difficulty = '1'
}
</script>

<template>
  <el-splitter>
    <el-splitter-panel size="50%">
      <el-card>
        <template #header>
          <el-button type="primary" @click="emit('submit')" :disabled="!check()">{{ label }}</el-button>
          <el-button type="default" @click="resetForm">重置</el-button>
          <el-button type="default" @click="JudgeCaseVisiable = true">测试用例</el-button>
          <el-button type="default" @click="JudgeSampleVisiable = true">样例</el-button>
          <el-dropdown>
            <el-button type="default">时间/空间</el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-input-number v-model="form.judge_config.time_limit">
                  <template #prepend>时间限制</template>
                  <template #append>秒</template>
                </el-input-number><br />
                <el-input-number v-model="form.judge_config.memory_limit">
                  <template #prepend>空间限制</template>
                  <template #append>MB</template>
                </el-input-number>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <el-dropdown>
            <el-button type="default">标签</el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-input-tag v-model="form.tags" draggable clearable delimiter="," placeholder="请输入标签"
                  tag-type="primary" />
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
        <el-input v-model="form.title"  placeholder="题目名称" />
        <el-input v-model="form.context" autosize placeholder="题目描述" type="textarea" />
        <el-select v-model="form.difficulty" placeholder="难度">
          <el-option label="C" value="1" />
          <el-option label="B" value="2" />
          <el-option label="A" value="3" />
          <el-option label="S" value="4" />
          <el-option label="SS" value="5" />
        </el-select>
        <el-input v-model="form.input_description" autosize placeholder="输入描述" type="textarea" />
        <el-input v-model="form.output_description" autosize placeholder="输出描述" type="textarea" />
        <el-input v-model="form.tip" placeholder="提示" autosize type="textarea" />
      </el-card>
    </el-splitter-panel>
    <el-splitter-panel :min="200">
      <MarkdownPreview :value="display" />
    </el-splitter-panel>
  </el-splitter>

  <DashBoardJudgeCase v-model:open="JudgeCaseVisiable" v-model:judge_case="form.judge_case" />
  <DashBoardJudgeSample v-model:open="JudgeSampleVisiable" v-model:judge_sample="form.judge_sample" />
</template>

<style scoped></style>
