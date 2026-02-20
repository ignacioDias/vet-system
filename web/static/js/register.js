const $LoginButton = document.querySelector(".login-button")
const $form = document.querySelector(".form-register")

$LoginButton.addEventListener("click", () => {
    window.location.href = '/login';
})

$form.addEventListener("submit", async (event) => {
    event.preventDefault()

    const formData = new FormData($form)
    const data = {
        name: formData.get("name"),
        email: formData.get("email"),
        dni: formData.get("DNI"),
        password: formData.get("password"),
        profilePicture: ""
    }

    try {
        const response = await fetch("/api/users", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        })

        if (response.ok) {
            alert("Registrado correctamente, ahora inicie sesión.")
            window.location.href = "/"
        } else {
            const errorText = await response.text()
            
            switch (response.status) {
                case 400:
                    if (errorText.includes("Name is required")) {
                        alert("Error: El nombre es obligatorio")
                    } else if (errorText.includes("Email is required")) {
                        alert("Error: El email es obligatorio")
                    } else if (errorText.includes("Invalid email format")) {
                        alert("Error: El formato del email no es válido")
                    } else if (errorText.includes("DNI is required")) {
                        alert("Error: El DNI es obligatorio")
                    } else if (errorText.includes("Invalid DNI")) {
                        alert("Error: El DNI no está autorizado para registrarse en el sistema")
                    } else if (errorText.includes("Invalid password")) {
                        alert("Error: La contraseña no cumple con los requisitos:\n- Mínimo 8 caracteres, máximo 72\n- Al menos una mayúscula\n- Al menos una minúscula\n- Al menos un número\n- Al menos un símbolo especial (! @ # $ % ^ & *)")
                    } else {
                        alert("Error: Datos inválidos. Verifica la información ingresada")
                    }
                    break
                case 500:
                    alert("Error del servidor: No se pudo completar el registro. Intenta nuevamente más tarde")
                    break
                case 429:
                    alert("Error: Demasiados intentos. Espera un momento antes de intentar nuevamente")
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