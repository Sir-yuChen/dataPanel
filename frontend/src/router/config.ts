import Home from '@/pages/home';
import { lazy } from 'react';

export interface IRouter {
    name?: string;
    redirect?: string;
    path: string;
    children?: Array<IRouter>;
    component: React.ComponentType;
}

export const router: Array<IRouter> = [
    //需要登录才能访问的页面
    // {
    //     path: '/',
    //     component: withPrivateRoute(lazy(() => import('@/pages/home'))), 
    //     children: [
    //         {
    //             path: 'chat',
    //             component: withPrivateRoute(lazy(() => import('@/pages/home')))
    //         },
    //         {
    //             path: 'address-book',
    //             component: withPrivateRoute(lazy(() => import('@/pages/home')))
    //         }
    //     ]
    // },

    {
        path: '/',
        component: Home,
    },
    //最低优先级，路由匹配，404
    {
        path: '*',
        component: lazy(() => import('@/pages/error')),
        redirect: '/'
    }
];
