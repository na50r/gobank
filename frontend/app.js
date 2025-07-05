const app = document.getElementById('app');
const API = 'http://localhost:3000';

function renderAccount(account = {}) {
    app.innerHTML = `
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
    app.innerHTML = `
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
    app.innerHTML = `
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
    app.innerHTML = `
    <h1>Make Transfer</h1>
    <form onsubmit="transfer(event)">
      <input name="to" placeholder="Number" required />
      <input name="amount" type="number" placeholder="Amount" required />
      <button class="btn" type="submit">Transfer</button>
    </form>
  `;
}

async function getAccount() {
    const number = Number(localStorage.getItem('number'));
    const token = localStorage.getItem('token');
    const account = await fetch(`${API}/account/${number}`, {
        headers: { 'x-jwt-token': `${token}` }
    }).then(r => r.json());
    renderAccount(account);
}

function logout() {
    localStorage.removeItem('number');
    localStorage.removeItem('token');
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
        renderAccount();
        const token = res.headers.get('x-jwt-token');
        localStorage.setItem('token', token);
        getAccount();
    } else {
        alert('Login failed');
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

    const res = await fetch(`${API}/transfer/${number}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'x-jwt-token': `${token}` },
        body: JSON.stringify(data)
    });
    if (res.ok) {
        alert('Transfer successful');
        getAccount();
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
