// 新建 src/hooks/useGlobalListeners.ts
import {useEffect, useCallback} from 'react';
import {EventsOn} from '@/wailsjs/runtime';
import NiceModal from '@ebay/nice-modal-react';
import {LoadDataModalComponent} from '@/components/modal/loadDataModal';
import {useShowMessage, useShowNotification} from '@/hooks/useShowMessage';

interface MessageEvent {
    dialogType: 'success' | 'error' | 'info' | 'warning';
    content: string;
    duration?: number;
}

const isMessageEvent = (data: unknown): data is MessageEvent => {
    return !!data && typeof data === 'object' &&
        'dialogType' in data &&
        ['success', 'error', 'info', 'warning'].includes((data as any).dialogType) &&
        'content' in data;
};

export const useGlobalListeners = () => {
    const showNotification = useShowNotification();
    const showMessage = useShowMessage();

    const handleLoadData = useCallback((data: unknown) => {
        try {
            NiceModal.show(LoadDataModalComponent);
        } catch (error) {
            showMessage('error', `Failed to open data loader: ${error instanceof Error ? error.message : 'Unknown error'}`);
        }
    }, [showMessage]);

    const handleMessageDialog = useCallback((data: unknown) => {
        try {
            if (isMessageEvent(data)) {
                showNotification(data.dialogType, data.content, data.duration);
            }
        } catch (error) {
            showMessage('error', `Failed to handle message dialog: ${error instanceof Error ? error.message : 'Unknown error'}`);
        }
    }, [showNotification, showMessage]);

    useEffect(() => {
        const modalListener = EventsOn("loadData", handleLoadData);
        const messageDialogsListener = EventsOn("messageDialogs", handleMessageDialog);

        return () => {
            modalListener?.();
            messageDialogsListener?.();
        };
    }, [handleLoadData, handleMessageDialog]);
};
