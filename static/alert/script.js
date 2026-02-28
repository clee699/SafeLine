// 全局变量
let currentConfigId = null;

// DOM元素
const modal = document.getElementById('config-modal');
const modalTitle = document.getElementById('modal-title');
const configForm = document.getElementById('config-form');
const closeBtn = document.querySelector('.close');
const addConfigBtn = document.getElementById('add-config');
const saveGlobalBtn = document.getElementById('save-global');
const configTypeSelect = document.getElementById('config-type');
const emailConfigSection = document.getElementById('email-config');
const wechatConfigSection = document.getElementById('wechat-config');
const configList = document.getElementById('config-list');

// API基础URL
const API_BASE_URL = '/api';

// 初始化页面
document.addEventListener('DOMContentLoaded', function() {
    loadGlobalSettings();
    loadAlertConfigs();
    
    // 绑定事件监听器
    bindEventListeners();
});

// 绑定事件监听器
function bindEventListeners() {
    // 模态框相关
    closeBtn.addEventListener('click', closeModal);
    window.addEventListener('click', function(event) {
        if (event.target == modal) {
            closeModal();
        }
    });
    
    // 表单提交
    configForm.addEventListener('submit', handleConfigSubmit);
    document.getElementById('cancel-form').addEventListener('click', closeModal);
    
    // 添加配置按钮
    addConfigBtn.addEventListener('click', openAddModal);
    
    // 保存全局设置
    saveGlobalBtn.addEventListener('click', saveGlobalSettings);
    
    // 配置类型切换
    configTypeSelect.addEventListener('change', handleConfigTypeChange);
}

// 加载全局设置
function loadGlobalSettings() {
    fetch(`${API_BASE_URL}/Option`) // 假设获取全局设置的API
        .then(response => response.json())
        .then(data => {
            // 解析全局设置
            const settings = data || {};
            
            // 设置全局开关
            const alertEnabled = settings[constants.AlertEnabled] === 'true';
            document.getElementById('alert-enabled').checked = alertEnabled;
            
            // 设置检查间隔
            const checkInterval = settings[constants.AlertCheckInterval] || '10';
            document.getElementById('check-interval').value = checkInterval;
            
            // 设置默认告警级别
            // 假设从配置文件读取
            document.getElementById('default-level').value = '2';
        })
        .catch(error => {
            console.error('加载全局设置失败:', error);
            alert('加载全局设置失败，请刷新页面重试');
        });
}

// 保存全局设置
function saveGlobalSettings() {
    const enabled = document.getElementById('alert-enabled').checked;
    const interval = document.getElementById('check-interval').value;
    const defaultLevel = document.getElementById('default-level').value;
    
    // 构建请求数据
    const settings = {
        [constants.AlertEnabled]: enabled ? 'true' : 'false',
        [constants.AlertCheckInterval]: interval
    };
    
    // 发送请求保存全局设置
    fetch(`${API_BASE_URL}/Option`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(settings)
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            alert('全局设置保存成功');
        } else {
            alert('全局设置保存失败: ' + (data.message || '未知错误'));
        }
    })
    .catch(error => {
        console.error('保存全局设置失败:', error);
        alert('保存全局设置失败，请刷新页面重试');
    });
}

// 加载告警配置列表
function loadAlertConfigs() {
    fetch(`${API_BASE_URL}/AlertConfig`)
        .then(response => response.json())
        .then(configs => {
            renderConfigList(configs);
        })
        .catch(error => {
            console.error('加载告警配置失败:', error);
            alert('加载告警配置失败，请刷新页面重试');
        });
}

// 渲染配置列表
function renderConfigList(configs) {
    configList.innerHTML = '';
    
    if (!configs || configs.length === 0) {
        configList.innerHTML = '<tr><td colspan="4" style="text-align: center; color: #6c757d;">暂无告警配置</td></tr>';
        return;
    }
    
    configs.forEach(config => {
        const row = document.createElement('tr');
        
        // 告警类型显示文本
        const typeText = config.Type === 1 ? '邮件告警' : '微信告警';
        
        // 状态显示
        const statusText = config.Enabled ? '启用' : '禁用';
        const statusClass = config.Enabled ? 'status enabled' : 'status disabled';
        
        row.innerHTML = `
            <td>${config.Name}</td>
            <td>${typeText}</td>
            <td><span class="${statusClass}">${statusText}</span></td>
            <td>
                <button class="btn btn-primary btn-sm edit-btn" data-id="${config.ID}">编辑</button>
                <button class="btn btn-danger btn-sm delete-btn" data-id="${config.ID}">删除</button>
            </td>
        `;
        
        // 添加编辑和删除事件监听器
        row.querySelector('.edit-btn').addEventListener('click', function() {
            openEditModal(config);
        });
        
        row.querySelector('.delete-btn').addEventListener('click', function() {
            deleteConfig(config.ID);
        });
        
        configList.appendChild(row);
    });
}

