import {makeAutoObservable} from 'mobx';
import {ConfigInfo} from "@/model/common";

export class ConfigStore {
    configInfo: ConfigInfo[] = [];

    setUserModel(config: ConfigInfo[]) {
        this.configInfo = config
    }

    constructor() {
        makeAutoObservable(this);
    }

    getConfigInfo(): ConfigInfo[] {
        return this.configInfo;
    }

    setConfigInfo(config: ConfigInfo[]) {
        this.configInfo = config
    }

    //新增修改指定key得value值 注意：
    setConfigValue(key: string, value: string) {
        const updateNestedConfig = (items: ConfigInfo[]) => {
            items.forEach((item) => {
                // 当前层级匹配直接修改
                if (item.key === key) {
                    item.value = value;
                }
                // 递归处理子节点（支持无限嵌套）
                if (item.children && item.children.length > 0) {
                    updateNestedConfig(item.children);
                }
            });
        };
        // 从根节点开始递归处理
        updateNestedConfig(this.configInfo);
    }

    getConfigValue(key: string): string | undefined {
        const getNestedConfig = (items: ConfigInfo[]): string | undefined => {
            if (!items) return undefined; // 处理空输入

            for (const item of items) { // 使用for...of支持提前返回
                if (item.key === key) {
                    return item.value; // 直接返回匹配值
                }
                if (item.children?.length) { // 可选链操作符处理空值
                    const childResult = getNestedConfig(item.children);
                    if (childResult !== undefined) { // 处理子级返回值
                        return childResult;
                    }
                }
            }
            return undefined; // 明确返回未找到
        };

        return getNestedConfig(this.configInfo);
    }

}
