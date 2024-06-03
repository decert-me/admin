import {
    MenuOutlined,
    BookOutlined,
    FolderOutlined,
    HomeOutlined,
    UserOutlined,
    TagsOutlined,
    ProfileOutlined,
    PartitionOutlined,
    GiftOutlined,
    TeamOutlined,
    AuditOutlined
  } from '@ant-design/icons';

export const menu = [
    // {
    //     label: "首页",
    //     key: "home",
    //     icon: <HomeOutlined />,
    // },
    {
        label: "教程管理",
        key: "tutorials/list",
        icon: <BookOutlined />,
    },
    {
        label: "标签管理",
        key: "tags",
        icon: <TagsOutlined />,
    },
    {
        label: "挑战管理",
        // key: "challenge",
        icon: <ProfileOutlined />,
        children: [
            {
                label: "挑战列表",
                key: "challenge/list",
                icon: <ProfileOutlined />,
            },
            {
                label: "挑战合辑",
                key: "challenge/compilation",
                icon: <PartitionOutlined />,
            },
            {
                label: "开放题评分",
                key: "challenge/openquest",
                icon: <AuditOutlined />,
            }
        ]
    },
    {
        label: "用户管理",
        // key: "challenge",
        icon: <TeamOutlined />,
        children: [
            {
                label: "用户列表",
                key: "user/list",
                icon: <ProfileOutlined />,
            },
            {
                label: "用户标签",
                key: "user/tag",
                icon: <TagsOutlined />,
            }
        ]
    },
    {
        label: "空投管理",
        key: "airdrop/list",
        icon: <GiftOutlined />,
    },
    {
        label: "成员管理",
        key: "personnel/list",
        icon: <UserOutlined />,
    },
]