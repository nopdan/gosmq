<template>
  <div class="keyboard">
    <div class="keyrow" v-for="line of keys">
      <div v-if="!isStraightKeyboard" :style="{ width: '' + line[1] + 'em' }"></div>
      <div class="key" v-for="k of line[0]" :style="{ backgroundColor: getBackgroundColor(k) }" :title="getKeyTitle(k)">
        <div class="keyname">{{ k.toUpperCase() }}</div>
        <div class="rate">{{ formatFloatToPercent(getRateFromData(k)) }}</div>
      </div>
    </div>
    <div class="keyrow">
      <div v-if="!isStraightKeyboard" style="width: 7.5em" :title="getKeyTitle('left_space')"></div>
      <div
        class="spacekey"
        :style="{ backgroundColor: getBackgroundColor('left_space') }"
        :title="getKeyTitle('left_space')"
      >
        <div class="keyname">左手空格</div>
        <div class="rate">
          {{ formatFloatToPercent(getRateFromData("left_space")) }}
        </div>
      </div>
      <div
        class="spacekey"
        :title="getKeyTitle('left_space')"
        :style="{ backgroundColor: getBackgroundColor('right_space') }"
      >
        <div class="keyname">右手空格</div>
        <div class="rate">
          {{ formatFloatToPercent(getRateFromData("right_space")) }}
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
const p = defineProps(["data", "isStraightKeyboard"]);

function getRateFromData(key: string) {
  return (p.data.Keys as any)?.["" + key]?.Rate;
}

function getBackgroundColor(key: string) {
  let rate = getRateFromData(key) ?? 0;
  if (rate) {
    rate = rate * 10;
  }
  return "rgba(239,68,68," + rate + ")";
}

function getKeyTitle(key: string) {
  let d = (p.data.Keys as any)?.[key];
  if (!d) {
    return "";
  }
  return `按键：${key}\n使用量：${d.Count}\n使用率：${d.Rate}`;
}

function formatFloatToPercent(src: number): string {
  if (isNaN(src) || src * 1 === 0) {
    return "";
  }
  src *= 100;
  return ("" + src).substring(0, 4);
}

const keys: Array<[string, number]> = [
  ["1234567890", 0],
  ["qwertyuiop", 0.6],
  ["asdfghjkl;'", 1.4],
  ["zxcvbnm,./", 2.8],
];
</script>

<style scoped>
.keyboard {
  transform: scale(0.8);
  user-select: none;
  margin: 1em;
}

.keyrow {
  display: flex;
  width: 38em;
  justify-content: left;
}

.key {
  min-width: 2.4em;
  height: 2.4em;
  margin: 4px 3px;
  padding: 1px 2px;
  border-radius: 4px;
  box-shadow: 0 1px 4px #71717a55;
  transition: transform 0.2s;
}

.spacekey {
  min-width: 8.7em;
  height: 2.4em;
  margin: 4px 3px;
  padding: 1px 2px;
  border-radius: 0 0 4px 4px;
  box-shadow: 0 1px 4px #71717a55;
  transition: transform 0.2s;
}
.key:hover {
  transform: scale(110%, 110%);
}

.keyname {
  font-family: "Gill Sans", "Gill Sans MT", Calibri, "Trebuchet MS", sans-serif;
  line-height: 1.6em;
}

.rate {
  font-family: "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
  font-size: 0.6em;
  color: #71717a;
  line-height: 1em;
}
</style>
