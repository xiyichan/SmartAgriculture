import Vue from 'vue'
import VueRouter from 'vue-router'
import Login from '../components/Login.vue'
import HomePage from '../components/HomePage.vue'
import Home from '../components/Home.vue'
import UserList from '../components/UserManagement/UserList.vue'
import AdminList from "../components/AdminList.vue";
import Register from "../components/Register.vue"
import Pi from "../components/Pi.vue"
Vue.use(VueRouter)

const routes = [{
        path: '/',
        name: 'Login',
        component: Login
    },
    {
        path: "/HomePage",
        name: "HomePage",
        component: HomePage,
        redirect: "/Home",
        children: [{
                path: "/Home",
                name: "Home",
                component: Home,
            },

            {
                path: "/UserList",
                name: "UserList",
                component: UserList,
            },
            {
                path: "/Pi",
                name: "Pi",
                component: Pi
            },
            {
                path: "/AdminList",
                name: "AdminList",
                component: AdminList,
            },
            {
                path: "/Register",
                name: "Register",
                component: Register,
            },
        ]
    }
]

const router = new VueRouter({
    routes
})
router.beforeEach((to, from, next) => {
    let token = localStorage.getItem("token");
    if (to.path === "/") {
        next();
    } else if (!token) {
        next("/");
    } else {
        next();
    }
});
export default router