<template>
  <div style="width: 1000px; height: 280px; margin: auto; overflow: hidden">
    <v-chart
      class="chart"
      :option="option"
      autoresize
      :update-options="updateoptions"
    />
  </div>
</template>

<script lang="ts" setup>
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { LineChart } from "echarts/charts";
import { TooltipComponent, GridComponent } from "echarts/components";
import VChart, { THEME_KEY } from "vue-echarts";

use([CanvasRenderer, LineChart, TooltipComponent, GridComponent]);
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

const p = defineProps([
  "data1",
  "data2",
  "names",
  "schemaName1",
  "schemaName2",
]);

const data1 = computed(() => p.data1);
const data2 = computed(() => p.data2);
const names = computed(() => p.names);
const s1 = computed(() => p.schemaName1);
const s2 = computed(() => p.schemaName2);
const updateoptions = {
  lazyUpdate: false,
  silent: false,
  notMerge: true,
};

const option = reactive({
  animation: false,
  xAxis: {
    data: names,
    type: "category",
    nameLocation: "middle",
  },
  tooltip: {
    trigger: "axis",
  },
  yAxis: clogy,
  series: [
    {
      name: s1,
      type: "line",
      data: data1,
      color: "#10b981",
      emphasis: {
        focus: "series",
      },
      lineStyle: {
        width: 4,
      },
    },
    {
      name: s2,
      type: "line",
      data: data2,
      color: "#fb7185",
      emphasis: {
        focus: "series",
      },
      lineStyle: {
        width: 4,
      },
    },
  ],
});
</script>

<style scoped></style>
