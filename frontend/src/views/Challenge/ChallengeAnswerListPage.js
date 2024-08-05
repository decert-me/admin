import { useEffect, useState } from "react";
import { useLocation, useNavigate, useParams } from "react-router-dom";
import { Button, Form, Input, Modal, Space, Table, message } from "antd";
import { getChallengeStatisticsSummary, getQuestAnswerList } from "../../request/api/quest";
import { useUpdateEffect } from "ahooks";
import JudgReviewModal from "./JudgReviewModal";
const { TextArea } = Input;

export default function ChallengeAnswerListPage() {
  const location = useLocation();
  const navigateTo = useNavigate();
  const [formProps] = Form.useForm();
  const { tokenId } = useParams();
  const [data, setData] = useState([]);
  const [form, setForm] = useState({tokenId: tokenId}); //  搜索
  const [totalObj, setTotalObj] = useState({});
  let [pageConfig, setPageConfig] = useState({
    page: 0,
    pageSize: 50,
    total: 0,
  });

  const columns = [
    {
      title: "挑战者地址",
      dataIndex: "address",
      render: (address) => (
        <a target="_blank" href={`${process.env.REACT_APP_LINK_URL || "https://decert.me"}/user/${address}`}>{address}</a>
      )
    },
    {
      title: "昵称",
      dataIndex: "name",
    },
    {
      title: "标签",
      dataIndex: "tags",
      ellipsis: true,
    },
    {
      title: "挑战名称",
      dataIndex: "title",
      render: (title, quest) => (
        <a target="_blank" href={`${process.env.REACT_APP_LINK_URL || "https://decert.me"}/quests/${quest.uuid}`}>{title}</a>
      )
    },
    {
      title: "挑战结果",
      dataIndex: "challenge_result",
    },
    {
      title: "领取NFT",
      dataIndex: "claimed",
      render: (claimed) => (claimed ? "是" : "否"),
    },
    {
      title: "得分/及格分",
      dataIndex: "score_detail",
    },
    {
      title: "评分详情",
      dataIndex: "annotation",
      ellipsis: true,
      render: (annotation, quest) => (
          <a onClick={() => info(quest)}>{annotation ? "查看" : ""}</a>
      )
    },
    {
      title: "挑战时间",
      dataIndex: "challenge_time",
      render: (time) =>
        time.indexOf("0001-01-01T") === -1
          ? time.replace("T", " ").split(".")[0].split("+")[0]
          : "-",
    },
  ];

  const info = ({uuid, address}) => {
    Modal.info({
        icon: <></>,
        width: "1200px",
        wrapClassName: "judg-modal",
        content: <JudgReviewModal uuid={uuid} address={address} />,
        footer: null,
        closable: true
    });
};

  function onFinish(params) {
    setForm({ ...params });
  }

  async function getList(page) {
    if (page) {
      pageConfig.page = page;
      setPageConfig({ ...pageConfig });
    }
    // 获取教程列表
    let res = await getQuestAnswerList({
      ...pageConfig,
      search_quest: form?.tokenId,
      search_tag: form?.tag,
      search_address: form?.addr,
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

  function getTotal() {
    getChallengeStatisticsSummary({
      "search_quest": form?.tokenId,
      "search_tag": form?.tag,
      "search_address": form?.addr,
      // "pass": true,
      // "claimed": false
    })
    .then(res => {
      if (res.code === 0) {
        setTotalObj(res.data);
      }
    })
  }

  function init(params) {
    pageConfig.page += 1;
    setPageConfig({ ...pageConfig });
    getList();
    getTotal();
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

  useEffect(() => {
    formProps.setFieldValue("tokenId", tokenId);
  }, []);

  useUpdateEffect(() => {
    const obj = { tokenId: tokenId || "" };
    formProps.setFieldValue("tokenId", tokenId);
    setForm({ ...obj });
  }, [tokenId]);

  return (
    <div className="challenge" key={location.pathname}>
      <div className="tabel-title">
        <h2>挑战详情统计</h2>
      </div>
      <Form
        name="horizontal_login"
        layout="inline"
        form={formProps}
        onFinish={onFinish}
      >
        <Form.Item label="挑战" name="tokenId">
          <Input />
        </Form.Item>
        <Form.Item label="标签" name="tag">
          <Input />
        </Form.Item>
        <Form.Item label="挑战者地址" name="addr">
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
          onChange: (page) => getList(page),
        }}
      />
      <div>
        <p>总计：</p>
        <br/>
        <Space size={50}>
          <p>挑战数量：{totalObj?.challenge_num}</p>
          <p>挑战人数：{totalObj?.challenge_user_num}</p>
          <p>成功/失败人数：{totalObj?.success_num}/{totalObj?.fail_num}</p>
          <p>领取/未领取人数：{totalObj?.claim_num}/{totalObj?.not_claim_num}</p>
        </Space>
        <br/>
      </div>
    </div>
  );
}
