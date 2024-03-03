<template>
  <n-descriptions bordered label-placement="top" size="small" :column="4">
    <n-descriptions-item>
      <template #label>左手使用量</template>
      {{ data.Hands.Left.Count }} <br />
      {{ ("" + data.Hands.Left.Rate * 100).substring(0, 4) + "%" }}
    </n-descriptions-item>
    <n-descriptions-item>
      <template #label>右手使用量</template>
      {{ data.Hands.Right.Count }} <br />
      {{ ("" + data.Hands.Right.Rate * 100).substring(0, 4) + "%" }}
    </n-descriptions-item>
    <n-descriptions-item>
      <template #label>同手击键</template>
      {{ data.Hands.Same.Count }} <br />
      {{ ("" + data.Hands.Same.Rate * 100).substring(0, 4) + "%" }}
    </n-descriptions-item>
    <n-descriptions-item>
      <template #label>异手击键</template>
      {{ data.Hands.Diff.Count }} <br />
      {{ ("" + data.Hands.Diff.Rate * 100).substring(0, 4) + "%" }}
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
