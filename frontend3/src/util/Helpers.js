import { navigateTo } from "../main.js";

export function cacheAccount(account) {
    localStorage.setItem('account', JSON.stringify(account));
}

export function deleteAccount() {
    localStorage.removeItem('account');
}

export function loadAccount() {
    const account = localStorage.getItem('account');
    if (account) {
        try {
            return JSON.parse(account);
        } catch {
            return null;
        }
    }
    return null;
}


export function accountInactive() {
    const nav = document.querySelector("nav");
    nav.innerHTML = `
        <a href="#/login" class="nav__link" data-link>Login</a>
        <a href="#/register" class="nav__link" data-link>Register</a>`
}

export function accountActive() {
    const nav = document.querySelector("nav");
    const number = localStorage.getItem("number")
    nav.innerHTML = `
        <a href="#/account/${number}" class="nav__link" data-link>Account</a>
        <a href="#/account/${number}/transfer" class="nav__link" data-link>Transfer</a>
        <a href="#/account/${number}/game" class="nav__link" data-link>Game</a>`
}

export function loggedIn() {
    if (localStorage.getItem('token')) {
        return true;
    }
    return false;
}

export function accountAccess(id) {
    const number = localStorage.getItem('number');
    if (id !== number) {
        if (loggedIn()) {
            alert('You are not logged in to this account');
            navigateTo(`#/account/${number}`);
            return false;
        }
        else {
            alert('You are not logged in');
            navigateTo("#/login");
            return false;
        }
    }
    return true;
}

export function notFound() {
    alert('Page not found');
    if (loggedIn()) {
        navigateTo(`#/account/${localStorage.getItem('number')}`);
    }
    else {
        navigateTo('#/');
    }
}

export function eventHandler(event) {
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
}