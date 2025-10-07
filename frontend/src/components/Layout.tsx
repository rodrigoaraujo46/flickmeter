import { Outlet } from "react-router";
import { Toaster } from "sonner";
import Navbar from "./Navbar";

export function Layout() {
    return (
        <>
            <Toaster />
            <Navbar />
            <main className="min-w-[320px] max-w-[1280px]">
                <Outlet />
            </main>
        </>
    );
}
