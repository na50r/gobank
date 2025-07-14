import { accountActive, accountInactive, deleteAccount } from "./Helpers.js";
import { navigateTo } from "../index.js";
export const API = document.body.dataset.apiUrl;

export async function callWithRefresh(endpoint, method, headers, body) {
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

export async function login(e) {
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
        navigateTo(`/account/${number}`);
        accountActive();
    } else {
        alert('Login failed');
    }
}


export async function getAccount(number) {
    const token = localStorage.getItem('token');
    const res = await callWithRefresh(`account/${number}`, 'GET', { 'Authorization': `${token}` }, null);
    if (res.ok) {
        const account = await res.json();
        return account;
    } else {
        alert('Unable to retrieve account information');
        const account = {};
        return account;
    }
}

export async function getImage(number) {
    const token = localStorage.getItem('token');
    const res = await callWithRefresh(`image/${number}`, 'GET', { 'Authorization': `${token}` }, null);
    var image = new Image();
    image.alt = 'Profile Picture';
    image.id = 'profile-pic';
    if (res.ok) {
        const resp = await res.json();
        const imageEnc = resp.image;
        image.src = 'data:image/png;base64,' + imageEnc;
        return image;
    } else {
        alert('Unable to retrieve image');
        return image;
    }
}

//In this version, logout does not make a call
//But usually, logout should be noticed by the server too somehow, so a call is required!
export function logout() {
    localStorage.removeItem('number');
    localStorage.removeItem('token');
    localStorage.removeItem('refresh_token');
    deleteAccount();
    navigateTo('/login');
    accountInactive();
}

export async function refreshAuth() {
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

export async function transfer(e) {
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
        navigateTo(`/account/${number}`);
    } else {
        alert('Transfer failed');
    }
}

export async function getElement(a, b) {
    const token = localStorage.getItem('token');
    const data = {
        a: a,
        b: b
    };
    const body = JSON.stringify(data);
    const res = await callWithRefresh(`element`, 'POST', { 'Authorization': `${token}` }, body);
    if (res.ok) {
        const resp = await res.json();
        console.log(resp);
        return resp.result;
    } else {
        return "Star";
    }
}

export async function register(e) {
    e.preventDefault();
    const form = e.target;
    const data = {
        first_name: form.first_name.value,
        last_name: form.last_name.value,
        password: form.password.value
    };

    const res = await fetch(`${API}/accounts`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });

    if (res.ok) {
        const account = await res.json();
        alert(`Registration successful! You account number was copied to your clipboard.`);
        navigator.clipboard.writeText(account.number);
        navigateTo('/login');
    } else {
        alert('Registration failed');
    }
}