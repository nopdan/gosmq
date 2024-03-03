<template>
  <n-flex justify="center">
    <n-tag type="success" closable v-for="text in textList" @close="removeText(text)">
      {{ cleanName(text.name || "") }}</n-tag>
  </n-flex>
  <div class="line">
    <span class="name">文章来源</span>
    <n-radio-group v-model:value="textConfig.source" name="textDataSource">
      <n-flex>
        <n-radio v-for="src in dataSource" :key="src.value" :value="src.value">
          {{ src.label }}
        </n-radio>
      </n-flex>
    </n-radio-group>
  </div>

  <div class="line">
    <span class="name">文章名称</span>
    <n-input v-model:value="textConfig.name" type="text" :disabled="textConfig.source === DataSource.Local" />
    <n-button style="margin-left: 12px" @click="addText" :disabled="disableText">添加</n-button>
  </div>

  <div class="line" v-if="textConfig.source === DataSource.Local">
    <span class="name">选择文章</span>
    <n-select v-model:value="textConfig.name" :options="textOptions" placeholder="请选择文章" />
  </div>
  <div class="line" v-if="textConfig.source === DataSource.Upload" style="min-height: 150px">
    <span class="name"></span>
    <n-upload :max="1" :action="host + '/upload'" @change="uploadText" v-model:file-list="textFileList">
      <n-upload-dragger>
        <div style="margin-bottom: 12px">
          <n-icon size="48" :depth="3">
            <archive-icon />
          </n-icon>
        </div>
        <n-text style="font-size: 16px">
          点击或者拖动文件到该区域来上传
        </n-text>
      </n-upload-dragger>
    </n-upload>
  </div>

  <div class="line" v-if="textConfig.source === DataSource.Clipboard" style="min-height: 150px">
    <span class="name"></span>
    <div v-if="textConfig.text" class="text-preview">
      {{ clipText(textConfig.text || "", 100) }}
    </div>
    <div class="clipboard" @click="pasteText">
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

  <!-- 下面是添加码表的部分 -->
  <n-divider />
  <n-flex justify="center">
    <n-tag type="info" closable v-for="dict in dictList" @close="removeDict(dict)">
      {{ cleanName(dict.name || "") }}</n-tag>
  </n-flex>
  <div class="line">
    <span class="name">码表来源</span>
    <n-radio-group v-model:value="dictConfig.source" name="dictDataSource">
      <n-flex>
        <n-radio v-for="src in dataSource" :key="src.value" :value="src.value">
          {{ src.label }}
        </n-radio>
      </n-flex>
    </n-radio-group>
  </div>

  <div class="line">
    <span class="name">码表名称</span>
    <n-input v-model:value="dictConfig.name" type="text" :disabled="dictConfig.source === DataSource.Local" />
    <n-button style="margin-left: 12px" @click="addDict" :disabled="disableDict">添加</n-button>
  </div>

  <div class="line" v-if="dictConfig.source === DataSource.Local">
    <span class="name">选择码表</span>
    <n-select v-model:value="dictConfig.name" :options="dictOptions" placeholder="请选择码表" />
  </div>
  <div class="line" v-if="dictConfig.source === DataSource.Upload" style="min-height: 150px">
    <span class="name"></span>
    <n-upload :max="1" :action="host + '/upload'" @change="uploadDict" v-model:file-list="dictFileList">
      <n-upload-dragger>
        <div style="margin-bottom: 12px">
          <n-icon size="48" :depth="3">
            <archive-icon />
          </n-icon>
        </div>
        <n-text style="font-size: 16px">
          点击或者拖动文件到该区域来上传
        </n-text>
      </n-upload-dragger>
    </n-upload>
  </div>

  <div class="line" v-if="dictConfig.source === DataSource.Clipboard" style="min-height: 150px">
    <span class="name"></span>
    <div v-if="dictConfig.text" class="text-preview">
      {{ clipText(dictConfig.text || "", 100) }}
    </div>
    <div class="clipboard" @click="pasteDict">
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
  <div class="line">
    <span class="name">码表格式</span>
    <n-radio-group v-model:value="dictConfig.format" name="dictFormat"
      :disabled="dictConfig.source === DataSource.Local">
      <n-flex>
        <n-radio v-for="format in formatOptions" :key="format.value" :value="format.value">
          {{ format.label }}
        </n-radio>
      </n-flex>
    </n-radio-group>
  </div>
  <div class="line">
    <span class="name">起顶码长</span>
    <n-input-number v-model:value="dictConfig.push" :min="0" style="width: 230px" :disabled="dictConfig.format === DictFormat.Default ||
      dictConfig.format === DictFormat.Jisu
      " />
  </div>
  <div class="line">
    <span class="name">选重键</span>
    <n-input v-model:value="dictConfig.keys" type="text" placeholder="选重键" style="width: 230px"
      :disabled="dictConfig.format === DictFormat.Default" />
  </div>
  <div class="line">
    <span class="name">只打单字</span>
    <n-switch v-model:value="dictConfig.single" />
  </div>
  <div class="line">
    <span class="name">空格偏好</span>
    <n-radio-group v-model:value="dictConfig.space" name="dictSpace">
      <n-flex>
        <n-radio v-for="space in spaceOptions" :key="space.value" :value="space.value">
          {{ space.label }}
        </n-radio>
      </n-flex>
    </n-radio-group>
  </div>
  <div class="line">
    <span class="name">匹配算法</span>
    <n-radio-group v-model:value="dictConfig.algo" name="dictAlgo" :disabled="dictConfig.single">
      <n-flex>
        <n-radio v-for="algo in algoOptions" :key="algo.value" :value="algo.value">
          {{ algo.label }}
        </n-radio>
      </n-flex>
    </n-radio-group>
  </div>

  <div class="line" style="justify-content: right">
    <n-button type="primary" @click="race" :disabled="textList.length === 0 || dictList.length === 0">开始比赛</n-button>
  </div>

  <n-drawer v-model:show="activeDrawer" :width="502" placement="bottom" height="100vh">
    <n-drawer-content :native-scrollbar="false" closable>
      <template #header>
        <div style="display: flex; align-items: center; justify-content: center;">
          <div style="width: 100%;margin-right: 10px;">文章</div>
          <n-select v-model:value="textConfig.name" :options="textOptions" placeholder="请选择"
            style="min-width: 500px;" />
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

