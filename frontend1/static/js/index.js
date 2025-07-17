import Login from "./views/Login.js";
import Register from "./views/Register.js";
import Account from "./views/Account.js";
import Transfer from "./views/Transfer.js";
import Game from "./views/Game.js";
import { accountActive, accountInactive, loggedIn, eventHandler, deleteAccount, notFound } from "./util/Helpers.js";
import { config } from "./config.js";
export const API = config.apiUrl;

const evtSource = new EventSource(`${API}/events`);
evtSource.addEventListener("msg", eventHandler);
evtSource.onerror = (err) => {console.log("sse error", err);};

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
        { path: "#/", view: Login },
        { path: "#/login", view: Login },
        { path: "#/register", view: Register },
        { path: "#/account/:id/transfer", view: Transfer },
        { path: "#/account/:id", view: Account },
        { path: "#/account/:id/game", view: Game }
    ]

    const potentialMatches = routes.map(route => {
        const currHash = location.hash || "#/";
        return { route: route, result: currHash.match(pathToRegex(route.path)) };
    })

    let match = potentialMatches.find(potentialMatch => potentialMatch.result !== null);
    if (!match) {
        notFound();
        return;
    }
    const view = new match.route.view(getParams(match));
    const container = await view.getHtml();
    const app = document.querySelector("#app");
    app.replaceChildren(container);
}

function startBehaviour() {
    if (location.pathname.endsWith("index.html")) {
        location.pathname = location.pathname.replace("index.html", "");
        navigateTo("#/");
    }

    if (loggedIn() && !location.pathname.startsWith("#/account")) {
        navigateTo(`#/account/${localStorage.getItem('number')}`);
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

export function navigateTo(url) {
    location.hash = url;
}

window.addEventListener("hashchange", router);
document.addEventListener("DOMContentLoaded", startBehaviour);
