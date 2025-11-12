import { DropdownMenuTrigger } from "@radix-ui/react-dropdown-menu";
import React, { useEffect, useRef, useState } from "react";
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
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar";
import { LucideCircleUserRound } from "lucide-react";

function Navbar() {
    const [showNav, setShowNav] = useState(false);
    const lastScrollY = useRef(5000);

    useEffect(() => {
        const handleScroll = () => {
            setShowNav(
                window.scrollY < lastScrollY.current || window.scrollY <= 0,
            );
            lastScrollY.current = window.scrollY;
        };

        handleScroll();

        window.addEventListener("scroll", handleScroll);
        return () => window.removeEventListener("scroll", handleScroll);
    }, []);

    return (
        <header
            className={`sticky top-3 z-50 mx-4 mt-3 mb-6 flex h-[5.5rem] min-w-[320px] flex-row items-center rounded-full bg-foreground px-8 text-primary-foreground transition-transform duration-1000 ${showNav ? "translate-y-0" : "-translate-y-[200%]"}`}
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
                <DropdownMenuTrigger>
                    <Avatar className="h-12 w-12">
                        <AvatarImage
                            referrerPolicy="no-referrer"
                            src={data.avatar_url}
                            alt={data.username}
                        />
                        <AvatarFallback
                            asChild
                            className="bg-transparent text-background"
                        >
                            <LucideCircleUserRound />
                        </AvatarFallback>
                    </Avatar>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                    <DropdownMenuItem asChild>
                        <Button
                            variant="ghost"
                            className="w-full justify-start"
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
                        >
                            Log out
                        </Button>
                    </DropdownMenuItem>
                </DropdownMenuContent>
            </DropdownMenu>
        );
    }

    return (
        <AuthDialog
            trigger={
                <Button variant="secondary" className="rounded-full">
                    Sign in
                </Button>
            }
        />
    );
}

export function AuthDialog({ trigger }: { trigger: React.ReactNode }) {
    return (
        <Dialog>
            <DialogTrigger asChild>{trigger}</DialogTrigger>
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
