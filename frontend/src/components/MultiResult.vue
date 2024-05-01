<script setup lang="ts">
import { Data, DataUtils } from "./Data";

const props = defineProps<{
  data: Data;
}>();

const rate = ref(true);

const res = computed(() => {
  if (!rate.value) {
    return props.data;
  }
  let cp = JSON.parse(JSON.stringify(props.data));
  cp.Commit.Word = util.commitRate(cp.Commit.Word);
  cp.Commit.Collision = util.commitRate(cp.Commit.Collision);
  cp.Char.Word = util.charRate(cp.Char.Word);
  cp.Char.Collision = util.charRate(cp.Char.Collision);
  cp.Keys.LeftHand = util.keyRate(cp.Keys.LeftHand);
  cp.Keys.RightHand = util.keyRate(cp.Keys.RightHand);
  cp.Pair.SameHand = util.pairRate(cp.Pair.SameHand);
  cp.Pair.DiffHand = util.pairRate(cp.Pair.DiffHand);
  cp.Pair.SameFinger = util.pairRate(cp.Pair.SameFinger);
  cp.Pair.DiffFinger = util.pairRate(cp.Pair.DiffFinger);
  cp.Pair.DoubleHit = util.pairRate(cp.Pair.DoubleHit);
  cp.Pair.TribleHit = util.pairRate(cp.Pair.TribleHit);
  cp.Pair.SingleSpan = util.pairRate(cp.Pair.SingleSpan);
  cp.Pair.MultiSpan = util.pairRate(cp.Pair.MultiSpan);
  cp.Pair.Staggered = util.pairRate(cp.Pair.Staggered);
  cp.Pair.Disturb = util.pairRate(cp.Pair.Disturb);
  cp.Pair.LeftToLeft = util.pairRate(cp.Pair.LeftToLeft);
  cp.Pair.LeftToRight = util.pairRate(cp.Pair.LeftToRight);
  cp.Pair.RightToLeft = util.pairRate(cp.Pair.RightToLeft);
  cp.Pair.RightToRight = util.pairRate(cp.Pair.RightToRight);
  return cp;
});

