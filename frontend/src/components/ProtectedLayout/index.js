import { Navigate, useLocation, useNavigate, useOutlet } from "react-router-dom";
import { useAuth } from "../../hooks/useAuth";
import { useEffect, useState } from "react";
import {
    MenuFoldOutlined,
    MenuUnfoldOutlined,
    GithubOutlined,
    CopyOutlined
} from '@ant-design/icons';
import { Layout, Menu, Button, theme, Dropdown, Space, message } from 'antd';
import "./index.scss";
import { menu } from "./menu";
import { useDisconnect } from "wagmi";
const { Header, Sider, Content } = Layout;


export const ProtectedLayout = () => {

  const { user, auth, logout } = useAuth();
  const outlet = useOutlet();
  const location = useLocation();
  const navigateTo = useNavigate();
  const { disconnect } = useDisconnect();

  const { token: { colorBgContainer }, } = theme.useToken();
  const [collapsed, setCollapsed] = useState(false);

  if (!user) {
    return <Navigate to="/" />;
  }

  const dropdownMenu = [
    {
      key: '1',
      label: (
        <Space onClick={copy}>
          {user.address.substring(0,5) + "..." + user.address.substring(38,42)}
          <CopyOutlined />
        </Space>
      ),
    },
    {
      key: '2',
      label: (
        <p style={{textAlign: "center"}} onClick={() => navigateTo(`/dashboard/personnel/edit?id=${user.id}`)}>编辑资料</p>
      )
    },
    {
      key: '3',
      label: (
        <p style={{color: "#ff4d55", textAlign: "center"}} onClick={goDisconnect}>
          断开连接
        </p>
      ),
    }
  ]
//   自动生成侧边栏
  // function menu() {
    // const arr = []
    // TODO: menu生成
    // 鉴权生成
    // auth.forEach((element, index) => {
    //     arr.push({
    //         key: element,
    //         icon: "",
    //         label: element
    //     })
    // });
    
  //   return arr
  // }

  function defaultSelectedKeys() {
    // 获取当前路由 menu
    const path = location.pathname;
    let key = "";
    let open = "";
    menu.forEach(e => {
      if (e.children) {
        const res = e.children.filter(ele => path.includes(ele.key))
        if (res.length !== 0) {
          key = res[0].key
          open = e.key;
        }
      }else{
        if (path.includes(e.key)) {
          key = e.key
        }
      }
    })
    return {
      key,
      open
    }
  }

  // 断开连接
  function goDisconnect() {
    disconnect();
    logout();
  }

  // 复制地址
  function copy(params) {
    navigator.clipboard.writeText(user.address)
    .then(() => {
      message.success("复制成功!")
    })
    .catch(err => {
      message.error("复制失败!")
    });
  }
  
  return (
    <Layout style={{height: "100vh"}} className="main">
        {/* 侧边栏 */}
      <Sider trigger={null} collapsible collapsed={collapsed} className="main-sidebar">
        <div className="demo-logo-vertical" >
            <img src={require("../../assets/logo.png")} alt="" />
        </div>
        <Menu
          theme="dark"
          mode="inline"
          defaultSelectedKeys={[defaultSelectedKeys().key]}  //   初始值
          defaultOpenKeys={[defaultSelectedKeys().open]}  //  默认展开
          items={menu}
          onClick={(item) => {navigateTo(item.key)}}  //   跳转
        />
      </Sider>
      {/* 正文 */}
      <Layout className="main-content">
        {/* 正文导航栏 */}
        <Header
          style={{
            padding: 0,
            background: colorBgContainer,
          }}
          className="content-navbar"
        >
            {/* menu闭合按钮 */}
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            style={{
              fontSize: '16px',
              width: 64,
              height: 64,
            }}
          />

          {/* 控制台 */}
          <div className="operate">
            {!!user && (
              // TODO: user展示 ===>
              <div className="user">
                <Dropdown
                  menu={{
                    items: dropdownMenu,
                  }}
                >
                  {
                    user.headerImg ? 
                      <div className="avatar">
                          <img src={process.env.REACT_APP_BASE_URL+"/"+user.headerImg} alt="" />
                      </div>
                      :
                      <div className="avatar">{user.username[0].toUpperCase()}</div>
                  }
                </Dropdown>
              </div>
            )}
          </div>
        </Header>
        {/* 正文内容 */}
        <Content
          style={{
            margin: '24px 16px',
            padding: 24,
            minHeight: 280,
            background: colorBgContainer,
            overflow: "auto"
          }}
          className="content-info"
        >
          {outlet}
        </Content>
      </Layout>
    </Layout>
  );
};
