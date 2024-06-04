import { useEffect, useState } from "react";
import { useLocation, useNavigate, useParams } from "react-router-dom";
import { Button, Form, Input, Table, message } from "antd";
import { getQuestAnswerList } from "../../request/api/quest";

export default function ChallengeAnswerListPage() {

  const location = useLocation();
  const navigateTo = useNavigate();
  const { tokenId } = useParams();
  const [data, setData] = useState([]);
  const [search_key, setSearch_key] = useState(""); //  搜索
  let [pageConfig, setPageConfig] = useState({
      page: 0, pageSize: 10, total: 0
  });

  const columns = [
    {
      title: "挑战名称",
      dataIndex: "title",
    },
    {
      title: "挑战者地址",
      dataIndex: "address",
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
      title: "领取NFT",
      dataIndex: "claimed",
      render: (claimed) => (
        claimed ? "是" : "否"
      )
    },
    {
      title: "得分/及格分",
      dataIndex: "score_detail",
    },
    {
      title: "批注",
      dataIndex: "annotation",
      ellipsis: true
    },
    {
        title: "挑战时间",
        dataIndex: "challenge_time",
        render: (time) => (
          time.indexOf("0001-01-01T") === -1 ?
          time.replace("T", " ").split(".")[0].split("+")[0]
          :"-"
      )
    },
  ];

  function init() {
    getQuestAnswerList({ id: tokenId }).then((res) => {
      const list = res.data || [];
      // setData([...list]);
      console.log(list);
    });
  }

  async function getList(page) {
    if (page) {
      pageConfig.page = page;
      setPageConfig({ ...pageConfig });
    }
    // 获取教程列表
    let res = await getQuestAnswerList({ 
      ...pageConfig, 
      "search_quest": tokenId,
      "search_tag": "",
      "search_address": "",
      // "pass": true,
      // "claimed": false
    });

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
    <div className="challenge" key={location.pathname}>
      <div className="tabel-title">
        <h2>挑战详情统计</h2>
      </div>
      <Form
        name="horizontal_login"
        layout="inline"
        //   onFinish={onFinish}
      >
        <Form.Item label="挑战" name="password">
          <Input />
        </Form.Item>
        <Form.Item label="标签" name="username">
          <Input />
        </Form.Item>
        <Form.Item label="挑战者地址" name="password">
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
      <Table
        columns={columns} 
        dataSource={data}         
        pagination={{
            current: pageConfig.page, 
            total: pageConfig.total, 
            pageSize: pageConfig.pageSize, 
            onChange: (page) => getList(page)
        }} 
      />
    </div>
  );
}
