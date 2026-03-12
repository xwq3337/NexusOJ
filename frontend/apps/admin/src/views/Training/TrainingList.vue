<script setup lang="ts">
import { ref } from 'vue';
const tableData = ref([]);

</script>
<template>
    <el-table :data="tableData" style="width: 100%">
      <el-table-column prop="id" label="ID" width="100"></el-table-column>
      <el-table-column prop="title" label="名称" width="180"></el-table-column>

    </el-table>
</template>
<style scoped lang="scss"></style>
