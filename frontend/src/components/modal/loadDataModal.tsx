import { LoadDataModalSubmit } from '@/api/commonApi';
import { ReqField } from '@/model/loadDtaModel';
import { Button, Checkbox, Flex, Form, FormProps, Input, Radio, Space, Spin } from 'antd';
import { useState } from 'react';

const baseStyle: React.CSSProperties = {
    width: '100%',
    height: 54,
};
export const LoadDataModalComponent = () => {
    const [loadDataForm] = Form.useForm();

    const [loadType, setLoadType] = useState<string>('default');
    const [loading, setLoading] = useState<boolean>(false);
    const tailLayout = {
        wrapperCol: { flex: 1, offset: 8, span: 16 },
    };

    const onReset = () => {
        loadDataForm.resetFields();
    };
    const onFinish: FormProps<ReqField>['onFinish'] = (values) => {
        console.log('Success:', values);
        setLoading(true);
        LoadDataModalSubmit(values).then((res) => {
            setLoading(false);
        });

    };

    return (
        <div style={{ width: '500px', height: '100%' }}>
            <Flex gap="middle" vertical >
                <Flex vertical>
                    <Spin spinning={loading} tip="配置加载中,请耐心等待">
                        <Form
                            form={loadDataForm}
                            layout="vertical"
                            requiredMark
                            onFinish={onFinish}
                            initialValues={{ loadDataType: 'default', dataSavePath: '/dataPanel', loadDataChecked: ['c'] }}
                        >
                            <Form.Item label="加载数据方式" name="loadDataType" required>
                                <Radio.Group defaultValue={'default'} onChange={(e) => { setLoadType(e.target.value) }}>
                                    <Radio.Button value="default">初始化默认配置及数据</Radio.Button>
                                    <Radio.Button value='customize'>导入历史配置及数据</Radio.Button>
                                </Radio.Group>
                            </Form.Item>
                            <Form.Item hidden={loadType == 'default' ? false : true}
                                required={loadType == 'default' ? true : false}
                                name="dataSavePath"
                                label="配置文件及数据目录"
                                tooltip="用来存放配置文件及应用所有数据">
                                <Input placeholder="input placeholder" />
                            </Form.Item>
                            <Form.Item hidden={loadType == 'default' ? false : true}
                                required={loadType == 'default' ? true : false}
                                name="loadDataChecked"
                                valuePropName="checked"
                                label="初始化数据类型"
                                tooltip="请选择需要初始化加载的数据(注意：为不影响应用使用，配置文件将同步加载，其它数据将在后台加载)">
                                <Checkbox.Group defaultValue={['c']} >
                                    <Checkbox value="c" disabled >配置文件</Checkbox>
                                    <Checkbox value="a">A股基础数据</Checkbox>
                                    <Checkbox value="h">港股基础数据</Checkbox>
                                    <Checkbox value="m">美股基础数据</Checkbox>
                                </Checkbox.Group>
                            </Form.Item>
                            {/* 提交 */}
                            <Form.Item {...tailLayout}>
                                <Space>
                                    <Button type="primary" htmlType="submit">
                                        提交
                                    </Button>
                                    <Button htmlType="button" onClick={onReset}>
                                        重置
                                    </Button>
                                </Space>
                            </Form.Item>
                        </Form>
                    </Spin>
                </Flex>
            </Flex>
        </div>
    )
}