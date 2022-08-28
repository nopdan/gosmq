<template>
  <article>
    <n-grid :cols="2" x-gap="6" y-gap="6">
      <n-gi>
        <input-dict @change-value="(e) => (result.dict1 = e)">
          <template #head>码表 A</template></input-dict
        ></n-gi
      >
      <n-gi>
        <input-dict @change-value="(e) => (result.dict2 = e)">
          <template #head>码表 B</template></input-dict
        ></n-gi
      >
      <n-gi span="2">
        <input-article @change-value="(e) => (result.article = e)" />
      </n-gi>
      <n-gi span="2">
        <input-racer-options @change-value="(e) => (result.racer = e)" />
      </n-gi>
      <n-gi span="2">
        <n-button v-if="checkToStart()" type="primary" round @click="$emit('start', result)">
          <template #icon>
            <n-icon>
              <send />
            </n-icon>
          </template>
          开始计算！</n-button
        ></n-gi
      >
    </n-grid>
  </article>
</template>
<script setup lang="ts">
import InputDict from "./InputDict.vue";
import InputArticle from "./InputArticle.vue";
import InputRacerOptions from "./InputRacerOptions.vue";
import { NGrid, NGi, NButton, NIcon } from "naive-ui";
import { Send } from "@vicons/tabler";
import { reactive } from "vue";

const e = defineEmits(["start"]);

// TODO: text area 太慢，提供拖动上传文件功能

function checkToStart() {
  if (
    result.dict1.name &&
    result.dict1.content &&
    result.dict2.name &&
    result.dict2.content &&
    result.article.name &&
    result.article.content
  ) {
    return true;
  } else false;
}

const result: any = reactive({
  dict1: {},
  dict2: {},
  article: {},
  racer: {},
});

function start() {}
</script>
<style scoped>
article {
  background-color: #f1f3f5;
  width: 100%;
  padding: 1rem 0;
}
.n-grid {
  max-width: 60rem;
  margin: auto;
  padding-top: 1rem;
}
</style>
