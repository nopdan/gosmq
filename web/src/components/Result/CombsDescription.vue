<template>
  <n-descriptions bordered label-placement="top" size="small" :column="4">
    <n-descriptions-item>
      <template #label>手指组合总数</template>
      {{ data.Combs.Count }}
    </n-descriptions-item>
    <n-descriptions-item>
      <template #label>平均当量</template>
      {{ data.Combs.Equivalent }}
    </n-descriptions-item>
    <n-descriptions-item v-for="i of items">
      <template #label>{{i[0]}}</template>
      {{ data.Combs?.[i[1]].Count }}<br />(
      {{ ("" + data.Combs?.[i[1]].Rate * 100).substring(0,4) + "%" }})
    </n-descriptions-item>
  </n-descriptions>
</template>
<script setup lang="ts">
import { NDescriptions, NDescriptionsItem } from "naive-ui";
const p = defineProps(["data"]);

const items = [
  ["二连击", "DoubleHit"],
  ["三连击", "TribleHit"],
  ["小跨排", "SingleSpan"],
  ["大跨排", "MultiSpan"],
  ["错手", "LongFingersDisturb"],
  ["小指干扰", "LittleFingersDisturb"],
];
</script>

<style>
.n-descriptions-table {
  text-align: left;
}
</style>
