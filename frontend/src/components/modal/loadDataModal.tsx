import {LoadDataModalSubmit} from '@/api/commonApi';
import {Checkbox, Flex, Form, Input, Modal, Radio} from 'antd';
import {useState} from 'react';
import NiceModal, {useModal} from "@ebay/nice-modal-react";
import {useShowMessage} from "@/hooks/useShowMessage";

export const LoadDataModalComponent = NiceModal.create(() => {
    const modal = useModal();
    const [loadDataForm] = Form.useForm();
    const [loadType, setLoadType] = useState<string>('default');
    const [confirmLoading, setConfirmLoading] = useState(false);
    const showMessage = useShowMessage();
    const onReset = () => {
        console.log("重置表单")
        loadDataForm.resetFields();
        showMessage("success", "重置成功");
    };
    const onFinish = () => {
        const values = loadDataForm.getFieldsValue();
        setConfirmLoading(true);
        LoadDataModalSubmit(values).then((res) => {
            setConfirmLoading(false);
            if (res.code == 200) {
                modal.remove();
            } else {
                //提示需要重新提交
                modal.reject();
            }
        });
    };
    return (
        <Modal
            title="资源加载"
            onOk={onFinish}
            visible={modal.visible}
            onCancel={onReset}
            centered={true}
            maskClosable={false}
            cancelText="重置"
            closable={false}
            okText="提交"
            confirmLoading={confirmLoading}
        >
            <Flex vertical>
                <Form
                    form={loadDataForm}
                    layout="vertical"
                    requiredMark
                    initialValues={{
                        loadDataType: 'default',
                        loadDataChecked: []
                    }}
                >
                    <Form.Item name="loadDataType" rules={[{required: true, message: '必选项'}]}>
                        <Radio.Group defaultValue={'default'} onChange={(e) => {
                            setLoadType(e.target.value)
                        }}>
                            <Radio.Button value="default">初始化默认配置及数据</Radio.Button>
                            <Radio.Button value='customize'>导入历史配置及数据</Radio.Button>
                        </Radio.Group>
                    </Form.Item>
                    <Form.Item
                        hidden={loadType == 'customize' ? false : true}
                        rules={[{
                            required: loadType == 'customize' ? true : false,
                            message: loadType == 'customize' ? '必选项' : '可选项'
                        }]}
                        name="dataSavePath"
                        label="历史数据目录"
                        tooltip="历史数据目录">
                        <Input/>
                    </Form.Item>
                    <Form.Item hidden={loadType == 'default' ? false : true}
                               rules={[{
                                   required: loadType == 'default' ? true : false,
                                   message: loadType == 'default' ? '必选项' : '可选项'
                               }]}
                               name="loadDataChecked"
                               valuePropName="checked"
                               label="初始化数据类型"
                               tooltip="请选择需要初始化加载的数据(注意：为不影响应用使用，配置文件将同步加载，其它数据将在后台加载)">
                        <Checkbox.Group>
                            <Checkbox value="c">配置文件</Checkbox>
                            <Checkbox value="a">A股基础数据</Checkbox>
                            <Checkbox value="h">港股基础数据</Checkbox>
                            <Checkbox value="m">美股基础数据</Checkbox>
                        </Checkbox.Group>
                    </Form.Item>
                </Form>
            </Flex>
        </Modal>
    )
});