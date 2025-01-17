import { Login } from '../wailsjs/go/main/App';
import { RunRPA } from '../wailsjs/go/main/App';
import { GetLogs } from '../wailsjs/go/main/App';
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
                    <button id="zhaopinButton">登陆智联招聘</button>
                    <button id="fetchResumeButton">抓取&下载简历</button>
                    <button id="updateOMOButton">更新至OMO</button>
                    <button id="logoutButton">退出登录</button>
                    <button id="toggleLogButton">显示日志</button>
                    <div id="logContainer" style="display: none; border: 1px solid #ccc; padding: 10px; max-height: 200px; overflow-y: auto;">
                        <pre id="logContent"></pre>
                    </div>
                `;

                // 清空日志内容
                const logContentElement = document.getElementById('logContent');
                logContentElement.innerText = '';

                document.getElementById('zhaopinButton').addEventListener('click', () => {
                    console.log("登陆智联招聘");
                });

                document.getElementById('fetchResumeButton').addEventListener('click', () => {
                    console.log("抓取&下载简历");
                    RunRPA();
                });

                document.getElementById('updateOMOButton').addEventListener('click', () => {
                    console.log("更新至OMO");
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