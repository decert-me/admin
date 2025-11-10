import { useEffect, useState } from "react";
import { Select, Table, Tag, Button, message } from "antd";
import { SettingOutlined } from "@ant-design/icons";
import { getTagList } from "../../request/api/userTags";
import { getBootcampChallengeStatistics, getEnabledBootcampChallenges } from "../../request/api/quest";
import ManageChallengesModal from "./ManageChallengesModal";

export default function BootcampChallengeStatistics() {
  const [selectedTagId, setSelectedTagId] = useState(""); // 当前选中的标签ID
  const [selectedTagName, setSelectedTagName] = useState(""); // 当前选中的标签名称
  const [availableTags, setAvailableTags] = useState([]); // 可用标签列表 {id, name}
  const [userData, setUserData] = useState([]); // 用户数据
  const [loading, setLoading] = useState(false);
  const [challenges, setChallenges] = useState([]); // 动态挑战列表
  const [challengeUuidMap, setChallengeUuidMap] = useState(new Map()); // 挑战UUID映射
  const [manageModalVisible, setManageModalVisible] = useState(false); // 管理弹窗显示状态

  // 获取可用的S标签（S2-S8）
  useEffect(() => {
    fetchAvailableTags();
    fetchEnabledChallenges();
  }, []);

  const fetchEnabledChallenges = async () => {
    try {
      const res = await getEnabledBootcampChallenges();
      if (res.code === 0) {
        const enabledChallenges = res.data || [];
        const titles = enabledChallenges.map(c => c.title);
        const uuidMap = new Map();
        enabledChallenges.forEach(c => {
          uuidMap.set(c.title, c.uuid);
        });
        setChallenges(titles);
        setChallengeUuidMap(uuidMap);
      } else {
        message.error(res.msg || "获取挑战列表失败");
      }
    } catch (error) {
      console.error("获取挑战列表失败:", error);
      message.error("获取挑战列表失败");
    }
  };

  const fetchAvailableTags = async () => {
    try {
      const res = await getTagList({
        page: 1,
        pageSize: 100
      });

      if (res.code === 0) {
        // 筛选出S2-S8的标签
        const sTags = (res.data?.list || [])
          .filter(tag => /^S[2-8]$/.test(tag.name))
          .map(tag => ({
            id: tag.id,
            name: tag.name
          }))
          .sort((a, b) => {
            const numA = parseInt(a.name.substring(1));
            const numB = parseInt(b.name.substring(1));
            return numB - numA; // 降序排列，最大的在前
          });

        setAvailableTags(sTags);
        // 默认选择第一个（最大的S值）
        if (sTags.length > 0) {
          setSelectedTagId(sTags[0].id);
          setSelectedTagName(sTags[0].name);
        }
      } else {
        message.error(res.msg || "获取标签列表失败");
      }
    } catch (error) {
      console.error("获取标签列表失败:", error);
      message.error("获取标签列表失败");
    }
  };

  // 从用户标签中提取分组信息（如S8-G1）
  const extractGroup = (tags) => {
    if (!tags) return "未分组";
    // 匹配 S8-G1, S8-G2 这样的分组标签
    const groupMatch = tags.match(new RegExp(`${selectedTagName}-G\\d+`));
    return groupMatch ? groupMatch[0] : "未分组";
  };

  // 获取用户数据和挑战完成情况
  const fetchData = async () => {
    if (!selectedTagId) return;

    setLoading(true);
    try {
      // 调用新的API，一次性获取所有数据
      const res = await getBootcampChallengeStatistics({
        tag_id: selectedTagId,
        challenges: challenges
      });

      if (res.code !== 0) {
        message.error(res.msg || "获取数据失败");
        setLoading(false);
        return;
      }

      const rawData = res.data || [];

      // 按用户分组数据
      const userMap = new Map();

      rawData.forEach(item => {
        const userId = item.user_id;

        if (!userMap.has(userId)) {
          userMap.set(userId, {
            key: item.address,
            address: item.address,
            name: item.name || item.address.slice(0, 8),
            group: extractGroup(item.tags),
            challengeStatus: {}
          });
        }

        const user = userMap.get(userId);

        // 只有当title不为空时才处理挑战状态
        if (item.title && challenges.includes(item.title)) {
          // status: 0=未提交, 1=未通过, 2=通过
          user.challengeStatus[item.title] = {
            status: item.status
          };
        }
      });

      const usersData = Array.from(userMap.values());

      // 按分组排序：先按分组名排序，未分组的放最后
      usersData.sort((a, b) => {
        // 未分组的排到最后
        if (a.group === "未分组" && b.group !== "未分组") return 1;
        if (a.group !== "未分组" && b.group === "未分组") return -1;
        if (a.group === "未分组" && b.group === "未分组") return 0;

        // 提取分组编号进行排序 (S8-G1 -> 1, S8-G2 -> 2)
        const getGroupNum = (group) => {
          const match = group.match(/G(\d+)/);
          return match ? parseInt(match[1]) : 0;
        };

        const groupNumA = getGroupNum(a.group);
        const groupNumB = getGroupNum(b.group);

        return groupNumA - groupNumB;
      });

      setUserData(usersData);
    } catch (error) {
      console.error("获取数据失败:", error);
      message.error("获取数据失败");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (selectedTagId && challenges.length > 0) {
      fetchData();
    }
  }, [selectedTagId, challenges]);

  // 渲染挑战完成状态标签
  const renderChallengeStatus = (challengeStatus) => {
    if (!challengeStatus) {
      return <Tag color="default">未提交</Tag>;
    }

    const { status } = challengeStatus;

    // status: 0=未提交, 1=未通过, 2=通过
    if (status === 2) {
      return <Tag color="success">通过</Tag>;
    } else if (status === 1) {
      return <Tag color="warning">未通过</Tag>;
    } else {
      return <Tag color="default">未提交</Tag>;
    }
  };

  // 构建表格列
  const buildColumns = () => {
    const baseColumns = [
      {
        title: "分组",
        dataIndex: "group",
        key: "group",
        width: 100,
        fixed: "left",
        render: (text, record, index) => {
          // 计算合并的行数
          const groupUsers = userData.filter(u => u.group === text);
          const firstIndex = userData.findIndex(u => u.group === text);

          if (index === firstIndex) {
            return {
              children: text || "-",
              props: {
                rowSpan: groupUsers.length,
              },
            };
          }
          return {
            children: text || "-",
            props: {
              rowSpan: 0,
            },
          };
        },
      },
      {
        title: "姓名/挑战",
        dataIndex: "name",
        key: "name",
        width: 120,
        fixed: "left",
        ellipsis: true,
      },
    ];

    // 添加挑战列
    const challengeColumns = challenges.map((challenge, index) => {
      const uuid = challengeUuidMap.get(challenge);
      const linkUrl = uuid ? `${process.env.REACT_APP_LINK_URL || "https://decert.me"}/quests/${uuid}` : null;

      return {
        title: linkUrl ? (
          <a target="_blank" rel="noopener noreferrer" href={linkUrl}>
            {challenge}
          </a>
        ) : (
          challenge
        ),
        dataIndex: `challenge_${index}`,
        key: `challenge_${index}`,
        width: 160,
        render: (_, record) => {
          const status = record.challengeStatus[challenge];
          return renderChallengeStatus(status);
        },
      };
    });

    return [...baseColumns, ...challengeColumns];
  };

  return (
    <div className="bootcamp-challenge-statistics">
      <div className="tabel-title">
        <h2>训练营挑战统计</h2>
      </div>

      <div style={{ marginBottom: 16, display: 'flex', alignItems: 'center', gap: '12px' }}>
        <span style={{ marginRight: 8 }}>选择训练营:</span>
        <Select
          value={selectedTagId}
          onChange={(tagId) => {
            setSelectedTagId(tagId);
            const tag = availableTags.find(t => t.id === tagId);
            if (tag) {
              setSelectedTagName(tag.name);
            }
          }}
          style={{ width: 200 }}
          options={availableTags.map(tag => ({
            label: tag.name,
            value: tag.id
          }))}
        />
        <Button
          icon={<SettingOutlined />}
          onClick={() => setManageModalVisible(true)}
        >
          管理训练营挑战
        </Button>
      </div>

      <div className="sticky-table-wrapper">
        <Table
          columns={buildColumns()}
          dataSource={userData}
          loading={loading}
          pagination={false}
          scroll={{ x: 1800 }}
          bordered
        />
      </div>

      <ManageChallengesModal
        visible={manageModalVisible}
        onClose={() => setManageModalVisible(false)}
        onUpdate={() => {
          fetchEnabledChallenges();
          if (selectedTagId) {
            fetchData();
          }
        }}
      />

      <style jsx>{`
        .sticky-table-wrapper {
          position: relative;
        }

        .sticky-table-wrapper :global(.ant-table-body) {
          overflow-x: auto !important;
          overflow-y: auto !important;
        }
      `}</style>
    </div>
  );
}
