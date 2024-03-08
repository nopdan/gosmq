<template>
  <div style="width: 700px; height: 280px; margin: auto; overflow: hidden">
    <v-chart class="chart" :option="option" autoresize :update-options="updateoptions" />
  </div>
</template>

<script lang="ts" setup>
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { BarChart } from "echarts/charts";
import { TooltipComponent, GridComponent } from "echarts/components";
import VChart, { THEME_KEY } from "vue-echarts";
import { reactive, provide, inject, computed } from "vue";

use([CanvasRenderer, BarChart, TooltipComponent, GridComponent]);
provide(THEME_KEY, "light");

const logY: any = inject("logY");
const clogy = computed(() => {
  if (logY.value) {
    return {
      type: "log",
      min: 1,
      logBase: 10,
    };
  } else {
    return {};
  }
});

const p = defineProps(["data1", "data2", "names", "schemaNames"]);

const d1 = computed(() => p.data1);
const d2 = computed(() => p.data2);
const n = computed(() => p.names);
const sn1 = computed(() => p.schemaNames[0]);
const sn2 = computed(() => p.schemaNames[1]);

const updateoptions = {
  lazyUpdate: false,
  silent: true,
  notMerge: true,
};

const option = reactive({
  xAxis: {
    data: n,
    type: "category",
  },
  tooltip: {
    trigger: "item",
    formatter: "{a}<br />{b}：{c}",
  },
  yAxis: clogy,
  emphasis: {
    itemStyle: {
      shadowBlur: 10,
      shadowOffsetX: 0,
      shadowColor: "rgba(0, 0, 0, 0.5)",
    },
  },
  series: [
    {
      name: sn1,
      type: "bar",
      data: d1,
      barGap: "20%",
      barCategoryGap: "30%",
      color: "#6ee7b7",
      emphasis: {
        focus: "series",
      },
    },
    {
      name: sn2,
      type: "bar",
      data: d2,
      color: "#fb7185",
      emphasis: {
        focus: "series",
      },
    },
  ],
});
</script>

<style scoped></style>
