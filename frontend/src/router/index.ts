import { createRouter, createWebHistory } from 'vue-router'
import BookListPage from '../pages/BookListPage.vue'
import BookFormPage from '../pages/BookFormPage.vue'

// Define app routes
const routes = [
  { path: '/', component: BookListPage },
  { path: '/book/new', component: BookFormPage },
  { path: '/book/:id', component: BookFormPage },
]

// Configure router
const router = createRouter({
  history: createWebHistory(),
  routes,
  linkActiveClass: 'active'
})

export default router
