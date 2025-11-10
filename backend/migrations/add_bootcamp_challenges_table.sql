-- 创建训练营挑战配置表
CREATE TABLE IF NOT EXISTS bootcamp_challenges (
    id SERIAL PRIMARY KEY,
    quest_id INTEGER NOT NULL,
    enabled BOOLEAN DEFAULT false,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(quest_id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_bootcamp_challenges_enabled ON bootcamp_challenges(enabled);

-- 添加注释
COMMENT ON TABLE bootcamp_challenges IS '训练营挑战配置表';
COMMENT ON COLUMN bootcamp_challenges.quest_id IS '挑战ID';
COMMENT ON COLUMN bootcamp_challenges.enabled IS '是否启用';
COMMENT ON COLUMN bootcamp_challenges.display_order IS '展示排序（数字越小越靠前）';
