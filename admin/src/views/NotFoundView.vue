<script setup lang="ts"></script>
<template>
  <a-result
    style="margin: 0 auto"
    status="404"
    title="404"
    sub-title="Sorry, the resource you requested cannot be found."
  >
    <template #extra>
      <a-button type="primary" @click="$router.push('/')">Back Home</a-button>
    </template>
  </a-result>
</template>

<style scoped></style>
