<script setup lang="ts">
import { ArchiveOutline as ArchiveIcon } from "@vicons/ionicons5";
import { ClipboardOutline as ClipboardIcon } from "@vicons/ionicons5";
import type { UploadFileInfo } from "naive-ui";

export interface TextConfig {
  source: string;
  name: string;
  path: string | null;
  /** 上传文件的 index */
  index: number | null;
  text: string | null;
}
const props = defineProps<{
  _type: "text" | "dict";
  configList: Array<TextConfig>;
  files: { text: string[]; dict: string[] };
}>();
const emit = defineEmits<{
  (e: "remove", idx: number): void;
  (e: "add", value: TextConfig): void;
  (e: "update", value: TextConfig): void;
}>();

/** 文章或者码表数据来源 */
const srcOptions = [
  { label: "本地文件", value: "local" },
  { label: "上传文件", value: "upload" },
  { label: "剪切板", value: "clipboard" },
];

const conf = reactive({ source: "local" } as TextConfig);
watch(conf, () => {
  emit("update", conf);
});

watch(
  () => conf.path,
  () => {
    conf.name = tidyPath(conf.path || "");
  },
);

const _Type = computed(() => {
  return props._type === "text" ? "文章" : "码表";
});

function tidyPath(path: string) {
  const index = path.lastIndexOf(props._type);
  let name = path;
  if (index != -1) {
    name = path.substring(index + 5);
  }
  name = name.replace(".txt", "");
  return name;
}

const opts = computed(() => {
  const tmp = props.files[props._type] as string[];
  return tmp.map((path: string) => {
    return {
      label: tidyPath(path),
      value: path,
    };
  });
});

/** 是否禁用“添加”按钮 */
const disabled = computed(() => {
  return (
    conf.name === "" ||
    (conf.source === "local" && conf.path == null) ||
    (conf.source === "upload" && conf.index == null) ||
    (conf.source === "clipboard" && conf.text == null)
  );
});

const fileList = ref<UploadFileInfo[]>([]);
function upload(options: { file: UploadFileInfo; fileList: Array<UploadFileInfo>; event?: Event }) {
  if (options.file.status !== "finished") {
    conf.index = null;
    return;
  }

  fetch("/file_index", {
    method: "GET",
  })
    .then((res) => {
      res.json().then((data) => {
        /** 上传完成 */
        conf.index = data.index;
        conf.name = options.file.name;
        add(conf);
        fileList.value = [];
        console.log(options.file.name, "上传完成，index: ", data.index);
      });
    })
    .catch((error) => {
      console.error("Error:", error);
    });
}

function paste() {
  navigator.clipboard
    .readText()
    .then((text) => {
      conf.text = text;
      if (!conf.name) {
        conf.name = clipText(text || "", 20);
      }
    })
    .catch((err) => console.error("读取剪贴板内容失败：", err));
}

/** 截取字符串前 n 个字符 */
function clipText(text: string, n: number) {
  if (text.length > n) {
    return text.slice(0, n) + "...";
  }
  return text;
}

function remove(index: number) {
  emit("remove", index);
}
function add(value: TextConfig) {
  emit("add", value);
  conf.index = null;
  conf.text = null;
}
</script>

<template>
  <n-flex justify="center">
    <n-tag type="success" closable v-for="(conf, idx) in configList" @close="remove(idx)">
      {{ conf.name }}
    </n-tag>
  </n-flex>
  <div class="line">
    <span class="name">{{ _Type }}来源</span>
    <n-radio-group v-model:value="conf.source" name="textDataSource">
      <n-flex>
        <n-radio v-for="src in srcOptions" :key="src.value" :value="src.value">
          {{ src.label }}
        </n-radio>
      </n-flex>
    </n-radio-group>
  </div>

  <div class="line">
    <span class="name">{{ _Type }}名称</span>
    <n-input v-model:value="conf.name" type="text" :disabled="conf.source === 'local'" />
    <n-button style="margin-left: 12px" @click="add(conf)" :disabled="disabled">添加</n-button>
  </div>

  <div class="line" v-if="conf.source === 'local'">
    <span class="name">选择{{ _Type }}</span>
    <n-select v-model:value="conf.path" :options="opts" placeholder="请选择文件" />
  </div>
  <div class="line" v-if="conf.source === 'upload'" style="min-height: 150px">
    <span class="name"></span>
    <n-upload :max="1" :action="'/upload'" @change="upload" v-model:file-list="fileList">
      <n-upload-dragger>
        <div style="margin-bottom: 12px">
          <n-icon size="48" :depth="3">
            <archive-icon />
          </n-icon>
        </div>
        <n-text style="font-size: 16px"> 点击或者拖动文件到该区域来上传 </n-text>
      </n-upload-dragger>
    </n-upload>
  </div>

  <div class="line" v-if="conf.source === 'clipboard'" style="min-height: 150px">
    <span class="name"></span>
    <div v-if="conf.text" class="text-preview">
      {{ clipText(conf.text || "", 100) }}
    </div>
    <div class="clipboard" @click="paste">
      <div class="empty">
        <div style="margin-bottom: 12px">
          <n-icon size="48" :depth="3">
            <clipboard-icon />
          </n-icon>
        </div>
        <n-text style="font-size: 16px"> 点击读取系统剪切板 </n-text>
      </div>
    </div>
  </div>
</template>
<style>
.clipboard {
  display: flex;
  justify-content: center;
  width: 100%;
  min-height: 120px;
  border: 1px dashed #e0e0e6;
  border-radius: 3px;
  padding: 14px;
  background-color: #fafafc;
  cursor: pointer;
  align-items: center;

  & .empty {
    text-align: center;
  }

  &:hover {
    border-color: #18a058;
  }

  &:active {
    background-color: #f0f0f0;
  }
}

.text-preview {
  color: #999;
  width: 600px;
  margin-right: 14px;
}
</style>
