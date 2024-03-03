<script setup lang="ts">
import ComparedBars from "./comparedBars.vue";
import { inject, computed } from "vue";
const schema1: any = inject("schema1");
const schema2: any = inject("schema2");
const s1 = computed(() => schema1.value.Name);
const s2 = computed(() => schema2.value.Name);
const d1 = computed(() => schema1.value.Collision.Dist);
const d2 = computed(() => schema2.value.Collision.Dist);

const labalName = computed(() => {
  const l1 = d1.value.length;
  const l2 = d1.value.length;
  const l = l1 > l2 ? l1 : l2;
  let result = [];
  for (let i = 0; i < l; i++) {
    result.push("" + (i + 1) + " 重");
  }
  return result;
});
</script>
<template>
  <compared-bars
    :data1="d1"
    :data2="d2"
    :names="labalName"
    :schema-names="[s1, s2]"
  />
</template>
