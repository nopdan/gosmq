<script setup lang="ts">
import { Data } from "./Data";

const props = defineProps<{
  _key: string;
  data: Data;
}>();
</script>

<template>
  <div class="heat-key-container">
    <div
      class="heat-key"
      :style="{
        'background-color': `rgb(255, 0, 0, ${(data.Dist.Key[_key] || 0) / data.Keys.Count / 0.13})`,
      }"
    >
      <div class="key-label">{{ _key.toUpperCase() }}</div>
      <div class="key-rate">
        {{ ((data.Dist.Key[_key] / data.Keys.Count) * 100 || 0).toFixed(2) }}
      </div>
    </div>
  </div>
</template>
<style>
div.heat-key-container {
  --my-wh: calc(33rem / 12);
  width: var(--my-wh);
  height: var(--my-wh);
}

div.heat-key {
  display: flex;
  flex-direction: column;
  align-items: center;
  border-radius: 5px;
  border: 1px solid #e0e0e0;
  margin: 0.125rem;

  & > .key-label {
    font-size: 0.9rem;
    color: #000;
    white-space: nowrap;
  }
  & > .key-rate {
    font-size: 0.6rem;
  }
}
</style>