<script setup lang="ts">
import { ArchiveOutline as ArchiveIcon } from "@vicons/ionicons5";
import { ClipboardOutline as ClipboardIcon } from "@vicons/ionicons5";
import type { UploadInst, UploadFileInfo } from "naive-ui";
const host = "http://127.0.0.1:7007";

/** 文章或者码表数据来源 */
enum DataSource {
  Local = "local",
  Upload = "upload",
  Clipboard = "clipboard",
}

const dataSource = [
  { label: "本地文件", value: DataSource.Local },
  { label: "上传文件", value: DataSource.Upload },
  { label: "剪切板", value: DataSource.Clipboard },
];

interface TextConfig {
  source: DataSource;
  name: string | null;
  /** 上传文件的 index */
  fileIndex: number | null;
  text: string | null;
}

const textConfig = reactive({
  source: DataSource.Local,
} as TextConfig);

/** 是否禁用“添加文章”按钮 */
const disableText = computed(() => {
  return (
    textConfig.name === "" ||
    (textConfig.source === DataSource.Upload && textConfig.fileIndex == null) ||
    (textConfig.source === DataSource.Clipboard && textConfig.text == null)
  );
});



/** 从剪切板读取文本 */
function pasteText() {
  navigator.clipboard
    .readText()
    .then((text) => {
      textConfig.text = text;
      if (!textConfig.name) {
        textConfig.name = clipText(text || "", 20);
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

const textList = reactive([] as TextConfig[]);
const textFileList = ref<UploadFileInfo[]>([]);

function addText() {
  // 深拷贝 textConfig
  const config = JSON.parse(JSON.stringify(textConfig)) as TextConfig;
  textList.push(config);
  console.log("添加文章", config);
  // 重置
  textConfig.fileIndex = null;
  textConfig.text = null;
}

watch(() => textConfig.fileIndex, (current) => {
  if (current === null) {
    textFileList.value = [];
    textConfig.name = null;
  }
})

/** 从文章列表中删除 */
function removeText(config: TextConfig) {
  textList.splice(textList.indexOf(config), 1);
  console.log("删除文章", config);
}

// 下面是码表的设置选项

enum DictFormat {
  /** 默认赛码表（gosmq 专用） */
  Default = "default",
  /** 极速赛码表 */
  Jisu = "jisu",
  /** 多多格式 */
  Duoduo = "duoduo",
  /** 极点格式 */
  Jidian = "jidian",
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

interface DictConfig extends TextConfig {
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
    label: "多多(Rime)",
    value: DictFormat.Duoduo,
  },
  {
    label: "极点",
    value: DictFormat.Jidian,
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

const dictConfig = reactive({
  source: DataSource.Local,
  format: DictFormat.Default,
  push: 4,
  keys: "_;'",
  single: false,
  algo: Algorithm.Greedy,
  space: SpacePreference.Both,
} as DictConfig);

/** 是否禁用“添加码表”按钮 */
const disableDict = computed(() => {
  return (
    dictConfig.name === "" ||
    (dictConfig.source === DataSource.Upload && dictConfig.fileIndex == null) ||
    (dictConfig.source === DataSource.Clipboard && dictConfig.text == null)
  );
});

watch(
  () => dictConfig.source,
  (current) => {
    if (current === DataSource.Local) {
      dictConfig.format = DictFormat.Default;
    }
  }
);

function pasteDict() {
  navigator.clipboard
    .readText()
    .then((text) => {
      dictConfig.text = text;
      if (!dictConfig.name) {
        dictConfig.name = clipText(text || "", 20);
      }
    })
    .catch((err) => console.error("读取剪贴板内容失败：", err));
}

const dictList = reactive([] as DictConfig[]);
const dictFileList = ref<UploadFileInfo[]>([]);

/** 添加码表 */
function addDict() {
  // 深拷贝 textConfig
  const config = JSON.parse(JSON.stringify(dictConfig)) as DictConfig;
  dictList.push(config);
  console.log("添加码表", config);
  // 重置
  dictConfig.fileIndex = null;
  dictConfig.text = null;
}

watch(() => dictConfig.fileIndex, (current) => {
  if (current === null) {
    dictFileList.value = [];
    dictConfig.name = null;
  }
})

function removeDict(config: DictConfig) {
  dictList.splice(dictList.indexOf(config), 1);
  console.log("删除码表", config.name, config);
}

const textArray = [
  "GB2312字集.txt",
  "small/一些词.txt",
  "small/四季.txt",
  "small/小张和小丽.txt",
  "small/常用单字中五百.txt",
  "small/常用单字前1500.txt",
  "small/常用单字前五百.txt",
  "small/常用单字后五百.txt",
  "small/形码评测加频1500.txt",
  "small/心的出口.txt",
  "small/经济数据.txt",
  "small/青青.txt",
  "《万族之劫》.txt",
  "《三 体》.txt",
  "《围城-钱钟书》.txt",
  "《夜的命名术》.txt",
  "《带着农场混异界》.txt",
  "《庆余年》.txt",
  "《时停499年》（校对版全本）作者：左手萝莉.txt",
  "《红楼梦》-曹雪芹.txt",
  "《西游记》-吴承恩.txt",
  "乌合之众.txt",
  "凡人修仙传之仙界篇(0).txt",
  "影响力.txt",
  "心情决定事情.txt",
  "极品全能高手_花都大少.txt",
  "红楼梦原著.txt",
  "那些热血飞扬的日子（整理版）.txt",
];
const textOptions = textArray.map((e) => {
  return {
    label: e,
    value: e,
  };
});

const dictArray = [
  "091五笔.txt",
  "091点儿2023春.txt",
  "092K.txt",
  "092五笔.txt",
  "092五笔M.txt",
  "1.txt",
  "2.txt",
  "3码单字前瞻版0.1.txt",
  "86王码五笔(含词组).txt",
  "NewWB-20230205.txt",
  "c42单字.txt",
  "flypy.txt",
  "null.txt",
  "倾心单字.txt",
  "倾心打词.txt",
  "哲豆圆满版.txt",
  "哲豆音形快版.txt",
  "小兮码.txt",
  "小可两笔.txt",
  "小鹤.txt",
  "新徐码gb18030_冰凌码表.txt",
  "星二S-1901b赛码表.txt",
  "星二单赛码表(右耍版)20180813.txt",
  "晨逸单字.txt",
  "灵形速影.txt",
  "虎码230303.txt",
  "西风瘦码.txt",
  "超强快码.txt",
  "逸码v20单.txt",
  "键道6赛码表.txt",
];
const dictOptions = dictArray.map((e) => {
  return {
    label: e,
    value: e,
  };
});

/** 去除文件名后缀和书名号 */
function cleanName(name: string): string {
  const index = name.lastIndexOf("/");
  if (index !== -1) {
    name = name.substring(index + 1);
  }
  return name.replace(/《?(.+?)》?(\.txt)?/g, "$1");
}

const activeDrawer = ref(false);
function race() {
  // 生成 formData
  const formData = new FormData();
  formData.append("text", JSON.stringify(textList));
  formData.append("dict", JSON.stringify(dictList));

  // post 请求 url
  const url = host + "/race";

  // fetch 发送 post 请求
  fetch(url, {
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

function uploadDict(options: {
  file: UploadFileInfo;
  fileList: Array<UploadFileInfo>;
  event?: Event;
}) {
  if (options.file.status !== "finished") {
    dictConfig.fileIndex = null;
    return;
  }
  /** 上传完成 */
  dictConfig.name = options.file.name;
  fetch(host + "/file_index", {
    method: "POST",
  })
    .then((res) => {
      if (res.ok) {
        return res.json();
      }
    })
    .then((data) => {
      dictConfig.fileIndex = data.index;
      console.log(options.file.name, "上传完成，index: ", data.index);
    })
    .catch((error) => {
      console.error("Error:", error);
    });
}

function uploadText(options: any) {
  if (options.file.status !== "finished") {
    textConfig.fileIndex = null;
    return;
  }
  /** 上传完成 */
  textConfig.name = options.file.name;
  fetch(host + "/file_index", {
    method: "POST",
  })
    .then((res) => {
      if (res.ok) {
        return res.json();
      }
    })
    .then((data) => {
      textConfig.fileIndex = data.index;
      console.log(options.file.name, "上传完成，index: ", data.index);
    })
    .catch((error) => {
      console.error("Error:", error);
    });
}

</script>
