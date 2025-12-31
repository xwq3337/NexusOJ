<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ref } from 'vue'
const router = useRouter()
const activeIndex = ref('')
const handleSelect = (index: string) => {
  activeIndex.value = index
  router.push({ name: index })
}

</script>
<template>
  <el-menu :default-active="activeIndex" class="el-menu-demo" mode="horizontal" @select="handleSelect">
    <el-menu-item index="blog-info">博客管理</el-menu-item>
    <el-menu-item index="blog-recycle">回收站</el-menu-item>
    <el-menu-item index="blog-verify">博客审核</el-menu-item>
    <el-sub-menu index="3">
      <template #title>Workspaces</template>
      <el-menu-item index="3-1">item one</el-menu-item>
      <el-sub-menu index="3-2">
        <template #title>item two</template>
        <el-menu-item index="3-2-1">item one</el-menu-item>
      </el-sub-menu>
    </el-sub-menu>

  </el-menu>
  <router-view />
</template>