const util = new DataUtils(props.data);

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
    <div class="title">
      <div class="name">
        {{ data.Info.DictName }}<span style="color: red; margin-left: 10px" v-if="data.Info.Single">*单</span>
      </div>
    </div>

    <table class="pure-table">
      <thead>
        <tr>
          <th>码长</th>
          <th>总键数</th>
          <th>上屏数</th>
          <th>匹配字数</th>
          <th>词条数</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>{{ data.Keys.CodeLen.toFixed(4) }}</td>
          <td>{{ data.Keys.Count }}</td>
          <td>{{ data.Commit.Count }}</td>
          <td>{{ data.Char.Count }}</td>
          <td>{{ data.Info.DictLen }}</td>
        </tr>
      </tbody>

      <thead>
        <tr>
          <th>打词</th>
          <th>打词字数</th>
          <th>选重</th>
          <th>选重字数</th>
          <th>缺字</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>{{ data.Commit.Word }}</td>
          <td>{{ data.Char.Word }}</td>
          <td>{{ data.Commit.Collision }}</td>
          <td>{{ data.Char.Collision }}</td>
          <td>{{ data.Han.Lacks }}</td>
        </tr>
      </tbody>
      <tbody>
        <tr>
          <td>{{ res.Commit.Word }}</td>
          <td>{{ res.Char.Word }}</td>
          <td>{{ res.Commit.Collision }}</td>
          <td>{{ res.Char.Collision }}</td>
          <td></td>
        </tr>
      </tbody>
    </table>
    <table class="pure-table">
      <thead>
        <tr>
          <th>双击</th>
          <th>三连击</th>
          <th>左手</th>
          <th>左右</th>
          <th>左左</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>{{ res.Pair.DoubleHit }}</td>
          <td>{{ res.Pair.TribleHit }}</td>
          <td>{{ res.Keys.LeftHand }}</td>
          <td>{{ res.Pair.LeftToRight }}</td>
          <td>{{ res.Pair.LeftToLeft }}</td>
        </tr>
      </tbody>
      <thead>
        <tr>
          <th>异手</th>
          <th>同指</th>
          <th>右手</th>
          <th>右左</th>
          <th>右右</th>
        </tr>
      </thead>
      <tr>
        <td>{{ res.Pair.DiffHand }}</td>
        <td>{{ res.Pair.SameFinger }}</td>
        <td>{{ res.Keys.RightHand }}</td>
        <td>{{ res.Pair.RightToLeft }}</td>
        <td>{{ res.Pair.RightToRight }}</td>
      </tr>
      <thead>
        <tr>
          <th>当量</th>
          <th>小跨排</th>
          <th>大跨排</th>
          <th>错手</th>
          <th>小指干扰</th>
        </tr>
      </thead>
      <tr>
        <td>{{ (data.Pair.Equivalent / data.Pair.Count).toFixed(4) }}</td>
        <td>{{ res.Pair.SingleSpan }}</td>
        <td>{{ res.Pair.MultiSpan }}</td>
        <td>{{ res.Pair.Staggered }}</td>
        <td>{{ res.Pair.Disturb }}</td>
      </tr>
    </table>

    <div class="dist">
      <table class="dist-header">
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
      <table class="dist-body">
        <tbody>
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
        </tbody>
      </table>
    </div>

    <n-flex justify="center" style="font-weight: bold; margin-bottom: 8px">按键频率 %</n-flex>
    <div class="heat-map">
      <div class="heat-line">
        <div v-for="v in '0123456789-='"><heat-map :data="data" :_key="v"></heat-map></div>
      </div>
      <div class="heat-line">
        <div v-for="v in 'qwertyuiop[]'"><heat-map :data="data" :_key="v"></heat-map></div>
      </div>
      <div class="heat-line">
        <div v-for="v in 'asdfghjkl;\''"><heat-map :data="data" :_key="v"></heat-map></div>
      </div>
      <div class="heat-line">
        <div v-for="v in 'zxcvbnm,./'"><heat-map :data="data" :_key="v"></heat-map></div>
      </div>
      <div class="heat-line">
        <finger-heat-map :data="data" :colspan="1" :name="'小指'" :num="1"></finger-heat-map>
        <finger-heat-map :data="data" :colspan="1" :name="'无名'" :num="2"></finger-heat-map>
        <finger-heat-map :data="data" :colspan="1" :name="'中指'" :num="3"></finger-heat-map>
        <finger-heat-map :data="data" :colspan="1" :name="'食指'" :num="4"></finger-heat-map>
        <finger-heat-map :data="data" :colspan="1" :name="'拇指'" :num="5"></finger-heat-map>
        <finger-heat-map :data="data" :colspan="1" :name="'拇指'" :num="6"></finger-heat-map>
        <finger-heat-map :data="data" :colspan="1" :name="'食指'" :num="7"></finger-heat-map>
        <finger-heat-map :data="data" :colspan="1" :name="'中指'" :num="8"></finger-heat-map>
        <finger-heat-map :data="data" :colspan="1" :name="'无名'" :num="9"></finger-heat-map>
        <finger-heat-map :data="data" :colspan="1" :name="'小指'" :num="10"></finger-heat-map>
      </div>
    </div>
  </div>
</template>

<style scoped>
div.heat-map {
  & .heat-line {
    display: flex;
    flex-wrap: wrap;
  }
}

.card {
  font-size: 1rem;
  display: flex;
  flex-direction: column;
  margin: 1rem 0.5rem;
  padding: 0.8rem 1.5rem 1.5rem;
  background-color: #fefefe;
  border: 1px solid #eee;
  border-radius: 10px;
  width: 33em;
}

.name {
  font-family: Baskerville, "Times New Roman", "Liberation Serif", STFangsong, FangSong, FangSong_GB2312, "CWTEX\-F",
    serif;
  font-size: 1.3rem;
  font-weight: bold;
  overflow-x: auto;
  white-space: nowrap;
  scrollbar-gutter: stable;
}

div.dist {
  display: flex;
  margin-top: 5px;
  max-width: 550px;
  align-items: start;
  height: 105px;

  & table {
    white-space: nowrap;
  }
}

table.dist-body {
  table-layout: fixed;
  display: flex;
  overflow-x: auto;

  & td {
    padding-left: 10px;
  }
}

.pure-table {
  border-collapse: collapse;
  border-spacing: 0;
  empty-cells: show;
  border: 1px solid #cbcbcb;
  margin: 5px 0;

  & td,
  th {
    border-left: 1px solid #cbcbcb;
    border-width: 0 0 0 1px;
    font-size: inherit;
    margin: 0;
    overflow: visible;
    padding: 3px 1rem;
    white-space: nowrap;
  }
}

.pure-table caption {
  color: #000;
  padding: 1rem 0;
  text-align: center;
}

.pure-table thead {
  background-color: rgb(224, 238, 225);
  color: #000;
  text-align: left;
  vertical-align: bottom;
}
</style>
