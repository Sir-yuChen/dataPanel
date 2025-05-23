import {App} from 'antd';
import useStore from "@/hooks/useStore";
import {useState} from "react";

// 全局提示弹窗封装
export const useShowMessage = () => {
    const {message} = App.useApp();
    return (type: 'success' | 'error' | 'warning' | 'info', text: string | unknown, duration?: number) => {
        if (typeof text === 'string') {
            message[type](text, duration || 1.5);
        } else if (text && typeof text === 'object') {
            message[type](JSON.stringify(text), duration || 1.5);
        }
    };
};
export const useShowNotification = () => {
    const {notification} = App.useApp();
    const DEFAULT_DURATION = 3;
    return (type: 'success' | 'error' | 'warning' | 'info',
            text: string | unknown,
            duration?: number, title?: string) => {
        try {
            const message = typeof text === 'string' ? text :
                JSON.stringify(text, (_, value) => {
                    if (value instanceof Error) return value.message;
                    return value;
                }, 2);

            notification[type]({
                message: title ?? "提示消息",
                duration: duration ?? DEFAULT_DURATION,
                description: message,
                placement: 'bottomRight',
                showProgress: true,
                className: 'notification-notice-custom-content'
            });
        } catch (error) {
            notification.error({
                message: '提示错误',
                description: error instanceof Error ? error.message : '未知错误'
            });
        }
    };
};