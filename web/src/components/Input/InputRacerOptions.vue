<template>
  <n-card>
    <n-form
      spellcheck="false"
      autocorrect="false"
      label-align="right"
      label-placement="left"
      label-width="auto"
      @change="$emit('changeValue', result)"
    >
      <n-form-item label="赛码器匹配算法：">
        <n-radio-group v-model:value="result.algorithm" >
          <n-space>
            <n-radio v-for="i in algorithms" :key="i[0]" :value="i[0]">
              {{ i[1] }}
            </n-radio></n-space
          >
        </n-radio-group>
      </n-form-item>
    </n-form>
  </n-card>
</template>
<script setup lang="ts">
import { NForm, NFormItem, NRadio, NRadioGroup, NSpace, NCard } from "naive-ui";
import { reactive, onMounted } from "vue";
const e = defineEmits(['changeValue'])
const result = reactive({
  algorithm: "strie",
});

const algorithms = [
  ["strie", "码表顺序查找"],
  ["trie", "最长词条匹配"],
];

onMounted(()=>{
    e('changeValue', result)
})
</script>

