<template>
  <n-space vertical>
    <n-form
      ref="formRef"
      :model="result"
      :rules="rules"
      label-placement="left"
      label-align="left"
      label-width="auto"
      require-mark-placement="right-hanging"
    >
      <n-form-item class="dictSelect" label="选择码表" path="path">
        <n-select
          v-model:value="result.path"
          placeholder="请选择码表"
          :options="props.files"
        />
      </n-form-item>
      <n-form-item label="码表格式" title="帮助程序分析文本数据。">
        <n-radio-group v-model:value="result.format">
          <n-radio-button v-for="i in formats" :key="i[0]" :value="i[0]">
            {{ i[1] }}</n-radio-button
          >
        </n-radio-group>
      </n-form-item>
      <n-form-item label="选重键" title="从首选开始，不限选重按键数量。">
        <n-input
          v-model:value="result.selectKeys"
          style="
            font-family: 'Courier New', Courier, monospace;
            text-align: left;
            max-width: 50%;
          "
        />
      </n-form-item>
      <n-form-item label="顶屏码长" title="多长编码后不自动添加空格键？">
        <n-input-number
          v-model:value="result.pushStart"
          :min="0"
          :max="20"
          style="max-width: 50%"
        />
      </n-form-item>
      <n-form-item label="只打单字" title="只用码表里的单字输入赛文。">
        <n-switch v-model:value="result.single" />
      </n-form-item>
      <n-form-item label="匹配算法">
        <n-radio-group v-model:value="result.alg" default-value="longest">
          <n-radio v-for="i in algs" :key="i[0]" :value="i[0]">
            {{ i[1] }}</n-radio
          >
        </n-radio-group>
      </n-form-item>
    </n-form>
  </n-space>
</template>
<script setup lang="ts">
import { FormInst } from "naive-ui";

export interface Dict {
  path: string;
  format: string;
  single: boolean;
  selectKeys: string;
  pushStart: number;
  alg: string;
}

const formRef = ref<FormInst | null>();
const result = reactive({
  path: null,
  format: "js",
  single: false,
  selectKeys: "_;'",
  pushStart: 4,
  alg: "longest",
});

const props = defineProps(["msg", "files", "idx"]);
props.msg.dicts[props.idx] = result;

const formats = [
  ["js", "极速"],
  ["dd", "多多"],
  ["jd", "极点"],
  ["bl", "冰凌"],
];

const algs = [
  ["longest", "最长匹配"],
  ["stable_trie", "按码表顺序"],
];

const rules = {
  path: {
    required: true,
    message: "请选择码表",
    trigger: ["blur"],
  },
};

onMounted(() => {});
</script>
