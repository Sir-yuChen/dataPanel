import { useState } from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import { GetHello } from '../wailsjs/go/exposed/HelloWails';
import { EventsOn } from '../runtime/runtime';

function App() {
    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [name, setName] = useState('');
    const updateName = (e: any) => setName(e.target.value);

    EventsOn("showSearch", () => {
        setName("被触发");
    });
    const greet = async () => {
        GetHello().then((result) => {
            setName(result);
        });
    }

    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo" />
            <div id="result" className="result">{resultText + name}</div>
            <div id="result" className="result">这是一个数据面板----测试</div>
            <div id="input" className="input-box">
                <button className="btn" onClick={greet}>Hello</button>
            </div>
        </div>
    )
}

export default App

