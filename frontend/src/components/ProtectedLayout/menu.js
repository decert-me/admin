import {
    MenuOutlined,
    BookOutlined,
    FolderOutlined,
    HomeOutlined,
    TagsOutlined,
    ProfileOutlined,
    PartitionOutlined
  } from '@ant-design/icons';

export const menu = [
    {
        label: "首页",
        key: "home",
        icon: <HomeOutlined />,
    },
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
            }
        ]
    }
]