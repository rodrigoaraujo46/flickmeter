import { useQuery, useQueryClient } from "@tanstack/react-query";
import { LucideEdit } from "lucide-react";
import React from "react";
import { toast } from "sonner";
import { fetchMyReview, fetchReviews } from "@/services/api/movies";
import { Button } from "./Button";
import { AuthDialog } from "./Navbar";
import { ReviewCard } from "./ReviewCard";
import { NewReviewDialog } from "./ReviewDialog";
import { Skeleton } from "./Skeleton";
import {
    Pagination,
    PaginationContent,
    PaginationItem,
    PaginationLink,
    PaginationNext,
    PaginationPrevious,
} from "./ui/pagination";

export function Reviews({ movieId }: { movieId: number }) {
    return (
        <section>
            <UserReview movieId={movieId} />
            <ReviewPages movieId={movieId} />
        </section>
    );
}

function UserReview({ movieId }: { movieId: number }) {
    const {
        error,
        isPending,
        data: review,
    } = useQuery({
        queryKey: ["movies", movieId, "reviews", "me"],
        queryFn: () => fetchMyReview(movieId),
        retry(failureCount, error) {
            if (error.cause === 401) {
                return false;
            }
            return failureCount < 3;
        },
    });

    if (isPending) {
        return <Skeleton className="h-80 w-full" />;
    }

    if (error) {
        if (error.cause !== 401) {
            toast.error("Couldn't load your review");
            return null;
        } else {
            return (
                <AuthDialog
                    trigger={
                        <Button className="ml-2" variant="default">
                            <LucideEdit />
                            Write Review
                        </Button>
                    }
                />
            );
        }
    }

    return (
        <>
            {review ? (
                <div className="h-80">
                    <ReviewCard review={review} />
                </div>
            ) : (
                <NewReviewDialog movieId={movieId} />
            )}
        </>
    );
}

function ReviewPages({ movieId }: { movieId: number }) {
    const [page, setPage] = React.useState(1);
    const queryClient = useQueryClient();
    const maxReviewsPerPage = 10;

    const {
        data: reviews,
        error,
        isPending,
    } = useQuery({
        queryKey: ["movies", movieId, "reviews", page],
        queryFn: () => fetchReviews(movieId, page),
        staleTime: 5 * 60 * 1000,
    });

    const { data: userReview } = useQuery({
        queryKey: ["movies", movieId, "reviews", "me"],
        queryFn: () => fetchMyReview(movieId),
        enabled: false,
    });

    React.useEffect(() => {
        if (reviews?.length === maxReviewsPerPage) {
            queryClient.prefetchQuery({
                queryKey: ["movies", movieId, "reviews", page + 1],
                queryFn: () => fetchReviews(movieId, page + 1),
            });
        }
    }, [reviews, page, queryClient, movieId]);

    if (isPending) {
        return <Skeleton className="h-[50rem] w-full" />;
    }

    if (error) {
        toast.error("Couldn't load reviews");
        return;
    }

    return (
        <>
            <hr className="my-10" />
            <div className="mt-4 flex w-full flex-row flex-wrap gap-4">
                <div className="flex w-full flex-row flex-wrap gap-4">
                    {reviews
                        .filter((review) => review.id !== userReview?.id)
                        .map((review) => (
                            <div
                                key={review.id}
                                className="h-80 flex-[1_1_calc(50%-1rem)]"
                            >
                                <ReviewCard review={review} />
                            </div>
                        ))}
                </div>
                {reviews && (
                    <Pagination className="mt-4">
                        <PaginationContent>
                            {page > 1 && (
                                <PaginationItem>
                                    <PaginationPrevious
                                        onClick={() => setPage(page - 1)}
                                    />
                                </PaginationItem>
                            )}
                            <PaginationItem>
                                <PaginationLink isActive>{page}</PaginationLink>
                            </PaginationItem>
                            {reviews.length === maxReviewsPerPage && (
                                <PaginationItem>
                                    <PaginationNext
                                        onClick={() => setPage(page + 1)}
                                    />
                                </PaginationItem>
                            )}
                        </PaginationContent>
                    </Pagination>
                )}
            </div>
        </>
    );
}
