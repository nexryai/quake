<script setup lang="ts">
import L from "leaflet";
import { onMounted } from "vue";

import "leaflet/dist/leaflet.css";
import {getQuakeScaleColor} from "@/color.ts";

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

onMounted(async () => {
    const map = L.map("map", {
        attributionControl: false,
        zoomControl: false
    }).setView([35.6895, 139.6917], 10); // 東京を中心に表示

    // タイルレイヤーを追加
    /*
    L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
        attribution: "© OpenStreetMap contributors",
    }).addTo(map.value);
    */
    // GeoJSON データ
    const geojson = await fetchMapData();

    const quakeData = await fetchQuakeData("20250304061539_0_VXSE53_010000");

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
});
</script>

<template>
    <div id="map" style="width: 100%; height: 500px"></div>
</template>

<style scoped>
#map {
    width: 100%;
    height: 500px;
}
</style>
