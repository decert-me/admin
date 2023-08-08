import { Route, Routes } from 'react-router-dom';
import './styles/App.css';
import { Redirect } from './components/Redirect';
import { ProtectedLayout } from './components/ProtectedLayout';
import LoginPage from './views/Login';
import ProfilePage from './views/Profile';
import SettingsPage from './views/Settings';
import AuthGuard from './components/AuthGuard';
import TutorialsPage from './views/Tutorials';

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
          path="tutorials" 
          element={<TutorialsPage />} 
        />
      </Route>
    </Routes>
  );
}

export default App;
