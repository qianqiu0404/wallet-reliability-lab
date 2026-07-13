import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: () => import('./pages/HomePage.vue') },
    { path: '/lab/:scenarioId', name: 'lab', component: () => import('./pages/LabPage.vue') },
    { path: '/architecture', name: 'architecture', component: () => import('./pages/ArchitecturePage.vue') },
    { path: '/evidence', name: 'evidence', component: () => import('./pages/EvidencePage.vue') },
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
  scrollBehavior: () => ({ top: 0 }),
})

export default router

