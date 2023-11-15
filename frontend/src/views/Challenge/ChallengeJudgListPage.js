import { Button, Table } from "antd"
import { useEffect, useState } from "react";
import { getUserOpenQuestList } from "../../request/api/judgment";
import { useNavigate } from "react-router-dom";



export default function ChallengeJudgListPage(params) {

    const navigateTo = useNavigate();
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
              <p>{address.substring(0,5) + "..." + address.substring(38,42)}</p>
            )
        },
        {
            title: '挑战ID',
            dataIndex: 'token_id',
            key: 'token_id',
        },
        {
            title: '提交时间',
            dataIndex: 'CreatedAt',
            key: 'CreatedAt',
            render: (CreatedAt) => (
                CreatedAt.replace("T", " ").split(".")[0]
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
            title: `状态:${status === 1 ? "待处理" : status === 2 ? "已处理" : status === 3 ? "空投失败" : "全部"}`,
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
                }}>{status === 2 ? "空投失败" : "待处理"}</p>
            )
        },
        {
            title: '分数',
            dataIndex: 'open_quest_score',
            key: 'open_quest_score',
            render: (open_quest_score, record) => (
                record.open_quest_review_status === 1 ? 
                "-"
                : open_quest_score + "分"
            )
        },
        {
            title: '操作',
            key: 'action',
            render: (_, quest) => (
                <Button
                  type="link" 
                  onClick={() => navigateTo(`/dashboard/challenge/openquest/judg/${quest.ID}`)}
                >编辑</Button>
            ),
        }
    ];
    
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