import { Login } from '../wailsjs/go/main/App';

window.addEventListener("DOMContentLoaded", () => {
    console.log("DOM fully loaded and parsed");

    const usernameElement = document.getElementById("username");
    const passwordElement = document.getElementById("password");
    const statusElement = document.getElementById("status");

    if (usernameElement && passwordElement) {
        // 设置默认账号和密码
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

            // 调用Wails的Login函数
            Login(username, password).then(result => {
                console.log("Login successful:", result);
                statusElement.innerText = result;

                // 登录成功后，更新页面内容
                document.getElementById('app').innerHTML = `
                    <h1>自动化处理简历</h1>
                    <button id="zhaopinButton">登陆智联招聘</button>
                    <button id="fetchResumeButton">抓取&下载简历</button>
                    <button id="updateOMOButton">更新至OMO</button>
                    <button id="logoutButton">退出登录</button>
                `;

                // 为新按钮添加事件监听器
                document.getElementById('zhaopinButton').addEventListener('click', () => {
                    console.log("登陆智联招聘");
                    // 在这里添加登陆智联招聘的逻辑
                });

                document.getElementById('fetchResumeButton').addEventListener('click', () => {
                    console.log("抓取&下载简历");
                    // 在这里添加抓取&下载简历的逻辑
                });

                document.getElementById('updateOMOButton').addEventListener('click', () => {
                    console.log("更新至OMO");
                    // 在这里添加更新至OMO的逻辑
                });

                document.getElementById('logoutButton').addEventListener('click', () => {
                    console.log("退出登录");
                    // 返回到登录页面
                    document.getElementById('app').innerHTML = `
                        <h1>Login</h1>
                        <input type="text" id="username" placeholder="Username">
                        <input type="password" id="password" placeholder="Password">
                        <button id="loginButton">登陆</button>
                        <p id="status"></p>
                    `;
                    // 重新绑定登录按钮事件
                    document.getElementById('loginButton').addEventListener('click', () => {
                        // 重新调用登录逻辑
                    });
                });

            }).catch(err => {
                console.error("Login failed:", err);
                statusElement.innerText = "Login failed: " + err;
            });
        });
    } else {
        console.error("Login button not found");
    }
});
