<script setup lang="ts">
import { defineAsyncComponent, ref } from "vue";
import Menu from "./components/Menu.vue";
const activeComponent = ref();
const handleUpdateComponent = (path: any) => {
  activeComponent.value = () => import(path);
};
</script>

<template>
  <el-container>
    <el-aside width="200px" style="height: 100vh">
      <Menu @update="handleUpdateComponent"/>
    </el-aside>
    <el-main>
      <component v-if="activeComponent" :is="defineAsyncComponent(activeComponent)" />
    </el-main>
  </el-container>
</template>

<style scoped></style>
