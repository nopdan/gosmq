<script setup lang="ts">
import AutoSortLine from "./AutoSortLine.vue";

const sortByAnother = ref(false);
const schema1: any = inject("schema1");
const schema2: any = inject("schema2");

const data = computed(() => {
  let result = [];
  for (let k of Object.keys(schema1.value.Keys)) {
    if (k in schema2.value.Keys) {
      result.push([k, schema1.value.Keys[k].Count, schema2.value.Keys[k].Count]);
    }
  }
  if (sortByAnother.value) {
    result = result.sort((a: any, b: any) => b[2] - a[2]);
  } else {
    result = result.sort((a: any, b: any) => b[1] - a[1]);
  }
  return result;
});

const d1 = computed(() => data.value.map((i) => i[1]));
const d2 = computed(() => data.value.map((i) => i[2]));
const n1 = computed(() => schema1.value.Name);
const n2 = computed(() => schema2.value.Name);
const labalName = computed(() => data.value.map((i) => i[0]));
</script>
<template>
  两方案的各按键使用次数，可以按照方案２排序：<n-switch v-model:value="sortByAnother" size="small" />
  <auto-sort-line :names="labalName" :data1="d1" :data2="d2" :schemaName1="n1" :schemaName2="n2" />
</template>
