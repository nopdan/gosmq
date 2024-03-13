<script setup lang="ts">
import { Data } from "./Data";

const props = defineProps<{
  _key: string;
  data: Data;
}>();
</script>

<template>
  <div
    class="heat-key"
    :style="{
      'background-color': `rgb(255, 0, 0, ${(data.Dist.Key[_key] || 0) / data.Keys.Count / 0.13})`,
    }"
  >
    <div style="font-size: 0.8em; color: #000; white-space: nowrap">{{ _key.toUpperCase() }}</div>
    <div style="font-size: 0.3em; color: #555; white-space: nowrap">
      {{ ((data.Dist.Key[_key] / data.Keys.Count) * 100 || 0).toFixed(2) }}
    </div>
  </div>
</template>
<style>
div.heat-key {
  display: flex;
  flex-direction: column;
  align-items: center;
  border-radius: 5px;
  border: 1px solid #e0e0e0;
  margin: 2px;
}
</style>
