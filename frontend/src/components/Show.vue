<script setup lang="ts">
import { Data, New2Old } from "./Data";
import { Close as CloseIcon } from "@vicons/ionicons5";

const props = defineProps<{
  result: Data[];
}>();

const active = ref(false);
watch(
  () => props.result,
  () => {
    active.value = true;
  },
);

const multi = ref(true);
</script>
<template>
  <n-drawer v-model:show="active" :width="502" placement="bottom" height="100vh">
    <n-drawer-content :native-scrollbar="false" header-style="padding: 0; display: flex; margin: auto">
      <template #header>
        <n-flex justify="space-between" style="width: 80vw; align-items: center; padding: 10px">
          <n-button @click="multi = !multi">风格</n-button>
          <span class="title">{{ props.result[0].Info.TextName }}</span>
          <n-button type="info" ghost @click="active = false">关闭</n-button>
        </n-flex>
      </template>
      <div>
        <div v-if="multi" style="display: flex; width: 100%; margin: auto; justify-content: center">
          <div v-for="data in props.result">
            <MultiResult :data="data"></MultiResult>
          </div>
        </div>
        <div v-else>
          <Result :result="props.result" style="min-width: 50em; max-width: 80em; margin: auto" />
        </div>
      </div>
    </n-drawer-content>
  </n-drawer>
</template>
<style scoped>
.title {
  font-size: 30px;
  font-weight: bold;
  font-family: Baskerville, Georgia, "Liberation Serif", "Kaiti SC", STKaiti, "AR PL UKai CN", "AR PL UKai HK",
    "AR PL UKai TW", "AR PL UKai TW MBE", "AR PL KaitiM GB", KaiTi, KaiTi_GB2312, DFKai-SB, "TW\-Kai", serif;
}
</style>
