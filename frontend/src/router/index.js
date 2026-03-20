import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  { path: '/', component: () => import('../views/Dashboard.vue') },
  { path: '/members', component: () => import('../views/MemberList.vue') },
  { path: '/members/create', component: () => import('../views/MemberForm.vue') },
  { path: '/members/:id/edit', component: () => import('../views/MemberForm.vue') },
  { path: '/clauses', component: () => import('../views/ClauseList.vue') },
  { path: '/clauses/create', component: () => import('../views/ClauseForm.vue') },
  { path: '/clauses/:id/edit', component: () => import('../views/ClauseForm.vue') },
  { path: '/templates', component: () => import('../views/TemplateList.vue') },
  { path: '/cases', component: () => import('../views/CaseList.vue') },
  { path: '/cases/create', component: () => import('../views/CaseCreate.vue') },
  { path: '/cases/:id', component: () => import('../views/CaseDetail.vue') },
  { path: '/settings', component: () => import('../views/Settings.vue') },
]

export default createRouter({
  history: createWebHashHistory(),
  routes,
})
