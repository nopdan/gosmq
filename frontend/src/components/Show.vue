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
    <n-drawer-content :native-scrollbar="false" header-style="padding: 0;" body-content-style="padding: 0;">
      <template #header>
        <div class="drawer-header">
          <div class="button-left">
            <n-button @click="multi = !multi">切换风格</n-button>
          </div>
          <div class="button-right">
            <n-button type="info" ghost @click="active = false">关闭</n-button>
          </div>
          <div class="title">{{ props.result[0].Info.TextName }}</div>
        </div>
      </template>
      <div v-if="multi" class="multi-res">
        <div v-for="data in props.result">
          <MultiResult :data="data"></MultiResult>
        </div>
      </div>
      <Result v-else :result="props.result" />
    </n-drawer-content>
  </n-drawer>
</template>
<style scoped>
.multi-res {
  display: flex;
  justify-content: safe center;
  margin: 0 auto;
  flex-wrap: wrap;
  overflow-x: auto;
}

.drawer-header {
  width: 100vw;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  height: 2.75em;

  & .title {
    clear: both;
    font-size: 1.75rem;
    font-weight: bold;
    font-family: Baskerville, Georgia, "Liberation Serif", "Kaiti SC", STKaiti, "AR PL UKai CN", "AR PL UKai HK",
      "AR PL UKai TW", "AR PL UKai TW MBE", "AR PL KaitiM GB", KaiTi, KaiTi_GB2312, DFKai-SB, "TW\-Kai", serif;
  }

  & .button-left {
    position: absolute;
    float: left;
    left: 0;
    margin-left: 1rem;
  }

  & .button-right {
    position: absolute;
    float: right;
    right: 0;
    margin-right: 1rem;
  }
}
</style>
