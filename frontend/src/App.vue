<script setup lang="ts">
import NewText from './components/NewText.vue';

</script>

<template>
  <n-message-provider>
    <div class="content">
      <n-card>
        <NewText />
      </n-card>
    </div>
  </n-message-provider>
</template>
<style>
.content {
  margin: 1em auto;
  max-width: 800px;
}
</style>
