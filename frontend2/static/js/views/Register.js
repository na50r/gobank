import AbstractView from "./AbstractView.js";
import { API } from "../index.js";
import { navigateToURL } from "../index.js";
import { h1Comp, formComp, inputComp, btnComp } from "../components/ui.js";


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
        navigateToURL('#/login');
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
    form.append(input1, input2, input3, btn);
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