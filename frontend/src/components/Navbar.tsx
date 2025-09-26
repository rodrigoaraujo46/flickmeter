import { useEffect, useRef, useState } from "react";
import { Button } from "./ui/button";

function Navbar() {
    const [showNav, setShowNav] = useState(true);
    const lastScrollY = useRef(0);

    useEffect(() => {
        const handleScroll = () => {
            const currentScrollY = window.scrollY; // read directly
            const diff = currentScrollY - lastScrollY.current;

            setShowNav(diff <= 0);

            lastScrollY.current = currentScrollY;
        };

        window.addEventListener("scroll", handleScroll);
        return () => window.removeEventListener("scroll", handleScroll);
    }, []);

    return (
        <header
            className={`sticky top-0 flex h-20 flex-row items-center bg-primary p-6 transition-transform delay-200 duration-500 ${showNav ? "translate-y-0" : "-translate-y-full"}`}
        >
            <a href="/" className="h-full text-primary-foreground">
                <p>FLICKMETER</p>
            </a>
            <Button className="ml-auto" variant="secondary">
                Sign in
            </Button>
        </header>
    );
}

export default Navbar;
