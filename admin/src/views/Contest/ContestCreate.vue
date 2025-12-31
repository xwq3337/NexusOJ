<script setup lang="ts">
import { ref } from 'vue';
import type { Dayjs } from 'dayjs';
type RangeValue = [Dayjs, Dayjs];
const TimeRange = ref<RangeValue>();

</script>

<template>
    <a-form>
    </a-form>
    <a-range-picker v-model:value="TimeRange" show-time />
</template>

<style scoped>

</style>
