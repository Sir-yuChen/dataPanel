import RouteRender from '@/router';
import '@/style.css';
import { App, ConfigProvider } from 'antd';
import zhCN from 'antd/es/locale/zh_CN';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <BrowserRouter>
        <ConfigProvider
            theme={{
                token: {
                    colorPrimary: '#28a770'
                }
            }}
            locale={zhCN}
        >
            <App>
                <RouteRender />
            </App>
        </ConfigProvider>
    </BrowserRouter>
)
