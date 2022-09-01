<script setup lang="ts">
import Input from "./components/Input.vue";
// This starter template is using Vue 3 <script setup> SFCs
// Check out https://vuejs.org/api/sfc-script-setup.html#script-setup
</script>

<template>
  <div class="title">
    <a href="https://github.com/cxcn/gosmq">赛码器 GO!</a>
  </div>
  <n-message-provider>
    <Input />
  </n-message-provider>
</template>

<style scoped>
.title {
  display: flex;
  justify-content: center;
}
.title a {
  font-size: 2em;
  color: black;
  margin: 0.5em 0 0.2em;
  text-decoration: none;
}
body {
  background-color: rgb(240, 244, 245);
}
</style>
