import React from "react";
import { cn } from "@/lib/utils";
import type { Movie } from "@/services/api/movies";
import { Skeleton } from "./Skeleton";

interface Props {
    className?: string;
    movie: Movie;
}

function MovieCard({ movie, className }: Props) {
    const [error, setError] = React.useState(false);

    return (
        <div className={cn("h-full w-full", className)}>
            {error ? (
                <Skeleton className="h-full w-40 rounded-lg" />
            ) : (
                <img
                    className="h-full w-full rounded-lg object-cover"
                    src={`https://image.tmdb.org/t/p/original/${movie.poster_path}`}
                    onError={() => setError(true)}
                    alt={`${movie.title} poster`}
                />
            )}
        </div>
    );
}

export default MovieCard;
