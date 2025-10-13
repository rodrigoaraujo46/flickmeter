type Movie = {
    adult: boolean;
    backdrop_path: string;
    id: number;
    title: string;
    original_language: string;
    original_title: string;
    overview: string;
    poster_path: string;
    media_type: string;
    genre_ids: number[];
    popularity: number;
    release_date: string;
    video: boolean;
    vote_average: number;
    vote_count: number;
};

async function TrendingMovies(weekly: boolean): Promise<Movie[]> {
    const res = await fetch(`/api/movies/trending?weekly=${weekly}`);

    if (!res.ok) {
        const errData = await res.json();
        throw new Error(
            errData.error || `Request failed with status ${res.status}`,
        );
    }

    return (await res.json()) as Movie[];
}

export { TrendingMovies, type Movie };
