import './style.css';
import './app.css';
import * as XLSX from 'xlsx';
import { Login, GetLogs, RunRPA, GetDownloadPath } from '../wailsjs/go/main/App';
import { UpdateOmo, SayHi } from '../wailsjs/go/main/App';
window.addEventListener("DOMContentLoaded", () => {
    console.log("DOM fully loaded and parsed");

    const usernameElement = document.getElementById("username");
    const passwordElement = document.getElementById("password");
    const statusElement = document.getElementById("status");

    // 添加记住密码复选框
    const rememberMeHTML = `
        <div class="remember-me">
            <input type="checkbox" id="rememberMe">
            <label for="rememberMe">记住账号密码</label>
        </div>
    `;
    passwordElement.insertAdjacentHTML('afterend', rememberMeHTML);
    const rememberMeCheckbox = document.getElementById("rememberMe");

    // 从 localStorage 获取保存的账号密码
    const savedUsername = localStorage.getItem("username");
    const savedPassword = localStorage.getItem("password");
    const rememberMe = localStorage.getItem("rememberMe") === "true";

    if (rememberMe && savedUsername && savedPassword) {
        usernameElement.value = savedUsername;
        passwordElement.value = savedPassword;
        rememberMeCheckbox.checked = true;
    } else {
        usernameElement.value = "";
        passwordElement.value = "";
        rememberMeCheckbox.checked = false;
    }

    const loginButton = document.getElementById("loginButton");
    if (loginButton) {
        loginButton.addEventListener("click", () => {
            const username = usernameElement.value;
            const password = passwordElement.value;
            const rememberMe = rememberMeCheckbox.checked;
            console.log(`尝试登录，用户名: ${username}`);

            // 根据复选框状态保存或清除账号密码
            if (rememberMe) {
                localStorage.setItem("username", username);
                localStorage.setItem("password", password);
                localStorage.setItem("rememberMe", "true");
            } else {
                localStorage.removeItem("username");
                localStorage.removeItem("password");
                localStorage.removeItem("rememberMe");
            }

            Login(username, password).then(result => {
                console.log("登录结果:", result);
                statusElement.innerText = result;
                
                // 只有在登录成功时才跳转到主界面
                if (result === "登录成功！") {
                    document.getElementById('app').innerHTML = `
                        <h1>自动化处理简历</h1>
                        <button id="sayHiButton">处理新招呼</button>
                        <button id="selectCampusButton">选择校区</button>
                        <button id="fetchResumeButton">抓取&下载简历</button>
                        <div class="form-group">
                            <label for="fileUpload">上传表格：</label>
                            <input type="file" id="fileUpload" accept=".csv, .xlsx, .xls">
                            <button id="uploadButton">上传</button>
                            <button id="updateOMOButton">更新至OMO</button>
                        </div>
                        
                        <div class="top-right-buttons">
                            <button id="toggleLogButton">显示日志</button>
                            <button id="logoutButton">退出登录</button>
                        </div>
                        <div id="logContainer" style="display: none; border: 1px solid #ccc; padding: 10px; max-height: 200px; overflow-y: auto;">
                            <pre id="logContent"></pre>
                        </div>
                        <div id="tableContainer"></div>
                    `;
                    
                    // 清空日志内容
                    const logContentElement = document.getElementById('logContent');
                    logContentElement.innerText = '';

                    document.getElementById('sayHiButton').addEventListener('click', () => {
                        console.log("打招呼");
                        SayHi().then(result => {
                            console.log("打招呼结果:", result);
                        }).catch(err => {
                            console.error("打招呼失败:", err);
                        });
                    });

                    let selectedCampuses = []; // 用于存储选中的校区

                    document.getElementById('selectCampusButton').addEventListener('click', () => {
                        console.log("选择校区");

                        const campuses = ["重庆", "杭州", "厦门", "广州", "北京", "天津", "郑州",
                                          "山西", "济南", "武汉", "南宁", "中山", "佛山", "深圳",
                                          "潍坊", "淄博", "苏州", "天津", "青岛", "上海", "西安",
                                          "长沙", "长春", "合肥", "南京", "成都", "东莞", "河北", "哈尔滨"];

                        let campusSelectionHTML = '<div id="campusSelection" style="display: flex; flex-wrap: wrap;">';
                        campuses.forEach(campus => {
                            campusSelectionHTML += `<label style="width: 20%;"><input type="checkbox" value="${campus}">${campus}</label>`;
                        });
                        campusSelectionHTML += '<button id="confirmCampusButton" style="width: 10%;">确定</button>';
                        campusSelectionHTML += '</div>';

                        document.getElementById('selectCampusButton').insertAdjacentHTML('afterend', campusSelectionHTML);

                        document.getElementById('confirmCampusButton').addEventListener('click', () => {
                            selectedCampuses = Array.from(document.querySelectorAll('#campusSelection input[type="checkbox"]:checked'))
                                                    .map(checkbox => checkbox.value);
                            const selectedText = selectedCampuses.length > 0 ? selectedCampuses.join(', ') : '无';
                            document.getElementById('campusSelection').innerHTML = `
                                <p>已选择校区：${selectedText}</p>
                                <button id="reselectCampusButton">重新选择</button>
                            `;

                            document.getElementById('reselectCampusButton').addEventListener('click', () => {
                                document.getElementById('campusSelection').remove();
                            });
                        });
                    });

                    document.getElementById('fetchResumeButton').addEventListener('click', () => {
                        if (selectedCampuses.length === 0) {
                            alert("请先选择校区！");
                            return;
                        }

                        console.log("抓取&下载简历");
                        RunRPA(selectedCampuses).then(() => {
                            console.log("RPA 抓取&下载简历运行成功");
                        }).catch(err => {
                            console.error("RPA 抓取&下载简历运行失败:", err);
                        });
                    });

                    let savedJsonData = null; // 用于保存上传的表格数据

                    document.getElementById('uploadButton').addEventListener('click', () => {
                        console.log("uploadButton");
                        const fileInput = document.getElementById('fileUpload');
                        const file = fileInput.files[0];
                        if (!file) {
                            alert("请选择一个文件！");
                            return;
                        }

                        const reader = new FileReader();
                        reader.onload = (e) => {
                            try {
                                const data = new Uint8Array(e.target.result);
                                const workbook = XLSX.read(data, { type: 'array' });
                                const firstSheetName = workbook.SheetNames[0];
                                const worksheet = workbook.Sheets[firstSheetName];

                                // 将工作表转换为 JSON 并保存
                                savedJsonData = XLSX.utils.sheet_to_json(worksheet, { header: 1 });
                                console.log("表格数据:", savedJsonData);

                                // 显示在前端
                                const tableContainer = document.getElementById('tableContainer');
                                displayTable(savedJsonData, tableContainer);
                            } catch (error) {
                                console.error("解析Excel文件出错:", error);
                                alert("解析Excel文件出错，请检查文件格式！错误信息：" + error.message);
                            }
                        };
                        reader.readAsArrayBuffer(file);
                    });

                    document.getElementById('updateOMOButton').addEventListener('click', () => {
                        if (!savedJsonData || savedJsonData.length <= 1) {
                            alert("请先上传并打开文件");
                            return;
                        }

                        // 添加结果列到表格
                        const tableContainer = document.getElementById('tableContainer');
                        const table = tableContainer.querySelector('table');
                        
                        // 检查并添加结果列
                        const firstHeader = table.rows[0].cells[0];
                        if (firstHeader.textContent !== '更新结果') {
                            // 为每一行添加新列
                            for (let i = 0; i < table.rows.length; i++) {
                                const newCell = table.rows[i].insertCell(0);
                                if (i === 0) {
                                    newCell.outerHTML = '<th>更新结果</th>';
                                }
                            }
                        }

                        // 字段映射表
                        const fieldMapping = {
                            // "序号": "row_no",
                            // "简历编号": "resumeNumber",
                            "意向课程": "regit_course",  
                            "手机": "mobile_phone",
                            "校区": "campus_id",         
                            "姓名": "name",
                            "性别": "gender",
                            "邮箱": "email",
                            "学历": "degree",           
                            "工作年限": "work_life",    
                            "应聘职位": "job_objective", 
                            "居住地": "domicile",       
                            "在职情况": "description",   
                            "来源": "source"
                        };

                        // 转换数据格式
                        const headers = savedJsonData[0];
                        const records = savedJsonData.slice(1).map((row, index) => {
                            const record = {};
                            headers.forEach((header, colIndex) => {
                                if (fieldMapping[header]) {  
                                    record[fieldMapping[header]] = row[colIndex] || '';
                                }
                            });
                            // record['row_no'] = index + 1;
                            return record;
                        });

                        console.log("发送的数据:", records);  

                        UpdateOmo(records).then((results) => {
                            // 更新每一行的结果
                            results.forEach((result, index) => {
                                const row = table.rows[index + 1];  
                                const resultCell = row.cells[0];
                                
                                // 显示状态
                                resultCell.textContent = `${result.msg_type || ''}`;
                                
                                // 根据状态设置样式
                                switch(result.msg_type) {
                                    case '失败':
                                        resultCell.className = 'status-error';
                                        break;
                                    case '重复':
                                        resultCell.className = 'status-duplicate';
                                        break;
                                    case '重置':
                                        resultCell.className = 'status-success';
                                        break;
                                    case '成功':
                                        resultCell.className = 'status-success';
                                        break;
                                }
                            });
                            
                            console.log("更新完成");
                        }).catch(err => {
                            console.error("更新失败:", err);
                            alert("更新失败，请检查控制台日志。");
                        });
                    });

                    document.getElementById('logoutButton').addEventListener('click', () => {
                        console.log("退出登录");
                        // 直接重新加载页面
                        window.location.reload();
                    });

                    document.getElementById('toggleLogButton').addEventListener('click', () => {
                        const logContainer = document.getElementById('logContainer');
                        if (logContainer.style.display === "none") {
                            GetLogs().then(data => {
                                const logContentElement = document.getElementById('logContent');
                                logContentElement.innerText = data;
                                logContainer.style.display = "block";
                                document.getElementById('toggleLogButton').innerText = "隐藏日志";
                            }).catch(err => console.error("获取日志失败:", err));
                        } else {
                            logContainer.style.display = "none";
                            document.getElementById('toggleLogButton').innerText = "显示日志";
                        }
                    });

                    function displayTableData(data) {
                        const tableContent = document.getElementById('tableContent');
                        let html = '<table border="1"><tr>';

                        // 表头
                        data[0].forEach(header => {
                            html += `<th>${header}</th>`;
                        });
                        html += '</tr>';

                        // 表格内容
                        for (let i = 1; i < data.length; i++) {
                            html += '<tr>';
                            data[i].forEach(cell => {
                                html += `<td>${cell}</td>`;
                            });
                            html += '</tr>';
                        }
                        html += '</table>';

                        tableContent.innerHTML = html;
                    }

                    setInterval(() => {
                        GetLogs().then(data => {
                            // 追加新日志内容
                            logContentElement.innerText += data + '\n';
                        }).catch(err => console.error("获取日志失败:", err));
                    }, 5000);

                } else {
                    alert("登录失败，请检查用户名和密码!");
                }
            }).catch(err => {
                alert("登录出错:", err);
            });
        });
    } else {
        console.error("登陆按钮出错！");
    }
});

