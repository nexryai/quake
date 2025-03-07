<script setup lang="ts">
import { ref, onMounted } from "vue";

import TsunamiHazardSign from "./icons/TsunamiHazardSign.vue";


const props = defineProps<{
    eventId: string
    isDebug?: boolean
}>();

// 津波予報
type TsunamiData = {
  expire: string | null;
  issue: TsunamiIssue;
  cancelled: boolean;
  areas: Area[];
};

type TsunamiIssue = {
  source: string;
  time: string;
  type: string;
};

type Area = {
  name: string;
  grade: TsunamiGrade;
  immediate: boolean;
  firstHeight: FirstHeight;
  maxHeight: MaxHeight;
};

type FirstHeight = {
  arrivalTime?: string;
  condition?: string;
};

type MaxHeight = {
  description?: string;
  value?: number;
};

// 津波警報の種類
enum TsunamiGrade {
    MajorWarning = "MajorWarning", // 大津波警報
    Warning = "Warning",           // 津波警報
    Watch = "Watch",               // 津波注意報
    Unknown = "Unknown"            // 不明
}

const isDummyData = ref(false);
const issue = ref<TsunamiIssue | null>(null);
const majorWarningAreas = ref<Area[]>([]);
const warningAreas = ref<Area[]>([]);
const watchAreas = ref<Area[]>([]);
const unknownAreas = ref<Area[]>([]);
const isCancelled = ref(false);

const fatalError = ref(false);
const isLoading = ref(true);

const fetchTsunamiData = async (id: string, isDebug = false): Promise<TsunamiData> => {
    try {
        const response = await fetch(!isDebug
            ? `https://quake-jade.vercel.app/api/events/details?id=${id}`
            : `https://quake-jade.vercel.app/api/events/details?id=${id}&debug=dummy`
        );

        if (response.status != 200) {
            throw new Error(`Status code was not 200: ${response.status}`);
        }

        return await response.json() as TsunamiData;
    } catch (e) {
        fatalError.value = true;
        console.error(`Failed to fetch quake data: ${e}`);

        throw("Abort");
    }
};


onMounted(async () => {
    isDummyData.value = props.isDebug;
    const tsunamiData = await fetchTsunamiData(props.eventId, props.isDebug);

    issue.value = tsunamiData.issue;
    isCancelled.value = tsunamiData.cancelled;

    for (const area of tsunamiData.areas) {
        switch (area.grade) {
            case TsunamiGrade.MajorWarning:
                majorWarningAreas.value.push(area);
                break;
            case TsunamiGrade.Warning:
                warningAreas.value.push(area);
                break;
            case TsunamiGrade.Watch:
                watchAreas.value.push(area);
                break;
            default:
                unknownAreas.value.push(area);
        }
    }

    isLoading.value = false;
});

</script>

