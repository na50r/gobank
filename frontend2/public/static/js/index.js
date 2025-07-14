import Login from "./views/Login.js";
import Register from "./views/Register.js";
import Account from "./views/Account.js";
import Transfer from "./views/Transfer.js";
import Game from "./views/Game.js";
import { deleteAccount } from "./views/Account.js";
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
                deleteAccount();
                console.log("Transaction detected, reloading account page");
                navigateTo(`/account/${number}`);
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
        <a href="/login" class="nav__link" data-link>Login</a>
        <a href="/register" class="nav__link" data-link>Register</a>`
}

function accountActive() {
    const nav = document.querySelector("nav");
    const number = localStorage.getItem("number")
    nav.innerHTML = `
        <a href="/account/${number}" class="nav__link" data-link>Account</a>
        <a href="/transfer" class="nav__link" data-link>Transfer</a>
        <a href="/game" class="nav__link" data-link>Game</a>`
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

const routes = [
    { path: "/", view: loggedIn() ? Account : Login },
    { path: "/login", view: Login },
    { path: "/register", view: Register },
    { path: "/transfer", view: Transfer },
    { path: "/account/:id", view: Account },
    { path: "/game", view: Game }
]

async function router() {
    const potentialMatches = routes.map(route => {
        const path = location.pathname
        return { route: route, result: path.match(pathToRegex(route.path)) }
    });
    let match = potentialMatches.find(potentialMatch => potentialMatch.result !== null);
    if (!match) {
        match = {
            route: routes[0],
            result: [location.pathname]
        }
    }
    const view = new match.route.view(getParams(match));
    const container = await view.getHtml();
    const app = document.querySelector("#app");
    app.replaceChildren(container);
}

export function navigateTo(url) {
    history.pushState(null, null, url);
    router();
}

function navBehaviour() {
    document.body.addEventListener("click", e => {
        if (e.target.matches("[data-link]")) {
            e.preventDefault();
            navigateTo(e.target.href);
        }
    })
}


function startBehaviour() {
    navBehaviour();
    if (loggedIn()) {
        const number = localStorage.getItem("number");
        navigateTo(`/account/${number}`);
    }
    router();
}

window.onload = () => {
    if (loggedIn()) {
        accountActive();
        return;
    }
    accountInactive();
};

window.addEventListener("popstate", router);
window.addEventListener("DOMContentLoaded", startBehaviour);
window.addEventListener("load", router);