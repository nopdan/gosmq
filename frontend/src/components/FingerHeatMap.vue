<script setup lang="ts">
import { Data } from "./Data";

const props = defineProps<{
  data: Data;
  name: string;
  num: number;
}>();
</script>

<template>
  <div class="heat-key-container">
    <div
      class="heat-key"
      :style="{
        'background-color': `rgb(80, 80, 255, ${(data.Dist.Finger[num] || 0) / data.Keys.Count / 0.3})`,
      }"
    >
      <div class="key-label">{{ name }}</div>
      <div class="key-rate">
        {{ ((data.Dist.Finger[num] / data.Keys.Count) * 100 || 0).toFixed(2) }}
      </div>
    </div>
  </div>
</template>
<style scoped></style>
