<script setup lang="ts">
import L from "leaflet";
import { onMounted, ref } from "vue";
import "leaflet/dist/leaflet.css";

const map = ref(null);

const fetchMapData = async () => {
    const response = await fetch("/JP20250304.geojson");

    return  await response.json();
};

onMounted(async () => {
    map.value = L.map("map", {
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

    // GeoJSON をマップに追加
    L.geoJSON(geojson, {
        style: function(feature) {
            return {
                color: "gray",  // 境界線の色
                weight: 2,
                opacity: 0.7,
                fillColor: feature.properties.code === "360" ? "red" : "white",
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
    }).addTo(map.value);
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
