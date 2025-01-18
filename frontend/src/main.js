import './style.css';
import './app.css';
import * as XLSX from 'xlsx';
import { Login, GetLogs, RunRPA } from '../wailsjs/go/main/App';
import { UpdateOmo } from '../wailsjs/go/main/App';
window.addEventListener("DOMContentLoaded", () => {
    console.log("DOM fully loaded and parsed");

    const usernameElement = document.getElementById("username");
    const passwordElement = document.getElementById("password");
    const statusElement = document.getElementById("status");

    if (usernameElement && passwordElement) {
        usernameElement.value = "wujinxuan@kelote.com";
        passwordElement.value = "chill000";
        console.log("Default username and password set");
    } else {
        console.error("Username or password element not found");
    }

    const loginButton = document.getElementById("loginButton");
    if (loginButton) {
        loginButton.addEventListener("click", () => {
            const username = usernameElement.value;
            const password = passwordElement.value;
            console.log(`Login attempt with username: ${username}, password: ${password}`);

            Login(username, password).then(result => {
                console.log("Login successful:", result);
                statusElement.innerText = result;

                document.getElementById('app').innerHTML = `
                    <h1>自动化处理简历</h1>
                    <button id="selectCampusButton">选择校区</button>
                    <button id="fetchResumeButton">抓取&下载简历</button>
                    <div class="form-group">
                        <label for="fileUpload">上传表格：</label>
                        <input type="file" id="fileUpload" accept=".csv, .xlsx, .xls">
                        <button id="uploadButton">上传</button>
                    </div>
                    
                    <button id="updateOMOButton">更新至OMO</button>
                    <button id="logoutButton">退出登录</button>
                    <button id="toggleLogButton">显示日志</button>
                    <div id="logContainer" style="display: none; border: 1px solid #ccc; padding: 10px; max-height: 200px; overflow-y: auto;">
                        <pre id="logContent"></pre>
                    </div>
                    <div id="tableContainer"></div>
                `;

                // 清空日志内容
                const logContentElement = document.getElementById('logContent');
                logContentElement.innerText = '';

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
                            let html = '<table border="1"><tr>';

                            // 表头
                            savedJsonData[0].forEach(header => {
                                html += `<th>${header}</th>`;
                            });
                            html += '</tr>';

                            // 表格内容
                            for (let i = 1; i < savedJsonData.length; i++) {
                                html += '<tr>';
                                savedJsonData[i].forEach(cell => {
                                    html += `<td>${cell}</td>`;
                                });
                                html += '</tr>';
                            }
                            html += '</table>';

                            tableContainer.innerHTML = html;
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
                        "序号": "row_no",
                        "简历编号": "resumeNumber",
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
                        record['row_no'] = index + 1;
                        return record;
                    });

                    console.log("发送的数据:", records);  // 调试用

                    UpdateOmo(records).then((result) => {
                        // 检查返回的结果
                        if (result && result.msg_type === '失败') {
                            // 更新表格中的结果列
                            const firstRow = table.rows[1];  // 第一个数据行
                            const resultCell = firstRow.cells[0];
                            resultCell.textContent = result.msg || result.msg_base || '更新失败';
                            resultCell.style.backgroundColor = '#ffebee';  // 失败用红色背景
                            
                            // 显示错误消息
                            alert(`更新失败：${result.msg || result.msg_base}`);
                        } else {
                            // 更新成功
                            const firstRow = table.rows[1];
                            const resultCell = firstRow.cells[0];
                            resultCell.textContent = '更新成功';
                            resultCell.style.backgroundColor = '#e8f5e9';  // 成功用绿色背景
                            alert("更新成功！");
                        }
                    }).catch(err => {
                        console.error("更新失败:", err);
                        alert("更新失败，请检查控制台日志。");
                    });
                });

                document.getElementById('logoutButton').addEventListener('click', () => {
                    console.log("退出登录");
                    document.getElementById('app').innerHTML = `
                        <h1>Login</h1>
                        <input type="text" id="username" placeholder="Username">
                        <input type="password" id="password" placeholder="Password">
                        <button id="loginButton">登陆</button>
                        <p id="status"></p>
                    `;
                    document.getElementById('loginButton').addEventListener('click', () => {
                        // 重新调用登录逻辑
                    });
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

            }).catch(err => {
                alert("Login failed:", err);
                statusElement.innerText = "Login failed: " + err;
            });
        });
    } else {
        console.error("Login button not found");
    }
});
