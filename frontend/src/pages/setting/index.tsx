import {useEffect, useState} from "react";
import {GetSetting, UpdateSetting} from "@/api/commonApi";
import {
    Button,
    Checkbox,
    ColorPicker,
    Divider,
    Image,
    Input,
    Radio,
    RadioChangeEvent,
    Select,
    Space,
    Spin,
    Switch,
    Tabs
} from "antd";
import {UploadOutlined} from '@ant-design/icons';
import {ConfigInfo} from "@/model/common";
import {useShowMessage} from "@/hooks/useShowMessage";
import TextArea from "antd/es/input/TextArea";
import useStore from "@/hooks/useStore";
import {observer} from "mobx-react";

const Setting = () => {
    const configStore = useStore("configStore")
    const [colorPickerOpen, setColorPickerOpen] = useState(false);
    const [loading, setLoading] = useState(true);
    const [backgroundImgValue, setBackgroundImgValue] = useState("");
    //是否展示保存按钮
    const [showSaveButton, setShowSaveButton] = useState(false);

    let showMessage = useShowMessage();
    //请求后端接口 获取设置信息
    useEffect(() => {
        GetSetting().then((res) => {
            if (res.code == 200) {
                configStore.setConfigInfo(res.data)
                console.log("configstore:", res.data)
            }
        }).finally(() => {
            setLoading(false); // 无论成功失败都关闭loading
        });
    }, [])
    //根据配置文件中得类型，选择操作组件
    const getComponent = (config: ConfigInfo) => {
        //colorPicker取色器 switch单选按钮 checkbox单选框 checkboxs多选框 input输入框 selects多选下拉框
        const render = () => {
            switch (config.showType) {
                case "colorPicker":
                    return renderColorPicker(config);
                case "switch":
                    return renderSwitch(config);
                case "checkboxs":
                    return renderCheckboxs(config);
                case "input":
                    return renderInput(config);
                case "selects":
                    return renderSelects(config);
                case "rb":
                    return renderRadio(config);
                case "image":
                    return renderImage(config);
            }
        }
        const renderChildren = () => {
            if (!config.children?.length) return null;
            return config.children.map((child, index) => (
                <div> {getComponent(child)}</div>
            ));
        };
        //如果config.children数组不为空，并且config.showType为空或不存在 则返回子组件
        if (config.children?.length && (!config.showType || config.showType == '')) {
            return (
                <Space direction="horizontal" split={<Divider type="vertical"/>}>
                    <div>{config.name}</div>
                    <Space direction="vertical" align={"baseline"} size={10}
                           style={{border: "1px solid #d9d9d9", padding: '10px'}}>
                        {renderChildren()}
                    </Space>
                </Space>
            )
        }
        return (
            <Space direction="vertical" size={15} align="start">
                {render()}
                {renderChildren()}
            </Space>
        );
    };
    const renderColorPicker = (config: ConfigInfo) => {
        const handleChange = (color: any) => {
            config.value = color.toHexString();
            setShowSaveButton(true);
        };
        return (
            <Space direction="horizontal" split={<Divider type="vertical"/>}>
                <div>{config.name}</div>
                <ColorPicker
                    disabled={config.modify == 1 ? false : true}
                    onChangeComplete={handleChange}
                    value={config.value}
                    open={colorPickerOpen}
                    onOpenChange={setColorPickerOpen}
                    showText
                />
            </Space>
        )
    };
    const renderSwitch = (config: ConfigInfo) => {
        const handleChange = (c: any) => {
            if (c) {
                config.value = "1"
            } else {
                config.value = "0"
            }
            setShowSaveButton(true);
        };
        return (
            <Space direction="horizontal" split={<Divider type="vertical"/>}>
                <div>{config.name}</div>
                <Switch checkedChildren="开启" disabled={config.modify == 1 ? false : true}
                        unCheckedChildren="关闭" value={config.value == "1" ? true : false}
                        defaultChecked onChange={handleChange}/>
            </Space>
        );
    }
    const renderCheckboxs = (config: ConfigInfo) => {
        const plainOptions = config.values.map((item) => {
            return {
                label: item.name.toString(),
                value: item.key.toString(),
            }
        });
        const onChange = (checkedValues: string[]) => {
            config.value = checkedValues.join(",");
            console.log('checked = ', checkedValues.join(","));
            setShowSaveButton(true);
        }
        return (
            <Space direction="horizontal" split={<Divider type="vertical"/>}>
                <div>{config.name}</div>
                <Checkbox.Group options={plainOptions} disabled={config.modify == 1 ? false : true}
                                defaultValue={config.value.split(",")} onChange={onChange}/>
            </Space>
        )
    }
    const renderRadio = (config: ConfigInfo) => {
        const plainOptions = config.values.map((item) => {
            return {
                label: item.name.toString(),
                value: item.key.toString(),
            }
        });
        const onChange = ({target: {value}}: RadioChangeEvent) => {
            config.value = value
            setShowSaveButton(true)
        };
        return (
            <Space direction="horizontal" split={<Divider type="vertical"/>}>
                <div>{config.name}</div>
                <Radio.Group options={plainOptions} onChange={onChange}
                             disabled={config.modify == 1 ? false : true}
                             value={config.value} optionType="button"
                             buttonStyle="solid"/>
            </Space>
        )
    }
    const renderInput = (config: ConfigInfo) => {
        const handleChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
            config.value = event.target.value;
            setShowSaveButton(true);
        };
        return (
            <Space direction="horizontal" split={<Divider type="vertical"/>}>
                <div>{config.name}</div>
                <TextArea placeholder='请输入' value={config.value}
                          disabled={config.modify == 1 ? false : true} style={{width: '100%', minWidth: '300px'}}
                          autoSize
                          onChange={handleChange}/>
            </Space>
        )
    }
    const renderSelects = (config: ConfigInfo) => {
        const handleChange = (value: string[]) => {
            config.value = value.join(",");
            setShowSaveButton(true);
        };
        const options = config.values.map((item) => {
            return {
                label: item.name.toString(),
                value: item.key.toString(),
            }
        });

        return (
            <Space direction="horizontal" split={<Divider type="vertical"/>}>
                <div>{config.name}</div>
                <Select
                    mode="tags"
                    disabled={config.modify == 1 ? false : true}
                    style={{minWidth: 280}}
                    defaultValue={config.value.split(",")}
                    onChange={handleChange}
                    tokenSeparators={[',']}
                    options={options}
                />
            </Space>
        )
    }
    const renderImage = (config: ConfigInfo) => {
        const handleFileUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
            const file = e.target.files?.[0];
            if (file) {
                // 新增文件大小校验（5MB = 10 * 1024 * 1024）
                if (file.size > 5 * 1024 * 1024) {
                    showMessage("error", "图片大小不能超过10MB");
                    e.target.value = ''; // 清空文件选择
                    return;
                }
                const reader = new FileReader();
                reader.onload = (event) => {
                    if (event.target?.result) {
                        setBackgroundImgValue(event.target.result as string)
                        setShowSaveButton(true);
                    }
                };
                reader.readAsDataURL(file);
            }
        };
        return (
            <Space direction="horizontal" split={<Divider type="vertical"/>}>
                <span>{config.name}</span>
                <Image
                    onClick={() => {
                        setBackgroundImgValue("")
                    }}
                    width={100}
                    src={config.value}
                    alt="背景图"
                    placeholder={true}
                    preview={{
                        imageRender: () => (
                            <Space direction="vertical" size={10} align="start">
                                <Image
                                    preview={false}
                                    src={config.value}
                                    alt="背景图预览"
                                    width={480}
                                    placeholder={true}
                                />
                                <div style={{
                                    width: 440,
                                    marginLeft: 15,
                                    whiteSpace: 'nowrap',
                                    textOverflow: 'ellipsis',
                                    overflow: 'hidden',
                                }}>{config.value}</div>
                            </Space>
                        ),
                        toolbarRender: () => (
                            <Space direction="horizontal" size={20}>
                                {/* 远程URL输入 */}
                                <Input
                                    style={{minWidth: 350, backgroundColor: '#a6a5a5'}}
                                    value={backgroundImgValue}
                                    onChange={(e) => {
                                        setBackgroundImgValue(e.target.value)
                                    }}
                                    placeholder="输入图片URL或上传本地图片"
                                    addonAfter={
                                        <label style={{cursor: 'pointer'}}>
                                            <input
                                                type="file"
                                                accept="image/*"
                                                onChange={handleFileUpload}
                                                style={{display: 'none'}}

                                            />
                                            <UploadOutlined/>
                                        </label>
                                    }
                                />
                                <Button
                                    onClick={() => {
                                        config.value = backgroundImgValue;
                                        setShowSaveButton(true);
                                    }}
                                >
                                    确认切换
                                </Button>
                            </Space>
                        )
                    }}
                />
            </Space>
        )
    }

    //保持最新配置文件
    const saveConfig = () => {
        const updateInfo = () => {
            setLoading(true);
            //调用后端保存最新得配置文件value
            UpdateSetting(configStore.getConfigInfo()).then(() => {
                setShowSaveButton(false);
                setLoading(false);
            }).catch(() => {
                setLoading(false);
                showMessage("error", "配置项保存失败,请重新操作");
            });
        }
        return (
            <Space>
                <Button type="primary" onClick={updateInfo}>保存修改配置</Button>
                <Button>取消</Button>
            </Space>
        )
    }

    return (
        <Spin spinning={loading} tip="配置信息加载中..." size="large">
            <Tabs
                tabPosition={'top'}
                tabBarExtraContent={showSaveButton ? saveConfig() : ''}
                items={(configStore.getConfigInfo() || []).map((config, i) => {
                    return {
                        label: config.name,
                        key: config.key,
                        children: (<div style={{
                            display: 'flex',
                            flexDirection: 'column',
                            minHeight: '470px',
                            backgroundColor: '#ededed',
                            padding: '20px',
                        }}>
                            {/*修改后得保存按钮和取消*/}
                            {getComponent(config)}
                        </div>),
                    };
                })}
            />
        </Spin>
    );
};
export default observer(Setting);
