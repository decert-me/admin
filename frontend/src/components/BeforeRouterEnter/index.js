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
    ChallengeCompilationModifyPage,
    ChallengeJudgPage,
    ChallengeJudgListPage
} from '../../views/Challenge';
import { PersonelEditPage, PersonelListPage } from '../../views/Personel';
import ChallengeAnswerListPage from '../../views/Challenge/ChallengeAnswerListPage';
import UserTagsPage from '../../views/User/UserTagsPage';
import UserListPage from '../../views/User/UserListPage';
import UserTagInfoPage from '../../views/User/UserTagInfoPage';
import UserTagAddPage from '../../views/User/UserTagAddPage';
import UserTagModify from '../../views/User/UserTagModify';
import ChallengerListPage from '../../views/Challenge/ChallengerListPage';
import UserTagUserPage from '../../views/User/UserTagUserPage';



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
                path="challenge/answer/list/:tokenId" 
                element={<ChallengeAnswerListPage />} 
                />
                <Route 
                path="challenge/answer/list" 
                element={<ChallengeAnswerListPage />} 
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
                <Route 
                path="challenge/openquest" 
                element={<ChallengeJudgListPage />} 
                />
                <Route 
                path="challenge/openquest/judg/:id" 
                element={<ChallengeJudgPage />} 
                />
                <Route 
                path="challenge/challenge/list" 
                element={<ChallengerListPage />} 
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

                {/* 用户管理 */}
                <Route 
                path="user/list/:tagid" 
                element={<UserTagUserPage />} 
                />
                <Route 
                path="user/list" 
                element={<UserListPage />} 
                />
                <Route 
                path="user/tag" 
                element={<UserTagsPage />} 
                />
                <Route 
                path="user/tag/add" 
                element={<UserTagInfoPage />} 
                />
                <Route 
                path="user/tag/modify/:id" 
                element={<UserTagInfoPage />} 
                />
                <Route 
                path="user/tag/adduser/:id" 
                element={<UserTagAddPage />} 
                />
                <Route 
                path="user/tag/modifyuser/:address" 
                element={<UserTagModify />} 
                />
            </Route>
        </Routes>
    )
}