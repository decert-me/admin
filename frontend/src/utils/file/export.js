export function exportJsonFile(jsonData, fileName) {
    // 创建包含JSON数据的Blob对象
    const blob = new Blob([JSON.stringify(jsonData)], { type: 'application/json' });
  
    // 创建一个<a>标签
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = fileName + '.json';
  
    // 模拟用户点击链接来触发文件下载
    link.click();
  
    // 清理临时元素
    // document.body.removeChild(link);
    // URL.revokeObjectURL(link.href);
}  