import { useQuery } from "@tanstack/react-query";
import React from "react";
import { Carousel, CarouselItem } from "@/components/Carousel";
import MovieCard from "@/components/MovieCard";
import { Skeleton } from "@/components/Skeleton";
import {
    Tabs,
    TabsContent,
    TabsContents,
    TabsList,
    TabsTrigger,
} from "@/components/ui/shadcn-io/tabs";
import { fetchTrendingMovies, type Movie } from "@/services/api/movies";

function Home() {
    return (
        <div className="flex w-full flex-row">
            <TrendingTabbed />
        </div>
    );
}

function TrendingTabbed() {
    const [activeTab, setActiveTab] = React.useState("Daily");
    return (
        <Tabs onValueChange={setActiveTab} value={activeTab} className="w-full">
            <div className="mb-3 flex flex-row items-center gap-5">
                <h2 className="font-bold text-4xl">Trending</h2>
                <TabsList>
                    <TabsTrigger value="Daily">Daily</TabsTrigger>
                    <TabsTrigger value="Weekly">Weekly</TabsTrigger>
                </TabsList>
            </div>
            <TabsContents>
                <TabsContent inert={activeTab !== "Daily"} value="Daily">
                    <Carousel>
                        <Trending />
                    </Carousel>
                </TabsContent>
                <TabsContent inert={activeTab !== "Weekly"} value="Weekly">
                    <Carousel>
                        <Trending duration="weekly" />
                    </Carousel>
                </TabsContent>
            </TabsContents>
        </Tabs>
    );
}

function Trending({ duration = "daily" }: { duration?: "daily" | "weekly" }) {
    const { data, isLoading } = useQuery({
        queryKey: ["trending", duration],
        queryFn: () => fetchTrendingMovies(duration === "weekly"),
    });

    return (
        <>
            {(isLoading ? Array.from<Movie>({ length: 20 }) : data)?.map(
                (movie, i) =>
                    isLoading ? (
                        <CarouselItem key={`${`${i}`}`}>
                            <Skeleton className="aspect-[2/3] h-full rounded-lg" />
                        </CarouselItem>
                    ) : (
                        <CarouselItem key={movie.id}>
                            <MovieCard movie={movie} />
                        </CarouselItem>
                    ),
            )}
        </>
    );
}

export default Home;
