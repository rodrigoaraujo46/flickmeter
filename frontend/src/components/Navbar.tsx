import { DropdownMenuTrigger } from "@radix-ui/react-dropdown-menu";
import { useQuery } from "@tanstack/react-query";
import { LucideCircleUserRound, LucideSearch } from "lucide-react";
import type React from "react";
import { useEffect, useRef, useState } from "react";
import { Link } from "react-router";
import { toast } from "sonner";
import { useCurrentUserQuery } from "@/hooks/useCurrentUserQuery";
import { searchMovies } from "@/services/api/movies";
import AuthForm from "./AuthForm";
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
import MoviePoster from "./MoviePoster";
import { Skeleton } from "./Skeleton";
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Item, ItemContent, ItemGroup, ItemSeparator } from "./ui/item";

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
            <div className="ml-auto flex h-full flex-row items-center justify-center gap-4 text-xs">
                <MovieSearch />
                <UserMenu />
            </div>
        </header>
    );
}

function MovieSearch() {
    const [query, setQuery] = useState("");
    const [open, setOpen] = useState(false);
    const [debouncedQuery, setDebouncedQuery] = useState("");

    useEffect(() => {
        const timer = setTimeout(() => setDebouncedQuery(query), 300);
        return () => clearTimeout(timer);
    }, [query]);

    const { data: searchResults } = useQuery({
        queryKey: ["movies", "search", debouncedQuery],
        queryFn: () => searchMovies(debouncedQuery),
        enabled: !!debouncedQuery,
    });

    const ref = useRef<HTMLDivElement>(null);
    useEffect(() => {
        function handleClickOutside(e: MouseEvent) {
            if (ref.current && !ref.current.contains(e.target as Node)) {
                setOpen(false);
            }
        }
        document.addEventListener("mousedown", handleClickOutside);
        return () =>
            document.removeEventListener("mousedown", handleClickOutside);
    }, []);

    return (
        <div className="relative w-[30rem]">
            <div className="relative">
                <LucideSearch
                    size={14}
                    className="-translate-y-[55%] absolute top-1/2 left-3 text-muted-foreground"
                />
                <Input
                    type="text"
                    placeholder="Search movies..."
                    onChange={(e) => setQuery(e.target.value)}
                    onFocus={(e) => {
                        setOpen(true);
                        setQuery(e.target.value);
                    }}
                    className="rounded-full bg-popover pl-10 text-popover-foreground"
                />
            </div>
            {searchResults && searchResults.length > 0 && (
                <ItemGroup
                    ref={ref}
                    onClick={() => setOpen(false)}
                    onBlur={(e) => {
                        if (
                            !e.currentTarget.contains(e.relatedTarget as Node)
                        ) {
                            setOpen(false);
                        }
                    }}
                    className={`${open ? "visible" : "invisible"} absolute top-full z-20 mt-2 max-h-[35rem] w-full overflow-auto rounded-lg bg-popover text-popover-foreground shadow-md outline-0 focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50`}
                >
                    {searchResults.map((movie, i) => (
                        <>
                            {i > 0 && <ItemSeparator />}
                            <Item
                                className="m-[2px] p-1 focus-within:ring-[3px] focus-within:ring-ring/50 focus-visible:border-ring"
                                key={movie.id}
                            >
                                <ItemContent>
                                    <Link
                                        className="flex h-24 flex-row items-center gap-3 rounded outline-0"
                                        to={`/movies/${movie.id}`}
                                    >
                                        <MoviePoster movie={movie} />
                                        <p className="font-semibold text-lg">
                                            {movie.title}
                                        </p>
                                    </Link>
                                </ItemContent>
                            </Item>
                        </>
                    ))}
                </ItemGroup>
            )}
        </div>
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
