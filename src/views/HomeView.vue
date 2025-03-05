<script setup lang="ts">
import { ref, onMounted } from "vue";

import QuakeMap from "@/components/QuakeMap.vue";

const eventId = ref("");
const isLoading = ref(true);

onMounted(async () => {
    const resp = await fetch("https://quake-jade.vercel.app/api/events");
    const events = (await resp.json() as { events: string[] }).events;
    eventId.value = events[0];

    isLoading.value = false;
});
</script>

<template>
  <main>
    <QuakeMap v-if=!isLoading :event-id="eventId" />
  </main>
</template>
