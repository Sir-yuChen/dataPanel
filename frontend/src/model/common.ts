//{ loadDataType: 'default', dataSavePath: '/dataPanel', loadDataChecked: ['c'] }
export type ReqField = {
    loadDataType: string;
    dataSavePath: string;
    loadDataChecked: string[];
};
//配置文件
export type ConfigInfo = {
    id: number;
    key: string;
    value: string;
    modify: number
    name: string;
    isShow: number;
    showType: string;
    create_at: string;
    update_at: string;
    values: ConfigValues[];
    children: ConfigInfo[];
};
export type ConfigValues = {
    key: string;
    name: string;
}