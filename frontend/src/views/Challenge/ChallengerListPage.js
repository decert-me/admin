import { Button, Form, Input, Modal, Table, message } from "antd";
import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { getChallengeUserStatistics } from "../../request/api/quest";
import ChallengerModal from "./ChallengerModal";

export default function ChallengerListPage() {

    const navigateTo = useNavigate();
    const location = useLocation();
    const [data, setData] = useState([]);
    const [search_key, setSearch_key] = useState(""); //  搜索
    let [pageConfig, setPageConfig] = useState({
        page: 0, pageSize: 10, total: 0
    });

    const columns = [
        {
          title: "挑战者地址",
          dataIndex: "address",
          render: (addr) => (
            <a onClick={() => info(addr)}>{addr.substring(0,5) + "..." + addr.substring(38,42)}</a>
          )
        },
        {
          title: "昵称",
          dataIndex: "nickname"
        },
        {
          title: "标签",
          dataIndex: "tags",
          ellipsis: true
        },
        {
            title: "挑战成功",
            dataIndex: "success_num"
        },
        {
            title: "挑战失败",
            dataIndex: "fail_num"
        },
        {
            title: "领取NFT",
            dataIndex: "claim_num"
        },
        {
            title: "未领取NFT",
            dataIndex: "not_claim_num"
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

    async function getList(page) {
        if (page) {
          pageConfig.page = page;
          setPageConfig({ ...pageConfig });
        }
        // 获取教程列表
        let res = await getChallengeUserStatistics({ ...pageConfig, search_address: "", search_tag: "" });
    
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

    // getTagUserList
    useEffect(() => {
        pageConfig = {
          page: 0,
          pageSize: 10,
          total: 0,
        };
        setPageConfig({ ...pageConfig });
        // if (location.search) {
        //   let serch = new URLSearchParams(location.search);
        //   search_key = serch.get("tokenId");
        //   setSearch_key(search_key);
        // }
        init();
    }, [location]);

    
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
                        //   onFinish={onFinish}
                    >
                        <Form.Item label="标签" name="tag">
                            <Input />
                        </Form.Item>
                        <Form.Item label="挑战者" name="tag">
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