<template>
    <div class="grades-table">
        <div v-if="issue" class="info">
            <div>
                <h1>津波情報</h1>
                <p>{{ issue.time }}</p>
            </div>
            <div>
                <p>発行元: <span>{{ issue.source }}</span></p>
            </div>
        </div>
        <div class="full-message">
            <div v-if="isCancelled">
                <span>津波警報は解除されました</span>
            </div>
        </div>
        <div class="grade" v-if="majorWarningAreas.length > 0">
            <div class="header major-warning">
                <span>大津波警報</span>
                <TsunamiHazardSign class="sign" />
            </div>
            <div class="grade-areas">
                <div v-for="area in majorWarningAreas" :key="area.name" class="area">
                    <div class="area-name">
                        <h3 class="important">{{ area.name }}</h3>
                    </div>
                    <div class="area-details">
                        <p v-if="area.firstHeight.arrivalTime" class="important">到達予想時刻: {{ area.firstHeight.arrivalTime }}</p>
                        <p v-if="area.immediate" class="important text-red">直ちに来襲すると予想</p>
                        <p v-if="area.firstHeight.condition" class="important">{{ area.firstHeight.condition }}</p>
                        <p>
                            予想される高さ:
                            <span v-if="area.maxHeight.value" class="text-bg-red">{{ area.maxHeight.value }}m</span>
                            <span v-if="area.maxHeight.description && !area.maxHeight.description.includes('ｍ')" class="text-bg-red">{{ area.maxHeight.description }}</span>
                            <span v-if="!area.maxHeight.value && !area.maxHeight.description">不明</span>
                        </p>
                    </div>
                </div>
            </div>
        </div>
        <div class="grade" v-if="warningAreas.length > 0">
            <div class="header warning">
                <span>津波警報</span>
                <TsunamiHazardSign class="sign" />
            </div>
            <div class="grade-areas">
                <div v-for="area in warningAreas" :key="area.name" class="area">
                    <div class="area-name">
                        <h3>{{ area.name }}</h3>
                    </div>
                    <div class="area-details">
                        <p v-if="area.firstHeight.arrivalTime" class="important">到達予想時刻: {{ area.firstHeight.arrivalTime }}</p>
                        <p v-if="area.firstHeight.condition" class="important">{{ area.firstHeight.condition }}</p>
                        <p>
                            予想される高さ:
                            <span v-if="area.maxHeight.value" class="important">{{ area.maxHeight.value }}m</span>
                            <span v-if="area.maxHeight.description && !area.maxHeight.description.includes('ｍ')" class="important">{{ area.maxHeight.description }}</span>
                            <span v-if="!area.maxHeight.value && !area.maxHeight.description">不明</span>
                        </p>
                    </div>
                </div>
            </div>
        </div>
        <div class="grade" v-if="watchAreas.length > 0">
            <div class="header watch">
                <span>津波注意報</span>
            </div>
            <div class="grade-areas">
                <div v-for="area in watchAreas" :key="area.name" class="area">
                    <div class="area-name">
                        <h3>{{ area.name }}</h3>
                    </div>
                    <div class="area-details">
                        <p v-if="area.firstHeight.arrivalTime">到達予想時刻: {{ area.firstHeight.arrivalTime }}</p>
                        <p v-if="area.firstHeight.condition">{{ area.firstHeight.condition }}</p>
                        <p>
                            予想される高さ:
                            <span v-if="area.maxHeight.value">{{ area.maxHeight.value }}m</span>
                            <span v-if="area.maxHeight.description && !area.maxHeight.description.includes('ｍ')">{{ area.maxHeight.description }}</span>
                            <span v-if="!area.maxHeight.value && !area.maxHeight.description">不明</span>
                        </p>
                    </div>
                </div>
            </div>
        </div>
        <div class="grade" v-if="unknownAreas.length > 0">
            <div class="header">
                <span>不明</span>
            </div>
            <div class="grade-areas">
                <div v-for="area in unknownAreas" :key="area.name" class="area">
                    <div class="area-name">
                        <h3>{{ area.name }}</h3>
                    </div>
                    <div class="area-details">
                        <p v-if="area.firstHeight.arrivalTime">到達予想時刻: {{ area.firstHeight.arrivalTime }}</p>
                        <p v-if="area.firstHeight.condition">{{ area.firstHeight.condition }}</p>
                        <p>
                            予想される高さ:
                            <span v-if="area.maxHeight.value">{{ area.maxHeight.value }}m</span>
                            <span v-if="area.maxHeight.description && !area.maxHeight.description.includes('ｍ')">{{ area.maxHeight.description }}</span>
                            <span v-if="!area.maxHeight.value && !area.maxHeight.description">不明</span>
                        </p>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.sign {
    width: 32px;
    height: 32px;
}

.text-red {
    color: red;
}

.text-bg-red {
    color: rgb(255, 255, 255);
    background-color: rgb(180, 0, 180);
    padding: 2px;
    border-radius: 4px;
    font-weight: bold;
}

.grades-table {
    margin: 16px;

    & .info {
        display: flex;
        justify-content: space-between;
        margin-bottom: 16px;

        & h1 {
            margin: 0;
        }
    }

    & .full-message {
        width: 100%;

        & div {
            margin: 156px auto 12px auto;
            text-align: center;
        }
    }

    & .grade {
        width: 500px;
        margin-bottom: 16px;

        & .grade-areas {
            & .area {
                display: flex;
                justify-content: space-between;
                margin-bottom: 16px;

                & .area-name {
                    & .important {
                        font-weight: bold;
                    }
                }

                & .area-details {
                    width: 300px;

                    & .important {
                        font-weight: bold;
                    }
                }
            }
        }

        & .header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 12px;
            border-radius: 10px;
            margin: 52px 0 32px 0;

            & span {
                font-size: 18px;
                font-weight: bold;
            }
        }

        & .major-warning {
            color: white;
            background-color: rgb(180, 0, 180);
        }

        & .warning {
            color: white;
            background-color: rgb(177, 0, 0);
        }

        & .watch {
            color: white;
            background-color: rgb(255, 149, 0);
        }

        & .unknown {
            color: black;
            background-color: lightgray;
        }
    }

    & h2 {
        margin: 16px 0;
    }

    & h3 {
        margin: 8px 0;
    }
}
</style>
