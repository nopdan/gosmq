<template>
  <n-descriptions bordered label-placement="top" size="small" :column="2">
    <n-descriptions-item>
      <template #label>方案名称</template>
      {{ data.Name }}
    </n-descriptions-item>
    <n-descriptions-item>
      <template #label>码表词条数量</template>
      {{ data.Basic.DictLen }}
    </n-descriptions-item>
    <n-descriptions-item>
      <template #label>缺字总数</template>
      {{ data.Basic.LackCount }}
    </n-descriptions-item>
    <n-descriptions-item>
      <template #label>非汉字总数</template>
      {{ data.Basic.NotHanCount }}
    </n-descriptions-item>
    <n-descriptions-item :span="2">
      <template #label>缺字</template>
      {{ data.Basic.Lack }}&nbsp;
    </n-descriptions-item>
    <n-descriptions-item :span="2">
      <template #label>非汉字字符</template>
      {{ data.Basic.NotHan }}
    </n-descriptions-item>
    <n-descriptions-item>
      <template #label>上屏总数</template>
      {{ data.Basic.Commits }}
    </n-descriptions-item>
  </n-descriptions>
</template>
<script setup lang="ts">
const p = defineProps(["data"]);
</script>

<style>
.n-descriptions-table {
  text-align: left;
}
</style>
