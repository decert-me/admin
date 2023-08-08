import { LockOutlined, UserOutlined, SafetyOutlined } from '@ant-design/icons';
import { Button, Checkbox, Col, Form, Input, Row } from 'antd';
import { useAuth } from "../../hooks/useAuth";
import "./index.scss";
import { useEffect, useState } from 'react';
import { userCaptcha, userLogin } from '../../request/api/user';

export default function LoginPage(params) {

    const { login } = useAuth();
    const [captcha, setCaptcha] = useState();

    const onFinish = (values) => {
        const { username, password, captcha } = values;
        userLogin({
            username, password, captcha, CaptchaId: captcha.captchaId
        })
        .then(res => {
            if (res.code === 0) {
                const { token, user } = res.data;
                login(token, user);
            }
        })
    };

    function getCaptcha() {
        // TODO: 获取验证码
        userCaptcha()
        .then(res => {
            if (res?.code === 0) {
                const { picPath, captchaId } = res.data;
                setCaptcha({picPath, captchaId})
            }
        })
    }

    useEffect(() => {
        getCaptcha();
    },[])

    return (
        <div className="login">
            <div className="login-content">
                <Form
                    name="normal_login"
                    className="login-form"
                    initialValues={{
                        remember: true,
                    }}
                    onFinish={onFinish}
                    >

                    {/* 用户名 */}
                    <Form.Item
                        name="username"
                        rules={[
                        {
                            required: true,
                            message: 'Please input your Username!',
                        },
                        ]}
                    >
                        <Input prefix={<UserOutlined />} placeholder="Username" />
                    </Form.Item>

                    {/* 密码 */}
                    <Form.Item
                        name="password"
                        rules={[
                        {
                            required: true,
                            message: 'Please input your Password!',
                        },
                        ]}
                    >
                        <Input
                        prefix={<LockOutlined />}
                        type="password"
                        placeholder="Password"
                        />
                    </Form.Item>

                    {/* 验证码 */}
                    <Form.Item>
                        <Row gutter={8}>
                        <Col span={12}>
                            <Form.Item
                            name="captcha"
                            noStyle
                            rules={[
                                {
                                required: true,
                                message: 'Please input the captcha you got!',
                                },
                            ]}
                            >
                            <Input 
                                prefix={<SafetyOutlined />}
                                placeholder="Captcha"
                            />
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            {
                                captcha &&
                                <img src={captcha.picPath} alt="" className="captcha" onClick={() => getCaptcha()} />
                            }
                        </Col>
                        </Row>
                    </Form.Item>

                    <Form.Item>
                        <Button type="primary" htmlType="submit" className="login-form-button">
                        Log in
                        </Button>
                    </Form.Item>
                    Or <a href="">register now!</a>

                </Form>
            </div>
        </div>
    )
}