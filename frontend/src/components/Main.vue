<script setup lang="ts">
import { Data } from "./Data";

// 文本
interface TextConfig {
  source: string;
  path: string;
  name: string | null;
  string: string | null;
}
/** 计算所有文本并合并结果 */
const merge = ref(false);
/** 忽略缺字和符号*/
const clean = ref(false);

/** 文本来源 */
const srcOptions = [
  { label: "本地文件", value: "local" },
  { label: "剪贴板", value: "clipboard" },
];

const text = reactive({
  source: "local",
} as TextConfig);

const textOpts = ref([]);

function paste(type: string) {
  navigator.clipboard
    .readText()
    .then((str) => {
      if (type === "text") {
        text.string = str;
      } else if (type === "dict") {
        dict.string = str;
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
  /** chai */
  Chai = "chai",
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

interface Dict extends TextConfig {
  /** 码表格式 */
  format: DictFormat;
  /** 顶屏模式 */
  pattern: string;
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
    label: "多多",
    value: DictFormat.Duoduo,
  },
  {
    label: "冰凌",
    value: DictFormat.Bingling,
  },
  {
    label: "chai",
    value: DictFormat.Chai,
  },
  {
    label: "小小 | 极点",
    value: DictFormat.Xiaoxiao,
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
    disabled: false,
  },
  // {
  //   label: "最短码长(慢)",
  //   value: Algorithm.Dynamic,
  //   disabled: true,
  // },
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
  source: "local",
  format: DictFormat.Default,
  pattern: "^.{4,}$",
  keys: "_;'",
  single: false,
  algo: Algorithm.Ordered,
  space: SpacePreference.Both,
} as Dict);

const dictList = reactive(new Array<Dict>());
const dictOpts = ref([]);
watch(
  () => text.source,
  () => {
    text.name = null;
    text.string = null;
  },
);
watch(
  () => text.string,
  () => {
    if (text.string) {
      text.name = clipText(text.string, 10);
    }
  },
);

const disableFormat = ref(true);

watch(
  () => dict.source,
  () => {
    dict.name = null;
    dict.string = null;
    if (dict.source === "local") {
      dict.format = DictFormat.Default;
      disableFormat.value = true;
    } else {
      disableFormat.value = false;
    }
  },
);

watch(
  () => dict.string,
  () => {
    if (dict.string) {
      dict.name = clipText(dict.string, 10);
    }
  },
);

watch(
  () => dict.single,
  () => {
    if (dict.single) {
      dict.algo = Algorithm.Ordered;
    }
  },
);

function addDict(dict: Dict): void {
  if (dict.source === "local") {
    dict.name = dict.path.replace(/(.+)(\\|\/)(.+)\.txt/g, "$3");
  }
  const d = Object.assign({}, dict);
  dictList.push(d);
  console.log("添加码表", d);
}

/** 从文章列表中删除 */
function removeDict(index: number): void {
  console.log("删除文章", dictList[index]);
  dictList.splice(index, 1);
}

function tidyPath(path: string, dir: string) {
  const index = path.indexOf(dir);
  let name = path;
  if (index !== -1) {
    name = path.substring(index + dir.length + 1);
  }
  name = name.replace(".txt", "");
  return name;
}

const textDir = ref("text");
const dictDir = ref("dict");
function fetchList() {
  fetch("/list", {
    method: "GET",
  })
    .then((response) => response.json())
    .then((data) => {
      textDir.value = data.textDir;
      dictDir.value = data.dictDir;
      textOpts.value = data.text.map((e: string) => {
        return {
          label: tidyPath(e, textDir.value),
          value: e,
        };
      });
      if (text.path == null) {
        text.path = data.text[0];
      }
      dictOpts.value = data.dict.map((e: string) => {
        return {
          label: tidyPath(e, dictDir.value),
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
  const data = JSON.stringify({
    clean: clean.value,
    merge: merge.value,
    text: text,
    dict: dictList,
  });
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
  <div>
    <div class="line">
      <span class="name">文本来源</span>
      <n-radio-group v-model:value="text.source" size="large">
        <n-flex>
          <n-radio v-for="src in srcOptions" :key="src.value" :value="src.value">
            {{ src.label }}
          </n-radio>
        </n-flex>
      </n-radio-group>
    </div>
    <div v-if="text.source === 'local'">
      <div class="line">
        <span class="name">选择文本</span>
        <n-select v-model:value="text.path" :options="textOpts" placeholder="请选择文本" :disabled="merge" />
      </div>
      <div class="line">
        <span class="name">全部文本</span>
        <n-switch v-model:value="merge"></n-switch>
      </div>
    </div>
    <div v-if="text.source === 'clipboard'">
      <div class="line">
        <span class="name">文本名</span>
        <n-input v-model:value="text.name"></n-input>
      </div>
      <div class="line">
        <span class="name"></span>
        <div style="display: flex">
          <n-button @click="paste('text')" style="margin-right: 10px">读取剪贴板</n-button>
        </div>
        <div v-if="text.string" class="text-preview">
          {{ clipText(text.string || "", 100) }}
        </div>
      </div>
    </div>
  </div>

  <!-- 下面是添加码表的部分 -->
  <n-divider style="margin: 10px auto" />
  <div>
    <n-flex justify="center">
      <n-tag type="success" closable v-for="(conf, idx) in dictList" @close="removeDict(idx)">
        {{ conf.name }}
      </n-tag>
    </n-flex>
    <div class="line">
      <span class="name">码表来源</span>
      <n-radio-group v-model:value="dict.source" size="large">
        <n-flex>
          <n-radio v-for="src in srcOptions" :key="src.value" :value="src.value">
            {{ src.label }}
          </n-radio>
        </n-flex>
      </n-radio-group>
    </div>
    <div v-if="dict.source === 'local'">
      <div class="line">
        <span class="name">选择码表</span>
        <n-select v-model:value="dict.path" :options="dictOpts" placeholder="请选择文件" />
      </div>
    </div>
    <div v-if="dict.source === 'clipboard'">
      <div class="line">
        <span class="name">码表名</span>
        <n-input v-model:value="dict.name"></n-input>
      </div>
      <div class="line">
        <span class="name"></span>
        <div style="display: flex">
          <n-button @click="paste('dict')" style="margin-right: 10px">读取剪贴板</n-button>
        </div>
        <div v-if="dict.string" class="text-preview">
          {{ clipText(dict.string || "", 100) }}
        </div>
      </div>
    </div>
    <div class="line">
      <span class="name">码表格式</span>
      <n-radio-group v-model:value="dict.format" size="large">
        <n-flex>
          <n-radio v-for="format in formatOptions" :key="format.value" :value="format.value" :disabled="disableFormat">
            {{ format.label }}
          </n-radio>
        </n-flex>
      </n-radio-group>
    </div>
    <div class="line">
      <span class="name">顶屏模式</span>
      <n-input
        v-model:value="dict.pattern"
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
      <n-button
        type="info"
        @click="addDict(dict)"
        :disabled="dict.source === 'clipboard' && !dict.string"
        style="width: 100px"
        ghost
      >
        添加</n-button
      >
    </n-flex>
  </div>

  <n-divider style="margin: 15px auto 0" />
  <div>
    <n-flex justify="space-between">
      <div class="line">
        <span class="name">忽略缺字和符号</span>
        <n-switch v-model:value="clean" />
      </div>
      <div class="line">
        <n-button type="primary" @click="race" :disabled="dictList.length === 0" :loading="racing">开始比赛</n-button>
      </div>
    </n-flex>
  </div>
  <Show :result="result"></Show>
</template>

<style scoped>
.line {
  display: flex;
  align-items: center;
  min-height: 30px;
  margin: 15px 0;

  & > .name {
    min-width: 80px;
    text-align: right;
    padding-right: 30px;
    color: #111;
    font-size: 1.1rem;
    font-weight: bold;
    font-family: Baskerville, "Times New Roman", "Liberation Serif", STFangsong, FangSong, FangSong_GB2312, "CWTEX\-F",
      serif;
  }
}

.text-preview {
  color: #999;
  padding: 0 15px;
}
</style>
