const $RegisterButton = document.querySelector(".register-button")
const $form = document.querySelector(".form-login")

$RegisterButton.addEventListener("click", () => {
    window.location.href = '/register';
})

$form.addEventListener("submit", async (event) => {
    event.preventDefault()

    const formData = new FormData($form)
    const data = {
        dni: formData.get("DNI"),
        password: formData.get("password")
    }

    try {
        const response = await fetch("/api/auth/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        })

        if (response.ok) {
            window.location.href = "/home"
        } else {
            const errorText = await response.text()
            
            switch (response.status) {
                case 400:
                    alert("Error: Datos incompletos. Verifica que hayas ingresado DNI y contraseña")
                    break
                case 401:
                    alert("Error: DNI o contraseña incorrectos. Verifica tus credenciales")
                    break
                case 429:
                    alert("Error: Demasiados intentos de inicio de sesión. Espera un momento antes de intentar nuevamente")
                    break
                case 500:
                    alert("Error del servidor: No se pudo completar el inicio de sesión. Intenta nuevamente más tarde")
                    break
                default:
                    alert("Error inesperado: " + errorText)
            }
        }
    } catch (error) {
        alert("Error al conectar con el servidor. Verifica tu conexión a internet")
        console.error(error)
    }
})
