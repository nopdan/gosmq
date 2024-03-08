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

const p = defineProps(["data1", "data2"]);
const d1 = ref(p.data1);
const d2 = ref(p.data2);
watch(p, () => {
  d1.value = p.data1;
  d2.value = p.data2;
});

function shiftEmptyItems(schema: any) {
  schema.Words.Dist.shift();
  schema.Collision.Dist.shift();
  schema.CodeLen.Dist.shift();
}
shiftEmptyItems(p.data1);
shiftEmptyItems(p.data2);
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
      <n-p> {{ p.data1.DictName }} VS {{ p.data2.DictName }}</n-p>
      <n-p>本报告中的条形图是否使用对数坐标轴？<n-switch v-model:value="logYAxis" size="small" /></n-p><br
    /></n-layout-header>
    <n-tabs animated type="line" :tabs-padding="100">
      <n-tab-pane name="basic" tab="基本">
        <n-h3>码表基本信息</n-h3>
        <n-grid :cols="2" :x-gap="16">
          <n-gi>
            <basic-description :data="p.data1" />
          </n-gi>
          <n-gi> <basic-description :data="p.data2" /> </n-gi
        ></n-grid>
      </n-tab-pane>
      <n-tab-pane name="efficiency" tab="效率">
        <n-h3>码长</n-h3>

        <n-grid :cols="2" :x-gap="16">
          <n-gi>
            <code-len-description :data="p.data1" />
          </n-gi>
          <n-gi>
            <code-len-description :data="p.data2" />
          </n-gi>
          <n-gi :span="2">
            <code-len-dist-bar />
          </n-gi>
        </n-grid>

        <n-h3>打词</n-h3>
        <n-grid :cols="2" :x-gap="16">
          <n-gi>
            <words-description :data="p.data1" />
          </n-gi>
          <n-gi> <words-description :data="p.data2" /> </n-gi>
          <n-gi :span="2">
            <word-dist-bar />
          </n-gi>
        </n-grid>

        <n-h3>重码</n-h3>

        <n-grid :cols="2" :x-gap="16">
          <n-gi>
            <collision-desciption :data="p.data1" />
          </n-gi>
          <n-gi>
            <collision-desciption :data="p.data2" />
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
            <result-key-heat-map :is-straight-keyboard="isStraightKeyboard" :data="p.data1"
          /></n-gi>
          <n-gi style="overflow: hidden">
            <result-key-heat-map :is-straight-keyboard="isStraightKeyboard" :data="p.data2" />
          </n-gi>
          <n-gi style="overflow: hidden" :span="2">
            <key-heat-sorted />
          </n-gi>
        </n-grid>
        <n-h3>手指组合</n-h3>

        <n-grid :cols="2" :x-gap="16">
          <n-gi> <combs-description :data="p.data1" /></n-gi>
          <n-gi> <combs-description :data="p.data2" /></n-gi>
          <n-gi :span="2">
            <combs-dist-bar />
          </n-gi>
        </n-grid>

        <n-h3>手指使用量</n-h3>
        <n-grid :cols="2" :x-gap="16">
          <n-gi>
            <fingers-description :data="p.data1" />
          </n-gi>
          <n-gi>
            <fingers-description :data="p.data2" />
          </n-gi>
          <n-gi> <finger-pie :data="d1" /> </n-gi>
          <n-gi><finger-pie :data="d2" /></n-gi>
        </n-grid>

        <n-h3>双手使用量</n-h3>
        <n-grid :cols="2" :x-gap="16">
          <n-gi><hands-description :data="p.data1" /></n-gi>
          <n-gi><hands-description :data="p.data2" /></n-gi>
          <n-gi><hand-comp :data="p.data1" /></n-gi>
          <n-gi><hand-comp :data="p.data2" /></n-gi>
        </n-grid>
      </n-tab-pane>
    </n-tabs>

    <n-layout-footer>
      <n-p> 以上数据仅供参考。方案的效率手感以实际打字体验为准。</n-p>
    </n-layout-footer>
  </n-layout>
</template>
