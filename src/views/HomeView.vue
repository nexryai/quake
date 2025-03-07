<script setup lang="ts">
import { ref, onMounted } from "vue";

import QuakeMap from "@/components/QuakeMap.vue";

const eventId = ref("");
const isDebug = ref(false);
const isLoading = ref(true);

onMounted(async () => {
    if (window.localStorage.getItem("DBG_DUMMY_EVENT_ID")) {
        isDebug.value = true;
        eventId.value = window.localStorage.getItem("DBG_DUMMY_EVENT_ID")!;
        isLoading.value = false;

        return;
    }

    const resp = await fetch("https://quake-jade.vercel.app/api/events");
    const events = (await resp.json() as { events: string[] }).events;
    // IDに_VXSE51または_VXSE53が含まれる最初のイベントを取得
    eventId.value = events.find((id) => id.includes("_VXSE51") || id.includes("_VXSE53")) ?? "";
    if (!eventId.value) {
        console.error("No event found");
        return;
    }

    isLoading.value = false;
});
</script>

<template>
    <main>
        <QuakeMap v-if=!isLoading :is-debug="isDebug" :event-id="eventId" />
    </main>
</template>
