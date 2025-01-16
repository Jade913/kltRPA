import './style.css';
import './app.css';

import { Login } from '../../wailsjs/go/main/App';

window.addEventListener("DOMContentLoaded", () => {
    const usernameElement = document.getElementById("username");
    const passwordElement = document.getElementById("password");
    const statusElement = document.getElementById("status");

    // 设置默认账号和密码
    usernameElement.value = "wujinxuan@kelote.com";
    passwordElement.value = "chill000";

    document.getElementById("loginButton").addEventListener("click", () => {
        const username = usernameElement.value;
        const password = passwordElement.value;
        
        // 调用Wails的Login函数
        Login(username, password).then(result => {
            statusElement.innerText = result;
        }).catch(err => {
            statusElement.innerText = "Login failed: " + err;
        });
    });
});

// import logo from './assets/images/logo-universal.png';
// import {Greet} from '../wailsjs/go/main/App';

// document.querySelector('#app').innerHTML = `
//     <img id="logo" class="logo">
//       <div class="result" id="result">Please enter your name below 👇</div>
//       <div class="input-box" id="input">
//         <input class="input" id="name" type="text" autocomplete="off" />
//         <button class="btn" onclick="greet()">Greet</button>
//       </div>
//     </div>
// `;
// document.getElementById('logo').src = logo;

// let nameElement = document.getElementById("name");
// nameElement.focus();
// let resultElement = document.getElementById("result");

// // Setup the greet function
// window.greet = function () {
//     // Get name
//     let name = nameElement.value;

//     // Check if the input is empty
//     if (name === "") return;

//     // Call App.Greet(name)
//     try {
//         Greet(name)
//             .then((result) => {
//                 // Update result with data back from App.Greet()
//                 resultElement.innerText = result;
//             })
//             .catch((err) => {
//                 console.error(err);
//             });
//     } catch (err) {
//         console.error(err);
//     }
// };

// import { Login } from '../wailsjs/go/main/App';

// // 设置默认账号和密码
// document.getElementById("username").value = "wujinxuan@kelote.com";
// document.getElementById("password").value = "chill000";

// document.getElementById("loginButton").onclick = function() {
//     const username = document.getElementById("username").value;
//     const password = document.getElementById("password").value;
    
//     Login(username, password).then(result => {
//         document.getElementById("status").innerText = result;
//     }).catch(err => {
//         document.getElementById("status").innerText = "Login failed: " + err;
//     });
// };