<script setup lang="ts">
import Main from "./components/Main.vue";
</script>

<template>
  <n-message-provider>
    <div class="content">
      <div class="header">昙花赛码器</div>
      <n-card>
        <Main />
      </n-card>
    </div>
  </n-message-provider>
</template>
<style>
.content {
  margin: auto;
  max-width: 700px;
  min-height: 90vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  font-family: -apple-system, "Noto Sans", "Helvetica Neue", Helvetica, "Nimbus Sans L", Arial, "Liberation Sans",
    "PingFang SC", "Hiragino Sans GB", "Noto Sans CJK SC", "Source Han Sans SC", "Source Han Sans CN", "Microsoft YaHei",
    "Wenquanyi Micro Hei", "WenQuanYi Zen Hei", "ST Heiti", SimHei, "WenQuanYi Zen Hei Sharp", sans-serif;
}

.header {
  background-image: linear-gradient(to right, orange, rgb(197, 4, 245));
  background-clip: text;
  color: transparent;
  font-size: xx-large;
  font-weight: bold;
  padding: 12px;
  font-family: Baskerville, Georgia, "Liberation Serif", "Kaiti SC", STKaiti, "AR PL UKai CN", "AR PL UKai HK",
    "AR PL UKai TW", "AR PL UKai TW MBE", "AR PL KaitiM GB", KaiTi, KaiTi_GB2312, DFKai-SB, "TW\-Kai", serif;
}
</style>
