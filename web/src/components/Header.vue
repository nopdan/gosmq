<script setup lang="ts">
import { ref } from "vue";
import { Puzzle } from "@vicons/tabler";
import { NIcon } from "naive-ui";
const needActiveSetDrawer = ref(false);
function activate() {
  needActiveSetDrawer.value = true;
}
</script>
<template>
  <n-icon><puzzle /></n-icon>
  <a
    href="https://github.com/cxcn/gosmq"
    target="_blank"
    style="text-decoration: none"
  >
    赛码器Go！</a
  >
</template>
