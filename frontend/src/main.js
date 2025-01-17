import { Login, GetLogs, RunRPA, GetLatestTable } from '../wailsjs/go/main/App';
// import { updateOMOData } from '../wailsjs/go/main/App';
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
                    <button id="importTableButton">导入表格</button>
                    <button id="updateOMOButton">更新至OMO</button>
                    <button id="logoutButton">退出登录</button>
                    <button id="toggleLogButton">显示日志</button>
                    <div id="logContainer" style="display: none; border: 1px solid #ccc; padding: 10px; max-height: 200px; overflow-y: auto;">
                        <pre id="logContent"></pre>
                    </div>
                    <div id="tableContent"></div>
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

                document.getElementById('updateOMOButton').addEventListener('click', () => {
                    const tableContent = document.getElementById('tableContent');
                    const rows = tableContent.querySelectorAll('tr');

                    if (rows.length <= 1) {
                        alert("请先打开文件");
                        return;
                    }

                    // 逐条传入OMO接口
                    for (let i = 1; i < rows.length; i++) {
                        const cells = rows[i].querySelectorAll('td');
                        const rowData = Array.from(cells).map(cell => cell.innerText);

                        // 调用后端接口
                        // updateOMO(rowData).then(() => {
                        //     console.log(`第${i}行数据更新成功`);
                        // }).catch(err => {
                        //     console.error(`第${i}行数据更新失败:`, err);
                        // });
                    }
                });

                // function updateOMO(rowData) {
                //     // 后端方法 updateOMOData
                //     return updateOMOData(rowData);
                // }

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
                            console.log("最新日志内容:", data);
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

                document.getElementById('importTableButton').addEventListener('click', () => {
                    // 创建一个隐藏的文件输入元素
                    const fileInput = document.createElement('input');
                    fileInput.type = 'file';
                    fileInput.accept = '.xlsx,.xls,.xlsm'; // 只接受Excel文件

                    fileInput.onchange = async (event) => {
                        try {
                            const file = event.target.files[0];
                            if (file) {
                                console.log(`选择的文件: ${file.name}`);
                                await importTableFromExcel(file);
                            }
                        } catch (error) {
                            console.error("文件处理出错:", error);
                            alert("文件处理出错，请重试！");
                        }
                    };

                    // 触发文件选择对话框
                    fileInput.click();
                });

                async function importTableFromExcel(file) {
                    const reader = new FileReader();
                    reader.onload = (e) => {
                        try {
                            const data = new Uint8Array(e.target.result);
                            const workbook = XLSX.read(data, { type: 'array' });
                            const firstSheetName = workbook.SheetNames[0];
                            const worksheet = workbook.Sheets[firstSheetName];

                            // 将工作表转换为 JSON
                            const jsonData = XLSX.utils.sheet_to_json(worksheet, { header: 1 });
                            console.log("表格数据:", jsonData);

                            // 显示在前端
                            displayTableData(jsonData);
                        } catch (error) {
                            console.error("解析Excel文件出错:", error);
                            alert("解析Excel文件出错，请检查文件格式！");
                        }
                    };
                    reader.readAsArrayBuffer(file);
                }

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
                        console.log("日志内容:", data);
                        // 追加新日志内容
                        logContentElement.innerText += data + '\n';
                    }).catch(err => console.error("获取日志失败:", err));
                }, 5000);

            }).catch(err => {
                console.error("Login failed:", err);
                statusElement.innerText = "Login failed: " + err;
            });
        });
    } else {
        console.error("Login button not found");
    }
});