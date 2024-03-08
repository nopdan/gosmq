<template>
  <n-descriptions bordered label-placement="top" size="small">
    <n-descriptions-item>
      <template #label>总码长</template>
      {{ data.CodeLen.Total }}
    </n-descriptions-item>
    <n-descriptions-item>
      <template #label>字均码长</template>
      {{ data.CodeLen.PerChar }}
    </n-descriptions-item>
  </n-descriptions>
  <n-descriptions bordered label-placement="top" size="small" :column="8">
    <n-descriptions-item v-for="(v, i) of data.CodeLen.Dist">
      <template #label>{{ i + 1 }} 码</template>
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
