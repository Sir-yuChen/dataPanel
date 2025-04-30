import { RootStore } from '@/store';
import { App as AntdApp, ConfigProvider } from 'antd';
import zhCN from 'antd/es/locale/zh_CN';
import { Provider as MobxProvider } from "mobx-react";
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './app';
import './style.css';

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    /* 使用MobX-React提供的Provider写入Store Context */
    <MobxProvider {...RootStore}>
        <BrowserRouter >
            <ConfigProvider
                theme={{
                    token: {
                        colorPrimary: '#28a770'
                    }
                }}
                locale={zhCN}
            >
                <AntdApp>
                    <App />
                </AntdApp>
            </ConfigProvider>
        </BrowserRouter>
    </MobxProvider>
)
