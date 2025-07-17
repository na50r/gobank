import AbstractView from "./AbstractView.js";
import { register } from "../util/Calls.js";
import { navigateTo } from "../main.js";
import { h1Comp, formComp, inputComp, btnComp, containerComp, btnBar } from "../components/Ui.js";

function renderRegister() {
    console.log(import.meta.env.VITE_API_URL);
    const container = containerComp();
    const h1 = h1Comp('Register');
    container.append(h1);
    const form = formComp(register);
    const input1 = inputComp('first_name', 'First Name');
    const input2 = inputComp('last_name', 'Last Name');
    const input3 = inputComp('password', 'Password', 'password');
    const btn = btnComp('Register', () => { }, 'submit');
    const btn2 = btnComp('Login', () => { navigateTo('#/login') });
    const bar = btnBar([btn, btn2]);
    form.append(input1, input2, input3, bar);
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