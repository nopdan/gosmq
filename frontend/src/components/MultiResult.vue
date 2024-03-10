<script setup lang="ts">
import { Data } from "./Data";

const props = defineProps<{
  data: Data;
}>();

const rate = ref(true);

const res = computed(() => {
  if (!rate.value) {
    return props.data;
  }
  let cp = JSON.parse(JSON.stringify(props.data));
  cp.Commit.Word = commitRate(cp.Commit.Word);
  cp.Commit.Collision = commitRate(cp.Commit.Collision);
  cp.Char.Word = charRate(cp.Char.Word);
  cp.Char.Collision = charRate(cp.Char.Collision);
  cp.Keys.LeftHand = keyRate(cp.Keys.LeftHand);
  cp.Keys.RightHand = keyRate(cp.Keys.RightHand);
  cp.Pair.SameHand = pairRate(cp.Pair.SameHand);
  cp.Pair.DiffHand = pairRate(cp.Pair.DiffHand);
  cp.Pair.SameFinger = pairRate(cp.Pair.SameFinger);
  cp.Pair.DiffFinger = pairRate(cp.Pair.DiffFinger);
  cp.Pair.DoubleHit = pairRate(cp.Pair.DoubleHit);
  cp.Pair.TribleHit = pairRate(cp.Pair.TribleHit);
  cp.Pair.SingleSpan = pairRate(cp.Pair.SingleSpan);
  cp.Pair.MultiSpan = pairRate(cp.Pair.MultiSpan);
  cp.Pair.Staggered = pairRate(cp.Pair.Staggered);
  cp.Pair.Disturb = pairRate(cp.Pair.Disturb);
  cp.Pair.LeftToLeft = pairRate(cp.Pair.LeftToLeft);
  cp.Pair.LeftToRight = pairRate(cp.Pair.LeftToRight);
  cp.Pair.RightToLeft = pairRate(cp.Pair.RightToLeft);
  cp.Pair.RightToRight = pairRate(cp.Pair.RightToRight);
  return cp;
});

function commitRate(count: number): string {
  return ((count / props.data.Commit.Count) * 100).toFixed(2) + "%";
}

function charRate(count: number): string {
  return ((count / props.data.Char.Count) * 100).toFixed(2) + "%";
}

function keyRate(count: number): string {
  return ((count / props.data.Keys.Count) * 100).toFixed(2) + "%";
}

function pairRate(count: number): string {
  return ((count / props.data.Pair.Count) * 100).toFixed(2) + "%";
}

function dist(dist: number[]) {
  let res = [] as any[];
  dist.forEach((count, index) => {
    if (count !== 0) {
      res.push({
        len: index,
        count: count,
      });
    }
  });
  return res;
}
</script>

