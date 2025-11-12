import { LucideStar } from "lucide-react";
import { Link } from "react-router";
import { cn } from "@/lib/utils";
import type { Movie } from "@/services/api/movies";
import MoviePoster from "./MoviePoster";

interface Props {
    className?: string;
    movie: Movie;
}

function MovieCard({ movie, className }: Props) {
    const date = new Date(movie.release_date);

    return (
        <Link
            className={cn(
                "group relative block h-full w-full overflow-hidden rounded-lg",
                className,
            )}
            to={`/movies/${movie.id}`}
        >
            <MoviePoster movie={movie} />
            <div className="pointer-events-none absolute bottom-0 flex h-30 w-full translate-y-full flex-col justify-between bg-black/85 p-4 font-bold text-white transition-transform duration-500 ease-out group-hover:translate-y-0 group-focus:translate-y-0">
                <p className="line-clamp-2 text-md">
                    {movie.title || movie.original_title}
                </p>
                <div className="flex flex-row justify-between">
                    <p className="ellipsis">{date.getFullYear()}</p>
                    <p className="flex items-center">
                        <LucideStar size={15} fill="gold" color="gold" />
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
