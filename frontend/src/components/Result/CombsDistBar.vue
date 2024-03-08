<script setup lang="ts">
import ComparedBars from "./comparedBars.vue";
import { computed, inject, reactive } from "vue";
const schema1: any = inject("schema1");
const schema2: any = inject("schema2");
const d1 = computed(() => [
  schema1.value.Combs.DoubleHit.Count,
  schema1.value.Combs.TribleHit.Count,
  schema1.value.Combs.SingleSpan.Count,
  schema1.value.Combs.MultiSpan.Count,
  schema1.value.Combs.LongFingersDisturb.Count,
  schema1.value.Combs.LittleFingersDisturb.Count,
]);
const d2 = computed(() => [
  schema2.value.Combs.DoubleHit.Count,
  schema2.value.Combs.TribleHit.Count,
  schema2.value.Combs.SingleSpan.Count,
  schema2.value.Combs.MultiSpan.Count,
  schema2.value.Combs.LongFingersDisturb.Count,
  schema2.value.Combs.LittleFingersDisturb.Count,
]);
const n1 = computed(() => schema1.value.Name);
const n2 = computed(() => schema2.value.Name);

function getLabalName() {
  return ["二连击", "三连击", "小跨排", "大跨排", "错手", "小指干扰"];
}
</script>
<template>
  <compared-bars :data1="d1" :data2="d2" :names="getLabalName()" :schema-names="[n1, n2]" />
</template>
