<script setup lang="ts">
import { RouterLink, RouterView } from "vue-router";
</script>

<template>
    <header>
        <nav>
            <span class="nav-title">Quake</span>
            <RouterLink to="/">Home</RouterLink>
            <RouterLink to="/about">About</RouterLink>
        </nav>
        <div class="header-right">
            <span>GitHub</span>
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

  & .header-right {
      margin: 1rem;
  }

  & nav {
    width: 100%;
    font-size: 12px;
    margin: 1rem 0 1rem 1rem;

    & .nav-title {
      font-size: 16px;
      font-weight: bold;
      margin-right: 32px;
    }
  }
}

#view {
  margin-top: 52px;
  width: 100%;
  height: 100%;
}
</style>
