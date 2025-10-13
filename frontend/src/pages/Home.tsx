import { useSuspenseQuery } from "@tanstack/react-query";
import { Suspense } from "react";
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
import { TrendingMovies } from "@/services/api/movies";

function Home() {
    return (
        <main>
            <div className="flex w-full flex-row">
                <TrendingTabbed />
            </div>
        </main>
    );
}

function TrendingTabbed() {
    const skeletons = Array.from({ length: 8 }, () => (
        <CarouselItem key={crypto.randomUUID()}>
            <Skeleton className="h-full w-40" />
        </CarouselItem>
    ));

    return (
        <Tabs className="w-full" defaultValue="Daily">
            <div className="mb-3 flex flex-row items-center gap-5">
                <h2 className="font-bold text-4xl">Trending</h2>
                <TabsList>
                    <TabsTrigger className="cursor-pointer" value="Daily">
                        Daily
                    </TabsTrigger>
                    <TabsTrigger className="cursor-pointer" value="Weekly">
                        Weekly
                    </TabsTrigger>
                </TabsList>
            </div>
            <TabsContents>
                <TabsContent value="Daily">
                    <Carousel>
                        <Suspense fallback={skeletons}>
                            <Trending />
                        </Suspense>
                    </Carousel>
                </TabsContent>
                <TabsContent value="Weekly">
                    <Carousel>
                        <Suspense fallback={skeletons}>
                            <Trending duration="weekly" />
                        </Suspense>
                    </Carousel>
                </TabsContent>
            </TabsContents>
        </Tabs>
    );
}

function Trending({ duration = "daily" }: { duration?: "daily" | "weekly" }) {
    const { data } = useSuspenseQuery({
        queryKey: ["trenging", duration],
        queryFn: () => TrendingMovies(duration === "weekly"),
    });

    return data.map((movie) => (
        <CarouselItem key={movie.id}>
            <MovieCard movie={movie} />
        </CarouselItem>
    ));
}

export default Home;
