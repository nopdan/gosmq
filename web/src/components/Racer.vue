<script setup lang="ts">
import InputForm from "./Input/Input.vue";
import Result from "./Result/Result.vue";
import { NSpace } from "naive-ui";
import { reactive, ref } from "vue";
import t1 from "./Result/test1.json";
import t2 from "./Result/test2.json";
import t3 from "./Result/test3.json";

async function start(e: any) {
  let res = await fetch("/api", {
    body: JSON.stringify(e),
    method: "POST",
    headers: {
      "Content-Type": "application/json;charset=utf-8",
    },
  });
  let r: any = {};
  if (res.ok) {
    try {
      r = await res.json();
    } catch (error) {
      alert("参数有误，请检查。\n" + error);
      hasResult.value = false
      return
    }
  }
  data1.value = r[0];
  data2.value = r[1];
  hasResult.value = true;
}
const hasResult = ref(false);
const data1 = ref(null);
const data2 = ref(null);


</script>

<template>
  <NSpace vertical>
    <input-form @start="start" />
    <result
      v-if="hasResult"
      :data1="data1"
      :data2="data2"
      style="min-width: 50em; max-width: 80em; margin: auto"
    />
  </NSpace>
</template>

<style scoped>
input-form {
  max-width: 80rem;
}
</style>
