<template>
  <n-descriptions bordered label-placement="top" size="small">
    <n-descriptions-item>
      <template #label>同指击键</template>
      {{ data.Fingers.Same.Count }} <br />
      {{ ("" + data.Fingers.Same.Rate * 100).substring(0, 4) + "%" }}
    </n-descriptions-item>
    <n-descriptions-item>
      <template #label>异指击键</template>
      {{ data.Fingers.Diff.Count }} <br />
      {{ ("" + data.Fingers.Diff.Rate * 100).substring(0, 4) + "%" }}
    </n-descriptions-item>
  </n-descriptions>
</template>
<script setup lang="ts">
import { NDescriptions, NDescriptionsItem } from "naive-ui";
const p = defineProps(["data"]);
</script>

<style>
.n-descriptions-table {
  text-align: left;
}
</style>