// 打开添加模态框
function openAddModal() {
    currentConfigId = null;
    modalTitle.textContent = '添加告警配置';
    configForm.reset();
    
    // 切换到邮件配置
    configTypeSelect.value = '1';
    emailConfigSection.style.display = 'block';
    wechatConfigSection.style.display = 'none';
    
    // 显示模态框
    modal.style.display = 'block';
}

// 打开编辑模态框
function openEditModal(config) {
    currentConfigId = config.ID;
    modalTitle.textContent = '编辑告警配置';
    
    // 填充表单数据
    document.getElementById('config-id').value = config.ID;
    document.getElementById('config-name').value = config.Name;
    document.getElementById('config-type').value = config.Type;
    document.getElementById('config-enabled').checked = config.Enabled;
    document.getElementById('alert-levels').value = config.AlertLevels || '';
    document.getElementById('attack-types').value = config.AttackTypes || '';
    
    // 解析配置JSON
    let configJson;
    try {
        configJson = JSON.parse(config.Config);
    } catch (e) {
        configJson = {};
    }
    
    // 根据类型显示对应配置
    if (config.Type === 1) {
        // 邮件配置
        emailConfigSection.style.display = 'block';
        wechatConfigSection.style.display = 'none';
        
        document.getElementById('smtp-server').value = configJson.smtp_server || '';
        document.getElementById('smtp-port').value = configJson.smtp_port || '25';
        document.getElementById('smtp-username').value = configJson.username || '';
        document.getElementById('smtp-password').value = configJson.password || '';
        document.getElementById('smtp-from').value = configJson.from || '';
        document.getElementById('smtp-to').value = configJson.to ? configJson.to.join(',') : '';
    } else {
        // 微信配置
        emailConfigSection.style.display = 'none';
        wechatConfigSection.style.display = 'block';
        
        document.getElementById('wechat-webhook').value = configJson.webhook_url || '';
    }
    
    // 显示模态框
    modal.style.display = 'block';
}

// 关闭模态框
function closeModal() {
    modal.style.display = 'none';
    configForm.reset();
    currentConfigId = null;
}

// 处理配置类型切换
function handleConfigTypeChange() {
    const type = configTypeSelect.value;
    
    if (type === '1') {
        // 邮件配置
        emailConfigSection.style.display = 'block';
        wechatConfigSection.style.display = 'none';
    } else {
        // 微信配置
        emailConfigSection.style.display = 'none';
        wechatConfigSection.style.display = 'block';
    }
}

// 处理配置提交
function handleConfigSubmit(event) {
    event.preventDefault();
    
    // 收集表单数据
    const id = document.getElementById('config-id').value;
    const name = document.getElementById('config-name').value;
    const type = parseInt(document.getElementById('config-type').value);
    const enabled = document.getElementById('config-enabled').checked;
    const alertLevels = document.getElementById('alert-levels').value;
    const attackTypes = document.getElementById('attack-types').value;
    
    // 构建配置JSON
    let configJson = {};
    
    if (type === 1) {
        // 邮件配置
        const to = document.getElementById('smtp-to').value.split(',').map(email => email.trim());
        configJson = {
            smtp_server: document.getElementById('smtp-server').value,
            smtp_port: parseInt(document.getElementById('smtp-port').value),
            username: document.getElementById('smtp-username').value,
            password: document.getElementById('smtp-password').value,
            from: document.getElementById('smtp-from').value,
            to: to
        };
    } else {
        // 微信配置
        configJson = {
            webhook_url: document.getElementById('wechat-webhook').value
        };
    }
    
    // 构建请求数据
    const configData = {
        ID: id ? parseInt(id) : 0,
        Name: name,
        Type: type,
        Enabled: enabled,
        Config: JSON.stringify(configJson),
        AlertLevels: alertLevels,
        AttackTypes: attackTypes
    };
    
    // 发送请求
    const url = currentConfigId ? `${API_BASE_URL}/AlertConfig` : `${API_BASE_URL}/AlertConfig`;
    const method = currentConfigId ? 'PUT' : 'POST';
    
    fetch(url, {
        method: method,
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(configData)
    })
    .then(response => response.json())
    .then(data => {
        if (data.success || data) {
            alert(currentConfigId ? '编辑配置成功' : '添加配置成功');
            closeModal();
            loadAlertConfigs();
        } else {
            alert('操作失败: ' + (data.message || '未知错误'));
        }
    })
    .catch(error => {
        console.error('操作失败:', error);
        alert('操作失败，请刷新页面重试');
    });
}

// 删除配置
function deleteConfig(id) {
    if (!confirm('确定要删除这个告警配置吗？')) {
        return;
    }
    
    fetch(`${API_BASE_URL}/AlertConfig?id=${id}`, {
        method: 'DELETE'
    })
    .then(response => response.json())
    .then(data => {
        if (data.success || data) {
            alert('删除配置成功');
            loadAlertConfigs();
        } else {
            alert('删除失败: ' + (data.message || '未知错误'));
        }
    })
    .catch(error => {
        console.error('删除配置失败:', error);
        alert('删除配置失败，请刷新页面重试');
    });
}
