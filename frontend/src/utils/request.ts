/* eslint-disable @typescript-eslint/no-explicit-any */
import {apiBaseURL} from '@/config';
import axios, {AxiosError, AxiosInstance, AxiosRequestConfig, AxiosResponse, InternalAxiosRequestConfig} from 'axios';
import {message} from "antd";

interface ApiResponse<T> {
    code: number;
    msg: string;
    data: T;
}

const createRequest = (config?: AxiosRequestConfig) => {
    const defaultConfig: AxiosRequestConfig = {
        baseURL: apiBaseURL,
        timeout: 6000
    };

    const instance: AxiosInstance = axios.create({
        ...defaultConfig,
        ...config
    });

    // 请求拦截器
    instance.interceptors.request.use(
        (config: InternalAxiosRequestConfig) => {
            const token = getToken();
            if (token) {
                config.headers.Authorization = token;
            }
            return config;
        },
        (error: AxiosError) => Promise.reject(error)
    );

    // 响应拦截器
    instance.interceptors.response.use(
        (response: AxiosResponse) => {
            if (response.data.code !== 200) {
                message.error(response.data.msg)
            }
            message.success(response.data.msg)
            return response;
        },
        (error: AxiosError) => {
            const apiError = error.response?.data;
            return Promise.reject(apiError);
        }
    );

    // Token 获取函数
    const getToken = (): string | null => {
        return JSON.parse(sessionStorage.getItem('authToken') || 'null');
    };

    // HTTP 方法
    const request = (config: AxiosRequestConfig): Promise<any> => {
        return instance.request(config);
    };

    const get = <TResponse = any>(
        url: string,
        config?: AxiosRequestConfig
    ): Promise<AxiosResponse<ApiResponse<TResponse>>> => {
        return instance.get(url, config);
    };

    const post = <TRequest = any, TResponse = any>(
        url: string,
        data?: TRequest,
        config?: AxiosRequestConfig
    ): Promise<AxiosResponse<ApiResponse<TResponse>>> => {
        return instance.post(url, data, config);
    };

    const put = <TRequest = any, TResponse = any>(
        url: string,
        data?: TRequest,
        config?: AxiosRequestConfig
    ): Promise<AxiosResponse<ApiResponse<TResponse>>> => {
        return instance.put(url, data, config);
    };

    const del = <TResponse = any>(
        url: string,
        config?: AxiosRequestConfig
    ): Promise<AxiosResponse<ApiResponse<TResponse>>> => {
        return instance.delete(url, config);
    };

    return {
        request,
        get,
        post,
        put,
        delete: del
    };
};

export default createRequest({});
