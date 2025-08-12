import { createRouter, createWebHistory } from 'vue-router'
import Home from './Home.vue'
// import Signup from './Signup.vue'
import Dashboard from './Dashboard.vue'

const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },

  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard
  }
  // ,
  // {
  //   path: '/signup',
  //   name: 'Signup',
  //   component: Signup
  // },

]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router 