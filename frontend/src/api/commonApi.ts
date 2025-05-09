import { ReqField } from "@/model/loadDtaModel";
import Request from '@/utils/request';

export const LoadDataModalSubmit = async (data: ReqField) => {
    const res = await Request.post<ReqField, any>(`/hello`, data);
    console.log("调用后端返回参数：", res);
    return res.data;
};