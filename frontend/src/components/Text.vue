<template>
  <n-form
    ref="formRef"
    :model="model"
    :rules="rules"
    label-placement="left"
    label-align="left"
    label-width="auto"
    require-mark-placement="right-hanging"
  >
    <div style="display: flex">
      <n-form-item v-if="model.flag" class="half" label="赛文名称" path="name">
        <n-input v-model:value="model.name" placeholder="请输入赛文名称……" />
      </n-form-item>
      <n-form-item v-else class="half" label="选择赛文" path="path">
        <n-select
          v-model:value="model.path"
          placeholder="请选择赛文"
          :options="props.files"
        />
      </n-form-item>
      <div style="margin: 0 auto"></div>
      <n-form
        :model="model"
        label-placement="left"
        label-align="left"
        label-width="auto"
      >
        <n-form-item label="手动输入赛文" path="flag">
          <n-switch v-model:value="model.flag" />
        </n-form-item>
      </n-form>
    </div>
    <n-form-item v-if="model.flag" label="赛文文本" path="plain">
      <n-input
        v-model:value="model.plain"
        placeholder="请输入赛文文本……"
        type="textarea"
        :autosize="{
          minRows: 3,
          maxRows: 5,
        }"
      />
    </n-form-item>
  </n-form>
</template>
<script setup lang="ts">
import { FormInst } from "naive-ui";

export interface Text {
  name: string,
  plain: string,
  path: string,
  flag: boolean,
}
const formRef = ref<FormInst | null>(null);
const model = reactive({
  name: null,
  plain: null,
  path: null,
  flag: false,
});
const props = defineProps(["msg", "files"]);
props.msg.text = model;

const rules = ref({
  path: {
    required: true,
    message: "请选择赛文",
    trigger: "blur",
  },
  plain: {
    required: true,
    message: "请输入文本",
    trigger: ["input", "blur"],
  },
});
</script>
<style>
.half {
  width: 46%;
}
</style>
