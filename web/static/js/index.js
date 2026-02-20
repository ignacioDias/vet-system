import { checkAuth } from "./auth.js";

document.addEventListener("DOMContentLoaded", async () => {
    const isLogged = await checkAuth();
    if (isLogged) {
        window.location.href = "/home";
        return;
    }
})
const $RegisterButton = document.querySelector(".register");
const $LoginButton = document.querySelector(".login");

$RegisterButton.addEventListener("click", () => {
    window.location.href = '/register';
})

$LoginButton.addEventListener("click", () => {
    window.location.href = '/login';
})
