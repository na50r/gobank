// Basic Vanilla JS
// Everything is within one div
const body = document.body
const div = document.createElement('div')
body.append(div)

const API = 'http://localhost:3000';
const evtSource = new EventSource(`${API}/stream`);

evtSource.onmessage = function (event) {
    const data = JSON.parse(event.data);
    switch (data.type) {
        case "transaction":
            const number = Number(localStorage.getItem('number'));
            console.log(`[Transaction] ${data.data.sender} -> ${data.data.amount} -> ${data.data.recipient}`);
            sender = Number(data.data.sender);
            recipient = Number(data.data.recipient);
            if (sender === number || recipient === number) {
                getAccount();
            }
            break;
        case "chat":
            console.log(`[Chat] ${data.data.name}: ${data.data.msg}`);
            break;
        default:
            console.log("Unknown event type", data);
    }
};

function renderAccount(account = {}) {
    div.innerHTML = `
    <h1>Account Details</h1>
    <table>
        <tr>
            <td>Number</td>
            <td>${account.number || 'Loading...'}</td>
        </tr>
        <tr>
            <td>First Name</td>
            <td>${account.first_name ?? 'Loading...'}</td>
        </tr>
        <tr>
            <td>Last Name</td>
            <td>${account.last_name ?? 'Loading...'}</td>
        </tr>
        <tr>
            <td>Balance</td>
            <td>${account.balance ?? 'Loading...'}</td>
        </tr>
    </table>
    <button class="btn" onclick="renderTransfer()">Make Transfer</button>
    <button class="btn" onclick="logout()">Logout</button>
  `;
}

function renderLogin() {
    div.innerHTML = `
    <h1>Login</h1>
    <form onsubmit="login(event)">
      <input name="number" placeholder="Number" required />
      <input name="password" type="password" placeholder="Password" required />
      <button class="btn" type="submit">Login</button>
      <button class="btn" onclick="renderRegister()">Register</button>
    </form>
  `;
}

function renderRegister() {
    div.innerHTML = `
    <h1>Register</h1>
    <form onsubmit="register(event)">
      <input name="first_name" placeholder="First Name" required />
      <input name="last_name" placeholder="Last Name" required />
      <input name="password" type="password" placeholder="Password" required />
      <button class="btn" type="submit">Register</button>
    </form>
  `;
}

function renderTransfer() {
    div.innerHTML = `
    <h1>Make Transfer</h1>
    <form onsubmit="transfer(event)">
      <input name="to" placeholder="Number" required />
      <input name="amount" type="number" placeholder="Amount" required />
      <button class="btn" type="submit">Transfer</button>
    </form>
  `;
}

async function refreshAuth() {

    const refresh_token = localStorage.getItem('refresh_token');
    data = {
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
        getAccount();
    } else if (res.status === 401) {
        alert('Session expired');
        logout();
    }
    else {
        alert('Refresh failed');
        logout();
    }
}

async function getAccount() {
    const number = Number(localStorage.getItem('number'));
    const token = localStorage.getItem('token');
    const res = await callWithRefresh(`account/${number}`, 'GET', { 'Authorization': `${token}` }, null);
    if (res.ok) {
        const account = await res.json();
        renderAccount(account);
    } else {
        alert('Get account failed');
        renderAccount();
    }
}

function logout() {
    localStorage.removeItem('number');
    localStorage.removeItem('token');
    localStorage.removeItem('refresh_token');
    renderLogin();
}

async function register(e) {
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
        alert(`Registration successful! You account number is ${account.number}`);
        renderLogin();
    } else {
        alert('Registration failed');
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
        getAccount();
    } else {
        alert('Login failed');
    }
}

async function callWithRefresh(endpoint, method, headers, body) {
    async function call() {
        // Update token
        const token = localStorage.getItem('token');
        headers['Authorization'] = token;
        const res = await fetch(`${API}/${endpoint}`, {
            method: method,
            headers: headers,
            body: body
        });
        return res;
    }
    res = await call();
    if (res.ok) {
        return res;
    } else if (res.status === 401) {
        await refreshAuth();
        return await call();
    } else {
        return res;
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
    } else {
        alert('Transfer failed');
    }
}

window.onload = () => {
    if (localStorage.getItem('token')) {
        getAccount();
        return;
    }
    renderLogin();
};
