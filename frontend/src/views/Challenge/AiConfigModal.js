import { Button, Form, Input, List, Modal, Popconfirm, Space, Switch, message, Badge, Statistic, Table, InputNumber, Spin } from "antd";
import { useState, useEffect } from "react";
import { PlusOutlined, EditOutlined, DeleteOutlined, SyncOutlined, HistoryOutlined, QuestionCircleOutlined } from '@ant-design/icons';
import ReactQuill from 'react-quill';
import 'react-quill/dist/quill.snow.css';
import {
    getAiJudgeConfigList,
    createAiJudgeConfig,
    updateAiJudgeConfig,
    deleteAiJudgeConfig,
    toggleAiJudgeConfig,
    toggleAutoGrading,
    getPendingGradeList,
    batchGradePreview,
    submitBatchGrade,
    getAiGradeHistory
} from '../../request/api/aiJudgeConfig';

export default function AiConfigModal({ open, onCancel }) {
    const [form] = Form.useForm();
    const [configs, setConfigs] = useState([]);
    const [editMode, setEditMode] = useState(false);
    const [editId, setEditId] = useState(null);
    const [showForm, setShowForm] = useState(false);
    const [loading, setLoading] = useState(false);
    const [pendingCount, setPendingCount] = useState(0);
    const [isGrading, setIsGrading] = useState(false);

    // 批量判题预览相关状态
    const [previewModalOpen, setPreviewModalOpen] = useState(false);
    const [previewResults, setPreviewResults] = useState([]);
    const [previewLoading, setPreviewLoading] = useState(false);
    const [submitting, setSubmitting] = useState(false);
    const [batchCount, setBatchCount] = useState(5); // 批量判题数量，默认5

    // AI判题历史相关状态
    const [historyModalOpen, setHistoryModalOpen] = useState(false);
    const [historyList, setHistoryList] = useState([]);
    const [historyLoading, setHistoryLoading] = useState(false);
    const [historyPagination, setHistoryPagination] = useState({ current: 1, pageSize: 10, total: 0 });

    // 使用说明弹窗
    const [guideModalOpen, setGuideModalOpen] = useState(false);

    // 富文本编辑器配置
    const modules = {
        toolbar: [
            [{ 'header': [1, 2, false] }],
            ['bold', 'italic', 'underline', 'strike', 'blockquote'],
            [{ 'list': 'ordered' }, { 'list': 'bullet' }],
            ['link', 'code-block'],
            ['clean']
        ],
    };

    const formats = [
        'header',
        'bold', 'italic', 'underline', 'strike', 'blockquote',
        'list', 'bullet',
        'link', 'code-block'
    ];

    // 加载配置列表
    const loadConfigs = async () => {
        try {
            const res = await getAiJudgeConfigList();
            if (res.code === 0) {
                setConfigs(res.data || []);
            }
        } catch (error) {
            message.error('加载配置失败');
        }
    };

    // 加载待判题数量
    const loadPendingCount = async () => {
        try {
            const res = await getPendingGradeList({ page: 1, pageSize: 1 });
            if (res.code === 0) {
                setPendingCount(res.data.total || 0);
            }
        } catch (error) {
            console.error('获取待判题数量失败', error);
        }
    };

    // 首次打开时加载
    useEffect(() => {
        if (open) {
            loadConfigs();
            loadPendingCount();
        }
    }, [open]);

    // 定时刷新待判题数量
    useEffect(() => {
        if (!open) return;

        const interval = setInterval(() => {
            loadPendingCount();
        }, 5000); // 每5秒刷新一次

        return () => clearInterval(interval);
    }, [open]);

    // 新增配置
    const handleAdd = () => {
        setShowForm(true);
        setEditMode(false);
        setEditId(null);
        form.resetFields();
    };

    // 编辑配置
    const handleEdit = (config) => {
        setShowForm(true);
        setEditMode(true);
        setEditId(config.id);
        form.setFieldsValue({
            title: config.title,
            config: config.config
        });
    };

    // 删除配置
    const handleDelete = async (id) => {
        try {
            const res = await deleteAiJudgeConfig({ id });
            if (res.code === 0) {
                message.success('删除成功');
                loadConfigs();
            } else {
                message.error(res.msg || '删除失败');
            }
        } catch (error) {
            message.error('删除失败');
        }
    };

    // 保存配置
    const handleSave = async () => {
        try {
            const values = await form.validateFields();
            setLoading(true);

            let res;
            if (editMode && editId) {
                res = await updateAiJudgeConfig({
                    id: editId,
                    ...values
                });
            } else {
                res = await createAiJudgeConfig(values);
            }

            if (res.code === 0) {
                message.success(editMode ? '更新成功' : '新增成功');
                setShowForm(false);
                form.resetFields();
                loadConfigs();
            } else {
                message.error(res.msg || '操作失败');
            }
        } catch (error) {
            if (error.errorFields) {
                return;
            }
            message.error('操作失败');
        } finally {
            setLoading(false);
        }
    };

    // 切换启用状态
    const handleToggle = async (id, currentEnabled) => {
        try {
            const res = await toggleAiJudgeConfig({ id });
            if (res.code === 0) {
                message.success(currentEnabled ? '已禁用' : '已启用');
                loadConfigs();
            } else {
                message.error(res.msg || '操作失败');
            }
        } catch (error) {
            message.error('操作失败');
        }
    };

    // 切换自动判题
    const handleToggleAutoGrading = async (id, currentAutoGrading) => {
        try {
            const res = await toggleAutoGrading({
                id,
                auto_grading: !currentAutoGrading
            });
            if (res.code === 0) {
                message.success(currentAutoGrading ? '自动判题已关闭' : '自动判题已开启');
                loadConfigs();
            } else {
                message.error(res.msg || '操作失败');
            }
        } catch (error) {
            message.error('操作失败');
        }
    };

    // 手动触发批量判题（预览模式）
    const handleBatchGrade = async () => {
        if (pendingCount === 0) {
            message.info('当前没有待判题的开放题');
            return;
        }

        setPreviewLoading(true);
        try {
            const res = await batchGradePreview({ max_count: batchCount });
            if (res.code === 0) {
                const results = res.data.results || [];
                const count = res.data.count || results.length || 0;
                setPreviewResults(results);

                if (count === 0) {
                    message.info('当前没有待判题的开放题');
                } else {
                    setPreviewModalOpen(true);
                    message.success(`已生成 ${count} 条判题结果，请检查后提交`);
                }
            } else {
                message.error(res.msg || '批量判题失败');
            }
        } catch (error) {
            message.error('批量判题失败');
        } finally {
            setPreviewLoading(false);
        }
    };

    // 提交批量判题结果
    const handleSubmitBatchGrade = async () => {
        if (previewResults.length === 0) {
            message.warning('没有需要提交的判题结果');
            return;
        }

        setSubmitting(true);
        try {
            const results = previewResults.map(item => ({
                record_id: item.record_id,
                answer_index: item.answer_index,
                score: item.score,
                annotation: item.annotation
            }));

            const res = await submitBatchGrade({ results });
            if (res.code === 0) {
                message.success(res.data.message || '提交成功');
                setPreviewModalOpen(false);
                setPreviewResults([]);
                loadPendingCount();
            } else {
                message.error(res.msg || '提交失败');
            }
        } catch (error) {
            message.error('提交失败');
        } finally {
            setSubmitting(false);
        }
    };

    // 打开AI判题历史
    const handleOpenHistory = () => {
        setHistoryModalOpen(true);
        loadHistory(1, 10);
    };

    // 加载AI判题历史
    const loadHistory = async (page = 1, pageSize = 10) => {
        setHistoryLoading(true);
        try {
            const res = await getAiGradeHistory({ page, pageSize });
            if (res.code === 0) {
                setHistoryList(res.data.list || []);
                setHistoryPagination({
                    current: res.data.page,
                    pageSize: res.data.pageSize,
                    total: res.data.total
                });
            }
        } catch (error) {
            message.error('获取历史记录失败');
        } finally {
            setHistoryLoading(false);
        }
    };

    // 取消编辑
    const handleCancelForm = () => {
        setShowForm(false);
        setEditMode(false);
        setEditId(null);
        form.resetFields();
    };

    // 获取当前启用且开启自动判题的配置
    const enabledConfig = configs.find(c => c.enabled);

    // 预览结果表格列定义
    const previewColumns = [
        {
            title: '题目',
            dataIndex: 'question_title',
            width: 200,
            ellipsis: true,
        },
        {
            title: '挑战',
            dataIndex: 'challenge_title',
            width: 150,
            ellipsis: true,
        },
        {
            title: '用户',
            dataIndex: 'address',
            width: 120,
            render: (addr) => `${addr.slice(0, 6)}...${addr.slice(-4)}`
        },
        {
            title: '分数',
            dataIndex: 'score',
            width: 150,
            render: (score, record) => (
                <InputNumber
                    min={0}
                    max={record.question_score}
                    value={score}
                    onChange={(val) => {
                        const newResults = [...previewResults];
                        const index = newResults.findIndex(r =>
                            r.record_id === record.record_id && r.answer_index === record.answer_index
                        );
                        if (index !== -1) {
                            newResults[index].score = val;
                            setPreviewResults(newResults);
                        }
                    }}
                    addonAfter={`/ ${record.question_score}`}
                />
            )
        },
        {
            title: '批注',
            dataIndex: 'annotation',
            render: (annotation, record) => (
                <Input.TextArea
                    value={annotation}
                    rows={2}
                    onChange={(e) => {
                        const newResults = [...previewResults];
                        const index = newResults.findIndex(r =>
                            r.record_id === record.record_id && r.answer_index === record.answer_index
                        );
                        if (index !== -1) {
                            newResults[index].annotation = e.target.value;
                            setPreviewResults(newResults);
                        }
                    }}
                    placeholder="不通过的理由（通过则留空）"
                />
            )
        },
        {
            title: '操作',
            width: 100,
            render: (_, record) => (
                <Button
                    type="link"
                    onClick={() => showDebugModal(record)}
                >
                    查看详情
                </Button>
            )
        }
    ];

    // 历史记录表格列定义
    const historyColumns = [
        {
            title: '判题时间',
            dataIndex: 'review_time',
            width: 180,
            render: (time) => time || '-'
        },
        {
            title: '判题方式',
            width: 100,
            render: () => (
                <Badge color="blue" text="AI判题" />
            )
        },
        {
            title: '题目',
            dataIndex: 'question_title',
            ellipsis: true,
        },
        {
            title: '挑战',
            dataIndex: 'challenge_title',
            ellipsis: true,
        },
        {
            title: '用户',
            dataIndex: 'address',
            width: 120,
            render: (addr) => `${addr.slice(0, 6)}...${addr.slice(-4)}`
        },
        {
            title: '分数',
            dataIndex: 'score',
            width: 100,
            render: (score, record) => `${score} / ${record.question_score}`
        },
        {
            title: '是否通过',
            dataIndex: 'pass',
            width: 100,
            render: (pass) => (
                <Badge
                    status={pass ? 'success' : 'error'}
                    text={pass ? '通过' : '未通过'}
                />
            )
        },
        {
            title: '批注',
            dataIndex: 'annotation',
            width: 150,
            ellipsis: true,
            render: (annotation) => annotation || '-'
        },
        {
            title: '操作',
            width: 100,
            render: (_, record) => (
                <Button
                    type="link"
                    onClick={() => showHistoryDetail(record)}
                >
                    查看详情
                </Button>
            )
        }
    ];

    // 显示调试详情
    const [debugModalOpen, setDebugModalOpen] = useState(false);
    const [debugInfo, setDebugInfo] = useState({});

    // 历史详情
    const [historyDetailModalOpen, setHistoryDetailModalOpen] = useState(false);
    const [historyDetail, setHistoryDetail] = useState({});

    const showDebugModal = (record) => {
        setDebugInfo({
            systemPrompt: record.system_prompt,
            userPrompt: record.user_prompt,
            rawResult: record.raw_result,
            userAnswer: record.user_answer,
            attachmentUrls: record.attachment_urls
        });
        setDebugModalOpen(true);
    };

    const showHistoryDetail = (record) => {
        setHistoryDetail(record);
        setHistoryDetailModalOpen(true);
    };

    return (
        <Modal
            title={
                <div style={{ display: 'flex', alignItems: 'center', gap: '30px' }}>
                    <span>AI判题配置</span>
                    <QuestionCircleOutlined
                        style={{ fontSize: 18, cursor: 'pointer', color: '#1890ff' }}
                        onClick={() => setGuideModalOpen(true)}
                        title="查看使用说明"
                    />
                </div>
            }
            open={open}
            onCancel={onCancel}
            width={900}
            footer={[
                <Button key="close" onClick={onCancel}>
                    关闭
                </Button>
            ]}
        >
            {/* 顶部统计信息 */}
            <div style={{
                marginBottom: 16,
                padding: 16,
                background: '#f5f5f5',
                borderRadius: 4,
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center'
            }}>
                <div style={{ display: 'flex', gap: 24 }}>
                    <Statistic
                        title="待判题数量"
                        value={pendingCount}
                        suffix="题"
                        valueStyle={{ color: pendingCount > 0 ? '#cf1322' : '#3f8600' }}
                    />
                    {enabledConfig && (
                        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                            <span>自动判题：</span>
                            <Switch
                                checked={enabledConfig.auto_grading}
                                onChange={() => handleToggleAutoGrading(enabledConfig.id, enabledConfig.auto_grading)}
                                checkedChildren="开"
                                unCheckedChildren="关"
                            />
                            {enabledConfig.auto_grading && (
                                <Badge status="processing" text="运行中" />
                            )}
                        </div>
                    )}
                </div>
                <Space>
                    <Button
                        icon={<HistoryOutlined />}
                        onClick={handleOpenHistory}
                    >
                        AI判题历史
                    </Button>
                    <Button
                        type="primary"
                        icon={<SyncOutlined spin={previewLoading} />}
                        onClick={handleBatchGrade}
                        loading={previewLoading}
                        disabled={pendingCount === 0 || !enabledConfig}
                    >
                        手动批量判题
                    </Button>
                    <InputNumber
                        min={1}
                        max={10}
                        value={batchCount}
                        onChange={setBatchCount}
                        style={{ width: 80 }}
                        placeholder="数量"
                        disabled={pendingCount === 0 || !enabledConfig}
                    />
                    <span style={{ color: '#999', fontSize: 12 }}>题</span>
                </Space>
            </div>

            <div style={{ marginBottom: 16 }}>
                <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    onClick={handleAdd}
                    disabled={showForm}
                >
                    新增配置
                </Button>
            </div>

            {!showForm && configs.length === 0 && (
                <div style={{ textAlign: 'center', padding: '40px 0', color: '#999' }}>
                    暂无配置，请点击"新增配置"按钮添加
                </div>
            )}

            {showForm && (
                <Form
                    form={form}
                    layout="vertical"
                    style={{
                        marginBottom: 24,
                        padding: 16,
                        border: '1px solid #d9d9d9',
                        borderRadius: 4,
                        background: '#fafafa'
                    }}
                >
                    <Form.Item
                        label="标题"
                        name="title"
                        rules={[{ required: true, message: '请输入配置标题' }]}
                    >
                        <Input placeholder="请输入配置标题" />
                    </Form.Item>

                    <Form.Item
                        label={
                            <div>
                                配置内容
                                <div style={{ color: '#999', fontWeight: 'normal', fontSize: 12, marginTop: 4 }}>
                                    请配置AI判题的ApiKey和对应的BaseUrl，格式为JSON
                                </div>
                            </div>
                        }
                        name="config"
                        rules={[
                            { required: true, message: '请输入配置内容' }
                        ]}
                    >
                        <ReactQuill
                            theme="snow"
                            modules={modules}
                            formats={formats}
                            placeholder={'请输入JSON格式的配置，例如：{"apiKey": "your-api-key", "baseUrl": "https://api.example.com"}'}
                            style={{ background: 'white' }}
                        />
                    </Form.Item>

                    <Form.Item style={{ marginBottom: 0 }}>
                        <Space>
                            <Button type="primary" onClick={handleSave} loading={loading}>
                                保存
                            </Button>
                            <Button onClick={handleCancelForm}>
                                取消
                            </Button>
                        </Space>
                    </Form.Item>
                </Form>
            )}

            {!showForm && configs.length > 0 && (
                <List
                    bordered
                    dataSource={configs}
                    renderItem={(item) => (
                        <List.Item
                            actions={[
                                <div key="enabled" style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                                    <span style={{ color: item.enabled ? '#52c41a' : '#999', fontSize: 12 }}>
                                        {item.enabled ? '已启用' : '未启用'}
                                    </span>
                                    <Switch
                                        checked={item.enabled}
                                        onChange={() => handleToggle(item.id, item.enabled)}
                                        checkedChildren="启用"
                                        unCheckedChildren="禁用"
                                    />
                                </div>,
                                <Button
                                    key="edit"
                                    type="link"
                                    icon={<EditOutlined />}
                                    onClick={() => handleEdit(item)}
                                >
                                    编辑
                                </Button>,
                                <Popconfirm
                                    key="delete"
                                    title="确认删除"
                                    description="确定要删除这个配置吗？"
                                    onConfirm={() => handleDelete(item.id)}
                                    okText="确定"
                                    cancelText="取消"
                                >
                                    <Button
                                        type="link"
                                        danger
                                        icon={<DeleteOutlined />}
                                    >
                                        删除
                                    </Button>
                                </Popconfirm>
                            ]}
                        >
                            <List.Item.Meta
                                title={
                                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                                        {item.title}
                                        {item.enabled && item.auto_grading && (
                                            <Badge status="processing" text="自动判题中" />
                                        )}
                                    </div>
                                }
                                description={
                                    <div
                                        dangerouslySetInnerHTML={{ __html: item.config }}
                                        style={{
                                            maxHeight: 100,
                                            overflow: 'auto',
                                            padding: 8,
                                            background: '#f5f5f5',
                                            borderRadius: 4,
                                            fontSize: 12
                                        }}
                                    />
                                }
                            />
                        </List.Item>
                    )}
                />
            )}

            {/* 批量判题预览弹窗 */}
            <Modal
                title="批量判题结果预览"
                open={previewModalOpen}
                onCancel={() => setPreviewModalOpen(false)}
                width={1200}
                footer={[
                    <Button key="cancel" onClick={() => setPreviewModalOpen(false)}>
                        取消
                    </Button>,
                    <Button
                        key="submit"
                        type="primary"
                        onClick={handleSubmitBatchGrade}
                        loading={submitting}
                    >
                        确认提交 ({previewResults.length} 条)
                    </Button>
                ]}
            >
                <div style={{ marginBottom: 16, color: '#666' }}>
                    请检查AI判题结果，可以修改分数和批注后再提交
                </div>
                <Table
                    columns={previewColumns}
                    dataSource={previewResults}
                    rowKey={(record) => `${record.record_id}_${record.answer_index}`}
                    pagination={false}
                    scroll={{ y: 500 }}
                />
            </Modal>

            {/* AI判题历史弹窗 */}
            <Modal
                title="AI判题历史"
                open={historyModalOpen}
                onCancel={() => setHistoryModalOpen(false)}
                width={1200}
                footer={[
                    <Button key="close" onClick={() => setHistoryModalOpen(false)}>
                        关闭
                    </Button>
                ]}
            >
                <Table
                    columns={historyColumns}
                    dataSource={historyList}
                    rowKey={(record) => `${record.record_id}_${record.answer_index}_${record.review_time}`}
                    loading={historyLoading}
                    pagination={{
                        ...historyPagination,
                        onChange: (page, pageSize) => loadHistory(page, pageSize)
                    }}
                />
            </Modal>

            {/* 调试详情弹窗 */}
            <Modal
                title="AI判题详情"
                open={debugModalOpen}
                onCancel={() => setDebugModalOpen(false)}
                width={900}
                footer={[
                    <Button key="close" onClick={() => setDebugModalOpen(false)}>
                        关闭
                    </Button>
                ]}
            >
                <div style={{ maxHeight: '70vh', overflow: 'auto' }}>
                    <div style={{ marginBottom: 20 }}>
                        <h3>🤖 System Prompt（系统提示词）</h3>
                        <pre style={{ background: '#f5f5f5', padding: 12, borderRadius: 4, whiteSpace: 'pre-wrap' }}>
                            {debugInfo.systemPrompt}
                        </pre>
                    </div>
                    <div style={{ marginBottom: 20 }}>
                        <h3>💬 User Prompt（用户提示词）</h3>
                        <pre style={{ background: '#f5f5f5', padding: 12, borderRadius: 4, whiteSpace: 'pre-wrap' }}>
                            {debugInfo.userPrompt}
                        </pre>
                    </div>
                    <div style={{ marginBottom: 20 }}>
                        <h3>📝 用户答案</h3>
                        <pre style={{ background: '#f5f5f5', padding: 12, borderRadius: 4, whiteSpace: 'pre-wrap' }}>
                            {debugInfo.userAnswer}
                        </pre>
                    </div>
                    {debugInfo.attachmentUrls && debugInfo.attachmentUrls.length > 0 && (
                        <div style={{ marginBottom: 20 }}>
                            <h3>📎 附件</h3>
                            <div style={{ background: '#f5f5f5', padding: 12, borderRadius: 4 }}>
                                {debugInfo.attachmentUrls.map((url, index) => (
                                    <div key={index}>{index + 1}. {url}</div>
                                ))}
                            </div>
                        </div>
                    )}
                    <div>
                        <h3>📥 AI返回的原始结果</h3>
                        <pre style={{ background: '#f5f5f5', padding: 12, borderRadius: 4, whiteSpace: 'pre-wrap' }}>
                            {debugInfo.rawResult}
                        </pre>
                    </div>
                </div>
            </Modal>

            {/* 历史详情弹窗 */}
            <Modal
                title="AI判题详情"
                open={historyDetailModalOpen}
                onCancel={() => setHistoryDetailModalOpen(false)}
                width={900}
                footer={[
                    <Button key="close" onClick={() => setHistoryDetailModalOpen(false)}>
                        关闭
                    </Button>
                ]}
            >
                <div style={{ lineHeight: '1.8' }}>
                    <div style={{ marginBottom: 20 }}>
                        <h3>📝 题目信息</h3>
                        <div style={{ background: '#f5f5f5', padding: 12, borderRadius: 4 }}>
                            <p><strong>挑战：</strong>{historyDetail.challenge_title}</p>
                            <p><strong>题目：</strong>{historyDetail.question_title}</p>
                            <p><strong>用户：</strong>{historyDetail.address}</p>
                            <p><strong>判题时间：</strong>{historyDetail.review_time}</p>
                        </div>
                    </div>
                    <div style={{ marginBottom: 20 }}>
                        <h3>✍️ 用户答案</h3>
                        <pre style={{ background: '#f5f5f5', padding: 12, borderRadius: 4, whiteSpace: 'pre-wrap', maxHeight: 300, overflow: 'auto' }}>
                            {historyDetail.user_answer}
                        </pre>
                    </div>
                    <div style={{ marginBottom: 20 }}>
                        <h3>📊 评分结果</h3>
                        <div style={{ background: '#f5f5f5', padding: 12, borderRadius: 4 }}>
                            <p><strong>得分：</strong>{historyDetail.score} / {historyDetail.question_score}</p>
                            <p>
                                <strong>是否通过：</strong>
                                <Badge
                                    status={historyDetail.pass ? 'success' : 'error'}
                                    text={historyDetail.pass ? '通过' : '未通过'}
                                    style={{ marginLeft: 8 }}
                                />
                            </p>
                        </div>
                    </div>
                    {historyDetail.annotation && (
                        <div style={{ marginBottom: 20 }}>
                            <h3>💬 批注</h3>
                            <pre style={{ background: '#fff3cd', padding: 12, borderRadius: 4, whiteSpace: 'pre-wrap', maxHeight: 300, overflow: 'auto' }}>
                                {historyDetail.annotation}
                            </pre>
                        </div>
                    )}
                </div>
            </Modal>

            {/* 使用说明弹窗 */}
            <Modal
                title="AI判题使用说明"
                open={guideModalOpen}
                onCancel={() => setGuideModalOpen(false)}
                width={800}
                footer={[
                    <Button key="close" type="primary" onClick={() => setGuideModalOpen(false)}>
                        知道了
                    </Button>
                ]}
            >
                <div style={{ lineHeight: '1.8' }}>
                    <h3>📋 功能说明</h3>
                    <p>AI判题功能可以自动批改开放题答案，大大提升评分效率。</p>

                    <h3>🔧 配置步骤</h3>
                    <ol>
                        <li><strong>添加AI配置</strong>：点击"新增配置"按钮，输入配置名称和AI API配置信息（JSON格式）</li>
                        <li><strong>启用配置</strong>：在配置列表中找到你创建的配置，点击"启用"开关</li>
                        <li><strong>选择判题模式</strong>：
                            <ul>
                                <li><strong>自动判题</strong>：开启后，系统会每分钟自动处理待判题的开放题（最多10题/分钟）</li>
                                <li><strong>手动判题</strong>：点击"手动批量判题"按钮，可以预览AI判题结果并修改后提交</li>
                            </ul>
                        </li>
                    </ol>

                    <h3>⚙️ 运行逻辑</h3>
                    <ul>
                        <li><strong>待判题统计</strong>：系统会统计所有已上架题目中未评分的开放题答案数量</li>
                        <li><strong>判题过程</strong>：
                            <ol>
                                <li>提取用户答案和附件信息</li>
                                <li>构建系统提示词（包含题目要求、总分、及格分）</li>
                                <li>调用AI API进行判题</li>
                                <li>解析AI返回结果，提取分数和批注</li>
                                <li>更新到数据库（手动模式需确认后提交）</li>
                            </ol>
                        </li>
                        <li><strong>自动判题</strong>：后台定时任务每分钟运行一次，自动处理待判题列表</li>
                        <li><strong>手动判题</strong>：一次最多处理10题，可以在提交前修改AI的评分和批注</li>
                    </ul>

                    <h3>📊 查看历史</h3>
                    <p>点击"AI判题历史"按钮可以查看所有已完成的AI判题记录，包括判题时间、分数、是否通过等信息。</p>

                    <h3>💡 使用建议</h3>
                    <ul>
                        <li>首次使用建议先用<strong>手动判题</strong>，检查AI判题质量</li>
                        <li>确认AI判题准确度后，可以开启<strong>自动判题</strong>提高效率</li>
                        <li>定期查看判题历史，确保评分公平公正</li>
                        <li>对于特殊题目，建议人工复核AI判题结果</li>
                    </ul>

                    <h3>⚠️ 注意事项</h3>
                    <ul>
                        <li>确保AI API配置正确，包含有效的apiKey和baseUrl</li>
                        <li>同一时间只能启用一个AI配置</li>
                        <li>AI判题结果仅供参考，重要评分建议人工复核</li>
                        <li>手动判题预览后必须点击"确认提交"才会保存到数据库</li>
                    </ul>
                </div>
            </Modal>
        </Modal>
    );
}
