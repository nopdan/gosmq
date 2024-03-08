<script setup lang="ts">
import Main from "./components/Main.vue";
</script>

<template>
	<n-message-provider>
		<div class="content">
			<n-card>
				<Main />
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
