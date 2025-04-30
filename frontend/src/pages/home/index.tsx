import useStore from '@/hooks/useStore';
import { Input, Select } from 'antd';
import { useState } from 'react';
import './index.css';

const { Search } = Input;
const { Option } = Select;
// 首页 搜索
const Home = () => {
    //校验页面数据是否加载完成
    const globalStore = useStore("globalStore")
    //搜索框加载按钮状态 loadingFlag
    const [loading, setLoading] = useState(false);
    //搜索结果集合
    const [searchResult, setSearchResult] = useState([]);

    const selectSearchBefore = (
        <Select defaultValue="name">
            <Option value="stockName">名称</Option>
            <Option value="stockCode">代码</Option>
        </Select>
    );
    //搜索框，请求后端接口
    const onSearch = (value: string) => {
        setLoading(true);
        // getSearchStockList(value).then((res) => {
        //     setLoading(false);
        //     searchResult.push(res);
        // });
    };


    return (
        <>
            <div className="home_all">
                {/* 搜索框 居中 */}
                <div className="home_search">
                    <Search addonBefore={selectSearchBefore} placeholder="请输入搜索关键字" enterButton="搜索" size="large" loading={loading} onSearch={onSearch} />
                </div>
                {/* 搜索结果展示 */}
            </div>
        </>
    );
};

export default Home;
