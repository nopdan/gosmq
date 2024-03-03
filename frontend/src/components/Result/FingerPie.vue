<template>
  <div style="width: 500px; height: 400px; margin: auto; overflow: hidden">
    <bicolor-pie
      :data="reformData"
      :names="[
        '右手拇指',
        '右手食指',
        '右手中指',
        '右手无名指',
        '右手小指',
        '左手小指',
        '左手无名指',
        '左手中指',
        '左手食指',
        '左手拇指',
      ]"
    />
  </div>
</template>

<script lang="ts" setup>
import { watch,ref,computed } from "vue";
import BicolorPie from "./BicolorPie.vue";
const p = defineProps(["data"]);

const d = ref(p.data)
watch(p,()=>{
  d.value = p.data
})
const reformData =computed(()=> {
  const orders = [6,7,8,9,0,1,2,3,4,5];
  let fingers = d.value.Fingers.Dist;
  return orders.map((e) => fingers[e].Count);
})
</script>
