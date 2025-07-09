import Login from "./views/Login.js";
import Register from "./views/Register.js";
import Account from "./views/Account.js";
import Transfer from "./views/Transfer.js";
import Game from "./views/Game.js";
export const API = document.body.dataset.apiUrl;

const evtSource = new EventSource(`${API}/stream`);
evtSource.onmessage = function (event) {
    const data = JSON.parse(event.data);
    switch (data.type) {
        case "transaction":
            const number = Number(localStorage.getItem('number'));
            console.log(`[Transaction] ${data.data.sender} -> ${data.data.amount} -> ${data.data.recipient}`);
            const sender = Number(data.data.sender);
            const recipient = Number(data.data.recipient);
            if (sender === number || recipient === number) {
                console.log("Transaction detected, reloading account page");
                location.hash = `#/account/${number}`;
                location.reload();
            }
            break;
        case "chat":
            console.log(`[Chat] ${data.data.name}: ${data.data.msg}`);
            break;
        default:
            console.log("Unknown event type", data);
    }
};


function accountInactive() {
    const nav = document.querySelector("nav");
    nav.innerHTML = `
        <a href="#/login" class="nav__link" data-link>Login</a>
        <a href="#/register" class="nav__link" data-link>Register</a>`
}

function accountActive() {
    const nav = document.querySelector("nav");
    const number = localStorage.getItem("number")
    nav.innerHTML = `
        <a href="#/account/${number}" class="nav__link" data-link>Account</a>
        <a href="#/transfer" class="nav__link" data-link>Transfer</a>
        <a href="#/game" class="nav__link" data-link>Game</a>`
}

function loggedIn() {
    if (localStorage.getItem('token')) {
        return true;
    }
    return false;
}

export { accountActive, accountInactive };


function pathToRegex(path) {
    return new RegExp("^" + path.replace(/\//g, "\\/").replace(/:\w+/g, "([^\\/]+)") + "$");
}


function getParams(match) {
    const values = match.result.slice(1);
    const keys = Array.from(match.route.path.matchAll(/:(\w+)/g)).map(result => result[1]);
    return Object.fromEntries(keys.map((key, i) => {
        return [key, values[i]];
    }))
}

async function router() {
    const routes = [
        { path: "#/", view: loggedIn() ? Account : Login },
        { path: "#/login", view: Login },
        { path: "#/register", view: Register },
        { path: "#/transfer", view: Transfer },
        { path: "#/account/:id", view: Account },
        { path: "#/game", view: Game }
    ]

    const potentialMatches = routes.map(route => {
        const currHash = location.hash || "#/";
        return { route: route, result: currHash.match(pathToRegex(route.path)) };
    })

    let match = potentialMatches.find(potentialMatch => potentialMatch.result !== null);
    if (!match) {
        match = {
            route: routes[0], result: [location.hash]
        }
    }
    const view = new match.route.view(getParams(match));
    const container = await view.getHtml();
    const app = document.querySelector("#app");
    app.replaceChildren(container);
}

function startBehaviour() {
    if (location.pathname.endsWith("index.html")) {
        location.pathname = location.pathname.replace("index.html", "");
        location.hash = "#/";
    }

    if (loggedIn()) {
        const number = localStorage.getItem("number");
        location.hash = `#/account/${number}`;
    }

    document.body.addEventListener("click", e => {
        if (e.target.matches("[data-link]")) {
            e.preventDefault();
            location.hash = e.target.getAttribute("href");
        }
    });
    router();
}

window.onload = () => {
    if (loggedIn()) {
        accountActive();
        return;
    }
    accountInactive();
};

window.addEventListener("hashchange", router);
document.addEventListener("DOMContentLoaded", startBehaviour);
