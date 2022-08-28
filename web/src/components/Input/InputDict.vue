<template>
  <n-card role="dialog" size="small">
    <template #header> <slot name="head"></slot> </template>
    <n-form
      ref="formRef"
      :model="result"
      spellcheck="false"
      autocorrect="false"
      label-align="right"
      label-placement="left"
      label-width="auto"
      @change="$emit('changeValue', result)"
      require-mark-placement="right-hanging"
      :rules="rules"
    >
      <n-form-item label="码表名称：" size="small" path="name" required>
        <n-input
          clearable
          size="small"
          placeholder="请输入名称…"
          style="text-align: left"
          v-model:value="result.name"
        />
      </n-form-item>
      <n-form-item label="码表文本：" size="small" path="content" required>
        <n-input
          v-model:value="result.content"
          type="textarea"
          placeholder="在这里粘贴文本…"
          clearable
          :autosize="{ minRows: 3, maxRows: 3 }"
          style="text-align: left;overflow-x: hidden;"
        />
      </n-form-item>
      <n-form-item
        label="码表格式："
        size="small"
        title="帮助程序分析文本数据。"
      >
        <n-radio-group
          v-model:value="result.format"
          size="small"
          default-value="js"
        >
          <n-space>
            <n-radio v-for="i in formats" :key="i[0]" :value="i[0]">
              {{ i[1] }}</n-radio
            ></n-space
          >
        </n-radio-group>
      </n-form-item>
      <n-form-item size="small" label="顶屏码长："
        ><n-input-number
          size="small"
          :default-value="4"
          v-model:value="result.commitLeng"
          title="多长编码后不自动添加空格键？"
          :min="0"
          :max="20"
          style="width: 8rem"
      /></n-form-item>

      <n-form-item
        label="选重键："
        size="small"
        title="从首选开始，不限选重按键数量。"
      >
        <n-input
          v-model:value="result.collidedKeys"
          style="
            font-family: 'Courier New', Courier, monospace;
            width: 8rem;
            text-align: left;
          "
          default-value="_2345"
        />
      </n-form-item>
      <n-form-item
        label="只打单字："
        size="small"
        title="只用码表里的单字输入赛文。"
        ><n-switch v-model:value="result.singleMode"
      /></n-form-item>
    </n-form>
  </n-card>
</template>
<script setup lang="ts">
export interface dictResultType {
  name: string;
  content: string;
  singleMode: boolean;
  commitLeng: number;
  collidedKeys: string;
  format: string;
}

import {
  NInput,
  NForm,
  NFormItem,
  NCard,
  NSwitch,
  NInputNumber,
  NRadio,
  NRadioGroup,
  FormItemRule,
  FormInst,
  NSpace,
} from "naive-ui";
import { reactive, onMounted, ref } from "vue";
const e = defineEmits(["changeValue"]);
const formRef = ref<FormInst | null>();
const result: dictResultType = reactive({
  name: "",
  content: "",
  singleMode: false,
  commitLeng: 4,
  collidedKeys: "_2345",
  format: "js",
});

const formats = [
  ["js", "极速"],
  ["dd", "多多"],
  ["jd", "极点"],
  ["bl", "冰凌"],
];

const rules = {
  content: {
    required: true,
    trigger: ["blur", "input"],
    validator: (rule: FormItemRule, value: string) => {
      if (value.length !== 0) {
        return true;
      } else {
        return new Error("必须填写码表");
      }
    },
  },
  name: {
    required: true,
    trigger: ["blur", "input"],
    validator: (rule: FormItemRule, value: string) => {
      if (value.length !== 0) {
        return true;
      } else {
        return new Error("必须填写名称");
      }
    },
  },
};

onMounted(() => {
  formRef?.value?.validate();
  e("changeValue", result);
});
</script>
<style scoped></style>
