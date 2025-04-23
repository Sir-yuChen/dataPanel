import { GetHello } from '@/wailsjs/go/exposed/HelloWails';
import { EventsOn } from '@/wailsjs/runtime/runtime';
import { useState } from 'react';
import styles from './index.module.less';

const Home = () => {
    //定义展示变量
    const [resText, setResText] = useState('');
    const geet = () => {
        GetHello().then((res) => {
            setResText(res);
        })
    }
    //后端回调函数EventsEmit
    EventsOn('showSearch', () => {
        setResText("回调成功");
    })

    return (
        <>
            <div className={styles.bgContainer}>
                <div>
                    你好，欢迎使用
                </div>
                <div>{resText}</div>
                <button onClick={geet}>点击触发</button>
            </div>
        </>
    );
};

export default Home;
