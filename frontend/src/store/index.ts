import {GlobalStore} from "@/store/globalStore";
import {UserStore} from "@/store/userStore";
import {ConfigStore} from "@/store/config";

/** 将每个Store实例化 */
export const RootStore = {
    globalStore: new GlobalStore(),
    userStore: new UserStore(),
    configStore: new ConfigStore(),
}
