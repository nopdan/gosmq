<template>
  <n-descriptions bordered label-placement="top" size="small">
    <n-descriptions-item>
      <template #label>选重次数</template>
      {{ data.Collision.Commits.Count }}
    </n-descriptions-item>
    <n-descriptions-item>
      <template #label>选重率</template>
      {{ ("" + data.Collision.Commits.Rate * 100).substring(0, 4) + "%" }}
    </n-descriptions-item>
    <n-descriptions-item>
      <template #label>选重出字总数</template>
      {{ data.Collision.Chars.Count }}
    </n-descriptions-item>
  </n-descriptions>
  <n-descriptions bordered label-placement="top" size="small" :column="9">
    <n-descriptions-item v-for="v in dist">
      <template #label>{{ v.key }} 重</template>
      {{ v.value }}
    </n-descriptions-item>
  </n-descriptions>
</template>
<script setup lang="ts">
import { NDescriptions, NDescriptionsItem } from "naive-ui";
const p = defineProps(["data"]);

let dist = new Array();
for (let i = 0; i < p.data.Collision.Dist.length; i++) {
  if (p.data.Collision.Dist[i] === 0) continue;
  dist.push({
    key: i + 1,
    value: p.data.Collision.Dist[i],
  });
}
</script>

<style>
.n-descriptions-table {
  text-align: left;
}
</style>
