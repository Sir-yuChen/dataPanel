import react from '@vitejs/plugin-react';
import path from 'path';
import {defineConfig} from 'vite';

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [react()],
    server: {
        // port: 3000,
        open: true,
        // proxy: { '/api': 'http://localhost:8089' }
    },
    resolve: {
        alias: {
            '@': path.resolve(__dirname, './src'),
            '@components': path.resolve(__dirname, './src/components')
        }
    },
    css: {
        preprocessorOptions: {
            less: {
                javascriptEnabled: true // 启用 Less 的 JS 解析
            }
        }
    }
})
