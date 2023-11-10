import { useConnect, useDisconnect, useSignMessage } from 'wagmi';
import { Button } from "antd";
import { useAuth } from "../../hooks/useAuth";
import "./index.scss";
import { authLoginSign, getLoginMessage } from '../../request/api/sign';

export default function LoginPage(params) {

    const { login } = useAuth();
    const { disconnect } = useDisconnect();
    const { connect, connectors, isLoading } = useConnect({
        onSuccess(data) {
            goSignature(data.account);
        },
        onError(err) {
            disconnect();
            goConnect();
            console.log("===>", err);
        }
    });
    const { signMessage, signMessageAsync } = useSignMessage()

    function goSignature(account) {
        new Promise(async(resolve, reject) => {
            // 获取nonce
            let message;
            await getLoginMessage({address: account})
            .then(res => {
                if (res.code === 0) {
                    message = res.data.loginMessage;
                }else{
                    reject();
                }
            })
            // 发起签名
            signMessageAsync({ message })
            .then(res => {
                resolve({message, signature: res});
            })
            .catch(err => {
                reject();
            })
        })
        .then(({message, signature}) => {
            authLoginSign({
                address: account,
                message,
                signature
            })
            .then(res => {
                const { token, user } = res.data;
                login(token, user);
            })
        })
        .catch(err => {
            // 登陆失败 => 断开连接
            disconnect();
        })
    }

    function goConnect(params) {
        connect({connector: connectors[0]});
    }

    return (
        <div className="login">
            <div className="login-content">
                <img src="https://learnblockchain.cn/css/default/metamask-login.svg" alt="" />
                <p>MetaMask 登录</p>
                <Button 
                    onClick={() => goConnect()}
                    disabled={!connectors[0].ready}
                    loading={isLoading}
                >连接钱包</Button>
            </div>
        </div>
    )
}