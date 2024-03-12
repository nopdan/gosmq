<script setup lang="ts">
import { ClipboardOutline as ClipboardIcon } from "@vicons/ionicons5";
import { Data } from "./Data";

interface Config {
  // 文本
  source: string;
  path: string;
  text: string;
  /** 计算所有文本并合并结果 */
  merge: boolean;
  /** 忽略缺字和符号*/
  clean: boolean;
}

/** 文本来源 */
const srcOptions = [
  { label: "本地文件", value: "local" },
  { label: "剪贴板", value: "clipboard" },
];

const conf = reactive({
  source: "local",
} as Config);

const textOpts = ref([]);

function paste() {
  navigator.clipboard
    .readText()
    .then((text) => {
      conf.text = text;
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

// 下面是码表的设置选项
enum DictFormat {
  /** 默认赛码表（gosmq 专用） */
  Default = "default",
  /** 极速赛码表 */
  Jisu = "jisu",
  /** 多多 */
  Duoduo = "duoduo",
  /** 冰凌 */
  Bingling = "bingling",
  /** 小小 */
  Xiaoxiao = "xiaoxiao",
}

enum Algorithm {
  /** 贪心匹配 */
  Greedy = "greedy",
  /** 按照码表顺序 */
  Ordered = "ordered",
  /** 动态规划，最短码长 */
  Dynamic = "dynamic",
}

enum SpacePreference {
  /** 互击 */
  Both = "both",
  /** 左手 */
  Left = "left",
  /** 右手 */
  Right = "right",
}

interface Dict {
  path: string;
  /** 码表格式 */
  format: DictFormat;
  /** 起顶码长 */
  push: number;
  /** 选重键 */
  keys: string;
  /** 是否只用码表里的单字 */
  single: boolean;
  /** 匹配算法 */
  algo: Algorithm;
  /** 空格偏好 */
  space: SpacePreference;
}

const formatOptions = [
  {
    label: "默认",
    value: DictFormat.Default,
  },
  {
    label: "极速赛码表",
    value: DictFormat.Jisu,
  },
  {
    label: "多多 | Rime",
    value: DictFormat.Duoduo,
  },
  {
    label: "小小 | 极点",
    value: DictFormat.Xiaoxiao,
  },
  {
    label: "冰凌",
    value: DictFormat.Bingling,
  },
];

const algoOptions = [
  {
    label: "按照码表顺序",
    value: Algorithm.Ordered,
  },
  {
    label: "贪心匹配",
    value: Algorithm.Greedy,
  },
  {
    label: "最短码长(慢)",
    value: Algorithm.Dynamic,
    disabled: true,
  },
];

const spaceOptions = [
  {
    label: "互击",
    value: SpacePreference.Both,
  },
  {
    label: "总是左手",
    value: SpacePreference.Left,
  },
  {
    label: "总是右手",
    value: SpacePreference.Right,
  },
];

const dict = reactive({
  format: DictFormat.Default,
  push: 4,
  keys: "_;'",
  single: false,
  algo: Algorithm.Ordered,
  space: SpacePreference.Both,
} as Dict);

const dictList = reactive(new Array<Dict>());
const dictOpts = ref([]);

function addDict(dict: Dict): void {
  const d = Object.assign({}, dict);
  dictList.push(d);
  console.log("添加码表", d);
}

/** 从文章列表中删除 */
function removeDict(index: number): void {
  console.log("删除文章", dictList[index]);
  dictList.splice(index, 1);
}

function tidyPath(path: string, suffix: string) {
  const index = path.lastIndexOf(suffix);
  let name = path;
  if (index != -1) {
    name = path.substring(index + 5);
  }
  name = name.replace(".txt", "");
  return name;
}

function fetchList() {
  fetch("/list", {
    method: "GET",
  })
    .then((response) => response.json())
    .then((data) => {
      textOpts.value = data.text.map((e: string) => {
        return {
          label: tidyPath(e, "text"),
          value: e,
        };
      });
      if (conf.path == null) {
        conf.path = data.text[0];
      }
      dictOpts.value = data.dict.map((e: string) => {
        return {
          label: tidyPath(e, "dict"),
          value: e,
        };
      });
      if (dict.path == null) {
        dict.path = data.dict[0];
      }
      console.log(data);
    })
    .catch((error) => {
      console.error("Error:", error);
    });
}

fetchList();

/** 总的结果 */
const racing = ref(false);
const result = ref<Data[]>([]);
async function race() {
  racing.value = true;
  // 生成 formData
  const formData = new FormData();
  const data = JSON.stringify(Object.assign({}, conf, { dict: dictList }));
  formData.append("data", data);
  console.log("data:", data);

  // post 请求 url
  // fetch 发送 post 请求
  await fetch("/race", {
    method: "POST",
    body: formData,
  })
    .then((response) => response.json())
    .then((data) => {
      fetchList();
      racing.value = false;
      result.value = data as Data[];
      console.log(data);
    })
    .catch((error) => {
      racing.value = false;
      console.error("Error:", error);
    });
}
</script>

<template>
  <div class="line">
    <span class="name">文本来源</span>
    <n-radio-group v-model:value="conf.source" size="large">
      <n-flex>
        <n-radio v-for="src in srcOptions" :key="src.value" :value="src.value">
          {{ src.label }}
        </n-radio>
      </n-flex>
    </n-radio-group>
  </div>
  <div v-if="conf.source === 'local'">
    <div class="line">
      <span class="name">选择文本</span>
      <n-select v-model:value="conf.path" :options="textOpts" placeholder="请选择文件" :disabled="conf.merge" />
    </div>
    <div class="line">
      <span class="name">全部文本</span>
      <n-switch v-model:value="conf.merge"></n-switch>
    </div>
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
        <n-text style="font-size: 16px"> 点击读取系统剪贴板 </n-text>
      </div>
    </div>
  </div>

  <!-- 下面是添加码表的部分 -->
  <n-divider />
  <div>
    <n-flex justify="center">
      <n-tag type="success" closable v-for="(conf, idx) in dictList" @close="removeDict(idx)">
        {{ tidyPath(conf.path, "dict") }}
      </n-tag>
    </n-flex>
    <div class="line">
      <span class="name">选择码表</span>
      <n-select v-model:value="dict.path" :options="dictOpts" placeholder="请选择文件" />
    </div>
    <div class="line">
      <span class="name">码表格式</span>
      <n-radio-group v-model:value="dict.format" size="large">
        <n-flex>
          <n-radio v-for="format in formatOptions" :key="format.value" :value="format.value">
            {{ format.label }}
          </n-radio>
        </n-flex>
      </n-radio-group>
    </div>
    <div class="line">
      <span class="name">起顶码长</span>
      <n-input-number
        v-model:value="dict.push"
        :min="0"
        style="width: 230px"
        :disabled="dict.format === DictFormat.Default || dict.format === DictFormat.Jisu"
      />
    </div>
    <div class="line">
      <span class="name">选重键</span>
      <n-input
        v-model:value="dict.keys"
        type="text"
        placeholder="选重键"
        style="width: 230px"
        :disabled="dict.format === DictFormat.Default"
      />
    </div>
    <div class="line">
      <span class="name">只打单字</span>
      <n-switch v-model:value="dict.single" />
    </div>
    <div class="line">
      <span class="name">空格偏好</span>
      <n-radio-group v-model:value="dict.space" size="large">
        <n-flex>
          <n-radio v-for="space in spaceOptions" :key="space.value" :value="space.value">
            {{ space.label }}
          </n-radio>
        </n-flex>
      </n-radio-group>
    </div>
    <div class="line">
      <span class="name">匹配算法</span>
      <n-radio-group v-model:value="dict.algo" size="large" :disabled="dict.single">
        <n-flex>
          <n-radio v-for="algo in algoOptions" :key="algo.value" :value="algo.value" :disabled="algo.disabled">
            {{ algo.label }}
          </n-radio>
        </n-flex>
      </n-radio-group>
    </div>
    <n-flex justify="center">
      <n-button type="info" @click="addDict(dict)" style="width: 100px" ghost> 添加</n-button>
    </n-flex>
  </div>

  <n-divider />
  <div class="line" style="justify-content: right">
    <span style="margin-right: 20px; display: flex; align-items: center">
      <span style="margin-right: 10px">忽略缺字和符号</span>
      <n-switch v-model:value="conf.clean" />
    </span>
    <n-button type="primary" @click="race" :disabled="dictList.length === 0" :loading="racing">开始比赛</n-button>
  </div>
  <Show :result="result"></Show>
</template>

<style>
.line {
  display: flex;
  align-items: center;
  min-height: 30px;
  margin: 20px 0;

  & > .name {
    min-width: 80px;
    text-align: right;
    padding-right: 30px;
    color: #111;
    font-size: 125%;
    font-weight: bold;
    font-family: Baskerville, "Times New Roman", "Liberation Serif", STFangsong, FangSong, FangSong_GB2312, "CWTEX\-F",
      serif;
  }
}

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
