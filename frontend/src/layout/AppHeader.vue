<template>
  <header class="header">
    <button class="menu-button" @click="$emit('toggle-sidebar')">☰</button>
    <h1>{{ title }}</h1>
  </header>
</template>

<script lang="ts">
export default {
  props: {
    title: {
      type: String,
      required: true
    }
  }
}
</script>

<style scoped>
.header {
  width: 100%;
  height: 4.5rem;
  background-color: #640032;
  display: flex;
  align-items: center;
  padding: 0 1rem;
}

h1 {
  color: white;
  margin: 0 auto;
  font-size: 1.5rem;
}

.menu-button {
  font-size: 1.5rem;
  background: none;
  border: none;
  color: white;
  cursor: pointer;
  display: none;
}

@media (max-width: 768px) {
  .menu-button {
    display: block;
  }

  h1 {
    color: white;
    margin: 0 auto;
    font-size: 1.2rem;
  }
}
</style>