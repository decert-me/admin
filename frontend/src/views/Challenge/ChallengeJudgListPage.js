import { Button, Modal, Space, Table } from "antd"
import { useEffect, useRef, useState } from "react";
import { getUserOpenQuestList } from "../../request/api/judgment";
import ChallengeJudgPage from "./ChallengeJudgPage";
import "./judg.scss";
import { copyToClipboard } from "../../utils/text/copyToClipboard";

const location = window.location.host;
const isTest = ((location.indexOf("localhost") !== -1) || (location.indexOf("192.168.1.10") !== -1)) ? false : true;
const host = isTest ? "https://decert.me" : "http://192.168.1.10:8087";
const { confirm } = Modal;

export default function ChallengeJudgListPage(params) {

    const [detailOpen, setDetailOpen] = useState(false);
    const [rateNum, setRateNum] = useState(0);
    const [questDetail, setQuestDetail] = useState();
    const [reviewStatus, setReviewStatus] = useState();     // true: 未评分 || false: 已评分
    
    let [tableLoad, setTableLoad] = useState();
    
    let [status, setStatus] = useState(1);
    let [data, setData] = useState([]);
    let [pageConfig, setPageConfig] = useState({
        page: 0, pageSize: 10, total: 0
    });

    const handleChange = (pagination, filters, sorter) => {
        const { pageSize } = pagination
        const newStatus = Array.isArray(filters.status) ? filters.status[0] : null;
        if (status !== newStatus) {
            status = newStatus;
            setStatus(newStatus);
            getList(1);
        }
        if (pageSize !== pageConfig.pageSize) {
            pageConfig.pageSize = pageSize;
            setPageConfig({...pageConfig});
            getList();
        }
    };

    const columns = [
        {
            title: "题目",
            key: 'title',
            dataIndex: "title",
            render: (title, quest) => (
                <p className="of-h pointer" onClick={() => window.open(`${host}/quests/${quest.token_id}`, "_blank")}>{title}</p>
            )
        },
        {
            title: "挑战编号",
            key: 'token_id',
            dataIndex: "token_id",
            ellipsis: true,
            render: (token_id) => (
                <p className="of-h pointer" onClick={() => window.open(`${host}/quests/${token_id}`, "_blank")}>{token_id}</p>
            )
        },
        {
            title: "待评分",
            key: 'to_review_count',
            dataIndex: "to_review_count"
        },
        {
            title: "最新提交时间",
            key: 'last_sumbit_time',
            dataIndex: "last_sumbit_time",
            render: (time) => (
                time.indexOf("0001-01-01T") === -1 ?
                time.replace("T", " ").split(".")[0]
                :"-"
            )
        },
        {
            title: "上次评分时间",
            key: 'last_review_time',
            dataIndex: "last_review_time",
            render: (time) => (
                time.indexOf("0001-01-01T") === -1 ?
                time.replace("T", " ").split(".")[0].split("+")[0]
                :"-"
            )
        },
        {
            title: "操作",
            key: "action",
            render: (_, quest) => (
              <Space size="middle">
                <Button
                  type="link"
                  className="p0"
                    // onClick={() => navigateTo(`/dashboard/user/tag/modifyuser/${tag.address}`)}
                    onClick={() => openDetail(quest, true)}
                >
                  评分
                </Button>
                <Button
                  type="link"
                  className="p0"
                    // onClick={() => navigateTo(`/dashboard/user/tag/modifyuser/${tag.address}`)}
                    onClick={() => openDetail(quest, false)}
                >
                  已评分
                </Button>
              </Space>
            ),
        },
    ];

    // 展示该题详情
    function openDetail({index, token_id}, isReview) {
        const obj = {index, token_id};
        setQuestDetail({...obj});
        setReviewStatus(isReview);
        setDetailOpen(true);
    }
    
    // 获取列表
    async function getList(page) {
        setTableLoad(true);
        if (page) {
            pageConfig.page = page;
            setPageConfig({...pageConfig});
        }
        await getUserOpenQuestList({
            open_quest_review_status: status,
            ...pageConfig
        })
        .then(res => {
            const list = res.data.list;
            data = list ? list : [];
            // 添加key
            data.forEach((ele, index) => {
                ele.key = ele.uuid + ele.index
            })
            setData([...data]);
            pageConfig.total = res.data.total;
            if (status === 1) {
                setRateNum(pageConfig.total);
            }
            setPageConfig({...pageConfig});
        })
        setTableLoad(false);
    }

    async function init() {
        pageConfig.page += 1;
        setPageConfig({...pageConfig});
        await getList();
    }

    useEffect(() => {
        init();
    },[])

    return (
        <div className="judg">
            <div className="tabel-title">
                <h2 style={{fontSize: 20}}>评分列表
                    {/* <span style={{fontSize: 14, fontWeight: 400, color: "#999999"}}>（待评分 {rateNum}）</span> */}
                </h2>
            </div>

            <Modal
                width={1177}
                className="judg-modal"
                open={detailOpen}
                footer={null}
                onCancel={() => {setDetailOpen(false)}}
            >
                <ChallengeJudgPage questDetail={questDetail} reviewStatus={reviewStatus} hideModal={() => setDetailOpen(false)} updateList={() => getList()} />
            </Modal>
            <Table
                columns={columns} 
                dataSource={data} 
                loading={tableLoad}
                onChange={handleChange}
                pagination={{
                    current: pageConfig.page, 
                    total: pageConfig.total, 
                    pageSize: pageConfig.pageSize, 
                    onChange: (page) => {
                        page !== pageConfig.page && getList(page)
                    }
                }} 
            />
        </div>
    )
}