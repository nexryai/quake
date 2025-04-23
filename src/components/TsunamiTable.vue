<script setup lang="ts">
import L from "leaflet";
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
    code: string,
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

const mapError = ref(false);
const fatalError = ref(false);
const isLoading = ref(true);
const mapIsLoading = ref(false);

const fetchMapData = async () => {
    const response = await fetch("/JP20250311T.geojson");
    return await response.json();
};

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

        throw ("Abort");
    }
};

const getColorByGrade = (grade: TsunamiGrade): string => {
    switch (grade) {
        case TsunamiGrade.MajorWarning:
            return "rgb(180, 0, 180)";
        case TsunamiGrade.Warning:
            return "rgb(177, 0, 0)";
        case TsunamiGrade.Watch:
            return "rgb(255, 149, 0)";
        default:
            return "gray";
    }
};


onMounted(async () => {
    const geojson = await fetchMapData();

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

    // areaのcodeをkeyとするgradeのmap
    const gradeMap = new Map<string, TsunamiGrade>();
    for (const area of tsunamiData.areas) {
        gradeMap.set(area.code, area.grade);
    }

    try {
        const map = L.map("map", {
            attributionControl: false,
            scrollWheelZoom: false,
            zoomControl: false
        }).setView([35.7894, 136], 4);

        L.geoJSON(geojson, {
            style: function (feature) {
                return {
                    // グレードによって色分け
                    color: getColorByGrade(gradeMap.get(feature!.properties.code) ?? TsunamiGrade.Unknown),
                    weight: 2.5,
                    opacity: 0.7,
                    fillColor: "white",
                    fillOpacity: 1
                };
            },
            onEachFeature: (feature, layer) => {
                if (feature.properties && feature.properties.name) {
                    layer.bindPopup(feature.properties.name);
                }
            },
            pointToLayer: (feature, latlng) => {
                return L.circleMarker(latlng, { radius: 8, color: "red" });
            },
        }).addTo(map);
    } catch (e) {
        mapError.value = true;
        console.error(`Failed to load map: ${e}`);
    }
});

</script>

<template>
    <div class="view">
        <div id="map">
            <div class="loading" v-if="isLoading || mapIsLoading">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon icon-tabler icons-tabler-outline icon-tabler-loader-3">
                    <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                    <path d="M3 12a9 9 0 0 0 9 9a9 9 0 0 0 9 -9a9 9 0 0 0 -9 -9" />
                    <path d="M17 12a5 5 0 1 0 -5 5" />
                </svg>
                <p v-if="isLoading">LOADING...</p>
                <p v-else-if="mapIsLoading">マップ準備中...</p>
            </div>
            <div class="map-error" v-if="mapError">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon icon-tabler icons-tabler-outline icon-tabler-map-x">
                    <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                    <path d="M14 19.5l-5 -2.5l-6 3v-13l6 -3l6 3l6 -3v9" />
                    <path d="M9 4v13" />
                    <path d="M15 7v6.5" />
                    <path d="M22 22l-5 -5" />
                    <path d="M17 22l5 -5" />
                </svg>
                <p>マップのロードに失敗しました</p>
            </div>
        </div>
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

.view {
    @media (min-width: 1110px) {
        display: flex;
        justify-content: space-between;
    }
}

#map {
    @media (min-width: 1110px) {
        margin: 64px;
        width: 320px;
    }

    height: 500px;

    border-radius: 10px;
    overflow: hidden;
    background: rgb(247, 247, 247);

    & .loading {
        z-index: 1000;
        width: 200px;
        margin: 200px auto 12px auto;
        text-align: center;

        & svg {
            display: inline-block;
            animation: rotate-z 2s linear infinite;
        }
    }

    & .map-error {
        position: relative;
        z-index: 9999;
        width: 100%;
        height: 100%;
        padding-top: 200px;
        text-align: center;
        background: rgba(230, 230, 230, 0.665);
        backdrop-filter: blur(18px);
    }
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
        width: 600px;
        margin-bottom: 16px;

        @media (max-width: 650px) {
            width: 100%;
        }

        & .grade-areas {
            margin-left: 12px;

            & .area {
                margin-bottom: 40px;

                @media (min-width: 650px) {
                    display: flex;
                    justify-content: space-between;
                    margin-bottom: 16px;
                }

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

@keyframes rotate-z {
    from {
        transform: rotateZ(0deg);
    }

    to {
        transform: rotateZ(360deg);
    }
}
</style>
