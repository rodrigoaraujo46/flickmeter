type User = {
    id: string;
    username: string;
    avatar_url: string;
};

async function fetchCurrentUser(): Promise<User | null> {
    const res = await fetch("/api/users/me", {
        credentials: "include",
    });

    switch (res.status) {
        case 200:
            return (await res.json()) as User;
        case 401:
            return null;
        case 500:
            throw new Error("Internal Server Error");
        default:
            throw new Error(`Unexpected response: ${res.status}`);
    }
}

async function logout() {
    const res = await fetch("/api/users/logout", {
        method: "POST",
        credentials: "include",
    });

    switch (res.status) {
        case 200:
            return;
        case 500:
            throw new Error("Internal Server Error");
        default:
            throw new Error(`Unexpected response: ${res.status}`);
    }
}

export { fetchCurrentUser, logout, type User };
