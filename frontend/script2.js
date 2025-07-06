// Basic Vanilla JS
const body = document.body
const API = 'http://localhost:3000';
const evtSource = new EventSource(`${API}/stream`);

window.onload = () => {
    if (localStorage.getItem('token')) {
        getAccount();
        return;
    }
    renderLogin();
};

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
    container = document.createElement('div');
    h1 = document.createElement('h1');
    h1.innerText = 'Account Details';
    container.append(h1);
    table = document.createElement('table');
    info = [
        { name: 'Number', value: account.number },
        { name: 'First Name', value: account.first_name },
        { name: 'Last Name', value: account.last_name },
        { name: 'Balance', value: account.balance }
    ];
    info.forEach((item) => {
        row = document.createElement('tr');
        col1 = document.createElement('td');
        col2 = document.createElement('td');
        col1.innerText = item.name;
        col2.innerText = item.value ?? 'Loading...';
        row.append(col1, col2);
        table.append(row);
    });
    container.append(table);
    btn1 = document.createElement('button');
    btn1.innerText = 'Make Transfer';
    btn1.onclick = renderTransfer;
    btn1.classList.add('btn');
    btn2 = document.createElement('button');
    btn2.innerText = 'Logout';
    btn2.onclick = logout;
    btn2.classList.add('btn');
    container.append(btn1, btn2);
    body.replaceChild(container, body.firstChild);
}

function renderLogin() {
    const container = document.createElement('div');
    const h1 = document.createElement('h1');
    h1.innerText = 'Login';
    container.append(h1);
    const form = document.createElement('form');
    form.addEventListener('submit', login);
    const input1 = document.createElement('input');
    input1.name = 'number';
    input1.placeholder = 'Number';
    input1.required = true;
    const input2 = document.createElement('input');
    input2.name = 'password';
    input2.type = 'password';
    input2.placeholder = 'Password';
    input2.required = true;
    const btn1 = document.createElement('button');
    btn1.innerText = 'Login';
    btn1.type = 'submit';
    btn1.classList.add('btn');
    const btn2 = document.createElement('button');
    btn2.innerText = 'Register';
    btn2.onclick = renderRegister;
    btn2.classList.add('btn');
    form.append(input1, input2, btn1, btn2);
    container.append(form);
    body.replaceChild(container, body.firstChild);
}

function renderRegister() {
    const container = document.createElement('div');
    const h1 = document.createElement('h1');
    h1.innerText = 'Register';
    container.append(h1);
    const form = document.createElement('form');
    form.addEventListener('submit', register);
    const input1 = document.createElement('input');
    input1.name = 'first_name';
    input1.placeholder = 'First Name';
    input1.required = true;
    const input2 = document.createElement('input');
    input2.name = 'last_name';
    input2.placeholder = 'Last Name';
    input2.required = true;
    const input3 = document.createElement('input');
    input3.name = 'password';
    input3.type = 'password';
    input3.placeholder = 'Password';
    input3.required = true;
    const btn = document.createElement('button');
    btn.innerText = 'Register';
    btn.type = 'submit';
    btn.classList.add('btn');
    form.append(input1, input2, input3, btn);
    container.append(form);
    body.replaceChild(container, body.firstChild);
}

function renderTransfer() {
    const container = document.createElement('div');
    const h1 = document.createElement('h1');
    h1.innerText = 'Make Transfer';
    container.append(h1);
    const form = document.createElement('form');
    form.addEventListener('submit', transfer);
    const input1 = document.createElement('input');
    input1.name = 'to';
    input1.placeholder = 'Number';
    input1.required = true;
    const input2 = document.createElement('input');
    input2.name = 'amount';
    input2.type = 'number';
    input2.placeholder = 'Amount';
    input2.required = true;
    const btn = document.createElement('button');
    btn.innerText = 'Transfer';
    btn.type = 'submit';
    btn.classList.add('btn');
    form.append(input1, input2, btn);
    container.append(form);
    body.replaceChild(container, body.firstChild);
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
