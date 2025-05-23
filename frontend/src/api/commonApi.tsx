import {ConfigInfo, ReqField} from "@/model/common";
import request from "@/utils/request";

export const LoadDataModalSubmit = async (data: ReqField) => {
    const res = await request.post<any>(`/common/loadData`, data);
    return res.data;
};
export const GetSetting = async () => {
    const res = await request.get<ConfigInfo[]>(`/common/appConfig`);
    return res.data;
};
export const UpdateSetting = async (data: ConfigInfo[])=>{
    const res = await request.post<any>(`/common/updateAppConfig`, data);
    return res.data;
}