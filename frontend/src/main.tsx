import {RootStore} from '@/store';
import {App as AntdApp, ConfigProvider} from 'antd';
import zhCN from 'antd/es/locale/zh_CN';
import {Provider as MobxProvider} from "mobx-react";
import {createRoot} from 'react-dom/client';
import {HashRouter} from 'react-router-dom';
import NiceModal from '@ebay/nice-modal-react';
import App from './app';
import './style.css';
import 'antd/dist/reset.css';
import React from 'react';
import ErrorBoundary from "@/components/errorBoundary";

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    /* 使用MobX-React提供的Provider写入Store Context */
    <MobxProvider {...RootStore}>
        {/*必须使用HashRouter路由,否则会由于路由切换页面找不到wails中得事件*/}
        <HashRouter>
            <ConfigProvider
                theme={{
                    token: {
                        colorPrimary: '#28a770'
                    }
                }}
                locale={zhCN}
            >
                <AntdApp>
                    <ErrorBoundary
                        fallback={<div style={{padding: 80, fontSize: 30, color: 'white'}}>⚠️
                            应用发生错误，请刷新页面</div>}
                    >
                        <NiceModal.Provider>
                            <App/>
                        </NiceModal.Provider>
                    </ErrorBoundary>
                </AntdApp>
            </ConfigProvider>
        </HashRouter>
    </MobxProvider>
)
