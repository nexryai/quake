<script setup lang="ts">
import { computed } from "vue";

import {getQuakeScaleColor} from "@/color.ts";

const props = defineProps<{
  scale: number
}>();

const getDisplayScale = (scale: number): string => {
    switch (scale) {
    case 10: return "1";
    case 20: return "2";
    case 30: return "3";
    case 40: return "4";
    case 45: return "5-";
    case 50: return "5+";
    case 55: return "6-";
    case 60: return "6+";
    case 70: return "7";
    default: return "?";
    }
};

const color = computed(() => getQuakeScaleColor(props.scale));
const displayScale = computed(() => getDisplayScale(props.scale));

</script>

<template>
    <div class="quake-scale-icon" :style="{ backgroundColor: color }">
        {{ displayScale }}
    </div>
</template>

<style scoped>
.quake-scale-icon {
    display: flex;
    justify-content: space-around;
    align-items: center;
    border-radius: 10px;
    font-size: 32px;
    width: 64px;
    height: 64px;
}
</style>