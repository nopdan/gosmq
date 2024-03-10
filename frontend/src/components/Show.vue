<script setup lang="ts">
import { Data } from "./Data";

const props = defineProps<{
  result: Data[][];
}>();

const active = ref(false);
watch(
  () => props.result,
  () => {
    active.value = true;
  },
);

const resOptions = computed(() => {
  return props.result.map((item, index) => {
    return {
      label: item[0].Info.TextName,
      value: index,
    };
  });
});
/** 某个文本的结果 */
const resIndex = ref(0);
const res = computed(() => {
  return props.result[resIndex.value];
});

const idx1 = ref(0);
const idx2 = ref(0);

const opts = computed(() => {
  return res.value.map((d: Data, index: number) => {
    return {
      label: d.Info.DictName,
      value: index,
    };
  });
});

watch(
  () => res,
  () => {
    if (res.value) {
      idx1.value = 0;
      idx2.value = Math.min(res.value.length - 1, 1);
    }
  },
);

const multi = ref(true);
</script>
<template>
  <n-drawer v-model:show="active" :width="502" placement="bottom" height="100vh">
    <n-drawer-content :native-scrollbar="false" closable>
      <template #header>
        <span style="width: auto">
          <!-- <n-select v-model:value="resIndex" :options="resOptions" placeholder="请选择" style="min-width: 500px" /> -->
          <span>标题</span>
        </span>
      </template>
      <div>
        <div v-if="multi" style="display: flex; width: 100%; margin: auto; justify-content: center">
          <div v-for="data in res">
            <MultiResult :data="data"></MultiResult>
          </div>
        </div>
        <div v-else>
          <div style="display: flex; align-items: center; justify-content: center">
            <n-select v-model:value="idx1" :options="opts" placeholder="请选择文件" style="max-width: 300px" />
            <span style="font: larger bold; padding: 0 20px"> VS </span>
            <n-select v-model:value="idx2" :options="opts" placeholder="请选择文件" style="max-width: 300px" />
          </div>

          <!-- <Result :data1="props.res[idx1]" :data2="props.res[idx2]"></Result> -->
          {{ JSON.stringify(res[idx1]) }}
          <p />
          {{ JSON.stringify(res[idx2]) }}
        </div>
      </div>
    </n-drawer-content>
  </n-drawer>
</template>
