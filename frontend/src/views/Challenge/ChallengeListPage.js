import { Button, Checkbox, Input, Modal, Popconfirm, Space, Spin, Switch, Table, Tooltip, message } from "antd";
import {
  ArrowLeftOutlined,
  SearchOutlined
} from '@ant-design/icons';
import React, { useEffect, useRef, useState } from "react";
import "./index.scss";
import { addQuestToCollection, deleteQuest, getCollectionQuestList, getQuestCollectionAddList, getQuestList, updateCollectionQuestSort, updateQuest, updateQuestStatus } from "../../request/api/quest";
import { format } from "../../utils/format";
import { useLocation, useNavigate, useParams } from "react-router-dom";
import { DndContext, PointerSensor, useSensor, useSensors } from '@dnd-kit/core';
import { restrictToVerticalAxis } from '@dnd-kit/modifiers';
import {
  arrayMove,
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { useRequest } from 'ahooks';
import InfiniteScroll from "../../components/InfiniteScroll";
import CustomLoading from "../../components/InfiniteScroll/CustomLoading";

const location = window.location.host;
const isTest = ((location.indexOf("localhost") !== -1) || (location.indexOf("192.168.1.10") !== -1)) ? false : true;
const host = isTest ? "https://decert.me" : "http://192.168.1.10:8087";
const opensea = isTest ? "https://opensea.io/assets/matic/0xc8E9cd4921E54c4163870092Ca8d9660e967B53d" : "https://testnets.opensea.io/assets/mumbai/0x66C54CB10Ef3d038aaBA2Ac06d2c25B326be8142"

const Row = (props) => {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: props['data-row-key'],
  });
  const style = {
    ...props.style,
    transform: CSS.Transform.toString(
      transform && {
        ...transform,
        scaleY: 1,
      },
    ),
    transition,
    cursor: 'move',
    ...(isDragging
      ? {
          position: 'relative',
          zIndex: 9999,
        }
      : {}),
  };
  return <tr {...props} ref={setNodeRef} style={style} {...attributes} {...listeners} />;
};

