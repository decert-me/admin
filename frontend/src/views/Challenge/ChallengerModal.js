import { Modal, Table, message, Input } from "antd";
import { useEffect, useState } from "react";
import { getQuestAnswerList } from "../../request/api/quest";
import JudgReviewModal from "./JudgReviewModal";

export default function ChallengerModal(props) {

    const {challenge} = props;
    const [data, setData] = useState([]);
    let [pageConfig, setPageConfig] = useState({
        page: 0, pageSize: 50, total: 0
    });

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

    const columns = [
        {
            title: "挑战名称",
            dataIndex: "title",
            ellipsis: true
        },
        {
          title: "挑战者地址",
          dataIndex: "address",
          render: (addr) => (
            <a target="_blank" href={`https://decert.me/user/${addr}`}>{addr.substring(0,5) + "..." + addr.substring(38,42)}</a>
          )
        },
        {
          title: "名称",
          dataIndex: "name"
        },
        {
          title: "标签",
          dataIndex: "tags",
          ellipsis: true
        },
        {
            title: "挑战结果",
            dataIndex: "pass",
            render: (pass) => (
                pass ? "成功" : "失败"
            )
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
            dataIndex: "score_detail"
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
            render: (time) => (
              time.indexOf("0001-01-01T") === -1 ?
              time.replace("T", " ").split(".")[0].split("+")[0]
              :"-"
          )
        },
    ];

    function changePageSize(props) {
        pageConfig.pageSize = props.pageSize;
        setPageConfig({...pageConfig});
        getList();
    }

    async function getList(page) {
        if (page) {
          pageConfig.page = page;
          setPageConfig({ ...pageConfig });
        }
        // 获取教程列表
        let res = await getQuestAnswerList({ ...pageConfig, ...challenge });
    
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
    }, [challenge]);

    

    return (
        <Table
            columns={columns} 
            dataSource={data}         
            onChange={(e) => changePageSize(e)}
            pagination={{
                current: pageConfig.page, 
                total: pageConfig.total, 
                pageSize: pageConfig.pageSize, 
                onChange: (page) => getList(page)
            }} />
    )
}