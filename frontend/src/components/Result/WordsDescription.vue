<template>
  <n-descriptions bordered label-placement="top" size="small" :column="4">
    <n-descriptions-item>
      <template #label>打词上屏总数</template>
      {{ data.Words.Commits.Count }}
    </n-descriptions-item>
    <n-descriptions-item>
      <template #label>打词率</template>
      {{ ("" + data.Words.Commits.Rate * 100).substring(0, 4) + "%" }}
    </n-descriptions-item>
    <n-descriptions-item>
      <template #label>出字总数</template>
      {{ data.Words.Chars.Count }}
    </n-descriptions-item>
    <n-descriptions-item>
      <template #label>首选次数</template>
      {{ data.Words.FirstCount }}
    </n-descriptions-item>
  </n-descriptions>

  <n-descriptions bordered label-placement="top" size="small" :column="6">
    <n-descriptions-item v-for="(v, i) of data.Words.Dist">
      <template #label>{{ i + 1 }} 字词</template>
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
