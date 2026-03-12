<script setup lang="ts">

const judge_case = defineModel<any[]>('judge_case')
const open = defineModel('open')
</script>

<template>
  <el-dialog title="测试用例" v-model="open" width="400">

    <div v-for="(sample, index) in judge_case" :key="index" style="padding: 10px 15px">
      <h3>#{{ index + 1 }}</h3>
      <el-input type="textarea" autosize v-model="sample.input" />
      <el-input type="textarea" autosize v-model="sample.expected" />
    </div>
    <el-button type="primary" @click="judge_case?.push({ input: '', expected: '' })">新增用例</el-button>
    <el-button type="primary" @click="judge_case = judge_case?.slice(0, -1) ?? []">删除末尾</el-button>
  </el-dialog>
</template>
