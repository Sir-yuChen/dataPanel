import { CommentOutlined, CustomerServiceOutlined } from '@ant-design/icons';
import { FloatButton } from "antd";

const FloatButtonsComponent = () => {

    return (
        <>
            <FloatButton.Group
                trigger="hover"
                type="primary"
                style={{ insetInlineEnd: 94 }}
                icon={<CustomerServiceOutlined />}
            >
                <FloatButton />
                <FloatButton icon={<CommentOutlined />} />
            </FloatButton.Group>
        </>
    )
}

export default FloatButtonsComponent