<script setup lang="ts">
import L from "leaflet";
import { ref, onMounted } from "vue";

import "leaflet/dist/leaflet.css";
import {getQuakeScaleColor} from "@/color.ts";
import QuakeScaleIcon from "@/components/core/QuakeScaleIcon.vue";

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

const fetchMapData = async () => {
    const response = await fetch("/JP20250304.geojson");
    return await response.json();
};

const fetchQuakeData = async (id: string): Promise<EarthquakeData> => {
    const response = await fetch(`https://quake-jade.vercel.app/api/events/details?id=${id}`);
    return await response.json() as EarthquakeData;
};

const quakeScale = ref(0);
const magnitude = ref(0);
const hypocenterLabel = ref("不明");

onMounted(async () => {
    // 東京を中心に表示

    // タイルレイヤーを追加
    /*
    L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
        attribution: "© OpenStreetMap contributors",
    }).addTo(map.value);
    */
    // GeoJSON データ
    const geojson = await fetchMapData();

    const quakeData = await fetchQuakeData("20250304061539_0_VXSE53_010000");

    quakeScale.value = quakeData.earthquake.maxScale;
    magnitude.value = quakeData.earthquake.hypocenter.magnitude;
    hypocenterLabel.value = quakeData.earthquake.hypocenter.name;

    // areaCodeごとの最大scaleを求める
    const areaScaleMap = new Map<string, number>();

    quakeData.points.forEach((point: Point) => {
        const { areaCode, scale } = point;
        const maxScale = areaScaleMap.get(areaCode);

        if (maxScale === undefined || scale > maxScale) {
            areaScaleMap.set(areaCode, scale);
        }
    });

    console.log(areaScaleMap);

    const map = L.map("map", {
        attributionControl: false,
        zoomControl: false
    }).setView([quakeData.earthquake.hypocenter.latitude, quakeData.earthquake.hypocenter.longitude], 7);

    // GeoJSON をマップに追加
    L.geoJSON(geojson, {
        style: function(feature) {
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

    L.marker({
        lat: quakeData.earthquake.hypocenter.latitude,
        lng: quakeData.earthquake.hypocenter.longitude
    }, { icon: L.icon({
        iconUrl: "/b4a90f43-d844-fd77-8f28-221910d52e3b.png",
        iconSize: [40, 40],
        iconAnchor: [20, 20],
        popupAnchor: [0, -40]
    })}).addTo(map);
});
</script>

<template>
    <div id="map" style="width: 100%; height: 500px"></div>
    <div class="quake-info">
        <div class="scale">
            <span class="scale-label">最大震度</span>
            <QuakeScaleIcon class="scale-icon" :scale=quakeScale />
        </div>

        <div class="hypocenter">
            <span class="hypocenter-label">震源 {{ hypocenterLabel }}</span> <br>
            <span class="magnitude-label">M{{ magnitude }}</span>
        </div>
    </div>
</template>

<style scoped>
#map {
    width: 100%;
    height: 500px;
    border-radius: 10px;
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

        & .magnitude-label {
            font-size: 20px;
        }
    }

}
</style>
