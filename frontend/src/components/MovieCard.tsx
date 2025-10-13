import { useState } from "react";
import { Link } from "react-router";
import { cn } from "@/lib/utils";
import type { Movie } from "@/services/api/movies";
import { Skeleton } from "./Skeleton";

interface Props {
    className?: string;
    movie: Movie;
}

function MovieCard({ movie, className }: Props) {
    const [error, setError] = useState(false);
    const date = new Date(movie.release_date);

    return error ? (
        <Skeleton className="h-full w-40 rounded-lg" />
    ) : (
        <Link
            className={cn(
                "group relative block h-full w-full overflow-hidden rounded-lg",
                className,
            )}
            to={`/movies/${movie.id}`}
        >
            <img
                className="h-full w-full object-cover transition-all"
                src={`https://image.tmdb.org/t/p/original/${movie.poster_path}`}
                onError={() => setError(true)}
                alt={`${movie.title} poster`}
            />
            <div className="pointer-events-none absolute bottom-0 flex h-30 w-full translate-y-full flex-col justify-between bg-black/85 p-4 font-bold text-white transition-transform group-hover:translate-y-0">
                <p className="line-clamp-2 text-md">
                    {movie.title || movie.original_title}
                </p>
                <div className="flex flex-row justify-between">
                    <p className="ellipsis">{date.getFullYear()}</p>
                    <p className="flex items-center">
                        <span className="mb-[2.5px]">⭐</span>
                        <span className="ml-2">
                            {movie.vote_average.toFixed(1)}
                        </span>
                    </p>
                </div>
            </div>
        </Link>
    );
}

export default MovieCard;
