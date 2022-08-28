<template>
  <n-descriptions bordered label-placement="top" size="small">
    <n-descriptions-item>
      <template #label>选重次数</template>
      {{ data.Collision.Commits.Count }}
    </n-descriptions-item>
    <n-descriptions-item>
      <template #label>选重率</template>
      {{ ("" + data.Collision.Commits.Rate * 100).substring(0,4) + "%" }}
    </n-descriptions-item>
    <n-descriptions-item >
      <template #label>选重出字总数</template>
      {{ data.Collision.Chars.Count }}
    </n-descriptions-item>
  </n-descriptions>
    <n-descriptions bordered label-placement="top" size="small" :column="9">
    <n-descriptions-item v-for="(v, i) of data.Collision.Dist">
      <template #label>{{ i+1 }} 重</template>
      {{ v }}
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
