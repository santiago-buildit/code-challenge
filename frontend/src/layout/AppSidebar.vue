<template>
  <aside :class="['sidebar', { 'sidebar-hidden': !visible }]">
    <nav>
      <ul>
        <li><router-link to="/" @click="$emit('close')">Books</router-link></li>
        <li><router-link to="/book/new" @click="$emit('close')">Add Book</router-link></li>
      </ul>
    </nav>
  </aside>
</template>

<script lang="ts">
export default {
  props: {
    visible: {
      type: Boolean,
      default: true,
    },
  },
}
</script>

<style scoped>
.sidebar {
  width: 200px;
  background-color: #f4f4f4;
  padding-top: 2rem;
  transition: transform 0.3s ease;
}

.sidebar ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.sidebar li {
  margin: 0;
}

a {
  display: block;
  width: 100%;
  padding: 0.75rem 1rem;
  color: #333;
  text-decoration: none;
}

a:hover {
  background-color: #ddd;
}

a.active {
  background-color: #ccc;
  font-weight: bold;
}

.close-button {
  display: none;
  background: none;
  border: none;
  font-size: 1.2rem;
  margin-top: 2rem;
}

/* Only on mobile can the sidebar be hidden */
@media (max-width: 768px) {
  .sidebar {
    position: absolute;
    height: 100vh;
    top: 5rem;
    left: 0;
    z-index: 1000;
    box-shadow: 2px 0 5px rgba(0, 0, 0, 0.1);
    transform: translateX(0);
  }

  .sidebar-hidden {
    transform: translateX(-100%);
  }

  .close-button {
    display: block;
  }
}
</style>
