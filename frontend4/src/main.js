import Login from "./views/Login.js";
import Register from "./views/Register.js";
import Account from "./views/Account.js";
import Transfer from "./views/Transfer.js";
import Game from "./views/Game.js";
import { loggedIn, notFound, eventHandler, accountInactive, accountActive } from "./util/Helpers.js";
export const API = "http://localhost:3000";

const evtSource = new EventSource(`${API}/events`);
evtSource.addEventListener("msg", eventHandler);
onerror = (err) => {console.log("sse error", err);};

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
  { path: "/", view: Login },
  { path: "/login", view: Login },
  { path: "/register", view: Register },
  { path: "/account/:id/transfer", view: Transfer },
  { path: "/account/:id", view: Account },
  { path: "/account/:id/game", view: Game }
]

async function router() {
  const potentialMatches = routes.map(route => {
    const path = location.pathname
    return { route: route, result: path.match(pathToRegex(route.path)) }
  });
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
  if (loggedIn() && !location.pathname.startsWith("/account")) {
    navigateTo(`/account/${localStorage.getItem('number')}`);
  }
  if (loggedIn() && !location.pathname === '/') {
    navigateTo(`/account/${localStorage.getItem('number')}`);
  }
  navBehaviour();
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