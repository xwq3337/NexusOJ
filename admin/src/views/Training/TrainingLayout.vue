<script setup lang="ts">
import { useRouter } from 'vue-router';
const router = useRouter();
const items = [
  {
    title: '信息列表',
    pathName: 'training-info'
  },{
    title: '创建训练',
    pathName: 'training-create'
  }
]
const handleSelect = (index: string) => {
  router.push({ name: index })
}
</script>

<template>
  <el-menu mode="horizontal" @select="handleSelect">
    <el-menu-item v-for="(i,index) in items" :key="index" :index="i.pathName">{{ i.title }}</el-menu-item>
  </el-menu>
  <RouterView style="margin-top: 20px;"/>
</template>
