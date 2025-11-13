import type { User } from "./users";

interface Genre {
    id: number;
    name: string;
}

interface ProductionCompany {
    id: number;
    logo_path: string;
    name: string;
    origin_country: string;
}

interface ProductionCountry {
    iso_3166_1: string;
    name: string;
}

interface SpokenLanguage {
    english_name: string;
    iso_639_1: string;
    name: string;
}

interface BelongsToCollection {
    id: number;
    name: string;
    poster_path: string;
    backdrop_path: string;
}

interface Video {
    id: string;
    iso_639_1: string;
    iso_3166_1: string;
    name: string;
    key: string;
    site: string;
    size: number; // defaults to 0
    type: string;
    official: boolean; // defaults to true
    published_at: string;
}

interface Review {
    id: number;
    movie_id: number;
    user_id: number;
    title: string;
    rating: number;
    review: string;
    created_at: string;
    updated_at?: string;
    user: User;
}

interface Movie {
    adult: boolean;
    backdrop_path: string;
    belongs_to_collection?: BelongsToCollection;
    budget: number;
    genres?: Genre[];
    homepage: string;
    id: number;
    imdb_id: string;
    original_language: string;
    original_title: string;
    overview: string;
    popularity: number;
    poster_path: string;
    production_companies?: ProductionCompany[];
    production_countries?: ProductionCountry[];
    release_date: string;
    revenue: number;
    runtime: number;
    spoken_languages?: SpokenLanguage[];
    status: string;
    tagline: string;
    title: string;
    video: boolean;
    videos?: Video[];
    vote_average: number;
    vote_count: number;
}

async function fetchTrendingMovies(weekly: boolean): Promise<Movie[]> {
    const res = await fetch(`/api/movies/trending?weekly=${weekly}`);
    if (!res.ok) {
        const error = await res.json();
        error.cause = res.status;
        throw error;
    }

    return (await res.json()) as Movie[];
}

async function searchMovies(query: string): Promise<Movie[]> {
    const res = await fetch(`/api/movies/search?query=${query}`);
    if (!res.ok) {
        const error = await res.json();
        error.cause = res.status;
        throw error;
    }

    return (await res.json()) as Movie[];
}

async function fetchVideos(id: number): Promise<Video[]> {
    const res = await fetch(`/api/movies/${id}/videos`);
    if (!res.ok) {
        const error = await res.json();
        error.cause = res.status;
        throw error;
    }

    return (await res.json()) as Video[];
}

async function fetchMovie(id: number): Promise<Movie> {
    const res = await fetch(`/api/movies/${id}`);
    if (!res.ok) {
        const error = await res.json();
        error.cause = res.status;
        throw error;
    }

    return (await res.json()) as Movie;
}

async function fetchMyReview(id: number): Promise<Review | null> {
    const res = await fetch(`/api/movies/${id}/reviews/me`);

    switch (res.status) {
        case 200:
            return (await res.json()) as Review;
        case 404:
            return null;
        default: {
            const error = await res.json();
            error.cause = res.status;
            throw error;
        }
    }
}

async function fetchReviews(id: number, page: number): Promise<Review[]> {
    const res = await fetch(`/api/movies/${id}/reviews?page=${page}`);
    if (!res.ok) {
        const error = await res.json();
        error.cause = res.status;
        throw error;
    }

    const data = await res.json();
    return Array.isArray(data) ? data : [];
}

async function saveReview(movieId: number, reviewId: number, review: Review) {
    const method = reviewId === 0 ? "POST" : "PATCH";
    const url =
        method === "POST"
            ? `/api/movies/${movieId}/reviews`
            : `/api/movies/${movieId}/reviews/${reviewId}`;

    const res = await fetch(url, {
        method: method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(review),
    });

    if (!res.ok) throw new Error("Failed to update review");
}

async function deleteReview(review: Review) {
    const res = await fetch(
        `/api/movies/${review.movie_id}/reviews/${review.id}`,
        {
            method: "DELETE",
        },
    );

    if (!res.ok) throw new Error("Failed to delete review");
}

export {
    deleteReview,
    saveReview,
    fetchMyReview,
    fetchReviews,
    fetchVideos,
    searchMovies,
    fetchMovie,
    fetchTrendingMovies,
    type Review,
    type Movie,
    type Video,
};