<template>
  <div class="card">
    <div class="cardBlock schemaInfo">
      <div class="dictName">{{ data.Info.DictName }}</div>
      <div class="isSingle" v-if="data.Info.Single">*单</div>
    </div>

    <div class="cardBlock">
      <table class="tableBasic">
        <tr>
          <th>码长</th>
          <th>总键数</th>
          <th>上屏数</th>
          <th>打词</th>
          <th>选重</th>
        </tr>
        <tr>
          <td>{{ data.Keys.CodeLen.toFixed(4) }}</td>
          <td>{{ data.Keys.Count }}</td>
          <td>{{ data.Commit.Count }}</td>
          <td>{{ res.Commit.Word }}</td>
          <td>{{ res.Commit.Collision }}</td>
        </tr>

        <tr>
          <th>词条数</th>
          <th>缺字</th>
          <th>匹配字数</th>
          <th>打词字数</th>
          <th>选重字数</th>
        </tr>
        <tr>
          <td>{{ data.Info.DictLen }}</td>
          <td>{{ data.Han.Lacks }}</td>
          <td>{{ data.Info.TextLen }}</td>
          <td>{{ res.Char.Word }}</td>
          <td>{{ res.Char.Collision }}</td>
        </tr>
      </table>
    </div>

    <div class="cardBlock">
      <table class="tableFeel">
        <tr>
          <th>左手</th>
          <th>左左</th>
          <th>左右</th>
          <th>右左</th>
          <th>右右</th>
        </tr>
        <tr>
          <td>{{ res.Keys.LeftHand }}</td>
          <td>{{ res.Pair.LeftToLeft }}</td>
          <td>{{ res.Pair.LeftToRight }}</td>
          <td>{{ res.Pair.RightToLeft }}</td>
          <td>{{ res.Pair.RightToRight }}</td>
        </tr>
        <tr>
          <th>右手</th>
          <th>同指</th>
          <th>同键</th>
          <th>小跨排</th>
          <th>大跨排</th>
        </tr>
        <tr>
          <td>{{ res.Keys.RightHand }}</td>
          <td>{{ res.Pair.SameFinger }}</td>
          <td>{{ res.Pair.DoubleHit }}</td>
          <td>{{ res.Pair.SingleSpan }}</td>
          <td>{{ res.Pair.MultiSpan }}</td>
        </tr>
        <tr>
          <th>当量</th>
          <th>异手</th>
          <th>异指</th>
          <th>错手</th>
          <th>小指干扰</th>
        </tr>
        <tr>
          <td>{{ (data.Pair.Equivalent / data.Pair.Count).toFixed(4) }}</td>
          <td>{{ res.Pair.DiffHand }}</td>
          <td>{{ res.Pair.DiffFinger }}</td>
          <td>{{ res.Pair.Staggered }}</td>
          <td>{{ res.Pair.Disturb }}</td>
        </tr>
      </table>
    </div>

    <div class="cardBlock stat">
      <table class="statHead">
        <tr>
          <th>码长</th>
        </tr>
        <tr>
          <th>词长</th>
        </tr>
        <tr>
          <th>选重</th>
        </tr>
      </table>
      <div class="statTable" style="white-space: nowrap; scrollbar-gutter: stable; overflow-x: auto">
        <table>
          <tr>
            <td v-for="v in dist(data.Dist.CodeLen)">
              <span style="color: #ff9933">{{ v.len }} </span><span style="color: #555">: </span>
              <span>{{ v.count }}</span>
            </td>
          </tr>
          <tr>
            <td v-for="v in dist(data.Dist.WordLen)">
              <span style="color: #ff9933">{{ v.len }} </span><span style="color: #555">: </span>
              <span>{{ v.count }}</span>
            </td>
          </tr>
          <tr>
            <td v-for="v in dist(data.Dist.Collision)">
              <span style="color: #ff9933">{{ v.len }} </span><span style="color: #555">: </span>
              <span>{{ v.count }}</span>
            </td>
          </tr>
        </table>
      </div>
    </div>

    <div class="cardBlock">
      <div class="keyRateTitle">按键频率 %</div>
      <table class="heatMap"></table>
    </div>
    <n-flex>
      <span>%</span>
      <n-switch v-model:value="rate"></n-switch
    ></n-flex>
  </div>
</template>

<style scoped>
.card {
  display: flex;
  font-size: medium;
  flex-direction: column;
  transition: 0.2s;
  width: 31rem;
  max-width: 400px;
  margin: 2px;
  padding: 10px;
  background-color: #fefefe;
  border: 1px solid #eee;
  border-radius: 10px;

  &:hover {
    box-shadow: 0 0 6px 0 #bbb;
  }
}

.dictName {
  overflow-x: auto;
  white-space: nowrap;
  scrollbar-gutter: stable;
}

.card .cardBlock.schemaInfo {
  display: flex;
  align-items: end;
  justify-content: start;
  font-size: larger;
  font-weight: bold;
  padding-left: 0.2rem;
}

.dictLen {
  padding-left: 1rem;
  font-size: 0.8em;
  color: lightslategray;
  white-space: nowrap;
}

.dictLenValue {
  font-size: 0.7em;
  font-weight: normal;
}

.cardBlock {
  padding: 0.5rem 0;
}

table {
  table-layout: fixed;
}

th {
  text-align: left;
  color: lightslategray;
}

td {
  white-space: nowrap;
}

.tableBasic {
  width: 100%;
}

.keyRateTitle {
  text-align: center;
  font-weight: bold;
  padding-bottom: 0.5rem;
}

.heatMap {
  width: 100%;
  border-spacing: 6px 6px;
}

.key {
  padding: 0.1rem;
  border-radius: 5px;
  box-shadow: 0 0 2px 0 #bbb;
}

td.key {
  text-align: center;
}

.heatMapRate {
  text-align: center;
  font-size: 0.5em;
  color: rgb(0, 0, 4, 0.4);
}

.fin {
  font-size: smaller;
}

.finHeatMap td {
  white-space: nowrap;
}

.tableFeel {
  width: 100%;
}

.cardBlock.stat {
  display: flex;
  align-items: flex-start;
}

.statHead th {
  padding-right: 0.5rem;
  white-space: nowrap;
}

.statTable {
  display: flex;
  overflow-x: auto;
}

.statTable > table td {
  padding-right: 0.2rem;
}
</style>
