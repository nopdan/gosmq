<template>
  <v-chart class="chart" :option="option" autoresize :update-options="updateoptions" />
</template>

<script lang="ts" setup>
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { PieChart } from "echarts/charts";
import { TooltipComponent } from "echarts/components";
import VChart, { THEME_KEY } from "vue-echarts";

use([CanvasRenderer, PieChart, TooltipComponent]);

provide(THEME_KEY, "light");

const p = defineProps<{
  data: number[];
  names: string[];
}>();
const d = ref(p.data);
const n = ref(p.names);

watch(p, () => {
  d.value = p.data;
  n.value = p.names;
});

const updateoptions = {
  lazyUpdate: false,
  silent: false,
  notMerge: true,
};

const colors = {
  red: ["#fb7185", "#f43f5e", "#e11d48", "#be123c"],
  green: ["#34d399", "#10b981", "#059669", "#047857"],
  orange: "#fdba74",
  lime: "#bef264",
};

function getColorList(colorCount: number) {
  switch (colorCount) {
    case 2:
      return [colors.red[1], colors.green[1]];
    case 4:
      return [colors.orange, colors.red[1], colors.green[1], colors.lime];
    case 6:
      return [colors.orange, colors.red[1], colors.red[2], colors.green[2], colors.green[1], colors.lime];
    case 8:
      return [
        colors.orange,
        colors.red[0],
        colors.red[1],
        colors.red[2],
        colors.green[2],
        colors.green[1],
        colors.green[0],
        colors.lime,
      ];
    case 10:
      return [
        colors.orange,
        colors.red[0],
        colors.red[1],
        colors.red[2],
        colors.red[3],
        colors.green[3],
        colors.green[2],
        colors.green[1],
        colors.green[0],
        colors.lime,
      ];
  }
  throw new Error("双色饼图必须含有偶数个项！");
}

const dataList = computed(() => {
  console.assert(n.value.length === d.value.length);
  let result = [];
  for (let i = 0; i < d.value.length; i++) {
    result.push({ value: d.value[i], name: n.value[i] });
  }
  return result;
});

const option = reactive({
  tooltip: {
    trigger: "item",
    formatter: "{b}：{c} ({d}%)",
  },
  series: [
    {
      center: ["50%", "45%"],
      type: "pie",
      data: dataList,
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowOffsetX: 0,
          shadowColor: "rgba(0, 0, 0, 0.5)",
        },
      },
      radius: "50%",
      color: getColorList(p.names.length),
    },
  ],
});
</script>

<style scoped>
.chart {
  height: 400px;
}
</style>
