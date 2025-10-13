import { Outlet } from "react-router";
import { Toaster } from "sonner";
import Navbar from "./Navbar";

export function Layout() {
    return (
        <>
            <Toaster />
            <Navbar />
            <main>
                <Outlet />
            </main>
        </>
    );
}
