<script setup lang="ts">
import { ref, onMounted } from "vue";

import TsunamiTable from "@/components/TsunamiTable.vue";

const eventId = ref<string | null>(null);
const isDebug = ref(false);
const isLoading = ref(true);

onMounted(async () => {
    if (window.localStorage.getItem("DBG_DUMMY_TSUNAMI_ID")) {
        isDebug.value = true;
        eventId.value = window.localStorage.getItem("DBG_DUMMY_TSUNAMI_ID")!;
        isLoading.value = false;

        return;
    }

    const resp = await fetch("https://quake-jade.vercel.app/api/events");
    const events = (await resp.json() as { events: string[] }).events;
    // IDにVTSE41が含まれる最初のイベントを取得
    eventId.value = events.find((id) => id.includes("_VTSE41")) ?? null;

    isLoading.value = false;
});
</script>

<template>
    <main>
        <TsunamiTable v-if="!isLoading && eventId" :is-debug="isDebug" :event-id="eventId" />
        <div v-else-if="!isLoading && !eventId" class="no-tsunami">
            <p>津波情報はありません。</p>
        </div>
    </main>
</template>


<style scoped>

.no-tsunami {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 500px;
}

</style>
