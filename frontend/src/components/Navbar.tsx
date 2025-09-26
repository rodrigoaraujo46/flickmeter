import { useEffect, useRef, useState } from "react";
import { Link } from "react-router";
import { Button } from "./ui/button";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "./ui/dialog";
import AuthForm from "./AuthForm";

function Navbar() {
    const [showNav, setShowNav] = useState(true);
    const lastScrollY = useRef(0);

    useEffect(() => {
        const handleScroll = () => {
            const currentScrollY = window.scrollY;
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
            <Link to="/" className="h-full text-primary-foreground">
                <p>FLICKMETER</p>
            </Link>
            <Dialog>
                <DialogTrigger asChild>
                    <Button className="ml-auto" variant="secondary">
                        Sign in
                    </Button>
                </DialogTrigger>
                <DialogContent className="sm:max-w-[425px]">
                    <DialogHeader>
                        <DialogTitle>Sign Up or Log In</DialogTitle>
                    </DialogHeader>
                    <AuthForm />
                </DialogContent>
            </Dialog>
        </header>
    );
}

export default Navbar;
