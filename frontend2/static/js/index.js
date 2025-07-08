import Login from "./views/Login.js";
import Register from "./views/Register.js";
import Account from "./views/Account.js";
import Transfer from "./views/Transfer.js";

const API_URL = "http://localhost:3000";
export const API = API_URL;

window.onload = () => {
    if (localStorage.getItem('token')) {
        accountActive();
        return;
    }
    accountInactive();
};

function accountInactive() {
    const nav = document.querySelector("nav");
    nav.innerHTML = `
        <a href="#/" class="nav__link" data-link>Login</a>
        <a href="#/register" class="nav__link" data-link>Register</a>`
}

function accountActive() {
    const nav = document.querySelector("nav");
    nav.innerHTML = `
        <a href="#/" class="nav__link" data-link>Account</a>
        <a href="#/transfer" class="nav__link" data-link>Transfer</a>`
}

export { accountActive, accountInactive };

const pathToRegex = path =>
    new RegExp("^" + path.replace(/\//g, "\\/").replace(/:\w+/g, "([^\\/]+)") + "$");


const getParams = match => {
    const values = match.result.slice(1);
    const keys = Array.from(match.route.path.matchAll(/:(\w+)/g)).map(result => result[1]);
    return Object.fromEntries(keys.map((key, i) => {
        return [key, values[i]];
    }))
}


const navigateTo = url => {
    history.pushState(null, null, url);
    router();
}

export const navigateToURL = navigateTo;

const router = async () => {
    const routes = [
        { path: "#/", view: Login },
        { path: "#/login", view: Login },
        { path: "#/register", view: Register },
        { path: "#/transfer", view: Transfer },
        { path: "#/account/:id", view: Account },
    ]

    // Test each route for potential match
    const potentialMatches = routes.map(route => {
        const currHash = location.hash || "#/";
        return { route: route, result: currHash.match(pathToRegex(route.path)) };
    })

    let match = potentialMatches.find(potentialMatch => potentialMatch.result !== null); 
    if (!match) {
        match = {
            route: routes[0], result: [location.pathname]
        }
    }
    const view = new match.route.view(getParams(match));
    const container = await view.getHtml();
    const app = document.querySelector("#app");
    app.replaceChildren(container);
};

window.addEventListener("popstate", router);

document.addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("click", e => {
        if (e.target.matches("[data-link]")) {
            e.preventDefault();
            navigateTo(e.target.href);
        }
    })
    router();
});