// 显示表格的函数
function displayTable(jsonData, tableContainer) {
    let html = '<table border="1"><tr>';

    // 表头
    jsonData[0].forEach(header => {
        html += `<th>${header}</th>`;
    });
    html += '</tr>';

    // 表格内容
    for (let i = 1; i < jsonData.length; i++) {
        html += '<tr>';
        jsonData[i].forEach(cell => {
            html += `<td>${cell}</td>`;
        });
        html += '</tr>';
    }
    html += '</table>';

    // 添加筛选栏
    const filterHTML = `
        <div class="filter-container">
            <div class="filter-item">
                <input type="checkbox" id="successOnly">
                <label for="successOnly">只显示成功结果</label>
            </div>
            <div class="filter-item">
                <input type="text" id="campusFilter" placeholder="输入校区">
            </div>
            <div class="filter-item">
                <input type="text" id="positionFilter" placeholder="输入岗位">
            </div>
            <button id="applyFilter">确认筛选</button>
            <button id="packResumeButton">打包简历</button>
        </div>
    `;

    tableContainer.innerHTML = filterHTML + html;

    // 添加筛选功能
    document.getElementById('applyFilter').addEventListener('click', async () => {
        const successOnly = document.getElementById('successOnly').checked;
        const campusFilter = document.getElementById('campusFilter').value.toLowerCase();
        const positionFilter = document.getElementById('positionFilter').value.toLowerCase();

        const table = tableContainer.querySelector('table');
        const rows = Array.from(table.querySelectorAll('tr')).slice(1); // 跳过表头
        
        // 获取下载路径
        const downloadPath = await GetDownloadPath();
        console.log("简历下载路径:", downloadPath);

        // 筛选和检查文件
        const filteredData = [];
        const headers = Array.from(table.querySelectorAll('th')).map(th => th.textContent);
        filteredData.push(headers);

        for (const row of rows) {
            const cells = row.cells;
            const status = cells[0].textContent.trim();
            const campus = cells[findColumnIndex('校区', table)].textContent.toLowerCase();
            const position = cells[findColumnIndex('应聘职位', table)].textContent.toLowerCase();

            const showRow = (!successOnly || status === '成功' || status === '重置') &&
                          (!campusFilter || campus.includes(campusFilter)) &&
                          (!positionFilter || position.includes(positionFilter));

            row.style.display = showRow ? '' : 'none';

            if (showRow) {
                const rowData = Array.from(cells).map(cell => cell.textContent);
                filteredData.push(rowData);
            
                // 获取简历相关信息
                const name = cells[findColumnIndex('姓名', table)].textContent;
                const jobPosition = cells[findColumnIndex('应聘职位', table)].textContent;
                const campus = cells[findColumnIndex('校区', table)].textContent;
            
                console.log("正在检查简历:", {
                    name: name,
                    position: jobPosition,
                    campus: campus
                });
            
                try {
                    const filePath = await window.go.main.App.CheckResumeFile(name, jobPosition, campus);
                    if (filePath) {
                        console.log(`✅ 定位到 ${name} 的简历在：${filePath}`);
                    } else {
                        console.log(`❌ 未找到 ${name} 的简历文件`);
                    }
                } catch (err) {
                    console.error(`检查 ${name} 的简历文件时出错:`, err);
                }
            }
            
        }

        console.log("筛选条件:", { successOnly, campusFilter, positionFilter });
        console.log("筛选后的数据:", filteredData);
    });

    // 在筛选栏中添加打包按钮的事件监听
    document.getElementById('packResumeButton').addEventListener('click', async () => {
        console.log("开始打包简历...");
        
        // 获取表格中的所有数据
        const table = document.querySelector('table');
        const rows = Array.from(table.querySelectorAll('tr')).slice(1); // 跳过表头
        
        // 按校区和岗位组织数据
        const groupedData = {};
        
        for (const row of rows) {
            const cells = row.cells;
            const name = cells[findColumnIndex('姓名', table)].textContent;
            const jobPosition = cells[findColumnIndex('应聘职位', table)].textContent;
            const campus = cells[findColumnIndex('校区', table)].textContent;
            
            // 检查简历文件
            try {
                const filePath = await window.go.main.App.CheckResumeFile(name, jobPosition, campus);
                if (filePath) {
                    // 初始化校区对象
                    if (!groupedData[campus]) {
                        groupedData[campus] = {};
                    }
                    // 初始化岗位数组
                    if (!groupedData[campus][jobPosition]) {
                        groupedData[campus][jobPosition] = [];
                    }
                    // 添加文件路径
                    groupedData[campus][jobPosition].push(filePath);
                }
            } catch (err) {
                console.error(`获取 ${name} 的简历文件时出错:`, err);
            }
        }
        
        console.log("已组织的简历数据:", groupedData);
        
        // 调用后端进行打包
        try {
            const result = await window.go.main.App.PackResumes(groupedData);
            if (result) {
                console.log("简历打包成功:", result);
                alert("简历打包成功！");
            }
        } catch (err) {
            console.error("打包简历时出错:", err);
            alert("打包简历时出错: " + err);
        }
    });
}

// 辅助函数：找到指定列名的索引
function findColumnIndex(columnName, table) {
    const headers = Array.from(table.querySelectorAll('th'));
    return headers.findIndex(th => th.textContent.includes(columnName));
}
