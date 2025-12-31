<script setup lang="ts">
import { onMounted } from "vue";
import { MenuList } from "../views";
const emit = defineEmits(["update"]);
onMounted(() => {
  handleClick("./views/Home.vue");
});
const handleClick = (path: string) => {
  emit("update", path);
};
</script>

<template>
  <el-menu default-active="0">
    <el-menu-item
      v-for="(item, index) in MenuList"
      :key="index"
      @click="handleClick(item.path)"
    >
      <span>{{ item.label }}</span>
    </el-menu-item>
  </el-menu>
</template>
