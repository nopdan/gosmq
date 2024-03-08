<script setup lang="ts">
import ComparedBars from "./comparedBars.vue";
const schema1: any = inject("schema1");
const schema2: any = inject("schema2");
function getData(schema: any) {
  return schema.CodeLen.Dist;
}

const name1 = computed(() => {
  return schema1.value.Name;
});
const name2 = computed(() => {
  return schema2.value.Name;
});

const data1 = computed(() => {
  return schema1.value.CodeLen.Dist;
});

const data2 = computed(() => {
  return schema2.value.CodeLen.Dist;
});

const names = computed(() => {
  const l1 = (data1.value as any).length;
  const l2 = (data2.value as any).length;
  const l = l1 > l2 ? l1 : l2;
  let result = [];
  for (let i = 0; i < l; i++) {
    result.push("" + (i + 1) + " 码");
  }
  return result;
});
</script>
<template>
  <compared-bars :data1="data1" :data2="data2" :names="names" :schema-names="[name1, name2]" />
</template>
