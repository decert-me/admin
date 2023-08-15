import { Route, Routes } from 'react-router-dom';
import './styles/App.css';
import { Redirect } from './components/Redirect';
import { ProtectedLayout } from './components/ProtectedLayout';
import LoginPage from './views/Login';
import ProfilePage from './views/Profile';
import SettingsPage from './views/Settings';
import AuthGuard from './components/AuthGuard';
import HomePage from './views/Home';
import { TutorialsAddPage, TutorialsBuildPage, TutorialsListPage, TutorialsModifyPage, } from './views/Tutorials';
import { TagsAddPage, TagsPage } from './views/Tags';

function App() {
  return (
    <Routes>
      {/* 错误地址重定向 */}
      <Route path="*" element={<Redirect />} />
      <Route path="/login" element={<LoginPage />} />
      <Route 
        path="/dashboard" 
        element={
          <AuthGuard>
            <ProtectedLayout/>
          </AuthGuard>
        }
      >
        <Route 
          path="profile" 
          element={<ProfilePage />} 
        />
        <Route 
          path="settings" 
          element={<SettingsPage />} 
        />
        <Route 
          path="home" 
          element={<HomePage />} 
        />

        {/* 教程 */}
        <Route 
          path="tutorials/list" 
          element={<TutorialsListPage />} 
        />
        <Route 
          path="tutorials/modify/:id" 
          element={<TutorialsModifyPage />} 
        />
        <Route 
          path="tutorials/add" 
          element={<TutorialsAddPage />} 
        />
        <Route 
          path="tutorials/build" 
          element={<TutorialsBuildPage />} 
        />

        {/* 标签 */}
        <Route 
          path="tags" 
          element={<TagsPage />} 
        />
        <Route 
          path="tags/add" 
          element={<TagsAddPage />} 
        />     


      </Route>
    </Routes>
  );
}

export default App;
