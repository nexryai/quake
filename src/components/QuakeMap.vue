<script setup lang="ts">
import L from "leaflet";
import { ref, onMounted, computed } from "vue";

import "leaflet/dist/leaflet.css";
import { getQuakeScaleColor } from "@/color.ts";
import QuakeScaleIcon from "@/components/core/QuakeScaleIcon.vue";
import UndrawError from "@/components/images/UndrawError.vue";

const props = defineProps<{
    eventId: string
    isDebug?: boolean
}>();

type EarthquakeData = {
    expire: null;
    issue: {
        source: string;
        time: string;
        type: string;
        correct: string;
    };
    earthquake: {
        time: string;
        hypocenter: {
            name: string;
            latitude: number;
            longitude: number;
            depth: number;
            magnitude: number;
        };
        maxScale: number;
        domesticTsunami: string;
        foreignTsunami: string;
    };
    points: Point[];
};

type Point = {
    code: string;
    areaCode: string;
    pref: string;
    addr: string;
    scale: number;
    isArea: boolean;
};

enum TsunamiStatus {
    UNKNOWN,
    WARNING,
    NOTICE,
    NONE
}

const tsunamiStatusStyle = computed(() => {
    switch (tsunamiStatus.value) {
        case TsunamiStatus.WARNING:
            return { backgroundColor: "#b20000", color: "white" };
        case TsunamiStatus.NOTICE:
            return { backgroundColor: "#FFD700", color: "black" };
        case TsunamiStatus.NONE:
            return { backgroundColor: "#f6f6f6", color: "black" };
        default:
            return { backgroundColor: "#727272", color: "white" };
    }
});

let isForeignQuake = false;

const isDummyData = ref(false);
const quakeDateTimeLabel = ref("");
const quakeScale = ref(0);
const tsunamiStatusLabel = ref("");
const tsunamiStatus = ref(TsunamiStatus.UNKNOWN);
const magnitude = ref<number | string>(0);
const hypocenterLabel = ref("不明");

const mapError = ref(false);
const fatalError = ref(false);
const isLoading = ref(true);
const mapIsLoading = ref(false);

const fetchMapData = async () => {
    const response = await fetch("/JP20250304.geojson");
    return await response.json();
};

const fetchQuakeData = async (id: string, isDebug = false): Promise<EarthquakeData> => {
    try {
        const response = await fetch(!isDebug
            ? `https://quake-jade.vercel.app/api/events/details?id=${id}`
            : `https://quake-jade.vercel.app/api/events/details?id=${id}&debug=dummy`
        );

        if (response.status != 200) {
            throw new Error(`Status code was not 200: ${response.status}`);
        }

        return await response.json() as EarthquakeData;
    } catch (e) {
        fatalError.value = true;
        console.error(`Failed to fetch quake data: ${e}`);

        throw("Abort");
    }
};

