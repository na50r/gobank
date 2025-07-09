import AbstractView from "./AbstractView.js";
import { API } from "../index.js";
import { h1Comp, formComp, inputComp, btnComp } from "../components/Ui.js";


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
        alert(`Registration successful! You account number was copied to your clipboard.`);
        navigator.clipboard.writeText(account.number);
        location.hash = '#/login';
    } else {
        alert('Registration failed');
    }
}

function renderRegister() {
    const container = document.createElement('div');
    const h1 = h1Comp('Register');
    container.append(h1);
    const form = formComp(register);
    const input1 = inputComp('first_name', 'First Name');
    const input2 = inputComp('last_name', 'Last Name');
    const input3 = inputComp('password', 'Password', 'password');
    const btn = btnComp('Register', () => { }, 'submit');
    const btn2 = btnComp('Login', () => { location.hash = '#/login'; });
    form.append(input1, input2, input3, btn, btn2);
    container.append(form);
    return container;
}

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Register");
    }

    async getHtml() {
        return renderRegister();
    }
}