<script setup lang="ts">

const open = defineModel('open')
const judge_sample = defineModel<any[]>('judge_sample')
</script>

<template>
  <el-dialog title="样例" v-model="open" width="400">

    <div v-for="(sample, index) in judge_sample" :key="index" style="padding: 10px 15px">
      <h3>#{{ index + 1 }}</h3>
      <el-input type="textarea" autosize v-model="sample.input" />
      <el-input type="textarea" autosize v-model="sample.expected" />
    </div>
    <el-button type="primary" @click="judge_sample?.push({ input: '', expected: '' })">新增样例</el-button>
    <el-button type="primary" @click="judge_sample = judge_sample?.slice(0, -1) ?? []">删除末尾</el-button>
  </el-dialog>
</template>
