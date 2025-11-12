import { useQuery } from "@tanstack/react-query";
import { LucideStar } from "lucide-react";
import { useParams } from "react-router";
import MoviePoster, { MoviePosterSkeleton } from "@/components/MoviePoster";
import { Reviews } from "@/components/Reviews";
import { Skeleton } from "@/components/Skeleton";
import VideoGallery, { VideoGallerySkeleton } from "@/components/VideoGallery";
import { fetchMovie } from "@/services/api/movies";

function minutesToTime(mins: number) {
    const hours = Math.floor(mins / 60);
    const minutes = mins % 60;

    const hoursS = hours > 0 ? `${hours}h ` : "";
    const minsS = minutes > 0 ? `${minutes}m` : "";

    return `${hoursS}${minsS}`;
}

function Movie() {
    const { id } = useParams();
    if (!id) throw "Not Found";

    const {
        data: movie,
        isPending,
        error,
    } = useQuery({
        queryKey: ["movies", id],
        queryFn: () => fetchMovie(id),
        throwOnError: true,
    });

    if (error) throw error;

    if (isPending) {
        return (
            <div className="flex flex-col">
                <Skeleton className="h-10 w-full" />
                <div className="mt-10 flex h-[40rem] flex-row gap-4">
                    <MoviePosterSkeleton />
                    <VideoGallerySkeleton />
                </div>
                <Skeleton className="mt-14 h-[50rem] w-full" />
            </div>
        );
    }

    return (
        <div className="flex flex-col gap-6">
            <div className="flex flex-row justify-between">
                <div className="flex flex-col gap-2">
                    <h1 className="font-extrabold text-4xl">
                        {movie.title || movie.original_title}
                    </h1>
                    <p>
                        {movie.release_date.split("-")[0]}
                        {" · "}
                        {minutesToTime(movie.runtime)}
                    </p>
                </div>
                <p className="flex items-center font-bold text-2xl">
                    <LucideStar fill="gold" color="gold" />
                    <span className="ml-2">
                        {movie.vote_average.toFixed(1)}
                    </span>
                </p>
            </div>
            <div className="flex h-[40rem] flex-row gap-4">
                <MoviePoster movie={movie} />
                <VideoGallery movieId={movie.id} />
            </div>
            <p className="px-1">
                {movie.genres?.map((genre) => (
                    <span
                        key={genre.id}
                        className="mr-2 mb-2 rounded-full bg-secondary px-3 py-1 font-semibold text-secondary-foreground text-sm"
                    >
                        {genre.name}
                    </span>
                ))}
            </p>
            <p className="px-2 text-lg">{movie.overview}</p>
            <Reviews movieId={movie.id} />
        </div>
    );
}

export default Movie;
