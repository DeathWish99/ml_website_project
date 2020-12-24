const form : HTMLFormElement | any = document.querySelector('#frmLogin')

form.onsubmit = () => {
    const data = new FormData(form)

    const username = data.get('uname') as string
    const password = data.get('pword') as string

    console.log(username + ' ' + password)
    return false
}