onMounted(async () => {
    const dispHypocenter = !props.eventId.includes("_VXSE51");

    // GeoJSON データ
    const geojson = await fetchMapData();

    // Debug flag
    isDummyData.value = props.isDebug;
    const quakeData = await fetchQuakeData(props.eventId, props.isDebug);

    quakeDateTimeLabel.value = quakeData.earthquake.time;
    quakeScale.value = quakeData.earthquake.maxScale;
    magnitude.value = dispHypocenter ? quakeData.earthquake.hypocenter.magnitude : "?";
    hypocenterLabel.value = dispHypocenter ? quakeData.earthquake.hypocenter.name : "調査中";

    //quakeData.earthquake.domesticTsunami = "Watch";

    switch (quakeData.earthquake.domesticTsunami) {
        case "None":
            tsunamiStatusLabel.value = "この地震による津波の心配はありません";
            tsunamiStatus.value = TsunamiStatus.NONE;
            break;
        case "Unknown":
            tsunamiStatusLabel.value = "津波への影響は不明です";
            break;
        case "Checking":
            tsunamiStatusLabel.value = "津波への影響は調査中です。注意して行動してください。";
            break;
        case "NonEffective":
            tsunamiStatusLabel.value = "若干の海面変動が予想されますが、被害の心配はありません。";
            tsunamiStatus.value = TsunamiStatus.NOTICE;
            break;
        case "Watch":
            tsunamiStatusLabel.value = "津波注意報が発表されました";
            tsunamiStatus.value = TsunamiStatus.WARNING;
            break;
        case "Warning":
            tsunamiStatusLabel.value = "津波警報等が発表されました";
            tsunamiStatus.value = TsunamiStatus.WARNING;
            break;
        default: tsunamiStatusLabel.value = "津波情報の取得に失敗しました";
            tsunamiStatus.value = TsunamiStatus.UNKNOWN;
    }

    // 国内でWarningがあるなら海外からの津波は無視
    if (quakeData.earthquake.foreignTsunami != "None" && quakeData.earthquake.foreignTsunami != "Unknown" && tsunamiStatus.value != TsunamiStatus.WARNING) {
        if (quakeData.earthquake.maxScale == -1) {
            isForeignQuake = true;
        }

        switch (quakeData.earthquake.foreignTsunami) {
            case "Checking":
                if (tsunamiStatus.value != TsunamiStatus.NOTICE) {
                    // NOTICEレベルの通知があるならそっちを優先する
                    tsunamiStatusLabel.value = "津波への影響は調査中です。";
                }
                break;
            case "NonEffectiveNearby":
                tsunamiStatusLabel.value = "震源の近傍で小さな津波が予想されますが、被害の心配はありません。";
                tsunamiStatus.value = TsunamiStatus.NOTICE;
                break;
            case "WarningPacific":
                tsunamiStatusLabel.value = "太平洋で津波が予想されます";
                tsunamiStatus.value = TsunamiStatus.WARNING;
                break;
            case "WarningPacificWide":
                tsunamiStatusLabel.value = "太平洋の広域で津波が予想されます";
                tsunamiStatus.value = TsunamiStatus.WARNING;
                break;
            case "WarningIndian":
                tsunamiStatusLabel.value = "インド洋で津波が予想されます";
                tsunamiStatus.value = TsunamiStatus.WARNING;
                break;
            case "WarningIndianWide":
                tsunamiStatusLabel.value = "インド洋の広域で津波が予想されます";
                tsunamiStatus.value = TsunamiStatus.WARNING;
                break;
            case "Potential":
                tsunamiStatusLabel.value = "この規模の地震では津波の可能性があります";
                tsunamiStatus.value = TsunamiStatus.NOTICE;
                break;
            default: tsunamiStatusLabel.value = "津波情報（海外）の取得に失敗しました";
                tsunamiStatus.value = TsunamiStatus.NOTICE;
        }
    }

    // areaCodeごとの最大scaleを求める
    const areaScaleMap = new Map<string, number>();

    quakeData.points.forEach((point: Point) => {
        const { areaCode, scale } = point;
        const maxScale = areaScaleMap.get(areaCode);

        if (maxScale === undefined || scale > maxScale) {
            areaScaleMap.set(areaCode, scale);
        }
    });

    mapIsLoading.value = true;
    isLoading.value = false;

    try {
        const map = L.map("map", {
            attributionControl: false,
            zoomControl: false
        }).setView([quakeData.earthquake.hypocenter.latitude, quakeData.earthquake.hypocenter.longitude], 7);

        if (isForeignQuake) {
            L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
                attribution: "© OpenStreetMap contributors",
            }).addTo(map);
        } else {
        // GeoJSON をマップに追加
            const geoJsonLayer = L.geoJSON(geojson, {
                style: function (feature) {
                    return {
                        color: "gray",  // 境界線の色
                        weight: 2,
                        opacity: 0.7,
                        fillColor: areaScaleMap.has(feature!.properties.code) ? getQuakeScaleColor(areaScaleMap.get(feature!.properties.code) || 0) : "white",
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

            if (!dispHypocenter) {
                let maxScale = -Infinity;
                let targetCode = null;
                let targetCenter = null;

                // 一番スケールが高い地域の `code` を取得
                for (const [code, scale] of areaScaleMap.entries()) {
                    if (scale > maxScale) {
                        maxScale = scale;
                        targetCode = code;
                    }
                }

                // 対象エリアの中心を取得
                geoJsonLayer.eachLayer(layer => {
                // @ts-ignore
                    if (layer.feature.properties.code === targetCode) {
                    // @ts-ignore
                        targetCenter = layer.getBounds().getCenter();
                    }
                });

                // ズーム倍率を手動で設定
                if (targetCenter) {
                    const zoomLevel = 6.5; // 任意のズームレベルに変更
                    map.setView(targetCenter, zoomLevel);
                }
            }
        }

        if (dispHypocenter) {
            L.marker({
                lat: quakeData.earthquake.hypocenter.latitude,
                lng: quakeData.earthquake.hypocenter.longitude
            }, {
                icon: L.icon({
                    iconUrl: "/b4a90f43-d844-fd77-8f28-221910d52e3b.png",
                    iconSize: [40, 40],
                    iconAnchor: [20, 20],
                    popupAnchor: [0, -40]
                })
            }).addTo(map);
        }
    } catch(e) {
        mapError.value = true;
        console.error(`Failed to load map: ${e}`);
    } finally {
        mapIsLoading.value = false;
    }
});
</script>

<template>
    <div class="debug-warning" v-if="isDebug">
        <span>DEBUG MODE - THIS IS DUMMY DATA</span>
    </div>
    <div class="fatal-error" v-if="fatalError">
        <UndrawError />
        <p>FATAL ERROR</p>
    </div>
    <div v-else>
        <div id="map" style="width: 100%; height: 500px">
            <div class="loading" v-if="isLoading || mapIsLoading">
                <svg xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-loader-3"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M3 12a9 9 0 0 0 9 9a9 9 0 0 0 9 -9a9 9 0 0 0 -9 -9" /><path d="M17 12a5 5 0 1 0 -5 5" /></svg>
                <p v-if="isLoading">LOADING...</p>
                <p v-else-if="mapIsLoading">マップ準備中...</p>
            </div>
            <div class="map-error" v-if="mapError">
                <svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-map-x"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M14 19.5l-5 -2.5l-6 3v-13l6 -3l6 3l6 -3v9" /><path d="M9 4v13" /><path d="M15 7v6.5" /><path d="M22 22l-5 -5" /><path d="M17 22l5 -5" /></svg>
                <p>マップのロードに失敗しました</p>
            </div>
        </div>
        <div class="quake-info" v-if="!isLoading">
            <div class="scale">
                <span class="scale-label">最大震度</span>
                <QuakeScaleIcon class="scale-icon" :scale=quakeScale />
            </div>

            <div class="hypocenter">
                <span class="hypocenter-label">震源 {{ hypocenterLabel }}</span> <br>
                <span class="time-label">{{ quakeDateTimeLabel }}</span><span class="magnitude-label">M{{ magnitude
                    }}</span>
            </div>
        </div>
        <div class="quake-info" v-if="!isLoading">
            <div class="tsunami-status" :style="tsunamiStatusStyle">
                <span class="tsunami-label">{{ tsunamiStatusLabel }}</span>
            </div>
        </div>
    </div>
</template>

<style scoped>
.debug-warning {
    position: fixed;
    z-index: 9999;
    top: 0;
    right: 0;
    background: #ffffff82;
    color: red;
    margin: 4px;
}

#map {
    width: 100%;
    height: 500px;
    border-radius: 10px;

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

.fatal-error {
    width: 320px;
    padding: 42px;
    margin: 100px auto 12px auto;
    text-align: center;

    & p {
        margin-top: 32px;
        font-size: 18px;
    }
}

.quake-info {
    display: flex;
    justify-content: space-between;

    margin: 32px 0 0 12px;

    & .scale {
        display: flex;
        align-items: center;

        width: 180px;
        justify-content: space-between;

        & .scale-label {
            font-size: 24px;
        }
    }

    & .hypocenter {
        text-align: right;

        & .hypocenter-label {
            font-size: 24px;
        }

        & .time-label {
            margin-right: 24px;
        }

        & .magnitude-label {
            font-size: 20px;
            font-weight: bold;
        }
    }

    & .tsunami-status {
        padding: 10px;
        border-radius: 10px;

        & .tsunami-label {
            font-weight: bold !important;
        }
    }

}

@keyframes rotate-z {
  from { transform: rotateZ(0deg); }
  to { transform: rotateZ(360deg); }
}
</style>
