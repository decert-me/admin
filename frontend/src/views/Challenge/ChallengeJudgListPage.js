import { Button, Modal, Table } from "antd"
import { useEffect, useRef, useState } from "react";
import { getUserOpenQuestList } from "../../request/api/judgment";
import ChallengeJudgPage from "./ChallengeJudgPage";
import "./judg.scss";
import { copyToClipboard } from "../../utils/text/copyToClipboard";

const location = window.location.host;
const isTest = ((location.indexOf("localhost") !== -1) || (location.indexOf("192.168.1.10") !== -1)) ? false : true;
const host = isTest ? "https://decert.me" : "http://192.168.1.10:8087";

export default function ChallengeJudgListPage(params) {

    const judgRef = useRef(null);
    const [isModalOpen, setIsModalOpen] = useState(false);
    let [selectQuest, setSelectQuest] = useState();
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
            title: 'ID',
            dataIndex: 'ID',
            key: 'ID'
        },
        {
            title: '地址',
            dataIndex: 'address',
            key: 'address',
            render: (address) => (
              <p onClick={() => copyToClipboard(address)} style={{cursor: "pointer", textDecoration: "underline"}}>{address.substring(0,5) + "..." + address.substring(38,42)}</p>
            )
        },
        {
            title: '挑战编号',
            dataIndex: 'token_id',
            key: 'token_id',
            render: (tokenId) => (
                <a className="underline" href={`${host}/quests/${tokenId}`} target="_blank">{tokenId}</a>
            )
        },
        {
            title: '提交时间',
            dataIndex: 'UpdatedAt',
            key: 'UpdatedAt',
            render: (UpdatedAt) => (
                UpdatedAt.replace("T", " ").split(".")[0]
            )
        },
        {
            title: '处理时间',
            dataIndex: 'open_quest_review_time',
            key: 'open_quest_review_time',
            render: (open_quest_review_time) => (
                open_quest_review_time.indexOf("0001-01-01T") === -1 ?
                <p>{open_quest_review_time.replace("T", " ").split(".")[0]}</p>
                :"-"
            )
        },
        {
            title: `状态:${status === 1 ? "待处理" : "已处理"}`,
            dataIndex: 'open_quest_review_status',
            key: 'status',
            filters: [
                { text: '待处理', value: 1 },
                { text: '已处理', value: 2 }
            ],
            filterMultiple: false,
            filteredValue: [status],
            render: (status) => (
                <p style={{
                    color: status === 2 ? "#35D6A6" : "000"
                }}>{status === 2 ? "已处理" : "待处理"}</p>
            )
        },
        {
            title: '操作',
            key: 'action',
            render: (_, quest) => (
                
                <Button
                  type="link" 
                  onClick={() => OpenJudgModal(quest)}
                >{quest.open_quest_review_status === 1 ? "评分" : "查看"}</Button>
            ),
        }
    ];

    // 判题弹窗
    function OpenJudgModal(quest) {
        selectQuest = quest;
        setSelectQuest({...selectQuest});
        setIsModalOpen(true);
    }
    
    // 获取列表
    function getList(page) {
        if (page) {
            pageConfig.page = page;
            setPageConfig({...pageConfig});
        }
        getUserOpenQuestList({
            open_quest_review_status: status,
            ...pageConfig
        })
        .then(res => {
            const list = res.data.list;
            data = list ? list : [];
            // 添加key
            data.forEach(ele => {
                ele.key = ele.ID
            })
            setData([...data]);
            pageConfig.total = res.data.total;
            setPageConfig({...pageConfig});
        })
    }

    function handleOk() {
        judgRef.current.confirm();
    }

    function onFinish() {
        setIsModalOpen(false)
        getList()
    }

    function init() {
        pageConfig.page += 1;
        setPageConfig({...pageConfig});
        getList();
    }

    useEffect(() => {
        init();
    },[])

    return (
        <div className="judg">
            <Modal
                width={1177}
                open={isModalOpen}
                className="judg-modal"
                onCancel={() => {setIsModalOpen(false)}}
                onOk={handleOk}
            >
                <ChallengeJudgPage ref={judgRef} selectQuest={selectQuest} onFinish={onFinish} />
            </Modal>
            <Table
                columns={columns} 
                dataSource={data} 
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