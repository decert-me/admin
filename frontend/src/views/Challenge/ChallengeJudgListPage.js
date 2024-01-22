import { Button, Modal, Table } from "antd"
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

    const judgRef = useRef(null);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [detailOpen, setDetailOpen] = useState(false);
    const [rateNum, setRateNum] = useState(0);
    let [isLoading, setIsLoading] = useState();
    let [tableLoad, setTableLoad] = useState();
    let [detail, setDetail] = useState([]);
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
                <p className="of-h pointer" onClick={() => openDetail(quest)}>{title}</p>
            )
        },
        {
            title: "挑战编号",
            key: 'token_id',
            dataIndex: "token_id",
            render: (token_id) => (
                <p className="pointer text-w-300" onClick={() => window.open(`${host}/quests/${token_id}`, "_blank")}>{token_id}</p>
            )
        },
        {
            title: "挑战者地址",
            key: 'address',
            dataIndex: "address",
            render: (address) => (
                <p className="pointer" onClick={() => window.open(`${host}/${address}`, "_blank")}>{address.substring(0,5) + "..." + address.substring(38,42)}</p>
            )
        },
        {
            title: "状态",
            key: 'status',
            dataIndex: "open_quest_review_status",
            filters: [
                { text: "全部", value: null },
                { text: "待评分", value: 1 },
                { text: "已评分", value: 2 },
            ],
            filterMultiple: false,
            filteredValue: [status],
            render: (status) => (
                <p style={{
                    color: status === 2 ? "#35D6A6" : "#9A9A9A",
                    fontWeight: 600
                }}>{status === 2 ? "已评分" : status === 1 ? "待评分" : "全部"}</p>
            )
        },
        {
            title: "提交时间",
            key: 'updated_at',
            dataIndex: "updated_at",
            render: (time) => (
                time.indexOf("0001-01-01T") === -1 ?
                time.replace("T", " ").split(".")[0]
                :"-"
            )
        },
        {
            title: "处理时间",
            key: 'open_quest_review_time',
            dataIndex: "open_quest_review_time",
            render: (time) => (
                time.indexOf("0001-01-01T") === -1 ?
                time.replace("T", " ").split(".")[0]
                :"-"
            )
        }
    ];

    // 展示该题详情
    function openDetail(quest) {
        detail = [quest];
        setDetail([...detail]);
        setDetailOpen(true);
    }

    // 判题弹窗
    function OpenJudgModal() {
        setIsModalOpen(true);
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
                ele.key = ele.updated_at + ele.index + index
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

    // 提交批改内容
    async function submitReview(params) {
        setIsLoading(true);
        await judgRef.current.confirm();
        setIsLoading(false);
    }

    function handleOk() {
        // 是否批改完
        const {flag, remain} = judgRef.current.isOver();
        if (flag) { 
            submitReview()
        }else{
            confirm({
                title: `还有${remain}道未评分，仍然提交？`,
                onOk() {
                    submitReview()
                },
            });
        }
    }

    function onFinish() {
        setIsModalOpen(false)
        getList()
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
                <h2 style={{fontSize: 20}}>评分列表<span style={{fontSize: 14, fontWeight: 400, color: "#999999"}}>（待评分 {rateNum}）</span></h2>
                {
                    data?.findIndex((e) => e.open_quest_review_status === 1) !== -1 &&
                    <Button id="hover-btn-full" className="btn-start" onClick={() => OpenJudgModal()}>开始评分</Button>
                }
            </div>
            <Modal
                width={1177}
                open={isModalOpen}
                className="judg-modal"
                onCancel={() => {setIsModalOpen(false)}}
                onOk={handleOk}
                okButtonProps={{
                    loading: isLoading
                }}
            >
                <ChallengeJudgPage ref={judgRef} rateNum={rateNum} status={status} pageNum={pageConfig.page} data={data} onFinish={onFinish} />
            </Modal>

            <Modal
                width={1177}
                className="judg-modal"
                open={detailOpen}
                footer={null}
                onCancel={() => {setDetailOpen(false)}}
            >
                <ChallengeJudgPage data={detail} />
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