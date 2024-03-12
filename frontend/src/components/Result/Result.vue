<script setup lang="ts">
import {
  NGrid,
  NGi,
  NH2,
  NH3,
  NSwitch,
  NP,
  NTabs,
  NTabPane,
  NLayout,
  NLayoutHeader,
  NLayoutFooter,
  NLayoutContent,
} from "naive-ui";
import ResultKeyHeatMap from "./ResultKeyHeat.vue";
import HandComp from "./HandComp.vue";
import FingerPie from "./FingerPie.vue";
import BasicDescription from "./BasicDescription.vue";
import CollisionDesciption from "./CollisionDescription.vue";
import WordsDescription from "./WordsDescription.vue";
import CombsDescription from "./CombsDescription.vue";
import CodeLenDescription from "./CodeLenDescription.vue";
import FingersDescription from "./FingersDescription.vue";
import HandsDescription from "./HandsDescription.vue";
import WordDistBar from "./WordsDistBar.vue";
import CollisionDistBar from "./CollisionDistBar.vue";
import CodeLenDistBar from "./CodeLenDistBar.vue";
import CombsDistBar from "./CombsDistBar.vue";
import KeyHeatSorted from "./KeyHeatSorted.vue";
import { Data, New2Old } from "../Data";
import { OldData } from "../OldData";

const props = defineProps<{
  result: Data[];
}>();

const idx1 = ref(0);
const idx2 = ref(0);
if (props.result.length > 1) {
  idx2.value = 1;
}

const opts = computed(() => {
  return props.result.map((d: Data, index: number) => {
    return {
      label: d.Info.DictName,
      value: index,
    };
  });
});

const d1 = computed(() => {
  let _old = New2Old(props.result[idx1.value]);
  shiftEmptyItems(_old);
  return _old;
});

const d2 = computed(() => {
  let _old = New2Old(props.result[idx2.value]);
  shiftEmptyItems(_old);
  return _old;
});

function shiftEmptyItems(schema: OldData) {
  schema.Words.Dist.shift();
  schema.Collision.Dist.shift();
  schema.CodeLen.Dist.shift();
}

const isStraightKeyboard = ref(false);
const logYAxis = ref(false);
provide("logY", logYAxis);
provide("schema1", d1);
provide("schema2", d2);
</script>
<template>
  <n-layout>
    <n-layout-header>
      <n-h2>赛码报告</n-h2>
      <n-p style="display: flex; align-items: center">
        <n-select v-model:value="idx1" :options="opts" style="max-width: 16em" />
        <span style="font: larger bold; padding: 0 20px"> VS </span>
        <n-select v-model:value="idx2" :options="opts" style="max-width: 16em" />
      </n-p>
      <n-p>本报告中的条形图是否使用对数坐标轴？<n-switch v-model:value="logYAxis" size="small" /></n-p><br
    /></n-layout-header>
    <n-tabs animated type="line" :tabs-padding="100">
      <n-tab-pane name="basic" tab="基本">
        <n-h3>码表基本信息</n-h3>
        <n-grid :cols="2" :x-gap="16">
          <n-gi>
            <basic-description :data="d1" />
          </n-gi>
          <n-gi> <basic-description :data="d2" /> </n-gi
        ></n-grid>
      </n-tab-pane>
      <n-tab-pane name="efficiency" tab="效率">
        <n-h3>码长</n-h3>

        <n-grid :cols="2" :x-gap="16">
          <n-gi>
            <code-len-description :data="d1" />
          </n-gi>
          <n-gi>
            <code-len-description :data="d2" />
          </n-gi>
          <n-gi :span="2">
            <code-len-dist-bar />
          </n-gi>
        </n-grid>

        <n-h3>打词</n-h3>
        <n-grid :cols="2" :x-gap="16">
          <n-gi>
            <words-description :data="d1" />
          </n-gi>
          <n-gi> <words-description :data="d2" /> </n-gi>
          <n-gi :span="2">
            <word-dist-bar />
          </n-gi>
        </n-grid>

        <n-h3>重码</n-h3>

        <n-grid :cols="2" :x-gap="16">
          <n-gi>
            <collision-desciption :data="d1" />
          </n-gi>
          <n-gi>
            <collision-desciption :data="d2" />
          </n-gi>
          <n-gi :span="2">
            <collision-dist-bar />
          </n-gi>
        </n-grid>
      </n-tab-pane>
      <n-tab-pane name="feeling" tab="手感">
        <n-h3>按键热力图 </n-h3>
        观察键盘的颜色，用得越多颜色越深，使用率就越高（百分数，省略了 %）。<br />若喜欢直角样式的键盘，可以点这里：<n-switch
          size="small"
          title="切换直排键盘"
          v-model:value="isStraightKeyboard"
        /><br /><br />

        <n-grid :cols="2" :x-gap="4">
          <n-gi style="overflow: hidden">
            <result-key-heat-map :is-straight-keyboard="isStraightKeyboard" :data="d1"
          /></n-gi>
          <n-gi style="overflow: hidden">
            <result-key-heat-map :is-straight-keyboard="isStraightKeyboard" :data="d2" />
          </n-gi>
          <n-gi style="overflow: hidden" :span="2">
            <key-heat-sorted />
          </n-gi>
        </n-grid>
        <n-h3>手指组合</n-h3>

        <n-grid :cols="2" :x-gap="16">
          <n-gi> <combs-description :data="d1" /></n-gi>
          <n-gi> <combs-description :data="d2" /></n-gi>
          <n-gi :span="2">
            <combs-dist-bar />
          </n-gi>
        </n-grid>

        <n-h3>手指使用量</n-h3>
        <n-grid :cols="2" :x-gap="16">
          <n-gi>
            <fingers-description :data="d1" />
          </n-gi>
          <n-gi>
            <fingers-description :data="d2" />
          </n-gi>
          <n-gi> <finger-pie :data="d1" /> </n-gi>
          <n-gi><finger-pie :data="d2" /></n-gi>
        </n-grid>

        <n-h3>双手使用量</n-h3>
        <n-grid :cols="2" :x-gap="16">
          <n-gi><hands-description :data="d1" /></n-gi>
          <n-gi><hands-description :data="d2" /></n-gi>
          <n-gi><hand-comp :data="d1" /></n-gi>
          <n-gi><hand-comp :data="d2" /></n-gi>
        </n-grid>
      </n-tab-pane>
    </n-tabs>

    <n-layout-footer>
      <n-p> 以上数据仅供参考。方案的效率手感以实际打字体验为准。</n-p>
    </n-layout-footer>
  </n-layout>
</template>
