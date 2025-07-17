import AbstractView from "./AbstractView.js";
import { h1Comp, formComp, inputComp, btnComp, containerComp, btnBar } from "../components/Ui.js";
import { login } from "../util/Calls.js";
import { navigateTo } from "../main.js";


function renderLogin() {
    const container = containerComp();
    const h1 = h1Comp('Login');
    container.append(h1);
    const form = formComp(login);
    const input1 = inputComp('number', 'Number', 'number');
    const input2 = inputComp('password', 'Password', 'password');
    const btn1 = btnComp('Login', () => { }, 'submit');
    const btn2 = btnComp('Register', () => {navigateTo('/register')});
    const bar = btnBar([btn1, btn2]);
    form.append(input1, input2, bar);
    container.append(form);
    return container;
}

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Login");
    }

    async getHtml() {
        if (this.params.id) {
            console.log("Id: " + this.params.id);
        }
        return renderLogin();
    }
}