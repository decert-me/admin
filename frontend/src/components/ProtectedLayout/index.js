import { Navigate, useLocation, useNavigate, useOutlet } from "react-router-dom";
import { useAuth } from "../../hooks/useAuth";
import { useState } from "react";
import {
    MenuFoldOutlined,
    MenuUnfoldOutlined,
    GithubOutlined
} from '@ant-design/icons';
import { Layout, Menu, Button, theme } from 'antd';
import "./index.scss";
const { Header, Sider, Content } = Layout;

export const ProtectedLayout = () => {

  const { user, auth, logout } = useAuth();
  const outlet = useOutlet();
  const location = useLocation();
  const navigateTo = useNavigate();

  const [collapsed, setCollapsed] = useState(false);
  const {
    token: { colorBgContainer },
  } = theme.useToken();

  if (!user) {
    return <Navigate to="/" />;
  }
//   自动生成侧边栏
  function menu() {
    const arr = []
    // TODO: menu生成
    // 鉴权生成
    // auth.forEach((element, index) => {
    //     arr.push({
    //         key: element,
    //         icon: "",
    //         label: element
    //     })
    // });
    
    return arr
  }

  return (
    <Layout style={{height: "100vh"}} className="main">
        {/* 侧边栏 */}
      <Sider trigger={null} collapsible collapsed={collapsed} className="main-sidebar">
        <div className="demo-logo-vertical" >
            <GithubOutlined />
        </div>
        <Menu
          theme="dark"
          mode="inline"
        //   初始值
          defaultSelectedKeys={[location.pathname.split("/").pop()]}
          items={menu()}
        //   跳转
          onClick={(item) => navigateTo(item.key)}
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
                <button key={"logout"} onClick={logout}>
                  <p >Logout</p>
                </button>
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
          }}
          className="content-info"
        >
          {outlet}
        </Content>
      </Layout>
    </Layout>
  );
};
