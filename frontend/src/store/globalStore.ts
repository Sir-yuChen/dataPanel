import { makeAutoObservable } from "mobx";

// 全局状态管理

export class GlobalStore {
    token: string = '';
    loginStatus: string = 'N';
    //数据是否加载完成标识
    isDataLoaded: boolean = true;
    //数据加载进度文本
    isDataLoadedText: string = '努力加载数据中,请稍等...';

    constructor() {
        makeAutoObservable(this);
        // makePersistable 数据持久化存储
        /* makePersistable(this, {  // 在构造函数内使用 makePersistable
             name:'userToken', // 保存的name，用于在storage中的名称标识，只要不和storage中其他名称重复就可以
             properties: ["token"], // 要保存的字段，这些字段会被保存在name对应的storage中，注意：不写在这里面的字段将不会被保存，刷新页面也将丢失：get字段例外。get数据会在数据返回后再自动计算
             storage: window.localStorage, // 保存的位置：看自己的业务情况选择，可以是localStorage，sessionstorage 常见的就是localStorage
             //还有一些其他配置参数，例如数据过期时间等等像storage这种字段可以配置在全局配置里
         }).then(
             action((persistStore) => { // persist 完成的回调，在这里可以执行一些拿到数据后需要执行的操作，如果在页面上要判断是否完成persist，使用 isHydrated
                 console.log("持久化后的数据",persistStore);
             }))*/
    }

    clear() {
        this.token = "";
        this.loginStatus = "N";
    }

    setIsDataLoadedText(text: string) {
        this.isDataLoadedText = text;
    }
    getIsDataLoadedText() {
        return this.isDataLoadedText;
    }

    setIsDataLoaded(falg: boolean) {
        this.isDataLoaded = falg;
    }
    getIsDataLoaded() {
        return this.isDataLoaded;
    }

    setToken(token: string) {
        this.token = token;
    }

    getToken() {
        return this.token;
    }

    setLoginStatus(loginStatus: string) {
        this.loginStatus = loginStatus;
    }

    getLoginStatus() {
        return this.loginStatus;
    }

}
