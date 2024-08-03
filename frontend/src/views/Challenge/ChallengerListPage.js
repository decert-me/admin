import { Button, Form, Input, Modal, Table, message } from "antd";
import { useEffect, useState } from "react";
import { getChallengeUserStatistics } from "../../request/api/quest";
import ChallengerModal from "./ChallengerModal";

export default function ChallengerListPage() {

    const [data, setData] = useState([]);
    const [form, setForm] = useState({}); //  搜索
    let [pageConfig, setPageConfig] = useState({
        page: 0, pageSize: 50, total: 0
    });

    const columns = [
        {
          title: "挑战者地址",
          dataIndex: "address",
          render: (addr) => (
            <a target="_blank" href={`https://decert.me/user/${addr}`}>{addr.substring(0,5) + "..." + addr.substring(38,42)}</a>
          )
        },
        {
          title: "昵称",
          dataIndex: "name"
        },
        {
          title: "标签",
          dataIndex: "tags",
          ellipsis: true
        },
        {
          title: "全部挑战",
          render: (_, user) => (
            <a onClick={() => info({search_address: user.address})}>{user.success_num + user.fail_num}</a>
          )
        },
        {
            title: "挑战成功",
            dataIndex: "success_num",
            render: (success_num, user) => (
              <a onClick={() => info({search_address: user.address, pass: true})}>{success_num}</a>
            )
        },
        {
            title: "挑战失败",
            dataIndex: "fail_num",
            render: (fail_num, user) => (
              <a onClick={() => info({search_address: user.address, pass: false})}>{fail_num}</a>
            )
        },
        {
            title: "领取NFT",
            dataIndex: "claim_num",
            render: (claim_num, user) => (
              <a onClick={() => info({search_address: user.address, claimed: true, pass: true})}>{claim_num}</a>
            )
        },
        {
            title: "未领取NFT",
            dataIndex: "not_claim_num",
            render: (not_claim_num, user) => (
              <a onClick={() => info({search_address: user.address, claimed: false, pass: true})}>{not_claim_num}</a>
            )
        }
    ];

    const info = (challenge) => {
        Modal.info({
            icon: <></>,
            width: "1200px",
            title: '',
            content: <ChallengerModal challenge={challenge} />,
            footer: null,
            maskClosable: true
        });
    };

    function onFinish(params) {
        setForm({...params});
    }

    async function getList(page) {
        if (page) {
          pageConfig.page = page;
          setPageConfig({ ...pageConfig });
        }
        // 获取教程列表
        let res = await getChallengeUserStatistics({ ...pageConfig, search_address: form?.challenger, search_tag: form?.tag });
    
        if (res.code === 0) {
          const list = res.data.list || [];
          // 添加key
          list.forEach((ele) => {
            ele.key = ele.id;
          });
          setData([...list]);
          pageConfig.total = res.data.total;
          setPageConfig({ ...pageConfig });
        } else {
          message.success(res.msg);
        }
    }

    function init(params) {
        pageConfig.page += 1;
        setPageConfig({ ...pageConfig });
        getList();
    }

    useEffect(() => {
        pageConfig = {
          page: 0,
          pageSize: 50,
          total: 0,
        };
        setPageConfig({ ...pageConfig });
        init();
    }, [form]);
    
    return (
        <div>
            <div className="tabel-title">
                <h2>挑战者</h2>
            </div>
            <div>
                <div className="operat">
                    <Form
                        name="horizontal_login"
                        layout="inline"
                        onFinish={onFinish}
                    >
                        <Form.Item label="标签" name="tag">
                            <Input />
                        </Form.Item>
                        <Form.Item label="挑战者" name="challenger">
                            <Input />
                        </Form.Item>
                        <Form.Item shouldUpdate>
                        {() => (
                            <Button type="primary" htmlType="submit">
                            搜索
                            </Button>
                        )}
                        </Form.Item>
                    </Form>
                </div>
            </div>
            <Table
                columns={columns} 
                dataSource={data}         
                pagination={{
                    current: pageConfig.page, 
                    total: pageConfig.total, 
                    pageSize: pageConfig.pageSize, 
                    onChange: (page) => getList(page)
                }} />
        </div>
    )
}