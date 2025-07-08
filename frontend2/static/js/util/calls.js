import { API, accountActive, accountInactive } from "../index.js";
import { navigateToURL} from "../index.js";


async function callWithRefresh(endpoint, method, headers, body) {
    async function call() {
        const token = localStorage.getItem('token');
        headers['Authorization'] = token;
        const res = await fetch(`${API}/${endpoint}`, {
            method: method,
            headers: headers,
            body: body
        });
        return res;
    }
    const res = await call();
    if (res.ok) {
        return res;
    } else if (res.status === 401) {
        await refreshAuth();
        return await call();
    } else {
        return res;
    }
}


async function login(e) {
    e.preventDefault();
    localStorage.setItem('number', e.target.number.value);
    const number = Number(e.target.number.value);
    const form = e.target;
    const data = {
        number: number,
        password: form.password.value
    };
    const res = await fetch(`${API}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });

    if (res.ok) {
        var resp = await res.json();
        var token = resp.token;
        var refresh_token = resp.refresh_token;
        localStorage.setItem('token', token);
        localStorage.setItem('refresh_token', refresh_token);
        alert('Login successful');
        navigateToURL(`#/account/${number}`);
        accountActive();
    } else {
        alert('Login failed');
    }
}


async function getAccount() {
    const number = Number(localStorage.getItem('number'));
    const token = localStorage.getItem('token');
    const res = await callWithRefresh(`account/${number}`, 'GET', { 'Authorization': `${token}` }, null);
    if (res.ok) {
        const account = await res.json();
        return account;
    } else {
        alert('Get account failed');
        const account = {};
        return account;
    }
}


function logout() {
    localStorage.removeItem('number');
    localStorage.removeItem('token');
    localStorage.removeItem('refresh_token');
    navigateToURL('#/login');
    accountInactive();
}

async function refreshAuth() {
    const refresh_token = localStorage.getItem('refresh_token');
    const data = {
        refresh_token: refresh_token
    };
    const res = await fetch(`${API}/refresh`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });
    if (res.ok) {
        var resp = await res.json();
        const token = resp.token;
        const refresh_token = resp.refresh_token;
        localStorage.setItem('token', token);
        localStorage.setItem('refresh_token', refresh_token);
    } else if (res.status === 401) {
        alert('Session expired');
        logout();
    }
    else {
        alert('Refresh failed');
        logout();
    }
}

async function transfer(e) {
    e.preventDefault();
    const number = Number(localStorage.getItem('number'));
    const token = localStorage.getItem('token');
    const form = e.target;
    const data = {
        recipient: Number(form.to.value),
        amount: Number(form.amount.value)
    };
    const res = await callWithRefresh(`transfer/${number}`, 'POST', { 'Content-Type': 'application/json', 'Authorization': `${token}` }, JSON.stringify(data));
    if (res.ok) {
        alert('Transfer successful');
        navigateToURL(`#/account/${number}`);
    } else {
        alert('Transfer failed');
    }
}


export { login, getAccount, logout, refreshAuth, transfer };



