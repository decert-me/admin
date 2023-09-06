import {
    LoadingOutlined,
} from '@ant-design/icons';
import { Spin } from 'antd';


export default function CustomLoading(props) {

    const { className } = props;

    return (
        <div className={`${className ? className : "loading"}`}>
            <Spin 
                indicator={
                    <LoadingOutlined
                        style={{
                        fontSize: 24,
                        }}
                        spin
                    />
                } 
            />
        </div>
    )
}