<script setup lang="ts">
import Request from '@/api/axios';
import { onMounted, ref } from 'vue';
 const tableData = ref([]);
onMounted(async ()=>{
  // 获取训练信息列表
  await Request.get({url:'/training/list'}).then((res) => {
    tableData.value = res.info;
    })
})
</script>
<template>
    <el-table :data="tableData" style="width: 100%">
      <el-table-column prop="id" label="ID" width="100"></el-table-column>
      <el-table-column prop="title" label="名称" width="180"></el-table-column>

    </el-table>
</template>
<style scoped lang="scss"></style>
