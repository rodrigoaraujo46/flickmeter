import { useState } from "react";
import type { Movie } from "@/services/api/movies";
import { Skeleton } from "./Skeleton";

export default function MoviePoster({ movie }: { movie: Movie }) {
    const [isLoading, setIsLoading] = useState(true);

    return (
        <div className="relative aspect-[2/3] h-full">
            {isLoading && <Skeleton className="absolute inset-0 rounded-lg" />}
            <img
                className={`h-full w-full rounded-lg object-cover ${
                    isLoading ? "hidden" : "block"
                }`}
                src={`https://image.tmdb.org/t/p/original/${movie.poster_path}`}
                onLoad={() => setIsLoading(false)}
                alt={`${movie.title} poster`}
            />
        </div>
    );
}

export function MoviePosterSkeleton() {
    return <Skeleton className="aspect-[2/3] h-full" />;
}
