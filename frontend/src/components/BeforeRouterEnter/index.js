import { Route, Routes } from 'react-router-dom';
import { Redirect } from '../Redirect';
import { ProtectedLayout } from '../ProtectedLayout';
import LoginPage from '../../views/Login';
import ProfilePage from '../../views/Profile';
import SettingsPage from '../../views/Settings';
import AuthGuard from '../AuthGuard';
import HomePage from '../../views/Home';
import { AirdropListPage } from '../../views/Airdrop';
import { 
    TutorialsAddPage, 
    TutorialsBuildLogPage, 
    TutorialsBuildPage, 
    TutorialsListPage, 
    TutorialsModifyPage, 
} from '../../views/Tutorials';
import { 
    TagsAddPage, 
    TagsModifyPage, 
    TagsPage 
} from '../../views/Tags';
import { 
    ChallengeAddPage, 
    ChallengeCompilationPage, 
    ChallengeListPage, 
    ChallengeModifyPage,
    ChallengeCompilationModifyPage
} from '../../views/Challenge';
import { PersonelEditPage, PersonelListPage } from '../../views/Personel';



export default function BeforeRouterEnter(params) {
    
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
                {/* <Route 
                path="tutorials/build" 
                element={<TutorialsBuildPage />} 
                />
                <Route 
                path="tutorials/buildlog/:id" 
                element={<TutorialsBuildLogPage />} 
                /> */}
                

                {/* 标签 */}
                <Route 
                path="tags" 
                element={<TagsPage />} 
                />
                <Route 
                path="tags/add" 
                element={<TagsAddPage />} 
                />
                <Route 
                path="tags/modify/:type/:id" 
                element={<TagsModifyPage />} 
                />


                {/* 挑战 */}
                <Route 
                path="challenge/list" 
                element={<ChallengeListPage />} 
                />
                <Route 
                path="challenge/list/:id" 
                element={<ChallengeListPage />} 
                />
                <Route 
                path="challenge/modify/:id/:tokenId" 
                element={<ChallengeModifyPage />} 
                />
                <Route 
                path="challenge/compilation"
                element={<ChallengeCompilationPage />} 
                />
                <Route 
                path="challenge/compilation/modify/:id"
                element={<ChallengeCompilationModifyPage />}
                />
                <Route 
                path="challenge/add" 
                element={<ChallengeAddPage />} 
                />

                {/* 空投 */}
                <Route 
                path="airdrop/list" 
                element={<AirdropListPage />} 
                />

                {/* 人员管理 */}
                <Route 
                path="personnel/list" 
                element={<PersonelListPage />} 
                />
                <Route 
                path="personnel/:type" 
                element={<PersonelEditPage />} 
                />
            </Route>
        </Routes>
    )
}