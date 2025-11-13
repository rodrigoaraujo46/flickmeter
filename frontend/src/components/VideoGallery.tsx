import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { toast } from "sonner";
import { fetchVideos } from "@/services/api/movies";
import { Skeleton } from "./Skeleton";

export default function VideoGallery({ movieId }: { movieId: number }) {
    const [currentVideoI, setCurrentVideoI] = useState(0);
    const [frameLoading, setFrameLoading] = useState(true);

    const prevVideo = () => {
        if (!videos) return;
        setFrameLoading(true);
        setCurrentVideoI((prev) => (prev === 0 ? videos.length - 1 : prev - 1));
    };

    const nextVideo = () => {
        if (!videos) return;
        setFrameLoading(true);
        setCurrentVideoI((prev) => (prev + 1) % videos.length);
    };

    const {
        data: videos,
        isPending,
        error,
    } = useQuery({
        queryKey: ["movies", movieId, "videos"],
        queryFn: () => fetchVideos(movieId),
    });

    if (isPending) {
        return <Skeleton className="h-full w-full rounded-lg" />;
    }

    if (error) {
        toast.error("Failed to load videos.");
        return <Skeleton className="h-full w-full rounded-lg" />;
    }

    if (!videos) {
        return <Skeleton className="h-full w-full rounded-lg" />;
    }

    const currentVideo = videos[currentVideoI];

    return (
        <div className="relative flex w-full items-center overflow-hidden rounded-lg">
            <PreviousVideoButton prevVideo={prevVideo} />
            {frameLoading && (
                <Skeleton className="absolute inset-0 h-full w-full" />
            )}
            <iframe
                className="h-full w-full"
                key={currentVideo?.id}
                title={currentVideo?.name}
                onLoad={() => setFrameLoading(false)}
                src={`https://www.youtube.com/embed/${currentVideo?.key}?autoplay=1&mute=1`}
                allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                allowFullScreen
            />
            <NextVideoButton nextVideo={nextVideo} />
        </div>
    );
}

function PreviousVideoButton({ prevVideo }: { prevVideo: () => void }) {
    return (
        <button
            className="absolute h-24 w-14 rounded-r-lg bg-black text-white/50 opacity-90 hover:text-white"
            type="button"
            onClick={prevVideo}
        >
            <svg
                className="ml-3.5"
                fill="none"
                height="24"
                viewBox="0 0 24 24"
                width="24"
            >
                <title>Previous Video</title>
                <path
                    d="M4 4C3.73 4 3.48 4.10 3.29 4.29C3.10 4.48 3 4.73 3 5V19C3 19.26 3.10 19.51 3.29 19.70C3.48 19.89 3.73 20 4 20C4.26 20 4.51 19.89 4.70 19.70C4.89 19.51 5 19.26 5 19V5C5 4.73 4.89 4.48 4.70 4.29C4.51 4.10 4.26 4 4 4ZM18.95 4.23L6 12.00L18.95 19.77C19.15 19.89 19.39 19.96 19.63 19.96C19.87 19.97 20.10 19.91 20.31 19.79C20.52 19.67 20.69 19.50 20.81 19.29C20.93 19.09 21.00 18.85 21 18.61V5.38C20.99 5.14 20.93 4.91 20.81 4.70C20.69 4.50 20.52 4.33 20.31 4.21C20.10 4.09 19.87 4.03 19.63 4.03C19.39 4.04 19.15 4.10 18.95 4.23Z"
                    fill="currentColor"
                ></path>
            </svg>
        </button>
    );
}

function NextVideoButton({ nextVideo }: { nextVideo: () => void }) {
    return (
        <button
            className="absolute right-0 h-24 w-14 rounded-l-lg bg-black text-white/50 opacity-90 hover:text-white"
            type="button"
            onClick={nextVideo}
        >
            <svg className="ml-3.5" fill="none" viewBox="0 0 24 24" width="24">
                <title>Next Video</title>
                <path
                    d="M20 20C20.26 20 20.51 19.89 20.70 19.70C20.89 19.51 21 19.26 21 19V5C21 4.73 20.89 4.48 20.70 4.29C20.51 4.10 20.26 4 20 4C19.73 4 19.48 4.10 19.29 4.29C19.10 4.48 19 4.73 19 5V19C19 19.26 19.10 19.51 19.29 19.70C19.48 19.89 19.73 20 20 20ZM5.04 19.77L18 12L5.04 4.22C4.84 4.10 4.60 4.03 4.36 4.03C4.12 4.03 3.89 4.09 3.68 4.21C3.47 4.32 3.30 4.49 3.18 4.70C3.06 4.91 2.99 5.14 3 5.38V18.61C2.99 18.85 3.06 19.08 3.18 19.29C3.30 19.50 3.47 19.67 3.68 19.79C3.89 19.90 4.12 19.96 4.36 19.96C4.60 19.96 4.84 19.89 5.04 19.77Z"
                    fill="currentColor"
                ></path>
            </svg>
        </button>
    );
}

export function VideoGallerySkeleton() {
    return <Skeleton className="h-full w-full rounded-lg" />;
}
