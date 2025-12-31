<script setup>

const value = ref('## Hello, VuePress!')
</script>


<template>
  <v-md-editor v-model="value" height="500px" />
</template>
