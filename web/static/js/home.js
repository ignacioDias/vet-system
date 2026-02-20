const form = document.getElementById("filtersForm");
const clearBtn = document.getElementById("clearFilters");
const consultationsContainer = document.querySelector(".consultations");

let currentPage = 1;
let totalPages = 1;

document.addEventListener("DOMContentLoaded", () => {
    fetchConsultations();
});

async function fetchConsultations(page = 1) {
    const formData = new FormData(form);

    const params = new URLSearchParams();

    if (formData.get("is_completed")) {
        params.append("is_completed", formData.get("is_completed"));
    }

    if (formData.get("dni")) {
        params.append("dni", formData.get("dni"));
    }

    if (formData.get("order")) {
        params.append("order", formData.get("order"));
    }

    params.append("page", page);

    try {
        const response = await fetch(`/api/consultations?${params.toString()}`);
        if (!response.ok) {
            throw new Error("Error fetching consultations");
        }

        const data = await response.json();

        currentPage = data.page;
        totalPages = data.total_pages;

        renderConsultations(data.data);
        renderPagination();

    } catch (error) {
        consultationsContainer.innerHTML = `<p>Error cargando consultas</p>`;
        console.error(error);
    }
}


function renderConsultations(consultations) {
    if (!consultations || consultations.length === 0) {
        consultationsContainer.innerHTML = "<p>No hay consultas.</p>";
        return;
    }

    let html = `
        <table>
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Dueño</th>
                    <th>Paciente</th>
                    <th>Motivo</th>
                    <th>Estado</th>
                    <th>Fecha</th>
                </tr>
            </thead>
            <tbody>
    `;

    consultations.forEach(c => {
        html += `
            <tr>
                <td>${c.id}</td>
                <td>${c.owner_dni || "-"}</td>
                <td>${c.patient_id}</td>
                <td>${c.reason}</td>
                <td>${c.is_completed ? "✔" : "Pendiente"}</td>
                <td>${formatDate(c.created_at)}</td>
            </tr>
        `;
    });

    html += `
            </tbody>
        </table>
    `;

    consultationsContainer.innerHTML = html;
}

function renderPagination() {
    const paginationHTML = `
        <div class="pagination">
            <button id="prevPage" ${currentPage <= 1 ? "disabled" : ""}>
                Anterior
            </button>

            <span>Página ${currentPage} de ${totalPages}</span>

            <button id="nextPage" ${currentPage >= totalPages ? "disabled" : ""}>
                Siguiente
            </button>
        </div>
    `;

    consultationsContainer.innerHTML += paginationHTML;

    document.getElementById("prevPage")?.addEventListener("click", () => {
        if (currentPage > 1) {
            fetchConsultations(currentPage - 1);
        }
    });

    document.getElementById("nextPage")?.addEventListener("click", () => {
        if (currentPage < totalPages) {
            fetchConsultations(currentPage + 1);
        }
    });
}

form.addEventListener("submit", (e) => {
    e.preventDefault();
    currentPage = 1;
    fetchConsultations(1);
});

// CLEAR FILTERS
clearBtn.addEventListener("click", () => {
    form.reset();
    currentPage = 1;
    fetchConsultations(1);
});

function formatDate(dateString) {
    if (!dateString) return "-";
    const date = new Date(dateString);
    return date.toLocaleDateString();
}