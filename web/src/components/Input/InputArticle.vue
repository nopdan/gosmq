<template>
  <n-card size="small">
    <n-form
      ref="formRef"
      :model="result"
      :rules="rules"
      spellcheck="false"
      autocorrect="false"
      label-align="right"
      label-placement="left"
      label-width="auto"
      @change="$emit('changeValue', result)"
      require-mark-placement="right-hanging"
    >
      <n-form-item label="赛文名称：" size="small" path="name" required>
        <n-input
          clearable
          size="small"
          placeholder="请输入名称…"
          style="text-align: left"
          v-model:value="result.name"
        />
      </n-form-item>
      <n-form-item label="赛文文本：" size="small" path="content" required>
        <n-input
          v-model:value="result.content"
          type="textarea"
          placeholder="在这里粘贴文本…"
          clearable
          :autosize="{ minRows: 3, maxRows: 8 }"
          style="text-align: left"
        />
      </n-form-item> </n-form
  ></n-card>
</template>
<script setup lang="ts">
import {
  NInput,
  NRadio,
  NRadioGroup,
  NForm,
  NFormItem,
  FormInst,
  NCard,
} from "naive-ui";
import { reactive, ref, onMounted } from "vue";
const e = defineEmits(["changeValue"]);
const result = reactive({
  name: "",
  content: "",
});
const rules = {
  name: {
    required: true,
    trigger: ["blur", "input"],
    validator: (rule: any, value: string) => {
      if (value.length !== 0) {
        return true;
      } else {
        return new Error("必须填写文章名");
      }
    },
  },
  content: {
    required: true,
    trigger: ["blur", "input"],
    validator: (rule: any, value: string) => {
      if (value.length !== 0) {
        return true;
      } else {
        return new Error("必须填写赛文");
      }
    },
  },
};
const formRef = ref<FormInst | null>();
onMounted(() => {
  e("changeValue", result);
  formRef?.value?.validate()
});
</script>
