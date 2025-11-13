import { useMutation, useQueryClient } from "@tanstack/react-query";
import { LucideEdit, LucideExternalLink, LucideTrash2 } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { useCurrentUserQuery } from "@/hooks/useCurrentUserQuery";
import { deleteReview, type Review } from "@/services/api/movies";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "./Dialog";
import { getRatingColor } from "./helpers";
import ReviewForm from "./ReviewForm";
import { Avatar, AvatarImage } from "./ui/avatar";
import { Button } from "./ui/button";

export function ReviewDialog({ review }: { review: Review }) {
    const [open, setOpen] = useState(false);
    const [mode, setMode] = useState<"view" | "edit" | "delete">("view");

    return (
        <Dialog
            open={open}
            onOpenChange={(state) => {
                setOpen(state);
                if (!state) setTimeout(() => setMode("view"), 150);
            }}
        >
            <DialogTrigger asChild>
                <Button className="font-semibold" variant="link">
                    FULL REVIEW
                    <LucideExternalLink />
                </Button>
            </DialogTrigger>

            <DialogContent className="flex max-h-[80vh] flex-col sm:max-w-[40%]">
                {mode === "view" && (
                    <ReviewView
                        review={review}
                        onEdit={() => setMode("edit")}
                        onDelete={() => setMode("delete")}
                    />
                )}

                {mode === "edit" && (
                    <ReviewEdit
                        review={review}
                        onClose={() => setMode("view")}
                        onCancel={() => setMode("view")}
                    />
                )}

                {mode === "delete" && (
                    <ReviewDelete
                        review={review}
                        onCancel={() => setMode("view")}
                        onDeleted={() => setOpen(false)}
                    />
                )}
            </DialogContent>
        </Dialog>
    );
}

function ReviewView({
    review,
    onEdit,
    onDelete,
}: {
    review: Review;
    onEdit: () => void;
    onDelete: () => void;
}) {
    const { data: user } = useCurrentUserQuery();
    return (
        <>
            <DialogHeader>
                <DialogTitle>REVIEW</DialogTitle>
            </DialogHeader>

            <div className="flex flex-col overflow-hidden p-2">
                <div className="flex justify-between">
                    <div>
                        <Button
                            className="px-0 pb-8 font-semibold"
                            variant="link"
                        >
                            <Avatar className="w-auto">
                                <AvatarImage
                                    referrerPolicy="no-referrer"
                                    src={review.user.avatar_url}
                                    alt={review.user.username}
                                />
                            </Avatar>
                            <p>{review.user.username}</p>
                        </Button>
                        {review.updated_at && (
                            <p>
                                {new Date(review.updated_at).toLocaleDateString(
                                    undefined,
                                    {
                                        year: "numeric",
                                        month: "short",
                                        day: "numeric",
                                    },
                                )}
                            </p>
                        )}
                        <h3 className="mt-4 truncate font-bold text-2xl">
                            {review.title}
                        </h3>
                    </div>

                    <div
                        className={`${getRatingColor(
                            review.rating,
                        )} flex aspect-square h-20 items-center justify-center rounded-lg text-center font-bold text-5xl text-primary-foreground`}
                    >
                        <p>{review.rating}</p>
                    </div>
                </div>

                <div className="mt-2 max-h-[60vh] overflow-y-auto overflow-x-hidden">
                    <p className="whitespace-pre-wrap break-words">
                        {review.review}
                    </p>
                </div>

                {user && user.id === review.user.id && (
                    <div className="mt-4 ml-auto flex gap-2">
                        <Button className="font-semibold" onClick={onEdit}>
                            EDIT REVIEW <LucideEdit />
                        </Button>
                        <Button
                            variant="destructive"
                            className="font-semibold"
                            onClick={onDelete}
                        >
                            DELETE REVIEW <LucideTrash2 />
                        </Button>
                    </div>
                )}
            </div>
        </>
    );
}
function ReviewEdit({
    review,
    onClose,
    onCancel,
}: {
    review: Review;
    onClose: () => void;
    onCancel: () => void;
}) {
    return (
        <>
            <DialogHeader>
                <DialogTitle>Edit Review</DialogTitle>
            </DialogHeader>
            <ReviewForm
                review={review}
                movieId={review.movie_id}
                onSuccess={onClose}
                onCancel={onCancel}
            />
        </>
    );
}

function ReviewDelete({
    review,
    onCancel,
    onDeleted,
}: {
    review: Review;
    onCancel: () => void;
    onDeleted: () => void;
}) {
    const queryClient = useQueryClient();

    const deleteReviewMutation = useMutation({
        mutationFn: deleteReview,
        onSuccess: () => {
            queryClient.invalidateQueries({
                queryKey: ["movies", review.movie_id, "reviews", "me"],
            });
        },
        onError: () => {
            toast.error("Failed to delete this review. Try again.");
        },
    });

    return (
        <>
            <DialogHeader>
                <DialogTitle>Delete Review</DialogTitle>
            </DialogHeader>
            <div className="flex flex-col gap-4">
                <p>
                    Are you sure you want to delete this review? This can’t be
                    undone.
                </p>
                <div className="ml-auto flex gap-2">
                    <Button variant="outline" onClick={onCancel}>
                        Cancel
                    </Button>
                    <Button
                        variant="destructive"
                        onClick={() => {
                            deleteReviewMutation.mutate(review);
                            onDeleted();
                        }}
                    >
                        Delete
                    </Button>
                </div>
            </div>
        </>
    );
}

export function NewReviewDialog({ movieId }: { movieId: number }) {
    const [open, setOpen] = useState(false);

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button className="ml-2" variant="default">
                    <LucideEdit />
                    Write Review
                </Button>
            </DialogTrigger>
            <DialogContent className="flex max-h-[80vh] flex-col sm:max-w-[40%]">
                <DialogHeader>
                    <DialogTitle>Write Review</DialogTitle>
                </DialogHeader>
                <ReviewForm
                    onSuccess={() => {
                        setOpen(false);
                    }}
                    movieId={movieId}
                />
            </DialogContent>
        </Dialog>
    );
}
