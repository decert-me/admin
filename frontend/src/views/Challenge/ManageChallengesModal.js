import { useEffect, useState } from "react";
import { Modal, Table, Switch, Input, InputNumber, message } from "antd";
import { SearchOutlined } from "@ant-design/icons";
import { getBootcampChallengeConfig, updateBootcampChallengeConfig } from "../../request/api/quest";
import { useRequest } from "ahooks";

export default function ManageChallengesModal({ visible, onClose, onUpdate }) {
  const [data, setData] = useState([]);
  const [filteredData, setFilteredData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [searchKey, setSearchKey] = useState("");
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 20,
    showSizeChanger: true,
    showTotal: (total) => `共 ${total} 条`,
  });

  const { run } = useRequest(changeSearch, {
    debounceWait: 500,
    manual: true,
  });

  useEffect(() => {
    if (visible) {
      setSearchKey("");
      setPagination({
        current: 1,
        pageSize: 20,
        showSizeChanger: true,
        showTotal: (total) => `共 ${total} 条`,
      });
      fetchConfig();
    }
  }, [visible]);

  const fetchConfig = async () => {
    setLoading(true);
    try {
      const res = await getBootcampChallengeConfig();
      if (res.code === 0) {
        const list = res.data || [];
        list.forEach(item => {
          item.key = item.quest_id;
        });
        setData(list);
        filterData(list, searchKey);
      } else {
        message.error(res.msg || "获取配置失败");
      }
    } catch (error) {
      console.error("获取配置失败:", error);
      message.error("获取配置失败");
    } finally {
      setLoading(false);
    }
  };

  // 过滤数据
  const filterData = (list, key) => {
    if (!key) {
      setFilteredData(list);
      return;
    }
    const filtered = list.filter(item =>
      item.title && item.title.toLowerCase().includes(key.toLowerCase())
    );
    setFilteredData(filtered);
  };

  // 搜索关键词变化时过滤数据
  useEffect(() => {
    filterData(data, searchKey);
  }, [searchKey, data]);

  // 搜索函数
  function changeSearch(value) {
    setSearchKey(value);
  }

  const handleToggle = async (questId, enabled) => {
    try {
      const res = await updateBootcampChallengeConfig({
        quest_id: questId,
        enabled: enabled,
        display_order: enabled ? (data.find(item => item.quest_id === questId)?.display_order || 0) : 0
      });

      if (res.code === 0) {
        message.success("更新成功");
        // 重新获取数据以保持正确的排序
        fetchConfig();
        // 通知父组件更新
        if (onUpdate) {
          onUpdate();
        }
      } else {
        message.error(res.msg || "更新失败");
      }
    } catch (error) {
      console.error("更新配置失败:", error);
      message.error("更新失败");
    }
  };

  const handleDisplayOrderChange = async (questId, displayOrder) => {
    try {
      const item = data.find(item => item.quest_id === questId);
      if (!item) return;

      const res = await updateBootcampChallengeConfig({
        quest_id: questId,
        enabled: item.enabled,
        display_order: displayOrder || 0
      });

      if (res.code === 0) {
        message.success("更新成功");
        // 重新获取数据以保持正确的排序
        fetchConfig();
        // 通知父组件更新
        if (onUpdate) {
          onUpdate();
        }
      } else {
        message.error(res.msg || "更新失败");
      }
    } catch (error) {
      console.error("更新配置失败:", error);
      message.error("更新失败");
    }
  };

  // 处理表格变化（包括分页）
  const handleTableChange = (paginationConfig) => {
    setPagination({
      ...pagination,
      current: paginationConfig.current,
      pageSize: paginationConfig.pageSize,
    });
  };

  const columns = [
    {
      title: "标题",
      dataIndex: "title",
      key: "title",
      ellipsis: true,
      render: (text, record) => (
        <a
          target="_blank"
          rel="noopener noreferrer"
          href={`${process.env.REACT_APP_LINK_URL || "https://decert.me"}/quests/${record.uuid}`}
        >
          {text}
        </a>
      ),
    },
    {
      title: "训练营挑战",
      dataIndex: "enabled",
      key: "enabled",
      width: 120,
      render: (enabled, record) => (
        <Switch
          checked={enabled}
          onChange={(checked) => handleToggle(record.quest_id, checked)}
        />
      ),
    },
    {
      title: "展示排序",
      dataIndex: "display_order",
      key: "display_order",
      width: 120,
      render: (displayOrder, record) => (
        <InputNumber
          min={0}
          precision={0}
          value={displayOrder || 0}
          disabled={!record.enabled}
          onChange={(value) => handleDisplayOrderChange(record.quest_id, value)}
          style={{ width: '100%' }}
        />
      ),
    },
  ];

  return (
    <Modal
      title="管理训练营挑战"
      open={visible}
      onCancel={onClose}
      footer={null}
      width={800}
      destroyOnClose
      bodyStyle={{ height: '620px', overflow: 'hidden' }}
    >
      <div style={{ marginBottom: 16 }}>
        <Input
          prefix={<SearchOutlined />}
          placeholder="搜索挑战标题"
          onChange={(e) => run(e.target.value)}
          allowClear
        />
      </div>
      <Table
        columns={columns}
        dataSource={filteredData}
        loading={loading}
        pagination={pagination}
        onChange={handleTableChange}
        scroll={{ y: 450 }}
      />
    </Modal>
  );
}
