import React from 'react';
import { Routes, Route, Link } from 'react-router-dom';
import { LoginForm, PageForm } from '@ant-design/pro-components';
import Login from './view/login'

// Outlet
const Home = () => <h2>🏠 首页</h2>;
const About = () => <h2>ℹ️ 关于页面</h2>;
const Users = () => <h2>👥 用户列表</h2>;
// const Login = () => <h2>🔑 登录</h2>;


function App() {
  return (
    <div style={{ fontFamily: 'sans-serif', padding: '20px' }}>
      <nav>
        {/* 使用 Link 组件进行导航，避免页面刷新 */}
        <ul style={{ listStyle: 'none', display: 'flex', gap: '10px' }}>
          <li><Link to="/">登录</Link></li>
          <li><Link to="/home">首页</Link></li>
          <li><Link to="/about">关于</Link></li>
          <li><Link to="/users">用户</Link></li>
        </ul>
      </nav>
      <hr />
      {/* 路由匹配区域 */}
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/home" element={<Home />} />
        <Route path="/about" element={<About />} />
        <Route path="/users" element={<Users />} />
      </Routes>
    </div>
  );
}

export default App;