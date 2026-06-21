// admin/src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '../pages/dashboard.vue'
import Databases from '../pages/database.vue'
import Sites from '../pages/sites.vue'
import Storage from '../pages/storage.vue'
import Docs from '../pages/docs.vue'

const routes = [
  {
    path: '/admin',
    name: 'Dashboard',
    component: Dashboard
  },
  {
    path: '/admin/sites',
    name: 'Sites',
    component: Sites
  },
  {
    path: '/admin/databases',
    name: 'Databases',
    component: Databases
  },
  {
    path: '/admin/storage',
    name: 'Storage',
    component: Storage
  },
  {
    path: '/admin/docs',
    name: 'Docs',
    component: Docs
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
