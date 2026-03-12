<script setup lang="ts">

const value = defineProps(['value'])

</script>
<template>
  <v-md-preview :text="value.value" height="400px"></v-md-preview>
</template>