export default function ChallengeListPage(params) {
  
    const { formatTimestamp } = format();
    const navigateTo = useNavigate();
    const location = useLocation();
    const { id } = useParams();
    const scrollRef = useRef(null);

    const [isModalOpen, setIsModalOpen] = useState(false);
    const [isOk, setIsOk] = useState(false);
    
    let [search_key, setSearch_key] = useState("");    //  搜索
    let [data, setData] = useState([]);
    let [isChange, setIsChange] = useState();
    let [modalPage, setModalPage] = useState({
      page: 0, pageSize: 20, total: 0
    });
    let [modalData, setModalData] = useState([]);
    let [checkedList, setCheckedList] = useState([]);
    let [pageConfig, setPageConfig] = useState({
        page: 0, pageSize: 10, total: 0
    });

    const { run } = useRequest(changeSearch, {
      debounceWait: 500,
      manual: true,
    });

    const columns = [
        {
          title: '权重',
          dataIndex: 'sort',
          render: (sort, record, index) => (
            <p>{sort}</p>
          )
        },
        {
          title: '挑战编号',
          dataIndex: 'tokenId',
          render: (tokenId, quest) => (
            id ?
            <Tooltip title="点此管理挑战">
              <a className="underline" href={`/dashboard/challenge/list?tokenId=${quest.tokenId}`} target="">{tokenId}</a>
            </Tooltip>
            :
            <a className="underline" href={`${host}/quests/${tokenId}`} target="_blank">{tokenId}</a>
          )
        },
        {
          title: 'NFT',
          dataIndex: 'metadata',
          render: ({image}, quest) => (
            quest.claim_num !== 0 ?
            <a href={`${opensea}/${quest.tokenId}`} target="_blank">
                <img src={image.replace("ipfs://", "https://ipfs.decert.me/")} alt="" style={{height: "53px"}} />
            </a>
            :
            <img src={image.replace("ipfs://", "https://ipfs.decert.me/")} alt="" style={{height: "53px"}} />
          )
        },
        {
          title: '标题',
          dataIndex: 'title',
          render: (title, quest) => (
            <Tooltip title="查看挑战网页效果">
              <a className="underline" href={`${host}/quests/${quest.tokenId}`} target="_blank">{title}</a>
            </Tooltip>
          )
        },
        {
            title: '发布者',
            dataIndex: 'creator',
            render: (creator) => (
                <a className="underline" href={`${host}/user/${creator}`} target="_blank">{creator.substring(0,5) + "..." + creator.substring(38,42)}</a>
            )
        },
        {
            title: '上架状态',
            dataIndex: 'status',
            render: (status, quest) => (
                <Switch
                    checkedChildren="已上架" 
                    unCheckedChildren="待上架" 
                    checked={status == 1 ? true : false}
                    onChange={(checked) => handleChangeStatus({status: checked ? 1 : 2, id: quest.id}, quest.key)}
                />
            )
        },
        {
            title: '难度',
            dataIndex: 'metadata',
            render: ({attributes}) => (
              <p>{attributes.difficulty === 0 ? "简单" : attributes.difficulty === 1 ? "中等" : attributes.difficulty === 2 ? "困难" : "/"}</p>
            )
        },
        {
            title: '时长',
            dataIndex: 'quest_data',
            render: ({estimateTime}) => (
              <p>{estimateTime ? (estimateTime / 60) + "min" : "/"}</p>
            )
        },
        {
            title: '铸造数量',
            dataIndex: 'claim_num',
            render: (claim_num) => (
              <p>{claim_num}</p>
            )
        },
        {
            title: '挑战人次',
            dataIndex: 'challenge_num',
            render: (challenge_num) => (
              <p>{challenge_num}次</p>
            )
        },
        {
            title: '创建时间',
            dataIndex: 'addTs',
            render: (addTs) => (
              <p>{formatTimestamp(addTs * 1000)}</p>
            )
        },
        {
            title: 'Action',
            key: 'action',
            render: (_, quest) => (
              <Space size="middle">
                {
                  id ?
                  // 移除合辑
                  <Popconfirm
                    title="移出合辑"
                    description="确定要移出该挑战吗?"
                    onConfirm={() => updateT(Number(quest.id))}
                    okText="确定"
                    cancelText="取消"
                  >
                    <Button 
                    type="link" 
                    className="p0"
                  >移出合辑</Button>
                  </Popconfirm>
                  :
                  <>
                    <Button 
                      type="link" 
                      className="p0"
                      onClick={() => navigateTo(`/dashboard/challenge/modify/${quest.id}/${quest.tokenId}`)}
                    >编辑</Button>
                    <Popconfirm
                      title="删除挑战"
                      description="确定要删除该挑战吗?"
                      onConfirm={() => deleteT(Number(quest.id))}
                      okText="确定"
                      cancelText="取消"
                    >
                      <Button 
                      type="link" 
                      className="p0"
                    >删除</Button>
                    </Popconfirm>
                  </>
                }
              </Space>
            ),
        }
    ];

    id && columns.splice(0,1);

    // 移出合辑
    async function updateT(selectId) {
      const selectData = data.filter(e => e.id === selectId)[0];
      const collection_id = selectData.collection_id.filter(e => e != id);
      await updateQuest({
        collection_id, id: selectId
      })
      .then(res => {
        if (res.code === 0) {
            message.success(res.msg);
        }
      })
      .catch(err => {
          message.error(err);
      })
      getList()
    }

    // 上下架
    function handleChangeStatus({status, id: paramsId}, key) {
      if (id) {
        return
      }
      const index = data.findIndex((item) => item.key === key);
      updateQuestStatus({id: paramsId, status})
      .then(res => {
        if (res.code === 0) {
          message.success(res.msg);
          data[index].status = status;
          setData([...data]);
        }
      })
    }

    // 删除教程
    async function deleteT(id) {
        await deleteQuest({id})
        .then(res => {
          if (res.code === 0) {
              message.success(res.msg);
          }
        })
        .catch(err => {
            message.error(err);
        })
        getList()
    }

    async function getList(page) {
        if (page) {
          pageConfig.page = page;
          setPageConfig({...pageConfig});
        }
        // 获取教程列表
        let res;
        if (id) {
          res = await getCollectionQuestList({id: Number(id)});
          getChallenge()
        }else{
          res = await getQuestList({...pageConfig, search_key});
        }

        if (res.code === 0) {
          const list = res.data.list;
          data = list ? list : [];
          // 添加key
          data.forEach(ele => {
            ele.key = ele.id
          })
          setData([...data]);
          pageConfig.total = res.data.total;
          setPageConfig({...pageConfig});
        }else{
            message.success(res.msg);
        }
    }

    const onDragEnd = ({ active, over }) => {
      if (active.id !== over?.id) {
        setData((previous) => {
          const activeIndex = previous.findIndex((i) => i.key === active.id);
          const overIndex = previous.findIndex((i) => i.key === over?.id);
          return arrayMove(previous, activeIndex, overIndex);
        });
        isChange = true;
        setIsChange(isChange);
      }
    };

    const sensors = useSensors(
      useSensor(PointerSensor, {
        activationConstraint: {
          // https://docs.dndkit.com/api-documentation/sensors/pointer#activation-constraints
          distance: 1,
        },
      }),
    );

    function getChallenge(params) {
      modalPage.page += 1;
      setModalPage({...modalPage})

      getQuestCollectionAddList({...modalPage, search_key})
      .then(res => {
        if (res.code === 0) {
          modalPage.total = res.data.total;
          setModalPage({...modalPage});
          const arr = res.data.list;
          const checked = [];
          arr.forEach(e => {
            e.checked = e.collection_id.indexOf(Number(id)) !== -1;
            e.checked && checked.push(e.id);
          })
          const checkboxArr = arr.map(e => {
            console.log(e);
            return {
              label: (
                <div style={{
                  width: "425px",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "space-between"
                }}>
                  <p>{e.tokenId}</p>
                  <p>{e.title}</p>
                  <div style={{
                    width: "100px",
                    height: "100px",
                    display: "flex",
                    alignItems: "center"
                  }}>
                    <img 
                      src={e.metadata.image.replace("ipfs://", "https://ipfs.decert.me/")} 
                      alt="" 
                      style={{
                        maxHeight: "100%",
                        maxWidth: "100%"
                      }} 
                    />
                  </div>
                </div>
              ),
              value: e.id
            }
          })
          modalData = modalData.concat(checkboxArr);
          setModalData([...modalData]);
          checkedList = checkedList.concat(checked);
          setCheckedList([...checkedList]);
          setIsOk(true);
        }
      })
    }

    function changeChecked(params) {
      checkedList = params;
      setCheckedList([...checkedList]);
    }

    function changeCollection() {
      addQuestToCollection({collection_id: Number(id), id: checkedList})
      .then(res => {
        if (res.code === 0) {
          message.success(res.msg);
          pageConfig = {
            page: 0, pageSize: 10, total: 0
          }
          modalPage.page = 0;
          setModalPage({...modalPage})
          modalData = [];
          setModalData([...modalData]);
          init();
          setIsModalOpen(false);
        }else{
          message.error(res.msg);
        }
      })
    }

    function changeSearch(params) {
      search_key = params;
      setSearch_key(search_key);
      modalPage = {
        page: 0, pageSize: 20, total: 0
      }
      modalData = [];
      setModalData([...modalData]);
      pageConfig = {
        page: 0, pageSize: 10, total: 0
      }
      init();
    }

    function init(params) {
        pageConfig.page += 1;
        setPageConfig({...pageConfig});
        getList()
    }

    useEffect(() => {
      pageConfig = {
        page: 0, pageSize: 10, total: 0
      }
      setPageConfig({...pageConfig});
      if (location.search) {
        let serch = new URLSearchParams(location.search);
        search_key = serch.get("tokenId");
        setSearch_key(search_key);
      }
      init();
    },[location])

    useEffect(() => {
      if (isChange) {
        setIsChange(false);
        updateCollectionQuestSort({collection_id: Number(id), id: data.map(e => Number(e.id))})
        .then(res => {
          if (res.code === 0) {
            message.success(res.msg);
          }else{
            message.error(res.msg);
          }
        })
      }
    },[data])
    
    return (
        <div className="challenge" key={location.pathname}>
            <div className="tabel-title">
              {
                id ? 
                  <Space style={{cursor: "pointer"}} onClick={() => navigateTo("/dashboard/challenge/compilation")}>
                    <ArrowLeftOutlined /><h2><span style={{color: "#999999"}}>合辑管理</span>/{decodeURIComponent(location.search.split("=")[1])}</h2>
                  </Space>
                :
                  <h2>挑战列表</h2>
              }
                <Space size="large">
                      {
                        id ? 
                        <Button 
                          type="primary"
                          onClick={() => setIsModalOpen(true)}
                          disabled={!isOk}
                        >
                          添加挑战
                        </Button>
                        :
                        <Input prefix={<SearchOutlined />} onChange={(e) => run(e.target.value)} />
                      }
                </Space>
            </div>
            <Modal 
              title="添加挑战"
              open={isModalOpen} 
              onOk={() => changeCollection()} 
              onCancel={() => setIsModalOpen(false)}
              closeIcon={null}
              okText="添加"
              cancelText="取消"
              bodyStyle={{
                textAlign: "center"
              }}
            >
              <Input prefix={<SearchOutlined />} onChange={(e) => run(e.target.value)} style={{marginBottom: "10px"}} />
              <div 
                ref={scrollRef}
                style={{
                  overflowY: "auto",
                  height: "300px",
                }}
              >
                <Checkbox.Group
                  options={modalData}
                  style={{
                    display: "flex",
                    flexDirection: "column",
                    gap: "10px"
                  }}
                  // defaultValue={checkedList}
                  value={checkedList}
                  onChange={e => changeChecked(e)}
                />
                {
                  modalPage.page * modalPage.pageSize < modalPage.total &&
                    <InfiniteScroll
                        scrollRef={scrollRef}
                        func={getChallenge}
                        isCustom={true}
                        components={(
                          <CustomLoading className="sbt-loading" />
                      )}
                    />
                }
              </div>
            </Modal>
            {
              id ? 
              <DndContext sensors={sensors} modifiers={[restrictToVerticalAxis]} onDragEnd={onDragEnd}>
              <SortableContext
                // rowKey array
                items={data.map((i) => i.key)}
                strategy={verticalListSortingStrategy}
              >
                <Table
                  components={{
                    body: {
                      row: Row,
                    },
                  }}
                  rowKey="key"
                  columns={columns}
                  dataSource={data}
                  pagination={false}
                />
              </SortableContext>
            </DndContext>
              :
              <Table 
                  // rowSelection={rowSelection} 
                  columns={columns} 
                  dataSource={data} 
                  rowClassName={(record) => record.top && "toTop"}
                  pagination={{
                      current: pageConfig.page, 
                      total: pageConfig.total, 
                      pageSize: pageConfig.pageSize, 
                      onChange: (page) => getList(page)
                  }} 
              />
            }
        </div>
    )
}