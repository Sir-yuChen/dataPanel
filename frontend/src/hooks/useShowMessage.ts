import { App } from 'antd';

// 全局提示弹窗封装
export const useShowMessage = () => {
    const { message } = App.useApp();
    return (type: 'success' | 'error' | 'warning' | 'info', text: string | unknown, duration?: number) => {
        if (typeof text === 'string') {
            message[type](text, duration || 1.5);
        } else if (text && typeof text === 'object') {
            message[type](JSON.stringify(text), duration || 1.5);
        }
    };
};

//数据加载全局弹窗
export const useShowLoadDataModal = () => {
    const { modal } = App.useApp();
    /**
   * 显示数据加载中的全局弹窗
   * @param content 弹窗内容（支持字符串或 React 节点）
   * @param onOk 确认回调
   * @param onCancel 取消回调
   */
    const showLoadDataModal = (props: {
        title: string;
        content: React.ReactNode;
        okText?: string;
        cancelText?: string;
        width?: number;
        isFooter?: boolean;
        onOk?: () => void;
        onCancel?: () => void;
    }) => {
        const { title, content, onOk, onCancel, okText, cancelText, width, isFooter } = props;

        modal.info({
            title: title,
            content: content,
            width: width ?? 500,
            onOk: () => {
                if (onOk) {
                    onOk();
                }
            },
            onCancel: () => {
                if (onCancel) {
                    onCancel();
                }
            },
            closable: true,
            maskClosable: false,
            centered: true, // 居中显示弹窗
            footer: null,
        });
    };

    return {
        showLoadDataModal,
        modal
    };
};

