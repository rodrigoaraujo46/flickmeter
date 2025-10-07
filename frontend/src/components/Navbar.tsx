import { DropdownMenuTrigger } from "@radix-ui/react-dropdown-menu";
import { useEffect, useRef, useState } from "react";
import { Link } from "react-router";
import { toast } from "sonner";
import { useCurrentUserQuery } from "@/hooks/useCurrentUserQuery";
import AuthForm from "./AuthForm";
import { Button } from "./Button";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "./Dialog";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
} from "./DropdownMenu";
import { Skeleton } from "./Skeleton";

function Navbar() {
    const [showNav, setShowNav] = useState(false);
    const lastScrollY = useRef(0);

    useEffect(() => {
        const handleScroll = () => {
            setShowNav(
                window.scrollY < lastScrollY.current || window.scrollY === 0,
            );
            lastScrollY.current = window.scrollY;
        };

        handleScroll();

        window.addEventListener("scroll", handleScroll);
        return () => window.removeEventListener("scroll", handleScroll);
    }, []);

    return (
        <header
            className={`sticky top-3 mx-4 mt-3 mb-6 flex h-[5.5rem] min-w-[320px] flex-row items-center rounded-full bg-primary px-8 text-primary-foreground transition-transform delay-200 duration-500 ${showNav ? "translate-y-0" : "-translate-y-full"}`}
        >
            <Link to="/">
                <p className="translate-y-1 font-extrabold font-logo text-4xl leading-none">
                    FLICKMETER
                </p>
            </Link>
            <div className="ml-auto flex h-full flex-row items-center justify-center text-xs">
                <UserMenu />
            </div>
        </header>
    );
}

function UserMenu() {
    const { data, error, isLoading, logoutMutation } = useCurrentUserQuery();
    if (error) {
        toast.error("Oops! We couldn’t load your info. ", {
            description: "Try again later",
        });
    }

    if (isLoading) {
        return <Skeleton className="h-12 w-12 rounded-full" />;
    }
    if (data) {
        return (
            <DropdownMenu>
                <DropdownMenuTrigger asChild>
                    <button
                        type="button"
                        className="cursor-pointer rounded-full"
                    >
                        <img
                            referrerPolicy="no-referrer"
                            className={`h-12 w-12 rounded-full object-cover`}
                            src={data.avatar_url}
                            alt="avatar"
                        />
                    </button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                    <DropdownMenuItem asChild className="h-full w-full">
                        <button
                            onClick={() => {
                                document.body.style.cursor = "wait";
                                logoutMutation.mutate(undefined, {
                                    onError: () =>
                                        toast.error(
                                            "Oops! We couldn’t log you out. ",
                                            { description: "Try again later" },
                                        ),
                                    onSettled: () => {
                                        document.body.style.cursor = "default";
                                    },
                                });
                            }}
                            type="button"
                            className="cursor-pointer"
                        >
                            Log out
                        </button>
                    </DropdownMenuItem>
                </DropdownMenuContent>
            </DropdownMenu>
        );
    }

    return (
        <Dialog>
            <DialogTrigger asChild>
                <Button variant="secondary" className="rounded-full">
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
    );
}

export default Navbar;
