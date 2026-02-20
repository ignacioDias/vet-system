export async function checkAuth() {
    try {
        const response = await fetch("/api/users/me", {
            method: "GET",
            credentials: "include",
            headers: {
                "Accept": "application/json"
            }
        });

        if (response.status === 401) {
            return false;
        }

        if (!response.ok) {
            console.error("Unexpected status:", response.status);
            return false;
        }

        return true;

    } catch (err) {
        console.error("Auth check failed:", err);
        return false;
    }
}