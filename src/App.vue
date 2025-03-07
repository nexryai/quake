<script setup lang="ts">
import { RouterLink, RouterView } from "vue-router";
</script>

<template>
    <header>
        <nav>
            <span class="nav-title">Quake</span>
            <RouterLink class="item" exact-active-class="active" to="/">Home</RouterLink>
            <RouterLink class="item" exact-active-class="active" to="/quake">地震一覧</RouterLink>
            <RouterLink class="item" exact-active-class="active" to="/tsunami">津波情報</RouterLink>
        </nav>
        <div class="header-right">
            <RouterLink class="item" exact-active-class="active" to="/about">About</RouterLink>
        </div>
    </header>
    <RouterView id="view" />
</template>

<style scoped>
header {
    display: flex;
    justify-content: space-between;
    position: fixed;
    z-index: 1000;
    top: 0;
    left: 0;
    width: 100vw;
    background: #fbfbfbb5;
    backdrop-filter: blur(18px);

    & nav {
        width: 100%;
        font-size: 13px;
        margin: 1rem 0 1rem 1rem;

        & .nav-title {
            font-size: 16px;
            font-weight: bold;
            margin-right: 32px;
        }

        & .item {
            margin-right: 8px;
            padding: 12px;
        }

        & .active {
            color: black;
            font-weight: bold;
            border-bottom: solid black 1px;
        }
    }

    & .header-right {
        margin: 1rem;

        & .active {
            color: black;
            font-weight: bold;
        }
    }
}

#view {
    margin-top: 52px;
    width: 100%;
    height: 100%;
}
</style>
