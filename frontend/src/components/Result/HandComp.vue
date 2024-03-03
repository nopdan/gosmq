<template>
  <bicolor-pie :data="getDataList()" :names="['右手→右手','右手→左手','左手→右手','左手→左手']"/>
</template>

<script lang="ts" setup>
import BicolorPie from "./BicolorPie.vue";
const p = defineProps(["data"]);
function getDataList(){
  let names = ['RightToRight','RightToLeft','LeftToRight','RightToRight']
  return names.map((e)=>p.data.Hands[e].Count)
}

</script>
