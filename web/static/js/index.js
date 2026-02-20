const $RegisterButton = document.querySelector(".register");
const $LoginButton = document.querySelector(".login");

$RegisterButton.addEventListener("click", () => {
    window.location.href = '/register';
})

$LoginButton.addEventListener("click", () => {
    window.location.href = '/login';
})
