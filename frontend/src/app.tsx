import FloatButtonsComponent from '@/components/floatButton';
import IconFont from '@/components/iconFont';
import useStore from '@/hooks/useStore';
import RouteRender from '@/router';
import {AntDesignOutlined} from '@ant-design/icons';
import {Affix, Avatar, Layout, Menu, Spin, theme} from 'antd';
import {To, useNavigate} from 'react-router-dom';
import {useEffect} from "react";
import NiceModal from "@ebay/nice-modal-react";
import {LoadDataModalComponent} from "@/components/modal/loadDataModal";
import {EventsOn} from "@/wailsjs/runtime";
import {useShowNotification} from "@/hooks/useShowMessage";
import {useGlobalListeners} from "@/hooks/useGlobalListeners";


const {Content, Footer, Sider} = Layout;

const App = () => {
    useGlobalListeners();
    const globalStore = useStore("globalStore");
    const navigate = useNavigate();
    const {token: {colorBgContainer, borderRadiusLG}} = theme.useToken();
    // 菜单点击处理
    const menuOnClick = (e: { key: To }) => navigate(e.key);
    return (
        <div style={{height: '100vh'}}>
            <Layout style={{minHeight: '100%'}}>
                <Sider collapsed={true} width={160}>
                    <Affix offsetTop={10}>
                        <div className="avatar_vertical">
                            <Avatar size={60} icon={<AntDesignOutlined/>}/>
                        </div>
                        <Menu
                            theme="dark"
                            mode="inline"
                            defaultSelectedKeys={['']}
                            onClick={menuOnClick}
                            items={[
                                {key: '/home', icon: <IconFont type="icon-shouye"/>, label: '首页'},
                                {key: '/my', icon: <IconFont type="icon-wode"/>, label: '我的'},
                                {key: '/setting', icon: <IconFont type="icon-shezhi"/>, label: '设置'},
                            ]}
                        />
                    </Affix>
                </Sider>
                <Layout>
                    <Spin
                        tip={globalStore.getIsDataLoadedText()}
                        delay={300}
                        spinning={!globalStore.getIsDataLoaded()}
                    >
                        <Content style={{margin: '10px 16px'}}>
                            <div
                                style={{
                                    padding: 20,
                                    background: colorBgContainer,
                                    borderRadius: borderRadiusLG,
                                }}
                            >
                                <FloatButtonsComponent/>
                                <RouteRender/>
                            </div>
                        </Content>
                    </Spin>
                    <Footer style={{textAlign: 'center'}}>
                        Ant Design ©{new Date().getFullYear()} Created by Ant UED
                    </Footer>
                </Layout>
            </Layout>
        </div>
    );
};

export default App;
