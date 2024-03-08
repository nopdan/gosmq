<script setup lang="ts">
import { TextConfig } from "./Text.vue";
import Text from "./Text.vue";

let route = "api/";
const cleanMode = ref(false);

const textList = reactive(new Array<TextConfig>());

function addText(config: TextConfig): void {
  const _new = Object.assign({}, config);
  textList.push(_new);
  console.log("添加文章", _new);
}

/** 从文章列表中删除 */
function removeText(index: number): void {
  console.log("删除文章", textList[index]);
  textList.splice(index, 1);
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

interface Dict extends TextConfig {
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
    label: "冰凌",
    value: DictFormat.Bingling,
  },
  {
    label: "小小 | 极点",
    value: DictFormat.Xiaoxiao,
  },
];

const algoOptions = [
  {
    label: "贪心匹配",
    value: Algorithm.Greedy,
  },
  {
    label: "按照码表顺序",
    value: Algorithm.Ordered,
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

function onDictChange(config: TextConfig) {
  Object.assign(dict, config);
}

const dict = reactive({
  source: "local",
  format: DictFormat.Default,
  push: 4,
  keys: "_;'",
  single: false,
  algo: Algorithm.Greedy,
  space: SpacePreference.Both,
} as Dict);

const dictList = reactive(new Array<Dict>());

function addDict(config: TextConfig): void {
  Object.assign(dict, config);
  const d = Object.assign({}, dict);
  dictList.push(d);
  console.log("添加码表", d);
}

/** 从文章列表中删除 */
function removeDict(index: number): void {
  console.log("删除文章", dictList[index]);
  dictList.splice(index, 1);
}

const localFiles = ref({
  text: [],
  dict: [],
});
function fetchList() {
  fetch(route + "list", {
    method: "GET",
  })
    .then((response) => response.json())
    .then((data) => {
      localFiles.value = data;
      console.log(data);
    })
    .catch((error) => {
      console.error("Error:", error);
    });
}

fetchList();

const activeDrawer = ref(false);
async function race() {
  // 生成 formData
  const formData = new FormData();
  const data = JSON.stringify({
    text: textList,
    dict: dictList,
    clean: cleanMode.value,
  });
  formData.append("data", data);
  console.log("data:", data);

  // post 请求 url
  // fetch 发送 post 请求
  await fetch(route + "race", {
    method: "POST",
    body: formData,
  })
    .then((response) => response.json())
    .then((data) => {
      console.log(data);
      activeDrawer.value = !activeDrawer.value;
    })
    .catch((error) => {
      console.error("Error:", error);
    });
}
</script>

<template>
  <Text :_type="'text'" :files="localFiles" :config-list="textList" @remove="removeText" @add="addText" />
  <!-- 下面是添加码表的部分 -->
  <n-divider />
  <div>
    <Text
      :_type="'dict'"
      :files="localFiles"
      :config-list="dictList"
      @remove="removeDict"
      @add="addDict"
      @update="onDictChange"
    />
    <div class="line">
      <span class="name">码表格式</span>
      <n-radio-group v-model:value="dict.format" name="dictFormat" :disabled="dict.source === 'local'">
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
      <n-radio-group v-model:value="dict.space" name="dictSpace">
        <n-flex>
          <n-radio v-for="space in spaceOptions" :key="space.value" :value="space.value">
            {{ space.label }}
          </n-radio>
        </n-flex>
      </n-radio-group>
    </div>
    <div class="line">
      <span class="name">匹配算法</span>
      <n-radio-group v-model:value="dict.algo" name="dictAlgo" :disabled="dict.single">
        <n-flex>
          <n-radio v-for="algo in algoOptions" :key="algo.value" :value="algo.value" :disabled="algo.disabled">
            {{ algo.label }}
          </n-radio>
        </n-flex>
      </n-radio-group>
    </div>
  </div>

  <div class="line" style="justify-content: right">
    <span style="margin-right: 20px; display: flex; align-items: center">
      <span style="margin-right: 10px">忽略缺字和符号</span>
      <n-switch v-model:value="cleanMode" />
    </span>
    <n-button type="primary" @click="race" :disabled="textList.length === 0 || dictList.length === 0"
      >开始比赛</n-button
    >
  </div>

  <n-drawer v-model:show="activeDrawer" :width="502" placement="bottom" height="100vh">
    <n-drawer-content :native-scrollbar="false" closable>
      <template #header>
        <div style="display: flex; align-items: center; justify-content: center">
          <div style="width: 100%; margin-right: 10px">文章</div>
          <!-- <n-select
						v-model:value="textConfig.name"
						:options="textOptions"
						placeholder="请选择"
						style="min-width: 500px"
					/> -->
        </div>
      </template>
      《斯通纳》是美国作家约翰·威廉姆斯在 1965 年出版的小说。

      <template #footer>
        <n-button>Footer</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<style>
.line {
  display: flex;
  align-items: center;
  min-height: 30px;
  margin: 10px 0;

  & .name {
    min-width: 84px;
    text-align: right;
    padding-right: 21px;
  }
}
</style>
