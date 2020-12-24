"use strict";
const form = document.querySelector('#frmLogin');
form.onsubmit = () => {
    const data = new FormData(form);
    const username = data.get('uname');
    const password = data.get('pword');
    console.log(username + ' ' + password);
    return false;
};
