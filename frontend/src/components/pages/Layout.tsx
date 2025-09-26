import { Outlet } from "react-router";
import Navbar from "../Navbar";

export function Layout() {
    return (
        <>
            <Navbar />
            <main className="pt-16">
                <Outlet />
            </main>
        </>
    );
}
