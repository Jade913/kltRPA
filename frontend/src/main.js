// import './style.css';
// import './app.css';

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

import { Login } from '../wailsjs/go/main/App';

document.getElementById("loginButton").onclick = function() {
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;
    
    Login(username, password).then(result => {
        document.getElementById("status").innerText = result;
    }).catch(err => {
        document.getElementById("status").innerText = "Login failed: " + err;
    });
};