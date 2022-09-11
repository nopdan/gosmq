<template>
  <div id="input">
    <div class="card">
      <Text :msg="options" :files="textList"></Text>
    </div>
    <div id="dicts">
      <Dict class="card half" :msg="options" :files="dictList" :idx="0"></Dict>
      <Dict class="card half" :msg="options" :files="dictList" :idx="1"></Dict>
    </div>
    <div id="buttons">
      <n-space>
        <n-button
          :loading="loading1"
          :disabled="illegel"
          type="primary"
          @click="start1"
          ghost
        >
          开始比赛 by 单单
        </n-button>
        <n-button
          :loading="loading2"
          :disabled="illegel"
          type="primary"
          @click="start2"
          ghost
        >
          开始比赛 by yyb
        </n-button>
      </n-space>
    </div>
  </div>
  <div>
    <div id="result" v-if="hasResult">
      <result
        :data1="data1"
        :data2="data2"
        style="min-width: 50em; max-width: 80em; margin: auto"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import Dict from "./Dict.vue";
import Text from "./Text.vue";
let options = ref({
  text: Text,
  dicts: [Dict, Dict],
});

const textList = new Array();
const dictList = new Array();
const loading1 = ref(false);
const loading2 = ref(false);
const illegel = computed(() => {
  return (
    !options.value.dicts[0].path ||
    !options.value.dicts[1].path ||
    !(options.value.text.flag
      ? options.value.text.plain
      : options.value.text.path)
  );
});

// 获取本地赛文和码表 ./text ./dict
async function getFiles(url: string, li: any[]) {
  // making a call to API
  const response = await fetch(url);
  // converting it to JSON format
  const data = await response.json();
  for (const v of data) {
    li.push({ label: v.split(".")[0], value: v });
  }
}

const dev = false;
const url_api = dev ? "http://localhost:7172/api" : "/api";
const url_result = dev ? "http://localhost:7172/result" : "/result";

onBeforeMount(() => {
  getFiles(dev ? "http://localhost:7172/texts" : "/texts", textList);
  getFiles(dev ? "http://localhost:7172/dicts" : "/dicts", dictList);
});

async function start1() {
  loading1.value = true;
  await fetch(url_api, {
    body: JSON.stringify(options.value),
    method: "POST",
    headers: {
      "Content-Type": "application/json;charset=utf-8",
    },
  });
  window.open(url_result);
  loading1.value = false;
}

async function start2() {
  loading2.value = true;
  let res = await fetch(url_api, {
    body: JSON.stringify(options.value),
    method: "POST",
    headers: {
      "Content-Type": "application/json;charset=utf-8",
    },
  });
  let r: any = {};
  if (res.ok) {
    try {
      r = await res.json();
    } catch (error) {
      alert("参数有误，请检查。\n" + error);
      hasResult.value = false;
      loading2.value = false;
      return;
    }
  }
  data1.value = r[0];
  data2.value = r[1];
  hasResult.value = true;
  loading2.value = false;
}
const hasResult = ref(false);
const data1 = ref(null);
const data2 = ref(null);
</script>
<style scoped>
#input {
  margin: auto;
  width: 72em;
}

#buttons {
  display: flex;
  justify-content: center;
  /* font-size: 25px;
    margin: 1em auto; */
}

#dicts {
  display: flex;
}

.card.half {
  width: 46%;
  margin: 1em 1%;
}

.card {
  margin: 1em 1% 0;
  padding: 1.5em;
  padding-bottom: 0%;
  box-shadow: 0 0 1px 0 #bbb;
  border-radius: 5px;
  background-color: white;
}

#result {
  background-color: white;
  margin-top: 2em;
  padding-top: 2em;
}
</style>
