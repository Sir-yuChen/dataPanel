import FloatButtonsComponent from '@/components/floatButton';
import IconFont from '@/components/iconFont';
import { LoadDataModalComponent } from '@/components/modal/loadDataModal';
import { useShowLoadDataModal, useShowMessage } from '@/hooks/useShowMessage';
import useStore from '@/hooks/useStore';
import RouteRender from '@/router';
import { EventsOn } from '@/wailsjs/runtime';
import {
    AntDesignOutlined
} from '@ant-design/icons';
import { Affix, Avatar, Layout, Menu, Spin, theme } from 'antd';
import { To, useNavigate } from 'react-router-dom';

const { Header, Content, Footer, Sider } = Layout;
const App = () => {
    const showMessage = useShowMessage()
    //校验页面数据是否加载完成
    const globalStore = useStore("globalStore")
    //路由跳转
    const navigate = useNavigate()

    const {
        token: { colorBgContainer, borderRadiusLG },
    } = theme.useToken();

    // 菜单跳转切换页面
    const menuOnClick = (e: { key: To; }) => {
        navigate(e.key)
    }
    //全局监听器
    EventsOn("messageDialogs", (re) => {
        if (re.dialogType) {
            if (re.duration && re.duration > 0) {
                showMessage(re.dialogType, re.content, re.duration)
            } else {
                showMessage(re.dialogType, re.content)
            }
        }
    })
    const showLoadData = useShowLoadDataModal()
    EventsOn("loadData", (re) => {
        showLoadData.showLoadDataModal({
            title: "加载数据配置",
            width: 600,
            content: <LoadDataModalComponent />,
            isFooter: true
        })
    })

    return (
        <div style={{ height: '100vh' }}>
            <Layout style={{ minHeight: '100%' }}>
                <Sider collapsed={true} width={160}>
                    <Affix offsetTop={10}>
                        <div className="avatar_vertical"  >
                            {/* 加载头像  默认为logo*/}
                            <Avatar size={60} icon={<AntDesignOutlined />} />
                        </div>
                        <Menu
                            theme="dark"
                            mode="inline"
                            defaultSelectedKeys={['']}
                            onClick={menuOnClick}
                            items={[
                                {
                                    key: '/home',
                                    icon: <IconFont type="icon-shouye" />,
                                    label: '首页',
                                },
                                {
                                    key: '/my',
                                    icon: <IconFont type="icon-wode" />,
                                    label: '我的',
                                },
                                {
                                    key: '/setting',
                                    icon: <IconFont type="icon-shezhi" />,
                                    label: '设置',
                                },
                            ]}
                        />
                    </Affix>
                </Sider>
                <Layout >
                    <Spin tip={globalStore.getIsDataLoadedText()} spinning={!globalStore.getIsDataLoaded()}>
                        <Content style={{ margin: '10px 16px' }}>
                            <div
                                style={{
                                    padding: 20,
                                    background: colorBgContainer,
                                    borderRadius: borderRadiusLG,
                                }}
                            >
                                {/*浮动按钮 */}
                                <FloatButtonsComponent />
                                {/* 路由出口 */}
                                <RouteRender />
                            </div>
                        </Content>
                    </Spin>
                    <Footer style={{ textAlign: 'center' }}>
                        Ant Design ©{new Date().getFullYear()} Created by Ant UED
                    </Footer>
                </Layout>
            </Layout>
        </div>
    );
};

export default App;
