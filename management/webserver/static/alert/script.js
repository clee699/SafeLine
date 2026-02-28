// 告警配置管理JavaScript

// 全局变量
let editingConfigId = null;

// 初始化页面
document.addEventListener('DOMContentLoaded', function() {
    // 绑定事件监听器
    bindEventListeners();
    
    // 加载全局告警开关状态
    loadGlobalAlertSwitch();
    
    // 加载告警配置列表
    loadAlertConfigs();
});

// 绑定事件监听器
function bindEventListeners() {
    // 添加配置按钮
    document.getElementById('add-config-btn').addEventListener('click', openAddModal);
    
    // 关闭模态框
    document.querySelector('.close').addEventListener('click', closeModal);
    document.getElementById('cancel-btn').addEventListener('click', closeModal);
    
    // 点击模态框外部关闭
    window.addEventListener('click', function(event) {
        const modal = document.getElementById('config-modal');
        if (event.target === modal) {
            closeModal();
        }
    });
    
    // 表单提交
    document.getElementById('config-form').addEventListener('submit', handleFormSubmit);
    
    // 全局告警开关
    document.getElementById('global-alert-switch').addEventListener('change', updateGlobalAlertSwitch);
}

// 加载全局告警开关状态
function loadGlobalAlertSwitch() {
    fetch('/api/options?key=alert_enabled')
        .then(response => response.json())
        .then(data => {
            if (data.code === 0 && data.data) {
                document.getElementById('global-alert-switch').checked = data.data.value === 'true';
            }
        })
        .catch(error => {
            console.error('加载全局告警开关失败:', error);
        });
}

// 更新全局告警开关状态
function updateGlobalAlertSwitch() {
    const enabled = document.getElementById('global-alert-switch').checked;
    
    fetch('/api/options', {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            key: 'alert_enabled',
            value: enabled.toString()
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.code !== 0) {
            // 恢复原状态
            document.getElementById('global-alert-switch').checked = !enabled;
            alert('更新全局告警开关失败: ' + data.msg);
        }
    })
    .catch(error => {
        console.error('更新全局告警开关失败:', error);
        // 恢复原状态
        document.getElementById('global-alert-switch').checked = !enabled;
        alert('更新全局告警开关失败');
    });
}

// 加载告警配置列表
function loadAlertConfigs() {
    fetch('/api/alert_config')
        .then(response => response.json())
        .then(data => {
            if (data.code === 0 && Array.isArray(data.data)) {
                renderConfigList(data.data);
            } else {
                console.error('加载告警配置失败:', data.msg);
            }
        })
        .catch(error => {
            console.error('加载告警配置失败:', error);
        });
}

// 渲染配置列表
function renderConfigList(configs) {
    const tbody = document.getElementById('config-list');
    tbody.innerHTML = '';
    
    configs.forEach(config => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td>${config.name}</td>
            <td>${config.type === 'email' ? '邮件告警' : '微信告警'}</td>
            <td><span class="status-badge ${config.enabled ? 'status-enabled' : 'status-disabled'}">${config.enabled ? '已启用' : '已禁用'}</span></td>
            <td>
                <button class="btn btn-secondary" onclick="editConfig(${config.id})">编辑</button>
                <button class="btn btn-danger" onclick="deleteConfig(${config.id})">删除</button>
            </td>
        `;
        tbody.appendChild(tr);
    });
}

// 打开添加配置模态框
function openAddModal() {
    editingConfigId = null;
    document.getElementById('modal-title').textContent = '添加告警配置';
    document.getElementById('config-form').reset();
    document.getElementById('config-enabled').checked = true;
    document.getElementById('config-id').value = '';
    document.getElementById('config-modal').style.display = 'block';
}

// 打开编辑配置模态框
function editConfig(id) {
    fetch(`/api/alert_config?id=${id}`)
        .then(response => response.json())
        .then(data => {
            if (data.code === 0 && data.data) {
                const config = data.data;
                editingConfigId = config.id;
                document.getElementById('modal-title').textContent = '编辑告警配置';
                document.getElementById('config-id').value = config.id;
                document.getElementById('config-name').value = config.name;
                document.getElementById('config-type').value = config.type;
                document.getElementById('config-enabled').checked = config.enabled;
                document.getElementById('config-settings').value = JSON.stringify(config.settings, null, 2);
                document.getElementById('config-modal').style.display = 'block';
            }
        })
        .catch(error => {
            console.error('加载配置详情失败:', error);
        });
}

// 关闭模态框
function closeModal() {
    document.getElementById('config-modal').style.display = 'none';
    editingConfigId = null;
}

// 处理表单提交
function handleFormSubmit(event) {
    event.preventDefault();
    
    // 验证JSON格式
    let settings;
    try {
        settings = JSON.parse(document.getElementById('config-settings').value);
    } catch (e) {
        alert('配置必须是有效的JSON格式');
        return;
    }
    
    const configData = {
        name: document.getElementById('config-name').value,
        type: document.getElementById('config-type').value,
        enabled: document.getElementById('config-enabled').checked,
        settings: settings
    };
    
    const url = editingConfigId ? `/api/alert_config` : `/api/alert_config`;
    const method = editingConfigId ? 'PUT' : 'POST';
    
    if (editingConfigId) {
        configData.id = editingConfigId;
    }
    
    fetch(url, {
        method: method,
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(configData)
    })
    .then(response => response.json())
    .then(data => {
        if (data.code === 0) {
            closeModal();
            loadAlertConfigs();
        } else {
            alert('保存配置失败: ' + data.msg);
        }
    })
    .catch(error => {
        console.error('保存配置失败:', error);
        alert('保存配置失败');
    });
}

// 删除配置
function deleteConfig(id) {
    if (confirm('确定要删除这个告警配置吗？')) {
        fetch(`/api/alert_config?id=${id}`, {
            method: 'DELETE'
        })
        .then(response => response.json())
        .then(data => {
            if (data.code === 0) {
                loadAlertConfigs();
            } else {
                alert('删除配置失败: ' + data.msg);
            }
        })
        .catch(error => {
            console.error('删除配置失败:', error);
            alert('删除配置失败');
        });
    }
}