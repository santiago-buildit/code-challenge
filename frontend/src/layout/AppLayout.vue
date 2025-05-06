<template>
  <div class="app-layout">
    <AppHeader title="Library Management System" @toggle-sidebar="showSidebar = !showSidebar" />
    <div class="layout-body">
      <AppSidebar :visible="showSidebar" @close="showSidebar = false" />
      <main class="main-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script lang="ts">
import AppHeader from './AppHeader.vue'
import AppSidebar from './AppSidebar.vue'

export default {
  components: { AppHeader, AppSidebar },
  data() {
    return {
      showSidebar: false,
    }
  },
}
</script>

<style scoped>
.app-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
  /* Prevents scroll on the entire layout */
}

.layout-body {
  display: flex;
  flex: 1;
  overflow: hidden;
  /* Prevents lateral scroll if sidebar overflows */
}

.main-content {
  flex: 1;
  padding: 1rem;
  overflow-y: auto;
  /* Allows scroll here */
  overflow-x: hidden;
}
</style>